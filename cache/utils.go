package cache

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func CacheAndWriteJson(w http.ResponseWriter, cacheName string, responseSupplier func() (interface{}, error)) error {
	var data []byte
	var has bool
	if data, has = GetData(cacheName); !has {
		resp, err := responseSupplier()
		if err != nil {
			return errors.Wrapf(err, "cache.CacheAndWriteJson -> responseSupplier")
		}
		var buff bytes.Buffer
		err = json.NewEncoder(&buff).Encode(resp)
		if err != nil {
			return errors.Wrapf(err, "cache.CacheAndWriteJson -> json.NewEncoder(&buff).Encode")
		}
		data = buff.Bytes()
		err = SetData(cacheName, data)
		if err != nil {
			return errors.Wrapf(err, "cache.CacheAndWriteJson -> SetData(%s)", cacheName)
		}
	}

	_, err := w.Write(data)
	if err != nil {
		return errors.Wrapf(err, "cache.CacheAndWriteJson -> w.Write")
	}
	return nil
}

func MustCacheAndWriteJson(w http.ResponseWriter, cacheName string, responseSupplier func() (interface{}, error)) {
	if err := CacheAndWriteJson(w, cacheName, responseSupplier); err != nil {
		panic(err)
	}
}
