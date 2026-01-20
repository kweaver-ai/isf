package interfaces

import "github.com/gin-gonic/gin"

//go:generate mockgen -package mock -source ../interfaces/driveradapters.go -destination ../interfaces/mock/mock_driveradapters.go
type PublicRESTHandler interface {
	// RegisterPublic 注册开放API
	RegisterPublic(oauthRouterGroup *gin.RouterGroup)
}

type PrivateRESTHandler interface {
	// RegisterPrivate 注册开放API
	RegisterPrivate(routerGroup *gin.RouterGroup)
}

type MQHandler interface {
	Subscribe()
}

type LogType int

const (
	LogType_Login      LogType = 1
	LogType_Management LogType = 2
	LogType_Operation  LogType = 3
)
