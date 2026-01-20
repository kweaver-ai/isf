package general

import (
	"encoding/json"
	"os"
	"reflect"

	"policy_mgnt/tapi/sharemgnt"
	"policy_mgnt/thrift"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"
)

// Management 一般策略管理
type Management interface {
	ListPolicy(int, int, []string) ([]models.Policy[[]byte], int, error)
	SetPolicyValue(map[string][]byte, bool) error
	SetPolicyState([]string, map[State]interface{}) error
}

type mgnt struct {
	sharemgnt thrift.ShareMgnt
}

// NewManagement 初始化管理实例
func NewManagement() (Management, error) {
	client, err := thrift.NewShareMgnt()
	if err != nil {
		return nil, err
	}
	return NewManagementWithClient(client), nil
}

// NewManagementWithClient 初始化管理实例
func NewManagementWithClient(client thrift.ShareMgnt) Management {
	return &mgnt{sharemgnt: client}
}

// ListPolicy 分页获取策略信息
//
// 返回策略列表，策略总数
func (m *mgnt) ListPolicy(start int, limit int, names []string) (result []models.Policy[[]byte], count int, err error) {
	var count64 int64
	if start < 0 {
		start = 0
	}
	if limit < -1 {
		limit = -1
	}

	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	if len(names) > 0 {
		err = db.Model(&models.Policy[[]byte]{}).Where("f_name in (?)", names).Count(&count64).Error
	} else {
		err = db.Model(&models.Policy[[]byte]{}).Count(&count64).Error
	}
	if err != nil {
		return
	}
	count = int(count64)

	query := db.Order("f_name")
	if limit != -1 {
		query = query.Offset(start).Limit(limit)
	} else {
		query = query.Offset(start).Limit(count)
	}

	if len(names) > 0 {
		query = query.Where("f_name in (?)", names)
	}
	err = query.Find(&result).Error
	return
}

func (m *mgnt) checkPolicyName(names []string) error {
	var unknownName []string
	for _, name := range names {
		if _, ok := getDefaultPolicy()[name]; !ok {
			unknownName = append(unknownName, name)
		}
	}
	if len(unknownName) > 0 {
		return errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": unknownName}})
	}
	return nil
}

// 通过json比较两个对象是否一致
// 先转为json字符串，再转成inteface{}
// 最后递归比较两个值
func objEqual(obj1 interface{}, obj2 interface{}) bool {
	var bytes1, bytes2 []byte
	var cmpObj1, cmpObj2 interface{}
	var err error
	bytes1, err = json.Marshal(obj1)
	if err != nil {
		return false
	}
	bytes2, err = json.Marshal(obj2)
	if err != nil {
		return false
	}
	err = json.Unmarshal(bytes1, &cmpObj1)
	if err != nil {
		return false
	}
	err = json.Unmarshal(bytes2, &cmpObj2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(cmpObj1, cmpObj2)
}

// JSON字符串到策略内容转换
func (m *mgnt) bytes2Value(name string, bytes []byte) (value PolicyValue, err error) {
	defaultValue, ok := getDefaultPolicy()[name]
	if !ok {
		err = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{name}}})
		return
	}
	value = newPolicyValue(defaultValue)
	var actualValue map[string]interface{}
	err = json.Unmarshal(bytes, &actualValue)
	params := []string{"value"}
	if err != nil {
		// 非object
		value = nil
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
		return
	}
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		// 类型不匹配
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
		return
	}
	// 检查字符串和策略内容结构是否一致
	if !objEqual(value, actualValue) {
		// 属性不匹配
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
		return
	}
	return
}

// SetPolicyValue 设置策略内容
func (m *mgnt) SetPolicyValue(params map[string][]byte, force bool) (err error) {
	var names []string
	dbType := os.Getenv("DB_TYPE")
	// 不存在策略
	for name := range params {
		names = append(names, name)
	}
	err = m.checkPolicyName(names)
	if err != nil {
		return
	}
	// 内容错误
	inputValues := make(map[string]PolicyValue, 5)
	for name, bytes := range params {
		v, serr := m.bytes2Value(name, bytes)
		if serr != nil {
			err = serr
			return
		}

		serr = v.CheckParams()
		if serr != nil {
			err = serr
			return
		}

		inputValues[name] = v
	}

	db, err := api.ConnectDB()
	if err != nil {
		return
	}
	if dbType == "DM8" {
		var result []models.Policy[string]
		// 根据名字查策略
		err = db.Where("f_name in (?)", names).Find(&result).Error
		if err != nil {
			return
		}

		if !force {
			// 已锁定策略检查
			var lockedNames []string
			for _, policy := range result {
				if policy.Locked {
					lockedNames = append(lockedNames, policy.Name)
				}
			}
			if len(lockedNames) > 0 {
				err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"policys": lockedNames}})
				return
			}
		}

		tx := db.Begin()
		for _, policy := range result {
			inputValue := inputValues[policy.Name]
			dbValue, serr := inputValue.Merge([]byte(policy.Value))
			if serr != nil {
				err = serr
				return
			}

			policy.Value = string(dbValue)
			err = m.ApplyPolicy(policy.Name, []byte(policy.Value))
			if err != nil {
				return err
			}
			err = db.Save(&policy).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
	} else {
		var result []models.Policy[[]byte]
		// 根据名字查策略
		err = db.Where("f_name in (?)", names).Find(&result).Error
		if err != nil {
			return
		}

		if !force {
			// 已锁定策略检查
			var lockedNames []string
			for _, policy := range result {
				if policy.Locked {
					lockedNames = append(lockedNames, policy.Name)
				}
			}
			if len(lockedNames) > 0 {
				err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"policys": lockedNames}})
				return
			}
		}

		tx := db.Begin()
		for _, policy := range result {
			inputValue := inputValues[policy.Name]
			dbValue, serr := inputValue.Merge(policy.Value)
			if serr != nil {
				err = serr
				return
			}

			policy.Value = dbValue
			err = m.ApplyPolicy(policy.Name, policy.Value)
			if err != nil {
				return err
			}
			err = db.Save(&policy).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
	}

	return nil
}

