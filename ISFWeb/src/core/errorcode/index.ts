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
     * 部门不存在
     */
    DepartmentNotExist = 400019002,

    /**
     * 用户组不存在
     */
    UserGruopsNotExist = 400019003,

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
     * 应用账户删除文档库管理权限时，某个权限不存在
     */
    UseAccountDocLibPermNotExist = 404000000,

    /**
     * 新建/编辑文档库访问策略时，文档库不存在
     */
    AccessPolicyDocLibNotExist = 400034001,

    /**
     * 新建/编辑文档库访问策略时，访问者（用户）不存在
     */
    AccessPolicyUsersNotExist = 400034002,

    /**
     * 新建/编辑文档库访问策略时，访问者（部门）不存在
     */
    AccessPolicyDepsNotExist = 400034003,

    /**
      * 新建/编辑文档库访问策略时，访问者（用户组）不存在
      */
    AccessPolicyUserGroupsNotExist = 400034004,

    /**
     * 应用账户已存在（获取OSS网关账户）
     */
    AppUserExists = 403031022,

    /**
     * 应用账户已删除（删除OSS网关账户）
     */
    AppUserNoExists = 404031005,

    /**
     * 用户不存在（共享模板）
     */
    UserNoExist = 400019001,

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
     * 空间站实验项目名称或标识冲突
     */
    ExperimentNameConflict = 400056402,

    /**
     * 空间站载荷信息名称冲突
     */
    RackNameConflict = 400055010,

    /**
     * 空间站载荷信息标识冲突
     */
    RackCodeConflict = 400055011,

    /**
     * 载荷被实验项目引用
     */
    RackRefError = 400055014,
}
