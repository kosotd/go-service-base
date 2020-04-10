package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kosotd/go-service-base/utils"
	"github.com/kosotd/go-service-base/utils/domain"
	"gotest.tools/assert"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	AddHandler("GET", "/get1", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
		writer.WriteHeader(http.StatusOK)
	})

	AddGetHandler("/get2", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
		writer.WriteHeader(http.StatusOK)
	})

	AddPostHandler("/post1", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
		writer.WriteHeader(http.StatusOK)
	})

	AddGinHandler("POST", "/post2", func(context *gin.Context) {
		context.Data(http.StatusOK, "application/text", []byte("ok"))
	})

	AddGinGetHandler("/get", func(c *gin.Context) {
		c.Data(http.StatusAccepted, "application/text", []byte("ok"))
	})

	AddGinGetHandler("/get-timeout", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.Data(http.StatusOK, "application/text", []byte("ok"))
	})

	AddGinGetHandler("/get-path/:param", func(c *gin.Context) {
		param := c.Param("param")
		header := c.GetHeader("header")
		c.Data(http.StatusOK, "application/text", []byte(param+";"+header))
	})

	AddGinGetHandler("/get-query", func(c *gin.Context) {
		param := c.Query("param")
		c.Data(http.StatusOK, "application/text", []byte(param))
	})

	AddGinPostHandler("/post-urlencoded", func(c *gin.Context) {
		param := c.PostForm("param")
		c.Data(http.StatusOK, "application/text", []byte(param))
	})

	AddGinPostHandler("/post-json", func(c *gin.Context) {
		var body map[string]string
		err := json.NewDecoder(c.Request.Body).Decode(&body)
		assert.NilError(t, err)
		c.Data(http.StatusOK, "application/text", []byte(body["key1"]+";"+body["key2"]))
	})

	go RunServer()

	assert.Equal(t, string(utils.MustDoRequest(domain.Method{Url: "http://localhost:8081/get1"})), "ok")
	assert.Equal(t, string(utils.MustDoRequest(domain.Method{Url: "http://localhost:8081/get2"})), "ok")
	assert.Equal(t, string(utils.MustDoRequest(domain.Method{Method: "POST", Url: "http://localhost:8081/post1"})), "ok")
	assert.Equal(t, string(utils.MustDoRequest(domain.Method{Method: "POST", Url: "http://localhost:8081/post2"})), "ok")

	status, resp, err := utils.DoRequest(domain.Method{
		Url: "http://localhost:8081/get",
	})
	assert.NilError(t, err)
	assert.Equal(t, string(resp), "ok")
	assert.Equal(t, status, http.StatusAccepted)

	_, _, err = utils.DoRequest(domain.Method{
		Url:     "http://localhost:8081/get-timeout",
		Timeout: 1 * time.Second,
	})
	assert.Error(t, err, "utils.DoRequest -> client.Do: Get http://localhost:8081/get-timeout: net/http: request canceled (Client.Timeout exceeded while awaiting headers)")

	status, resp, err = utils.DoRequest(domain.Method{
		Url:        "http://localhost:8081/get-path/{1}",
		Method:     "GET",
		ParamsType: "path-params",
		Params:     map[string]string{"1": "param-value"},
		Headers:    map[string]string{"header": "header-value"},
	})
	assert.NilError(t, err)
	assert.Equal(t, string(resp), "param-value;header-value")
	assert.Equal(t, status, http.StatusOK)

	_, resp, err = utils.DoRequest(domain.Method{
		Url:        "http://localhost:8081/get-query",
		Method:     "GET",
		ParamsType: "query-params",
		Params:     map[string]string{"param": "query-param-value"},
	})
	assert.NilError(t, err)
	assert.Equal(t, string(resp), "query-param-value")

	_, resp, err = utils.DoRequest(domain.Method{
		Url:        "http://localhost:8081/post-urlencoded",
		Method:     "POST",
		ParamsType: "form-urlencoded",
		Params:     map[string]string{"param": "post-urlencoded-value"},
	})
	assert.NilError(t, err)
	assert.Equal(t, string(resp), "post-urlencoded-value")

	body := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	bytes, err := json.Marshal(body)
	assert.NilError(t, err)
	_, resp, err = utils.DoRequest(domain.Method{
		Url:        "http://localhost:8081/post-json",
		Method:     "POST",
		ParamsType: "json-body",
		Body:       bytes,
	})
	assert.NilError(t, err)
	assert.Equal(t, string(resp), "value1;value2")

	_, _, err = utils.DoRequest(domain.Method{
		ParamsType: "unknown",
	})
	assert.Error(t, err, "utils.DoRequest -> params_type.FromString: unknown params type")
}
