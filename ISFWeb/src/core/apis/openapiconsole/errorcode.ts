/**
 * 错误码定义
 * 命名规则：操作／对象 + 状态
 * 常用状态表示：
 * * 非法值 Illegal
 * * 无效值 Invalid
 * * 操作对象不存在（被删除） Inaccessible
 * * 操作对象缺失（未配置） Missing
 * * 操作对象冲突 Conflict
 * * 超出允许范围 Exceeded
 * * 操作被限制 Restricted
 * * 操作对象不完整 Incomplete
 */
export const enum ErrorCode {

	/**
	 * 请求参数错误
	 */
	InvalidRequest = 400000000,

	/**
	 * 对象不存在（或已被删除）
	 */
	ResourceInaccessible = 404014000,

	/**
	 * 资源冲突（对象已存在）
	 */
	ResourceConflict = 409014000,

	/**
	 * 文档域连接失败
	 */
	DomainLinkFalied = 400014205,

	/**
	 * 策略绑定了文档域
	 */
	PolicyBoundDomain = 409014221,

	/**
	 * 不能添加指定的域为子域/平级域
	 */
	DomainAddFailed = 409014201,

	/**
	 * 该文档域已是父域，不允许添加为子域
	 */
	ParentCannotBeAddedAsChild = 409014202,

	/**
	 * 该文档域已被'xx'添加为子域，不允许再次被添加为子域
	 */
	ChildCannotBeAddedAsChild = 409014203,

	/*
	 * 当前域为子域，不能添加子域或添加其父域为平级域
	 */
	DomainIsChild = 409014204,

	/**
	 * 编目模板不存在
	 */
	MetaDataTemplateNotExist = 404003032,

	/**
	 * 编目模板已存在
	 */
	MetaDataTemplateExist = 403003206,

	/**
	 * 编目模板key已存在
	 */
	MetaDataTemplateKeyExist = 409003001,

	/**
	 * 编目模板个数已达上限
	 */
	MetaDataTemplateNubLimit = 403003209,

	/**
	 * 编目模板属性不存在
	 */
	MetaDataFieldNotExist = 404003035,

	/**
	 * 编目模板选项不存在
	 */
	MetaDataOptionNotExist = 404003036,

	/**
	 * 编目模板属性/选项拖拽多个不存在
	 */
	MetaDataFieldOrOptionNotExist = 403003217,

	/**
	 * 编目模板属性超过限制条数
	 */
	MetaDataFieldNumLimit = 403003211,

	/**
	 * 编目模板选项超过限制条数
	 */
	MetaDataOptionNumLimit = 403003212,

	/**
	 * 编目模板属性已存在
	 */
	MetaDataFieldExist = 403003215,

	/**
	 * 编目模板选项已存在
	 */
	MetaDataOptionExist = 403003207,

	/**
	 * 对象不存在（或已被删除）
	 */
	ResourceInaccessibleByPolicy = 404013000,

	/**
	 * 资源冲突（对象已存在）
	 */
	ResourceConflictByPolicy = 409013000,

	/**
	 * 预览/下载策略 资源冲突（对象已存在）
	 */
	ResourceConflictByReadPolicy = 409029000,

	/**
	 * 请求过多
	 */
	TooManyRequestsByPolicy = 429013000,

	/**
	 * 无权限进行此操作
	 */
	NoPermissionToOperateByPolicy = 403013000,

	/**
	 * 服务不可用
	 */
	ServerNotAvailable = 502,

	/**
	 * 内部错误
	 */
	InternalError = 500000000,

	/**
	 * 编目属性至少存在一个
	 */
	FieldAtLeastOne = 403003213,

	/**
	 * 编目选项至少存在一个
	 */
	OptionAtLeastOne = 403003214,

	/**
	 * 已存在关联文档库
	 */
	DocAlreadyExist = 409002002,

	/**
	 * 文档库存在于系统回收站
	 */
	DocInRecycle = 409002006,

	/**
	 * 部门文档库存在于系统回收站
	 */
	DepDocInRecycle = 409002005,

	/**
	 * 文档库不存在
	 */
	DocNotExist = 404002005,

	/**
	 * 用户不存在
	 */
	UserNotExist = 400019001,

	/**
	 * 部门不存在
	 */
	DepartmentsNotExist = 400019002,

	/**
	 * 用户组不存在
	 */
	UserGruopsNotExist = 400019003,

	/**
	 * 应用账户不存在
	 */
	UseApplicationAccountNotExist = 400019005,

	/**
	 * 库名称与用户显示名同名
	 */
	DocNameSameAsUserName = 409002004,

	/**
	 * 存在同名的文档库
	 */
	DocNameConflict = 409002003,

	/**
	* 部门文档库已存在
	*/
	DepartmentDocLibExist = 409002001,

	/**
	 * 配额空间大于管理员可配置的空间
	 */
	QuotaGreaterThanAvailable = 403002001,

	/**
	 * 编辑的配额小与已占用的配额
	 */
	QuotaLessThanUsed = 403002154,

	/**
	 * 存储位置不可用
	 */
	StorageUnAvailable = 403002220,

	/**
	 * 存储位置不存在
	 */
	StorageNotExist = 400002016,

	/*
	 * token无效或过期
	 */
	TokenExpire = 401000000,

	/*
	 * token无效或过期
	 */
	TokenEmpty = 401001001,

	/**
	 * 登录策略密码输入错误限制
	 */
	InvalidSecretPasswdErr = 201242,

	/**
	 * 认证凭据类型与添加域类型不符
	 */
	CredentialNotAvailable = 409014205,

	/**
	 * 认证凭据失效或已被使用
	 */
	CredentialInvalidOrUsed = 409014206,

	/**
	 * 认证凭据失效
	 */
	CredentialInvalid = 401014201,

	/**
	 * 认证凭据已被删除
	 */
	CredentialNotExist = 403014000,

	/**
	 * 文档库组织管理员无权操作
	 */
	NoPermission = 403001002,

	/**
	 * 目标文档域不存在
	 */
	TargetDomainNotExist = 400014234,

	/**
	 * 源文档库不存在
	 */
	SourceDocLibNotExist = 400014235,

	/**
	 * 同步源一样，同步模式不一样
	 */
	SyncPatternConflict = 409014236,

	/**
	 * 同步过程中文档库不存在
	 */
	DocLibNotExist = 403002024,

	/**
	 * 配额不足
	 */
	InsufficientQuota = 403002104,

	/**
	 * 同步计划正在进行中, 无法删除
	 */
	SyncPlanInProcessing = 409014230,

	/**
	 * 点击立即同步的同步计划正在同步进行中
	 */
	ConflictResource = 409015000,

	/**
	 * 用户组不存在
	 */
	UserGroupNotExist = 404019001,

	/**
	 * 用户组名称已存在
	 */
	UserGroupNameExist = 409019001,

	/**
	 * 预览/下载策略资源不存在（或已被删除）
	 */
	ResourceNotExistByReadPolicy = 404029000,

	/**
	 * 策略不存在（索引策略，杀毒策略或者杀毒策略对应的文档库）
	 */
	StrategyNotExist = 404034000,

	/**
	 * 无权限处理特殊索引策略
	 */
	NoPermissionToHandle = 403034000,

	/**
	 * 审核流程失效
	 */
	WorkflowAuditInvalid = 400014301,

	/**
	 * 杀毒策略已存在
	 */
	AntivirusStrategyExist = 409034000,

	/**
	 * 库类型不存在
	 */
	NotExistDocLibType = 400002301,

	/**
	 * 库类型下存在文档库
	 */
	ExitDocLib = 400002302,

	/**
	* 库类型已存在
	*/
	ExistDocLibType = 409002301,

	/**
	 * 非法的第三方插件
	 */
	IllegalPlugin = 500045002,

	/**
	 * 第三方应用不存在
	 */
	AppNotExist = 400045002,

	/**
	 * 限速设置所选部门不存在
	 */
	DepInexistence = 400047003,

	/**
	 * osspolicy相关服务无法连接
	 */
	OssPolicyErr = 500047000,

	/**
	 * 被迁移的个人文档库不存在
	 */
	HandoverUsersDocLibNotExist = 400046001,

	/**
	 * 目标文档库不存在
	 */
	TargetDocLibNotExist = 400046002,

	/**
	 * 权限接收者不存在
	 */
	PermReceiverNotExist = 400046003,

	/**
	 * 目标文档库配额空间不足
	 */
	TargetDocLibOutOfFreeSpace = 409046001,

	/**
	 * 个人文档库迁移内部错误
	 */
	UsersDocLibHandoverErr = 500046001,

	/**
	 * 权限迁移内部错误
	 */
	PermHandoverErr = 500046002,

	/**
	 * 应用账户编辑用户交接权限时，应用账户已有该权限
	 */
	UseAccountUserTransferPermExist = 409046000,

	/**
	 * 应用账户删除用户交接权限时，该权限不存在
	 */
	UseAccountUserTransferPermNotExist = 404046000,

	/**
	  * 编辑文档库访问策略时，策略不存在
	  */
	AccessPolicyNotExist = 404001301,

	/**
	 * 应用账户已存在（获取OSS网关账户）
	 */
	AppUserExists = 403031022,

	/**
	 * 应用账户已删除（删除OSS网关账户）
	 */
	AppUserNoExists = 404031005,

	/**
	 * 适配者已存在适配模板（共享模板）
	 */
	UserRepat = 409001302,

	/*
	 * 审核流程失效
	 */
	WorkflowInvalid = 400002303,

	/**
	 * 文件或目录不存在
	 */
	DocumentNotExist = 404002006,

	/**
	 * 编目模板不存在
	 */
	TemplateNotExist = 404055162,

	/**
	 * 编目模板属性不存在
	 */
	FieldsNotExist = 404055163,

	/**
	 * 设置列表属性个性化文档库不存在
	 */
	DobLibsNotExist = 404055164,

	/**
	 * 操作模板不存在
	 */
	ActionTemplateNotExist = 404055171,

	/**
	 * 操作模板名称已存在
	 */
	ActionTemplateNameExist = 409055172,

	/**
	 * 操作配置不存在
	 */
	ActionConfigNotExist = 404055173,

	/**
	 * 操作配置名称已存在
	 */
	ActionConfigNameExist = 409055174,

	/**
	 * 选择的组织结构中的成员有不存在的
	 *     此处前端根据后端返回的错误码和不存在的成员类型和成员ID（detail中包含）来组装提示内容，不使用后端返回的description
	 *     接口返回中detail的格式：
	 *         detail:{"not_exists_user_ids":["user_id_xxx1","user_id_xxx2"],"not_exists_dep_ids":["dep_id_xxx1","dep_id_xxx2"],"not_exists_user_group_ids":["user_group_id_xxx1","user_group_id_xxx2"]}
	 */
	ActionConfigMenmberNotExist = 404055175,

	/**
	 * 同一库类型或文档库只允许基于一种应用规则配置策略
	 */
	ActionConfigScopeNotOnly = 409055176,

	/**
	 * 编目模板有不存在的
	 */
	ActionConfigMetadataTemplateNotExist = 404055178,

	/**
	 * 操作模板有不存在的
	 */
	ActionConfigTemplateNotExist = 404055180,

	/**
	 * 文档库有不存在的
	 */
	ActionConfigDocLibNotExist = 404055181,

	/**
	  * 文档库有不存在的
	  */
	AllDocLibApplicationExist = 409055182,

	/**
	 * 策略名称重复
	 */
	PolicyNameConflict = 409001301,

	/**
	 * 安全策略配置重叠（文档库和适用范围）
	 */
	PolicyResourceConflict = 409001303,

	/**
	 * 权限申请策略名称重复
	 */
	PermPolicyNameConflict = 400055020,

	/**
	 * 权限策略用户不存在
	 */
	PermUserNotExist = 400055021,

	/**
	 * 权限策略文档库不存在
	 */
	PermDoclibNotExist = 400055022,

	/**
	 * 权限申请策略资源冲突
	 */
	PermPolicyResourceConflict = 409000000,

	/**
	 * 资源/目标不存在 ||
	 * 应用账户删除文档库管理权限时，某个权限不存在
	 */
	ResourceNotExist = 404000000,

	/**
	 * 文档流转范围限制为空
	 */
	DocFlowLimitEmpty = 400014302,

	/**
	 * 开启默认同步计划(个人文档库收发件箱)时存在个人文档库实时同步计划
	 */
	UserDocLibLiveSyncExist = 409014237,

	/**
	 * 访问被拒绝
	 * 文档库权限接收方密级低于移交方，无法交接
	 */
	Forbidden = 403000000,

	/**
	 * 对方密级低,无法交接个人文档库
	 */
	HaveNoSufficientLevel = 403046001,

	/**
	 * 策略管控对象冲突
	 */
	StrategyResourceConflict = 409056001,

	/**
	 * 策略用户为空
	 */
	StrategyResourceEmpty = 404056001,

	/**
	 * 推荐配置场景分类不存在
	 */
	RecommendedSceneGroupNotExist = 404062201,

	/**
	 * 推荐配置场景分类名称重复
	 */
	RecommendedSceneGroupNameExist = 409062202,

	/**
	 * 推荐配置场景分类名称重复
	 */
	RecommendedSceneGroupUsed = 409062203,

	/**
	 * 推荐配置场景不存在
	 */
	RecommendedSceneNotExist = 404062211,

	/**
	 * 推荐配置场景标识重复
	 */
	RecommendedSceneKeyExist = 409062212,

	/**
	 * 推荐配置场景已被使用，无法删除
	 */
	RecommendedSceneUsed = 409062213,

	/**
	 * 推荐配置场景策略不存在
	 */
	RecommendedSceneStrategyNotExist = 404062221,

	/**
	 * 推荐配置场景策略标识重复
	 */
	RecommendedSceneStrategyKeyExist = 409062222,

	/**
	 * 推荐配置场景策略状态无效
	 */
	RecommendedSceneStrategyStatusInvalid = 400062403,

	/**
	 * 推荐配置场景策略配置无效
	 */
	RecommendedSceneStrategyConfigInvalid = 400062404,
}

