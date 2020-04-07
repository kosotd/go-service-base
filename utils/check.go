package utils

import (
	"log"
)

func FailIfError(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}