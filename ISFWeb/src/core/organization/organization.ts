import { includes } from 'lodash';
import { getDepartmentOfUsers, getDepartmentOfUsersCount, getDepParentPathById } from '../thrift/sharemgnt'
import __ from './locale'

/**
 * 组织树节点类型
 */
export enum NodeType {
    /**
     * 组织
     */
    ORGANIZATION,

    /**
     * 部门
     */
    DEPARTMENT,

    /**
     * 用户
     */
    USER,

    /**
     * 未分配组
     */
    UNDISTRIBUTED,
}

/**
 * 选中项数据类型
 */
export enum MixinNodeType {
    /**
     * 用户
     */
    User = 'user',

    /**
     * 应用账号
     */
    AppAccount = 'app',

    /**
     * 组织
     */
    Organization = 'organization',

    /**
     * 部门
     */
    Department = 'department',

    /**
     * 用户组
     */
    Group = 'group',

    /**
     * 未分配组
     */
    UnDistributed = 'undistributed',
}

/**
 * 额外根节点
 */
export interface ExtraRoot {
    /**
     * 名称
     */
    name: string;

    /**
     * id
     */
    id: string;

    /**
     * 是否是叶子节点
     */
    isLeaf: boolean;

    /**
     * 是否是文档库节点
     */
    isDocLib: boolean;
}

/**
 * 格式化节点（搜索结果/组织树节点）信息
 */
export interface FormatedNodeInfo {
    /**
     * 节点id
     */
    id: string;

    /**
     * 节点名称
     */
    name: string;

    /**
     * 节点类型
     */
    type: NodeType | any;

    /**
     * 用户登录名
     */
    account?: string;

    /**
     * 父路径
     */
    parent_path?: string;

    /**
     * 源数据
     */
    original: Record<string, any>;
}

/**
 * 获取节点类型
 * @param node 节点
 * @return 返回节点类型
 */
export function getNodeType(node: any): NodeType {
    switch (true) {
        case node.parentDepartId === '' || node.is_root || node.isOrganization || (node.hasOwnProperty('parent_dep_path') && node.parent_dep_path === ''):
            return NodeType.ORGANIZATION;

        case node.hasOwnProperty('responsiblePersons') || node.hasOwnProperty('is_root') || (node.hasOwnProperty('parent_dep_path')):
            return NodeType.DEPARTMENT;

        case !!(node.loginName || node.user || node.account):
            return NodeType.USER;

        case node.isUndistributedNode:
            return NodeType.UNDISTRIBUTED;

        default:
            return node.type

    }
}

/**
 * 判断节点是否是叶子节点
 * @param node 节点
 * @param selectType 可选用户范围
 * @return 返回是否是子节点
 */
export function isLeaf(node: any, selectType: ReadonlyArray<NodeType>): boolean {
    if (node.hasOwnProperty('isLeaf')) {
        return node.isLeaf
    }

    switch (getNodeType(node)) {
        case NodeType.ORGANIZATION:
        case NodeType.DEPARTMENT:
            /**
             * 组织/部门，一样的判断逻辑：
             * 1. 在selectType包含用户或部门的前提下：
             *    1）有子部门，则不是叶子节点
             *    2）selectType 包含用户，且有用户，则不是叶子节点
             * 2. 其他情况，是叶子节点
             */
            return !(
                (includes(selectType, NodeType.USER) || includes(selectType, NodeType.DEPARTMENT))
                && (
                    (node.depart_existed || node.subDepartmentCount)
                    || (includes(selectType, NodeType.USER) && (node.user_existed || node.subUserCount))
                )
            )

        case NodeType.USER:
            return true

        case NodeType.UNDISTRIBUTED:
            return !node.subUserCount
    }
}

/**
 * 获取节点图标
 */
export function getNodeIcon(node: any): { code: string } {
    switch (getNodeType(node)) {
        case NodeType.ORGANIZATION:
            return {
                code: '\uf008',
            }

        case NodeType.DEPARTMENT:
            return {
                code: '\uf009',
            }

        case NodeType.USER:
            return {
                code: '\uf007',
            }
    }
}

