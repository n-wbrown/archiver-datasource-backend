package main
// This file intended to handle the boilerplate "best-practices" setup recommended by the grafana-plugin-sdk's example file "sample-plugin.go". The implementation of the archive querying logic will be elsewhere.

import ( 
	"context"
	// "encoding/json"
	// "math/rand"
	"net/http"
	// "time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	// "github.com/grafana/grafana-plugin-sdk-go/backend/log"
	//"github.com/grafana/grafana-plugin-sdk-go/data"
)

func newArchiverDataSource() datasource.ServeOpts {
    // Create a new instance manager

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

    // create response struct
    response := backend.NewQueryDataResponse()

    // IMPLEMENT HERE

    return response, nil
}

func (td ArchiverDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
    var status = backend.HealthStatusOk
    var message = "This is a fake success message"

    return &backend.CheckHealthResult{
        Status:     status,
        Message:    message,
    }, nil
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
