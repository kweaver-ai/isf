package apiserver

import (
	"fmt"
	"net/http"

	"policy_mgnt/decision"
	"policy_mgnt/utils/models"

	"policy_mgnt/apiserver/jsonschema"

	"github.com/gin-gonic/gin"

	cerrors "policy_mgnt/utils/gocommon/v2/errors"
	chttp "policy_mgnt/utils/gocommon/v2/http"
	clog "policy_mgnt/utils/gocommon/v2/log"
)

type decisionHandler struct {
	pmgnt  decision.PolicyDecision
	logger clog.Logger
}

func newDecisionHandler() (*decisionHandler, error) {
	pmgnt := decision.NewPolicyDecision()
	logger := clog.NewLogger()

	return newDecisionHandlerWithMgnt(pmgnt, logger), nil
}

func newDecisionHandlerWithMgnt(pmgnt decision.PolicyDecision, logger clog.Logger) *decisionHandler {
	return &decisionHandler{
		pmgnt:  pmgnt,
		logger: logger,
	}
}

// RegisterPrivateAPI 内部API
func (d *decisionHandler) AddRouters(group *gin.RouterGroup) {
	group.POST("/network/allow", chttp.JsonSchemaValidationMiddleware(jsonschema.AcccessorNetworkSchema), d.netWorkDecision)
	group.GET("/sign_in_policy/client_restriction", d.clientSignDecision)
	group.GET("/opa/data/policy", d.readOPAData)
}

func (d *decisionHandler) netWorkDecision(c *gin.Context) {
	var decisionObjs models.Accessor
	err := chttp.ParseBody(c, &decisionObjs)
	if err != nil {
		chttp.ErrorResponse(c, fmt.Errorf("decision: ParseBody: %w", err))
		return
	}
	ok, err := d.pmgnt.NetWorkDecision(c.Request.Context(), decisionObjs)
	if err != nil {
		chttp.ErrorResponse(c, fmt.Errorf("decision: Decision: %w", err))
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"result": ok})
}

func (d *decisionHandler) clientSignDecision(c *gin.Context) {
	clientType, ok := c.GetQuery("client_type") // client_type为必传参数, 不传参数报错400
	if !ok {
		chttp.ErrorResponse(c, cerrors.ErrBadRequestPublic(&cerrors.ErrorInfo{Cause: "client_type is nil."}))
		return
	}
	res, err := d.pmgnt.ClientSignDecision(c.Request.Context(), clientType)
	if err != nil {
		chttp.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"result": res})
}

func (d *decisionHandler) readOPAData(c *gin.Context) {
	path, ok := c.GetQuery("path") // path为非必传参数, 不传参数默认全部
	if !ok {
		path = "/"
	}

	data, err := d.pmgnt.ReadOPAData(c.Request.Context(), path)
	if err != nil {
		d.logger.Errorln(err)
		chttp.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, data)
}
