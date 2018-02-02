package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Verbose is used to log verbose commands
var Verbose = log.New(ioutil.Discard, "", 0)

// Fail will log the error and quit the application
func Fail(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// ValidateConfigs makes sure that all config options are set and if not
// it kills the program
func ValidateConfigs(keys ...string) error {
	for _, key := range keys {
		if viper.Get(key) == nil {
			return (errors.New("missing required config option: " + key))
		}
	}
	return nil
}
