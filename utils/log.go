package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
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

func LogErrorAndSetStatus(w http.ResponseWriter, code int, err error) {
	LogError("%v", err)
	http.Error(w, fmt.Sprintf("%s %v", time.Now().Format("2006/01/02 15:04:05"), errors.Cause(err)), code)
}

func LogAndSetStatusIfError(w http.ResponseWriter, code int, err error) {
	if err != nil {
		LogErrorAndSetStatus(w, code, err)
	}
}

func LogAndSetStatusIfRecover(w http.ResponseWriter, code int) {
	if r := recover(); r != nil {
		LogErrorAndSetStatus(w, code, errors.New(fmt.Sprint(r)))
	}
}
