package config

import (
	"fmt"
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

func MustInt(key string) int {
	val, ok := conf[key]
	if !ok {
		panic(fmt.Sprintf("property %s not found", key))
	}
	if res, ok := val.(float64); ok {
		return int(res)
	}
	if res, ok := val.(int); ok {
		return res
	}
	panic(fmt.Sprintf("property %s is not number", key))
}

func MustInt64(key string) int64 {
	val, ok := conf[key]
	if !ok {
		panic(fmt.Sprintf("property %s not found", key))
	}
	if res, ok := val.(float64); ok {
		return int64(res)
	}
	if res, ok := val.(int64); ok {
		return res
	}
	panic(fmt.Sprintf("property %s is not number", key))
}

func MustString(key string) string {
	val, ok := conf[key]
	if !ok {
		panic(fmt.Sprintf("property %s not found", key))
	}
	res, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("property %s is not string", key))
	}
	return res
}

func MustStringList(key string) []string {
	val, ok := conf[key]
	if !ok {
		panic(fmt.Sprintf("property %s not found", key))
	}
	if res, ok := val.([]string); ok {
		return res
	}
	if res, ok := val.([]interface{}); ok {
		strs := make([]string, 0, len(res))
		for _, r := range res {
			str, ok := r.(string)
			if !ok {
				panic(fmt.Sprintf("property %s is not string list", key))
			}
			strs = append(strs, str)
		}
		return strs
	}
	panic(fmt.Sprintf("property %s is not string list", key))
}

func MustDuration(key string) time.Duration {
	val, ok := conf[key]
	if !ok {
		panic(fmt.Sprintf("property %s not found", key))
	}
	resStr, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("property %s is not string", key))
	}
	duration, err := time.ParseDuration(resStr)
	if err != nil {
		panic(fmt.Sprintf("property %s must have duration format", key))
	}
	return duration
}