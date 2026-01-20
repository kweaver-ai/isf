package helpers

import (
	"context"

	"github.com/gin-gonic/gin"

	"AuditLog/common/enums"
	"AuditLog/infra/cmp/langcmp"
)

// GetLangFromHeader 获取语言
// 1. 从gin header中获取语言
// 2. 不存在时，使用系统设置的语言
func GetLangFromHeader(c *gin.Context) (lang langcmp.Lang) {
	lang = langcmp.NewFromStr(GetXLanguage(c))
	if lang == "" {
		lang = langcmp.NewLangCmp().GetSysDefaultLang()
	}

	return
}

func GetLangFromCtx(ctx context.Context) (lang langcmp.Lang) {
	vInter := ctx.Value(enums.VisitLangCtxKey.String())
	if vInter == nil {
		// 不存在时，使用系统设置的语言
		lang = langcmp.NewLangCmp().GetSysDefaultLang()
		return
	}

	if v, ok := vInter.(langcmp.Lang); ok {
		lang = v
	} else {
		panic("GetLangFromCtx:ctx.Value(enums.VisitLangCtxKey) is not langcmp.Lang")
	}

	return
}
