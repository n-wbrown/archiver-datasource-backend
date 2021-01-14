package main
// This file intended to handle the boilerplate "best-practices" setup recommended by the grafana-plugin-sdk's example file "sample-plugin.go". The implementation of the archive querying logic will be elsewhere.

import ( 
	"context"
	"encoding/json"
	// "math/rand"
	"net/http"
    "net/url"
	"time"
    "reflect"
    "strings"
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
    log.DefaultLogger.Debug("Starting QueryData of newArchiverDataSource")
    log.DefaultLogger.Debug("QueryData", "request", req)
    log.DefaultLogger.Debug("QueryData.PluginContext", "PluginContext", req.PluginContext)
    log.DefaultLogger.Debug("QueryData.PluginContext type", "PluginContext_type", reflect.TypeOf(req.PluginContext))
    log.DefaultLogger.Debug("Plugintype.DataSourceInstanceSettings", "settings", req.PluginContext.DataSourceInstanceSettings.URL)

    // create response struct
    response := backend.NewQueryDataResponse()
    // IMPLEMENT HERE
    for idx, q := range req.Queries {
        log.DefaultLogger.Debug("index:", idx)
        log.DefaultLogger.Debug("query:", q)

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
    log.DefaultLogger.Debug("query.JSON unmarshalled", "qm", qm)
    log.DefaultLogger.Debug("qm.Target", "qm.Target", qm.Target)


    // let's extract all the relevant fields here:

    // data from query
    log.DefaultLogger.Debug("query.TimeRange.From",    "value",    query.TimeRange.From)
    log.DefaultLogger.Debug("query.TimeRange.To",      "value",    query.TimeRange.To)
    log.DefaultLogger.Debug("query.QueryType",         "value",    query.QueryType)
    log.DefaultLogger.Debug("query.MaxDataPoints",     "value",    query.MaxDataPoints)
    //log.DefaultLogger.Debug("query.Interval",        "value",    query.Interval)
    // data from unmarshaled JSON
    // log.DefaultLogger.Debug("qm.Datasource",        "value",    qm.Datasource)
    log.DefaultLogger.Debug("qm.Target",               "value",    qm.Target)
    log.DefaultLogger.Debug("qm.Regex",                "value",    qm.Regex)
    // log.DefaultLogger.Debug("qm.Operator",          "value",    qm.Operator)
    //log.DefaultLogger.Debug("qm.Functions",          "value",    qm.Functions)

    // data from original request's PluginContext
    log.DefaultLogger.Debug("pluginctx.DataSourceInstanceSettings.URL", "value",    pluginctx.DataSourceInstanceSettings.URL)

    // Log a warning if 'Format' is empty
    if qm.Format == "" {
        log.DefaultLogger.Warn("format is empty. defaulting to time series")
    }

    queryUrl := BuildQueryUrl(query, pluginctx, qm)
    archiverSingleQuery(queryUrl)

    // create data frame response
    frame := data.NewFrame("response")

    //add the time dimension
    frame.Fields = append(frame.Fields,
        data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
    )

    // add values 
    frame.Fields = append(frame.Fields, 
        data.NewField("values", nil, []int64{10, 20}),
    )

    // add the frames to the response
    response.Frames = append(response.Frames, frame)

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

func BuildQueryUrl(query backend.DataQuery, pluginctx backend.PluginContext, qm archiverQueryModel) string {
    // Build the URL to query the archiver built from Grafana's configuration
    // Set some constants
    TIME_FORMAT := "2006-01-02T15:04:05.000Z"
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

    // assemble the query of the URL and attach it to u
    query_vals :=  make(url.Values)
    query_vals["pv"] = []string{qm.Target} 
    query_vals["from"] = []string{query.TimeRange.From.Format(TIME_FORMAT)}
    query_vals["to"] = []string{query.TimeRange.To.Format(TIME_FORMAT)}
    query_vals["donotchunk"] = []string{""}
    u.RawQuery = query_vals.Encode()

    // Display the result
    log.DefaultLogger.Debug("u", "value", u)
    log.DefaultLogger.Debug("u.String", "value", u.String())
    return u.String()
}

type singleData struct {
   times []time.Time
   values []float64
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

func archiverSingleQuery( queryUrl string) singleData {
    // Take the unformatted response from the http GET request and turn it into rows of timeseries data
    var sD singleData

    // Make the GET request
    httpResponse, getErr := http.Get(queryUrl)
    if getErr != nil {
        log.DefaultLogger.Warn("Get request has failed", "Error", getErr)
        return sD
    }

    // Convert get request response to variable and close the file
    jsonAsBytes, ioErr := ioutil.ReadAll(httpResponse.Body)
    httpResponse.Body.Close()
    if ioErr != nil {
        log.DefaultLogger.Warn("Parsing of incoming data has failed", "Error", ioErr)
        return sD
    }

    log.DefaultLogger.Debug("Raw data", "value", string(jsonAsBytes))

    // Convert received data to JSON
    var data []ArchiverResponseModel
    jsonError := json.Unmarshal(jsonAsBytes, &data)
    if jsonError != nil {
        log.DefaultLogger.Warn("Conversion of incoming data to JSON has failed", "Error", jsonError)
        return sD
    }

    log.DefaultLogger.Debug("Data as JSON", "value", data)



    return sD
}
