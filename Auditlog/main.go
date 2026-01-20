package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/ar_log"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/ar_trace"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/resource"
	"github.com/gin-gonic/gin"

	"AuditLog/boot"
	"AuditLog/common"
	"AuditLog/common/conf"
	"AuditLog/common/constants"
	"AuditLog/common/constants/persconsts"
	"AuditLog/common/global"
	"AuditLog/common/helpers"
	"AuditLog/common/helpers/hlartrace"
	"AuditLog/common/helpers/redishelper"
	"AuditLog/common/utils"
	"AuditLog/drivenadapters"
	"AuditLog/drivenadapters/db"
	"AuditLog/drivenadapters/httpaccess/doccenter"
	"AuditLog/drivenadapters/httpaccess/ossgateway"
	"AuditLog/drivenadapters/httpaccess/usermgnt"
	"AuditLog/drivenadapters/redisaccess"
	"AuditLog/drivenadapters/thrift"
	"AuditLog/driveradapters"
	"AuditLog/driveradapters/api/middleware"
	private "AuditLog/driveradapters/api/private"
	public "AuditLog/driveradapters/api/public"
	mq "AuditLog/driveradapters/mq"
	oprlogmq "AuditLog/driveradapters/mq/operation_log"
	"AuditLog/gocommon/api"
	"AuditLog/infra"
	"AuditLog/interfaces"
	"AuditLog/logics"
)

// build info
var (
	branchName      string
	buildTime       string
	commitID        string
	goCommonVersion string
	buildInfo       *helpers.BuildInfo
)

type auditLog struct {
	healthHandler  interfaces.PrivateRESTHandler
	logHandler     interfaces.PrivateRESTHandler
	historyHandler interfaces.PrivateRESTHandler
	activeHandler  interfaces.PrivateRESTHandler

	healthPubHandler interfaces.PublicRESTHandler
	// oprLogPubHandler   interfaces.PublicRESTHandler
	logStrategyHandler interfaces.PublicRESTHandler
	historyLogHandler  interfaces.PublicRESTHandler
	activeLogHandler   interfaces.PublicRESTHandler
	activeLogV2Handler interfaces.PublicRESTHandler

	mqHandler       interfaces.MQHandler
	oprLogMqHandler interfaces.MQHandler

	privateServer *http.Server
	publicServer  *http.Server

	dlmLock interfaces.DLM
}

func (a *auditLog) Start() {
	// 1. 设置信号监听
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// 2. mq
	if !helpers.IsLocalDev() {
		a.mqHandler.Subscribe()
	}

	a.oprLogMqHandler.Subscribe()

	// 3. http api
	// 3.1 内部接口
	go a.privateAPIServe()

	// 3.2 外部接口
	go a.publicAPIServe()

	// 4 初始化审计日志报表
	// 初始化报表迁移至doc-center

	// 5 初始化历史日志转存
	go a.initDumpLog()

	// 6. 打印build信息
	go printBuildInfo()

	// 7. 等待关闭信号
	<-ctx.Done()
	log.Println("shutdown signal received")

	// 8. 执行优雅关闭
	a.gracefulShutdown()
}

