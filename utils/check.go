package utils

import (
	"fmt"
	"log"
)

func FailIfError(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func FailIfNil(v interface{}, msg string) {
	if v == nil {
		log.Fatalf(fmt.Sprintf("variable is nil: %s", msg))
	}
}
