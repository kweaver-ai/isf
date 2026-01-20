/**
 * oauth认证错误状态码
 */
export enum ErrorCode {
    /**
     * 无异常
     */
    Normal,

    /**
     * 无账号
     */
    NoAccount,

    /**
     * 无密码
     */
    NoPassword,

    /**
     * 无图形验证码
     */
    NoCaptcha,

    /**
     * 无动态密码
     */
    NoDynamicPassword,

    /**
     * 无旧密码
     */
    NoOldPassword,

    /**
     * 无新密码
     */
    NoNewPassword,

    /**
     * 无确认密码
     */
    NoConfirmPassword,

    /**
     * 新密码与确认密码不一致
     */
    NewConfirmInconsitent,

    /**
     * 新密码等于初始密码
     */
    NewIsInitial,

    /**
     * 新密码等于旧密码
     */
    NewIsOld,

    /**
     * 旧密码不正确
     */
    OldPasswordInvalid,

    /**
     * 手机号或邮箱为空
     */
    BlankNumber,

    /**
     * 手机号为空
     */
    BlankPhone,

    /**
     * 邮箱为空
     */
    BlankEmail,

    /**
     * 手机号或邮箱格式不正确
     */
    InCorrect,

    /**
     * 手机号有误
     */
    CellphoneError,

    /**
     * 邮箱有误
     */
    EmailError,

    /**
     * 发送失败，管理员已关闭忘记密码重置功能
     */
    CloseForgetPasswordResetBySend,

    /**
     * 未启用短信验证
     */
    SMSClose,

    /**
     * 未启用邮箱验证
     */
    EmailClose,

    /**
     * 手机号绑定的用户非本地用户
     */
    SMSUserNoLocal,

    /**
     * 邮箱绑定的用户非本地用户
     */
    EmailUserNoLocal,

    /**
     * 手机号绑定的用户被管控（不允许自主修改密码）
     */
    SMSUserControlled,

    /**
     * 邮箱绑定的用户被管控（不允许自主修改密码）
     */
    EmailUserControlled,

    /**
     * 旧密码错误次数多，账号锁定
     */
    OldPasswordInvalidLocked,

    /**
     * 短信验证认证时，账号或密码错误
     */
    PasswordChange,

    /**
     * 网络未连接
     */
    NoNetwork,

    /**
     * 无效账户
     */
    invalid_account,

    /**
     * 用户被禁用
     */
    disable_user,

    /**
     * 密码找回功能未开启
     */
    unable_pwd_retrieval,

    /**
     * 非本地用户
     */
    non_local_user,

    /**
     * 密码管控开启
     */
    enable_pwd_control,

    /**
     * 用户未绑定手机号
     */
    UNBOUND_PHONENUMBER,

    /**
     * 用户未绑定邮箱
     */
    UNBOUND_EMAIL,

    /**
     * 用户未绑定手机号或邮箱
     */
    UNBOUND_PHONENUMBER_AND_EMAIL,

    /**
     * 未输入用户账号
     */
    NOACCOUNT,

    /**
     * 未配置第三方参数
     */
    MissSignOutConfig = 404001111,

    /**
     * 不支持的接口
     */
    URINotExists = 400002001,

    /**
     * 非法参数
     */
    ParametersIllegal = 400000000,

    /**
     * 非法uri参数
     */
    URIFormatIllegal = 400000000,

    /**
     * 非法json字符串
     */
    JSONFormatIllegal = 400000000,

    /**
     * 非法权限值类型
     */
    PermConfigIllegal = 400001005,

    /**
     * 非法访问者
     */
    AccessorIllegal = 400006,

    /**
     * 登录超时，access_token校验失败
     */
    TokenExpired = 401001001,

    /**
     * 账号或密码不正确
     */
    AuthFailed = 401001003,

    /**
     * 账号被禁用
     */
    UserDisabled = 401001004,

    /**
     * 加密凭证无效
     */
    EncryptionInvalid = 401002005,

    /**
     * 普通用户角色，无法登录控制台
     */
    UserLoginFailed = 401001005,

    /**
     * 管理员角色，无法登录客户端
     */
    AdminLoginFailed = 401001006,

    /**
     * 域控被禁用或删除
     */
    DomainDisabled = 401001007,

    /**
     * 用户不存在于域控
     */
    UserNotInDomain = 401001008,