func (a *auditLog) privateAPIServe() {
	server := gin.New()

	gin.SetMode(gin.ReleaseMode)
	server.Use(gin.Recovery())

	auditLogUrlPrefix := "/api/audit-log/v1"

	// 1 health check
	healthGroup := server.Group(auditLogUrlPrefix)
	a.healthHandler.RegisterPrivate(healthGroup)

	// 2 middleware
	mws := []gin.HandlerFunc{
		middleware.SetLangFromHEADER(),
	}

	// 3 audit-log group
	group := server.Group(auditLogUrlPrefix)
	group.Use(mws...)

	a.logHandler.RegisterPrivate(group)
	a.historyHandler.RegisterPrivate(group)
	a.activeHandler.RegisterPrivate(group)

	// 4 个性化
	persGroup := server.Group(fmt.Sprintf("/api/%s/v1", persconsts.PersSvcName))
	persGroup.Use(mws...)

	// 5 run server
	a.privateServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", common.SvcConfig.ServicePrivatePort),
		Handler: server.Handler(),
	}

	go func() {
		if err := a.privateServer.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (a *auditLog) publicAPIServe() {
	// 1. new server
	server := gin.New()
	gin.SetMode(gin.ReleaseMode)

	if helpers.IsLocalDev() {
		gin.SetMode(gin.DebugMode)
	}

	server.Use(gin.Recovery())

	auditLogUrlPrefix := "/api/audit-log/v1"

	// 1.1 health check
	healthGroup := server.Group(auditLogUrlPrefix)
	a.healthPubHandler.RegisterPublic(healthGroup)

	// 1.2 active log v2
	activeLogV2Group := server.Group(auditLogUrlPrefix)
	a.activeLogV2Handler.RegisterPublic(activeLogV2Group)

	// 2. middleware
	trace := hlartrace.NewARTrace()
	mws := []gin.HandlerFunc{
		api.MiddlewareTrace(trace),
		middleware.SetLangFromHEADER(),
	}

	if !helpers.IsLocalDev() {
		o := api.NewOAuth2()
		mws = append(mws, api.Oauth2Middleware(o, trace))
	}

	mws = append(mws, gin.Logger())
	mws = append(mws, middleware.SetUserIDToCtx())
	mws = append(mws, middleware.SetUserTokenToCtx())
	mws = append(mws, middleware.SetVisitorUserToCtx())

	// 4. audit-log group
	group := server.Group(auditLogUrlPrefix)
	group.Use(mws...)

	// 4.1 register public api
	// a.oprLogPubHandler.RegisterPublic(group)
	a.logStrategyHandler.RegisterPublic(group)
	a.historyLogHandler.RegisterPublic(group)
	a.activeLogHandler.RegisterPublic(group)

	// 5. 个性化 group
	persGroup := server.Group(fmt.Sprintf("/api/%s/v1", persconsts.PersSvcName))
	persGroup.Use(mws...)

	// 6. run server
	a.publicServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", common.SvcConfig.ServicePublicPort),
		Handler: server.Handler(),
	}
	go func() {
		if err := a.publicServer.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (a *auditLog) gracefulShutdown() {
	// 1. 创建关闭超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 解锁DLM
	if a.dlmLock != nil {
		if err := a.dlmLock.UnLock(common.DeptPersFeatureCronLock); err != nil {
			log.Printf("[%s] perFeature cron dlm unlock error: %v", common.DeptPersFeatureCronLock, err)
		}

		if err := a.dlmLock.UnLock(common.DeptPersFeatureLock); err != nil {
			log.Printf("[%s] dept perFeature dlm unlock error: %v", common.DeptPersFeatureLock, err)
		}

		if err := a.dlmLock.UnLock(constants.DumpLogLockKey); err != nil {
			log.Printf("[%s] dump log dlm unlock error: %v", constants.DumpLogLockKey, err)
		}

		log.Println("shutting down dlm...")
	}

	// 2. 关闭 HTTP 服务器
	if a.privateServer != nil {
		log.Println("shutting down private server...")

		if err := a.privateServer.Shutdown(ctx); err != nil {
			log.Printf("private server shutdown error: %v", err)
		}
	}

	if a.publicServer != nil {
		log.Println("shutting down public server...")

		if err := a.publicServer.Shutdown(ctx); err != nil {
			log.Printf("public server shutdown error: %v", err)
		}
	}
}

// 初始化历史日志转存
func (a *auditLog) initDumpLog() {
	dumpLog := logics.NewDumpLog()
	dumpLog.InitDumpLog(context.Background())
}

func main() {
	defer ar_trace.ShutdownTracer()

	// 初始化配置
	conf.InitJsonConf()

	global.BuildInfo = &helpers.BuildInfo{
		BranchName:      branchName,
		BuildTime:       buildTime,
		CommitID:        commitID,
		GoCommonVersion: goCommonVersion,
		Other:           map[string]string{},
	}

	boot.Init()

	// 初始化AR日志
	serviceName := utils.GetEnv(constants.SvcNameEnvKey, "audit-log")
	global.SvcName = serviceName

	resource.SetServiceVersion(helpers.BuildVersion(getBuildInfo()))

	if helpers.IsLocalDev() {
		ar_log.InitLogger("yaml", "ob-app-config-log", serviceName)
	} else {
		ar_log.InitLogger("cm", "anyshare-telemetry-sdk", serviceName)
	}

	defer func() {
		if ar_log.Logger != nil {
			ar_log.Logger.Close()
		}
	}()

	// 1. 注入drivenadapter
	// 1.1 logger
	drivenadapters.SetLogger(common.SvcConfig.Logger)

	//1.2 db
	//dbConfig := &sqlx.DBConfig{}
	//_ = common.Configure(dbConfig, "dbrw.yaml")
	//dbPool, _ := sqlx.NewDB(dbConfig)
	dbPool := infra.NewDBPool()
	drivenadapters.SetDBPool(dbPool)

	// 1.3 http client
	HTTPClient := api.NewHttpClient()
	drivenadapters.SetHTTPClient(HTTPClient)

	// 2. 注入logics
	// 2.1 logger
	logics.SetLogger(common.SvcConfig.Logger)

	// 2.2 logRepo
	loginLogRepo := db.NewLoginLog()
	mgntLogRepo := db.NewManagementLog()
	operLogRepo := db.NewOperationLog()
	historyRepo := db.NewHisotryLog()
	logStrategyRepo := db.NewLogStrategy()
	logScopeStrategyRepo := db.NewScopeStrategy()
	userMgntRepo := usermgnt.NewUserMgnt()
	shareMgntRepo := thrift.NewShareMgnt()
	docCenter := doccenter.NewDocCenter()
	ossGateway := ossgateway.NewOssGateway()
	// 2.7 tracer
	tracer := api.NewARTrace()
	logics.SetTracer(tracer)

	logics.SetLoginLogRepo(loginLogRepo)
	logics.SetMgntLogRepo(mgntLogRepo)
	logics.SetOperLogRepo(operLogRepo)
	logics.SetHistoryRepo(historyRepo)
	logics.SetLogStrategyRepo(logStrategyRepo)
	logics.SetLogScopeStrategyRepo(logScopeStrategyRepo)
	logics.SetUserMgntRepo(userMgntRepo)
	logics.SetShareMgntRepo(shareMgntRepo)
	logics.SetDocCenterRepo(docCenter)
	logics.SetOssGateway(ossGateway)

	// 2.3 db pool
	logics.SetDBPool(dbPool)

	// 2.4 db outbox
	dbOutbox := db.NewOutbox()
	logics.SetDBOutbox(dbOutbox)

	// 2.5 mq client
	mqClient := api.NewMQClient()
	logics.SetMQClient(mqClient)

	// 2.6 redis client
	redisClient := redishelper.GetRedisClient()
	drivenadapters.SetRedisClient(redisClient)
	logics.SetRedisClient(redisClient)

	// 2.8 dlm
	dlm := redisaccess.NewDLM()
	logics.SetDLM(dlm)

	// 3. 启动服务
	a := &auditLog{
		healthHandler:  private.NewHealthHandler(),
		logHandler:     private.NewLogHandler(),
		historyHandler: private.NewHistoryHandler(),
		activeHandler:  private.NewActiveHandler(),

		healthPubHandler: public.NewHealthHandler(),
		// oprLogPubHandler:   public.NewOperationLogHandler(),
		logStrategyHandler: public.NewLogStrategyHandler(),
		historyLogHandler:  public.NewHistoryLogHandler(),
		activeLogHandler:   public.NewActiveLogHandler(),
		activeLogV2Handler: driveradapters.NewActiveLogV2Handler(),

		mqHandler:       mq.NewMQHandler(),
		oprLogMqHandler: oprlogmq.NewOprLogMqHandler(),
		dlmLock:         dlm,
	}

	a.Start()
}

func printBuildInfo() {
	//nolint:gomnd
	helpers.PrintBuildInfo(time.Millisecond*300, getBuildInfo())
}

func getBuildInfo() *helpers.BuildInfo {
	if buildInfo == nil {
		buildInfo = global.BuildInfo

		if helpers.IsLocalDev() {
			buildInfo = &helpers.BuildInfo{
				BranchName:      "mock_branch",
				BuildTime:       "mock_build_time",
				CommitID:        "mock_commit_id",
				GoCommonVersion: "mock_go_common_version",
				Other: map[string]string{
					"mock_key": "mock_value",
				},
			}
		}
	}

	return buildInfo
}
