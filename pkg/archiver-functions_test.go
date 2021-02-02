package main

import (
    "fmt"
    "testing"
)

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

func TestGetParametersByName(t *testing.T) {
    var tests = []struct{
        input FunctionDescriptorQueryModel
        targetArg string
        output *string
    }{
		{
		    input: FunctionDescriptorQueryModel{
		        Def: FuncDefQueryModel{
		            Fake: nil,
		            Category: "Options",
		            DefaultParams: InitRawMsg(`[900]`),
		            Name: "binInterval",
                    Params: []FuncDefParamQueryModel{
		                {
		                    Name: "interval",
		                    Options: nil,
		                    Type: "int",
		                },
		            },
		        },
		        Params:[]string{"900"},
		    },
            targetArg: "interval",
            output: InitString("900"),
		},
        {
		    input: FunctionDescriptorQueryModel{
		        Def: FuncDefQueryModel{
		            Fake: nil,
		            Category: "Filter Series",
		            DefaultParams: InitRawMsg(`[5 avg]`),
		            Name: "bottom",
		            Params: []FuncDefParamQueryModel{
		                {
		                    Name: "number",
		                    Options: nil,
		                    Type: "int",
		                },
		                {
		                    Name: "value",
		                    Options: &[]string{"avg", "min", "max", "absoluteMin", "absoluteMax" ,"sum"},
		                    Type: "string",
		                },
		            },
		        },
		        Params: []string{"5", "avg"},
		    },
            targetArg: "value",
            output: InitString("avg"),
        },
        {
		    input: FunctionDescriptorQueryModel{
		        Def: FuncDefQueryModel{
		            Fake: nil,
		            Category: "Transform",
		            DefaultParams: InitRawMsg(`[100]`),
		            Name: "offset",
		            Params: []FuncDefParamQueryModel{
		                {
		                    Name: "delta",
		                    Options: nil,
		                    Type: "float",
		                },
		            },
		        },
		        Params:[]string{"100"},
		    },
            targetArg: "delta",
            output: InitString("100"),
        },
        {
		    input: FunctionDescriptorQueryModel{
		        Def: FuncDefQueryModel{
		            Fake: nil,
		            Category: "Transform",
		            DefaultParams: InitRawMsg(`[]`),
		            Name: "delta",
		            Params: []FuncDefParamQueryModel{},
		        },
		        Params:[]string{},
		    },
            targetArg: "delta",
            output: nil,
		},
    }

    for tdx, testCase := range tests {
        testName := fmt.Sprintf("case %d: %v", tdx, testCase.output)
        t.Run(testName, func(t *testing.T) {
            result, err := testCase.input.GetParametersByName(testCase.targetArg)
            if testCase.output == nil {
                if err == nil {
                    t.Errorf("An error was expected but not received. Bad output: %v", result)
                }
                if result != "" {
                    t.Errorf("Expected output was \"\" but something else was received. Bad output: %v", result)
                }
            } else {
                if err != nil {
                    t.Errorf("An error was received but not expected. Bad output: %v", result)
                }
                if result != *testCase.output {
                    t.Errorf("Incorrect result: got %v, want %v", result, *testCase.output)
                }
            }
        })
    }
}

func TestApplyFunctions(t *testing.T) {
    t.Skipf("Test not implemented")
}

func TestFunctionSelector(t *testing.T) {
    t.Skipf("Test not implemented")
}



