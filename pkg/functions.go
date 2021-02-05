package main

import (
    "fmt"
	//"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func Sample(input string) bool {
    fmt.Println(input)
    return true
}

// Utilities 

// Transform functions
func Scale(allData []SingleData, factor float64) []SingleData {
    newData := make([]SingleData, len(allData))
    for ddx, oneData := range allData {
        newValues := make([]float64, len(oneData.Values))
        for idx, val := range oneData.Values {
            newValues[idx] = val * factor
        }
        newSd := SingleData{
            Times: oneData.Times,
            Values: newValues,
        }
        newData[ddx] = newSd
    }
    return newData
}

func Offset(allData []SingleData, delta float64) []SingleData {
    newData := make([]SingleData, len(allData))
    for ddx, oneData := range allData {
        newValues := make([]float64, len(oneData.Values))
        for idx, val := range oneData.Values {
            newValues[idx] = val + delta
        }
        newSd := SingleData{
            Times: oneData.Times,
            Values: newValues,
        }
        newData[ddx] = newSd
    }
    return newData
}

// Array to Scalar Functions

// Filter Series Functions 

// Sort Functions 
