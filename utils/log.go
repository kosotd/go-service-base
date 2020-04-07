package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

var ll = 0

func SetLogLevel(logLevel int) {
	ll = logLevel
}

func LogError(format string, v ...interface{}) {
	if ll >= 0 {
		log.Printf("ERROR: "+format, v...)
	}
}

func LogInfo(format string, v ...interface{}) {
	if ll >= 1 {
		log.Printf("INFO: "+format, v...)
	}
}

func LogDebug(format string, v ...interface{}) {
	if ll >= 2 {
		log.Printf("DEBUG: "+format, v...)
	}
}

func LogTrace(format string, v ...interface{}) {
	if ll >= 3 {
		log.Printf("TRACE: "+format, v...)
	}
}

func LogAndSetStatus(w http.ResponseWriter, code int, err error) {
	LogError(err.Error())
	http.Error(w, http.StatusText(code), code)
}

func LogAndSetStatusIfError(w http.ResponseWriter, code int, err error) {
	if err != nil {
		LogAndSetStatus(w, code, err)
	}
}

func LogAndSetStatusIfRecover(w http.ResponseWriter, code int) {
	if r := recover(); r != nil {
		LogAndSetStatus(w, code, errors.New(fmt.Sprint(r)))
	}
}
