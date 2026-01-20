package persrec_db

import (
	"context"
	"sync"
	"time"

	"AuditLog/common/helpers"
	"AuditLog/common/helpers/dbhelper"
	persrecrepoi "AuditLog/interfaces/drivenadapter/idbaccess/persrec"
	"AuditLog/models/persistence/persrecpo"
)

var (
	svcConfigRepoOnce sync.Once
	svcConfigRepoImpl persrecrepoi.IPersSvcConfigRepo
)

// svcConfigRepo 服务配置仓库结构体，实现了IPersSvcConfigRepo接口
// 用于管理服务级别的键值对配置
// 继承了RepoBase的基础数据库操作能力
type svcConfigRepo struct {
	*RepoBase
}

// NewSvcConfigRepo 创建服务配置仓库实例
// 使用单例模式确保全局只有一个服务配置仓库实例
// 参数:
//   - base: 基础仓库实例，提供数据库连接等基础设施
//
// 返回:
//   - persrecrepoi.IPersSvcConfigRepo: 服务配置仓库接口实现
func NewSvcConfigRepo(base *RepoBase) persrecrepoi.IPersSvcConfigRepo {
	svcConfigRepoOnce.Do(func() {
		svcConfigRepoImpl = &svcConfigRepo{
			RepoBase: base,
		}
	})

	return svcConfigRepoImpl
}

// Set 设置服务配置项的键值对
// 如果配置项不存在则创建，如果存在则更新
// 参数:
//   - ctx: 上下文，用于链路追踪
//   - key: 配置项的键
//   - val: 配置项的值
//
// 返回:
//   - error: 操作过程中的错误信息
func (repo *svcConfigRepo) Set(ctx context.Context, key, val string) (err error) {
	ctx, span := repo.arTracer.AddInternalTrace(ctx)
	defer func() {
		repo.arTracer.TelemetrySpanEndIgnoreDBNotFound(span, err)
	}()

	// 1. 开启事务
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer helpers.TxRollbackOrCommit(tx, &err, repo.logger)

	// 2. 事务操作
	sr := dbhelper.TxSr(tx, repo.logger)
	po := &persrecpo.SvcConfigPo{}

	// 2.1 查询
	err = sr.FromPo(po).
		WhereEqual("f_key", key).
		FindOne(po)

	if helpers.IsSqlNotFound(err) {
		// 2.2 不存在则插入
		po.Key = key
		po.Value = val
		po.CreatedAt = time.Now().Unix()
		_, err = sr.FromPo(po).
			InsertStruct(po)
	} else if err == nil {
		// 2.3 存在则更新
		po.Value = val
		po.UpdatedAt = time.Now().Unix()
		_, err = sr.FromPo(po).
			WhereEqual("f_key", key).
			SetUpdateFields(
				[]string{"f_value", "f_updated_at"},
			).
			UpdateByStruct(po)
	}

	return
}

// Get 获取服务配置项的值
// 如果配置项不存在，返回空字符串和空错误
// 参数:
//   - ctx: 上下文，用于链路追踪
//   - key: 配置项的键
//
// 返回:
//   - val: 配置项的值，如果不存在则为空字符串
//   - error: 查询过程中的错误信息
func (repo *svcConfigRepo) Get(ctx context.Context, key string) (val string, err error) {
	_, span := repo.arTracer.AddInternalTrace(ctx)
	defer func() {
		repo.arTracer.TelemetrySpanEndIgnoreDBNotFound(span, err)
	}()

	po := &persrecpo.SvcConfigPo{}
	err = dbhelper.NewSQLRunner(repo.db, repo.logger).
		FromPo(po).
		WhereEqual("f_key", key).
		FindOne(po)

	if helpers.IsSqlNotFound(err) {
		err = nil
		return
	}

	if err != nil {
		return
	}

	val = po.Value

	return
}
