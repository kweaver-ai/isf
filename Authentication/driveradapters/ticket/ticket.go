// Package ticket 协议层
package ticket

import (
	"net/http"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"Authentication/common"
	ticketSchema "Authentication/driveradapters/jsonschema/ticket_schema"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	ticket "Authentication/logics/ticket"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	RegisterPublic(engine *gin.Engine)
}

type restHandler struct {
	ticket             interfaces.LogicsTicket
	createTicketSchema *gojsonschema.Schema
}

var (
	once sync.Once
	r    RESTHandler
)

// NewRESTHandler 创建ticket RESTHandler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		createTicketSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(ticketSchema.TicketSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		r = &restHandler{
			ticket:             ticket.NewTicket(),
			createTicketSchema: createTicketSchema,
		}
	})
	return r
}

// RegisterPublic 注册开放API
func (r *restHandler) RegisterPublic(engine *gin.Engine) {
	engine.POST("/api/authentication/v1/ticket", observable.MiddlewareTrace(common.SvcARTrace), r.createTicket)
}

func (r *restHandler) createTicket(c *gin.Context) {
	var err error
	var jsonV map[string]interface{}
	if err = util.ValidateAndBindGin(c, r.createTicketSchema, &jsonV); err != nil {
		rest.ReplyError(c, err)
		return
	}
	reqInfo := &interfaces.TicketReq{}
	reqInfo.ClientID = jsonV["client_id"].(string)
	reqInfo.RefreshToken = jsonV["refresh_token"].(string)

	visitor := &interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	ticketID, err := r.ticket.CreateTicket(c, visitor, reqInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := map[string]interface{}{
		"ticket": ticketID,
	}
	rest.ReplyOK(c, http.StatusOK, resInfo)
}
