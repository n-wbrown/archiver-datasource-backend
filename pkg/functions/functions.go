package functions

import (
    "fmt"
	//"github.com/grafana/grafana-plugin-sdk-go/backend/log"
    "github.com/n-wbrown/archiver-datasource-backend/pkg/archiver"
)

func Sample(input string) bool {
    fmt.Println(input)
    return true
}

// Utilities 

// Transform functions

func Offset(data []SingleData, delta float64) []SingleData {
    return data
}

// Array to Scalar Functions

// Filter Series Functions 

// Sort Functions 
