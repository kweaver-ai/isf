// Package persrec_db 提供个性化推荐相关的数据库操作功能
package persrec_db

import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/common/helpers/hlartrace"
	"AuditLog/gocommon/api"
)

// RepoBase 是数据库仓储层的基础结构体
// 包含了数据库操作所需的基本组件：
// - logger: 日志记录器
// - db: 数据库连接实例
// - arTracer: 链路追踪器
type RepoBase struct {
	logger api.Logger
	db     *sqlx.DB

	arTracer hlartrace.Tracer
}

// NewRepoBase 创建一个新的RepoBase实例
// 参数:
//   - logger: 日志记录器，用于记录操作日志
//   - db: 数据库连接实例，用于执行数据库操作
//   - arTracer: 链路追踪器，用于追踪请求链路
//
// 返回:
//   - *RepoBase: 返回RepoBase实例指针
func NewRepoBase(logger api.Logger, db *sqlx.DB, arTracer hlartrace.Tracer) *RepoBase {
	return &RepoBase{
		logger:   logger,
		db:       db,
		arTracer: arTracer,
	}
}