    /**
     * 登录设备被禁用
     */
    DeviceDisaled = 401001009,

    /**
     * 账号已绑定设备，无法登录web客户端
     */
    DeviceBinded = 401001011,

    /**
     * 密码已失效，非管控
     */
    PasswordFailure = 401001012,

    /**
     * 密码系数过低
     */
    PasswordInsecure = 401001013,

    /**
     * 弱密码格式错误
     */
    PasswordInvalid = 401001014,

    /**
     * 强密码格式错误
     */
    PasswordWeak = 401001015,

    /**
     * 外部用户不支持修改密码
     */
    PasswordChangeNotSupported = 401001016,

    /**
     * 初始密码登录
     */
    PasswordIsInitial = 401001017,

    /**
     * 密码输错一次
     */
    PasswordInvalidOnce = 401001018,

    /**
     * 密码输错两次
     */
    PasswordInvalidTwice = 401001019,

    /**
     * 非涉密模式，账号被锁定
     */
    PasswordInvalidLocked = 401001020,

    /**
     * 产品未许可
     */
    LicenseInvalid = 401001021,

    /**
     * 连接LDAP服务器失败
     */
    DomainServerUnavailable = 401001022,

    /**
     * 当前账号在另一地点登录，被迫下线
     */
    AccountDuplicatedLogin = 401001025,

    /**
     * 无修改密码权限
     */
    PasswordRestricted = 401001026,

    /**
     * 密码过期，管控
     */
    PasswordExpired = 401001027,

    /**
     * 分站点模式，无法登录
     */
    LoginSiteInvalid = 401001028,

    /**
     * 产品过期
     */
    LicenseExpired = 401001029,

    /**
     * 未完成初始化配置
     */
    Uninitialized = 401001030,

    /**
     * ip网段限制，无法登录
     */
    IPRestricted = 401001031,

    /**
     * 涉密模式，账号被锁定
     */
    AccountLocked = 401032,

    /**
     * 管理员禁止此类客户端登录
     */
    ClientRestricted = 401001033,

    /**
     * 新密码为初始密码
     */
    NewPasswordIsInitial = 401001035,

    /**
     * 网络环境改变
     */
    NetworkChanged = 401001036,

    /**
     * 验证码为空
     */
    VCodeMissing = 401001037,

    /**
     * 验证码过期
     */
    VCodeExpired = 401001038,

    /**
     * 验证码错误
     */
    VCodeInvalid = 401001039,

    /**
     * 用户禁用，请激活
     */
    NeedAction = 401001040,

    /**
     * 用户激活，请登录
     */
    UserActivated = 401001041,

    /**
     * 手机号不合法
     */
    PhoneNumberInvalid = 401001042,

    /**
     * 短信验证码不正确
     */
    CaprchaWrong = 401001044,

    /**
     * 短信验证码已过期
     */
    CaprchaOverstayed = 401001045,

    /**
     * 发送验证码失败
     */
    SendCaprchaFail = 401001046,

    /**
     * 短信激活未开启
     */
    NotOpenActivated = 401001047,

    /**
     * 邮箱不合法
     */
    EmailCorrectFormat = 401001048,

    /**
     * 激活失败
     */
    FailActivated = 401001050,

    /**
     * 关闭身份证号登录
     */
    IDAuthDisabled = 401001051,

    /**
     * 邮箱未绑定
     */
    UnboundEmail = 401001054,

    /**
     * 手机号未绑定
     */
    UnboundPhone = 401001055,

    /**
     * 验证码错误次数过多
     */
    VcodeErrorTimesTooMany = 401001056,

    /**
     * 动态密码为空
     */
    OTPRequest = 401001057,

    /**
     * 未绑定手机
     */
    PhoneNotUnbound = 401001058,

    /**
     * 动态密码服务异常
     */
    OTPServerExceptions = 401001059,

    /**
     * 手机号变更
     */
    PhoneModified = 401001060,

    /**
     * 动态密码不正确
     */
    OTPInvalid = 401001061,

    /**
     * 动态密码过期
     */
    OTPExpired = 401001062,

    /**
     * 动态密码错误次数过多
     */
    OTPErrorTimesTooMany = 401001063,

    /**
     * 短信服务器异常
     */
    SMSServerExceptions = 401001064,

