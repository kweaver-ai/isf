import { OpenAPI } from '../..'

/**
 * 用户范围
 */
export enum UserType {
    /**
     * 用户
     */
    User = 'user',

    /**
     * 部门
     */
    Department = 'department',

    /**
     * 部门
     */
    Dept = 'dept',

    /**
     * 用户组
     */
    Group = 'group',

    /**
     * 所有用户
     */
    All = 'all',
}

/**
 * 更新模式
 */
export enum Mode {
    /**
     * 强制更新
     */
    Mandatory,

    /**
     * 非强制更新
     */
    NonMandatory,
}

/**
 * 终端类型
 */
export enum ClientType {
    Windows = 'win',
    Linux = 'linux',
    MAC = 'mac',
    iOS = 'ios',
    Android = 'android',
    OfficePlugin = 'office',
}

/**
 * 用户
 */
export interface User {
    /**
     * 用户ID
     */
    id: string;

    /**
     * 用户类型
     */
    type: UserType;

    /**
     * 用户名称
     */
    name: string;
}

/**
 * 策略
 */
export interface Strategy {
    /**
     * 策略ID
     */
    id: string;

    /**
     * 策略名称
     */
    name: string;

    /**
     * 用户范围
     */
    users: ReadonlyArray<User>;

    /**
     * 更新模式
     */
    mode: Mode | null;

    /**
     * 是否静默升级
     */
    silence?: boolean | null;

    /**
     * 客户端类型列表
     */
    client_type: ReadonlyArray<ClientType>;

    /**
     * 备注
     */
    remark: string;
}

/**
 * 获取策略列表参数
 */
interface GetUpdateStrategyParams {
    /**
     * 页面偏移量
     */
    offset: number;

    /**
     * 单页显示记录总数
     */
    limit: number;

    /**
     * 关键字
     */
    keyword?: string;

    /**
     * 客户端更新模式
     */
    mode?: Mode;

    /**
     * 策略生效客户端类型
     */
    client_type?: ClientType;
}

/**
 * 获取策略列表结果
 */
interface GetUpdateStrategyResult {
    /**
     * 记录总数
     */
    total_count: number;

    /**
     * 策略记录列表
     */
    entries: ReadonlyArray<Strategy>;
}

export type GetUpdateStrategy = OpenAPI<GetUpdateStrategyParams, GetUpdateStrategyResult>

export type CreateUpdateStrategy = OpenAPI<Partial<Strategy>, {id: string; invalid_users: ReadonlyArray<User>}>

export type DeleteUpdateStrategy = OpenAPI<string, void>

export type EditUpdateStrategy = OpenAPI<{id: string; data: Partial<Strategy>}, {invalid_users: ReadonlyArray<User>}>