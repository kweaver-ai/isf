import { consolehttp } from '../../../openapiconsole';

/**
 * 获取应用账户列表
 */
export const getUseAccountMgnt: Core.APIs.Console.Account.GetUseAccountMgnt = ({ limit, offset, direction, sort, keyword }, options?) => {
    return consolehttp('get', ['user-management', 'v1', 'apps'], undefined, { limit, offset, direction, sort, keyword }, { options })
}

/**
 * 注册应用账户
 */
export const createUseAccountMgnt: Core.APIs.Console.Account.CreateUseAccountMgnt = ({ name, password }, options?) => {
    return consolehttp('post', ['user-management', 'v1', 'apps'], { name, password }, {}, { options })
}

/**
 * 删除应用账户
 */
export const delUseAccountMgnt: Core.APIs.Console.Account.DelUseAccountMgnt = ({ id }, options?) => {
    return consolehttp('delete', ['user-management', 'v1', 'apps', id], undefined, {}, { options })
}

/**
 * 更新应用账户
 */
export const setUseAccountMgnt: Core.APIs.Console.Account.SetUseAccountMgnt = ({ id, fields, name, password }, options?) => {
    return consolehttp('put', ['user-management', 'v1', 'apps', id, fields], { name, password }, {}, { options })
}

/**
 * 生成token
 */
export const generateToken = ({id}, options?) => {
    return consolehttp('post', ['user-management', 'v1', 'console', 'app-tokens'], { id }, {}, { options })
}

/**
 * 获取应用账户用户交接权限信息
 */
export const getUserTransferPerm: Core.APIs.Console.Account.GetUserTransferPerm = (options?) => {
    return consolehttp('get', ['user-transfer', 'v1', 'perm', 'app'], undefined, {}, options);
}

/**
 * 应用账户增加用户交接权限
 */
export const addUserTransferPermById: Core.APIs.Console.Account.AddUserTransferPermById = ({ id }, options?) => {
    return consolehttp('put', ['user-transfer', 'v1', 'perm', 'app', id], undefined, {}, options);
}

/**
 * 删除应用账户用户交接权限信息
 */
export const deleteUserTransferPerm: Core.APIs.Console.Account.DeleteUserTransferPerm = ({ id }, options?) => {
    return consolehttp('delete', ['user-transfer', 'v1', 'perm', 'app', id], undefined, {}, options);
}

/**
 * 获取指定应用账户用户交接权限
 */
export const getUserTransferPermById: Core.APIs.Console.Account.GetUserTransferPermById = ({ id }, options?) => {
    return consolehttp('get', ['user-transfer', 'v1', 'perm', 'app', id], undefined, {}, options);
}

/**
 * 获取应用账户文档域管理权限信息
 */
export const getDocDomainPerm: Core.APIs.Console.Account.GetDocDomainPerm = (options?) => {
    return consolehttp('get', ['document-domain-management', 'v1', 'domain', 'perm', 'app'], undefined, {}, options);
}

/**
 * 应用账户增加文档域管理权限
 */
export const addDocDomainPermById: Core.APIs.Console.Account.AddDocDomainPermById = ({ id }, options?) => {
    return consolehttp('put', ['document-domain-management', 'v1', 'domain', 'perm', 'app', id], undefined, {}, options);
}

/**
 * 删除应用账户文档域管理权限信息
 */
export const deleteDocDomainPerm: Core.APIs.Console.Account.DeleteDocDomainPerm = ({ id }, options?) => {
    return consolehttp('delete', ['document-domain-management', 'v1', 'domain', 'perm', 'app', id], undefined, {}, options);
}

/**
 * 获取指定应用账户文档域管理权限信息
 */
export const getDocDomainPermById: Core.APIs.Console.Account.GetDocDomainPermById = ({ id }, options?) => {
    return consolehttp('get', ['document-domain-management', 'v1', 'domain', 'perm', 'app', id], undefined, {}, options);
}

/**
 * 获取应用账户组织架构管理权限信息
 */
export const getOrgManagePerm: Core.APIs.Console.Account.GetOrgManagePerm = (app_id, org_manage_type, options?) => {
    return consolehttp('get', ['user-management', 'v1', 'app-perms', app_id, org_manage_type], undefined, {}, { options })
}

/**
 * 指定应用账户的组织架构管理权限
 */
export const setOrgManagePerm: Core.APIs.Console.Account.SetOrgManagePerm = ({ app_id, org_manage_type, allowed }, options?) => {
    return consolehttp('put', ['user-management', 'v1', 'app-perms', app_id, org_manage_type], [{ subject: app_id, object: org_manage_type, perms: allowed }], undefined, { options })
}

/**
 * 删除指定应用账户的组织架构管理
 */
export const delOrgManagePerm: Core.APIs.Console.Account.DelOrgManagePerm = ({ app_id, org_manage_type }, options?) => {
    return consolehttp('delete', ['user-management', 'v1', 'app-perms', app_id, org_manage_type], undefined, {}, { options })
}

/**
 * 获取所有具备获取任意用户访问令牌权限的应用账户
 */
export const getUserTokenPermList: Core.APIs.Console.Account.GetUserTokenPermList = (options?) => {
    return consolehttp('get', ['authentication', 'v1', 'access-token-perm', 'app'], undefined, {}, options);
}

/**
 * 配置应用账户获取任意用户访问令牌的权限
 */
export const addUserTokenPermById: Core.APIs.Console.Account.AddUserTokenPermById = ({ id }, options?) => {
    return consolehttp('put', ['authentication', 'v1', 'access-token-perm', 'app', id], undefined, {}, options);
}

/**
 * 删除应用账户获取任意用户访问令牌的权限
 */
export const deleteUserTokenPerm: Core.APIs.Console.Account.DeleteUserTokenPerm = ({ id }, options?) => {
    return consolehttp('delete', ['authentication', 'v1', 'access-token-perm', 'app', id], undefined, {}, options);
}