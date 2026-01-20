package general

import (
	"encoding/json"
	"reflect"

	"policy_mgnt/utils/errors"

	"policy_mgnt/utils/gocommon/api"
)

// PolicyValue policy value interface
type PolicyValue interface {
	Name() string
	CheckParams() error
	Merge([]byte) ([]byte, error)
}

func newPolicyValue(value PolicyValue) PolicyValue {
	policyType := reflect.TypeOf(value)
	// 解指针方案来自：https://github.com/jinzhu/copier indirect
	for policyType.Kind() == reflect.Ptr {
		policyType = policyType.Elem()
	}
	return reflect.New(policyType).Interface().(PolicyValue)
}

// PasswordStrengthMeter 密码强度检测策略
type PasswordStrengthMeter struct {
	Enable bool `json:"enable"`
	Length int  `json:"length"`
}

// Name policy name
func (v *PasswordStrengthMeter) Name() string {
	return "password_strength_meter"
}

// CheckParams 检查参数
func (v *PasswordStrengthMeter) CheckParams() (err error) {
	n := v.Length
	// 如果为弱密码，不检查长度
	if !v.Enable {
		return
	}
	// 如果为强密码，限制密码长度的左区间在8至99（右区间默认为100）
	if n > 99 || n < 8 {
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"password_strength_meter"}}, Cause: "Strong password minimum length must ben an interger between 8 and 99."})
		return
	}
	return
}

// Merge args:oldByte 为数据库的旧数据，v为传入的新数据
func (v *PasswordStrengthMeter) Merge(oldByte []byte) (outputByte []byte, err error) {
	// enable为false，length保持原来的数据不变
	if !v.Enable {
		var oldValue PasswordStrengthMeter
		err = json.Unmarshal(oldByte, &oldValue)
		if err != nil {
			return
		}
		// oldValue.Length为数据库中的数据
		v.Length = oldValue.Length
		outputByte, err = json.Marshal(v)
		if err != nil {
			return
		}
		return
	}
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}

// MultiFactorAuth 多因子认证策略
type MultiFactorAuth struct {
	Enable             bool `json:"enable"`
	ImageVcode         bool `json:"image_vcode"`
	PasswordErrorCount int  `json:"password_error_count"`
	SMSVcode           bool `json:"sms_vcode"`
	OTP                bool `json:"otp"`
}

// Name policy name
func (v *MultiFactorAuth) Name() string {
	return "multi_factor_auth"
}

// CheckParams 检查参数
func (v *MultiFactorAuth) CheckParams() (err error) {
	n := v.PasswordErrorCount
	if n > 99 || n < 0 {
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"multi_factor_auth"}}, Cause: "password_error_count must be an interger between 0 and 99."})
		return
	}

	// 限制同时只能开一个开关（临时功能，后续删除）
	switchSlice := []bool{v.ImageVcode, v.SMSVcode, v.OTP}
	flag := false
	for _, v := range switchSlice {
		if flag && v {
			err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"multi_factor_auth"}}, Cause: "mutil-factor auth can only enable one at a time."})
			return
		}
		flag = v || flag
	}

	return
}

// Merge 传入数据库旧的数据，修改后重新存入
func (v *MultiFactorAuth) Merge(oldByte []byte) (outputByte []byte, err error) {
	// 总开关关闭或者总开关打开、图形验证码关闭时，错误次数使用旧数据
	if !v.Enable || (v.Enable && !v.ImageVcode) {
		var oldValue MultiFactorAuth
		err = json.Unmarshal(oldByte, &oldValue)
		if err != nil {
			return
		}
		// oldValue为数据库中的数据
		v.PasswordErrorCount = oldValue.PasswordErrorCount
		outputByte, err = json.Marshal(v)
		if err != nil {
			return
		}
		return
	}
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}

