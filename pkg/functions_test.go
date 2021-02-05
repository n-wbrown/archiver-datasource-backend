package main

import (
    "fmt"
    "testing"
)

func BlankTest(t *testing.T) {
    t.Skipf("Blank test")
}


// Utilites


// Transform funcitons

func TestScale(t *testing.T) {
    var tests = []struct{
        inputSd []SingleData
        delta float64
        output []SingleData
    }{
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
            delta: 2,
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{2,2,4,6,10,16},
                },
            },
        },
    }
    for tdx, testCase := range tests {
        testName := fmt.Sprintf("case %d: %v", tdx, testCase.output)
        t.Run(testName, func(t *testing.T) {
            result := Scale(testCase.inputSd, testCase.delta)
            SingleDataCompareHelper(result, testCase.output, t)
        })
    }
}

func TestOffset(t *testing.T) {
    var tests = []struct{
        inputSd []SingleData
        delta float64
        output []SingleData
    }{
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
            delta: 2,
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{3,3,4,5,7,10},
                },
            },
        },
    }
    for tdx, testCase := range tests {
        testName := fmt.Sprintf("case %d: %v", tdx, testCase.output)
        t.Run(testName, func(t *testing.T) {
            result := Offset(testCase.inputSd, testCase.delta)
            SingleDataCompareHelper(result, testCase.output, t)
        })
    }
}

func TestDelta(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestFluctuation(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestMovingAverage(t *testing.T) {
    t.Skipf("Not Implemeneted")
}


// Array to Scalar Functions

func TestToScalarByAvg(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestToScalarByMax(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestToScalarByMin(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestToScalarBySum(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestToScalarByMed(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestToScalarByStd(t *testing.T) {
    t.Skipf("Not Implemeneted")
}


// Filter Series Functions 

func TestTop(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestBottom(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestExclude(t *testing.T) {
    t.Skipf("Not Implemeneted")
}


// Sort Functions 

func TestSortByAvg(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestSortByMax(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestSortByMin(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestSortBySum(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestSortByAbsMax(t *testing.T) {
    t.Skipf("Not Implemeneted")
}

func TestSorByAbsMin(t *testing.T) {
    t.Skipf("Not Implemeneted")
}
