package config

import (
	"gotest.tools/assert"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoadFileConfiguration(t *testing.T) {
	conf = make(map[string]interface{})

	loadFileConfiguration("./config.json")
	assert.Equal(t, ServerPort(), "8082")
	assert.Equal(t, CacheExpiration(), 6 * time.Minute)
	assert.Equal(t, CacheUpdatePeriod(), 2 * time.Minute)
	assert.Equal(t, LogLevel(), 2)
	assert.Equal(t, BuildMode(), "release")
	assert.Equal(t, reflect.DeepEqual(AllowedOrigins(), []string{"http://localhost"}), true)
	assert.Equal(t, reflect.DeepEqual(Databases(), []string{"OracleDb;oracle:user/password@localhost:1521/orcl"}), true)
	assert.Equal(t, MustInt("Test"), 1)
	assert.Equal(t, MustInt64("Test"), int64(1))
	assert.Equal(t, MustString("Test1"), "test")
	assert.Equal(t, reflect.DeepEqual(MustStringList("Test2"), []string{"test"}), true)
	assert.Equal(t, MustDuration("Test3"), 1 * time.Minute)

	conf = make(map[string]interface{})

	loadFileConfiguration("./empty_config.json")
	assert.Equal(t, ServerPort(), "8081")
	assert.Equal(t, CacheExpiration(), 5 * time.Minute)
	assert.Equal(t, CacheUpdatePeriod(), 1 * time.Minute)
	assert.Equal(t, LogLevel(), 1)
	assert.Equal(t, BuildMode(), "debug")
	assert.Equal(t, reflect.DeepEqual(AllowedOrigins(), []string{"*"}), true)
	var strs []string
	assert.Equal(t, reflect.DeepEqual(Databases(), strs), true)
}

func TestLoadEnvConfiguration(t *testing.T) {
	conf = make(map[string]interface{})

	_ = os.Remove("config.env")
	loadEnvConfiguration()
	assert.Equal(t, ServerPort(), "8081")
	assert.Equal(t, CacheExpiration(), 5 * time.Minute)
	assert.Equal(t, CacheUpdatePeriod(), 1 * time.Minute)
	assert.Equal(t, LogLevel(), 1)
	assert.Equal(t, BuildMode(), "debug")
	assert.Equal(t, reflect.DeepEqual(AllowedOrigins(), []string{"*"}), true)
	var strs []string
	assert.Equal(t, reflect.DeepEqual(Databases(), strs), true)

	conf = make(map[string]interface{})

	content := `
SERVER_PORT=8082
CACHE_EXPIRATION=6m
CACHE_UPDATE_PERIOD=2m
LOG_LEVEL=2
BUILD_MODE=release
ALLOWED_ORIGINS=["http://localhost"]
DATABASES=["OracleDb;oracle:user/password@localhost:1521/orcl"]
TEST=1
TEST1=test
TEST2=["test"]
TEST3="1m"
TEST4=2
`
	_ = ioutil.WriteFile("config.env", []byte(content), os.ModePerm)
	defer func() { _ = os.Remove("config.env") }()
	loadEnvConfiguration()
	LoadEnvironment(func(helper EnvHelper) map[string]interface{} {
		return map[string]interface{}{
			"Test": helper.GetEnvInt("TEST", 0),
			"Test1": helper.GetEnvString("TEST1", ""),
			"Test2": helper.GetEnvStringList("TEST2", []string{}),
			"Test3": helper.GetEnvString("TEST3", ""),
			"Test4": helper.GetEnvInt64("TEST4", 0),
		}
	})

	assert.Equal(t, ServerPort(), "8082")
	assert.Equal(t, CacheExpiration(), 6 * time.Minute)
	assert.Equal(t, CacheUpdatePeriod(), 2 * time.Minute)
	assert.Equal(t, LogLevel(), 2)
	assert.Equal(t, BuildMode(), "release")
	assert.Equal(t, reflect.DeepEqual(AllowedOrigins(), []string{"http://localhost"}), true)
	assert.Equal(t, reflect.DeepEqual(Databases(), []string{"OracleDb;oracle:user/password@localhost:1521/orcl"}), true)
	assert.Equal(t, MustInt("Test"), 1)
	assert.Equal(t, MustString("Test1"), "test")
	assert.Equal(t, reflect.DeepEqual(MustStringList("Test2"), []string{"test"}), true)
	assert.Equal(t, MustDuration("Test3"), 1 * time.Minute)
	assert.Equal(t, MustInt64("Test4"), int64(2))
}