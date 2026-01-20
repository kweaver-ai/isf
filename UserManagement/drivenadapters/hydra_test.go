package drivenadapters

import (
	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"UserManagement/common"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHydra(h *http.Client, s *httptest.Server) *hydra {
	return &hydra{
		log:          common.NewLogger(),
		client:       h,
		adminAddress: s.URL,
	}
}

//nolint:lll
func TestRegister(t *testing.T) {
	Convey("Register", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}

			mapInfo := make(map[string]interface{})
			err = jsoniter.Unmarshal(body, &mapInfo)
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}

			clientInfo, _ := jsoniter.Marshal(gin.H{
				"client_id":     "xxx-xxx-xxx-xxx",
				"client_secret": "some-secret",
			})

			tmpErr, _ := jsoniter.Marshal(gin.H{
				"error":             "invalid_client_metadata",
				"error_description": "The value of one of the Client Metadata fields is invalid and the server has rejected this request. Note that an Authorization Server MAY choose to substitute a valid value for any requested parameter of a Client's Metadata. Field client_secret must contain a secret that is at least 6 characters long.",
			})

			if mapInfo["client_name"].(string) != "" {
				w.WriteHeader(201)
				_, _ = w.Write(clientInfo)
			} else {
				w.WriteHeader(400)
				_, _ = w.Write(tmpErr)
			}
		}))
		defer ts.Close()

		Convey("register failed", func() {
			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts)
			_, err := mHydra.Register("", "some-secret", 10)
			assert.NotEqual(t, err, nil)
		})

		Convey("register success", func() {
			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts)
			_, err := mHydra.Register("test", "some-secret", 10)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDelete(t *testing.T) {
	Convey("Delete", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
			_, _ = w.Write([]byte("success\n"))
		}))
		defer ts.Close()

		Convey("delete success", func() {
			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts)
			err := mHydra.Delete("xxx")
			assert.Equal(t, err, nil)
		})
	})
}

func TestUpdate(t *testing.T) {
	Convey("Update", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}

			mapInfo := make([]interface{}, 0)
			err = jsoniter.Unmarshal(body, &mapInfo)
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}

			clientInfo, _ := jsoniter.Marshal(gin.H{
				"client_id":     "xxx-xxx-xxx-xxx",
				"client_secret": "some-secret"})

			w.WriteHeader(200)
			_, _ = w.Write(clientInfo)
		}))
		defer ts.Close()

		Convey("Update success", func() {
			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts)
			err := mHydra.Update("xxx", "test", "some-secret")
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteConsentAndLogin(t *testing.T) {
	Convey("DeleteConsentAndLogin", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Convey("DeleteConsentAndLogin userID empty", func() {
			ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(204)
				_, _ = w.Write([]byte("success\n"))
			}))
			defer ts1.Close()

			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts1)
			err := mHydra.DeleteConsentAndLogin("xxx", "")
			assert.Equal(t, err, nil)
		})

		Convey("DeleteConsentAndLogin 204 success", func() {
			ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(204)
				_, _ = w.Write([]byte("success\n"))
			}))
			defer ts1.Close()

			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts1)
			err := mHydra.DeleteConsentAndLogin("xxx", "xxxxx")
			assert.Equal(t, err, nil)
		})

		Convey("DeleteConsentAndLogin 404 success", func() {
			ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
				_, _ = w.Write([]byte("success\n"))
			}))
			defer ts2.Close()

			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts2)
			err := mHydra.DeleteConsentAndLogin("xxx", "xxxxx")
			assert.Equal(t, err, nil)
		})

		Convey("DeleteConsentAndLogin 200 fail", func() {
			ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("success\n"))
			}))
			defer ts2.Close()

			client := httpclient.NewRawHTTPClient()
			mHydra := newHydra(client, ts2)
			err := mHydra.DeleteConsentAndLogin("xxx", "xxxxx")
			assert.NotEqual(t, err, nil)
		})
	})
}
