package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"policy_mgnt/test"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setUpPageParamsRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		offset, limit := getPageParams(c)
		if c.IsAborted() {
			return
		}
		c.JSON(200, gin.H{
			"offset": offset,
			"limit":  limit,
		})
	})
	return router
}

// 分页默认参数
func TestPageParamsDefault(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	router := setUpPageParamsRouter()

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, string(`{"offset":0,"limit":20}`), string(resp.Body.Bytes()))
}

// 分页 offset 参数最小值
func TestPageParamsStartMin(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	router := setUpPageParamsRouter()
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/?offset=0", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, `{"offset":0,"limit":20}`, string(resp.Body.Bytes()))
}

// 分页 limit 参数最小值
func TestPageParamsLimitMin(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	router := setUpPageParamsRouter()
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/?limit=1", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, `{"offset":0,"limit":1}`, string(resp.Body.Bytes()))
}

// 分页 offset 参数最大值
func TestPageParamsLimitMax(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	router := setUpPageParamsRouter()
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/?limit=1000", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, `{"offset":0,"limit":1000}`, string(resp.Body.Bytes()))
}

// 分页参数不合法
func TestPageParamsInvalid(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	var tests = []struct {
		offset  interface{}
		limit   interface{}
		wanterr *api.Error
	}{
		// offset is not int
		{"a", 1, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"offset"}}})},
		// offset < 0
		{-1, 1, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"offset"}}})},
		// limit is not int
		{1, "a", errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"limit"}}})},
		// limit < 1
		{1, 0, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"limit"}}})},
		// limit > 1000
		{1, 1001, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"limit"}}})},
	}

	for _, each := range tests {
		router := setUpPageParamsRouter()
		req := httptest.NewRequest("GET", fmt.Sprintf("/?offset=%v&limit=%v", each.offset, each.limit), nil)
		test.CoverError(t, req, router, each.wanterr)
	}
}

func setUpListResponseRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		listResponse(c, 10, []string{
			"test1", "test2", "test3",
		})
	})
	return router
}

// 多条数据返回
func TestListResponse(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	router := setUpListResponseRouter()
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	data := test.AssertListResponse(t, resp.Body.Bytes(), 10)
	assertData, _ := json.Marshal(data)
	assert.JSONEq(t, `["test1","test2","test3"]`, string(assertData))
}

func setUErrorResponseRouter(err error) *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		errorResponse(c, err)
	})
	return router
}

// 错误响应
func TestErrorResponse(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	err := &api.Error{
		Code:    500000000,
		Message: "Server unavailable",
		Cause:   "Unknown",
		Detail: map[string]interface{}{
			"cause1": "result",
			"cause2": "result",
			"cause3": "result",
			"cause4": "result",
		},
	}

	router := setUErrorResponseRouter(err)
	req := httptest.NewRequest("GET", "/", nil)
	test.CoverError(t, req, router, err)
}

// 错误响应，未知错误
func TestErrorResponseUnknown(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	err := fmt.Errorf("error")
	router := setUErrorResponseRouter(err)
	req := httptest.NewRequest("GET", "/", nil)
	assertErr := errors.ErrInternalServerErrorPublic(&api.ErrorInfo{Cause: err.Error()})
	test.CoverError(t, req, router, assertErr)
}

func setUpParamArrayRouter(key string) *gin.Engine {
	router := gin.Default()
	router.GET("/:"+key, func(c *gin.Context) {
		c.JSON(200, paramArray(c, key))
	})
	return router
}

func TestParamArray(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	key := "name"
	router := setUpParamArrayRouter(key)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/abc", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, string(resp.Body.Bytes()), `["abc"]`)

	resp = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/abc,bcd,cde", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, string(resp.Body.Bytes()), `["abc","bcd","cde"]`)
}

func TestParseBody(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	// example with models.NetworkRestriction
	var params models.NetworkRestriction
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		c.JSON(400, parseBody(c, &params))
	})
	// 传入body字段错误
	var tests = []struct {
		input    []byte
		wantcode int
		params   []string
	}{
		{[]byte(`123123`), 400, []string{""}},
		{[]byte(`{"name": 1, "start_ip": "", "end_ip": "", "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"name"}},
		{[]byte(`{"name": "", "start_ip": 1, "end_ip": "", "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"start_ip"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": 1, "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"end_ip"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": 1, "netmask": "", "net_type": ""}`), 400, []string{"ip_address"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": "", "netmask": 1, "net_type": ""}`), 400, []string{"netmask"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": "", "netmask": "", "net_type": 1}`), 400, []string{"net_type"}},
	}

	for _, test := range tests {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", bytes.NewReader(test.input))
		router.ServeHTTP(resp, req)
		assert.Equal(t, test.wantcode, resp.Code)
		// 状态码不为200，判断detail
		if test.wantcode != 200 {
			if test.wantcode == 400 {
				expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
		}
	}
}

// func TestNotifyUpdatePolicy(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	client := mock_api.NewMockProduce(ctrl)
// 	defer ctrl.Finish()

// 	mockTopic := "policy-data-topic"
// 	mockMessage := []byte(`{"policy_source_service": "policy-management"}`)
// 	client.EXPECT().Publish(mockTopic, mockMessage).Return(nil)

// 	err := notifyUpdatePolicy(client)
// 	assert.Equal(t, nil, err)
// }
