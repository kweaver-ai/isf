package logics

import (
	"database/sql"
	"regexp"
	"sort"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/mitchellh/mapstructure"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"

	gerrors "github.com/kweaver-ai/go-lib/error"
)

var (
	conOnce sync.Once
	con     *config
)

type config struct {
	db      interfaces.DBConfig
	pool    *sqlx.DB
	logger  common.Logger
	role    interfaces.LogicsRole
	eacpLog interfaces.DrivenEacpLog
	ob      interfaces.LogicsOutbox
}

// NewConfig 创建新的Config对象
func NewConfig() *config {
	conOnce.Do(func() {
		con = &config{
			db:      dbConfig,
			pool:    dbPool,
			logger:  common.NewLogger(),
			role:    NewRole(),
			eacpLog: dnEacpLog,
			ob:      NewOutbox(OutboxBusinessConfig),
		}

		con.ob.RegisterHandlers(outboxDefaultPWDModifiedLog, con.sendDefaultPWDModifiedAuditLog)
		con.ob.RegisterHandlers(outboxCSFLevelEnumInitedLog, con.sendCSFLevelEnumInitedAuditLog)
		con.ob.RegisterHandlers(outboxCSFLevelEnum2InitedLog, con.sendCSFLevelEnum2InitedAuditLog)
	})
	return con
}

