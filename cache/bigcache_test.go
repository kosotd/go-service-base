package cache

import (
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestBigCache(t *testing.T) {
	err := SetData("test", []byte("data"))
	assert.NilError(t, err)
	MustSetData("test1", []byte("data1"))
	res, ok := GetData("test")
	assert.Equal(t, ok, true)
	assert.Equal(t, string(res), "data")
	res, ok = GetData("test1")
	assert.Equal(t, ok, true)
	assert.Equal(t, string(res), "data1")

	time.Sleep(3 * time.Second)
	res, ok = GetData("test")
	assert.Equal(t, ok, false)
	res, ok = GetData("test1")
	assert.Equal(t, ok, false)
}