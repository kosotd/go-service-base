package config

const serverPortEnvKey = "SERVER_PORT"
const cacheExpirationEnvKey = "CACHE_EXPIRATION"
const cacheUpdatePeriodEnvKey = "CACHE_UPDATE_PERIOD"
const logLevelEnvKey = "LOG_LEVEL"
const buildModeEnvKey = "BUILD_MODE"
const allowedOriginsEnvKey = "ALLOWED_ORIGINS"
const databasesEnvKey = "DATABASES"

const serverPortDefault = "8081"
const cacheExpirationDefault = "5m"
const cacheUpdatePeriodDefault = "1m"
const logLevelDefault = 1
const buildModeDefault = "debug"
var allowedOriginsDefault = []string{"*"}
var databasesDefault []string

const serverPort = "ServerPort"
const cacheExpiration = "CacheExpiration"
const cacheUpdatePeriod = "CacheUpdatePeriod"
const logLevel = "LogLevel"
const buildMode = "BuildMode"
const allowedOrigins = "AllowedOrigins"
const databases = "Databases"