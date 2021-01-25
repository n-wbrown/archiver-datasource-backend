package main

import (
    "fmt"
    "testing"
)

func TestBuildQueryUrl(t *testing.T) {
}

func TestArchiverSingleQuery(t *testing.T) {
}

func TestArchiverSingleQueryParser(t *testing.T) {
}

func TestBuildRegexUrl(t *testing.T) {
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
