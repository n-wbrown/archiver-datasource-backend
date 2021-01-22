package main
// This file intended to handle the boilerplate "best-practices" setup recommended by the grafana-plugin-sdk's example file "sample-plugin.go". The implementation of the archive querying logic will be elsewhere.

import ( 
	"context"
	"encoding/json"
	// "math/rand"
	"net/http"
    "net/url"
	"time"
    //"reflect"
    "strings"
    "strconv"
    "io/ioutil"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func newArchiverDataSource() datasource.ServeOpts {
    // Create a new instance manager
    log.DefaultLogger.Debug("Starting newArchiverDataSource")

    im := datasource.NewInstanceManager(newArchiverDataSourceInstance)
    ds := &ArchiverDatasource{
        im: im,
    }

    return datasource.ServeOpts{
        QueryDataHandler:   ds,
        CheckHealthHandler: ds,
    }
}

type ArchiverDatasource struct {
    // Structure defined by grafana-plugin-sdk-go. Implements QueryData and CheckHealth.
    im instancemgmt.InstanceManager
}

func (td *ArchiverDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
    // Structure defined by grafana-plugin-sdk-go. QueryData should unpack the req argument into individual queries.
    //log.DefaultLogger.Debug("Starting QueryData of newArchiverDataSource")
    //log.DefaultLogger.Debug("QueryData", "request", req)
    //log.DefaultLogger.Debug("QueryData.PluginContext", "PluginContext", req.PluginContext)
    //log.DefaultLogger.Debug("QueryData.PluginContext type", "PluginContext_type", reflect.TypeOf(req.PluginContext))
    //log.DefaultLogger.Debug("Plugintype.DataSourceInstanceSettings", "settings", req.PluginContext.DataSourceInstanceSettings.URL)

    // create response struct
    response := backend.NewQueryDataResponse()
    // IMPLEMENT HERE
    for _, q := range req.Queries {
        // log.DefaultLogger.Debug("index:", idx)
        // log.DefaultLogger.Debug("query:", q)

        res := td.query(ctx, q, req.PluginContext)

        // save the response in a hashmap
        // based on with RefID as identifier
        response.Responses[q.RefID] = res
    }

    return response, nil
}

func (td *ArchiverDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
    var status = backend.HealthStatusOk
    var message = "This is a fake success message"

    return &backend.CheckHealthResult{
        Status:     status,
        Message:    message,
    }, nil
}

type archiverQueryModel struct {
    // It's not apparent to me where these two originate from but they do appear to be necessary
    Format string `json:"format"`
    Constant json.Number `json:"constant"` // I don't know what this is for yet
    QueryText string `json:"queryText"` // deprecated

    // Parameters added in AAQuery's extension of DataQuery
    Target string `json:"target"` //This will be the PV as entered by the user, or regex searching for PVs 
    Alias string `json:"alias"` // What to refer to the data as in the table - I think this only for the frontend rn
    AliasPattern string `json:"aliasPattern"` // use for collecting a large number of returned values 
    Operator string `json:"operator"` // ?
    Regex bool `json:"regex"` // configured by the user's setting of the "Regex" field in the panel
    Functions json.RawMessage `json:"functions"` // collection of functions to applied to the data by the archiver

    // Parameters from DataQuery
    RefId string `json:"refId"`
    Hide *bool `json:"hide"`
    Key *string `json:"string"`
    QueryType *string `json:"queryType"`
    DataTopic *string `json:"dataTopic"` //??
    Datasource *string `json:"datasource"` // comes back empty -- investigate further 
}

func (td *ArchiverDatasource) query(ctx context.Context, query backend.DataQuery, pluginctx backend.PluginContext) backend.DataResponse {
    // log.DefaultLogger.Debug("Executing Query",     "query",               query)
    // log.DefaultLogger.Debug("query.RefID",         "query.RefID",         query.RefID)
    // log.DefaultLogger.Debug("query.QueryType",     "query.QueryType",     query.QueryType)
    // log.DefaultLogger.Debug("query.MaxDataPoints", "query.MaxDataPoints", query.MaxDataPoints)
    // log.DefaultLogger.Debug("query.Interval",      "query.Interval",      query.Interval)
    // log.DefaultLogger.Debug("query.TimeRange",     "query.TimeRange",     query.TimeRange)
    // log.DefaultLogger.Debug("query.JSON",          "query.JSON",          query.JSON)
    // log.DefaultLogger.Debug("pluginctx",           "pluginctx",           pluginctx)


    // Unmarshal the json into our queryModel
    var qm archiverQueryModel

    response := backend.DataResponse{}

    response.Error = json.Unmarshal(query.JSON, &qm)
    if response.Error != nil {
        return response
    }
    // log.DefaultLogger.Debug("query.JSON unmarshalled", "qm", qm)
    // log.DefaultLogger.Debug("qm.Target", "qm.Target", qm.Target)


    // let's extract all the relevant fields here:

    // data from query
    // log.DefaultLogger.Debug("query.TimeRange.From",    "value",    query.TimeRange.From)
    // log.DefaultLogger.Debug("query.TimeRange.To",      "value",    query.TimeRange.To)
    // log.DefaultLogger.Debug("query.QueryType",         "value",    query.QueryType)
    // log.DefaultLogger.Debug("query.MaxDataPoints",     "value",    query.MaxDataPoints)
    // log.DefaultLogger.Debug("query.Interval",        "value",    query.Interval)
    // data from unmarshaled JSON
    // log.DefaultLogger.Debug("qm.Datasource",        "value",    qm.Datasource)
    // log.DefaultLogger.Debug("qm.Target",               "value",    qm.Target)
    // log.DefaultLogger.Debug("qm.Regex",                "value",    qm.Regex)
    // log.DefaultLogger.Debug("qm.Operator",          "value",    qm.Operator)
    // log.DefaultLogger.Debug("qm.Functions",          "value",    qm.Functions)

    // data from original request's PluginContext
    // log.DefaultLogger.Debug("pluginctx.DataSourceInstanceSettings.URL", "value",    pluginctx.DataSourceInstanceSettings.URL)

    // Log a warning if 'Format' is empty
    if qm.Format == "" {
        log.DefaultLogger.Warn("format is empty. defaulting to time series")
    }

    // make the query and compile the results into a SingleData instance
    responseData := make([]SingleData, 0)
    targetPvList := make([]string,0) 
    if qm.Regex {
        // If the user is using a regex to specify the PVs, parse and resolve the regex expression first

        // assemble the list of PVs to be queried for
        regexUrl := BuildRegexUrl(qm.Target, pluginctx)
        regexQueryResponse, _ := ArchiverRegexQuery(regexUrl)
        targetPvList, _ = ArchiverRegexQueryParser(regexQueryResponse)
    } else {
        // If a regex is not being used, only check for listed PVs
        targetPvList = IsolateBasicQuery(qm.Target)
    }

    // execute the individual queries
    for _, targetPv := range targetPvList {
        parsedResponse, _ := ExecuteSingleQuery(targetPv, query, pluginctx, qm)
        responseData = append(responseData, parsedResponse)
    }

    // for each query response, compile the data into response.Framse
    for _, singleResponse := range responseData {
        // create data frame response
        frame := data.NewFrame("response")

        //add the time dimension
        frame.Fields = append(frame.Fields,
            data.NewField("time", nil, singleResponse.Times),
        )

        // add values 
        frame.Fields = append(frame.Fields, 
            data.NewField("values", nil, singleResponse.Values),
        )

        // add the frames to the response
        response.Frames = append(response.Frames, frame)
    }

    return response
}

type archiverInstanceSettings struct {
	httpClient *http.Client
}

func newArchiverDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
    // Adheres to structure defined by grafana-plugin-sdk-go
    return &archiverInstanceSettings{
		httpClient: &http.Client{},
	}, nil
}

