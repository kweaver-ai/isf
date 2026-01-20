package general

import (
	"encoding/json"
	"testing"

	"policy_mgnt/tapi/sharemgnt"
	"policy_mgnt/test"
	"policy_mgnt/test/mock_thrift"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func mockShareMgnt(t *testing.T) (*gomock.Controller, *mock_thrift.MockShareMgnt) {
	ctrl := gomock.NewController(t)
	mgnt := mock_thrift.NewMockShareMgnt(ctrl)
	return ctrl, mgnt
}

// 检查策略是否存在
func Test_mgnt_checkPolicyName(t *testing.T) {
	m := &mgnt{}

	var names []string
	var err, assertErr error

	names = []string{"unknown", "multi_factor_auth"}
	err = m.checkPolicyName(names)
	assertErr = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"unknown"}}})
	assert.Equal(t, assertErr, err)

	names = []string{"client_restriction", "multi_factor_auth", "unknown1", "unknown2"}
	err = m.checkPolicyName(names)
	assertErr = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"unknown1", "unknown2"}}})
	assert.Equal(t, err, assertErr)
}

// objEqual 检查,传入的不同
func Test_objNotEqual(t *testing.T) {
	right := ClientRestriction{
		Windows: true,
		Mac:     true,
	}
	otherRight := ClientRestriction{
		Windows: false,
		Mac:     true,
	}
	var tests = []struct {
		arg1 interface{}
		arg2 interface{}
		want bool
	}{
		{right, right, true},
		{right, otherRight, false},
	}
	for _, test := range tests {
		if got := objEqual(test.arg1, test.arg2); got != test.want {
			t.Errorf("input is (%v,%v), wanted %v, bug got %v", test.arg1, test.arg2, test.want, got)
		}
	}
}

// 参数范围校验
func Test_mgnt_CheckParams(t *testing.T) {
	// case name = "multi_factor_auth"
	// 参数不合法(PasswordErrorCount范围为0-99)
	value := MultiFactorAuth{
		Enable:             true,
		ImageVcode:         true,
		PasswordErrorCount: -1,
		SMSVcode:           false,
		OTP:                false,
	}
	var m PolicyValue
	m = &value
	err := m.CheckParams()
	expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"multi_factor_auth"}}, Cause: "password_error_count must be an interger between 0 and 99."})
	assert.Equal(t, expectErr, err)

	value.PasswordErrorCount = 100
	err = m.CheckParams()
	expectErr = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"multi_factor_auth"}}, Cause: "password_error_count must be an interger between 0 and 99."})
	assert.Equal(t, expectErr, err)

	// 合法
	value.PasswordErrorCount = 10
	err = m.CheckParams()
	assert.Nil(t, err)

	// case 多因子认证，同时开启一个以上开关，报错
	wrongMFAValue := MultiFactorAuth{
		Enable:             true,
		ImageVcode:         true,
		PasswordErrorCount: 8,
		SMSVcode:           false,
		OTP:                true,
	}
	m = &wrongMFAValue
	err = m.CheckParams()
	expectErr = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"multi_factor_auth"}}, Cause: "mutil-factor auth can only enable one at a time."})
	assert.Equal(t, expectErr, err)

	// case name = "password_strength_meter"
	// 参数不合法
	value1 := PasswordStrengthMeter{
		Enable: true,
		Length: 0,
	}
	m = &value1
	err = m.CheckParams()
	expectErr = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"password_strength_meter"}}, Cause: "Strong password minimum length must ben an interger between 8 and 99."})
	assert.Equal(t, expectErr, err)

	value1.Length = 100
	err = m.CheckParams()
	expectErr = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"password_strength_meter"}}, Cause: "Strong password minimum length must ben an interger between 8 and 99."})
	assert.Equal(t, expectErr, err)

	// 合法
	value1.Length = 8
	err = m.CheckParams()
	assert.Nil(t, err)

	// password_strength_meter的enable为false时，不限制length
	value1.Enable = false
	value1.Length = -1
	err = m.CheckParams()
	assert.Nil(t, err)
}

