package main

import (
    "errors"
    "fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func OperatorValidator(input string) bool {
    // return true if the operator given by the user is a valid, recognized operator

    // copied from the types.ts specification
    RECOGNIZED_OPERATORS := []string{
        "firstSample",
        "lastSample",
        "firstFill",
        "lastFill",
        "mean",
        "min",
        "max",
        "count",
        "ncount",
        "nth",
        "median",
        "std",
        "jitter",
        "ignoreflyers",
        "flyers",
        "variance",
        "popvariance",
        "kurtosis",
        "skewness",
        "raw",
        "last",
    }
    for _, entry := range RECOGNIZED_OPERATORS {
        if entry == input {
            return true
        }
    }
    return false
}

func (qm ArchiverQueryModel) IdentifyFunctionsByName(targetName string) []FunctionDescriptorQueryModel {
    // create a slice of the the FunctionDescrporQueryModels that have the type of name targetName in order
    response := make([]FunctionDescriptorQueryModel, 1, 1) 
    for _, entry := range qm.Functions {
        if entry.Def.Name == targetName {
            response = append(response, entry)
        }
    }
    return response
}

func CreateOperatorQuery(qm ArchiverQueryModel) (string, error) {
    // Create the Prefix in the query to specify the operator

    // Skip any unrecognized operators 
    if ! OperatorValidator(qm.Operator) {
        errMsg := fmt.Sprintf("%v is not a recognized operator", qm.Operator)
        log.DefaultLogger.Debug("Error parsing query", "message", errMsg)
        return "", errors.New(errMsg)
    }

    // No operators are necessary in this case
    if qm.Operator == "" || qm.Operator == "raw" {
        return "", nil
    }

    qm.IdentifyFunctionsByName("binInterval")
    return "", nil

}
