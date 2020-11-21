//+build mage

package main

import (
	"fmt"
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// Hello prints a message (shows that you can define custom Mage targets).
func Hello() {
	fmt.Println("hello plugin developer!")
}

// vtest() runs go test but uses the verbose option and does not stop on failure 
func Vtest() error {
	if err := sh.RunV("go", "test", "-v", "./pkg/..."); err != nil {
		return err
	}
	return nil
}

// Default configures the default target.
var Default = build.BuildAll
