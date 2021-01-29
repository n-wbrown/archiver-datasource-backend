package main

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
