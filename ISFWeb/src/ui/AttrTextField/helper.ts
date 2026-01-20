
/**
 * 展开项显示状态
 */
export enum ExpandStatus {
    // 隐藏
    HIDE = 0,

    // 显示展开
    EXPAND = 1,

    // 显示收起
    COLLAPSE = 2,
}

/**
 * 文件属性类型
 */
export enum AttrType {
    // 层级
    LEVEL = 0,
    // 枚举
    ENUM = 1,
    // 数字
    NUMBER = 2,
    // 文本
    TEXT = 3,
    // 时间
    TIME = 4,
}