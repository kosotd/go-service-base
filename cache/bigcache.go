package cache

import (
	"github.com/allegro/bigcache"
	"github.com/kosotd/go-service-base/config"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
)

var cache *bigcache.BigCache

func init() {
	defaultConfig := bigcache.DefaultConfig(config.CacheExpiration())
	defaultConfig.CleanWindow = config.CacheUpdatePeriod()
	var err error
	cache, err = bigcache.NewBigCache(defaultConfig)
	utils.FailIfError(errors.Wrapf(err, "error init bigcache"))
}

func MustSetData(key string, data []byte) {
	if err := cache.Set(key, data); err != nil {
		panic(errors.Wrapf(err, "error cache set data with key: %s", key))
	}
}

func SetData(key string, data []byte) error {
	if err := cache.Set(key, data); err != nil {
		return errors.Wrapf(err, "cache.SetData -> cache.Set(%s)", key)
	}
	return nil
}

func GetData(key string) ([]byte, bool) {
	if data, err := cache.Get(key); err == nil {
		return data, true
	}
	return nil, false
}

func Close() {
	if cache != nil {
		_ = cache.Close()
	}
}
