package ticket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	ticketSchema "Authentication/driveradapters/jsonschema/ticket_schema"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"

	. "github.com/smartystreets/goconvey/convey"
)

func newRESTHandler(loTicket interfaces.LogicsTicket) RESTHandler {
	createTicketSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(ticketSchema.TicketSchemaStr))
	return &restHandler{
		ticket:             loTicket,
		createTicketSchema: createTicketSchema,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

//nolint:lll
func TestCreateTicket(t *testing.T) {
	Convey("TestCreateTicket", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loTicket := mock.NewMockLogicsTicket(ctrl)

		common.InitARTrace("test")
		testRestHandler := newRESTHandler(loTicket)
		testRestHandler.RegisterPublic(r)
		target := "/api/authentication/v1/ticket"

		Convey("failed, client_id is empty string", func() {
			reqBody := map[string]interface{}{
				"client_id":     "",
				"refresh_token": "clq2q/1bz0Sw0TF+zeZugkdTfoLY7CiNR6LRcGoZs18Sl+ZXEidp8jluhSAsglaGWNsL6JtCthuHstjwpP1ML1c/rxe5ml6FgPQdZV7KSqrerOmb5NoOF+wuEkXwGbU/dv2UkM33pg/WxwyRw25yu1w+AG7io2j9sbCbD9HZaowPscrAOVheBvlRd1RZsVfSzazWvYkegu3nCRyrIMpj/U/jcuktqM/qzpZB6C9BTu7dIuXXR7lqf01Qdo6YArjQdnLje0VPJt12MfXzR7OxFMzF9kWsXkSBgLuz8KMuW2aMwVlgzesjS89VbsnOaEBf912qoBtuGx0h4X5e34eQ3g==",
			}
			reqBodyByte, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBodyByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("failed, refresh_token is empty string", func() {
			reqBody := map[string]interface{}{
				"client_id":     "0c8839f4-894c-452c-ae96-911df5e04c64",
				"refresh_token": "",
			}
			reqBodyByte, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBodyByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("failed, unknown error", func() {
			tmpErr := fmt.Errorf("unknown error")
			reqBody := map[string]interface{}{
				"client_id":     "0c8839f4-894c-452c-ae96-911df5e04c64",
				"refresh_token": "clq2q/1bz0Sw0TF+zeZugkdTfoLY7CiNR6LRcGoZs18Sl+ZXEidp8jluhSAsglaGWNsL6JtCthuHstjwpP1ML1c/rxe5ml6FgPQdZV7KSqrerOmb5NoOF+wuEkXwGbU/dv2UkM33pg/WxwyRw25yu1w+AG7io2j9sbCbD9HZaowPscrAOVheBvlRd1RZsVfSzazWvYkegu3nCRyrIMpj/U/jcuktqM/qzpZB6C9BTu7dIuXXR7lqf01Qdo6YArjQdnLje0VPJt12MfXzR7OxFMzF9kWsXkSBgLuz8KMuW2aMwVlgzesjS89VbsnOaEBf912qoBtuGx0h4X5e34eQ3g==",
			}
			reqBodyByte, _ := json.Marshal(reqBody)
			loTicket.EXPECT().CreateTicket(gomock.Any(), gomock.Any(), gomock.Any()).Return("", tmpErr)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBodyByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			reqBody := map[string]interface{}{
				"client_id":     "0c8839f4-894c-452c-ae96-911df5e04c64",
				"refresh_token": "clq2q/1bz0Sw0TF+zeZugkdTfoLY7CiNR6LRcGoZs18Sl+ZXEidp8jluhSAsglaGWNsL6JtCthuHstjwpP1ML1c/rxe5ml6FgPQdZV7KSqrerOmb5NoOF+wuEkXwGbU/dv2UkM33pg/WxwyRw25yu1w+AG7io2j9sbCbD9HZaowPscrAOVheBvlRd1RZsVfSzazWvYkegu3nCRyrIMpj/U/jcuktqM/qzpZB6C9BTu7dIuXXR7lqf01Qdo6YArjQdnLje0VPJt12MfXzR7OxFMzF9kWsXkSBgLuz8KMuW2aMwVlgzesjS89VbsnOaEBf912qoBtuGx0h4X5e34eQ3g==",
			}
			reqBodyByte, _ := json.Marshal(reqBody)
			loTicket.EXPECT().CreateTicket(gomock.Any(), gomock.Any(), gomock.Any()).Return("01HVQWAXXCSA98YT2BY6Q0KKRS", nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBodyByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			resParamByte, _ := io.ReadAll(result.Body)

			var jsonV map[string]interface{}
			err := jsoniter.Unmarshal(resParamByte, &jsonV)
			assert.Equal(t, err, nil)
			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, jsonV["ticket"].(string), "01HVQWAXXCSA98YT2BY6Q0KKRS")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
