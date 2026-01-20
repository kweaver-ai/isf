package helpers

import "os"

const (
	EnvPrefix     = "AUDIT_LOG_"
	EnvIsLocalDev = EnvPrefix + "LOCAL_DEV"  // AUDIT_LOG_LOCAL_DEV
	isDebugMode   = EnvPrefix + "DEBUG_MODE" // AUDIT_LOG_DEBUG_MODE

	isSQLPrint = EnvPrefix + "SQL_PRINT" // AUDIT_LOG_SQL_PRINT

	projPath = EnvPrefix + "PROJECT_PATH" // AUDIT_LOG_PROJECT_PATH
)

var mockIsLocalDev bool

func IsLocalDev() bool {
	return os.Getenv(EnvIsLocalDev) == "true" || mockIsLocalDev
}

func IsAaronLocalDev() bool {
	return os.Getenv(EnvIsLocalDev+"_AARON") == "true"
}

func SetIsLocalDev() {
	mockIsLocalDev = true
}

func IsDebugMode() bool {
	return os.Getenv(isDebugMode) == "true"
}

func IsOprLogShowLogForDebug() bool {
	return IsDebugMode()
}

func IsSQLPrint() bool {
	return os.Getenv(isSQLPrint) == "true"
}

func ProjectPathByEnv() string {
	return os.Getenv(projPath)
}
