package main

import ( 
	// "context"
	// "encoding/json"
	// "math/rand"
	// "net/http"
	// "time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func newArchiverDataSource() datasource.ServeOpts {
    // Create a new instance manager

    //
    im := datasource.NewInstanceManger(newArchiverDataSourceInstance)
    ds := &ArchiverDatasource{
        im: im,
    }
