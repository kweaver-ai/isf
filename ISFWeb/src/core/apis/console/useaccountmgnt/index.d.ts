import { string } from "prop-types";

declare namespace Core {
    namespace APIs {
        namespace Console {
            namespace Account {
                /**
                 * 获取应用账户列表
                 */
                type GetUseAccountMgnt = Core.APIs.OpenAPI<{
                    /**
                     * 加载限制
                     */
                    limit: number;

                    /**
                     * 加载起始点
                     */
                    offset: number;

                    /**
                     * 排序升降
                     */
                    direction: string;

                    /**
                     * 排序类型
                     */
                    sort: string;
                }, {
                    Account: ReadonlyArray<{
                        id: string;
                        name: string;
                    }>,
                    total_count: number;
                }>;

                /**
                 * 注册应用账户
                 */
                type CreateUseAccountMgnt = Core.APIs.OpenAPI<{
                    /**
                      * 账户名称
                      */
                    name: string;

                    /**
                     * 账户密码
                     */
                    password: string;
                }, { id: string }>;

                /**
                 * 删除应用账户
                 */
                type DelUseAccountMgnt = Core.APIs.OpenAPI<{
                    /**
                     * 账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 编辑应用账户
                 */
                type SetUseAccountMgnt = Core.APIs.OpenAPI<{
                    /**
                     * 需要编辑的字段: name,password
                     */
                    fields: string;

                    /**
                     * 账户id
                     */
                    id: string;

                    /**
                     * 账户名称
                     */
                    name?: string;

                    /**
                     * 账户密码
                     */
                    password?: string;
                }, void>;

                /**
                 * 获取应用账户权限信息
                 */
                type GetUsePermissions = Core.APIs.OpenAPI<
                    void,
                    Core.APIs.Console.NetworkList>

                /**
                 * 指定应用账户的文档库类型配置权限
                 */
                type SetUsePermission = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    app_id: string;

                    /**
                     * 文档库类型
                     */
                    doc_lib_type: string;

                    /**
                     * 权限
                     */
                    allowed: ReadonlyArray<string>;

                    /**
                    * 到期时间
                    */
                    expires_at: string;
                }, void>

                /**
                * 删除应用账户的文档库类型配置权限
                */
                type DelUsePermission = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    app_id: string;

                    /**
                     * 文档库类型
                     */
                    doc_lib_type: string;
                }, void>

                /**
                 * 获取应用账户用户交接权限信息
                 */
                type GetUserTransferPerm = Core.APIs.OpenAPI<
                    void,
                    ReadonlyArray<string>>;

                /**
                 * 应用账户增加用户交接权限
                 */
                type AddUserTransferPermById = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 删除应用账户用户交接权限信息
                 */
                type DeleteUserTransferPerm = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 获取指定应用账户用户交接权限
                 */
                type GetUserTransferPermById = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 获取应用账户文档域管理权限信息
                 */
                type GetDocDomainPerm = Core.APIs.OpenAPI<
                    void,
                    ReadonlyArray<string>>;

                /**
                 * 应用账户增加文档域管理权限
                 */
                type AddDocDomainPermById = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 删除应用账户文档域管理权限信息
                 */
                type DeleteDocDomainPerm = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 获取指定应用账户文档域管理权限信息
                 */
                type GetDocDomainPermById = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 获取应用账户组织架构管理权限信息
                 */
                type GetOrgManagePerm = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    app_id: string;

                    /**
                     * 组织架构管理权限类型
                     */
                    org_manage_type: Readonly<string>;
                },  void>

                /**
                 * 指定应用账户的组织架构管理权限
                 */
                type SetOrgManagePerm = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    app_id: string;

                    /**
                     * 组织架构管理权限类型
                     */
                    org_manage_type: string;

                    /**
                     * 权限信息
                     */
                    allowed: ReadonlyArray<string>;
                }, void>

                /**
                * 删除应用账户的组织架构管理权限
                */
                type DelOrgManagePerm = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    app_id: string;

                    /**
                     * 组织架构管理权限类型
                     */
                    org_manage_type: string;
                }, void>

                /**
                 * 获取所有具备获取任意用户访问令牌权限的应用账户
                 */
                type GetUserTokenPermList = Core.APIs.OpenAPI<
                    void,
                    ReadonlyArray<string>>;

                /**
                 * 配置应用账户获取任意用户访问令牌的权限
                 */
                type AddUserTokenPermById = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;

                /**
                 * 删除应用账户获取任意用户访问令牌的权限
                 */
                type DeleteUserTokenPerm = Core.APIs.OpenAPI<{
                    /**
                     * 应用账户id
                     */
                    id: string;
                }, void>;
            }
        }
    }
}