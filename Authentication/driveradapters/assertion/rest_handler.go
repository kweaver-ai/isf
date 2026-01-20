// Package assertion 协议层
package assertion

import (
	"net/http"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"Authentication/common"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	"Authentication/logics/assertion"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPublic 注册内外部API
	RegisterPublic(engine *gin.Engine)
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type restHandler struct {
	assertion interfaces.Assertion
	hydra     interfaces.Hydra
}

var (
	once sync.Once
	a    RESTHandler
)

// NewRESTHandler 创建context handler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		a = &restHandler{
			assertion: assertion.NewAssertion(),
			hydra:     util.NewHydra(),
		}
	})

	return a
}

// RegisterPublic 注册外部API
func (a *restHandler) RegisterPublic(engine *gin.Engine) {
	engine.GET("/api/authentication/v1/jwt", observable.MiddlewareTrace(common.SvcARTrace), a.getAssertionByUserID)
}

// RegisterPublic 注册内部API
func (a *restHandler) RegisterPrivate(engine *gin.Engine) {
	engine.POST("/api/authentication/v1/token-hook", a.tokenHook)
}

func (a *restHandler) getAssertionByUserID(c *gin.Context) {
	// token内省
	visitor, err := util.Verify(c, a.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	userID, hasUserID := c.GetQuery("user_id")
	if !hasUserID || userID == "" {
		err = rest.NewHTTPError("invalid user_id", rest.BadRequest, map[string]interface{}{"params": []string{0: "user_id"}})
		rest.ReplyError(c, err)
		return
	}

	res, err := a.assertion.GetAssertionByUserID(c, &visitor, userID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusOK, map[string]interface{}{"assertion": res})
}

func (a *restHandler) tokenHook(c *gin.Context) {
	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}
	jsonReq := jsonV.(map[string]interface{})

	var result map[string]interface{}
	assertionReq, ok := jsonReq["request"].(map[string]interface{})["payload"].(map[string]interface{})["assertion"]
	if ok {
		clientID, ok := jsonReq["request"].(map[string]interface{})["client_id"].(string)
		if !ok {
			err := rest.NewHTTPError("invalid request.client_id", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
		res, err := a.assertion.TokenHook(assertionReq.([]interface{})[0].(string), clientID)
		if err != nil {
			rest.ReplyError(c, err)
			return
		}
		session := map[string]interface{}{"access_token": res}
		result = map[string]interface{}{
			"session": session,
		}
		rest.ReplyOK(c, http.StatusOK, result)
	} else {
		rest.ReplyOK(c, http.StatusNoContent, nil)
	}
}
