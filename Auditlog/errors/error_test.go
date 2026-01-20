package errors

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"AuditLog/infra/cmp/langcmp"
)

func TestMain(m *testing.M) {
	langcmp.NewLangCmp().SetSysDefLang(string(langcmp.ZhCN))
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	测试用例 := []struct {
		名称      string
		语言      string
		错误码     int
		原因      string
		详情      interface{}
		期望描述    string
		期望HTTP码 int
	}{
		{
			名称:      "中文错误",
			语言:      "zh",
			错误码:     400001001,
			原因:      "参数错误",
			详情:      map[string]string{"field": "name"},
			期望描述:    i18ns[400001001].Description[langcmp.ZhCN],
			期望HTTP码: 400,
		},
		{
			名称:      "英文错误",
			语言:      "en",
			错误码:     500001001,
			原因:      "internal error",
			详情:      nil,
			期望描述:    i18ns[500001001].Description[langcmp.ZhCN],
			期望HTTP码: 500,
		},
		{
			名称:      "空语言默认值",
			语言:      "",
			错误码:     404001001,
			原因:      "not found",
			详情:      nil,
			期望描述:    i18ns[404001001].Description[langcmp.NewLangCmp().GetSysDefaultLang()],
			期望HTTP码: 404,
		},
	}

	for _, tc := range 测试用例 {
		t.Run(tc.名称, func(t *testing.T) {
			err := New(tc.语言, tc.错误码, tc.原因, tc.详情)

			assert.Equal(t, tc.错误码, err.Code())
			assert.Equal(t, tc.期望HTTP码, err.HTTPCode())
			assert.Equal(t, tc.期望描述, err.description)
			assert.Equal(t, tc.原因, err.cause)
			assert.Equal(t, tc.详情, err.detail)
		})
	}
}

func TestNewCtx(t *testing.T) {
	ctx := context.Background()
	// 在上下文中设置语言
	ctx = context.WithValue(ctx, "language", langcmp.ZhCN)

	err := NewCtx(ctx, 400001001, "参数错误", nil)

	assert.Equal(t, 400001001, err.Code())
	assert.Equal(t, 400, err.HTTPCode())
	assert.Equal(t, i18ns[400001001].Description[langcmp.ZhCN], err.description)
}

func TestErrorResp_MarshalJSON(t *testing.T) {
	err := New("zh", 400001001, "参数错误", map[string]string{"field": "name"})

	jsonBytes, marshalErr := json.Marshal(err)
	assert.NoError(t, marshalErr)

	var result map[string]interface{}
	unmarshalErr := json.Unmarshal(jsonBytes, &result)

	assert.NoError(t, unmarshalErr)
	assert.Equal(t, float64(400001001), result["code"])
	assert.Equal(t, "参数错误", result["cause"])
	assert.Equal(t, err.description, result["description"])
	assert.Equal(t, err.description, result["message"])
	assert.Equal(t, err.solution, result["solution"])
}

func TestErrorResp_Error(t *testing.T) {
	err := New("zh_cn", 400001001, "参数错误", nil)

	errStr := err.Error()

	assert.Contains(t, errStr, "code: 400001001")
	assert.Contains(t, errStr, "Cause: 参数错误")
}