// UpdateConfig 设置配置信息
func (con *config) UpdateConfig(visitor *interfaces.Visitor, rg map[interfaces.ConfigKey]bool, config *interfaces.Config) (err error) {
	bDefaultPwd := rg[interfaces.UserDefaultPWD]
	bCSFLevelEnum := rg[interfaces.CSFLevelEnum]
	bCSFLevelEnum2 := rg[interfaces.CSFLevel2Enum]
	err = con.checkAuth(visitor, bDefaultPwd, bCSFLevelEnum, bCSFLevelEnum2)
	if err != nil {
		return
	}

	// 默认密码需大于等于6位，小于等于100, 支持数字、字符和特定特殊字符
	if bDefaultPwd {
		if visitor != nil && visitor.ID != "" {
			config.UserDefaultPWD, err = decodeRSA(config.UserDefaultPWD, RSA2048)
			if err != nil {
				return err
			}
		}
		out := con.checkIsValidDefaultPWD(config.UserDefaultPWD)
		if !out {
			err = rest.NewHTTPError("default password formt error", rest.BadRequest, nil)
			return
		}
	}

	// 检查密级枚举是否合法
	err = con.checkCSFLevelEnum(config, bCSFLevelEnum, bCSFLevelEnum2)
	if err != nil {
		return
	}

	// 获取事务处理器
	tx, err := con.pool.Begin()
	if err != nil {
		return
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				con.logger.Errorf("UpdateConfig Transaction Commit Error:%v", err)
				return
			}

			// 记录修改用户初始密码成功日志
			if visitor != nil && visitor.ID != "" && (bDefaultPwd || bCSFLevelEnum || bCSFLevelEnum2) {
				con.ob.NotifyPushOutboxThread()
			}

		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				con.logger.Errorf("UpdateConfig Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 密码配置
	err = con.handlePwdInfo(rg, config)
	if err != nil {
		return
	}

	// 更新配置
	err = con.updateConfig(rg, config, tx)
	if err != nil {
		return
	}

	// 添加审计日志
	if visitor != nil && visitor.ID != "" {
		err = con.sendAuditLog(visitor, bDefaultPwd, bCSFLevelEnum, bCSFLevelEnum2, config, tx)
		if err != nil {
			return
		}
	}

	return
}

// updateConfig 更新配置
func (con *config) updateConfig(rg map[interfaces.ConfigKey]bool, config *interfaces.Config, tx *sql.Tx) (err error) {
	if rg[interfaces.UserDefaultPWD] {
		tempMap := make(map[interfaces.ConfigKey]bool)
		tempMap[interfaces.UserDefaultNTLMPWD] = true
		tempMap[interfaces.UserDefaultSha2PWD] = true
		tempMap[interfaces.UserDefaultDESPWD] = true
		tempMap[interfaces.UserDefaultMd5PWD] = true

		tempConfig := interfaces.Config{
			UserDefaultNTLMPWD: config.UserDefaultNTLMPWD,
			UserDefaultSha2PWD: config.UserDefaultSha2PWD,
			UserDefaultDESPWD:  config.UserDefaultDESPWD,
			UserDefaultMd5PWD:  config.UserDefaultMd5PWD,
		}
		err = con.db.SetConfig(tempMap, &tempConfig, tx)
		if err != nil {
			return
		}
	}
	if rg[interfaces.CSFLevelEnum] {
		rg[interfaces.CSFLevelEnum] = true
		err = con.db.SetShareMgntConfig(interfaces.CSFLevelEnum, &interfaces.Config{CSFLevelEnum: config.CSFLevelEnum}, tx)
		if err != nil {
			return
		}
	}
	if rg[interfaces.CSFLevel2Enum] {
		err = con.db.SetShareMgntConfig(interfaces.CSFLevel2Enum, &interfaces.Config{CSFLevel2Enum: config.CSFLevel2Enum}, tx)
		if err != nil {
			return
		}
	}

	return nil
}

// sendAuditLog 发送审计日志
func (con *config) sendAuditLog(visitor *interfaces.Visitor, bDefaultPwd, bCSFLevelEnum, bCSFLevelEnum2 bool, config *interfaces.Config, tx *sql.Tx) (err error) {
	if bDefaultPwd {
		contentJSON := make(map[string]interface{})
		contentJSON["visitor"] = *visitor
		err = con.ob.AddOutboxInfo(outboxDefaultPWDModifiedLog, contentJSON, tx)
		if err != nil {
			return
		}
	}

	if bCSFLevelEnum {
		contentJSON := make(map[string]interface{})
		contentJSON["visitor"] = *visitor
		contentJSON["csf_level_enum"] = con.enumToStringSlice(config.CSFLevelEnum)
		err = con.ob.AddOutboxInfo(outboxCSFLevelEnumInitedLog, contentJSON, tx)
		if err != nil {
			return
		}
	}

	if bCSFLevelEnum2 {
		contentJSON := make(map[string]interface{})
		contentJSON["visitor"] = *visitor
		contentJSON["csf_level_enum2"] = con.enumToStringSlice(config.CSFLevel2Enum)
		err = con.ob.AddOutboxInfo(outboxCSFLevelEnum2InitedLog, contentJSON, tx)
		if err != nil {
			return
		}
	}
	return nil
}

// checkCSFLevelEnum 检查密级枚举是否合法
func (con *config) checkCSFLevelEnum(config *interfaces.Config, bCSFLevelEnum, bCSFLevelEnum2 bool) (err error) {
	// 密级枚举需检查是否为空
	var currentConfig interfaces.Config
	if bCSFLevelEnum && len(config.CSFLevelEnum) == 0 {
		err = gerrors.NewError(gerrors.PublicBadRequest, "csf level enum is empty")
		return
	}
	if bCSFLevelEnum2 && len(config.CSFLevel2Enum) == 0 {
		err = gerrors.NewError(gerrors.PublicBadRequest, "csf level2 enum is empty")
		return
	}

	// 如果已经设置密级枚举，则不能再次设置
	if bCSFLevelEnum || bCSFLevelEnum2 {
		currentConfig, err = con.GetConfig(map[interfaces.ConfigKey]bool{interfaces.CSFLevelEnum: true, interfaces.CSFLevel2Enum: true})
		if err != nil {
			return
		}
	}

	if bCSFLevelEnum && len(currentConfig.CSFLevelEnum) > 0 {
		err = gerrors.NewError(gerrors.PublicConflict, "csf level enum already set")
		return
	}
	if bCSFLevelEnum2 && len(currentConfig.CSFLevel2Enum) > 0 {
		err = gerrors.NewError(gerrors.PublicConflict, "csf level2 enum already set")
		return
	}

	return nil
}

func (con *config) enumToStringSlice(enum map[string]int) []string {
	// 按照int的大小 从小到大形成name数组
	names := make([]string, 0)
	for k := range enum {
		names = append(names, k)
	}
	sort.Slice(names, func(i, j int) bool {
		return enum[names[i]] < enum[names[j]]
	})
	return names
}

// checkAuth 检查权限
func (con *config) checkAuth(visitor *interfaces.Visitor, bDefaultPwd, bCSFLevelEnum, bCSFLevelEnum2 bool) (err error) {
	if visitor != nil && visitor.ID != "" {
		// 获取用户角色信息
		var roleIDs map[interfaces.Role]bool
		roleIDs, err = getRolesByUserID(con.role, visitor.ID)
		if err != nil {
			return err
		}

		// 默认密码配置权限检查, 只支持超级管理员和安全管理员
		if bDefaultPwd &&
			!roleIDs[interfaces.SystemRoleSuperAdmin] &&
			!roleIDs[interfaces.SystemRoleSecAdmin] {
			return gerrors.NewError(gerrors.PublicForbidden, "this user do not has the authority")
		}

		// 密级枚举配置权限检查, 支持四种管理员
		if (bCSFLevelEnum || bCSFLevelEnum2) &&
			(!roleIDs[interfaces.SystemRoleSuperAdmin] &&
				!roleIDs[interfaces.SystemRoleSysAdmin] &&
				!roleIDs[interfaces.SystemRoleSecAdmin] &&
				!roleIDs[interfaces.SystemRoleAuditAdmin]) {
			return gerrors.NewError(gerrors.PublicForbidden, "this user do not has the authority")
		}

		return nil
	}
	return nil
}

// handlePwdInfo 处理密码
func (con *config) handlePwdInfo(rg map[interfaces.ConfigKey]bool, config *interfaces.Config) (err error) {
	// 密码配置
	if _, ok := rg[interfaces.UserDefaultPWD]; ok {
		// 密码加密
		config.UserDefaultNTLMPWD = encodeNtlm(config.UserDefaultPWD)
		config.UserDefaultSha2PWD = encodeSha2(config.UserDefaultPWD)

		var tempCode string
		tempCode, err = encodeDes(config.UserDefaultPWD, PKCS5Padding)
		if err != nil {
			return
		}
		config.UserDefaultDESPWD = tempCode

		tempCode, err = encodeMD5(config.UserDefaultPWD)
		if err != nil {
			return
		}
		config.UserDefaultMd5PWD = tempCode

		rg[interfaces.UserDefaultDESPWD] = true
		rg[interfaces.UserDefaultMd5PWD] = true
		rg[interfaces.UserDefaultNTLMPWD] = true
		rg[interfaces.UserDefaultSha2PWD] = true
	}
	return
}

func (con *config) sendDefaultPWDModifiedAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		con.logger.Errorf("sendDefaultPWDModifiedAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = con.eacpLog.OpSetDefaultPWDLog(&v)
	if err != nil {
		con.logger.Errorf("sendDefaultPWDModifiedAuditLog err:%v", err)
	}
	return err
}

// sendCSFLevelEnumInitedAuditLog 发送密级枚举初始化审计日志
func (con *config) sendCSFLevelEnumInitedAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		con.logger.Errorf("sendCSFLevelEnumInitedAuditLog mapstructure.Decode err:%v", err)
		return
	}

	csfLevelEnum := make([]string, 0)
	if data, ok := info["csf_level_enum"].([]interface{}); ok {
		for _, v := range data {
			csfLevelEnum = append(csfLevelEnum, v.(string))
		}
	} else {
		con.logger.Errorf("sendCSFLevelEnumInitedAuditLog csf_level_enum is not a string array")
		return
	}

	err = con.eacpLog.OpSetCSFLevelEnumLog(&v, csfLevelEnum)
	if err != nil {
		con.logger.Errorf("sendCSFLevelEnumInitedAuditLog err:%v", err)
	}
	return err
}

