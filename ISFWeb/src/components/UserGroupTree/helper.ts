export enum Type {
    /**
     * 用户类型
     */
    User = 'user',

    /**
     * 部门类型
     */
    Department = 'department',

    /**
     * 用户组类型
     */
    Group = 'group',
}

/**
 * 用户组信息
 */
export interface UserGroup {
    /**
     * 用户组id
     */
    id: string;

    /**
     * 用户组名称
     */
    name: string;

    /**
     * 文本
     */
    node?: string;
}

/**
 * 搜索结果
 */
export interface Result {
    /**
     * 类型
     */
    type: Type;

    /**
     * 用户组id
     */
    id: string;

    /**
     * 用户组名
     */
    name: string;

    /**
     * 源数据
     */
    origin: UserGroup;

}

export const Limit = 50;

export const formatUserGroup = (userGroup: UserGroup): Result => {
    const { id, name } = userGroup;

    return {
        type: Type.Group,
        id,
        name,
        origin: userGroup,
    }
}