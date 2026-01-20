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

func TestObligationTypeRestHandler_Set(t *testing.T) {
	Convey("set", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligationType := mock.NewMockObligationType(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationTypeRestHandler{
			obligationType: mockObligationType,
			hydra:          mockHydra,
			setSchemaStr: newJSONSchema(`{
				"type": "object",
				"required": ["name", "schema", "applicable_resource_types"],
				"properties": {
					"name": {"type": "string"},
					"description": {"type": "string"},
					"schema": {},
					"default_value": {},
					"ui_schema": {},
					"applicable_resource_types": {
						"type": "object",
						"required": ["unlimited"],
						"properties": {
							"unlimited": {"type": "boolean"},
							"resource_types": {"type": "array"}
						}
					}
				}
			}`),
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			reqBody := map[string]any{
				"name":   "test-obligation",
				"schema": map[string]any{"type": "string"},
				"applicable_resource_types": map[string]any{
					"unlimited": true,
				},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader(body))
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

			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("成功设置 - unlimited为true", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			reqBody := map[string]any{
				"name":          "test-obligation",
				"description":   "test description",
				"schema":        map[string]any{"type": "string"},
				"default_value": "default",
				"ui_schema":     map[string]any{"widget": "text"},
				"applicable_resource_types": map[string]any{
					"unlimited": true,
				},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("成功设置 - 指定资源类型", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			reqBody := map[string]any{
				"name":   "test-obligation",
				"schema": map[string]any{"type": "string"},
				"applicable_resource_types": map[string]any{
					"unlimited": false,
					"resource_types": []map[string]any{
						{
							"id": "doc",
							"applicable_operations": map[string]any{
								"unlimited": true,
							},
						},
					},
				},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("成功设置 - 指定操作", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			reqBody := map[string]any{
				"name":   "test-obligation",
				"schema": map[string]any{"type": "string"},
				"applicable_resource_types": map[string]any{
					"unlimited": false,
					"resource_types": []map[string]any{
						{
							"id": "doc",
							"applicable_operations": map[string]any{
								"unlimited": false,
								"operations": []map[string]any{
									{"id": "read"},
									{"id": "write"},
								},
							},
						},
					},
				},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("参数错误 - unlimited为false但缺少resource_types", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			reqBody := map[string]any{
				"name":   "test-obligation",
				"schema": map[string]any{"type": "string"},
				"applicable_resource_types": map[string]any{
					"unlimited": false,
				},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("参数错误 - operations unlimited为false但缺少operations", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			reqBody := map[string]any{
				"name":   "test-obligation",
				"schema": map[string]any{"type": "string"},
				"applicable_resource_types": map[string]any{
					"unlimited": false,
					"resource_types": []map[string]any{
						{
							"id": "doc",
							"applicable_operations": map[string]any{
								"unlimited": false,
							},
						},
					},
				},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("Set方法失败", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("set failed"))

			reqBody := map[string]any{
				"name":   "test-obligation",
				"schema": map[string]any{"type": "string"},
				"applicable_resource_types": map[string]any{
					"unlimited": true,
				},
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/authorization/v1/obligation-types/test-id", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

//nolint:dupl
func TestObligationTypeRestHandler_Delete(t *testing.T) {
	Convey("delete", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligationType := mock.NewMockObligationType(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationTypeRestHandler{
			obligationType: mockObligationType,
			hydra:          mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodDelete, "/api/authorization/v1/obligation-types/test-id", nil)
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
			mockObligationType.EXPECT().Delete(gomock.Any(), gomock.Any(), "test-id").Return(nil)

			req := httptest.NewRequest(http.MethodDelete, "/api/authorization/v1/obligation-types/test-id", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusNoContent)
		})

		Convey("Delete方法失败", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Delete(gomock.Any(), gomock.Any(), "test-id").Return(errors.New("delete failed"))

			req := httptest.NewRequest(http.MethodDelete, "/api/authorization/v1/obligation-types/test-id", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestObligationTypeRestHandler_Get(t *testing.T) {
	Convey("get", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligationType := mock.NewMockObligationType(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationTypeRestHandler{
			obligationType: mockObligationType,
			hydra:          mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligation-types?offset=0&limit=10", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("参数错误 - 无效的offset", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligation-types?offset=invalid&limit=10", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("成功获取列表", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(2, []interfaces.ObligationTypeInfo{
				{
					ID:           "obl1",
					Name:         "Obligation 1",
					Description:  "Description 1",
					Schema:       map[string]any{"type": "string"},
					DefaultValue: "default1",
					UiSchema:     map[string]any{"widget": "text"},
					ResourceTypeScope: interfaces.ObligationResourceTypeScopeInfo{
						Unlimited: true,
					},
				},
				{
					ID:           "obl2",
					Name:         "Obligation 2",
					Description:  "Description 2",
					Schema:       map[string]any{"type": "number"},
					DefaultValue: nil,
					UiSchema:     nil,
					ResourceTypeScope: interfaces.ObligationResourceTypeScopeInfo{
						Unlimited: false,
						Types: []interfaces.ObligationResourceTypeScope{
							{
								ResourceTypeID:   "doc",
								ResourceTypeName: "Document",
								OperationsScope: interfaces.ObligationOperationsScopeInfo{
									Unlimited: true,
								},
							},
						},
					},
				},
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligation-types?offset=0&limit=10", nil)
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

		Convey("Get方法失败", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil, errors.New("get failed"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligation-types?offset=0&limit=10", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestObligationTypeRestHandler_GetByID(t *testing.T) {
	Convey("getByID", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligationType := mock.NewMockObligationType(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationTypeRestHandler{
			obligationType: mockObligationType,
			hydra:          mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligation-types/test-id", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("成功获取", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().GetByID(gomock.Any(), gomock.Any(), "test-id").Return(interfaces.ObligationTypeInfo{
				ID:           "test-id",
				Name:         "Test Obligation",
				Description:  "Test Description",
				Schema:       map[string]any{"type": "string"},
				DefaultValue: "default",
				UiSchema:     map[string]any{"widget": "text"},
				ResourceTypeScope: interfaces.ObligationResourceTypeScopeInfo{
					Unlimited: false,
					Types: []interfaces.ObligationResourceTypeScope{
						{
							ResourceTypeID:   "doc",
							ResourceTypeName: "Document",
							OperationsScope: interfaces.ObligationOperationsScopeInfo{
								Unlimited: false,
								Operations: []interfaces.ObligationOperation{
									{ID: "read", Name: "Read"},
									{ID: "write", Name: "Write"},
								},
							},
						},
					},
				},
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligation-types/test-id", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(resp["id"], ShouldEqual, "test-id")
			So(resp["name"], ShouldEqual, "Test Obligation")
		})

		Convey("GetByID方法失败", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().GetByID(gomock.Any(), gomock.Any(), "test-id").Return(interfaces.ObligationTypeInfo{}, errors.New("get failed"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/obligation-types/test-id", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestObligationTypeRestHandler_QueryObligationTypes(t *testing.T) {
	Convey("queryObligationTypes", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockObligationType := mock.NewMockObligationType(ctrl)
		mockHydra := mock.NewMockHydra(ctrl)

		handler := &obligationTypeRestHandler{
			obligationType: mockObligationType,
			hydra:          mockHydra,
		}

		handler.RegisterPublic(r)

		Convey("验证失败 - 无权限", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active: false,
			}, gerrors.NewError(gerrors.PublicUnauthorized, "unauthorized"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligation-types?resource_type_id=doc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("参数错误 - 缺少resource_type_id", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligation-types", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("成功查询 - 不指定operation_ids", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[string][]interfaces.ObligationTypeInfo{
				"read": {
					{
						ID:           "obl1",
						Name:         "Obligation 1",
						Description:  "Description 1",
						Schema:       map[string]any{"type": "string"},
						DefaultValue: "default1",
						UiSchema:     map[string]any{"widget": "text"},
					},
				},
				"write": {
					{
						ID:           "obl2",
						Name:         "Obligation 2",
						Description:  "Description 2",
						Schema:       map[string]any{"type": "number"},
						DefaultValue: nil,
						UiSchema:     nil,
					},
				},
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligation-types?resource_type_id=doc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp []map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(len(resp), ShouldEqual, 2)
		})

		Convey("成功查询 - 指定operation_ids", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[string][]interfaces.ObligationTypeInfo{
				"read": {
					{
						ID:           "obl1",
						Name:         "Obligation 1",
						Description:  "Description 1",
						Schema:       map[string]any{"type": "string"},
						DefaultValue: "default1",
						UiSchema:     map[string]any{"widget": "text"},
					},
				},
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligation-types?resource_type_id=doc&operation_ids=read", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)

			var resp []map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			So(err, ShouldBeNil)
			So(len(resp), ShouldEqual, 1)
		})

		Convey("Query方法失败", func() {
			mockHydra.EXPECT().Introspect("test-token").Return(interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user1",
			}, nil)
			mockObligationType.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("query failed"))

			req := httptest.NewRequest(http.MethodGet, "/api/authorization/v1/query-obligation-types?resource_type_id=doc", nil)
			req.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestObligationTypeRestHandler_ResourceTypeScopeInfoToString(t *testing.T) {
	Convey("resourceTypeScopeInfoToString", t, func() {
		handler := &obligationTypeRestHandler{}

		Convey("unlimited为true", func() {
			info := interfaces.ObligationResourceTypeScopeInfo{
				Unlimited: true,
			}
			result := handler.resourceTypeScopeInfoToString(info)
			So(result["unlimited"], ShouldEqual, true)
			So(len(result["resource_types"].([]any)), ShouldEqual, 0)
		})

		Convey("unlimited为false - 操作unlimited为true", func() {
			info := interfaces.ObligationResourceTypeScopeInfo{
				Unlimited: false,
				Types: []interfaces.ObligationResourceTypeScope{
					{
						ResourceTypeID:   "doc",
						ResourceTypeName: "Document",
						OperationsScope: interfaces.ObligationOperationsScopeInfo{
							Unlimited: true,
						},
					},
				},
			}
			result := handler.resourceTypeScopeInfoToString(info)
			So(result["unlimited"], ShouldEqual, false)
			resourceTypes := result["resource_types"].([]any)
			So(len(resourceTypes), ShouldEqual, 1)
			rt := resourceTypes[0].(map[string]any)
			So(rt["id"], ShouldEqual, "doc")
			So(rt["name"], ShouldEqual, "Document")
			ops := rt["applicable_operations"].(map[string]any)
			So(ops["unlimited"], ShouldEqual, true)
		})

		Convey("unlimited为false - 操作unlimited为false", func() {
			info := interfaces.ObligationResourceTypeScopeInfo{
				Unlimited: false,
				Types: []interfaces.ObligationResourceTypeScope{
					{
						ResourceTypeID:   "doc",
						ResourceTypeName: "Document",
						OperationsScope: interfaces.ObligationOperationsScopeInfo{
							Unlimited: false,
							Operations: []interfaces.ObligationOperation{
								{ID: "read", Name: "Read"},
								{ID: "write", Name: "Write"},
							},
						},
					},
				},
			}
			result := handler.resourceTypeScopeInfoToString(info)
			So(result["unlimited"], ShouldEqual, false)
			resourceTypes := result["resource_types"].([]any)
			So(len(resourceTypes), ShouldEqual, 1)
			rt := resourceTypes[0].(map[string]any)
			ops := rt["applicable_operations"].(map[string]any)
			So(ops["unlimited"], ShouldEqual, false)
			operations := ops["operations"].([]any)
			So(len(operations), ShouldEqual, 2)
		})

		Convey("多个资源类型", func() {
			info := interfaces.ObligationResourceTypeScopeInfo{
				Unlimited: false,
				Types: []interfaces.ObligationResourceTypeScope{
					{
						ResourceTypeID:   "doc",
						ResourceTypeName: "Document",
						OperationsScope: interfaces.ObligationOperationsScopeInfo{
							Unlimited: true,
						},
					},
					{
						ResourceTypeID:   "file",
						ResourceTypeName: "File",
						OperationsScope: interfaces.ObligationOperationsScopeInfo{
							Unlimited: false,
							Operations: []interfaces.ObligationOperation{
								{ID: "download", Name: "Download"},
							},
						},
					},
				},
			}
			result := handler.resourceTypeScopeInfoToString(info)
			So(result["unlimited"], ShouldEqual, false)
			resourceTypes := result["resource_types"].([]any)
			So(len(resourceTypes), ShouldEqual, 2)
		})
	})
}
