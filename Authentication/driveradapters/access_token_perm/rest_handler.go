// Package accesstokenperm 协议层
package accesstokenperm

import (
	"net/http"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"

	"Authentication/common"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	accesstokenperm "Authentication/logics/access_token_perm"

	"github.com/gin-gonic/gin"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
	// RegisterPublic 注册外部API
	RegisterPublic(engine *gin.Engine)
}

type restHandler struct {
	accessTokenPerm interfaces.AccessTokenPerm
	hydra           interfaces.Hydra
}

var (
	once sync.Once
	a    RESTHandler
)

// NewRESTHandler 创建context handler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		a = &restHandler{
			accessTokenPerm: accesstokenperm.NewAccessTokenPerm(),
			hydra:           util.NewHydra(),
		}
	})

	return a
}

// RegisterPublic 注册内部API
func (a *restHandler) RegisterPrivate(engine *gin.Engine) {
	engine.PUT("/api/authentication/v1/access-token-perm/app/:app_id", observable.MiddlewareTrace(common.SvcARTrace), a.setAppAccessTokenPermPvt)
	engine.DELETE("/api/authentication/v1/access-token-perm/app/:app_id", observable.MiddlewareTrace(common.SvcARTrace), a.deleteAppAccessTokenPermPvt)
	engine.GET("/api/authentication/v1/access-token-perm/app", observable.MiddlewareTrace(common.SvcARTrace), a.getAllAppAccessTokenPermPvt)
}

// RegisterPublic 注册外部API
func (a *restHandler) RegisterPublic(engine *gin.Engine) {
	engine.PUT("/api/authentication/v1/access-token-perm/app/:app_id", a.setAppAccessTokenPermPub)
	engine.DELETE("/api/authentication/v1/access-token-perm/app/:app_id", a.deleteAppAccessTokenPermPub)
	engine.GET("/api/authentication/v1/access-token-perm/app", a.getAllAppAccessTokenPermPub)
}

func (a *restHandler) setAppAccessTokenPermPvt(c *gin.Context) {
	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	a.setAppAccessTokenPerm(c, &visitor)
}

func (a *restHandler) setAppAccessTokenPermPub(c *gin.Context) {
	// token内省
	visitor, err := util.Verify(c, a.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	a.setAppAccessTokenPerm(c, &visitor)
}

func (a *restHandler) setAppAccessTokenPerm(c *gin.Context, visitor *interfaces.Visitor) {
	appID := c.Param("app_id")
	err := a.accessTokenPerm.SetAppAccessTokenPerm(c, visitor, appID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (a *restHandler) deleteAppAccessTokenPermPvt(c *gin.Context) {
	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	a.deleteAppAccessTokenPerm(c, &visitor)
}

func (a *restHandler) deleteAppAccessTokenPermPub(c *gin.Context) {
	// token内省
	visitor, err := util.Verify(c, a.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	a.deleteAppAccessTokenPerm(c, &visitor)
}

func (a *restHandler) deleteAppAccessTokenPerm(c *gin.Context, visitor *interfaces.Visitor) {
	appID := c.Param("app_id")
	err := a.accessTokenPerm.DeleteAppAccessTokenPerm(c, visitor, appID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (a *restHandler) getAllAppAccessTokenPermPvt(c *gin.Context) {
	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	a.getAllAppAccessTokenPerm(c, &visitor)
}

func (a *restHandler) getAllAppAccessTokenPermPub(c *gin.Context) {
	// token内省
	visitor, err := util.Verify(c, a.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	a.getAllAppAccessTokenPerm(c, &visitor)
}

func (a *restHandler) getAllAppAccessTokenPerm(c *gin.Context, visitor *interfaces.Visitor) {
	permApps, err := a.accessTokenPerm.GetAllAppAccessTokenPerm(c, visitor)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusOK, permApps)
}
