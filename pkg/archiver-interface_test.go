package main

import (
    "fmt"
    "math"
	"io/ioutil"
    "testing"
)

func TestBuildQueryUrl(t *testing.T) {
    t.Skipf("Test not implemented")
}

func TestArchiverSingleQuery(t *testing.T) {
    t.Skipf("Test not implemented")
}

func TestArchiverSingleQueryParser(t *testing.T) {
    ARCHIVER_FLOAT_PRECISION := 1e-18
    type responseParams struct{
        length int
        firstVal float64
        lastVal float64
    }

    var dataNames = []struct{
        fileName string
        output responseParams
    }{
        {fileName: "test_data/good_query_response_01.JSON", output: responseParams{length: 612, firstVal: 0.005249832756817341, lastVal: 0.005262143909931183}},
    }

    type testData struct {
        input []byte
        output responseParams
    }

    var tests []testData
    for _, entry := range dataNames {
        fileData, err := ioutil.ReadFile(entry.fileName)
        if err != nil {
            t.Fatalf("Failed to load test data: %v", err)
        }
        tests = append(tests, testData{input: fileData, output: entry.output})
    }

    for idx, testCase := range tests {
        testName := fmt.Sprintf("Case: %d", idx)
        t.Run(testName, func(t *testing.T) {
            // result := testCase.output
            result, err := ArchiverSingleQueryParser(testCase.input)
            if err != nil {
                t.Fatalf("An unexpected error has occurred")
            }
            if len(result.Times) != len(result.Values){
                t.Fatalf("Lengths of Times and Values differ - Times: %v Values: %v", len(result.Times), len(result.Values))
            }
            resultLength := len(result.Times)
            if resultLength != testCase.output.length {
                t.Fatalf("Lengths differ - Wanted: %v Got: %v", testCase.output.length, resultLength)
            }
            if math.Abs(result.Values[0] - testCase.output.firstVal) > ARCHIVER_FLOAT_PRECISION {
                t.Fatalf("First values differ - Wanted: %v Got: %v", testCase.output.firstVal, result.Values[0])
            }
            if math.Abs(result.Values[resultLength-1] - testCase.output.lastVal) > ARCHIVER_FLOAT_PRECISION {
                t.Fatalf("Last values differ - Wanted: %v Got: %v", testCase.output.lastVal, result.Values[resultLength-1])
            }
        })
    }

}

func TestBuildRegexUrl(t *testing.T) {
    t.Skipf("Test not implemented")
}

func TestArchiverRegexQueryParser(t *testing.T) {
    var tests = []struct{
        input []byte
        output []string
    }{
        {input: []byte("[\"MR1K1:BEND:PIP:1:PMON\",\"MR1K3:BEND:PIP:1:PMON\"]"), output: []string{"MR1K1:BEND:PIP:1:PMON","MR1K3:BEND:PIP:1:PMON"}},
        {input: []byte("[\"MR1K3:BEND:PIP:1:PMON\"]"), output: []string{"MR1K3:BEND:PIP:1:PMON"}},
        {input: []byte("[]"), output: []string{}},
    }

    for idx, testCase := range tests {
        testName := fmt.Sprintf("%d: %s, %s", idx, testCase.input, testCase.output)
        t.Run(testName, func(t *testing.T) {
            // result := testCase.output
            result, err := ArchiverRegexQueryParser(testCase.input)
            if err != nil {
                t.Fatalf("An unexpected error has occurred")
            }
            if len(result) != len(testCase.output) {
                t.Fatalf("Lengths differ - Wanted: %v Got: %v", testCase.output, result)
            }
            for idx, _ := range(testCase.output) {
                if testCase.output[idx] != result[idx] {
                    t.Errorf("got %v, want %v", result, testCase.output)
                }
            }
        })
    }
}

func TestIsolateBasicQuery(t *testing.T) {
    var tests = []struct{
        inputUnparsed string
        output []string
    }{
        {inputUnparsed: "(this:is:1|this:is:2)", output: []string{"this:is:1", "this:is:2"}},
        {inputUnparsed: "(this:is:1)", output: []string{"this:is:1"}},
        {inputUnparsed: "this:is:1", output: []string{"this:is:1"}},
    }

    for idx, testCase := range tests {
        testName := fmt.Sprintf("%d: %s, %s", idx, testCase.inputUnparsed, testCase.output)
        t.Run(testName, func(t *testing.T) {
            // result := testCase.output
            result := IsolateBasicQuery(testCase.inputUnparsed)
            if len(result) != len(testCase.output) {
                t.Fatalf("Lengths differ - Wanted: %v Got: %v", testCase.output, result)
            }
            for idx, _ := range(testCase.output) {
                if testCase.output[idx] != result[idx] {
                    t.Errorf("got %v, want %v", result, testCase.output)
                }
            }
        })
    }
}
