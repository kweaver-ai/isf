package logics

import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-redis/redis/v8"

	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
)

var (
	logger               api.Logger                      = nil
	loginLogRepo         interfaces.LogRepo              = nil
	mgntLogRepo          interfaces.LogRepo              = nil
	operLogRepo          interfaces.LogRepo              = nil
	historyRepo          interfaces.HistoryRepo          = nil
	userMgntRepo         interfaces.UserMgntRepo         = nil
	shareMgntRepo        interfaces.ShareMgntRepo        = nil
	docCenterRepo        interfaces.DocCenterRepo        = nil
	ossGateway           interfaces.OssGatewayRepo       = nil
	logStrategyRepo      interfaces.LogStrategyRepo      = nil
	logScopeStrategyRepo interfaces.LogScopeStrategyRepo = nil
	redisClient          redis.Cmdable                   = nil
	mqClient             api.MQClient                    = nil
	dbOutbox             interfaces.DBOutbox             = nil
	dbPool               *sqlx.DB                        = nil
	tracer               api.Tracer                      = nil
	dlmLock              interfaces.DLM                  = nil
)

func SetLogger(i api.Logger) {
	logger = i
}

func SetLoginLogRepo(i interfaces.LogRepo) {
	loginLogRepo = i
}

func SetMgntLogRepo(i interfaces.LogRepo) {
	mgntLogRepo = i
}

func SetOperLogRepo(i interfaces.LogRepo) {
	operLogRepo = i
}

func SetHistoryRepo(i interfaces.HistoryRepo) {
	historyRepo = i
}

func SetUserMgntRepo(i interfaces.UserMgntRepo) {
	userMgntRepo = i
}

func SetShareMgntRepo(i interfaces.ShareMgntRepo) {
	shareMgntRepo = i
}

func SetDocCenterRepo(i interfaces.DocCenterRepo) {
	docCenterRepo = i
}

func SetOssGateway(i interfaces.OssGatewayRepo) {
	ossGateway = i
}

func SetLogStrategyRepo(i interfaces.LogStrategyRepo) {
	logStrategyRepo = i
}

func SetLogScopeStrategyRepo(i interfaces.LogScopeStrategyRepo) {
	logScopeStrategyRepo = i
}

func SetRedisClient(i redis.Cmdable) {
	redisClient = i
}

func SetMQClient(i api.MQClient) {
	mqClient = i
}

func SetDBOutbox(i interfaces.DBOutbox) {
	dbOutbox = i
}

func SetDBPool(i *sqlx.DB) {
	dbPool = i
}

func SetTracer(i api.Tracer) {
	tracer = i
}

func SetDLM(i interfaces.DLM) {
	dlmLock = i
}