// ApplyPolicy应用策略
func (m *mgnt) ApplyPolicy(name string, value []byte) (err error) {
	switch name {
	case "password_strength_meter":
		err = m.ApplyPasswordStrengthMeterByValue(value)
	case "multi_factor_auth":
		err = m.ApplyMultiFactorAuthByValue(value)
	case "system_protection_levels":
		err = m.ApplyPasswordConfigByLevel(value)
	}

	if err != nil {
		return
	}
	return nil
}

// ApplyPolicy应用密码长度策略
func (m *mgnt) ApplyPasswordStrengthMeterByValue(value []byte) (err error) {
	var actualValue PasswordStrengthMeter
	err = json.Unmarshal(value, &actualValue)
	if err != nil {
		return
	}

	// 获取当前的密码策略
	res, err := m.sharemgnt.GetPasswordConfig()
	if err != nil {
		return
	}

	// 修改密码长度策略，其余保持不变
	var config sharemgnt.NcTUsrmPasswordConfig
	config.StrongStatus = actualValue.Enable
	config.ExpireTime = res.ExpireTime
	config.LockStatus = res.LockStatus
	config.PasswdErrCnt = res.PasswdErrCnt
	config.PasswdLockTime = res.PasswdLockTime
	config.StrongPwdLength = int32(actualValue.Length)
	err = m.sharemgnt.ApplyPasswordStrengthMeter(config)
	if err != nil {
		return
	}
	return nil
}

// ApplyPolicy应用多因子认证
func (m *mgnt) ApplyMultiFactorAuthByValue(value []byte) (err error) {
	var actualValue MultiFactorAuth
	err = json.Unmarshal(value, &actualValue)
	if err != nil {
		return
	}

	// 图形验证码配置
	var config sharemgnt.NcTVcodeConfig
	if actualValue.Enable {
		config.IsEnable = actualValue.ImageVcode
		config.PasswdErrCnt = int32(actualValue.PasswordErrorCount)
	} else {
		config.IsEnable = false
		config.PasswdErrCnt = 0
	}
	err = m.sharemgnt.ApplyImageVcode(config)
	if err != nil {
		return
	}

	// 短信验证码、动态密码配置
	configValueMap := make(map[string]bool)
	if actualValue.Enable {
		configValueMap["auth_by_sms"] = actualValue.SMSVcode
		configValueMap["auth_by_OTP"] = actualValue.OTP
	} else {
		configValueMap["auth_by_sms"] = false
		configValueMap["auth_by_OTP"] = false
	}
	configValue, err := json.Marshal(configValueMap)
	if err != nil {
		return
	}
	configKey := "dualfactor_auth_server_status"
	err = m.sharemgnt.ApplyMultiFactorAuth(configKey, string(configValue))
	if err != nil {
		return
	}
	return nil
}

// ApplyPasswordConfigByLevel根据系统密级应用密码配置策略
func (m *mgnt) ApplyPasswordConfigByLevel(value []byte) (err error) {
	var actualValue SystemProtectionLevels
	err = json.Unmarshal(value, &actualValue)
	if err != nil {
		return
	}

	// 获取当前的密码策略
	res, err := m.sharemgnt.GetPasswordConfig()
	if err != nil {
		return
	}
	var config sharemgnt.NcTUsrmPasswordConfig
	if actualValue.Level == 1 {
		if res.ExpireTime == -1 || res.ExpireTime > 30 {
			config.ExpireTime = 30
		} else {
			config.ExpireTime = res.ExpireTime
		}
	} else if actualValue.Level == 2 {
		if res.ExpireTime == -1 || res.ExpireTime > 7 {
			config.ExpireTime = 7
		} else {
			config.ExpireTime = res.ExpireTime
		}
	} else {
		if res.ExpireTime == -1 || res.ExpireTime > 3 {
			config.ExpireTime = 3
		} else {
			config.ExpireTime = res.ExpireTime
		}
	}
	config.StrongStatus = res.StrongStatus
	config.LockStatus = res.LockStatus
	config.PasswdErrCnt = res.PasswdErrCnt
	config.PasswdLockTime = res.PasswdLockTime
	config.StrongPwdLength = res.StrongPwdLength
	err = m.sharemgnt.ApplyPasswordStrengthMeter(config)
	if err != nil {
		return
	}
	return nil
}

// SetPolicyState 设置策略状态
func (m *mgnt) SetPolicyState(names []string, states map[State]interface{}) (err error) {
	err = m.checkPolicyName(names)
	if err != nil {
		return
	}

	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	var result []models.Policy[[]byte]
	err = db.Select("f_name,f_locked").Where("f_name in (?)", names).Find(&result).Error
	if err != nil {
		return
	}

	values := make(map[string]interface{})
	for state, value := range states {
		// 映射数据库字段
		switch state {
		case StateLocked:
			values["f_locked"] = value
		}
	}

	tx := db.Begin()
	for _, policy := range result {
		err = tx.Model(&policy).Updates(values).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	return nil
}
