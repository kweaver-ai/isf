package dbaccess

import (
	"fmt"
	"sync"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
)

type flowClean struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	fcOnce sync.Once
	fc     *flowClean
)

// NewFlowClean 创建FlowClean操作对象
func NewFlowClean() *flowClean {
	fcOnce.Do(func() {
		fc = &flowClean{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})

	return fc
}

// CleanExpiredRefresh 清理过期refresh time 有效期 单位秒
func (fc *flowClean) CleanExpiredRefresh(t int64) (err error) {
	cutoffTime := common.Now().Add(-time.Duration(t) * time.Second)
	var num int64 = 1
	for num > 0 {
		num, err = fc.cleanExpiredRefresh(cutoffTime)
		if err != nil {
			return
		}
		fc.logger.Debugf("dbaccess CleanExpiredRefresh end, num: %d", num)
	}

	return
}

// 分批删除过期refresh_token 信息
func (fc *flowClean) cleanExpiredRefresh(t time.Time) (num int64, err error) {
	// 一般limit不过万  先设置成10000
	nLimit := 10000
	dbName := common.GetDBName("hydra_v2")
	strSQL := "delete from %s.hydra_oauth2_refresh where requested_at < ? limit ?"
	strSQL = fmt.Sprintf(strSQL, dbName)
	result, err := fc.db.Exec(strSQL, t, nLimit)
	if err != nil {
		fc.logger.Errorf("cleanExpiredRefresh  exec failed, err: %v", err)
		return
	}

	// 获取受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fc.logger.Errorf("cleanExpiredRefresh get rowsAffected failed, err: %v", err)
		return
	}

	return rowsAffected, nil
}

// GetAllExpireFlowIDs 获取所有过期的flow_id
func (fc *flowClean) GetAllExpireFlowIDs(t int64) (ids []string, err error) {
	cutoffTime := common.Now().Add(-time.Duration(t) * time.Second)

	dbName := common.GetDBName("hydra_v2")
	strSQL := `
	SELECT a.login_challenge, IFNULL(b.challenge_id, ''), a.requested_at FROM %s.hydra_oauth2_flow AS a
	LEFT JOIN %s.hydra_oauth2_refresh AS b ON a.consent_challenge_id = b.challenge_id
	WHERE a.requested_at < ? AND b.challenge_id IS NULL `
	strSQL = fmt.Sprintf(strSQL, dbName, dbName)
	rows, err := fc.db.Query(strSQL, cutoffTime)
	if err != nil {
		fc.logger.Errorln(err, strSQL)
		return
	}

	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				a.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				fc.logger.Errorln(closeErr)
			}
		}
	}()

	ids = make([]string, 0)
	for rows.Next() {
		var loginChallenge, challengeID, requestedAt string
		err = rows.Scan(&loginChallenge, &challengeID, &requestedAt)
		if err != nil {
			return nil, err
		}
		ids = append(ids, loginChallenge)
	}
	return
}

// CleanFlow 删除过期flow信息
func (fc *flowClean) CleanFlow(ids []string) (err error) {
	// 分割ID 防止sql过长
	splitedIDs := SplitArray(ids)

	for _, ids := range splitedIDs {
		tmpErr := fc.cleanFlow(ids)
		if tmpErr != nil {
			return tmpErr
		}

		fc.logger.Debugf(fmt.Sprintf("dbaccess CleanFlow end , num :%v", len(ids)))
	}

	return nil
}

// CleanFlow 分批删除过期flow信息
func (fc *flowClean) cleanFlow(ids []string) (err error) {
	if len(ids) == 0 {
		return nil
	}

	set, argIDs := GetFindInSetSQL(ids)

	dbName := common.GetDBName("hydra_v2")
	strSQL := "delete from %s.hydra_oauth2_flow where login_challenge in ("
	strSQL += set
	strSQL += ")"
	strSQL = fmt.Sprintf(strSQL, dbName)
	_, err = fc.db.Exec(strSQL, argIDs...)
	if err != nil {
		fc.logger.Errorf("cleanFlow  exec failed, err: %v", err)
		return
	}

	return nil
}