export const enum PublicErrorCode {
    /**
     * 未授权或已过期
     */
    Unauthorized = 'Public.Unauthorized',

    /**
	 * 参数错误
	 */
	BadRequest = 'Public.BadRequest',

	/**
	 * 没有权限
	 */
	Forbidden = 'Public.Forbidden',

	/**
	 * 未找到
	 */
	NotFound = 'Public.NotFound',

	/**
	 * 冲突
	 */
	Conflict = 'Public.Conflict',

	/*
	 * 系统错误
	 */
	InternalServerError = 'Public.InternalServerError',

	/**
	 * 服务不可用
	 */
	ServiceUnavailable = 'Public.ServiceUnavailable',
}

export const enum UserManagementErrorCode {
    /*
	 * 用户组成员不存在
	 */
	GroupMemberNotExisted = 'UserManagement.BadRequest.UserNotFound',
	
	/**
	 * app不存在
	 */
	AppNotFound = 'UserManagement.NotFound.AppNotFound',

	/**
	 * groups用户组不存在
	 */
	GroupNotFound = 'UserManagement.NotFound.GroupNotFound',

	/**
	 * app名称已存在
	 */
	AppConflict = 'UserManagement.Conflict.AppConflict',

	/**
	 * 用户组名称已存在
	 */
	GroupConflict = 'UserManagement.Conflict.GroupConflict',

	/**
	 * 用户组部门不存在
	 */
	DepartmentNotExisted = 'UserManagement.BadRequest.DepartmentNotFound',

	/**
	 * 用户组不存在
	 */
	UserGroupNotFound = 'UserManagement.BadRequest.GroupNotFound',

	/**
	 * 应用账号不存在
	 */
	AppAccountNotFound = 'UserManagement.BadRequest.AppNotFound',
}

export const enum PolicyMgntErrorCode {
    /**
     * 策略不存在
     */
    ProductAuthorizedNotValid = 'PolicyMgnt.BadRequest.ProductAuthorizedNotValid',
}
