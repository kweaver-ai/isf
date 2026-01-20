package errors

import (
	"fmt"

	"AuditLog/infra/cmp/langcmp"
)

type I18nMap map[int]struct {
	Description map[langcmp.Lang]string
	Solution    map[langcmp.Lang]string
}

// i18ns 是一个示例的国际化错误信息映射表
var i18ns = I18nMap{
	BadRequestErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "请求错误，请稍后重试。",
			langcmp.ZhTW: "請求錯誤，請稍後重試。",
			langcmp.En:   "Bad request. You can try again later.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	UnauthorizedErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "授权错误，请检查您的用户名或密码后重新登录。",
			langcmp.ZhTW: "授權錯誤，請檢查您的使用者名稱或密碼后重新登入。",
			langcmp.En:   "Unauthorized. Please check your username or password and try again.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	ForbiddenErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "您没有权限访问当前资源。",
			langcmp.ZhTW: "您沒有權限存取當前資源。",
			langcmp.En:   "You are not allowed to access this resource.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	ResourceNotFoundErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "当前资源已不存在，请稍后重试或联系管理员。",
			langcmp.ZhTW: "當前資源已不存在，請稍後重試或聯絡管理員。",
			langcmp.En:   "This resource doesn't exist. You can try again later or contact Admin.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	MethodNotAllowedErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "请求方法错误。请稍后重试。如果长时间仍无反应，您可以联系管理员。",
			langcmp.ZhTW: "請求方法錯誤。請稍後重試。如果長時間仍無反應，您可以聯絡管理員。",
			langcmp.En:   "Method error. You can try again later or contact Admin.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	ConflictErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "当前服务器请求冲突，您可以稍后重试。",
			langcmp.ZhTW: "當前伺服器請求衝突，您可以稍後重試。",
			langcmp.En:   "Conflict detected in the server. You can try again later.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	TooManyRequestsErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "请求过多，您可以稍后重试。",
			langcmp.ZhTW: "請求過多，您可以稍後重試。",
			langcmp.En:   "Too many requests. You can try again later.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	InternalErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "当前服务器存在内部错误或错误配置",
			langcmp.ZhTW: "當前伺服器存在內部錯誤或錯誤設定",
			langcmp.En:   "An internal error or wrong configuration has been detected",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	ScopeStrategyConflictErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "该策略已存在。",
			langcmp.ZhTW: "該策略已存在。",
			langcmp.En:   "This policy already exists.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	PasswordRequiredErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "密码不允许为空。",
			langcmp.ZhTW: "密碼不允許為空。",
			langcmp.En:   "Password is required.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	PasswordInvalidErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "密码必须同时包含 大小写英文字母 与 数字，允许包含 ~!%#$@-_. 字符，长度范围 10~100 个字符，请重新输入",
			langcmp.ZhTW: "密碼必須同時包含 大小寫英文字母與 數位，允許包含 ~!%#$@-_. 字元，長度範圍10~100個字元，請重新輸入",
			langcmp.En:   "Password must contain numbers, letters in both uppercase and lowercase within 10~100 characters, allowing for special characters( ~!%#$@-_. ), please re-enter",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
	ScopeStrategyNotFoundErr: {
		Description: map[langcmp.Lang]string{
			langcmp.ZhCN: "该日志策略已不存在。",
			langcmp.ZhTW: "該日誌策略已不存在。",
			langcmp.En:   "The policy does not exist.",
		},
		Solution: map[langcmp.Lang]string{
			langcmp.ZhCN: "",
			langcmp.ZhTW: "",
			langcmp.En:   "",
		},
	},
}

func RegisterI18ns(i18nMap I18nMap) {
	for k, v := range i18nMap {

		if _, ok := i18ns[k]; ok {
			panic(fmt.Sprintf("error code %d already exists", k))
		}

		i18ns[k] = v
	}
}
