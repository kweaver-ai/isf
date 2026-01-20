package helpers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"AuditLog/common/enums"
)

func TestGetXLanguage(t *testing.T) {
	// 设置测试用例
	tests := []struct {
		name          string
		languageValue string
		expected      string
	}{
		{
			name:          "正常语言值",
			languageValue: "zh-CN",
			expected:      "zh-CN",
		},
		{
			name:          "空语言值",
			languageValue: "",
			expected:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个新的 gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置请求头
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set(enums.XLanguage, tt.languageValue)

			// 调用测试函数
			result := GetXLanguage(c)

			// 验证结果
			assert.Equal(t, tt.expected, result)
		})
	}
}
