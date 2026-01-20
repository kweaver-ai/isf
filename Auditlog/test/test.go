package test

import (
	"os"
	"testing"

	"AuditLog/gocommon/api"

	"github.com/gin-gonic/gin"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// SetUpGin set up gin for unit test
func SetUpGin(t *testing.T) func(*testing.T) {
	oldMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func(t *testing.T) {
		gin.SetMode(oldMode)
	}
}

// SetUpEnv set environment with test
func SetUpEnv(t *testing.T, env map[string]string) func(*testing.T) {
	oldEnv := make(map[string]string)
	for key, value := range env {
		oldEnv[key] = os.Getenv(key)
		os.Setenv(key, value)
	}
	return func(t *testing.T) {
		for key, value := range oldEnv {
			os.Setenv(key, value)
		}
	}
}

// SetUpDB database mock
func SetUpDB(t *testing.T) func(*testing.T) {
	dbFile, err := os.CreateTemp("", "anyshare_test")
	if err != nil {
		t.Fatal(err)
	}
	dbFile.Close()
	env := map[string]string{
		"DB_DRIVER": "sqlite3",
		"DB_URL":    dbFile.Name(),
	}
	teardown := SetUpEnv(t, env)

	return func(t *testing.T) {
		teardown(t)
		api.DisconnectDB()
		os.Remove(dbFile.Name())
	}
}
