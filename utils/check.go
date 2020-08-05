package utils

import (
	"log"
)

func FailIfError(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func FailIfNil(v interface{}, format string, args ...interface{}) {
	if v == nil {
		log.Fatalf(format, args...)
	}
}

func Fail(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
