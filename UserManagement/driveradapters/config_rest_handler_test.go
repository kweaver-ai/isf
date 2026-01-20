package driveradapters

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
	"UserManagement/logics"

	gerrors "github.com/kweaver-ai/go-lib/error"
)

var (
	strENUS                 = "en-us"
	strGetConfigsPrivateURL = "/api/user-management/v1/configs/"
)

func TestNewConfigRESTHandler(t *testing.T) {
	sqlDB, _, err := sqlx.New()
	assert.Equal(t, err, nil)
	logics.SetDBPool(sqlDB)

	data := NewConfigRESTHandler()
	assert.NotEqual(t, data, nil)
}

func TestUpdateConfig(t *testing.T) {
	Convey("更新配置信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockLogicsConfig(ctrl)
		setConfigSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(setConfigSchemaStr))
		assert.Equal(t, err, nil)
		lConfig := &configRestHandler{
			config:            c,
			strDefaultUserPWd: "default_user_pwd",
			setConfigSchema:   setConfigSchema,
		}
		lConfig.RegisterPrivate(r)

		target := strGetConfigsPrivateURL
		Convey("参数不合法，fields存在非法参数", func() {
			reqParam := RegisterParam{
				Name:     "test",
				Password: "123",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			tempTarget := target + "xxxx"
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "invalid fields type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("参数不合法，fields存在非法参数，csf_level_enum", func() {
			reqParam := map[string]interface{}{
				"csf_level_enum": []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 5,
					},
				},
				"csf_level2_enum": []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 5,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			tempTarget := target + "csf_level_enum,csf_level2_enum"
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "invalid fields type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body不为json", func() {
			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer([]byte("xxx")), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内不存在特定配置参数", func() {
			reqParam := RegisterParam{
				Name:     "test",
				Password: "123",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "body.default_user_pwd is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内参数类型错误", func() {
			reqParam := map[string]interface{}{
				lConfig.strDefaultUserPWd: true,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "default_user_pwd: Invalid type. Expected: string, given: boolean")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tempErr := gerrors.NewError(gerrors.PublicConflict, "用户不存在")
		Convey("UpdateConfig错误", func() {
			reqParam := map[string]interface{}{
				lConfig.strDefaultUserPWd: "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			c.EXPECT().UpdateConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr)
			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, tempErr.Description)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			reqParam := map[string]interface{}{
				lConfig.strDefaultUserPWd: "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			c.EXPECT().UpdateConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:dupl,funlen
func TestUpdateManageConfig(t *testing.T) {
	Convey("管理员更新配置信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockLogicsConfig(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		setConfigSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(setConfigSchemaStr))
		assert.Equal(t, err, nil)
		lConfig := &configRestHandler{
			config:            c,
			strDefaultUserPWd: "default_user_pwd",
			hydra:             hydra,
			setConfigSchema:   setConfigSchema,
			strCSFLevelEnum:   "csf_level_enum",
			strCSFLevel2Enum:  "csf_level2_enum",
		}
		lConfig.RegisterPublic(r)

		target := "/api/user-management/v1/management/configs/"

		tokenInfo := interfaces.TokenIntrospectInfo{
			Active: false,
		}
		Convey("token无效，报错401", func() {
			reqParam := RegisterParam{
				Name:     "test",
				Password: "123",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + "xxxx"
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tokenInfo.Active = true
		Convey("参数不合法，fields存在非法参数", func() {
			reqParam := RegisterParam{
				Name:     "test",
				Password: "123",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + "xxxx"
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "invalid fields type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body不为json", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer([]byte("xxx")), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内不存在特定配置参数", func() {
			reqParam := RegisterParam{
				Name:     "test",
				Password: "123",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "body.default_user_pwd is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内参数类型错误", func() {
			reqParam := map[string]interface{}{
				lConfig.strDefaultUserPWd: true,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "default_user_pwd: Invalid type. Expected: string, given: boolean")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level2_enum在url 不在request的body里面", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevelEnum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 5,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevel2Enum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "body.csf_level2_enum is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level_enum在url 不在request的body里面", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevel2Enum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 51,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevelEnum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "body.csf_level_enum is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level_enum的value限制超过16", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevelEnum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 51,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevelEnum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf_level_enum.0.value: Must be less than or equal to 16")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level_enum的vakue限制小于5", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevelEnum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 4,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevelEnum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf_level_enum.0.value: Must be greater than or equal to 5")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level2_enum的value限制超过63", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevel2Enum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 63,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevel2Enum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf_level2_enum.0.value: Must be less than or equal to 62")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level2_enum的value限制小于51", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevel2Enum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 50,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevel2Enum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf_level2_enum.0.value: Must be greater than or equal to 51")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level_enum的name重复", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevelEnum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 5,
					},
					map[string]interface{}{
						"name":  "公开",
						"value": 5,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevelEnum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf level enum name is not unique")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level2_enum的name重复", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevel2Enum: []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 51,
					},
					map[string]interface{}{
						"name":  "公开",
						"value": 51,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevel2Enum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf level enum name is not unique")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level_enum的value重复", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevelEnum: []interface{}{
					map[string]interface{}{
						"name":  "公开1",
						"value": 5,
					},
					map[string]interface{}{
						"name":  "公开",
						"value": 5,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevelEnum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf level enum value is not unique")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("csf_level2_enum的value重复", func() {
			reqParam := map[string]interface{}{
				lConfig.strCSFLevel2Enum: []interface{}{
					map[string]interface{}{
						"name":  "公开1",
						"value": 51,
					},
					map[string]interface{}{
						"name":  "公开",
						"value": 51,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + lConfig.strCSFLevel2Enum
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "csf level enum value is not unique")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tempErr := rest.NewHTTPError("用户不存在", 400019001, nil)
		Convey("UpdateConfig错误", func() {
			reqParam := map[string]interface{}{
				lConfig.strDefaultUserPWd: "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			c.EXPECT().UpdateConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr)
			tempTarget := target + lConfig.strDefaultUserPWd
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, tempErr.Description)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			reqParam := map[string]interface{}{
				lConfig.strDefaultUserPWd: "xxx",
				"csf_level_enum": []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 5,
					},
				},
				"csf_level2_enum": []interface{}{
					map[string]interface{}{
						"name":  "公开",
						"value": 51,
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			c.EXPECT().UpdateConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			tempTarget := target + lConfig.strDefaultUserPWd + ",csf_level_enum,csf_level2_enum"
			result, _ := mockRequest(false, "PUT", tempTarget, bytes.NewBuffer(reqParamByte), r)

			respBody, _ := io.ReadAll(result.Body)
			fmt.Println(string(respBody))
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestCheckPWDValid(t *testing.T) {
	Convey("管理员检查密码是否符合格式", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockLogicsConfig(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		lConfig := &configRestHandler{
			config:            c,
			strDefaultUserPWd: "default_user_pwd",
			hydra:             hydra,
			mapLangs: map[string]interfaces.LangType{
				"zh_CN": interfaces.LTZHCN,
				"zh_TW": interfaces.LTZHTW,
				"en_US": interfaces.LTENUS,
			},
		}
		lConfig.RegisterPublic(r)

		target := "/api/user-management/v1/management/default-pwd-valid"

		tokenInfo := interfaces.TokenIntrospectInfo{
			Active: false,
		}

		Convey("token无效，报错401", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tokenInfo.Active = true
		Convey("无password参数，报错400", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("xxxx", rest.Forbidden, nil)
		Convey("检查格式报错", func() {
			common.SvcConfig.Lang = "zh_TW"
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			c.EXPECT().CheckDefaultPWD(gomock.Any(), gomock.Any()).AnyTimes().Return(false, "", testErr)

			tempTarget := target + "?password="
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		common.SvcConfig.Lang = "zh_CN"
		Convey("返回错误", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			c.EXPECT().CheckDefaultPWD(gomock.Any(), gomock.Any()).AnyTimes().Return(false, "xxxxxx", nil)

			tempTarget := target + "?password="
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam["result"], false)
			assert.Equal(t, respParam["err_msg"], "xxxxxx")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		common.SvcConfig.Lang = "en_US"
		Convey("返回正常", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			c.EXPECT().CheckDefaultPWD(gomock.Any(), gomock.Any()).AnyTimes().Return(true, "", nil)

			tempTarget := target + "?password="
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam["result"], true)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetLanguage(t *testing.T) {
	Convey("获取x-language优先语言", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lConfig := &configRestHandler{}

		Convey("按照顺序命中，测试中文优先", func() {
			common.SvcConfig.Lang = "aa"
			out := lConfig.GetLanguage("fr-CH,zh_cn,en-US")
			assert.Equal(t, out, "zh_CN")
		})

		Convey("按照顺序命中，测试英文优先", func() {
			common.SvcConfig.Lang = "aa"
			out := lConfig.GetLanguage("fr-CH,en-US,zh-CN")
			assert.Equal(t, out, "en_US")
		})

		Convey("测试使用优先使用匹配", func() {
			common.SvcConfig.Lang = "aa"
			out := lConfig.GetLanguage("fr-CH,zh_CN")
			assert.Equal(t, out, "zh_CN")
		})

		Convey("按照默认匹配", func() {
			common.SvcConfig.Lang = strENUS
			out := lConfig.GetLanguage("fr-CH")
			assert.Equal(t, out, "en_US")
		})

		Convey("按照默认匹配1", func() {
			common.SvcConfig.Lang = strENUS
			out := lConfig.GetLanguage("")
			assert.Equal(t, out, "en_US")
		})
	})
}

func TestGetConfigs(t *testing.T) {
	Convey("获取配置信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockLogicsConfig(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		lConfig := &configRestHandler{
			config:            c,
			strDefaultUserPWd: "default_user_pwd",
			hydra:             hydra,
			strCSFLevelEnum:   "csf_level_enum",
			strCSFLevel2Enum:  "csf_level2_enum",
			strShowCSFLevel2:  "show_csf_level2",
		}
		lConfig.RegisterPublic(r)

		target := strGetConfigsPrivateURL

		tokenInfo := interfaces.TokenIntrospectInfo{
			Active: false,
		}
		Convey("token无效，报错401", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + "csf_level_enum,csf_level2_enum"
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tokenInfo.Active = true
		Convey("参数不合法，fields存在非法参数", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + "xxxxcsf_level_enum,csf_level2_enum"
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "invalid fields type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := gerrors.NewError(gerrors.PublicConflict, "xxxx")
		Convey("GetConfig报错", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			c.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{}, testErr)

			tempTarget := target + "csf_level_enum,csf_level2_enum"
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, testErr.Description)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			c.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{
				CSFLevelEnum: map[string]int{
					"公开": 5,
					"秘密": 6,
					"机密": 7,
					"绝密": 8,
				},
				CSFLevel2Enum: map[string]int{
					"公开1": 51,
					"秘密1": 52,
					"机密1": 53,
					"绝密1": 54,
				},
				ShowCSFLevel2: true,
			}, nil)

			tempTarget := target + "csf_level_enum,csf_level2_enum,show_csf_level2"
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			temp1 := respParam["csf_level_enum"].([]interface{})
			assert.Equal(t, len(temp1), 4)
			temp11 := temp1[0].(map[string]interface{})
			assert.Equal(t, temp11["name"], "公开")
			assert.Equal(t, temp11["value"], float64(5))
			temp12 := temp1[1].(map[string]interface{})
			assert.Equal(t, temp12["name"], "秘密")
			assert.Equal(t, temp12["value"], float64(6))
			temp13 := temp1[2].(map[string]interface{})
			assert.Equal(t, temp13["name"], "机密")
			assert.Equal(t, temp13["value"], float64(7))
			temp14 := temp1[3].(map[string]interface{})
			assert.Equal(t, temp14["name"], "绝密")
			assert.Equal(t, temp14["value"], float64(8))
			temp2 := respParam["csf_level2_enum"].([]interface{})
			assert.Equal(t, len(temp2), 4)
			temp21 := temp2[0].(map[string]interface{})
			assert.Equal(t, temp21["name"], "公开1")
			assert.Equal(t, temp21["value"], float64(51))
			temp22 := temp2[1].(map[string]interface{})
			assert.Equal(t, temp22["name"], "秘密1")
			assert.Equal(t, temp22["value"], float64(52))
			temp23 := temp2[2].(map[string]interface{})
			assert.Equal(t, temp23["name"], "机密1")
			assert.Equal(t, temp23["value"], float64(53))
			temp24 := temp2[3].(map[string]interface{})
			assert.Equal(t, temp24["name"], "绝密1")
			assert.Equal(t, temp24["value"], float64(54))
			temp3 := respParam["show_csf_level2"].(bool)
			assert.Equal(t, temp3, true)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetConfigPrivates(t *testing.T) {
	Convey("获取配置信息（内部接口）", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockLogicsConfig(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		lConfig := &configRestHandler{
			config:            c,
			strDefaultUserPWd: "default_user_pwd",
			hydra:             hydra,
			strCSFLevelEnum:   "csf_level_enum",
			strCSFLevel2Enum:  "csf_level2_enum",
			strShowCSFLevel2:  "show_csf_level2",
		}
		lConfig.RegisterPrivate(r)

		target := strGetConfigsPrivateURL

		tokenInfo := interfaces.TokenIntrospectInfo{
			Active: false,
		}

		tokenInfo.Active = true
		Convey("参数不合法，fields存在非法参数", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			tempTarget := target + "show_csf_level2"
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "invalid fields type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := gerrors.NewError(gerrors.PublicConflict, "xxxx")
		Convey("GetConfig报错", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			c.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{}, testErr)

			tempTarget := target + "csf_level_enum,csf_level2_enum"
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respBody, _ := io.ReadAll(result.Body)
			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, testErr.Description)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			c.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{
				CSFLevelEnum: map[string]int{
					"公开": 5,
					"秘密": 6,
					"机密": 7,
					"绝密": 8,
				},
				CSFLevel2Enum: map[string]int{
					"公开1": 51,
					"秘密1": 52,
					"机密1": 53,
					"绝密1": 54,
				},
				ShowCSFLevel2: true,
			}, nil)

			tempTarget := target + "csf_level_enum,csf_level2_enum"
			result, _ := mockRequest(false, "GET", tempTarget, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			temp1 := respParam["csf_level_enum"].([]interface{})
			assert.Equal(t, len(temp1), 4)
			temp11 := temp1[0].(map[string]interface{})
			assert.Equal(t, temp11["name"], "公开")
			assert.Equal(t, temp11["value"], float64(5))
			temp12 := temp1[1].(map[string]interface{})
			assert.Equal(t, temp12["name"], "秘密")
			assert.Equal(t, temp12["value"], float64(6))
			temp13 := temp1[2].(map[string]interface{})
			assert.Equal(t, temp13["name"], "机密")
			assert.Equal(t, temp13["value"], float64(7))
			temp14 := temp1[3].(map[string]interface{})
			assert.Equal(t, temp14["name"], "绝密")
			assert.Equal(t, temp14["value"], float64(8))
			temp2 := respParam["csf_level2_enum"].([]interface{})
			assert.Equal(t, len(temp2), 4)
			temp21 := temp2[0].(map[string]interface{})
			assert.Equal(t, temp21["name"], "公开1")
			assert.Equal(t, temp21["value"], float64(51))
			temp22 := temp2[1].(map[string]interface{})
			assert.Equal(t, temp22["name"], "秘密1")
			assert.Equal(t, temp22["value"], float64(52))
			temp23 := temp2[2].(map[string]interface{})
			assert.Equal(t, temp23["name"], "机密1")
			assert.Equal(t, temp23["value"], float64(53))
			temp24 := temp2[3].(map[string]interface{})
			assert.Equal(t, temp24["name"], "绝密1")
			assert.Equal(t, temp24["value"], float64(54))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
