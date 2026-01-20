package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"policy_mgnt/utils"

	"policy_mgnt/utils/gocommon/api"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type listResp struct {
	Count int           `json:"count"`
	Data  []interface{} `json:"data"`
}

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
	dbFile, err := os.CreateTemp("", "policy_mgnt_test")
	if err != nil {
		t.Fatal(err)
	}
	dbFile.Close()
	env := map[string]string{
		"DB_DRIVER": "sqlite3",
		"DB_URL":    dbFile.Name(),
	}
	teardown := SetUpEnv(t, env)
	utils.InitDB()

	return func(t *testing.T) {
		teardown(t)
		api.DisconnectDB()
		os.Remove(dbFile.Name())
	}
}

// AssertListResponse assert body is list response
func AssertListResponse(t *testing.T, respBody []byte, assertCount int) []interface{} {
	var resp listResp
	err := json.Unmarshal(respBody, &resp)
	assert.Nil(t, err)
	assert.Equal(t, resp.Count, assertCount)
	return resp.Data
}

// AssertError assery response is specify error
func AssertError(t *testing.T, resp *httptest.ResponseRecorder, expectErr *api.Error) {
	cStr := strconv.Itoa(expectErr.Code)
	code, _ := strconv.Atoi(cStr[:3])
	assert.Equal(t, resp.Code, code)
	value, _ := json.Marshal(expectErr)
	assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
}

// CoverError cover request error
func CoverError(t *testing.T, req *http.Request, router *gin.Engine, expectErr *api.Error) {
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	AssertError(t, resp, expectErr)
}
