package main
// This file intended to handle the boilerplate "best-practices" setup recommended by the grafana-plugin-sdk's example file "sample-plugin.go". The implementation of the archive querying logic will be elsewhere.

import ( 
	"context"
	"encoding/json"
	// "math/rand"
	"net/http"
	"time"

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
    // create response struct
    response := backend.NewQueryDataResponse()
    // IMPLEMENT HERE
    for idx, q := range req.Queries {
        log.DefaultLogger.Debug("index:", idx)
        log.DefaultLogger.Debug("query:", q)

        res := td.query(ctx, q)

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
    Format string `json:"format"`
}

func (td *ArchiverDatasource) query(ctx context.Context, query backend.DataQuery) backend.DataResponse {
    log.DefaultLogger.Debug("Executing Query",     "query",               query)
    log.DefaultLogger.Debug("query.RefID",         "query.RefID",         query.RefID)
    log.DefaultLogger.Debug("query.QueryType",     "query.QueryType",     query.QueryType)
    log.DefaultLogger.Debug("query.MaxDataPoints", "query.MaxDataPoints", query.MaxDataPoints)
    log.DefaultLogger.Debug("query.Interval",      "query.Interval",      query.Interval)
    log.DefaultLogger.Debug("query.TimeRange",     "query.TimeRange",     query.TimeRange)
    log.DefaultLogger.Debug("query.JSON",          "query.JSON",          query.JSON)


    // Unmarshal the json into our queryModel
    var qm archiverQueryModel

    response := backend.DataResponse{}

    response.Error = json.Unmarshal(query.JSON, &qm)
    if response.Error != nil {
        return response
    }

    // Log a warning if 'Format' is empty
    if qm.Format == "" {
        log.DefaultLogger.Warn("format is empty. defaulting to time series")
    }

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