func BuildQueryUrl(target string, query backend.DataQuery, pluginctx backend.PluginContext, qm archiverQueryModel) string {
    // Build the URL to query the archiver built from Grafana's configuration
    // Set some constants
    TIME_FORMAT := "2006-01-02T15:04:05.000-07:00"
    JSON_DATA_URL := "data/getData.qw"

    // Unpack the configured URL for the datasource and use that as the base for assembling the query URL
    u, err := url.Parse(pluginctx.DataSourceInstanceSettings.URL)
        if err != nil {
        log.DefaultLogger.Warn("err", "err", err)
    }

    // amend the incomplete path
    var pathBuilder strings.Builder
    pathBuilder.WriteString(u.Path)
    pathBuilder.WriteString("/")
    pathBuilder.WriteString(JSON_DATA_URL)
    u.Path = pathBuilder.String()

    // location := time.Now().Location()
    // zone, offset := time.Now().Zone()
    // location, _ := time.Location(zone)
    // query.TimeRange.From.l

    //log.DefaultLogger.Debug("TimeRange.From TZ", "TZ", query.TimeRange.From.Location())
    //log.DefaultLogger.Debug("TimeRange.Format", "TZ", query.TimeRange.From.Format(TIME_FORMAT))
    //log.DefaultLogger.Debug("TimeRange.In", "TZ", query.TimeRange.From.In(location))
    //log.DefaultLogger.Debug("TimeRange.To TZ", "TZ", query.TimeRange.To.Location())

    /*
    new_time_from := time.Date(
        query.TimeRange.From.Year(),
        query.TimeRange.From.Month(),
        query.TimeRange.From.Day(),
        query.TimeRange.From.Hour(),
        query.TimeRange.From.Minute(),
        query.TimeRange.From.Second(),
        query.TimeRange.From.Nanosecond(),
        location)
    */


    //log.DefaultLogger.Debug("Local Timezone", "zone", zone)
    //log.DefaultLogger.Debug("Local Timezone", "offset", offset)

    // assemble the query of the URL and attach it to u
    query_vals :=  make(url.Values)
    query_vals["pv"] = []string{target} 
    query_vals["from"] = []string{query.TimeRange.From.Format(TIME_FORMAT)}
    query_vals["to"] = []string{query.TimeRange.To.Format(TIME_FORMAT)}
    query_vals["donotchunk"] = []string{""}
    u.RawQuery = query_vals.Encode()

    // Display the result
    // log.DefaultLogger.Debug("u", "value", u)
    // log.DefaultLogger.Debug("u.String", "value", u.String())
    return u.String()
}

