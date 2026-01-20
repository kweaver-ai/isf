// Package dbaccess 数据访问层
package dbaccess

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authorization/common"
	"Authorization/interfaces"
)

type obligation struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	obligationOnce    sync.Once
	obligationService *obligation
)

// NewObligationTemplate 创建数据库对象
func NewObligation() *obligation {
	obligationOnce.Do(func() {
		obligationService = &obligation{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})
	return obligationService
}

func (o *obligation) Add(ctx context.Context, info *interfaces.ObligationInfo) (err error) {
	curTime := common.GetCurrentMicrosecondTimestamp()
	// info.Config 转成string
	configByte, err := json.Marshal(info.Value)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	strSQL := "insert into " + common.GetDBName(databaseName) + ".t_obligation(f_id, f_type_id, f_name, f_description, f_value, f_created_at, f_modified_at) values(?,?,?,?,?,?,?)"
	_, err = o.db.Exec(strSQL, info.ID, info.TypeID, info.Name, info.Description, string(configByte), curTime, curTime)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	return
}

func (o *obligation) Update(ctx context.Context, obligationID, name string, nameChanged bool, description string, descriptionChanged bool, value any, valueChanged bool) (err error) {
	var args []any
	dbName := common.GetDBName(databaseName)
	sqlStr := "update %s.t_obligation set "
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	if nameChanged {
		args = append(args, name)
		sqlStr += "f_name = ?, "
	}

	if descriptionChanged {
		args = append(args, description)
		sqlStr += "f_description = ?, "
	}

	if valueChanged {
		configByte, err := json.Marshal(value)
		if err != nil {
			o.logger.Errorln(err)
			return err
		}
		args = append(args, string(configByte))
		sqlStr += "f_value = ?, "
	}

	curTime := common.GetCurrentMicrosecondTimestamp()
	sqlStr += "f_modified_at = ? where f_id = ? "
	args = append(args, curTime, obligationID)

	if _, err := o.db.Exec(sqlStr, args...); err != nil {
		o.logger.Errorf("err: %v, sqlStr: %s, args: %v", err, sqlStr, args)
		return err
	}
	return nil
}

func (o *obligation) Delete(ctx context.Context, obligationID string) (err error) {
	strSQL := "delete from " + common.GetDBName(databaseName) + ".t_obligation where f_id = ?"
	_, err = o.db.Exec(strSQL, obligationID)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	return
}

func (o *obligation) GetByID(ctx context.Context, obligationID string) (info interfaces.ObligationInfo, err error) {
	strSQL := "select f_id, f_type_id, f_name, f_description, f_value, f_created_at, f_modified_at from " + common.GetDBName(databaseName) + ".t_obligation where f_id = ?"
	// 帮我写一个逻辑 生成 sql 语句
	var args []any
	args = append(args, obligationID)
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var configStr string
		err = rows.Scan(&info.ID, &info.TypeID, &info.Name, &info.Description, &configStr, &info.CreatedTime, &info.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return
		}
		err = json.Unmarshal([]byte(configStr), &info.Value)
		if err != nil {
			o.logger.Errorln(err)
			return
		}
	}
	return
}

func (o *obligation) Get(ctx context.Context, info *interfaces.ObligationSearchInfo) (count int, resultInfos []interfaces.ObligationInfo, err error) {
	var countRows *sql.Rows
	countRows, err = o.db.Query("select count(1) from " + common.GetDBName(databaseName) + ".t_obligation")
	if err != nil {
		o.logger.Errorln(err)
		return 0, nil, err
	}

	for countRows.Next() {
		err = countRows.Scan(&count)
		if err != nil {
			o.logger.Errorln(err)
			return 0, nil, err
		}
	}

	if countRows != nil {
		if countRowsErr := countRows.Err(); countRowsErr != nil {
			o.logger.Errorln(countRowsErr)
		}
		if closeErr := countRows.Close(); closeErr != nil {
			o.logger.Errorln(closeErr)
		}
	}

	strSQL := "select f_id, f_type_id, f_name, f_description, f_value, f_created_at, f_modified_at from " + common.GetDBName(databaseName) + ".t_obligation limit ? offset ?"
	// 帮我写一个逻辑 生成 sql 语句
	var args []any
	args = append(args, info.Limit, info.Offset)
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return 0, nil, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var resultInfo interfaces.ObligationInfo
		var configStr string
		err = rows.Scan(&resultInfo.ID, &resultInfo.TypeID, &resultInfo.Name, &resultInfo.Description, &configStr, &resultInfo.CreatedTime, &resultInfo.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return 0, nil, err
		}
		err = json.Unmarshal([]byte(configStr), &resultInfo.Value)
		if err != nil {
			o.logger.Errorln(err)
			return 0, nil, err
		}
		resultInfos = append(resultInfos, resultInfo)
	}
	return
}

func (o *obligation) GetByObligationTypeIDs(ctx context.Context, obligationTypeIDMap map[string]bool) (resultInfos map[string][]interfaces.ObligationInfo, err error) {
	resultInfos = make(map[string][]interfaces.ObligationInfo)
	if len(obligationTypeIDMap) == 0 {
		return nil, nil
	}
	o.logger.Debugf("Query: obligationTypeIDMap %+v", obligationTypeIDMap)
	IDS := make([]string, 0, len(obligationTypeIDMap))
	for typeID := range obligationTypeIDMap {
		IDS = append(IDS, typeID)
	}
	typeIDSet, typeIDGroup := getFindInSetSQL(IDS)
	strSQL := "select f_id, f_type_id, f_name, f_description, f_value, f_created_at, f_modified_at from " + common.GetDBName(databaseName) +
		".t_obligation  where f_type_id in (" + typeIDSet + ")"
	var args []any
	args = append(args, typeIDGroup...)
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return nil, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var resultInfo interfaces.ObligationInfo
		var configStr string
		err = rows.Scan(&resultInfo.ID, &resultInfo.TypeID, &resultInfo.Name, &resultInfo.Description, &configStr, &resultInfo.CreatedTime, &resultInfo.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		err = json.Unmarshal([]byte(configStr), &resultInfo.Value)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		resultInfos[resultInfo.TypeID] = append(resultInfos[resultInfo.TypeID], resultInfo)
	}
	return
}

func (o *obligation) GetByIDs(ctx context.Context, ids []string) (resultInfos []interfaces.ObligationInfo, err error) {
	resultInfos = make([]interfaces.ObligationInfo, 0)
	if len(ids) == 0 {
		return nil, nil
	}
	IDSet, IDGroup := getFindInSetSQL(ids)
	strSQL := "select f_id, f_type_id, f_name, f_description, f_value, f_created_at, f_modified_at from " + common.GetDBName(databaseName) +
		".t_obligation  where f_id in (" + IDSet + ")"
	var args []any
	args = append(args, IDGroup...)
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return nil, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var resultInfo interfaces.ObligationInfo
		var configStr string
		err = rows.Scan(&resultInfo.ID, &resultInfo.TypeID, &resultInfo.Name, &resultInfo.Description, &configStr, &resultInfo.CreatedTime, &resultInfo.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		err = json.Unmarshal([]byte(configStr), &resultInfo.Value)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		resultInfos = append(resultInfos, resultInfo)
	}
	return
}
