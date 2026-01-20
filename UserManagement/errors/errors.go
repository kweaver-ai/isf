// Package errors 服务错误
package errors

import (
	"github.com/kweaver-ai/go-lib/rest"
)

// 服务错误码
const (
	strPrefix                 = "UserManagement."
	UserNotFound              = 400019001
	StrBadRequestUserNotFound = strPrefix + "BadRequest" + ".UserNotFound"

	DepartmentNotFound              = 400019002
	StrBadRequestDepartmentNotFound = strPrefix + "BadRequest" + ".DepartmentNotFound"

	GroupNotFound              = 400019003
	StrBadRequestGroupNotFound = strPrefix + "BadRequest" + ".GroupNotFound"

	ContactorNotFound = 400019004

	AppNotFound              = 400019005
	StrBadRequestAppNotFound = strPrefix + "BadRequest" + ".AppNotFound"

	Forbidden = 403019001

	NotFound                 = 404019001
	StrNotFoundAppNotFound   = strPrefix + "NotFound" + ".AppNotFound"
	StrNotFoundGroupNotFound = strPrefix + "NotFound" + ".GroupNotFound"

	Conflict         = 409019001
	StrConflictApp   = strPrefix + "Conflict" + ".AppConflict"
	StrConflictGroup = strPrefix + "Conflict" + ".GroupConflict"

	// 匿名账户相关错误码使用sharedlink 007，与之前的逻辑一致
	AnonymityNotFound        = 400007001
	AnonymityWrongPassword   = 403007001
	AnonymityReachLimitTimes = 403007002
)

var (
	errorI18n = map[int]map[string]string{
		UserNotFound: {
			rest.Languages[0]: "用户不存在",
			rest.Languages[1]: "使用者不存在",
			rest.Languages[2]: "This user does not exist",
		},
		DepartmentNotFound: {
			rest.Languages[0]: "部门不存在",
			rest.Languages[1]: "部門不存在",
			rest.Languages[2]: "This department does not exist",
		},
		ContactorNotFound: {
			rest.Languages[0]: "联系人组不存在",
			rest.Languages[1]: "联系人组不存在",
			rest.Languages[2]: "This contactor group does not exist",
		},
		GroupNotFound: {
			rest.Languages[0]: "用户组不存在",
			rest.Languages[1]: "用戶組不存在",
			rest.Languages[2]: "This group does not exist",
		},
		NotFound: {
			rest.Languages[0]: "数据不存在",
			rest.Languages[1]: "数据不存在",
			rest.Languages[2]: "Not Found",
		},
		Conflict: {
			rest.Languages[0]: "数据已存在",
			rest.Languages[1]: "数据已存在",
			rest.Languages[2]: "Conflict",
		},
		Forbidden: {
			rest.Languages[0]: "禁止访问",
			rest.Languages[1]: "禁止訪問",
			rest.Languages[2]: "Forbidden",
		},
		AnonymityWrongPassword: {
			rest.Languages[0]: "访问密码不正确",
			rest.Languages[1]: "存取密碼錯誤",
			rest.Languages[2]: "Wrong password",
		},
		AnonymityReachLimitTimes: {
			rest.Languages[0]: "访问次数已达上限",
			rest.Languages[1]: "存取次數已達上限",
			rest.Languages[2]: "Wrong password",
		},
		AnonymityNotFound: {
			rest.Languages[0]: "数据不存在",
			rest.Languages[1]: "数据不存在",
			rest.Languages[2]: "Not Found",
		},
	}
)

func init() {
	rest.Register(errorI18n)
}
