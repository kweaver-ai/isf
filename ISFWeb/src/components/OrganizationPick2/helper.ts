
import { NodeType, ExtraRoot } from '@/core/organization';

export interface OrganizationPick2Props {
    /**
     * 用户id
     */
    userid?: string;

    /**
     * 搜索框预设文案
     */
    placeholder?: string;

    /**
     * 是否自动聚焦
     */
    autoFocus?: boolean;

    /**
     * 搜索框宽度
     */
    searchBoxwWidth?: number;

    /**
     *  组织树是否禁用
     */
    disabled?: boolean;

    /**
     * 选择节点的类型
     */
    selectType: ReadonlyArray<NodeType>;

    /**
     * 提示描述
     */
    describeTip?: (() => React.ReactNode) | string;

    /**
     * 已选项
     */
    selections: ReadonlyArray<any>;

    /**
     * 额外的根节点
     */
    extraRoots?: ReadonlyArray<ExtraRoot>;

    /**
     * 传入数据时的转换规则
     */
    converterIn?: (selected: Node) => any;

    /**
     * 传出数据时的转换规则
     */
    convererOut?: (selected: Node) => any;

    /**
     * 选中项改变
     */
    onRequestChangeSelections: (selections: ReadonlyArray<any>) => void;

    /**
     * 格式化选中项显示
     */
    formatSelectedItem?: (item: Node) => React.ReactNode;
}

export interface OrganizationPick2State {
    /**
     * 已选项
     */
    selections: ReadonlyArray<any>;
}

/**
 * 格式化之后的节点信息
 */
export interface Node {
    /**
     * 名称
     */
    name: string;

    /**
     * 部门id
     */
    id: string;

    /**
     * 类型
     */
    type: number;

    /**
     * 全路径
     */
    path?: string;

    /**
     * 原始数据
     */
    original?: any;
}

/**
 * 组织树里的节点信息
 */
export interface TreeDepNode {
    /**
     * 部门id
     */
    id: string;

    /**
     * 是否是组织
     */
    isOrganization?: boolean;

    /**
     * 部门名称
     */
    name: string;

    /**
     * 部门负责人
     */
    responsiblePersons: ReadonlyArray<any>;

    /**
     * 子部门数量
     */
    subDepartmentCount: number;

    /**
     * 部门用户数
     */
    subUserCount: number;

    /**
     * 其他
     */
    [key: string]: any;
}

/**
 * 组织树中用户节点信息
 */
export interface TreeUserNode {
    /**
     * 直属部门的信息
     */
    directDeptInfo: {
        /**
         * 直属部门名称
         */
        departmentName: string;

        /**
         * 直属部门负责人
         */
        responsiblePersons: ReadonlyArray<any>;

        /**
         * 直属部门id
         */
        departmentId: string;
    };

    /**
     * 用户id
     */
    id: string;

    /**
     * 用户信息
     */
    user: {
        /**
         * 显示名
         */
        displayName: string;

        /**
         * 登录名
         */
        loginName: string;

        /**
         * 其他
         */
        [key: string]: any;
    };

    /**
     * 其他
     */
    [key: string]: any;
}

/**
 * 搜索框选择中的部门信息
 */
export interface SearchDepSelectDepInfo {
    /**
     * 部门id
     */
    departmentId: string;

    /**
     * 部门名称
     */
    departmentName: string;

    /**
     * 父部门id
     */
    parentDepartId: string;

    /**
     * 父部门名称
     */
    parentDepartName: string;

    /**
     * 父部门路径
     */
    parentPath: string;

    /**
     * 部门负责人
     */
    responsiblePersons: ReadonlyArray<any>;

    /**
     * 其他
     */
    [key: string]: any;
}

/**
 * 搜索框选择中的用户信息
 */
export interface SearchDepSelectUserInfo {
    /**
     * 父部门id
     */
    departmentIds: ReadonlyArray<string>;

    /**
     * 父部门名称
     */
    departmentNames: ReadonlyArray<string>;

    /**
     * 父部门路径
     */
    departmentPaths: ReadonlyArray<string>;

    /**
     * 显示名
     */
    displayName: string;

    /**
     * 用户id
     */
    id: string;

    /**
     * 登录名
     */
    loginName: string;

    /**
     * 其他
     */
    [key: string]: any;
}