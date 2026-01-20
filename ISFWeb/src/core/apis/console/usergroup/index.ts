import { consolehttp } from '../../../openapiconsole'

/**
 * 获取用户组管理列表
 */
export const getUserGroups: Core.APIs.Console.UserGroup.GetUserGroups = ({ direction, sort, offset, limit, keyword }, options?) => {
    return consolehttp('get', ['user-management', 'v1', 'management', 'groups'], undefined, { direction, sort, offset, limit, keyword }, options)
}

/**
 * 根据id获取用户组详情
 */
export const getUserGroupById: Core.APIs.Console.UserGroup.GetUserGroupById = ({ id }, options?) => {
    return consolehttp('get', ['user-management', 'v1', 'management', 'groups', id], undefined, {}, options)
}

/**
 * 创建用户组
 */
export const createUserGroup: Core.APIs.Console.UserGroup.CreateUserGroup = (param, options?) => {
    return consolehttp('post', ['user-management', 'v1', 'management', 'groups'], param, {}, options)
}

/**
 * 编辑用户组
 */
export const editUserGroup: Core.APIs.Console.UserGroup.EditUserGroup = ({ id, fields, name, notes }, options?) => {
    return consolehttp('put', ['user-management', 'v1', 'management', 'groups', id, fields], { name, notes }, {}, options)
}

/**
 * 删除用户组
 */
export const deleteUserGroup: Core.APIs.Console.UserGroup.DeleteUserGroup = ({ id }, options?) => {
    return consolehttp('delete', ['user-management', 'v1', 'management', 'groups', id], undefined, {}, options)
}

/**
 * 获取用户组成员列表
 */
export const getGroupMembers: Core.APIs.Console.UserGroup.GetGroupMembers = (id, { direction, sort, offset, limit, keyword }, options?) => {
    return consolehttp('get', ['user-management', 'v1', 'management', 'group-members', id], undefined, { direction, sort, offset, limit, keyword }, options)
}

/**
 * 获取用户组成员(搜索模糊匹配用户显示名)
 */
export const getGroupMembersByUserMatch: Core.APIs.Console.UserGroup.GetGroupMembersByUserMatch = ({ group_id, key, offset, limit }) => {
    return consolehttp('get', ['user-management', 'v1', 'console', 'search-users-in-group'], undefined, { group_id, key, offset, limit })
}

/**
 * 添加成员
 */
export const addGroupMembers: Core.APIs.Console.UserGroup.AddGroupMembers = (id, { members }, options?) => {
    return consolehttp('post', ['user-management', 'v1', 'management', 'group-members', id], { method: 'POST', members }, {}, options)
}

/**
 * 删除成员
 */
export const deleteGroupMembers: Core.APIs.Console.UserGroup.DeleteGroupMembers = (id, { members }, options?) => {
    return consolehttp('post', ['user-management', 'v1', 'management', 'group-members', id], { method: 'DELETE', members }, {}, options)
}

/**
 * 获取用户组列表
 */
export const getGroups = ({ offset, limit }, options?): Core.APIs.Console.UserGroup.GetGroups => {
    return consolehttp('get', ['user-management', 'v1', 'groups'], undefined, { offset, limit }, options)
}

/**
 * 获取部门下的子用户信息
 */
export const searchInGroup = ({ keyword, type, offset, limit }, options?): Core.APIs.Console.UserGroup.SearchInGroup => {
    return consolehttp('get', ['user-management', 'v1', 'search-in-group'], undefined, { keyword, type, offset, limit }, options)
}