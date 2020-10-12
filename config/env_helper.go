package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type EnvHelper interface {
	GetEnvString(key, def string) string
	GetEnvInt(key string, def int) int
	GetEnvInt64(key string, def int64) int64
	GetEnvStringList(key string, def []string) []string
	GetEnvMap(key string, def map[string]interface{}) map[string]interface{}
}

type envHelperImpl struct{}

func (envHelperImpl) GetEnvString(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}

func (envHelperImpl) GetEnvInt(key string, def int) int {
	if value, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}
	}
	return def
}

func (envHelperImpl) GetEnvInt64(key string, def int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if val, err := strconv.ParseInt(value, 0, 64); err == nil {
			return val
		}
	}
	return def
}

func (envHelperImpl) GetEnvStringList(key string, def []string) []string {
	if value, ok := os.LookupEnv(key); ok {
		var list []string
		if err := json.Unmarshal([]byte(value), &list); err == nil {
			return list
		}
	}
	return def
}

func (envHelperImpl) GetEnvMap(key string, def map[string]interface{}) map[string]interface{} {
	if value, ok := os.LookupEnv(key); ok {
		var res map[string]interface{}
		if err := json.Unmarshal([]byte(value), &res); err == nil {
			return res
		}
	}
	return def
}
