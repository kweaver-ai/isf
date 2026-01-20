//nolint:gocritic,funlen
package driveradapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	gerrors "github.com/kweaver-ai/go-lib/error"

	"Authorization/interfaces"
	"Authorization/interfaces/mock"
)

func TestObligationRestHandler_Add(t *testing.T) {
	Convey("add", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligation := mock.NewMockLogicsObligation(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationRestHandler{
			obligation: mockObligation,
			hydra:      mockHydra,
			addSchemaStr: newJSONSchema(`{
				"type": "object",
				"required": ["type_id", "name", "value"],
				"properties": {
					"type_id": {"type": "string"},
					"name": {"type": "string"},
					"description": {"type": "string"},
					"value": {}
				}
			}`),
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			reqBody := map[string]any{
				"type_id": "type1",
				"name":    "test-obligation",
				"value":   "test-value",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/authorization/v1/obligations", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("JSON格式错误", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/authorization/v1/obligations", bytes.NewReader([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("缺少必填字段", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			reqBody := map[string]any{
				"name": "test-obligation",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/authorization/v1/obligations", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("成功添加义务 - 无描述", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return("obligation-id-1", nil)

			reqBody := map[string]any{
				"type_id": "type1",
				"name":    "test-obligation",
				"value":   "test-value",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/authorization/v1/obligations", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusCreated)
			So(w.Header().Get("Location"), ShouldEqual, "/api/authorization/v1/obligations/obligation-id-1")

			var resp map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(resp["id"], ShouldEqual, "obligation-id-1")
		})

		Convey("成功添加义务 - 包含描述", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return("obligation-id-2", nil)

			reqBody := map[string]any{
				"type_id":     "type1",
				"name":        "test-obligation",
				"description": "test description",
				"value":       map[string]any{"key": "value"},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/authorization/v1/obligations", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusCreated)
		})

		Convey("添加失败 - 逻辑层错误", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("database error"))

			reqBody := map[string]any{
				"type_id": "type1",
				"name":    "test-obligation",
				"value":   "test-value",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/authorization/v1/obligations", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestObligationRestHandler_Update(t *testing.T) {
	Convey("update", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligation := mock.NewMockLogicsObligation(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationRestHandler{
			obligation: mockObligation,
			hydra:      mockHydra,
			updateSchemaStr: newJSONSchema(`{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"description": {"type": "string"},
					"value": {}
				}
			}`),
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			reqBody := map[string]any{
				"name": "updated-name",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/name", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("JSON格式错误", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/name", bytes.NewReader([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("成功更新名称", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Update(gomock.Any(), gomock.Any(), "obligation-id-1", "updated-name", true, "", false, nil, false).Return(nil)

			reqBody := map[string]any{
				"name": "updated-name",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/name", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("成功更新描述", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Update(gomock.Any(), gomock.Any(), "obligation-id-1", "", false, "updated-description", true, nil, false).Return(nil)

			reqBody := map[string]any{
				"description": "updated-description",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/description", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("成功更新值", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Update(gomock.Any(), gomock.Any(), "obligation-id-1", "", false, "", false, "new-value", true).Return(nil)

			reqBody := map[string]any{
				"value": "new-value",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/value", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("成功更新多个字段", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Update(gomock.Any(), gomock.Any(), "obligation-id-1", "updated-name", true, "updated-description", true, nil, false).Return(nil)

			reqBody := map[string]any{
				"name":        "updated-name",
				"description": "updated-description",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/name,description", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("字段存在但请求体缺少该字段", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			reqBody := map[string]any{
				"description": "updated-description",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/name", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("更新失败 - 逻辑层错误", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("database error"))

			reqBody := map[string]any{
				"name": "updated-name",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligations/obligation-id-1/name", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

//nolint:dupl
func TestObligationRestHandler_Delete(t *testing.T) {
	Convey("delete", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligation := mock.NewMockLogicsObligation(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationRestHandler{
			obligation: mockObligation,
			hydra:      mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodDelete, "/api/authorization/v1/obligations/obligation-id-1", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("成功删除", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Delete(gomock.Any(), gomock.Any(), "obligation-id-1").Return(nil)

			req := httptest.NewRequest(http.MethodDelete, "/api/authorization/v1/obligations/obligation-id-1", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("删除失败 - 逻辑层错误", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Delete(gomock.Any(), gomock.Any(), "obligation-id-1").Return(errors.New("database error"))

			req := httptest.NewRequest(http.MethodDelete, "/api/authorization/v1/obligations/obligation-id-1", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestObligationRestHandler_GetByID(t *testing.T) {
	Convey("getByID", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligation := mock.NewMockLogicsObligation(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationRestHandler{
			obligation: mockObligation,
			hydra:      mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations/obligation-id-1", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("成功获取义务", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			expectedObligation := interfaces.ObligationInfo{
				ID:          "obligation-id-1",
				TypeID:      "type1",
				Name:        "test-obligation",
				Description: "test description",
				Value:       "test-value",
			}

			mockObligation.EXPECT().GetByID(gomock.Any(), gomock.Any(), "obligation-id-1").Return(expectedObligation, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations/obligation-id-1", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(resp["id"], ShouldEqual, "obligation-id-1")
			So(resp["type_id"], ShouldEqual, "type1")
			So(resp["name"], ShouldEqual, "test-obligation")
			So(resp["description"], ShouldEqual, "test description")
			So(resp["value"], ShouldEqual, "test-value")
		})

		Convey("获取失败 - 义务不存在", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().GetByID(gomock.Any(), gomock.Any(), "non-existent-id").Return(interfaces.ObligationInfo{}, gerrors.NewError(gerrors.PublicNotFound, "not found"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations/non-existent-id", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})
	})
}

func TestObligationRestHandler_Get(t *testing.T) {
	Convey("get", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligation := mock.NewMockLogicsObligation(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationRestHandler{
			obligation: mockObligation,
			hydra:      mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("成功获取义务列表 - 默认参数", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			obligations := []interfaces.ObligationInfo{
				{
					ID:          "obligation-id-1",
					TypeID:      "type1",
					Name:        "test-obligation-1",
					Description: "description-1",
					Value:       "value-1",
				},
				{
					ID:          "obligation-id-2",
					TypeID:      "type2",
					Name:        "test-obligation-2",
					Description: "description-2",
					Value:       "value-2",
				},
			}

			mockObligation.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(2, obligations, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(resp["total_count"], ShouldEqual, 2)
			entries := resp["entries"].([]any)
			So(len(entries), ShouldEqual, 2)
		})

		Convey("成功获取义务列表 - 指定分页参数", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			obligations := []interfaces.ObligationInfo{
				{
					ID:     "obligation-id-3",
					TypeID: "type3",
					Name:   "test-obligation-3",
					Value:  "value-3",
				},
			}

			mockObligation.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, obligations, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations?offset=10&limit=10", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(resp["total_count"], ShouldEqual, 10)
			entries := resp["entries"].([]any)
			So(len(entries), ShouldEqual, 1)
		})

		Convey("成功获取义务列表 - 空列表", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			mockObligation.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, []interfaces.ObligationInfo{}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(resp["total_count"], ShouldEqual, 0)
			entries := resp["entries"].([]any)
			So(len(entries), ShouldEqual, 0)
		})

		Convey("获取失败 - 逻辑层错误", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil, errors.New("database error"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("参数错误 - offset非数字", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligations?offset=abc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

func TestObligationRestHandler_QueryObligation(t *testing.T) {
	Convey("queryObligation", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligation := mock.NewMockLogicsObligation(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationRestHandler{
			obligation: mockObligation,
			hydra:      mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligations?resource_type_id=doc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("缺少必填参数 resource_type_id", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligations", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("成功查询 - 仅指定资源类型", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			result := map[string][]interfaces.ObligationInfo{
				"read": {
					{
						ID:          "obligation-id-1",
						TypeID:      "type1",
						Name:        "test-obligation-1",
						Description: "description-1",
						Value:       "value-1",
					},
				},
				"write": {
					{
						ID:          "obligation-id-2",
						TypeID:      "type2",
						Name:        "test-obligation-2",
						Description: "description-2",
						Value:       "value-2",
					},
				},
			}

			mockObligation.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligations?resource_type_id=doc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp []map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(len(resp), ShouldEqual, 2)
		})

		Convey("成功查询 - 指定操作ID", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			result := map[string][]interfaces.ObligationInfo{
				"read": {
					{
						ID:     "obligation-id-1",
						TypeID: "type1",
						Name:   "test-obligation-1",
						Value:  "value-1",
					},
				},
			}

			mockObligation.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligations?resource_type_id=doc&operation_ids=read", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp []map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(len(resp), ShouldEqual, 1)
			So(resp[0]["operation_id"], ShouldEqual, "read")
		})

		Convey("成功查询 - 指定义务类型ID", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			result := map[string][]interfaces.ObligationInfo{
				"read": {
					{
						ID:     "obligation-id-1",
						TypeID: "type1",
						Name:   "test-obligation-1",
						Value:  "value-1",
					},
				},
			}

			mockObligation.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligations?resource_type_id=doc&obligation_type_ids=type1", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("成功查询 - 指定多个参数", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			result := map[string][]interfaces.ObligationInfo{
				"read": {
					{
						ID:     "obligation-id-1",
						TypeID: "type1",
						Name:   "test-obligation-1",
						Value:  "value-1",
					},
				},
			}

			mockObligation.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)

			req := httptest.NewRequest(http.MethodGet,
				"/api/authorization/v1/query-obligations?resource_type_id=doc&operation_ids=read&operation_ids=write&obligation_type_ids=type1&obligation_type_ids=type2", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("查询失败 - 逻辑层错误", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligation.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligations?resource_type_id=doc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("成功查询 - 空结果", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			result := map[string][]interfaces.ObligationInfo{}

			mockObligation.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligations?resource_type_id=doc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp []map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(len(resp), ShouldEqual, 0)
		})
	})
}