type SingleData struct {
   Times []time.Time
   Values []float64
}

type ArchiverResponseModel struct {
    // Structure for unpacking the JSON response from the Archiver
    Meta struct{
        Name string `json:"name"`
        Waveform bool `json:"waveform"`
        EGU string `json:"EGU"`
        PREC json.Number `json:"PREC"`
    } `json:"meta"`
    Data []struct{
        Millis *json.Number`json:"millis,omitempty"`
        Nanos *json.Number`json:"nanos,omitempty"`
        Secs *json.Number`json:"secs,omitempty"`
        Val json.Number `json:"val"`
    } `json:"data"`
}

func ArchiverSingleQuery(queryUrl string) ([]byte, error){
    // Take the unformatted response from the http GET request and turn it into rows of timeseries data
    var jsonAsBytes []byte

    // Make the GET request
    httpResponse, getErr := http.Get(queryUrl)
    if getErr != nil {
        log.DefaultLogger.Warn("Get request has failed", "Error", getErr)
        return jsonAsBytes, getErr
    }

    // Convert get request response to variable and close the file
    jsonAsBytes, ioErr := ioutil.ReadAll(httpResponse.Body)
    httpResponse.Body.Close()
    if ioErr != nil {
        log.DefaultLogger.Warn("Parsing of incoming data has failed", "Error", ioErr)
        return jsonAsBytes, ioErr
    }

    // log.DefaultLogger.Debug("Raw data", "value", string(jsonAsBytes))
    return jsonAsBytes, nil
}

func ArchiverSingleQueryParser(jsonAsBytes []byte) (SingleData, error){
    // Convert received data to JSON
    var sD SingleData
    var data []ArchiverResponseModel
    jsonErr := json.Unmarshal(jsonAsBytes, &data)
    if jsonErr != nil {
        log.DefaultLogger.Warn("Conversion of incoming data to JSON has failed", "Error", jsonErr)
        return sD, jsonErr
    }

    // log.DefaultLogger.Debug("Data as JSON", "value", data)

    // Build output data block
    dataSize := len(data[0].Data)

    // initialize the slices with their final size so append operations are not necessary
    sD.Times = make([]time.Time, dataSize, dataSize)
    sD.Values = make([]float64, dataSize, dataSize)

    for idx, dataPt := range data[0].Data {

        millisCache, millisErr := dataPt.Millis.Int64()
        if millisErr != nil {
            log.DefaultLogger.Warn("Conversion of millis to int64 has failed", "Error", millisErr)
        }
        // use convert to nanoseconds
        sD.Times[idx] = time.Unix(0, 1e6 * millisCache)
        valCache, valErr := dataPt.Val.Float64()
        if valErr != nil {
            log.DefaultLogger.Warn("Conversion of val to float64 has failed", "Error", valErr)
        }
        sD.Values[idx] = valCache
    }
    // log.DefaultLogger.Debug("SingleData block", "Data", sD)
    return sD, nil
}

