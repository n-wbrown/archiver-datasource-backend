package main
// This file intended to handle the boilerplate "best-practices" setup recommended by the grafana-plugin-sdk's example file "sample-plugin.go". The implementation of the archive querying logic will be elsewhere.

import ( 
	"context"
	"encoding/json"
	// "math/rand"
	"net/http"
    //"reflect"

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
    var qm ArchiverQueryModel

    response := backend.DataResponse{}

    response.Error = json.Unmarshal(query.JSON, &qm)
    if response.Error != nil {
        return response
    }
    log.DefaultLogger.Debug("query.JSON unmarshalled", "qm", qm)
    log.DefaultLogger.Debug("Functions", "f", qm.Functions)
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

    // Apply Functions to the data
    // ApplyFunctions(responseData)

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

