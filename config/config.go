package config

import (
	"encoding/json"
	"flag"
	"github.com/joho/godotenv"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
	"os"
)

var conf map[string]interface{}
var envHelper envHelperImpl

func init() {
	conf = make(map[string]interface{})
	fileName := flag.String("config", "", "Full path to config file")
	flag.Parse()

	if utils.NotEmpty(*fileName) {
		loadFileConfiguration(*fileName)
	} else {
		loadEnvConfiguration()
	}

	utils.SetLogLevel(LogLevel())
}

func loadFileConfiguration(file string) {
	configFile, err := os.Open(file)
	utils.FailIfError(errors.Wrapf(err, "error open config file"))
	defer utils.CloseSafe(configFile)
	err = json.NewDecoder(configFile).Decode(&conf)
	utils.FailIfError(errors.Wrap(err, "error decode config json"))
	loadDefault()
}

func loadDefault() {
	if _, ok := conf[serverPort]; !ok {
		conf[serverPort] = serverPortDefault
	}
	if _, ok := conf[cacheExpiration]; !ok {
		conf[cacheExpiration] = cacheExpirationDefault
	}
	if _, ok := conf[cacheUpdatePeriod]; !ok {
		conf[cacheUpdatePeriod] = cacheUpdatePeriodDefault
	}
	if _, ok := conf[logLevel]; !ok {
		conf[logLevel] = logLevelDefault
	}
	if _, ok := conf[buildMode]; !ok {
		conf[buildMode] = buildModeDefault
	}
	if _, ok := conf[allowedOrigins]; !ok {
		conf[allowedOrigins] = allowedOriginsDefault
	}
	if _, ok := conf[databases]; !ok {
		conf[databases] = databasesDefault
	}
}

func loadEnvConfiguration() {
	_ = godotenv.Load("./config.env")
	conf[serverPort] = envHelper.GetEnvString(serverPortEnvKey, serverPortDefault)
	conf[cacheExpiration] = envHelper.GetEnvString(cacheExpirationEnvKey, cacheExpirationDefault)
	conf[cacheUpdatePeriod] = envHelper.GetEnvString(cacheUpdatePeriodEnvKey, cacheUpdatePeriodDefault)
	conf[logLevel] = envHelper.GetEnvInt(logLevelEnvKey, logLevelDefault)
	conf[buildMode] = envHelper.GetEnvString(buildModeEnvKey, buildModeDefault)
	conf[allowedOrigins] = envHelper.GetEnvStringList(allowedOriginsEnvKey, allowedOriginsDefault)
	conf[databases] = envHelper.GetEnvStringList(databasesEnvKey, databasesDefault)
}

func LoadEnvironment(loadFunc func(EnvHelper) map[string]interface{}) {
	env := loadFunc(envHelper)
	for k, v := range env {
		conf[k] = v
	}
}