// 测试部分需要修改的参数
func Test_mgnt_Merge(t *testing.T) {
	// case: 强密码
	oldValue := PasswordStrengthMeter{
		Enable: true,
		Length: 8,
	}
	newValue := PasswordStrengthMeter{
		Enable: false,
		Length: -1,
	}
	oldByte, err := json.Marshal(oldValue)

	dbByte, err := newValue.Merge(oldByte)
	assert.Nil(t, err)
	expectValue := PasswordStrengthMeter{
		Enable: false,
		Length: 8,
	}
	exp, err := json.Marshal(expectValue)
	assert.Equal(t, exp, dbByte)

	// case:多因子认证
	oldValue1 := MultiFactorAuth{
		Enable:             true,
		ImageVcode:         true,
		PasswordErrorCount: 10,
		SMSVcode:           false,
		OTP:                false,
	}
	newValue1 := MultiFactorAuth{
		Enable:             true,
		ImageVcode:         false,
		PasswordErrorCount: 20,
		SMSVcode:           false,
		OTP:                false,
	}
	oldByte1, err := json.Marshal(oldValue1)
	dbByte1, err := newValue1.Merge(oldByte1)
	assert.Nil(t, err)
	expectValue1 := MultiFactorAuth{
		Enable:             true,
		ImageVcode:         false,
		PasswordErrorCount: 10,
		SMSVcode:           false,
		OTP:                false,
	}
	exp1, err := json.Marshal(expectValue1)
	assert.Equal(t, exp1, dbByte1)
}

// JSON字符串到策略配置转换
func Test_mgnt_bytes2Value(t *testing.T) {
	name := "password_strength_meter"
	bytes := []byte(`{"enable":false,"length":8}`)
	expectValue := &PasswordStrengthMeter{
		Enable: false,
		Length: 8,
	}

	m := &mgnt{}
	value, err := m.bytes2Value(name, bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectValue, value)
}

// 策略配置转换，未知策略
func Test_mgnt_bytes2ValueUnknown(t *testing.T) {
	name := "unknown"
	bytes := []byte(`{"enable":false,"length":8}`)

	m := &mgnt{}
	value, err := m.bytes2Value(name, bytes)
	expectErr := errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{name}}})
	assert.Nil(t, value)
	assert.Equal(t, expectErr, err)
}

// 策略配置转换，不正确策略内容
func Test_mgnt_bytes2ValueIncorrect(t *testing.T) {
	var bytes []byte
	var value PolicyValue
	var err error
	m := &mgnt{}
	name := "password_strength_meter"
	expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"value"}}})
	// 非object
	bytes = []byte(`"enable":false`)
	value, err = m.bytes2Value(name, bytes)
	assert.Nil(t, value)
	assert.Equal(t, expectErr, err)
}

// 获取所有策略信息
func Test_mgnt_ListPolicy(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	m, _ := NewManagement()

	var result []models.Policy[[]byte]
	var count int
	var err error

	// start 0 limit 1
	result, count, err = m.ListPolicy(0, 1, []string{})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "client_restriction")

	// start 0 limit -10
	result, count, err = m.ListPolicy(0, -10, []string{})
	assert.Nil(t, err)
	assert.Equal(t, len(result), count)

	// start -1 limit 1
	result, count, err = m.ListPolicy(-1, 1, []string{})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "client_restriction")

	// start 1 limit 1
	result, count, err = m.ListPolicy(1, 1, []string{})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "multi_factor_auth")

	// start 10000 limit 1
	result, count, err = m.ListPolicy(10000, 1, []string{})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 0)
}

// 获取指定策略
func Test_mgnt_ListPolicyByName(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	m, _ := NewManagement()

	var result []models.Policy[[]byte]
	var count int
	var err error

	// name存在
	result, count, err = m.ListPolicy(0, 1, []string{"client_restriction"})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "client_restriction")
	assert.Equal(t, len(result), count)

	// name不存在
	result, count, err = m.ListPolicy(0, 1, []string{"not exist"})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 0)
	assert.Equal(t, len(result), 0)

	// 多个name都存在, limit为1
	result, count, err = m.ListPolicy(0, 1, []string{"client_restriction", "multi_factor_auth"})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)

	// 多个name都存在, limit为-1, start=0
	result, count, err = m.ListPolicy(0, -1, []string{"client_restriction", "multi_factor_auth"})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 2)
	assert.Equal(t, len(result), count)

	// 多个name都存在, limit为-1, start=1
	result, count, err = m.ListPolicy(1, -1, []string{"client_restriction", "multi_factor_auth"})
	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
}

