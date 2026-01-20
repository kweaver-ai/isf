package common

import (
	"github.com/kweaver-ai/go-lib/rest"
)

//nolint:nolintlint
const (
	// InvalidAccountORPassword 账户密码错误
	InvalidAccountORPassword int = 401020108
	// UserDisabled 账户被禁用
	UserDisabled int = 401020109
	// ThirdPartyAuthNotOpen 第三方验证没有开启
	ThirdPartyAuthNotOpen int = 401020607
	// // NeedThirdOAuth 需要通过第三方认证网站进行登录
	// NeedThirdOAuth = 401020610
	// ForbiddenLogin 禁止登录
	ForbiddenLogin int = 401020414
	// DomainNotExist 域不存在
	DomainNotExist int = 401020402
	// DomainDisabled 用户域被禁用
	DomainDisabled int = 401020413
	// DomainUserNotExist 域用户不存在
	DomainUserNotExist int = 401020404
	// DomainServerUnavailable 连接域服务器失败
	DomainServerUnavailable int = 401020411
	// // ProductNotAuthorized 产品尚未授权
	// ProductNotAuthorized = 401020519
	// // ProductHasExpired 产品授权已过期
	// ProductHasExpired = 401020524
	// PasswordExpire 密码过期
	PasswordExpire int = 401020127
	// PasswordNotSafe 密码不符合强密码要求
	PasswordNotSafe int = 401020128
	// PasswordISInitial 初始密码
	PasswordISInitial int = 401020129
	// // PWDFirstFailed 第一次密码错误
	// PWDFirstFailed = 401020131
	// // PWDSecondFailed 第二次密码错误
	// PWDSecondFailed = 401020132
	// PWDThirdFailed 到达密码错误最大次数，帐号将被锁定(非涉密模式)
	PWDThirdFailed = 401020135
	// AccountLocked 账户被锁定
	AccountLocked int = 401020130
	// CannotConnectThirdPartyServer 不能连接到第三方认证服务
	CannotConnectThirdPartyServer int = 401020608
	// UserNotExist 用户不存在
	UserNotExist int = 401020110
	// // WrongPassord 用户密码错误
	// WrongPassord = 401020134
	// ControledPasswordExpire 密码管控状态下，密码过期
	ControledPasswordExpire int = 401020137
	// // CannotLoginSlaveSite 不能登录分站点
	// CannotLoginSlaveSite = 401201239
	// ImageVCodeISNULL 验证码为空
	ImageVCodeISNULL int = 401020149
	// ImageVCodeTimeout 验证码超时
	ImageVCodeTimeout int = 401020150
	// ImageVCodeISWrong 验证码出错
	ImageVCodeISWrong int = 401020151
	// // InsufficientSystemResources 系统资源不足
	// InsufficientSystemResources = 401022211
	// UserNotActivate 账户未激活
	UserNotActivate int = 401022602
	// OTPWrong 动态密码错误
	OTPWrong int = 401020162
	// OTPTimeout 动态密码已过期
	OTPTimeout int = 401020166
	// OTPTooManyWrongTime 动态错误次数过多
	OTPTooManyWrongTime int = 401020167
	// ImageVCodeMoreThanTheLimie 验证码输入已达到限定次数
	ImageVCodeMoreThanTheLimie int = 401020161
	// MFAOTPServerError MFA动态密码服务器异常
	MFAOTPServerError int = 401020163
	// ThirdPluginInterError 第三方插件内部错误
	ThirdPluginInterError int = 401020617
	// // MFAConfigError 多因子认证配置错误
	// MFAConfigError = 401020618
	// FailedThirdConfig 第三方认证配置错误
	FailedThirdConfig int = 401020602
)

