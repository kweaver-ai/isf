package locale

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"AuditLog/common/conf"
	"AuditLog/infra/cmp/langcmp"
)

func TestMain(m *testing.M) {
	langcmp.NewLangCmp().SetSysDefLang(string(langcmp.ZhCN))
	os.Exit(m.Run())
}

func TestGetI18nWithDefault(t *testing.T) {
	// 准备测试数据
	testI18nMap := map[string]map[langcmp.Lang]string{
		"test_key": {
			langcmp.ZhCN: "测试",
			langcmp.ZhTW: "測試",
			langcmp.En:   "test",
		},
	}

	tests := []struct {
		name     string
		ctx      context.Context
		key      string
		i18nMap  map[string]map[langcmp.Lang]string
		expected string
	}{
		{
			name:     "nil context",
			ctx:      nil,
			key:      "test_key",
			i18nMap:  testI18nMap,
			expected: testI18nMap["test_key"][langcmp.NewLangCmp().GetSysDefaultLang()],
		},
		{
			name:     "with gin context - ZH",
			ctx:      createGinContextWithLang("zh-CN"),
			key:      "test_key",
			i18nMap:  testI18nMap,
			expected: "测试",
		},
		{
			name:     "with gin context - EN",
			ctx:      createGinContextWithLang("EN-US"),
			key:      "test_key",
			i18nMap:  testI18nMap,
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getI18nWithDefault(tt.ctx, tt.key, tt.i18nMap)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetI18nCtx(t *testing.T) {
	ctx := createGinContextWithLang("zh-CN")
	result := GetI18nCtx(ctx, "log_role")
	assert.NotEmpty(t, result)
}

func TestGetRCLogLevelI18n(t *testing.T) {
	ctx := createGinContextWithLang("zh-CN")
	result := GetRCLogLevelI18n(ctx, "rc_log_level_info")
	assert.NotEmpty(t, result)
}

func TestGetRCLogLoginI18n(t *testing.T) {
	conf.InitJsonConf()
	ctx := createGinContextWithLang("zh-CN")
	result := GetRCLogLoginI18n(ctx, 1)
	assert.NotEmpty(t, result)
}

func TestGetRCLogMgntI18n(t *testing.T) {
	conf.InitJsonConf()
	ctx := createGinContextWithLang("zh-CN")
	result := GetRCLogMgntI18n(ctx, 1)
	assert.NotEmpty(t, result)
}

func TestGetRCLogOpI18n(t *testing.T) {
	conf.InitJsonConf()
	ctx := createGinContextWithLang("zh-CN")
	result := GetRCLogOpI18n(ctx, 1)
	assert.NotEmpty(t, result)
}

// 辅助函数：创建带有语言设置的 Gin Context
func createGinContextWithLang(lang string) context.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Header.Set("X-Language", lang)
	return c
}
