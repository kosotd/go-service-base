package config

import (
	"github.com/pkg/errors"
	"time"
)

func ServerPort() string {
	return MustString(serverPort)
}

func CacheExpiration() time.Duration {
	return MustDuration(cacheExpiration)
}

func CacheUpdatePeriod() time.Duration {
	return MustDuration(cacheUpdatePeriod)
}

func LogLevel() int {
	return MustInt(logLevel)
}

func BuildMode() string {
	return MustString(buildMode)
}

func AllowedOrigins() []string {
	return MustStringList(allowedOrigins)
}

func Databases() []string {
	return MustStringList(databases)
}

func Int(key string) (int, error) {
	val, ok := conf[key]
	if !ok {
		return 0, errors.Errorf("property %s not found", key)
	}
	if res, ok := val.(float64); ok {
		return int(res), nil
	}
	if res, ok := val.(int); ok {
		return res, nil
	}
	return 0, errors.Errorf("property %s is not number", key)
}

func MustInt(key string) int {
	if res, err := Int(key); err != nil {
		panic(err)
	} else {
		return res
	}
}

func Int64(key string) (int64, error) {
	val, ok := conf[key]
	if !ok {
		return 0, errors.Errorf("property %s not found", key)
	}
	if res, ok := val.(float64); ok {
		return int64(res), nil
	}
	if res, ok := val.(int); ok {
		return int64(res), nil
	}
	if res, ok := val.(int64); ok {
		return res, nil
	}
	return 0, errors.Errorf("property %s is not number", key)
}

func MustInt64(key string) int64 {
	if res, err := Int64(key); err != nil {
		panic(err)
	} else {
		return res
	}
}

func String(key string) (string, error) {
	val, ok := conf[key]
	if !ok {
		return "", errors.Errorf("property %s not found", key)
	}
	res, ok := val.(string)
	if !ok {
		return "", errors.Errorf("property %s is not string", key)
	}
	return res, nil
}

func MustString(key string) string {
	if res, err := String(key); err != nil {
		panic(err)
	} else {
		return res
	}
}

func StringList(key string) ([]string, error) {
	val, ok := conf[key]
	if !ok {
		return nil, errors.Errorf("property %s not found", key)
	}
	if res, ok := val.([]string); ok {
		return res, nil
	}
	if res, ok := val.([]interface{}); ok {
		strs := make([]string, 0, len(res))
		for _, r := range res {
			str, ok := r.(string)
			if !ok {
				return nil, errors.Errorf("property %s is not string list", key)
			}
			strs = append(strs, str)
		}
		return strs, nil
	}
	return nil, errors.Errorf("property %s is not string list", key)
}

func MustStringList(key string) []string {
	if res, err := StringList(key); err != nil {
		panic(err)
	} else {
		return res
	}
}

func Duration(key string) (time.Duration, error) {
	val, ok := conf[key]
	if !ok {
		return 0, errors.Errorf("property %s not found", key)
	}
	resStr, ok := val.(string)
	if !ok {
		return 0, errors.Errorf("property %s is not string", key)
	}
	duration, err := time.ParseDuration(resStr)
	if err != nil {
		return 0, errors.Errorf("property %s must have duration format", key)
	}
	return duration, nil
}

func MustDuration(key string) time.Duration {
	if res, err := Duration(key); err != nil {
		panic(err)
	} else {
		return res
	}
}
