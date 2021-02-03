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

func Offset(data []archiver.SingleData, delta float64) []archiver.SingleData {
    newData := make([]archiver.SingleData, len(data))
    for idx, d := range data {
        newValues := make([]float64, len(data))
        for idx, val := range d.Values {
            newValues[idx] = d.Values + delta
        }
        newSd := archiver.SingleData{
            Times: data[idx].Times,
            Values: newValues,
        }
        newData[idx] = newSd
    }
    return data
}

// Array to Scalar Functions

// Filter Series Functions 

// Sort Functions 
