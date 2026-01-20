import { consolehttp } from '../../../openapiconsole'
import { UpdateDefaultConfigs, CheckDefaultPwd, GetDepartmentRoots, DeleteDepartment } from './types'

/**
 * 更新初始化配置
 */
export const updateDefaultConfigs: UpdateDefaultConfigs = (data, options) => {
    return consolehttp('put', ['user-management', 'v1', 'management', 'configs', Object.keys(data || {}).join(',')], data, {}, options)
}

/**
 * 预检测密码格式
 */
export const checkDefaultPwd: CheckDefaultPwd = ({ password }, options) => {
    return consolehttp('get', ['user-management', 'v1', 'management', 'default-pwd-valid'], undefined, { password }, options)
}

/**
 * 获取根部门信息
 */
export const getDepartmentRoots: GetDepartmentRoots = ({ departmentId, fields, role, offset, limit })=> {
    return consolehttp('get', ['user-management', 'v1', 'department-members', departmentId, fields.join(',')], {}, { role, offset, limit }, {})
}

/**
 * 删除组织、部门
 */
export const deleteDepartment: DeleteDepartment = (id, options?) => {
    return consolehttp('delete', ['user-management', 'v1', 'management', 'departments', id], undefined, {}, options)
}

/**
 * 根据条件筛选用户
 */
export const searchUsers = ({ role, department_id, code, name, account, manager_name, direct_department_code, position, offset, limit, fields = ['name', 'account', 'code', 'remark', 'roles', 'csf_level', 'csf_level2','auth_type', 'priority', 'manager', 'position', 'frozen', 'created_at', 'enabled', 'parent_deps'].join(',') }: { role: string; department_id?: string; code?: string; name?: string; account?: string; manager_name?: string; direct_department_code?: string; position?: string; offset?: number; limit?: number; fields?: string }, options?) => {
    return consolehttp('get', ['user-management', 'v1', 'console', 'search-users', fields], undefined, { department_id, code, name, account, manager_name, direct_department_code, position, offset, limit, role }, options )
}

/**
 * 部门列举
 */
export const getDepartments = ({ department_id, role, offset, limit = 100, fields = ['departments'].join(',') }: {department_id: string; role: string; offset: number; limit?: number; fields?: string}) => {
    return consolehttp('get', ['user-management', 'v1', 'management', 'department-members', department_id, fields], undefined, { role, offset, limit })
}

/**
 * 部门搜索
 */
export const searchDepartments = ({ role, code, name, manager_name, direct_department_code, remark, enabled, offset, limit = 100, fields = ['name', 'code', 'remark', 'manager', 'enabled', 'parent_deps', 'email'].join(',') }: { role: string; code?: number; name?: string; manager_name?: string; direct_department_code?: number; remark?: string; enabled?: boolean; offset: number; limit?: number; fields?: string }) => {
    return consolehttp('get', ['user-management', 'v1', 'console', 'search-departments', fields], undefined, { role, code, name, manager_name, direct_department_code, remark, enabled, offset, limit })
}

/**
 * 获取配置
 */
export const getLevelConfig = ({ fields }: {fields?: string }) => {
    return consolehttp('get', ['user-management', 'v1', 'configs', fields], undefined, {})
}

/**
 * 更新配置
 */
export const setLevelConfig = ({ fields, csf_level_enum, csf_level2_enum }: {fields?: string, csf_level_enum?: Array<{name: string; value: number}>, csf_level2_enum?: Array<{name: string; value: number}> }) => {
    return consolehttp('put', ['user-management', 'v1', 'management', 'configs', fields], { csf_level_enum, csf_level2_enum }, {})
}   