// Package drivenadapters 消息队列
package drivenadapters

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

type mockOSSRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func TestGetDownloadURL(t *testing.T) {
	Convey("GetDownloadURL", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockDnHTTPClient(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		ossMock := &ossGateWay{
			log:    common.NewLogger(),
			client: c,
			trace:  trace,
		}

		testErr := errors.New("xxx")
		var outInfo interface{}

		ctx := context.Background()
		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("get error", func() {
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, testErr)
			_, err := ossMock.GetDownloadURL(ctx, &interfaces.Visitor{}, "", "")
			assert.Equal(t, err, testErr)
		})

		tempMap := make(map[string]string)
		tempMap["sdad"] = "sdadxxxssss"
		tempData := mockOSSRequest{
			Method:  "sdad",
			URL:     "xxads",
			Headers: tempMap,
		}
		tempBytes, err1 := jsoniter.Marshal(tempData)
		assert.Equal(t, err1, nil)

		err1 = jsoniter.Unmarshal(tempBytes, &outInfo)
		assert.Equal(t, err1, nil)
		Convey("success", func() {
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, nil)
			out, err := ossMock.GetDownloadURL(ctx, &interfaces.Visitor{}, "", "")
			assert.Equal(t, err, nil)
			assert.Equal(t, tempData.URL, out)
		})
	})
}

func TestUploadFile(t *testing.T) {
	Convey("UploadFile", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockDnHTTPClient(ctrl)
		bSuccess := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}

			if bSuccess {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("success\n"))
			} else {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}
		}))
		defer ts.Close()

		trace := mock.NewMockTraceClient(ctrl)
		ossMock := &ossGateWay{
			log:          common.NewLogger(),
			client:       c,
			rawClient:    httpclient.NewRawHTTPClient(),
			adminAddress: ts.URL,
			trace:        trace,
		}

		testErr := errors.New("xxx")
		var outInfo interface{}

		ctx := context.Background()
		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("get error", func() {
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, testErr)
			err := ossMock.UploadFile(ctx, &interfaces.Visitor{}, "", "", nil)
			assert.Equal(t, err, testErr)
		})

		h, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("%v", err)
		}

		tempMap := make(map[string]string)
		tempMap["sdad"] = "sdadxxx2313"
		tempData := mockOSSRequest{
			Method:  "POST",
			URL:     "http://" + h.Host,
			Headers: tempMap,
		}
		tempBytes, err1 := jsoniter.Marshal(tempData)
		assert.Equal(t, err1, nil)

		err1 = jsoniter.Unmarshal(tempBytes, &outInfo)
		assert.Equal(t, err1, nil)

		Convey("rawclient error", func() {
			bSuccess = false
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, nil)
			err := ossMock.UploadFile(ctx, &interfaces.Visitor{}, "", "", nil)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			bSuccess = true
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, nil)
			err := ossMock.UploadFile(ctx, &interfaces.Visitor{}, "", "", nil)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteFile(t *testing.T) {
	Convey("DeleteFile", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockDnHTTPClient(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		bSuccess := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}

			if bSuccess {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("success\n"))
			} else {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("error\n"))
			}
		}))
		defer ts.Close()

		ossMock := &ossGateWay{
			log:          common.NewLogger(),
			client:       c,
			rawClient:    httpclient.NewRawHTTPClient(),
			adminAddress: ts.URL,
			trace:        trace,
		}

		testErr := errors.New("xxx")
		var outInfo interface{}

		ctx := context.Background()
		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("get error", func() {
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, testErr)
			err := ossMock.DeleteFile(ctx, &interfaces.Visitor{}, "", "")
			assert.Equal(t, err, testErr)
		})

		h, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("%v", err)
		}

		tempMap := make(map[string]string)
		tempMap["sdad"] = "sdadxxx"
		tempData := mockOSSRequest{
			Method:  "POST",
			URL:     "http://" + h.Host,
			Headers: tempMap,
		}
		tempBytes, err1 := jsoniter.Marshal(tempData)
		assert.Equal(t, err1, nil)

		err1 = jsoniter.Unmarshal(tempBytes, &outInfo)
		assert.Equal(t, err1, nil)

		Convey("rawclient error", func() {
			bSuccess = false
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, nil)
			err := ossMock.DeleteFile(ctx, &interfaces.Visitor{}, "", "")
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			bSuccess = true
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, nil)
			err := ossMock.DeleteFile(ctx, &interfaces.Visitor{}, "", "")
			assert.Equal(t, err, nil)
		})
	})
}

func TestConvertOSSResponse(t *testing.T) {
	Convey("convertOSSResponse", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mock.NewMockTraceClient(ctrl)

		c := mock.NewMockDnHTTPClient(ctrl)
		ossMock := &ossGateWay{
			log:    common.NewLogger(),
			client: c,
			trace:  trace,
		}

		var outInfo interface{}
		tempMap := make(map[string]string)
		tempMap["sdad"] = "sdadxxx"
		tempData := mockOSSRequest{
			Method:  "POST",
			URL:     "http://",
			Headers: tempMap,
		}
		tempBytes, err1 := jsoniter.Marshal(tempData)
		assert.Equal(t, err1, nil)

		Convey("success1", func() {
			err1 = jsoniter.Unmarshal(tempBytes, &outInfo)
			assert.Equal(t, err1, nil)

			out := ossMock.convertOSSResponse(outInfo)
			assert.Equal(t, tempData.Method, out.method)
			assert.Equal(t, tempData.URL, out.url)
			assert.Equal(t, tempData.Headers["sdad"], out.headers["sdad"])
		})
	})
}

func TestGetLocalOSSInfo(t *testing.T) {
	Convey("TestGetLocalOSSInfo", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := mock.NewMockDnHTTPClient(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		ossMock := &ossGateWay{
			log:       common.NewLogger(),
			client:    c,
			rawClient: httpclient.NewRawHTTPClient(),
			trace:     trace,
		}

		testErr := errors.New("xxx")
		var outInfo interface{}

		ctx := context.Background()
		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("get error", func() {
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfo, testErr)
			_, err := ossMock.GetLocalEnabledOSSInfo(ctx, &interfaces.Visitor{})
			assert.Equal(t, err, testErr)
		})

		data := `[
			{
			"default": true,
			"enabled": true,
			"id": "EB2772D4196047168395633A79857325t",
			"name": "xxx"
			}
			]`
		var tempInterface interface{}
		err := jsoniter.Unmarshal([]byte(data), &tempInterface)
		assert.Equal(t, err, nil)

		Convey("success", func() {
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempInterface, nil)
			infos, err := ossMock.GetLocalEnabledOSSInfo(ctx, &interfaces.Visitor{})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(infos), 1)
			assert.Equal(t, infos[0].BDefault, true)
			assert.Equal(t, infos[0].ID, "EB2772D4196047168395633A79857325t")
		})
	})
}
