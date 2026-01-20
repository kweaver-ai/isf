/**
 * 被选择项数据对象
 */
export interface DataItem {
    /**
     * 类型
     */
    type: string;

    /**
     * id
     */
    id: string;

    /**
     * 名称
     */
    name: string;
}

/**
 * 操作类型（个人/部门）
 */
export enum DocType {
    /**
     * 个人
     */
    User,

    /**
     * 部门
     */
    Department,
}

/**
 * 文档库
 */
export const DocTypeText = {
    [DocType.User]: '个人文档库',
    [DocType.Department]: '部门文档库',
}

/**
 * 选择范围的类型
 */
export enum ScopeType {
    /**
     * 所有文档库
     */
    All,

    /**
     * 自定义选择
     */
    Custom,
}

/**
 * 显示内容
 */
export enum ConfigStatus {
    /**
     * 无显示内容
     */
    None,

    /**
     * 显示选择范围类型
     */
    ScopePicker,

    /**
     * 显示自定义范围
     */
    CustomPicker,
}
