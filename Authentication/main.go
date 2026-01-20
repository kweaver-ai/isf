package main

import (
	"fmt"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"Authentication/common"
	"Authentication/dbaccess"
	"Authentication/drivenadapters"
	da "Authentication/driveradapters"
	accesstokenperm "Authentication/driveradapters/access_token_perm"
	"Authentication/driveradapters/assertion"
	"Authentication/driveradapters/audit"
	"Authentication/driveradapters/conf"
	"Authentication/driveradapters/login"
	mq "Authentication/driveradapters/mq_handler"
	"Authentication/driveradapters/probe"
	"Authentication/driveradapters/register"
	"Authentication/driveradapters/session"
	"Authentication/driveradapters/sms"
	"Authentication/driveradapters/ticket"
	"Authentication/logics"
	flowclean "Authentication/logics/flow_clean"
)

// authentication 认证对象
type authentication struct {
	log                    common.Logger
	probeHandler           probe.RESTHandler
	loginHandler           login.RESTHandler
	sessionHandler         session.RESTHandler
	registerHandler        register.RESTHandler
	confHandler            conf.RESTHandler
	accessTokenPermHandler accesstokenperm.RESTHandler
	assertionHandler       assertion.RESTHandler
	mqHandler              mq.MQHandler
	smsHandler             sms.RESTHandler
	ticketHandler          ticket.RESTHandler
	auditHandler           audit.RESTHandler
}

// Start 启动服务
func (a *authentication) Start() {
	a.log.Infoln("start authentication server")

	// 定时清理context记录
	da.StartCleanThread()

	// 定时清理断言记录
	logics.StartCronThread()

	gin.SetMode(gin.ReleaseMode)

	// mq订阅
	a.mqHandler.Subscribe()

	go func() {
		engine := gin.New()
		engine.Use(gin.Recovery())
		engine.UseRawPath = true
		_ = engine.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})

		// 注册开放API
		a.loginHandler.RegisterPublic(engine)
		a.registerHandler.RegisterPublic(engine)
		a.confHandler.RegisterPublic(engine)
		a.assertionHandler.RegisterPublic(engine)
		a.accessTokenPermHandler.RegisterPublic(engine)
		a.smsHandler.RegisterPublic(engine)
		a.ticketHandler.RegisterPublic(engine)

		// 注册开放端口探针
		a.probeHandler.RegisterPublic(engine)

		if err := engine.Run(fmt.Sprintf("%s:%d", common.SvcConfig.SvcHost, common.SvcConfig.SvcPublicPort)); err != nil {
			a.log.Errorln(err)
		}
	}()

	go func() {
		engine := gin.New()
		engine.Use(gin.Recovery())
		engine.UseRawPath = true
		_ = engine.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})

		// 注册内部API
		a.sessionHandler.RegisterPrivate(engine)
		a.confHandler.RegisterPrivate(engine)
		a.accessTokenPermHandler.RegisterPrivate(engine)
		a.assertionHandler.RegisterPrivate(engine)
		a.loginHandler.RegisterPrivate(engine)
		a.auditHandler.RegisterPrivate(engine)

		// 注册内部端口探针
		a.probeHandler.RegisterPrivate(engine)

		if err := engine.Run(fmt.Sprintf("%s:%d", common.SvcConfig.SvcHost, common.SvcConfig.SvcPrivatePort)); err != nil {
			a.log.Errorln(err)
		}
	}()
}

func main() {
	// 配置注入
	common.InitConfig()

	// 设置错误码语言
	rest.Register(common.ErrorI18n)
	rest.SetLang(common.SvcConfig.Lang)

	// 初始化ARTrace实例
	common.InitARTrace("authentication")

	// 配置log等级
	svcLog := common.NewLogger()
	svcLog.SetLevel(common.SvcConfig.LogLevel)

	// dbPool 依赖注入
	dbPool := common.NewDBPool()
	dbaccess.SetDBPool(dbPool)

	// dbTracePool依赖注入
	dbTracePool := common.NewDBTracePool()
	logics.SetDBTracePool(dbTracePool)
	dbaccess.SetDBTracePool(dbTracePool)

	// dbaccess 依赖注入
	logics.SetDBRegister(dbaccess.NewRegister())
	logics.SetDBSession(dbaccess.NewSession())
	logics.SetDBConf(dbaccess.NewConf())
	logics.SetDBAccessTokenPerm(dbaccess.NewAccessTokenPerm())
	logics.SetDBLogin(dbaccess.NewLogin())
	logics.SetDBPool(dbPool)
	logics.SetDBOutbox(dbaccess.NewOutbox())
	logics.SetDBAnonymousSMS(dbaccess.NewAnonymousSMS())
	logics.SetDBTicket(dbaccess.NewTicket())
	logics.SetDBFlowClean(dbaccess.NewFlowClean())
	logics.SetDBUnorderedOutbox(dbaccess.NewUnorderedOutbox())

	// drivenadapters 依赖注入
	logics.SetDnHydraAdmin(drivenadapters.NewHydraAdmin())
	logics.SetDnHydraPublic(drivenadapters.NewHydraPublic())
	logics.SetDnUserManagement(drivenadapters.NewUserManagement())
	logics.SetDnEacp(drivenadapters.NewEacp())
	logics.SetDnShareMgnt(drivenadapters.NewShareMgnt())
	logics.SetDnEacpLog(drivenadapters.NewEacpLog())
	logics.SetDnMessageBroker(drivenadapters.NewMessageBroker())

	// 初始化清理flow实例
	_ = flowclean.NewFlowClean()

	server := &authentication{
		log:                    common.NewLogger(),
		sessionHandler:         session.NewRESTHandler(),
		loginHandler:           login.NewRESTHandler(),
		probeHandler:           probe.NewRESTHandler(),
		registerHandler:        register.NewRESTHandler(),
		confHandler:            conf.NewRESTHandler(),
		accessTokenPermHandler: accesstokenperm.NewRESTHandler(),
		assertionHandler:       assertion.NewRESTHandler(),
		mqHandler:              mq.NewMQHandler(),
		smsHandler:             sms.NewRESTHandler(),
		ticketHandler:          ticket.NewRESTHandler(),
		auditHandler:           audit.NewRESTHandler(),
	}

	server.Start()

	select {}
}
