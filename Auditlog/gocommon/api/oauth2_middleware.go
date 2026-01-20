package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/log"
	"github.com/gin-gonic/gin"
)

// extraBearerToken 获取
func extraBearerToken(req *http.Request) (token string, err *Error) {
	hdr := req.Header.Get("Authorization")
	if hdr == "" {
		err = ErrUnauthorization(&ErrorInfo{Cause: "access_token empty"})
		return
	}

	// Example: Bearer xxxx
	tokenList := strings.SplitN(hdr, " ", 2)
	if len(tokenList) != 2 || strings.ToLower(tokenList[0]) != "bearer" {
		err = ErrUnauthorization(&ErrorInfo{Cause: "access_token invalid"})
		return
	}
	return tokenList[1], nil
}

// 获取用户角色、密级、显示名
func getUserRoles(userId string) (userRoles []string, csfLevel float64, name, deptPaths string, err error) {
	logger := NewTelemetryLogger(os.Stdout, log.InfoLevel, &LogOptionServiceInfo{
		Name: GetEnv("SERVICE_NAME", "audit-log"),
	})
	tagMarker := func(level field.Field, _ string, message field.Field) (tags []string) {
		if string(level.(field.StringField)) == "Error" {
			return []string{GetEnv("SERVICE_NAME", "audit-log")}
		}
		return []string{}
	}
	logger.AddTagMarker(tagMarker)
	var resp *http.Response
	var respBodyByte []byte
	var respData []map[string]interface{}
	url := getOwnersEndpoint(userId)
	resp, err = http.Get(url)
	if err != nil {
		logger.Errorf("ERROR: GetUserRoles:%v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Warnf("ERROR: GetUserRoles: Response = %v", resp.StatusCode)
		err = fmt.Errorf("ERROR: GetUserRoles: Response = %v", resp.StatusCode)
		return
	}
	respBodyByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("ERROR: GetUserRoles:%v\n", err)
		return
	}
	err = json.Unmarshal(respBodyByte, &respData)
	if err != nil {
		logger.Errorf("ERROR: GetUserRoles:%v\n", err)
		return
	}
	for _, v := range respData[0]["roles"].([]interface{}) {
		userRoles = append(userRoles, v.(string))
	}
	csfLevel = respData[0]["csf_level"].(float64)
	name = respData[0]["name"].(string)
	parentDeps := respData[0]["parent_deps"].([]interface{})
	// 用户部门信息
	deptPathsSlice := make([]string, 0)

	for _, deptInfo := range parentDeps {
		deptAllPaths := make([]string, 0)
		for _, item := range deptInfo.([]interface{}) {
			deptAllPaths = append(deptAllPaths, item.(map[string]interface{})["name"].(string))
		}

		deptPath := strings.Join(deptAllPaths, "/")
		deptPathsSlice = append(deptPathsSlice, deptPath)
	}
	deptPaths = strings.Join(deptPathsSlice, ", ")

	return
}

// Oauth2Middleware oauth2认证
func Oauth2Middleware(o OAuth2, trace Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		trace.SetInternalSpanName("Hydra认证")
		newCtx, span := trace.AddInternalTrace(c)
		defer func() { trace.TelemetrySpanEnd(span, err) }()

		token, apiErr := extraBearerToken(c.Request)
		if apiErr != nil {
			ErrorResponse(c, apiErr)
			return
		}
		result, err := o.Introspection(newCtx, token, nil)
		c.Set("userId", result.Subject)
		c.Set("userToken", token)
		c.Set("ip", result.Extra["login_ip"])
		var cType, udid, visitorType string
		if result.ClientID != result.Subject {
			cType = AssertString(result.Extra["client_type"])
			udid = AssertString(result.Extra["udid"])
			if v, ok := result.Extra["visitor_type"].(string); ok {
				switch v {
				case "realname":
					visitorType = "authenticated_user"
				case "anonymous":
					visitorType = "anonymous_user"
				}
			}
		}
		c.Set("clientType", cType)
		c.Set("udid", udid)
		c.Set("visitorType", visitorType)
		if err != nil {
			// TODO: 记日志
			ErrorResponse(c, ErrInternalServerErrorPublic(&ErrorInfo{Cause: "introspection failed"}))
			return
		}

		if !result.Active {
			ErrorResponse(c, ErrUnauthorization(&ErrorInfo{Cause: "access_token does not active"}))
			return
		}
		if result.ClientID != result.Subject {
			userRoles, csfLevel, name, deptPaths, _ := getUserRoles(result.Subject)
			c.Set("userRoles", userRoles)
			c.Set("csfLevel", csfLevel)
			c.Set("name", name)
			c.Set("accountType", "user")
			c.Set("dept_paths", deptPaths)
		} else {
			c.Set("userRoles", nil)
			c.Set("name", "")
			c.Set("accountType", "app")
		}
	}
}
