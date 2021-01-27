package main

import (
    "fmt"
    "testing"
)

func TestOperatorValidator(t *testing.T) {
    var tests = []struct{
        input string
        output bool
    }{
        {input: "firstSample",     output: true},
        {input: "lastFill",        output: true},
        {input: "lastFill_16",     output: false},
        {input: "snakes",          output: false},
    }
    for idx, testCase := range tests {
        testName := fmt.Sprintf("%d: %s, %t", idx, testCase.input, testCase.output)
        t.Run(testName, func(t *testing.T) {
            result := OperatorValidator(testCase.input)
            if testCase.output != result {
                t.Errorf("got %v, want %v", result, testCase.output)
            }
        })
    }
}

func TestIdentifyFunctionsByName(t *testing.T) {
    var tests = []struct{
        input ArchiverQueryModel
        inputString string
        output []FunctionDescriptorQueryModel
    }{
        {
            input: ArchiverQueryModel{Functions: []FunctionDescriptorQueryModel{{Def: FuncDefQueryModel{Name: "binInterval"}}}},
            inputString: "binInterval",
            output: []FunctionDescriptorQueryModel{{Def: FuncDefQueryModel{Name: "binInterval"}}},
        },
        {
            input: ArchiverQueryModel{Functions: []FunctionDescriptorQueryModel{{Def: FuncDefQueryModel{Name: "binInterval"}}}},
            inputString: "binInterval-fake",
            output: []FunctionDescriptorQueryModel{},
        },
    }
    for idx, testCase := range tests {
        testName := fmt.Sprintf("%d: %v, %v", idx, testCase.input, testCase.output)
        t.Run(testName, func(t *testing.T) {
            result := testCase.input.IdentifyFunctionsByName(testCase.inputString)
            for idx, out := range result {
                if out.Def.Name != result[idx].Def.Name {
                    t.Errorf("got %v, want %v", out, result[idx])
                }
            }
        })
    }

}
