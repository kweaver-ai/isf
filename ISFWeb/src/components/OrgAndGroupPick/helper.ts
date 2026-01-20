import { ProtLevel } from '@/core/systemprotectionlevel'
import { getUserInfoCache, UserInfo } from '@/core/user'
import { SystemRoleType } from '@/core/role/role'

/**
 * 组件类型
 */
export enum Type {
    /**
     * 用户
     */
    User,

    /**
     * 部门
     */
    Department,
}

/**
 * 组件tab类型
 */
export enum TabType {
    /**
     * 组织
     */
    Org = 'org',

    /**
     * 用户组
     */
    Group = 'group',

    /**
     * 匿名用户
     */
    Anonymous = 'anonymous',

    /**
     * 应用账户
     */
    App = 'app',
}

/**
 * 选中项数据类型
 */
export enum SelectionType {
    /**
     * 用户
     */
    User = 'user',

    /**
     * 部门
     */
    Department = 'department',

    /**
     * 用户组
     */
    Group = 'group',

    /**
     * 匿名用户
     */
    Anonymous = 'anonymous',

    /**
     * 应用账户
     */
    App = 'app',
}

/**
 * 选中项数据
 */
export interface Selection {
    /**
     * 类型
     */
    type: SelectionType;

    /**
     * 数据id
     */
    id: string;

    /**
     * 名称
     */
    name: string;

    /**
     * 源数据
     */
    original?: any;
}

/**
 * 选中节点数据
 */
export interface NodeData {
    /**
     * 用户id
     */
    id?: string;

    /**
     * 部门id
     */
    departmentId?: string;

    /**
     * 用户显示名
     */
    displayName?: string;

    /**
     * 部门名
     */
    name?: string;

    /**
     * 部门名
     */
    departmentName?: string;

    /**
     * 用户数据
     */
    user?: {
        displayName: string;
    };

    /**
     * 数据的类型
     */
    type?: SelectionType;
}

/**
* 过滤tab，屏蔽角色不足的
*/
export function filterTabType(level: ProtLevel, tabType: ReadonlyArray<TabType>): ReadonlyArray<TabType> {
    let disabledType: Array<TabType> = []

    if (level === ProtLevel.MoreConfidential) {
        // 系统保护等级为机密级增强，屏蔽用户组
        disabledType = [TabType.Group]
    } else {
        // 超级管理员和系统管理员、安全管理员、审计管理员 以外的管理员，屏蔽用户组+应用账户
        const roles = getUserInfoCache(UserInfo.roles) as ReadonlyArray<{ id: SystemRoleType }>
        const AllowRoleType = [SystemRoleType.Supper, SystemRoleType.Admin, SystemRoleType.Securit, SystemRoleType.Audit]

        if (!roles.some(({ id }) => AllowRoleType.includes(id))) {
            disabledType = [TabType.Group, TabType.App]
        }
    }

    return tabType.filter((tab) => !disabledType.includes(tab))
}