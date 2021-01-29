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
        {
            input: ArchiverQueryModel{
                Functions: []FunctionDescriptorQueryModel{
                    {
                        Def: FuncDefQueryModel{
                            Category: "Options",
                            DefaultParams: InitRawMsg(`[16]`),
                            Name: "binInterval",
                            Params:[]FuncDefParamQueryModel{
                                {Name:"interval", Type: "int"},
                            },
                        },
                        Params: []string{"[16]",},
                    },
                },
                Operator: "mean",
            },
            inputString: "binInterval",
            output: []FunctionDescriptorQueryModel{{Def: FuncDefQueryModel{Name: "binInterval"}}},
        },
    }
    for tdx, testCase := range tests {
        testName := fmt.Sprintf("%d: %v, %v", tdx, testCase.input, testCase.output)
        t.Run(testName, func(t *testing.T) {
            result := testCase.input.IdentifyFunctionsByName(testCase.inputString)
            if len(result) != len(testCase.output) {
                t.Errorf("lengths differ: got %v, want %v", len(result), len(testCase.output))
            }
            fmt.Println(result)
            for idx, out := range result {
                if out.Def.Name != testCase.output[idx].Def.Name {
                    t.Errorf("got %v, want %v", out, result[idx])
                }
            }
        })
    }
}

func TestCreateOperatorQuery(t *testing.T) {
    var tests = []struct{
        input ArchiverQueryModel
        output string
    }{
        {
            input: ArchiverQueryModel{
                Functions: []FunctionDescriptorQueryModel{
                    {
                        Def: FuncDefQueryModel{
                            Category: "Options",
                            DefaultParams: InitRawMsg(`[16]`),
                            Name: "binInterval",
                            Params:[]FuncDefParamQueryModel{
                                {Name:"interval", Type: "int"},
                            },
                        },
                        Params: []string{"[16]",},
                    },
                },
                Operator: "mean",
            },
            output: "mean_16",
        },
        {
            input: ArchiverQueryModel{
                Functions: []FunctionDescriptorQueryModel{
                    {
                        Def: FuncDefQueryModel{
                            Category: "Options",
                            DefaultParams: InitRawMsg(`[16]`),
                            Name: "binInterval",
                            Params:[]FuncDefParamQueryModel{
                                {Name:"interval", Type: "int"},
                            },
                        },
                        Params: []string{"[16]",},
                    },
                },
                Operator: "raw",
            },
            output: "",
        },
    }
    for idx, testCase := range tests {
        testName := fmt.Sprintf("%d: %v, %v", idx, testCase.input, testCase.output)
        t.Run(testName, func(t *testing.T) {
            result, err := CreateOperatorQuery(testCase.input)
            if err != nil {
                t.Errorf("Error received %v", err)
            }
            if testCase.output != result {
                t.Errorf("got %v, want %v", result, testCase.output)
            }
        })
    }
}
