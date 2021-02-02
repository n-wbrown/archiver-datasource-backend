package main

import (
    "errors"
    "fmt"
    "github.com/n-wbrown/archiver-datasource-backend/pkg/functions"
	//"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (qm ArchiverQueryModel) IdentifyFunctionsByName(targetName string) []FunctionDescriptorQueryModel {
    // create a slice of the the FunctionDescrporQueryModels that have the type of name targetName in order
    response := make([]FunctionDescriptorQueryModel, 0, 0) 
    for _, entry := range qm.Functions {
        if entry.Def.Name == targetName {
            response = append(response, entry)
        }
    }
    return response
}

func (fdqm FunctionDescriptorQueryModel) GetParametersByName (target string) (string, error) {
    // Provide the argument value for the function given its name.
    //  If multiple are received, only return the first. This should never happen. 
    if len(fdqm.Params) < len(fdqm.Def.Params) {
        errMsgLen := fmt.Sprintf("List of arguments exceeded the number of arguments provided (got %v wanted %v)", len(fdqm.Params), len(fdqm.Def.Params))
        return "", errors.New(errMsgLen)
    }
    for idx, def := range fdqm.Def.Params {
        if def.Name == target {
            return fdqm.Params[idx], nil
        }
    }
    errMsg := fmt.Sprintf("Not able to identify argument %v in function %v", target, fdqm.Def.Name)
    return "", errors.New(errMsg)
}


func ApplyFunctions(responseData []SingleData) []SingleData {
    // iterate through the list of functions
    functions.Sample("yo")
    return responseData
}
/*
//  Call this one FunctionSelector again? 
func ApplyIndividualFunction(responseData []SingleData, fdqm FunctionDescriptorQueryModel) []SingleData{
    // Based on the name (as a string) of the function, select the actual function to be used 
}

one little func for each function, function will contain it's own parameter extractor and appear of the form:
func FunctionFunctionName(responseData []SingleData, fdqm FunctionDescriptorQueryModel) []SingleData {
}
*/