// sendCSFLevelEnum2InitedAuditLog 发送密级2枚举初始化审计日志
func (con *config) sendCSFLevelEnum2InitedAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		con.logger.Errorf("sendDefaultPWDModifiedAuditLog mapstructure.Decode err:%v", err)
		return
	}

	csfLevel2Enum := make([]string, 0)
	if data, ok := info["csf_level_enum2"].([]interface{}); ok {
		for _, v := range data {
			csfLevel2Enum = append(csfLevel2Enum, v.(string))
		}
	}

	err = con.eacpLog.OpSetCSFLevel2EnumLog(&v, csfLevel2Enum)
	if err != nil {
		con.logger.Errorf("sendCSFLevelEnum2InitedAuditLog err:%v", err)
	}
	return err
}

// GetPWDRetrievalConfig 获取密码找回配置信息
func (con *config) GetConfig(rg map[interfaces.ConfigKey]bool) (config interfaces.Config, err error) {
	return con.db.GetConfig(rg)
}

// GetConfigFromOption 获取配置信息
func (con *config) GetConfigFromOption(rg map[interfaces.ConfigKey]bool) (config interfaces.Config, err error) {
	return con.db.GetConfigFromOption(rg)
}

// CheckDefaultPWD 检查用户初始密码格式
func (con *config) CheckDefaultPWD(visitor *interfaces.Visitor, pwd string) (result bool, msg string, err error) {
	// 检查权限，只支持超级管理员和安全管理员
	err = con.checkDefaultPWDAuthority(visitor.ID)
	if err != nil {
		return false, "", err
	}

	// 密码解密
	strPwd, err := decodeRSA(pwd, RSA2048)
	if err != nil {
		return false, "", err
	}

	// 格式检查
	result = con.checkIsValidDefaultPWD(strPwd)
	if !result {
		msg = loadString(visitor.LangType, "IDS_DEFAULT_PWD_INVALID")
	}
	return
}

// checkDefaultPWDAuthority 检测用户密码配置管理权限
func (con *config) checkDefaultPWDAuthority(userID string) (err error) {
	// 获取用户角色信息
	roleIDs, err := getRolesByUserID(con.role, userID)
	if err != nil {
		return err
	}

	// 超级管理员或者系统管理员角色拥有管理权限
	if roleIDs[interfaces.SystemRoleSuperAdmin] ||
		roleIDs[interfaces.SystemRoleSecAdmin] {
		return nil
	}
	return rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
}

// checkIsValidDefaultPWD 密码可以包含ASCII字符的任意组合，特殊字符为~!%#$@-_.， 且长度[6-100]
func (con *config) checkIsValidDefaultPWD(pwd string) bool {
	str := "^([\x20-\x7E]{6,100})$"
	reg := regexp.MustCompile(str)
	return reg.MatchString(pwd)
}
