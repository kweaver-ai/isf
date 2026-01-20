import { NodeType, getNodeType, getNodeIconByType } from '@/core/organization';

/**
 * 组织节点
 */
interface OrganizationNode {
    /**
     * 名称
     */
    name: string;

    /**
     * 子部门的数量
     */
    subDepartmentCount: number;

    /**
     * 是否是组织
     */
    isOrganization: boolean;

    /**
     * 其他
     */
    [key: string]: any;
}

/**
 * 部门节点
 */
interface DepartmentNode {
    /**
     * 名称
     */
    name: string;

    /**
     * 子部门数量
     */
    subDepartmentCount: number;

    /**
     * 部门负责人
     */
    responsiblePerson: ReadonlyArray<UserNode>;

    /**
     * 其他
     */
    [key: string]: any;
}

/**
 * 用户节点
 */
interface UserNode {
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
 * 每次加载的用户数量
 */
export const UsersLimit = 150;

/**
 * 每次加载的用户起始页
 */
export const DefaultPage = 1;

/**
 * 根据节点获取图标
 * @param node 节点
 * @return 返回图标字体代码以及base64图片编码
 */
export function getIcon(node: any): { code: string } {
    return getNodeIconByType(getNodeType(node))
}

/**
 * 根据组织架构节点获取节点名称
 * @param node 组织架构节点
 */
export function getNodeName(node: OrganizationNode & DepartmentNode & UserNode): string {
    switch (getNodeType(node)) {
        case NodeType.ORGANIZATION:
        case NodeType.DEPARTMENT:
        case NodeType.UNDISTRIBUTED:
            return node.name;

        case NodeType.USER:
            return node.user.displayName;
    }
}