var (
	// ErrorI18n 错误码国际化
	ErrorI18n = map[int]map[string]string{
		InvalidAccountORPassword: {
			rest.Languages[0]: "用户名或密码不正确",
			rest.Languages[1]: "使用者名稱或密碼不正確",
			rest.Languages[2]: "invalid account or password",
		},
		UserDisabled: {
			rest.Languages[0]: "用户已禁用，请联系管理员。",
			rest.Languages[1]: "使用者已停用，請聯繫管理員。",
			rest.Languages[2]: "Your account has been disabled, please contact admin.",
		},
		ThirdPartyAuthNotOpen: {
			rest.Languages[0]: "第三方认证功能未开启",
			rest.Languages[1]: "協力廠商認證功能未開啟",
			rest.Languages[2]: "Third-party authentication is not open",
		},
		// NeedThirdOAuth: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		ForbiddenLogin: {
			rest.Languages[0]: "用户没有权限登录，请确认",
			rest.Languages[1]: "使用者沒有登入權限，請確認",
			rest.Languages[2]: "This user is not allowed to login. Please check.",
		},
		DomainNotExist: {
			rest.Languages[0]: "域用户所属的域控不存在",
			rest.Languages[1]: "網域使用者所屬的網域控制站不存在",
			rest.Languages[2]: "User domain controller does not exist",
		},
		DomainDisabled: {
			rest.Languages[0]: "域用户所属的域已被禁用",
			rest.Languages[1]: "網域使用者所屬的網域已被停用",
			rest.Languages[2]: "User domain controller has been disabled",
		},
		DomainUserNotExist: {
			rest.Languages[0]: "当前的用户已不存在于LDAP服务器中",
			rest.Languages[1]: "當前的使用者已不存在於LDAP伺服器中",
			rest.Languages[2]: "This user does not exist in LDAP server.",
		},
		DomainServerUnavailable: {
			rest.Languages[0]: "连接LDAP服务器失败，请检查域控ip是否正确，或者域控制器是否开启！",
			rest.Languages[1]: "連接LDAP伺服器失敗，請檢查網域控制站ip是否正確，或者網域控制站是否開啟！",
			rest.Languages[2]: "Connect LDAP server failure. Please check domain controller ip or enable domain controller.",
		},
		// ProductNotAuthorized: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		// ProductHasExpired: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		PasswordExpire: {
			rest.Languages[0]: "您的登录密码已失效，是否立即修改密码？",
			rest.Languages[1]: "您的登入密碼已失效，是否立即變更密碼？",
			rest.Languages[2]: "Password expired, change password now?",
		},
		PasswordNotSafe: {
			rest.Languages[0]: "您的密码安全系数过低，是否立即修改密码？",
			rest.Languages[1]: "您的密碼安全係數過低，是否立即變更密碼？",
			rest.Languages[2]: "Weak password, change password now?",
		},
		PasswordISInitial: {
			rest.Languages[0]: "密码是初始密码",
			rest.Languages[1]: "密碼是初始密碼",
			rest.Languages[2]: "password is initial password",
		},
		// PWDFirstFailed: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		// PWDSecondFailed: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		// PWDThirdFailed: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		AccountLocked: {
			rest.Languages[0]: "账号已被锁定，请%s分钟后重试",
			rest.Languages[1]: "帳戶已被鎖定，請%s分鐘後重試",
			rest.Languages[2]: "Your account has been locked, please try again %s minutes later.",
		},
		CannotConnectThirdPartyServer: {
			rest.Languages[0]: "不能连接到第三方认证服务",
			rest.Languages[1]: "不能連接到協力廠商認證服務",
			rest.Languages[2]: "Connect to third-party authentication server failure",
		},
		UserNotExist: {
			rest.Languages[0]: "用户不存在",
			rest.Languages[1]: "使用者不存在",
			rest.Languages[2]: "User does not exist",
		},
		// WrongPassord: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		ControledPasswordExpire: {
			rest.Languages[0]: "你的密码已过期, 请联系管理员",
			rest.Languages[1]: "你的密碼已過期, 請聯繫管理員",
			rest.Languages[2]: "Password expired, please contact admin.",
		},
		// CannotLoginSlaveSite: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		ImageVCodeISNULL: {
			rest.Languages[0]: "请输入验证码",
			rest.Languages[1]: "請輸入驗證碼",
			rest.Languages[2]: "Please enter verification code",
		},
		ImageVCodeTimeout: {
			rest.Languages[0]: "您输入的验证码已失效",
			rest.Languages[1]: "您輸入的驗證碼已失效",
			rest.Languages[2]: "Invalid verification code",
		},
		ImageVCodeISWrong: {
			rest.Languages[0]: "您输入的验证码有误",
			rest.Languages[1]: "您輸入的驗證碼有誤",
			rest.Languages[2]: "Incorrect verification code",
		},
		// InsufficientSystemResources: {
		// 	rest.Languages[0]: "用户已禁用，请联系管理员。",
		// 	rest.Languages[1]: "使用者已停用，請聯繫管理員。",
		// 	rest.Languages[2]: "Your account has been disabled, please contact admin.",
		// },
		UserNotActivate: {
			rest.Languages[0]: "您的账号已被禁用，是否立即激活？",
			rest.Languages[1]: "您的帳戶已被停用，是否立即啟動",
			rest.Languages[2]: "User has been disabled，Whether to activate immediately？",
		},
		OTPWrong: {
			rest.Languages[0]: "动态密码错误",
			rest.Languages[1]: "動態密碼錯誤",
			rest.Languages[2]: "Wrong one time password.",
		},
		OTPTimeout: {
			rest.Languages[0]: "动态密码错误",
			rest.Languages[1]: "動態密碼錯誤",
			rest.Languages[2]: "Wrong one time password.",
		},
		OTPTooManyWrongTime: {
			rest.Languages[0]: "动态密码错误",
			rest.Languages[1]: "動態密碼錯誤",
			rest.Languages[2]: "Wrong one time password.",
		},
		ImageVCodeMoreThanTheLimie: {
			rest.Languages[0]: "验证码输入已达到限定次数",
			rest.Languages[1]: "驗證碼輸入已達到限定次數",
			rest.Languages[2]: "Too many input attempts",
		},
		MFAOTPServerError: {
			rest.Languages[0]: "不能连接到第三方认证服务",
			rest.Languages[1]: "不能連接到協力廠商認證服務",
			rest.Languages[2]: "Connect to third-party authentication server failure",
		},
		ThirdPluginInterError: {
			rest.Languages[0]: "插件内部错误: %s",
			rest.Languages[1]: "外掛程式內部錯誤: %s",
			rest.Languages[2]: "Third plugin inter error: %s.",
		},
		// MFAConfigError: {
		// 	rest.Languages[0]: "您输入的验证码有误",
		// 	rest.Languages[1]: "您輸入的驗證碼有誤",
		// 	rest.Languages[2]: "Incorrect verification code",
		// },
		FailedThirdConfig: {
			rest.Languages[0]: "第三方认证模块加载失败: %s",
			rest.Languages[1]: "協力廠商認證模組載入失敗: %s",
			rest.Languages[2]: "Failed to invoke the third-party authentication tool: %s",
		},
	}
)