func BuildRegexUrl(regex string, pluginctx backend.PluginContext) string {
    // Construct the request URL for the regex search of PVs and return it as a string
    REGEX_URL := "bpl/getMatchingPVs"
    REGEX_MAXIMUM_MATCHES := 1000

    // Unpack the configured URL for the datasource and use that as the base for assembling the query URL
    u, err := url.Parse(pluginctx.DataSourceInstanceSettings.URL)
        if err != nil {
        log.DefaultLogger.Warn("err", "err", err)
    }

    // amend the incomplete path
    var pathBuilder strings.Builder
    pathBuilder.WriteString(u.Path)
    pathBuilder.WriteString("/")
    pathBuilder.WriteString(REGEX_URL)
    u.Path = pathBuilder.String()

    // assemble the query of the URL and attach it to u
    query_vals :=  make(url.Values)
    query_vals["regex"] = []string{regex}
    query_vals["limit"] = []string{strconv.Itoa(REGEX_MAXIMUM_MATCHES)}
    u.RawQuery = query_vals.Encode()

    // log.DefaultLogger.Debug("u.String", "value", u.String())
    return u.String()
}

func ArchiverRegexQuery(queryUrl string) ([]byte, error){
    // Make the GET request  for the JSON list of matching PVs, parse it, and return a list of strings
    var jsonAsBytes []byte

    // Make the GET request
    httpResponse, getErr := http.Get(queryUrl)
    if getErr != nil {
        log.DefaultLogger.Warn("Get request has failed", "Error", getErr)
        return jsonAsBytes, getErr
    }

    // Convert get request response to variable and close the file
    jsonAsBytes, ioErr := ioutil.ReadAll(httpResponse.Body)
    httpResponse.Body.Close()
    if ioErr != nil {
        log.DefaultLogger.Warn("Parsing of incoming data has failed", "Error", ioErr)
        return jsonAsBytes, ioErr
    }
    return jsonAsBytes, nil
}

func ArchiverRegexQueryParser(jsonAsBytes []byte) ([]string, error){
    // Convert received data to JSON
    var pvList []string
    jsonErr := json.Unmarshal(jsonAsBytes, &pvList)
    if jsonErr != nil {
        log.DefaultLogger.Warn("Conversion of incoming data to JSON has failed", "Error", jsonErr)
        return pvList, jsonErr
    }

    return pvList, nil
}

func ExecuteSingleQuery(target string, query backend.DataQuery, pluginctx backend.PluginContext, qm archiverQueryModel) (SingleData, error) {
    // wrap together the individual operations build a query, execute the query, and compile the data into a singleData structure
    // target: This is the PV to be queried for. As the "query" argument may be a regular expression, the specific PV desired must be specified
    queryUrl := BuildQueryUrl(target, query, pluginctx, qm)
    queryResponse, _ := ArchiverSingleQuery(queryUrl)
    parsedResponse, _ := ArchiverSingleQueryParser(queryResponse)
    return parsedResponse, nil
}

func IsolateBasicQuery(unparsed string) []string {
    // Non-regex queries can request multiple PVs using this syntax: (PV:NAME:1|PV:NAME:2|...)
    // This function takes queries in this format and breaks them up into a list of individual PVs
    unparsed_clean := strings.TrimSpace(unparsed)
    if unparsed_clean[0] != '(' || unparsed_clean[len(unparsed_clean)-1] != ')' {
        // if the statement doesn't have the parentheses, no parsing is necessary
        return []string{unparsed_clean}
    }
    // remove leading and following parentheses
    unparsed_clean = unparsed_clean[1:len(unparsed_clean)-1]
    result := strings.Split(unparsed_clean, "|")
    return result
}

