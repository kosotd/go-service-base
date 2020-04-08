package cache

import (
	"gotest.tools/assert"
	"net/http"
	"testing"
	"time"
)

type responseWriterImpl struct {
	callCount int
}

func (w *responseWriterImpl) Header() http.Header {
	return nil
}

func (w *responseWriterImpl) Write([]byte) (int, error) {
	return 0, nil
}

func (w *responseWriterImpl) WriteHeader(statusCode int) {}

func TestCacheAndWriteJson(t *testing.T) {
	writer := responseWriterImpl{}
	getResp := func() (interface{}, error) {
		writer.callCount++
		return nil, nil
	}
	err := CacheAndWriteJson(&writer, "test_cache", getResp)
	assert.NilError(t, err)
	assert.Equal(t, writer.callCount, 1)

	err = CacheAndWriteJson(&writer, "test_cache", getResp)
	assert.NilError(t, err)
	assert.Equal(t, writer.callCount, 1)

	time.Sleep(3 * time.Second)
	err = CacheAndWriteJson(&writer, "test_cache", getResp)
	assert.NilError(t, err)
	assert.Equal(t, writer.callCount, 2)
}
