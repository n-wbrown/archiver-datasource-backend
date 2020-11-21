package main

import (
    "testing"
)

func TestWorks(t *testing.T) {
    // This function tests nothing
}

func TestFails(t *testing.T) {
    // This function and always fails
    t.Error("This is an example of test failure")
}
