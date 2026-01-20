// Package dbaccess 数据访问层 conf_db
package dbaccess

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
)

type conf struct {
	db         *sqlx.DB
	dbTrace    *sqlx.DB
	trace      observable.Tracer
	logger     common.Logger
	confKeyMap map[interfaces.ConfigKey]string
}

var (
	cOnce sync.Once
	c     *conf
)

// NewConf 创建实名链接数据库对象
func NewConf() *conf {
	cOnce.Do(func() {
		c = &conf{
			db:      dbPool,
			dbTrace: dbTracePool,
			trace:   common.SvcARTrace,
			logger:  common.NewLogger(),
		}
		c.tableInfoMapInit()
	})
	return c
}

func (c *conf) tableInfoMapInit() {
	c.confKeyMap = map[interfaces.ConfigKey]string{
		interfaces.RememberFor:     "remember_for",
		interfaces.RememberVisible: "remember_visible",

		interfaces.EnableIDCardLogin:  "id_card_login_status",
		interfaces.EnablePWDLock:      "enable_pwd_lock",
		interfaces.EnableThirdPWDLock: "enable_third_pwd_lock",
		interfaces.PWDErrCnt:          "pwd_err_cnt",
		interfaces.PWDLockTime:        "pwd_lock_time",
		interfaces.VCodeConfig:        "vcode_login_config",
		interfaces.SMSExpiration:      "anonymous_sms_expiration",
		interfaces.LimitRedirectURI:   "oauth2_redirect_uris",

		interfaces.TriSystemStatus: "enable_tri_system_status",
	}
}

// GetConfig 获取认证配置
func (c *conf) GetConfig(configKeys map[interfaces.ConfigKey]bool) (cfg interfaces.Config, err error) {
	dbName := common.GetDBName("authentication")
	strSQL := "SELECT `f_key`,`f_value` FROM `" + dbName + "`.`t_conf` WHERE `f_key` in (%s)"

	return c.getConfig(context.Background(), configKeys, strSQL)
}

// GetConfigFromShareMgnt 获取认证配置
func (c *conf) GetConfigFromShareMgnt(ctx context.Context, configKeys map[interfaces.ConfigKey]bool) (cfg interfaces.Config, err error) {
	c.trace.SetClientSpanName("数据访问层-获取认证配置")
	newCtx, span := c.trace.AddClientTrace(ctx)
	defer func() { c.trace.TelemetrySpanEnd(span, err) }()

	dbName := common.GetDBName("sharemgnt_db")
	strSQL := "SELECT `f_key`,`f_value` FROM `" + dbName + "`.`t_sharemgnt_config` WHERE `f_key` in (%s)"

	return c.getConfig(newCtx, configKeys, strSQL)
}

func (c *conf) getConfig(ctx context.Context, configKeys map[interfaces.ConfigKey]bool, strSQL string) (cfg interfaces.Config, err error) {
	if len(configKeys) == 0 {
		return
	}

	groupsStr := make([]string, 0)
	configs := make([]interface{}, 0, len(configKeys))
	for k := range configKeys {
		configs = append(configs, c.confKeyMap[k])
		groupsStr = append(groupsStr, "?")
	}
	groupStr := strings.Join(groupsStr, ",")
	rows, err := c.dbTrace.QueryContext(ctx, fmt.Sprintf(strSQL, groupStr), configs...)
	if err != nil {
		c.logger.Errorln(err, strSQL, configs)
		return
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				c.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				c.logger.Errorln(closeErr)
			}
		}
	}()

	tmpKVMap := make(map[string]string)
	for rows.Next() {
		var fKey, fValue string
		err = rows.Scan(&fKey, &fValue)
		if err != nil {
			return
		}
		tmpKVMap[fKey] = fValue
	}

	return c.convertToConfig(configKeys, tmpKVMap)
}

