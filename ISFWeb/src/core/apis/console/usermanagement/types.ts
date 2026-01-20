import { OpenAPI } from '../..'
import { UserManagementErrorCode } from '../../openapiconsole/errorcode';

/**
 * 配置
 */
export enum Configs {
    /**
     * 用户初始密码
     */
    DefaultUserPwd = 'default_user_pwd',
}

/**
 * 更新初始化配置
 */
export type UpdateDefaultConfigs = Core.APIs.OpenAPI<Record<Configs, any>, void>;

/**
 * 检测初始密码格式
 */
export type CheckDefaultPwd = Core.APIs.OpenAPI<{ password: string }, { result: boolean; err_msg?: string }>;

/**
 * 删除组织、部门
 */
export type DeleteDepartment = OpenAPI<string, void>

/**
 * 用户、部门、用户组、应用账号不存在
 */
export const notExisted = [UserManagementErrorCode.GroupMemberNotExisted, UserManagementErrorCode.DepartmentNotExisted, UserManagementErrorCode.UserGroupNotFound, UserManagementErrorCode.AppAccountNotFound]

/**
 * 部门成员列举请求接口参数
 */
interface GetDepartmentRootsParam {
    /**
     * 部门id
     */
    departmentId: string;

    /**
     * 获取成员信息字段名
     */
    fields: ReadonlyArray<string>;

    /**
     * 用户角色
     */
    role: string;

    /**
     * 获取数据起始下标
     */
    offset?: number;

    /**
     * 获取数据量
     */
    limit?: number;
}

/**
 * 用户信息数据
 */
interface UserEntry {
    /**
     * 部门id
     */
    id: string;

    /**
     * 部门名称
     */
    name: string;

    /**
     * 固定为“user”
     */
    type: string;
}

/**
 * 部门信息数据
 */
export interface DepartmentEntry {
    /**
     * 部门id
     */
    id: string;

    /**
     * 部门名称
     */
    name: string;

    /**
     * 固定为“department”
     */
    type: string;

    /**
     * 是否为根部门
     */
    is_root: boolean;

    /**
     * 是否拥有子用户
     */
    user_existed: boolean;

    /**
     * 是否拥有子部门
     */
    depart_existed: boolean;
}

/**
 * 部门成员列举请求接口响应
 */
interface GetDepartmentRootsResult {
    /**
     * 用户信息分页数据
     */
    users?: {
        /**
         * 总条目数（忽略offset和limit）
         */
        total_count: number;

        /**
         * 用户条目列表
         */
        entries: ReadonlyArray<UserEntry>;
    };

    /**
     * 部门信息分页数据
     */
    departments?: {
        total_count: number;
        entries: ReadonlyArray<DepartmentEntry>;
    };
}

export type GetDepartmentRoots = OpenAPI<GetDepartmentRootsParam, GetDepartmentRootsResult>

export type LevelConfig =  {
    /**
     * 密级列表
     */
    csf_level_enum?: Array<{name: string; value: string}>;

    /**
     * 密级2列表
     */
    csf_level2_enum?: Array<{name: string; value: string}>;    

    /**
     * 是否展示密级2
     */
    show_csf_level2?: boolean;
}