package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Verbose is used to log verbose commands
var Verbose = log.New(ioutil.Discard, "", 0)

// Fail will log the error and quit the application
func Fail(err error) {
	fmt.Println(err)
	os.Exit(1)
}