    /**
     * 第三方服务内部错误
     */
    ThirdServiceInternalError = 401001065,

    /**
     * 登录认证配置错误
     */
    AuthServerExceptions = 401001066,

    /**
     * 第三方认证模块导入失败
     */
    ThridServerImportError = 401001067,

    /**
     * 第三方认证未开启
     */
    CASDisabled = 403001028,

    /**
     * 无法验证ticket
     */
    TicketInvalid = 403001029,

    /**
     * 用户未导入anyshare
     */
    UserNotFound = 403001030,

    /**
     * 第三方服务器认证失败
     */
    ThirdPartyValidateFail = 403001057,

    /**
     * 登录外部应用失败
     */
    ExtLoginFailed = 403001106,

    /**
     * 用户被冻结
     */
    AccountFrozen = 403001171,

    /**
     * 邮箱地址非法
     */
    EmailInvalid = 403001185,

    /**
     * SMTP服务器未设置
     */
    SMTPConfigMissing = 404001021,

    /**
     * SMTP服务器存在未知错误
     */
    SMTPUnknownError = 404001022,

    /**
     * SMTP服务器不可用
     */
    SMTPInaccessible = 404001023,

    /**
     * 文档域资源不足
     */
    InsufficientDomainResources = 404001027,

    /**
     * 发送验证码服务器未开启
     */
    SendVcodeServerUnavailable = 404001028,

    /**
     * 短信服务器未设置
     */
    SMSPConfigMissing = 404001029,

    /**
     * 动态密码服务器未设置
     */
    OTPServerNotEnabled = 404001030,

    /**
     * 短信服务器未正常启用
     */
    SMSServerNotEnabled = 404001031,

    /**
     * http方法错误
     */
    HTTPMethodError = 405002001,

    /**
     * 内部错误
     */
    InternalError = 500000000,

    /**
     * 服务器版本不支持该客户端
     */
    ServerClientMismatch = 500001007,

    /**
     * 不支持post外其它方法
     */
    HTTPNotPOST = 501001,

    /**
     * 服务器繁忙
     */
    ServiceBusy = 503001,

    /**
     * csrf验证失败
     */
    CsrfFailed = 403041000,

    /**
     * 连接eacp服务失败
     */
    ConnectEacpFailed = 500041001,

    /**
     * 连接hydra服务失败
     */
    ConnectHydraFailed = 500041002,

    /**
     * 连接authentication服务失败
     */
    ConnectAuthenticationFailed = 500041003,

    /**
     * 参数不合法,缺少code参数
     */
    INVALID_NO_CODE = 400041000,

    /**
     * 参数不合法,缺少consent_challenge参数
     */
    INVALID_NO_CONSENT_CHALLENGE = 400041001,

    /**
     * 参数不合法,缺少logout_challenge参数
     */
    INVALID_NO_LOGOUT_CHALLENGE = 400041002,

    /**
     * 参数不合法,state验证失败，可能cookie与querystring中state值不一致、state值不是以/开始的字符串等等
     */
    INVALID_STATE = 400041003,

    /**
     * 参数不合法,challenge或remember参数验证失败
     */
    INVALID_CHALLENGE_OR_REMEMBER = 400041004,

    /**
     * 用户不存在
     */
    USER_DOES_NOT_EXIST = 401001023,

    /**
     * hydra服务异常处理
     */
    Hydra400 = 400,
    Hydra401 = 401,
    Hydra403 = 403,
    Hydra404 = 404,
    Hydra409 = 409,
    Hydra500 = 500,
    Hydra503 = 503,
}

/**
 * 错误信息提示
 * @param errorStatus 错误状态码
 * @param t 全球化函数
 * @param value 提示插入值
 * @param defaultMessage 默认的错误提示语
 */
export const getErrorMessage = (
    errorStatus: ErrorCode,
    t: any,
    value?: { [key: string]: string },
    defaultMessage?: string
): string => {
    switch (errorStatus) {
        case ErrorCode.PasswordWeak:
        case ErrorCode.PasswordInvalidLocked:
        case ErrorCode.OldPasswordInvalidLocked:
            return t(`signin-error-${errorStatus}`, { ...value });
        default:
            return typeof ErrorCode[errorStatus] === "string"
                ? t(`signin-error-${errorStatus}`)
                : defaultMessage
                ? defaultMessage
                : t("unknown");
    }
};