func (c *conf) convertToConfig(configKeys map[interfaces.ConfigKey]bool, kvMap map[string]string) (cfg interfaces.Config, err error) {
	for k := range configKeys {
		switch k {
		case interfaces.RememberFor:
			var intV int
			intV, err = strconv.Atoi(kvMap[c.confKeyMap[interfaces.RememberFor]])
			if err != nil {
				return
			}
			cfg.RememberFor = intV
		case interfaces.RememberVisible:
			var boolV bool
			boolV, err = strconv.ParseBool(kvMap[c.confKeyMap[interfaces.RememberVisible]])
			if err != nil {
				return
			}
			cfg.RememberVisible = boolV
		case interfaces.EnableIDCardLogin:
			var boolV bool
			boolV, err = strconv.ParseBool(kvMap[c.confKeyMap[interfaces.EnableIDCardLogin]])
			if err != nil {
				return
			}
			cfg.EnableIDCardLogin = boolV
		case interfaces.EnablePWDLock:
			var boolV bool
			boolV, err = strconv.ParseBool(kvMap[c.confKeyMap[interfaces.EnablePWDLock]])
			if err != nil {
				return
			}
			cfg.EnablePWDLock = boolV
		case interfaces.EnableThirdPWDLock:
			var boolV bool
			boolV, err = strconv.ParseBool(kvMap[c.confKeyMap[interfaces.EnableThirdPWDLock]])
			if err != nil {
				return
			}
			cfg.EnableThirdPWDLock = boolV
		case interfaces.PWDErrCnt:
			var intV int
			intV, err = strconv.Atoi(kvMap[c.confKeyMap[interfaces.PWDErrCnt]])
			if err != nil {
				return
			}
			cfg.PWDErrCnt = intV
		case interfaces.PWDLockTime:
			var intV int
			intV, err = strconv.Atoi(kvMap[c.confKeyMap[interfaces.PWDLockTime]])
			if err != nil {
				return
			}
			cfg.PWDLockTime = intV
		case interfaces.VCodeConfig:
			var jsonObj interface{}
			err = jsoniter.Unmarshal([]byte(kvMap[c.confKeyMap[interfaces.VCodeConfig]]), &jsonObj)
			if err != nil {
				return
			}
			cfg.VCodeConfig.Enable = jsonObj.(map[string]interface{})["isEnable"].(bool)
			cfg.VCodeConfig.PWDErrCnt = int(jsonObj.(map[string]interface{})["passwdErrCnt"].(float64))
		case interfaces.SMSExpiration:
			var smsExpiration int
			smsExpiration, err = strconv.Atoi(kvMap[c.confKeyMap[interfaces.SMSExpiration]])
			if err != nil {
				c.logger.Errorln("failed to convert smsExpiration to int, err:", err)
				return
			}
			cfg.SMSExpiration = smsExpiration
		case interfaces.LimitRedirectURI:
			if kvMap[c.confKeyMap[interfaces.LimitRedirectURI]] == "" {
				cfg.LimitRedirectURI = map[string]bool{}
				continue
			}

			var jsonObj []interface{}
			err = jsoniter.Unmarshal([]byte(kvMap[c.confKeyMap[interfaces.LimitRedirectURI]]), &jsonObj)
			if err != nil {
				return
			}

			tempLimitRedirectURI := make(map[string]bool, 0)
			for _, obj := range jsonObj {
				tempLimitRedirectURI[obj.(string)] = true
			}
			cfg.LimitRedirectURI = tempLimitRedirectURI
		case interfaces.TriSystemStatus:
			var boolV bool
			boolV, err = strconv.ParseBool(kvMap[c.confKeyMap[interfaces.TriSystemStatus]])
			if err != nil {
				return
			}
			cfg.TriSystemStatus = boolV
		default:
			// 此项不应该被匹配，如果匹配到此项则代表遍历项存在杂项
			return cfg, errors.New("this error is unexpected")
		}
	}
	return
}

// SetConfig 设置认证配置
func (c *conf) SetConfig(configKeys map[interfaces.ConfigKey]bool, cfg interfaces.Config) (err error) {
	strSQL := "update t_conf set f_value = ? where f_key = ?"

	var tx *sql.Tx
	tx, err = c.db.Begin()
	if err != nil {
		return
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			err = tx.Commit()
			if err != nil {
				c.logger.Errorf("Transaction Commit Error:%v", err)
				return
			}
		default:
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				c.logger.Errorf("Transaction Rollback Error:%v", rollbackErr)
			}
		}
	}()

	var value interface{}
	for k := range configKeys {
		switch k {
		case interfaces.RememberFor:
			value = strconv.Itoa(cfg.RememberFor)
		case interfaces.RememberVisible:
			value = strconv.FormatBool(cfg.RememberVisible)
		case interfaces.SMSExpiration:
			value = strconv.Itoa(cfg.SMSExpiration)
		default:
			// 此项不应该被匹配，如果匹配到此项则代表遍历项存在杂项
			return errors.New("this error is unexpected")
		}
		_, err = tx.Exec(strSQL, value, c.confKeyMap[k])
		if err != nil {
			c.logger.Errorln(err, strSQL)
			return
		}
	}

	return
}
