package helpers

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"AuditLog/common/enums"
	"AuditLog/infra/cmp/langcmp"
)

func TestMain(m *testing.M) {
	langcmp.NewLangCmp().SetSysDefLang(string(langcmp.ZhCN))
	os.Exit(m.Run())
}

func TestGetLangFromHeader(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() *gin.Context
		wantLang  langcmp.Lang
	}{
		{
			name: "从header获取语言成功",
			setupFunc: func() *gin.Context {
				c, _ := gin.CreateTestContext(nil)
				c.Request = &http.Request{
					Header: make(http.Header),
				}
				c.Request.Header.Set("X-Language", "zh-CN")
				return c
			},
			wantLang: langcmp.ZhCN,
		},
		{
			name: "header中无语言时使用默认语言",
			setupFunc: func() *gin.Context {
				c, _ := gin.CreateTestContext(nil)
				c.Request = &http.Request{
					Header: make(http.Header),
				}
				return c
			},
			wantLang: langcmp.NewLangCmp().GetSysDefaultLang(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupFunc()
			got := GetLangFromHeader(c)
			assert.Equal(t, tt.wantLang, got)
		})
	}
}

func TestGetLangFromCtx(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() context.Context
		wantLang  langcmp.Lang
		wantPanic bool
	}{
		{
			name: "从context获取语言成功",
			setupFunc: func() context.Context {
				return context.WithValue(context.Background(), enums.VisitLangCtxKey, langcmp.Lang("zh_cn"))
			},
			wantLang: langcmp.ZhCN,
		},
		{
			name: "context中无语言时使用默认语言",
			setupFunc: func() context.Context {
				return context.Background()
			},
			wantLang: langcmp.NewLangCmp().GetSysDefaultLang(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupFunc()

			if tt.wantPanic {
				assert.Panics(t, func() {
					GetLangFromCtx(ctx)
				})

				return
			}

			got := GetLangFromCtx(ctx)
			assert.Equal(t, tt.wantLang, got)
		})
	}
}
