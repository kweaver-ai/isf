package locale

import (
	"context"

	"github.com/gin-gonic/gin"

	"AuditLog/common/conf"
	"AuditLog/common/helpers"
	"AuditLog/infra/cmp/langcmp"
)

// getI18nWithDefault 获取带默认值的国际化信息
func getI18nWithDefault(ctx context.Context, key string, i18nMap map[string]map[langcmp.Lang]string) string {
	var lang langcmp.Lang
	if ctx != nil {
		ginCtx, ok := ctx.(*gin.Context)
		if ok {
			lang = helpers.GetLangFromHeader(ginCtx)
		} else {
			lang = langcmp.NewLangCmp().GetSysDefaultLang()
		}
	} else {
		lang = langcmp.NewLangCmp().GetSysDefaultLang()
	}
	return i18nMap[key][lang]
}

// GetI18nCtx 获取国际化信息
func GetI18nCtx(ctx context.Context, key string) string {
	return getI18nWithDefault(ctx, key, i18ns)
}

// GetRCLogLevelI18n 获取报表中心日志等级国际化信息
func GetRCLogLevelI18n(ctx context.Context, key string) string {
	return getI18nWithDefault(ctx, key, rcLevelI18n)
}

// getI18nWithDefault 获取带默认值的国际化信息
func getI18nWithDefaultByInt(ctx context.Context, key int, i18nMap map[int]map[langcmp.Lang]string) string {
	var lang langcmp.Lang
	if ctx != nil {
		ginCtx, ok := ctx.(*gin.Context)
		if ok {
			lang = helpers.GetLangFromHeader(ginCtx)
		} else {
			lang = langcmp.NewLangCmp().GetSysDefaultLang()
		}
	} else {
		lang = langcmp.NewLangCmp().GetSysDefaultLang()
	}
	return i18nMap[key][lang]
}

// GetRCLogObjTypeI18n 获取报表中心日志对象类型国际化信息
func GetRCLogObjTypeI18n(ctx context.Context, key int) string {
	return getI18nWithDefaultByInt(ctx, key, conf.MapObjectTypeLang)
}

// GetRCLogLoginI18n 获取报表中心登录日志国际化信息
func GetRCLogLoginI18n(ctx context.Context, key int) string {
	return getI18nWithDefaultByInt(ctx, key, conf.MapLoginOperTypeLang)
}

// GetRCLogMgntI18n 获取报表中心管理日志国际化信息
func GetRCLogMgntI18n(ctx context.Context, key int) string {
	return getI18nWithDefaultByInt(ctx, key, conf.MapManageOperTypeLang)
}

// GetRCLogOpI18n 获取报表中心操作日志国际化信息
func GetRCLogOpI18n(ctx context.Context, key int) string {
	return getI18nWithDefaultByInt(ctx, key, conf.MapOperOperTypeLang)
}