// 设置多因子登录
func Test_mgnt_SetPolicyValueAuth(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	name := "multi_factor_auth"
	value := MultiFactorAuth{
		Enable:             true,
		ImageVcode:         true,
		PasswordErrorCount: 10,
		SMSVcode:           false,
		OTP:                false,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	ctrl, client := mockShareMgnt(t)
	defer ctrl.Finish()

	config := sharemgnt.NcTVcodeConfig{
		IsEnable:     true,
		PasswdErrCnt: 10,
	}
	client.EXPECT().ApplyImageVcode(config).Return(nil)

	configValueMap := make(map[string]bool)
	configValueMap["auth_by_sms"] = false
	configValueMap["auth_by_OTP"] = false
	configValue, _ := json.Marshal(configValueMap)
	configKey := "dualfactor_auth_server_status"
	client.EXPECT().ApplyMultiFactorAuth(configKey, string(configValue)).Return(nil)

	mgnt := NewManagementWithClient(client)
	err := mgnt.SetPolicyValue(param, false)
	assert.Nil(t, err)

	db, _ := api.ConnectDB()
	var actualValue models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&actualValue)
	assert.JSONEq(t, string(jsonValue), string(actualValue.Value))
}

// 设置密码长度
func Test_mgnt_SetPolicyValuePwd(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	name := "password_strength_meter"
	value := PasswordStrengthMeter{
		Enable: true,
		Length: 18,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	ctrl, client := mockShareMgnt(t)
	defer ctrl.Finish()

	var tlocktime int32 = 60
	config := sharemgnt.NcTUsrmPasswordConfig{
		StrongStatus:    true,
		ExpireTime:      -1,
		LockStatus:      false,
		PasswdErrCnt:    5,
		PasswdLockTime:  &tlocktime,
		StrongPwdLength: 18,
	}
	client.EXPECT().GetPasswordConfig().Return(&config, nil)
	client.EXPECT().ApplyPasswordStrengthMeter(config).Return(nil)

	mgnt := NewManagementWithClient(client)
	err := mgnt.SetPolicyValue(param, false)
	assert.Nil(t, err)

	db, _ := api.ConnectDB()
	var actualValue models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&actualValue)
	assert.JSONEq(t, string(jsonValue), string(actualValue.Value))
}

// 设置系统保护密级
func Test_mgnt_SetPasswordConfigByLevel(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	name := "system_protection_levels"
	value := SystemProtectionLevels{
		Level: 1,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	ctrl, client := mockShareMgnt(t)
	defer ctrl.Finish()

	var tlocktime int32 = 60
	config := sharemgnt.NcTUsrmPasswordConfig{
		StrongStatus:    true,
		ExpireTime:      30,
		LockStatus:      false,
		PasswdErrCnt:    5,
		PasswdLockTime:  &tlocktime,
		StrongPwdLength: 18,
	}
	client.EXPECT().GetPasswordConfig().Return(&config, nil)
	client.EXPECT().ApplyPasswordStrengthMeter(config).Return(nil)

	mgnt := NewManagementWithClient(client)
	err := mgnt.SetPolicyValue(param, false)
	assert.Nil(t, err)

	db, _ := api.ConnectDB()
	var actualValue models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&actualValue)
	assert.JSONEq(t, string(jsonValue), string(actualValue.Value))
}

// 设置网段白名单开关
func Test_mgnt_SetNetworkResitriction(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	name := "network_restriction"
	value := NetworkResitriction{
		IsEnabled: true,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	ctrl, client := mockShareMgnt(t)
	defer ctrl.Finish()

	mgnt := NewManagementWithClient(client)
	err := mgnt.SetPolicyValue(param, false)
	assert.Nil(t, err)

	db, _ := api.ConnectDB()
	var actualValue models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&actualValue)
	assert.JSONEq(t, string(jsonValue), string(actualValue.Value))
}

// 设置网段未绑定网段用户可访问开关
func Test_mgnt_SetNoNetworkPolicyAccessor(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	name := "no_network_policy_accessor"
	value := NoNetworkPolicyAccessor{
		IsEnabled: true,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	ctrl, client := mockShareMgnt(t)
	defer ctrl.Finish()

	mgnt := NewManagementWithClient(client)
	err := mgnt.SetPolicyValue(param, false)
	assert.Nil(t, err)

	db, _ := api.ConnectDB()
	var actualValue models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&actualValue)
	assert.JSONEq(t, string(jsonValue), string(actualValue.Value))
}

