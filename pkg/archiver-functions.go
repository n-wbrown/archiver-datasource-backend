package main

import (
    "errors"
    "fmt"
    "strconv"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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

func (fdqm FunctionDescriptorQueryModel) GetParamTypeByName (target string) (string, error) {
    // Provide the argument value for the function given its name.
    //  If multiple are received, only return the first. This should never happen. 
    if len(fdqm.Params) < len(fdqm.Def.Params) {
        errMsgLen := fmt.Sprintf("List of arguments exceeded the number of arguments provided (got %v wanted %v)", len(fdqm.Params), len(fdqm.Def.Params))
        return "", errors.New(errMsgLen)
    }
    for idx, def := range fdqm.Def.Params {
        if def.Name == target {
            return fdqm.Def.Params[idx].Type, nil
        }
    }
    errMsg := fmt.Sprintf("Not able to identify type of %v in function %v", target, fdqm.Def.Name)
    return "", errors.New(errMsg)
}

func (fdqm FunctionDescriptorQueryModel) ExtractorBase (target string, targetType string ) (string, error) {
    // confirm parameter type
    typeStr, getErr := fdqm.GetParamTypeByName(target)
    badReturn := ""
    fmt.Println(typeStr)
    if getErr != nil {
        errMsg := fmt.Sprintf("Failed to obtain parameter type for %v", target)
        log.DefaultLogger.Warn(errMsg)
        return badReturn, errors.New(errMsg)
    }
    if typeStr != targetType {
        // Warn but continue
        errMsg := fmt.Sprintf("Type %v not expected", typeStr)
        log.DefaultLogger.Warn(errMsg)
    }

    // attempt to locate and return the function's argument
    valueStr, getErr := fdqm.GetParametersByName(target)
    fmt.Println(valueStr)
    if getErr != nil {
        errMsg := fmt.Sprintf("Failed to obtain parameter %v", target)
        log.DefaultLogger.Warn(errMsg)
        return badReturn, errors.New(errMsg)
    }
    return valueStr, nil
}


func (fdqm FunctionDescriptorQueryModel) ExtractParamInt (target string) (int, error) {
    var result int

    // get string for argument's value and check type
    response, extractErr := fdqm.ExtractorBase(target, "int")
    if extractErr != nil {
        return result, extractErr
    }

    // convert to int
    result, conversionErr := strconv.Atoi(response)
    if conversionErr != nil {
        errMsg := fmt.Sprintf("Failed to convert %v to int", target)
        log.DefaultLogger.Warn(errMsg)
        return result, errors.New(errMsg)
    }

    return result, nil
}

func (fdqm FunctionDescriptorQueryModel) ExtractParamFloat64 (target string) (float64, error) {
    var result float64

    // get string for argument's value and check type
    response, extractErr := fdqm.ExtractorBase(target, "float")
    if extractErr != nil {
        return result, extractErr
    }

    // convert to int
    result, conversionErr := strconv.ParseFloat(response, 64)
    if conversionErr != nil {
        errMsg := fmt.Sprintf("Failed to convert %v to float64", target)
        log.DefaultLogger.Warn(errMsg)
        return result, errors.New(errMsg)
    }

    return result, nil
}

func (fdqm FunctionDescriptorQueryModel) ExtractParamString (target string) (string, error) {
    var result string

    // get string for argument's value and check type
    response, extractErr := fdqm.ExtractorBase(target, "string")
    if extractErr != nil {
        return result, extractErr
    }

    return response, nil
}

func ApplyFunctions(responseData []SingleData, qm ArchiverQueryModel) []SingleData {
    // iterate through the list of functions
    return responseData
}

func FunctionSelector(responseData []SingleData, fdqm FunctionDescriptorQueryModel) error{
    // Based on the name (as a string) of the function, select the actual function to be used
    // Note: This changes responseData inplace 
    name := fdqm.Def.Name
    // category := fdqm.Def.Category

    switch name {
        case "offset":
            delta := 3.3
            responseData = Offset(responseData, delta)
        default:
            errMsg := fmt.Sprintf("Function %v is not a recognized function", name)
            log.DefaultLogger.Warn(errMsg)
            return errors.New(errMsg)
    }
    return nil
}
