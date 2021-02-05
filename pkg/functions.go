package main

import (
    "fmt"
    "math"
    "errors"
    "sort"
	//"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)


// Utilities 

func FilterIndexer(allData []SingleData, value string) ([]float64, error) {
    rank := make([]float64, len(allData))
    for idx, sData := range allData {
        data := sData.Values
        switch value {
            case "avg":
                var total float64
                for _, val := range data {
                    total += val
                }
                rank[idx] = total/float64(len(data))
            case "min":
                var low_cache float64
                first_run := true
                for _, val := range data {
                    if first_run {
                        low_cache = val
                        first_run = false
                    }
                    if low_cache > val {
                        low_cache = val
                    }
                }
                rank[idx] = low_cache
            case "max":
                var high_cache float64
                first_run := true
                for _, val := range data {
                    if first_run {
                        high_cache = val
                        first_run = false
                    }
                    if high_cache < val {
                        high_cache = val
                    }
                }
                rank[idx] = high_cache
            case "absoluteMin":
                var low_cache float64
                first_run := true
                for _, originalVal := range data {
                    val := math.Abs(originalVal)
                    if first_run {
                        low_cache = val
                        first_run = false
                    }
                    if low_cache > val {
                        low_cache = val
                    }
                }
                rank[idx] = low_cache
            case "absoluteMax":
                var high_cache float64
                first_run := true
                for _, originalVal := range data {
                    val := math.Abs(originalVal)
                    if first_run {
                        high_cache = val
                        first_run = false
                    }
                    if high_cache < val {
                        high_cache = val
                    }
                }
                rank[idx] = high_cache
            case "sum":
                var total float64
                for _, val := range data {
                    total += val
                }
                rank[idx] = total
            default:
                errMsg := fmt.Sprintf("Value %v not recognized", value)
                return rank, errors.New(errMsg)
        }
    }
    return rank, nil
}



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
type SingleDataOrder struct {
    sD SingleData
    rank float64
}

func Bottom(allData []SingleData, number int, value string) ([]SingleData, error) {
    newData := make([]SingleData, 0, len(allData))
    rank, idxErr  := FilterIndexer(allData, value)
    if idxErr != nil {
        return allData, idxErr
    }
    if len(rank) != len(allData) {
        errMsg := fmt.Sprintf("Length of data (%v) and indexes (%v)differ", len(allData), len(rank))
        return allData, errors.New(errMsg)
    }
    order := make([]SingleDataOrder, len(allData))
    for idx, _ := range allData {
        order[idx] = SingleDataOrder{
            sD: allData[idx],
            rank: rank[idx],
        }
    }
    sort.SliceStable(order, func(i, j int) bool {
        return order[i].rank < order[j].rank
    })
    for idx, _ := range order {
        if idx >= number { break }
        newData = append(newData, order[idx].sD)
    }

    return newData, nil
}


// Sort Functions 