// 设置设备禁止登陆
func Test_mgnt_SetPolicyValueDevice(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	name := "client_restriction"
	value := ClientRestriction{
		Windows: true,
		Mac:     true,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	ctrl, client := mockShareMgnt(t)
	defer ctrl.Finish()

	mgnt := NewManagementWithClient(client)
	err := mgnt.SetPolicyValue(param, false)
	assert.Nil(t, err)

	db, _ := api.ConnectDB()
	var actualValue models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&actualValue)
	assert.JSONEq(t, string(jsonValue), string(actualValue.Value))
}

// 设置不存在策略
func Test_mgnt_SetPolicyValueNotFound(t *testing.T) {
	name := "client_restriction"
	value := ClientRestriction{
		Windows: true,
		Mac:     true,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name:      jsonValue,
		"unknown": []byte(`{"enable":false}`),
	}

	mgnt := NewManagementWithClient(nil)
	err := mgnt.SetPolicyValue(param, false)
	assertErr := errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"unknown"}}})
	assert.Equal(t, err, assertErr)
}

// 设置错误策略内容
func Test_mgnt_SetPolicyValueIncorrectValue(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	name := "client_restriction"
	param := map[string][]byte{
		name: []byte(`{"windows":true,"mac":true}`),
	}

	mgnt := NewManagementWithClient(nil)
	err := mgnt.SetPolicyValue(param, false)
	assertErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"value"}}})
	assert.Equal(t, err, assertErr)
}

// 设置已锁定策略
func Test_mgnt_SetPolicyValueLocked(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()
	db, _ := api.ConnectDB()

	name := "client_restriction"
	value := ClientRestriction{
		Windows: true,
		Mac:     true,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	db.Model(models.Policy[[]byte]{Name: name}).Update("f_locked", true)

	mgnt := NewManagementWithClient(nil)

	err := mgnt.SetPolicyValue(param, false)
	assertErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"policys": []string{name}}})
	assert.Equal(t, err, assertErr)
	var policy models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&policy)
	var actualValue ClientRestriction
	json.Unmarshal(policy.Value, &actualValue)
	assert.Equal(t, actualValue.Windows, false)
	assert.Equal(t, actualValue.Mac, false)
}

// 设置已锁定策略，强制
func Test_mgnt_SetPolicyValueLockedForce(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()
	db, _ := api.ConnectDB()

	name := "client_restriction"
	value := ClientRestriction{
		Windows: true,
		Mac:     true,
	}
	jsonValue, _ := json.Marshal(value)
	param := map[string][]byte{
		name: jsonValue,
	}

	ctrl, client := mockShareMgnt(t)
	defer ctrl.Finish()

	mgnt := NewManagementWithClient(client)

	db.Model(models.Policy[[]byte]{}).UpdateColumn("f_locked", true)

	err := mgnt.SetPolicyValue(param, true)
	assert.Nil(t, err)

	var actualValue models.Policy[[]byte]
	db.Where("f_name = ?", name).Take(&actualValue)
	assert.JSONEq(t, string(jsonValue), string(actualValue.Value))
}

// 设置策略状态
func Test_mgnt_SetPolicyState(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	names := []string{"client_restriction", "multi_factor_auth"}
	states := map[State]interface{}{
		StateLocked: true,
	}

	mgnt := NewManagementWithClient(nil)
	err := mgnt.SetPolicyState(names, states)
	assert.Nil(t, err)

	var actualValue []models.Policy[[]byte]
	db, _ := api.ConnectDB()
	db.Where("f_name in (?)", names).Find(&actualValue)

	for _, val := range actualValue {
		assert.True(t, val.Locked)
	}
}

// 设置不存在策略状态
func Test_mgnt_SetPolicyStateNotFound(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	CreateDefaultPolicy()

	namesFull := []string{"client_restriction", "multi_factor_auth", "unknown"}
	namesActual := []string{"client_restriction", "multi_factor_auth"}
	states := map[State]interface{}{
		StateLocked: true,
	}

	mgnt := NewManagementWithClient(nil)
	err := mgnt.SetPolicyState(namesFull, states)
	assertErr := errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"unknown"}}})
	assert.Equal(t, err, assertErr)

	var actualValue []models.Policy[[]byte]
	db, _ := api.ConnectDB()
	db.Where("f_name in (?)", namesActual).Find(&actualValue)

	for _, val := range actualValue {
		assert.False(t, val.Locked)
	}
}
