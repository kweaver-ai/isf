package dbaccess

import (
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
)

type dbHydra struct {
	db                *sqlx.DB
	batchSize         int
	defaultTimeFormat string
	logger            common.Logger
}

// DistributedLock 业务类型
const (
	_ = iota
	// DisLockBussinessAssertionClear 清理失效的断言
	DisLockBussinessAssertionClear
)

var (
	hydraOnce sync.Once
	hydra     *dbHydra
)

// NewDBHydra 创建新的DBHydra对象
func NewDBHydra() *dbHydra {
	hydraOnce.Do(func() {
		hydra = &dbHydra{
			db:                dbPool,
			batchSize:         1000,
			defaultTimeFormat: "2006-01-02 15:04:05",
			logger:            common.NewLogger(),
		}
	})

	return hydra
}

func (h *dbHydra) DeleteExpiredAssertions() (affectedRows int64, err error) {
	hydraDB := common.GetDBName("hydra_v2")
	authDB := common.GetDBName("authentication")
	// 获取nid
	nid := ""
	strSQL := "select id from %s.networks limit 1"
	strSQL = fmt.Sprintf(strSQL, hydraDB)
	if err = h.db.QueryRow(strSQL).Scan(&nid); err != nil {
		return
	}

	// 获取分布式锁
	tx, err := h.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		switch err {
		case nil:
			// 提交事务
			err = tx.Commit()
			if err != nil {
				h.logger.Errorf("Transaction Commit Failed:%v", err)
			}
		default:
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				h.logger.Errorf("Transaction Rollback Error:%v", rollbackErr)
			}
		}
	}()
	getLockSQL := "select f_business_type from %s.t_distributed_lock where f_business_type = ? for update"
	getLockSQL = fmt.Sprintf(getLockSQL, authDB)
	rows, err := tx.Query(getLockSQL, DisLockBussinessAssertionClear)
	if err != nil {
		h.logger.Errorln(err)
		return 0, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				h.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				h.logger.Errorln(closeErr)
			}
		}
	}()

	// 清理失效的断言
	strSQL = "delete from %s.hydra_oauth2_jti_blacklist WHERE nid = ? AND expires_at < ? limit ?"
	strSQL = fmt.Sprintf(strSQL, hydraDB)
	timeStr := common.Now().UTC().Format(h.defaultTimeFormat)
	result, err := h.db.Exec(strSQL, nid, timeStr, h.batchSize)
	if err != nil {
		return 0, err
	}
	affectedRows, err = result.RowsAffected()
	if err != nil {
		h.logger.Errorf("RowsAffected() err=%v\n", err)
		return 0, nil
	}

	return affectedRows, nil
}
