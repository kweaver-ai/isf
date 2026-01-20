declare namespace Core {
    namespace APIs {
        namespace EACHTTP {

            // ---
            // 可复用的数据结构
            // ---

            /**
             * 访问者类型定义
             */
            type AccessorType = 'user' | 'department' | 'contact';

            /**
             * 部门路径
             */
            interface BasicDepInfo {
                /**
                 * 部门路径
                 */
                name: string;
                path: string
            }

            /**
            * 获取用户路径信息
            */
            interface BasicUserInfo {
                /**
                 * 用户路径信息
                 */
                directdepinfos: ReadonlyArray<string>,
                deppath: string
            }

            /**
             * 搜索组织架构获取搜索条目
             */
            interface DepartmentSearchCountResult {
                count: number
            }

            /**
             * 访问者
             */
            type Accessor = {
                /**
                 * 访问者id
                 */
                accessorid: string;

                /**
                 * 访问者名字
                 */
                accessorname: string;

                /**
                 * 访问者类型
                 */
                accessortype: AccessorType;
            }

            /**
             * 权限值
             */
            type PermValue = {
                /**
                 * 允许的权限
                 */
                allowvalue: number;

                /**
                 * 拒绝的权限
                 */
                denyvalue: number;
            }

            /**
             * 审核信息
             */
            type Audit = {
                /**
                 * 审核者名称
                 */
                auditorname: string;

                /**
                 * 审核结果
                 */
                auditresult: boolean;

                /**
                 * 审核意见
                 */
                auditmsg: string;
            }

            /**
             * 文档类型
             */
            type DocType = 'userdoc' | 'sharedoc' | 'customdoc';

            // ---
            //  服务端配置相关
            // ---

            /**
             * 第三方认证配置参数
             */
            interface ThirdAuthConfig {
                /**
                 * 隐藏登录界面
                 */
                hideLogin?: boolean;

                /**
                 * 隐藏第三方登录界面
                 */
                hideThirdLogin?: boolean;

                /**
                 * 其他第三方认证配置
                 */
                [key: string]: any;
            }

            /**
             * 第三方认证配置
             */
            interface ThirAuth {
                /**
                 * 唯一第三方认证系统
                 */
                id?: string;

                /**
                 * 第三方认证系统的配置参数
                 */
                config: Readonly<ThirdAuthConfig>;
            }

            /**
             * 第三方标密系统配置，如果未开启，则不会有
             */
            interface ThirdCSFSysConfig {
                /**
                 *  第三方标密系统唯一标识
                 */
                id: string;

                /**
                 *  仅上传已标密文件
                 */
                only_upload_classified: boolean;

                /**
                 *  仅共享已标密文件
                 */
                only_share_classified: boolean;
            }

            /**
             * 表示windows ad单点登录相关配置信息
             */
            interface WindowsADSSO {
                /**
                 * 表示是否开启了windows ad单点登录
                 */
                is_enabled: boolean;
            }

            /**
             * 外部应用
             */
            interface ExtAPP {
                /**
                 * 在线表格
                 */
                enable_chaojibiaoge: boolean;
            }

            /**
             * 配置信息
             */
            interface CacheConfig {
                /**
                 * 缓存清除时间间隔
                 */
                interval: number;

                /**
                 * 缓存清除空间上限
                 */
                size: number;
            }

            /**
             * 配置信息
             */
            interface CacheConfigs {
                /**
                 * 缓存相关配置信息
                 */
                cache: ReadonlyArray<CacheConfig>
            }

            /**
             * OEM 配置项
             */
            interface OEMConfig {
                /**
                 * 权限配置，是否允许给密级低于文件密级的用户配置权限
                 */
                allowauthlowcsfuser: boolean;

                /**
                 * 权限配置，是否允许配置所有者
                 */
                allowowner: boolean;

                /**
                 * 客户端退出时，是否强制清除缓存
                 */
                clearcache: boolean;

                /**
                 *  "客户端超时退出时间，单位为分钟
                 * -1表示不开启超时退出
                 * 其它值表示多少分钟后退出
                 * 例如30表示30分钟超时后退出
                 */
                clientlogouttime: number;

                /**
                 * 客户端配置权限时，默认权限有效天数
                 * -1为无限期
                 */
                defaultpermexpireddays: number;

                /**
                 * 是否开启文件传输限制功能
                 */
                enablefiletransferlimit: boolean;

                /**
                 * 是否启用onedrive跳转(默认为false)
                 */
                enableonedrive: boolean;

                /**
                 * 是否开启共享审核, true表示开启
                 * false表示关闭
                 */
                enableshareaudit: boolean;

                /**
                 * 是否显示用户协议（默认为false）
                 */
                enableuseragreement: boolean;

                /**
                 * 客户端是否隐藏缓存设置
                 */
                hidecachesetting: boolean;

                /**
                 * 权限配置，截至时间是否支持无限期
                 */
                indefiniteperm: boolean;

                /**
                 * 密码配置，密码允许的最大过期天数，
                 * -1为永久
                 */
                maxpassexpireddays: number;

                /**
                 * OWAS的预览url，如果未配置则为空字符串
                 */
                owasurl: string;

                /**
                 * 登录配置，是否允许记住用户名和密码
                 */
                rememberpass: boolean;
            }

            /**
             * 密级枚举信息
             */
            interface CSFLevelEnum {
                [key: string]: number
            }

            /**
             * 服务器配置信息
             */
            interface Config {
                /**
                 * 返回入口文档视图模式
                 * 1 代表旧视图， 
                 * 2 代表新视图（默认）
                 */
                entrydoc_view_config: number;

                /**
                 * 客户端是否使用https
                 */
                https: boolean;

                /**
                 * 是否使用文件提取码
                 */
                enable_link_access_code: boolean;

                /**
                 * 是否开启涉密模式
                 */
                enable_secret_mode: boolean;

                /**
                 * 是否允许共享邀请
                 */
                enable_invitation_share: boolean;

                /**
                 * 是否开启文件水印
                 */
                enable_doc_watermark: boolean;

                /**
                 * 是否开启外链配置水印
                 */
                enable_outlink_watermark: boolean;

                /**
                 * 是否开启文件评论
                 */
                enable_doc_comment: boolean;

                /**
                 * 自动锁提醒
                 */
                auto_lock_remind: boolean;

                /**
                 * 第三方修改密码地址
                 */
                third_pwd_modify_url: string;

                /**
                 * 最大标签数
                 */
                tag_max_num: number;

                /**
                 * 服务端版本
                 */
                server_version: string;

                /**
                 * 禁用某种设备类型登陆
                 */
                forbid_ostype: number;

                /**
                 * Oem配置信息
                 */
                oemconfig: Readonly<OEMConfig>;

                /**
                 * 表示开启了第三方认证，如果未开启，则不会有
                 */
                thirdauth?: Readonly<ThirAuth>;

                /**
                 * 表示windows ad单点登录相关配置信息
                 */
                windows_ad_sso: Readonly<WindowsADSSO>

                /**
                 * 密级枚举信息
                 */
                csf_level_enum: Readonly<CSFLevelEnum>;

                /**
                 * 第三方标密系统配置，如果未开启，则不会有
                 */
                third_csfsys_config?: Readonly<ThirdCSFSysConfig>;

                /**
                 * 外部应用
                 */
                extapp: Readonly<ExtAPP>;
            }

            // ---
            //  用户与认证相关
            // ---

            /**
             * 鉴权信息
             */
            interface AuthInfo {
                /**
                 * 登陆超时
                 */
                expires: number;

                /**
                 * 是否需要二次验证
                 */
                needmodifypassword: boolean;

                /**
                 * 与userid一起验证请求的合法性
                 */
                tokenid: string;

                /**
                 * 唯一标识用户的ID
                 */
                userid: string;
            }

            /**
             * Auth2中的鉴权信息
             */
            interface Auth2Info {
                /**
                 * account对应的用户id，用来后续请求时附带到query_string中
                 */
                user_id: string;

                /**
                 * 服务器颁发的token，用来后续请求时附带到query_string中
                 */
                token_id: string;

                /**
                 * token过期时间，单位为秒，根据登录时token_type而不同
                 */
                expires_in: number;

                /**
                 * refresh_token，有效期为60天，如果access_token过期，需要拿refresh_token换取新的token_id，refresh后旧的refresh_token会失效
                 */
                refresh_token: string;
            }

            /**
             * 第三方登录鉴权信息
             */
            interface ThirdAuthInfo {
                /**
                 * 唯一标识用户的ID
                 */
                userid: string;

                /**
                 * 与userid一起验证请求的合法性
                 */
                tokenid: string;

                /**
                 * token的有效期，单位为秒
                 */
                expires: number;
            }

            /**
             * 登录web（信任的第三方应用appid）
             */
            interface ExtLoginWebInfo {
                /**
                 * AnyShare的web链接，通过该地址可以访问文档
                 */
                weburl: string;
            }

            /**
             * 登录外部应用（集成到anyshare）
             */
            interface ExtAppInfo {
                /**
                 * 返回登录信息
                 */
                value: any;
            }

            /**
             * 刷新身份凭证有效期
             */
            interface RefreshTokenInfo {
                /**
                 * 刷新以后token的有效期，单位为秒
                 */
                expires: number;
            }

            /**
             * 二次认证
             */
            interface SecurityValidateResult {
                /**
                 * 验证结果
                 */
                result: boolean;
            }

            // ---
            // 用户管理相关
            // ---

            /**
             * 直属部门信息
             */
            interface Directdeptinfo {
                /**
                 * 部门id
                 */
                depid: string;

                /**
                 * 部门名称
                 */
                name: string;
            }

            /**
             * 用户信息
             */
            interface UserInfo {

                /**
                 * 用户id
                 */
                userid: string;

                /**
                 * 账户名
                 */
                account: string;

                /**
                 * 显示名
                 */
                name: string;

                /**
                 * 邮箱
                 */
                mail: string;

                /**
                 * 密级
                 */
                csflevel: number;

                /**
                 * 允许打印，拷屏
                 */
                leakproofvalue: number;

                /**
                 * 密码管控
                 */
                pwdcontrol: number;

                /**
                 * 用户类型
                 */
                usertype: number;

                /**
                 * 直属部门信息
                 */
                directdeptinfos: Array<Directdeptinfo>;

                /**
                 * 角色类型
                 */
                roletypes: Array<string>;

                /**
                 * 二次验证
                 */
                needsecondauth: boolean;

                /**
                 * 冻结状态
                 * true：冻结 false：未冻结
                 */
                freezestatus: boolean;

                /**
                 * 同意使用协议状态
                 * true：同意 false：未同意
                 */
                agreedtotermsofuse: boolean;
            }

            /**
             * 同意用户使用协议
             */
            interface AgreedResult {
                /**
                 * 执行结果
                 */
                result: string;
            }

            // ---
            // 配额管理相关
            // ---

            /**
             * 单个文档配额使用情况
             */
            interface QuotaInfo {
                /**
                 * 文档类型[“userdoc”]
                 */
                doctype: string;

                /**
                 * 类型名字
                 */
                typename: string;

                /**
                 * 文档id
                 */
                docid: string;

                /**
                 * 文档名
                 */
                docname: string;

                /**
                 * 配额
                 */
                quota: number;

                /**
                 * 已用空间
                 */
                used: number;
            }

            /**
             * 用户配额信息
             */
            interface QuotaInfos {
                /**
                 * 所有文档的配额信息
                 */
                quotainfos: ReadonlyArray<QuotaInfo>;
            }

            /**
             * 根据cid路径获取的配额信息
             */
            interface QuotaInfosByCid {
                /**
                 * 配额
                 */
                quota: number;

                /**
                 * 已用空间
                 */
                used: number;
            }

            /**
            * OAuth信息
            */
            interface OAuthInfo {
                /**
                 * 完整的OAuth认证地址
                 */
                authurl: string;

                /**
                * 是否开启
                */
                isenabled: boolean;

                /**
                * 认证服务器的地址
                */
                authserver: string;

                /**
                * 第三方认证完毕后，跳转到AnyShare服务器的地址
                */
                redirectserver: string;
            }

            // ---
            // 文件锁相关
            // ---

            /**
             * 尝试文件锁结果
             */
            interface TryLock {
                /**
                 * 锁定成功返回true，被其他用户锁定返回false
                 */
                issucceed: boolean;

                /**
                 * 锁定者id
                 */
                lockerid?: string;

                /**
                 * 锁定者用户名
                 */
                lockeraccount?: string;

                /**
                 * 锁定者显示名
                 */
                lockername?: string;
            }

            /**
             * 锁信息
             */
            interface RefreshLock {
                /**
                 * 文件id
                 */
                docid: string;

                /**
                 * 0表示刷新成功，1表示失败
                 */
                state: number;

                /**
                 * state为1时会返回异常
                 */
                errmsg?: string;
            }

            /**
             * 刷新锁
             */
            interface Refresh {
                lockinfos: Array<LockInfo>
            }

            /**
             * 锁信息
             */
            interface LockInfo {
                /**
                 * 是否已被锁定
                 */
                islocked: boolean;

                /**
                 * 锁定者id
                 */
                lockerid?: string;

                /**
                 * 锁定者用户名
                 */
                lockeraccount?: string;

                /**
                 * 锁定者显示名
                 */
                lockername?: string;
            }

            /**
             * 当前用户所有锁定的文件信息
             */
            interface LockedFileInfos {
                /**
                 * 文件信息列表
                 */
                docinfos: Array<LockFileBasicInfo>
            }

            /**
             * 锁定文件基本信息
             */
            interface LockFileBasicInfo {
                /**
                 * 文件id
                 */
                docid: string;

                /**
                 * 文件路径
                 */
                path: string;

                /**
                 * 文件锁创建时间(微秒时间戳)
                 */
                locktime: number;

                /**
                 * 锁定者id
                 */
                lockerid: string;

                /**
                 * 锁定者用户名
                 */
                lockeraccount: string;

                /**
                 * 锁定者显示名
                 */
                lockername: string;
            }

            // ---
            // 审核相关
            // ---

            /**
             * 共享审核消息
             */
            type ShareApproval = {
                /**
                 * 1.新增
                 * 2.编辑
                 * 3.删除
                 */
                optype: 1 | 2 | 3;

                /**
                 * 截止时间
                 */
                endtime: number;
            } & Accessor & PermValue

            /**
             * 外链审核消息
             */
            interface LinkApproval {
                /**
                 * 1.新增
                 * 2.编辑
                 */
                optype: 1 | 2;

                /**
                 * 允许权限
                 */
                perm: number;

                /**
                 * 外链密码
                 */
                password: string;

                /**
                 * 截止时间
                 */
                endtime: number;
            }

            /**
             * 所有者审核消息
             */
            type OwnerApproval = {
                /**
                 * 1.新增
                 * 3.删除
                 */
                optype: 1 | 3;

                /**
                 * 截止时间
                 */
                endtime: number;
            } & Accessor

            /**
             * 修改密级审核消息
             */
            interface CSFApproval {
                /**
                 * 要申请的密级大小
                 */
                applycsflevel: number;
            }

            /**
             * 更改密级申请具体信息（完整）
             */
            type CSFApprovalComplete = {
                /**
                 * 申请时用户密级
                 */
                usercsflevel: number;
            } & CSFApproval & Accessor

            /**
             * 申请基本信息
             */
            interface ApplyBaseInfo {
                /**
                 * 申请记录id
                 */
                applyid: string;

                /**
                 * 文档gns路径
                 */
                docid: string;

                /**
                 * 文档名称
                 */
                docname: string;

                /**
                 * 申请时间（微秒）
                 */
                createdate: number;

                /**
                 * 是否是文件夹
                 */
                isdir: boolean;

                /**
                 * 文件密级,5~15，如果是文件夹，则为0
                 */
                csflevel: number;
            }

            /**
             * 审核基本信息
             */
            interface AuditBaseInfo {
                /**
                * 流程名称
                */
                processname: string;

                /**
                * 审核模式：
                * 1. 同级审核
                * 2. 汇签审核
                * 3. 逐级审核
                */
                audittype: number;

                /**
                * 所有审核员名称列表 string
                */
                auditornames: ReadonlyArray<string>
            }

            /**
             * 返回的数量
             */
            interface AuditCount {
                /**
                * 返回数
                */
                count: number;
            }

            /**
             * 我的共享申请
             */
            interface MyApply extends ApplyBaseInfo {
                /**
                 * 申请类型
                 * 1. 表示共享申请
                 * 2. 表示外链申请
                 * 3. 表示所有者申请
                 * 5. 表示修改密级申请
                 */
                apptype: 1 | 2 | 3 | 5;

                /**
                 * 所有的审核员名称
                 */
                auditornames: ReadonlyArray<string>;

                /**
                 * 具体的申请信息，根据apptype，所包含的字段不同
                 */
                detail: ShareApproval | LinkApproval | OwnerApproval | CSFApprovalComplete;
            }

            /**
             * 我的共享申请
             */
            interface MyApplys {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<MyApply>;
            }

            /**
             * 待审核的消息
             */
            interface ApplyApproval extends ApplyBaseInfo {
                /**
                 * 申请者名称
                 */
                sharename: string;

                /**
                 * 申请类型
                 * 1. 表示共享申请
                 * 2. 表示外链申请
                 * 3. 表示所有者申请
                 * 5. 表示修改密级申请
                 */
                apptype: 1 | 2 | 3 | 5;

                /**
                 * 申请详细信息
                 */
                detail: ShareApproval | LinkApproval | OwnerApproval | CSFApproval;
            }

            /**
             * 待审核的消息
             */
            interface ApplyApprovals {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<ApplyApproval>;
            }

            /**
            * 所有流程信息
            */
            interface DocProcess extends AuditBaseInfo {
                /**
                 * 流程id
                 */
                processid: string;

                /**
                 * 发布目的路径: AnyShare://文档库1/dir
                 */
                destpath: string;
            }

            /**
             * 所有流程信息
             */
            interface DocProcesses {
                /**
                 * 流程信息
                 */
                processinfos: ReadonlyArray<DocProcess>;
            }

            /**
             * 流程审核信息
             */
            interface DocApproval extends ApplyBaseInfo, AuditBaseInfo {
                /**
                 * 申请者名称
                 */
                creatorname: string;

                /**
                 * 申请理由
                 */
                applymsg: string;

                /**
                 * 当前审核员在auditornames中的位置
                 * 同级审核该字段值为0
                 */
                auditprogress: number;
            }

            /**
             * 待审核的流程消息
             */
            interface DocApprovals {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<DocApproval>;
            }

            /**
             * 发起流程申请
             */
            interface PublishDocs {
                /**
                 * 流程审核模式
                 */
                audittype: number;

                /**
                 * 结果[“ok”]
                 */
                result: string;
            }

            /**
            * 申请中的流程信息
            */
            interface DocApply extends ApplyBaseInfo, AuditBaseInfo {
                /**
                 * 发起流程时的理由
                 */
                applymsg: string;

                /**
                 * 当前审核员在auditornames中的位置
                 * 同级审核该字段值为0
                 */
                auditprogress: number;
            }

            /**
             * 申请中的流程信息
             */
            interface DocApplys {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<DocApply>;
            }

            /**
            * 流程申请历史
            */
            interface ApplyHistory extends ApplyBaseInfo, AuditBaseInfo {
                /**
                 * 发起流程时的理由
                 */
                applymsg: string;

                /**
                 * 结合auditorname,小于等于该值的表示审批通过的审核员，从0开始
                 */
                approveindex: number;

                /**
                 * 结合auditorname，该位置表示否决的人在auditornames中的位置
                 * 从0开始，-1表示没有人否决
                 */
                vetoindex: number;
            }

            /**
             * 流程申请历史
             */
            interface ApplyHistorys {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<ApplyHistory>;
            }

            /**
            * 流程审核历史
            */
            interface ApproveHistory extends ApplyBaseInfo, AuditBaseInfo {
                /**
                 *申请
                 */
                creatorname: string;

                /**
                 * 结合auditorname,小于等于该值的表示审批通过的审核员，从0开始
                 */
                approveindex: number;

                /**
                 * 结合auditorname，该位置表示否决的人在auditornames中的位置
                 * 从0开始，-1表示没有人否决
                 */
                vetoindex: number;
            }

            /**
             * 流程申请历史
             */
            interface ApproveHistorys {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<ApproveHistory>;
            }

            /**
            * 共享申请历史
            */
            interface ShareApplyHistory {
                /**
                 * 结合auditorname,等于该值的表示审批通过的审核员，从0开始
                 */
                approveindex: number;

                /**
                 * 结合auditorname，该位置表示否决的人在auditornames中的位置
                 * 从0开始，-1表示没有人否决
                 */
                vetoindex: number;

                /**
                * 审核员名称
                 */
                auditornames: ReadonlyArray<string>;

                /**
                 * 文档名称
                 */
                docname: string;

                /**
                 * 申请时间
                 */
                requestdate: number;

                /**
                 * 申请类型：
                 * 1. 共享申请
                 * 2. 外链申请
                 * 3. 所有者申请
                 * 5. 更改密级申请
                 */
                apptype: number;

                /**
                 * 是否是文件夹
                 */
                isdir: boolean;

                /**
                 * 文件密级,5~15，如果是文件夹，则为0
                 */
                csflevel: number;

                /**
                 * 具体的申请信息，根据apptype，所包含的字段不同
                 */
                detail: ShareApproval | LinkApproval | OwnerApproval | CSFApprovalComplete;
            }

            /**
             * 流程申请历史
             */
            interface ShareApplyHistorys {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<ShareApplyHistory>;
            }

            /**
            * 共享审核历史    
            */
            interface ShareApproveHistory extends ShareApplyHistory {
                /**
                * 共享者名称
                */
                requestername: string;
            }

            /**
             * 共享审核历史
             */
            interface ShareApproveHistorys {
                /**
                 * 申请信息
                 */
                applyinfos: ReadonlyArray<ShareApproveHistory>;
            }

            /**
             * 待审核的条目
             */
            interface Pendingapprovalscounts {
                /**
                 * 待审核总数
                 */
                counts: number;

                /**
                 * 共享待审核数目(包含定密审核的数量)
                 */
                pemcount: number;

                /**
                 * 流程待审核数目
                 */
                docauditcount: number;

                /**
                 * 待审核的定密审核数量
                 */
                csfauditcount: number;

                /**
                 * 待审核的共享审核数量（不包含定密审核的数量）
                 */
                shareauditcount: number;
            }

            /**
             * 内外网数据交换文档接收位置
             */
            interface Docrecvarea {
                /**
                 * 接收区名称
                 */
                recvareaname: string;

                /**
                 * 接收区id
                 */
                recvareaid: string;
            }

            /**
             * 内外网交换模糊匹配目标文档库名
             */
            interface DstDocname {
                /**
                 * 文档库名
                 */
                docname: string;

                /**
                 * 用户登录名，仅个人文档有此参数
                 */
                loginname: string;
            }

            // ---
            // OEM配置
            // ---

            /**
             * oem语言段
             */
            interface OEMSectionInfo {
                /**
                 * OEM图片资源
                 */
                'background.png': string;

                /**
                 * OEM图片资源
                 */
                'favicon.ico': string;

                /**
                 * OEM图片资源
                 */
                'logo.png': string;

                /**
                 * OEM图片资源
                 */
                'org.png': string;

                /**
                 * 帮助文档链接地址
                 */
                helper: string;

                /**
                 * 组织机构
                 */
                organization: string;

                /**
                 * 产品名称
                 */
                product: string;

                /**
                 * 官方链接
                 */
                site: string;

                /**
                 * 版权信息
                 */
                copyright: string;
            }

            /**
             * AnyShare段字段
             */
            interface OEMGlobalInfo {
                /**
                 * 是否开启用户许可协议
                 */
                userAgreement: string;

                /**
                 * 用户许可协议
                 */
                agreementText: string;

                /**
                 * 允许简体中文
                 */
                allowCn: string;

                /**
                 * 允许繁体中文
                 */
                allowTw: string;

                /**
                * 允许英文语言
                */
                allowEn: string;

                /**
                 * 开放mac客户端下载
                 */
                mac: string;

                /**
                 * 开放android客户端下载
                 */
                android: string;

                /**
                 * 开放iOS客户端下载
                 */
                ios: string;

                /**
                 * 开放Windows客户端下载
                 */
                windows: string;

                /**
                 * 开放Office插件下载
                 */
                office: string;

                /**
                 * 基本颜色值
                 */
                theme: string;

                /**
                * Web client功能视图
                */
                webclientTabs: string;

                /**
                 * Web console功能视图
                 */
                webconsoleTabs: string;

                /**
                 * 显示产品信息
                 */
                showProduct: string;

                /**
                 * 显示型号信息
                 */
                showHardware: string;

                /**
                 * 显示授权信息
                 */
                showLicense: string;

                /**
                 * 显示版权信息
                 */
                showCopyright: string;
            }

            /**
             *  Oem信息
             */
            interface OEMInfo extends OEMSectionInfo, OEMGlobalInfo {
            }

            /**
             * 水印配置
             */
            interface DocWatermarkConfig {
                /**
                 * 用户水印
                 */
                user: Readonly<{
                    /**
                     * 启用
                     */
                    enabled: boolean;

                    /**
                     * 用户名名称颜色
                     */
                    color: string;

                    /**
                     * 用户名字体大小
                     */
                    fontSize: number;

                    /**
                     * 用户名水印版式
                     */
                    layout: 0 | 1;

                    /**
                     * 用户名透明度
                     */
                    transparency: number;
                }>;

                /**
                 * 文本水印
                 */
                text: Readonly<{
                    /**
                     * 启用文字水印
                     */
                    enabled: boolean;

                    /**
                     * 文字水印颜色
                     */
                    color: string;

                    /**
                     * 文字水印内容
                     */
                    content: string;

                    /**
                     * 文字水印字体大小
                     */
                    fontSize: number;

                    /**
                     * 文字水印版式
                     */
                    layout: 0 | 1;

                    /**
                     * 文字水印透明度
                     */
                    transparency: number;
                }>;

                /**
                 * 图片水印
                 */
                image: Readonly<{
                    /**
                     * 启用图片水印
                     */
                    enabled: boolean;

                    /**
                     * 图片水印base64编码
                     */
                    src: string;

                    /**
                     * 图片水印版式
                     */
                    layout: 0 | 1;

                    /**
                     * 图片水印透明度
                     */
                    transparency: number;
                }>;
            }

            interface SiteOfficeOnlineInfo {
                /**
                 * 本机IP
                 */
                ip: string;

                /**
                 * 在线预览地址
                 */
                office: string;

                /**
                 * 在线编辑地址
                 */
                wopi: string;
            }

            // ---
            // 客户端下载
            // ---

            /**
            * 客户端是否需要检查更新
            */
            interface UpdateInfo {
                /**
                 * 客户端是否需要更新
                 * “ok”表示客户端版本是最新的
                 * “optional”表示客户端可选择更新
                 * “forced”表示客户端必须进行更新
                 */
                result: string;

                /**
                * 可用客户端安装包的版本信息
                */
                latest: string;

                /**
                * 可用客户端安装包的相对下载路径，需要加上服务器的ip后才能进行下载
                */
                downloadurl: string;
            }

            /**
             * 客户端下载地址
             */
            interface ClientURL {
                /**
                 * 下载地址
                 */
                URL: string;
            }

            // ---
            // 服务器地址
            // ---

            /**
             * 服务器地址配置
             */
            interface HostInfo {
                /**
                 * 服务器域名或IP
                 */
                host: string;

                /**
                 * 服务器站点名称
                 */
                name: string;

                /**
                 * 服务器HTTP端口
                 */
                port: number;

                /**
                 * 服务器HTTPS端口
                 */
                https_port: number;
            }

            // ---
            // 权限相关
            // ---

            /**
             * 权限检查结果
             */
            interface PermCheckResult {
                /**
                 * 0. 表示OK 
                 * 1. 表示未配置该权限
                 * 2. 表示拒绝该权限
                 */
                result: number;
            }

            /**
             * 获取权限配置信息返回信息
             */
            type PermInfo = {

                /**
                 * 权限id
                 */
                id: string;

                /**
                 * 是否是允许权限
                 */
                isallowed: boolean;

                /**
                 * 权限值
                 */
                permvalue: number;

                /**
                 * 继承路径
                 */
                inheritpath: string;

                /**
                 * 有效期
                 */
                endtime: number;
            } & Accessor

            /**
             * 获取权限配置信息返回信息
             */
            interface PermInfos {
                perminfos: ReadonlyArray<PermInfo>;
            }

            /**
             * 检查单个权限
             */
            interface checkInfo {
                /**
                 * 需要检查文档的docid
                 */
                docid: string;

                /**
                 * 权限值
                 */
                perm: 1 | 2 | 4 | 8 | 16 | 32 | 64;

                /**
                 * userid
                 */
                userid?: string;
            }

            /**
             * 批量设置权限返回结果
             */
            interface PermSetResult {
                /**
                 * 设置结果
                 * 0 表示成功
                 * 1 表示提交申请
                 */
                result: number;
            }

            /**
             * 访问者的最终权限
             */
            type PermResult = Accessor & PermValue;

            /**
             * 所有者的最终权限
             */
            type OwnerResult = Accessor;

            /**
             * 访问者的最终权限
             */
            interface PermResults {
                permresults: ReadonlyArray<PermResult>;
                ownerresults: ReadonlyArray<OwnerResult>;
            }

            /**
             * 检查所有权限
             */
            interface CheckAllResults {
                /**
                 * 允许的权限值
                 */
                allowvalue: number;

                /**
                 * 拒绝的权限值
                 */
                denyvalue: number;
            }

            /**
             * 用户已共享的文档
             */
            interface SharedDoc {
                /**
                 * 文档路径
                 */
                path: string;

                /**
                 * 文档的id
                 */
                docid: string;

                /**
                 * 文件大小，文件夹则为-1
                 */
                size: number;

                /**
                 * 如果是文件，返回由客户端设置的文件本地修改时间；若未设置，返回modified的值
                 */
                client_mtim: number;
            }

            /**
             * 用户已共享的文档
             */
            interface SharedDocs {
                docinfos: ReadonlyArray<SharedDoc>
            }

            /**
             * 内链共享模板返回结果
             */
            interface InternalLinkTemplate {
                /**
                 * 允许设置的权限
                 */
                allowperm: number;

                /**
                 * 默认设置的权限
                 */
                defaultperm: number;

                /**
                 * 允许设定所有者
                 */
                allowowner: boolean;

                /**
                 * 默认允许设置所有者
                 */
                defaultowner: boolean;

                /**
                 * 是否限制有效期
                 */
                limitexpiredays: boolean;

                /**
                 * 默认有效期
                 */
                allowexpiredays: number;
            }

            /**
             * 外链共享模板返回结果
             */
            interface ExternalLinkTemplate {

                /**
                 * 可设定的访问权限
                 */
                allowperm: number;

                /**
                 * 默认访问权限
                 */
                defaultperm: number;

                /**
                 * 是否限制外链有效期
                 */
                limitexpiredays: boolean;

                /**
                 * 有效期
                 */
                allowexpiredays: number;

                /**
                 * 是否强制使用访问密码
                 */
                accesspassword: boolean;

                /**
                 * 是否限制外链打开次数
                 */
                limitaccesstimes: boolean;

                /**
                 * 打开次数
                 */
                allowaccesstimes: number;
            }

            /**
             * 共享文档开关配置
             */
            interface ShareDocConfigResult {
                /**
                 * 开启个人文档内链共享
                 */
                enable_user_doc_inner_link_share: boolean;

                /**
                 * 开启个人文档外链共享
                 */
                enable_user_doc_out_link_share: boolean;

            }

            // ---
            // 所有者管理相关
            // ---

            /**
             * 是否是所有者
             */
            interface OwnerCheckResult {
                /**
                 * 是否是所有者
                 */
                isowner: boolean;
            }

            /**
             * 所有者信息
             */
            interface OwnerInfo {

                /**
                 * userid
                 */
                userid: string;

                /**
                 * 账号
                 */
                account: string;

                /**
                 * 显示名
                 */
                name: string;

                /**
                 * 继承路径
                 */
                inheritpath: string
            }

            /**
             * 所有者信息列表
             */
            interface OwnerInfos {
                ownerinfos: ReadonlyArray<OwnerInfo>
            }

            /**
             * 批量设置权限返回结果
             */
            interface OwnerSetResult {
                /**
                 * 设置结果
                 * 0 表示成功或没有变化
                 * 1 表示提交申请
                 */
                result: number;
            }

            // ---
            // 消息
            // ---

            /**
             * 消息结构
             */
            interface Message {
                /**
                 * 消息id
                 */
                id: string;

                /**
                 * 消息类型
                 */
                type: number;

                /**
                 * 发送者名字
                 */
                sender: string;

                /**
                 * 接受者名字，外链则为空字符串
                 */
                accessorname: string;

                /**
                 * 访问者类型
                 */
                accessortype: AccessorType;

                /**
                 * 操作时间
                 */
                time: number;

                /**
                 * 内链，相对地址
                 */
                url: string;

                /**
                 * 是否是文件夹
                 */
                isdir: boolean;

                /**
                 * 是否已读
                 */
                isread: boolean;

                /**
                 * 密级信息，文件夹为 0
                 */
                csf: number;
            }

            /**
             * 共享申请消息
             * 消息类型为 1,2,5,6,10,11
             */
            interface ShareApplyMessage extends Message {
                /**
                 * 有效期
                 */
                end?: number;

                /**
                 * 允许权限
                 */
                allowvalue?: number;

                /**
                 * 拒绝权限
                 */
                denyvalue?: number;
            }

            /**
             * 外链申请消息
             * 消息类型为 9,14
             */
            interface LinkApplyMessage extends Message {
                /**
                 * 有效期
                 */
                end?: number;

                /**
                 * 外链权限
                 */
                perm?: number;
            }

            /**
             * 共享审核消息
             * 消息类型为 10,11,12,13,14
             */
            interface ShareApprovalMessage extends Message, Audit {
            }

            /**
             * 流程申请消息
             * 消息类型 15
             */
            interface WorkflowApplyMessage extends Message {
                /**
                 * 申请者名字
                 */
                requester?: string;

                /**
                 * 上一级审核员名称
                 */
                previousauditor?: string;

                /**
                 * 流程名称
                 */
                processname?: string;

                /**
                 * 审核类型
                 */
                docaudittype?: number;
            }

            /**
             * 流程进度消息
             * 消息类型 16
             */
            type WorkflowProgressMessage = Message & Audit & {

                /**
                 * 流程名称
                 */
                processname?: string;

                /**
                 * 审核类型
                 */
                docaudittype?: number;

                /**
                 * 下一级审核者名字
                 */
                nextauditor?: string;

                /**
                 * 下一级审核员级别
                 */
                nextlevel?: number;

            }

            /**
             * 流程结果消息
             * 消息类型为 17
             */
            interface WorkflowResultMessage extends Message {

                /**
                 * 流程名称
                 */
                processname?: string;

                /**
                 * 审核类型
                 */
                docaudittype?: number;

                /**
                 * 审核者名字
                 */
                auditorname?: string;

                /**
                 * 审核结果
                 */
                auditresult?: boolean;

                /**
                 * 审核意见
                 */
                auditmsg?: string;
            }

            /**
             * 普通消息
             * 消息类型为 18
             */
            interface SimpleMessage extends Message {

                /**
                 * 消息内容
                 */
                content: string;
            }

            /**
             * 密级申请消息
             * 消息类型为 19
             */
            interface CSFApplyMessage extends Message {

                /**
                 * 申请密级
                 */
                applycsflevel: number;
            }

            /**
             * 密级结果消息
             * 消息类型为 20
             */
            interface CSFResultMessage extends Message {
                /**
                 * 审核者名字
                 */
                auditorname?: string;

                /**
                 * 审核结果
                 */
                auditresult?: boolean;

                /**
                 * 审核意见
                 */
                auditmsg?: string;
            }

            /**
             * 隔离区消息
             * 消息类型为 21,22
             */
            interface QuarantineMessage extends Message {
                /**
                 * 创建者
                 */
                creator: string;

                /**
                 * 修改者
                 */
                modifier: string;
            }

            /**
             * 杀毒推送消息
             * 消息类型为 23
             */
            interface AntiVirusMessage extends Message {
                /**
                 * 杀毒管理员
                 */
                antivirusadmin?: string;

                /**
                 * 杀毒类型 
                 * 1 表示隔离
                 * 2 表示修复
                 * 3 表示还原
                 */
                antivirusop?: number;
            }

            /**
             * 消息列表
             */
            interface Messages {
                /**
                 * 最后接收通知的事件，可供下次调用get使用
                 */
                stamp: number;

                /**
                 * 消息列表
                 */
                msgs: ReadonlyArray<ShareApplyMessage | LinkApplyMessage | ShareApprovalMessage | WorkflowApplyMessage | WorkflowProgressMessage | WorkflowResultMessage | SimpleMessage | CSFApplyMessage | CSFResultMessage | QuarantineMessage | AntiVirusMessage>
            }

            // ---
            // 联系人
            // ---

            /**
             * 联系人组
             */
            interface ContactGroup {
                /**
                 * 分组id
                 */
                id: string;

                /**
                 * 创建人id
                 */
                createrid: string;

                /**
                 * 分组名称
                 */
                groupname: string;

                /**
                 * 联系人个数
                 */
                count: number;

            }

            /**
             * 联系人组
             */
            interface ContactGroups {
                /**
                 * 联系人组
                 */
                groups: ReadonlyArray<ContactGroup>;
            }

            /**
             * 联系人
             */
            interface ContactUser {
                /**
                 * 联系人id
                 */
                userid: string;

                /**
                 * 联系人帐号
                 */
                account: string;

                /**
                 * 联系人名称
                 */
                name: string;

                /**
                 * 邮箱
                 */
                mail: string;
            }

            /**
             * 搜索出的联系人
             */
            interface SearchedContactUser extends ContactUser {
                /**
                 * 所属联系人组id
                 */
                groupid: string;

                /**
                 * 所属联系人组名字
                 */
                groupname: string;
            }

            /**
             * 搜索出的联系人
             */
            interface SearchedContactUsers {
                userinfos: Array<SearchedContactUser>;
            }

            /**
             * 搜索联系人组获取搜索数目
             */
            interface ContactSearchCountResult {
                count: number;
            }

            /**
             * 联系人
             */
            interface ContactUsers {
                /**
                 * 联系人
                 */
                userinfos: Array<ContactUser>;
            }

            /**
             * 搜索联系人结果
             */
            interface ContactSearchResult {
                userinfos: ReadonlyArray<SearchedContactUsers>;
                groups: ReadonlyArray<ContactGroups>;
            }

            /**
             * 联系人
             */
            interface ContactUser2 {
                /**
                 * 联系人显示名
                 */
                username: string;

                /**
                 * 邮箱
                 */
                email: string;

                /**
                 * 部门
                 */
                departname: string;

                /**
                 * 联系人id
                 */
                userid: string;
            }

            /**
             * 联系人
             */
            interface ContactUsers2 {
                userinfos: Array<ContactUsers2>;
            }

            /**
             * 搜索出的联系人
             */
            interface SearchedContactUser2 {
                /**
                * 联系人id
                */
                userid: string;

                /**
                 * 联系人帐号
                 */
                account: string;

                /**
                 * 联系人名称
                 */
                name: string;

                /**
                 * 群组id
                 */
                groupid: string;

                /**
                 * 群组名
                 */
                groupname: string;
            }

            /**
             * 搜索出的联系人
             */
            interface SearchedContactUsers2 {
                userinfos: Array<SearchedContactUser2>;
            }

            // ---
            // 部门管理相关
            // ---

            /**
             * 部门信息
             */
            interface Department {

                /**
                 * 部门id
                 */
                depid: string;

                /**
                 * 部门名
                 */
                name: string;

                /**
                 * 是否可配置权限
                 */
                isconfigable: boolean;
            }

            /**
             * 部门
             */
            interface Departments {
                /**
                 * 部门
                 */
                depinfos: ReadonlyArray<Department>;
            }

            /**
             * 部门下的用户
             */
            interface DepartmentUser {
                /**
                 * 用户id
                 */
                userid: string;

                /**
                 * 用户帐号
                 */
                account: string;

                /**
                 * 用户显示名
                 */
                name: string;

                /**
                 * 用户邮箱
                 */
                mail: string;

            }

            /**
             * 用户信息
             */
            interface DepartmentUsers {
                /**
                 * 用户
                 */
                userinfos: ReadonlyArray<DepartmentUser>;
            }

            /**
             * 部门下搜素出的用户
             */
            interface SearchedDepartmentUser extends DepartmentUser {
                /**
                 * 所属部门id
                 */
                depid: string;

                /**
                 * 所属部门名
                 */
                depname: string;
            }

            /**
             * 在组织下搜索到的用户和部门信息
             */
            interface DepartmentSearchResult {
                userinfos: ReadonlyArray<SearchedDepartmentUser>,
                depinfos: ReadonlyArray<Department>
            }

            // ---
            // 共享邀请
            // ---

            /**
             * 共享邀请配置
             */
            interface InvitationBaseInfo {

                /**
                 * 邀请链接id
                 */
                invitationid: string;

                /**
                 * 邀请链接的有效期，-1为无限期
                 */
                invitationendtime: number;

                /**
                 * 权限
                 */
                perm: number;

                /**
                 * 权限有效期
                 */
                permendtime: number;
            }

            /**
             * 共享邀请备注信息
             */
            interface InvitationNote {

                /**
                 * 邀请图标
                 */
                image: string;

                /**
                 * 描述
                 */
                description: string;
            }

            /**
             * 共享邀请详细信息 
             */
            interface InvitationInfo {

                /**
                 * 邀请链接的有效期，-1为无限期
                 */
                invitationendtime: number;

                /**
                 * 权限
                 */
                perm: number;

                /**
                 * 权限有效期
                 */
                permendtime: number;

                /**
                 * 邀请图标
                 */
                image: string;

                /**
                 * 描述
                 */
                description: string;

                /**
                 * 文件名
                 */
                docname: string;

                /**
                 * 是否是文件夹
                 */
                isdir: boolean;
            }

            /**
             * 共享邀请加入后返回
             */
            interface InvitationJoinResult {
                /**
                 * 文档名称
                 */
                docname: string;

                /**
                 * 本次加入是否成功
                 */
                result: boolean;
            }

            // ---
            // 入口文档
            // ---

            /**
             * 入口文档（根据文档类型）
             */
            interface EntryDocByType {
                /**
                * 文档id
                */
                docid: string;

                /**
                 * 文档类型
                 */
                doctype: DocType;

                /**
                 * 文档类型名
                 * 文档库显示为文档库类型
                 */
                typename: string;

                /**
                 * 文档名
                 */
                docname: string;

                /**
                 * 版本
                 */
                otag: string;

                /**
                 * 文档大小，-1 表示文件夹
                 */
                size: number;

                /**
                 * 最后修改时间
                 */
                modified: number;

                /**
                 * 客户端修改时间
                 */
                client_mtime?: number;

                /**
                 * 权限
                 */
                attr?: number;
            }

            /**
             * 入口文档
             */
            interface EntryDoc extends EntryDocByType {
                /**
                 * 分类视图类型
                 */
                view_type?: number;

                /**
                 * 分类视图名
                 */
                view_name?: string;

                /**
                 * 视图下的文档类型
                 */
                view_doctype?: number;

                /**
                 * 视图下的文档类型名
                 */
                view_doctypename?: string;
            }

            /**
             * 入口文档列表
             */
            interface EntryDocs {
                docinfos: ReadonlyArray<EntryDoc>;
            }

            /**
             * 入口文档（根据文档类型）
             */
            interface EntryDocsByType {
                docinfos: ReadonlyArray<EntryDocByType>
            }

            /**
             * 视图类型
             */
            interface View {
                /**
                 * 视图类型
                 * 1. 用户文档，包括共享文档
                 * 3. 文档库
                 * 4. 归档库
                 */
                view_type: number;

                /**
                 * 视图名称
                 */
                view_name: string;

                /**
                 * 视图doctype
                 */
                view_doctype: number;

                /**
                 * doctype名称
                 */
                doc_type: string;
            }

            /**
             * 视图列表
             */
            interface Views {
                viewsinfo: ReadonlyArray<View>;
            }

            /**
             * 已退出文档信息
             */
            interface QuitInfo {
                /**
                 * 文件或文件夹的id
                 */
                docid: string;

                /**
                * 入口文档的类型
                */
                doctype: DocType;

                /**
                * 文档类型的显示名称
                */
                typename: string;

                /**
                * 文档的显示名称
                */
                docname: string;
            }

            /**
             * 已退出文档列表
             */
            interface QuitInfos {
                docinfos: ReadonlyArray<QuitInfo>
            }

            // ---
            // 群组文档管理相关
            // ---

            /**
             * 新建群组文档
             */
            interface GroupCreateResult {
                /**
                 * 文档id
                 */
                docid: string;

                /**
                 * 文档类型名
                 */
                typename: string;
            }

            // ---
            // 管理文档相关
            // ---

            /**
             * 文档信息
             */
            interface DocInfo {
                /**
                 * 文档id
                 */
                docid: string;

                /**
                 * 文档类型
                 */
                doctype: Core.APIs.EACHTTP.DocType;

                /**
                 * 文档类型名
                 * 文档库显示为文档库类型
                 */
                typename: string;

                /**
                 * 文档名
                 */
                docname: string;

                /**
                 * 版本
                 */
                otag: string;

                /**
                 * 文档大小，-1 表示文件夹
                 */
                size: number
            }

            /**
             * 获取文档列表（个人文档）
             */
            interface ManageDocs {
                docinfos: ReadonlyArray<DocInfo>;
            }

            /**
             * 文件夹锁信息
             */
            interface DirLockInfo {
                /**
                 * 是否已被锁定
                 */
                islocked: boolean;
            }

            // ---
            // CA认证相关
            // ---

            /**
             * CA信息（cainfo详情）
             */
            interface CAInfo {
                /**
                 * 厂商ID
                 */
                vender?: string;

                /**
                * 厂商描述
                */
                description?: string;

                /**
                * CA服务器信息
                */
                server?: string;

                /**
                * CA服务器分配的appid
                */
                appid?: string;

                /**
                * CA服务器分配的appkey
                */
                appkey?: string;

            }

            /**
             * CA信息
             */
            interface CAInfos {
                /**
                 * 是否开启了CA认证，true表示开启，false表示关闭
                 */
                enable: boolean;

                /**
                 * 描述CA信息，如果服务器未设置过CA信息，该字段不会出现
                 */
                cainfo?: Readonly<CAInfo>
            }

            // ---
            // PKI认证
            // ---

            /**
             * 获取的original
             */
            interface OriginalInfo {
                /**
                 * 服务器返回的original值
                 */
                original: string;
            }

            /**
             * 使用PKI认证
             */
            interface AuthenInfo {
                /**
                 * 唯一标识用户的ID
                 */
                userid: string;

                /**
                * 与userid一起验证请求的合法性
                */
                tokenid: string;

                /**
                * 获取到的token的有效期，单位为秒
                */
                expires: number;
            }

            // ---
            // 发现管理相关
            // ---

            /**
             * 文档发现状态
             */
            interface DocStatus {
                /**
                 * 是否开启
                 */
                status: boolean;
            }

            /**
             * 用户开启发现的文档
             */
            interface EnabledInfo {
                /**
                 * 文档路径
                 */
                path: string;

                /**
                 * 文档的id
                 */
                docid: string;

                /**
                 * 文件大小，文件夹则为-1
                 */
                size: number;

                /**
                 * 如果是文件，返回由客户端设置的文件本地修改时间
                 * 若未设置，返回modified的值
                 */
                client_mtim: number;
            }

            /**
             * 用户开启发现的文档
             */
            interface EnabledInfos {
                docinfos: ReadonlyArray<EnabledInfo>
            }

            // ---
            // 登录设备管理
            // ---

            /**
            * 所有设备信息
            */
            interface DeviceInfo {
                /**
                 * 设备名称
                 */
                name: string;

                /**
                 * 设备os类型
                 */
                ostype: string;

                /**
                 * 设备类型
                 */
                devicetype: number;

                /**
                 * 设备唯一标识号
                 */
                udid: string;

                /**
                 * 上次登录ip
                 */
                lastloginip: string;

                /**
                 * 上次登录时间
                 */
                lastlogintime: number;

                /**
                 * 缓存数据擦除状态
                 * 0表示正常；1表示正在请求擦除
                 */
                eraseflag: number;

                /**
                 * 上次成功清除缓存的时间
                 */
                lasterasetime: number;

                /**
                 * 设备禁用状态
                 * 0表示未禁用；1表示已禁用
                 */
                disableflag: number;

                /**
                 * 设备登录状态
                 * 0表示未登录；1表示已登录
                 */
                loginflag: number;

                /**
                 * 设备绑定状态
                 * 0表示未绑定；1表示已绑定
                 */
                bindflag: number;
            }

            /**
             * 所有设备信息
             */
            interface DeviceInfos {
                deviceinfos: ReadonlyArray<DeviceInfo>
            }

            /**
             * 设备状态
             */
            interface DeviceStatus {
                /**
                 * 数据擦除状态
                 * 0表示正常；1表示正在请求擦除
                 */
                eraseflag: number;

                /**
                * 设备禁用状态
                * 0表示正常；1表示处于禁用状态
                */
                disableflag: number;
            }

            /**
             * 获取验证码返回
             */
            interface VcodeInfo {
                /**
                 * 验证码唯一标识符
                 */
                uuid: string;

                /**
                 * 编码后的验证码字符串
                 */
                vcode: string;
            }

            /**
             * 入口文档信息返回结果
             */
            interface EntryDocInfo extends EntryDocByType {
                /**
                 * 站点信息
                 */
                ossInfo: {

                    /**
                     *  对象存储ID
                     */
                    ossId: string;
                    /**
                     *  对象存储名称
                     */
                    ossName: string;

                    /**
                     * 对象存储归属站点名称
                     */
                    siteName: string;

                    /**
                     * 对象存储状态
                     */
                    enable: boolean;
                }

                /**
                 * 是否需要下载水印
                 * 文件夹：只检查cid是否需要下载水印
                 * 文件：类型不支持加水印时返回false
                 */
                downloadwatermark: boolean;

            }

            /**
             * 服务器时间
             */
            interface ServerTime {
                /**
                 * 时间戳
                 */
                time: number;
            }

            /**
             * 获取登录验证码返回信息
             */
            interface AuthVcode {
                /**
                 * 用户的手机号或者邮箱
                 */
                authinfo: string;

                /**
                 * 向短信或者邮箱发送验证码的间隔（单位/秒）
                 */
                sendinterval: number;
            }
        }
    }
}