// ClientRestriction 客户端限制策略
type ClientRestriction struct {
	PCWEB     bool `json:"pc_web"`
	MobileWEB bool `json:"mobile_web"`
	Windows   bool `json:"windows"`
	Mac       bool `json:"mac"`
	Android   bool `json:"android"`
	IOS       bool `json:"ios"`
	Linux     bool `json:"linux"`
}

// Name policy name
func (v *ClientRestriction) Name() string {
	return "client_restriction"
}

// CheckParams 检查参数
func (v *ClientRestriction) CheckParams() (err error) {
	return
}

// Merge 传入数据库旧的数据，修改后重新存入
func (v *ClientRestriction) Merge(oldByte []byte) (outputByte []byte, err error) {
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}

// UserDocumentSharing 个人文档共享开关策略
type UserDocumentSharing struct {
	Anyshare bool `json:"anyshare"`
	HTTP     bool `json:"http"`
}

// Name policy name
func (v *UserDocumentSharing) Name() string {
	return "user_document_sharing"
}

// CheckParams 检查参数
func (v *UserDocumentSharing) CheckParams() (err error) {
	return
}

// Merge 传入数据库旧的数据，修改后重新存入
func (v *UserDocumentSharing) Merge(oldByte []byte) (outputByte []byte, err error) {
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}

// UserDocument 新建个人文档策略
type UserDocument struct {
	Create bool    `json:"create"`
	Size   float32 `json:"size"`
}

// Name policy name
func (v *UserDocument) Name() string {
	return "user_document"
}

// CheckParams 检查参数
func (v *UserDocument) CheckParams() (err error) {
	return
}

// Merge 传入数据库旧的数据，修改后重新存入
func (v *UserDocument) Merge(oldByte []byte) (outputByte []byte, err error) {
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}

// NetworkResitriction 访问者网段
type NetworkResitriction struct {
	IsEnabled bool `json:"is_enabled"`
}

// Name policy name
func (v *NetworkResitriction) Name() string {
	return "network_restriction"
}

// CheckParams 检查参数
func (v *NetworkResitriction) CheckParams() (err error) {
	return
}

// Merge 传入数据库旧的数据，修改后重新存入
func (v *NetworkResitriction) Merge(oldByte []byte) (outputByte []byte, err error) {
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}

// NoPolicyAccessor 未配置访问者网段策略访问者
type NoNetworkPolicyAccessor struct {
	IsEnabled bool `json:"is_enabled"`
}

// Name policy name
func (v *NoNetworkPolicyAccessor) Name() string {
	return "no_network_policy_accessor"
}

// CheckParams 检查参数
func (v *NoNetworkPolicyAccessor) CheckParams() (err error) {
	return
}

// Merge 传入数据库旧的数据，修改后重新存入
func (v *NoNetworkPolicyAccessor) Merge(oldByte []byte) (outputByte []byte, err error) {
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}

// SystemProtectionLevels 系统保护密级
type SystemProtectionLevels struct {
	Level int `json:"level"`
}

// Name policy name
func (v *SystemProtectionLevels) Name() string {
	return "system_protection_levels"
}

// CheckParams 检查参数
func (v *SystemProtectionLevels) CheckParams() (err error) {
	n := v.Level
	// 保护等级为1，2，3
	if n > 3 || n < 1 {
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"system_protection_levels"}}, Cause: "System Protection Level must ben an interger between 1 and 3."})
		return
	}
	return
}

// Merge 传入数据库旧的数据，修改后重新存入
func (v *SystemProtectionLevels) Merge(oldByte []byte) (outputByte []byte, err error) {
	var oldValue SystemProtectionLevels
	err = json.Unmarshal(oldByte, &oldValue)
	if err != nil {
		return
	}
	// 系统保护密级不能设置为更低的密级
	if v.Level < oldValue.Level {
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"system_protection_levels"}}, Cause: "System Protection Level not set lower."})
		return
	}
	outputByte, err = json.Marshal(v)
	if err != nil {
		return
	}
	return
}
