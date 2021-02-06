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
    var tests = []struct{
        inputSd []SingleData
        output []SingleData
    }{
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
            output: []SingleData{
                {
                    Times: TimeArrayHelper(1,6),
                    Values: []float64{0,1,1,2,3},
                },
            },
        },
    }
    for tdx, testCase := range tests {
        testName := fmt.Sprintf("case: %d", tdx)
        t.Run(testName, func(t *testing.T) {
            result := Delta(testCase.inputSd)
            SingleDataCompareHelper(result, testCase.output, t)
        })
    }
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
    var tests = []struct{
        inputSd []SingleData
        number int
        value string
        output []SingleData
    }{
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "avg",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "min",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,81},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "max",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,81},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{0,10,20,30,50,80},
                },
            },
            number: 1,
            value: "absoluteMin",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "absoluteMax",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10},
                },
            },
            number: 1,
            value: "sum",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
        },
    }
    for tdx, testCase := range tests {
        testName := fmt.Sprintf("case %d: %v", tdx, testCase.value)
        t.Run(testName, func(t *testing.T) {
            result, err := Top(testCase.inputSd, testCase.number, testCase.value)
            if err != nil {
                t.Errorf("Error not expected %v", err)
            }
            SingleDataCompareHelper(result, testCase.output, t)
        })
    }
}

func TestBottom(t *testing.T) {
    var tests = []struct{
        inputSd []SingleData
        number int
        value string
        output []SingleData
    }{
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "avg",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "min",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "max",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{0,10,20,30,50,80},
                },
            },
            number: 1,
            value: "absoluteMin",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{0,10,20,30,50,80},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10,10,20,30,50,80},
                },
            },
            number: 1,
            value: "absoluteMax",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
            },
        },
        {
            inputSd: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{1,1,2,3,5,8},
                },
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10},
                },
            },
            number: 1,
            value: "sum",
            output: []SingleData{
                {
                    Times: TimeArrayHelper(0,6),
                    Values: []float64{10},
                },
            },
        },
    }
    for tdx, testCase := range tests {
        testName := fmt.Sprintf("case %d: %v", tdx, testCase.value)
        t.Run(testName, func(t *testing.T) {
            result, err := Bottom(testCase.inputSd, testCase.number, testCase.value)
            if err != nil {
                t.Errorf("Error not expected %v", err)
            }
            SingleDataCompareHelper(result, testCase.output, t)
        })
    }
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
