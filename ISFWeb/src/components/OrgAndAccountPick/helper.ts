/**
 * 组件Tab栏类型
 */
export enum TabType {

    /**
     * 应用账户
     */
    AppAccount = 'app',

    /**
     * 组织结构
     */
    Org = 'org',
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
     * 未分配组
     */
    UnDistributed = 'undistributed',
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
 * 获取节点类型
 * @param node 节点
 * @return 返回节点类型
 */
export function getNodeType(node: any): SelectionType {
    switch (true) {
        case node.parentDepartId === '' || node.is_root:
            return SelectionType.Organization;

        case node.hasOwnProperty('responsiblePersons'):
            return SelectionType.Department;

        case !!(node.loginName || node.user):
            return SelectionType.User;

        case node.isUndistributedNode:
            return SelectionType.UnDistributed;
    }
}