/**
 * 根据节点类型获取图标
 */
export function getNodeIconByType(nodeType: NodeType): { code: string } {
    switch (nodeType) {
        case NodeType.ORGANIZATION:
            return {
                code: '\uf008',
            }

        case NodeType.DEPARTMENT:
            return {
                code: '\uf009',
            }

        case NodeType.USER:
            return {
                code: '\uf01f',
            }

        case NodeType.UNDISTRIBUTED:
            return {
                code: '\uf009',
            }

    }
}

/**
 * 获取组织树部门节点内的用户
 */
export const getSubUsers = async ({ node, page, limit, isShowDisabledUsers = true }) => {
    if (!node.hasOwnProperty('subUserCount')) {
        const count = await getDepartmentOfUsersCount([node.id])
        node.subUserCount = count
    }

    const users = await getDepartmentOfUsers([node.id, (page - 1) * limit, limit])

    if (!isShowDisabledUsers) {
        // 过滤被禁用的用户
        const availableUsers = users.filter((userInfo) => (userInfo.user && userInfo.user.status !== 1))

        const disabledCount = users.length - availableUsers.length
        node.subUserCount -= disabledCount

        // 如果被禁用的用户数量超过50，继续加载下一页
        if (disabledCount >= 50) {
            const nextPage = page + 1
            const next = await getSubUsers({ node, page: nextPage, limit, isShowDisabledUsers })

            return {
                users: [...availableUsers, ...next.users],
                page: next.page,
            }
        }

        return {
            users: availableUsers,
            page,
        }
    }

    return {
        users,
        page,
    }
}

/**
 * 格式化组织树节点
 */
export const formatNodes = async (nodes: any[]): Promise<FormatedNodeInfo[]> => {
    let depPathMap = {}

    try {
        const depIds = new Set(nodes.map((node) => node.directDeptInfo ? node.directDeptInfo.departmentId : node.id))
        const depPaths = await getDepParentPathById(Array.from(depIds))

        depPaths.forEach((dep) => {
            depPathMap[dep.departmentId] = dep.parentPath
        })
    } catch{ }

    return nodes.map((node) => {
        const type = getNodeType(node)
        let parentPath = ''

        if (type === NodeType.DEPARTMENT || type === NodeType.USER) {
            parentPath = type === NodeType.USER ?
                node.directDeptInfo.departmentId === '-1' ?
                    __('未分配组')
                    : depPathMap[node.directDeptInfo.departmentId] ?
                        depPathMap[node.directDeptInfo.departmentId] + '/' + node.directDeptInfo.departmentName
                        : node.directDeptInfo.departmentName
                : (depPathMap[node.id] || '')
        }

        return {
            id: node.id,
            name: node.name || node.user.displayName,
            type,
            ...node.user ? { account: node.user.loginName } : {},
            parent_path: parentPath,
            original: node,
        }
    })
}

/**
 * 已选列表悬浮显示节点路径
 */
export const getDepName = (item: FormatedNodeInfo) => {
    try {
        switch (true) {
            case (item.type === NodeType.DEPARTMENT || item.type === MixinNodeType.Department) && !!item.parent_path:
                return `${item.name}(${__('部门')})-${item.parent_path}`
            case (item.type === NodeType.USER || item.type === MixinNodeType.User) && !!item.parent_path:
                return `${item.name}(${item.account})-${item.parent_path}`
            case item.type === MixinNodeType.AppAccount:
                return `${item.name}${__('（应用账户）')}`
            case item.type === MixinNodeType.Group:
                return `${item.name}${__('（用户组）')}`
            default:
                return item.name
        }
    } catch {
        return item.name
    }
}

/**
 * 节点类型映射为选中类型
 */
export const nodeTypeMaptoMixinType = (type: NodeType) => {
    switch (type) {
        case NodeType.ORGANIZATION:
            return MixinNodeType.Organization;
        case NodeType.DEPARTMENT:
            return MixinNodeType.Department;
        case NodeType.USER:
            return MixinNodeType.User;
        case NodeType.UNDISTRIBUTED:
            return MixinNodeType.UnDistributed;
        default:
            return type
    }
}