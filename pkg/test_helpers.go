package main

import (
    "encoding/json"
    "time"
)

func MultiReturnHelperParseDuration(result time.Duration, err error) time.Duration {
    return result
}

func MultiReturnHelperParse(result time.Time, err error) time.Time {
    return result
}

func InitString(value string) *string{
    new_string := value
    return &new_string
}

func InitRawMsg(value string) *json.RawMessage{
    new_msg := json.RawMessage(value)
    return &new_msg
}
