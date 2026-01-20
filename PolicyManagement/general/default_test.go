package general

import (
	"encoding/json"
	"testing"

	"policy_mgnt/test"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/stretchr/testify/assert"
)

// 默认策略创建
func TestCreateDefaultPolicy(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	err := CreateDefaultPolicy()
	assert.Nil(t, err)

	db, _ := api.ConnectDB()

	var policies []models.Policy[[]byte]
	db.Find(&policies)
	assert.NotEmpty(t, policies)

	for name, content := range getDefaultPolicy() {
		found := false
		for _, policy := range policies {
			if name == policy.Name {
				found = true
				c, _ := json.Marshal(content)
				assert.JSONEq(t, string(c), string(policy.Default))
				break
			}
		}
		assert.Truef(t, found, "Policy %s not in db", name)
	}
}

// 默认策略追加
func TestCreateDefaultPolicyAppend(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	db, _ := api.ConnectDB()

	// 添加一条到数据库
	defaultValue := getDefaultPolicy()
	for name, value := range defaultValue {
		data, _ := json.Marshal(value)
		policy := models.Policy[[]byte]{
			Name:    name,
			Default: data,
			Value:   data,
			Locked:  false,
		}
		db.Create(policy)
		break
	}

	err := CreateDefaultPolicy()
	assert.Nil(t, err)

	var policies []models.Policy[[]byte]
	db.Find(&policies)

	for name, content := range defaultValue {
		found := false
		for _, policy := range policies {
			if name == policy.Name {
				found = true
				c, _ := json.Marshal(content)
				assert.JSONEq(t, string(c), string(policy.Default))
				break
			}
		}
		assert.Truef(t, found, "Policy %s not in db", name)
	}
}

// 客户端登录策略值不包含linux
func TestCreateDefaultPolicyClientRestriction(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	db, _ := api.ConnectDB()

	// 添加不包含linux的客户端登录选项值至数据库
	clientResValue := map[string]bool{
		"pc_web":     false,
		"mobile_web": false,
		"windows":    false,
		"mac":        false,
		"android":    false,
		"ios":        false,
	}
	lientResValueByte, _ := json.Marshal(clientResValue)
	policy := models.Policy[[]byte]{
		Name:    "client_restriction",
		Default: lientResValueByte,
		Value:   lientResValueByte,
		Locked:  false,
	}

	err := db.Create(policy).Error
	assert.Nil(t, err)

	err = CreateDefaultPolicy()
	assert.Nil(t, err)

	var result models.Policy[[]byte]
	err = db.Where("f_name = ?", "client_restriction").Find(&result).Error
	assert.Nil(t, err)
	var defaultConfig map[string]bool
	json.Unmarshal(result.Default, &defaultConfig)
	assert.Equal(t, result.Name, "client_restriction")
	assert.Equal(t, defaultConfig["linux"], false)
}
