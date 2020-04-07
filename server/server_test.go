package server

import (
	"github.com/kosotd/go-service-base/utils"
	"gotest.tools/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	AddHandler("GET", "/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("ok"))
		assert.NilError(t, err)
		writer.WriteHeader(http.StatusOK)
	})
	go RunServer()

	resp, err := http.Get("http://localhost:8081/")
	assert.NilError(t, err)
	defer utils.CloseSafe(resp.Body)

	bytes, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err)
	assert.Equal(t, string(bytes), "ok")
}
