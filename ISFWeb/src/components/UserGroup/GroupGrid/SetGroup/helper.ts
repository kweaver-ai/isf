import __ from './locale';

/**
 * 校验状态
 */
export enum ValidateState {
    /**
     * 正常
     */
    Normal,

    /**
     * 为空
     */
    Empty,

    /**
     * 输入不合法
     */
    InvalidName,

    /**
     * 名称已存在
     */
    NameConfilct,
}

/**
 * 校验提示信息
 */
export const ValidateMessages = {
    [ValidateState.Empty]: __('此项不允许为空。'),
    [ValidateState.InvalidName]: __('用户组名称不能包含 空格 或 \\ / : * ? " < > | 特殊字符，长度不能超过128个字符。'),
    [ValidateState.NameConfilct]: __('此名称已存在。'),
}

/**
 * 用户组
 */
export interface UserGroupType {
    /**
     * id
     */
    id: string;

    /**
     * 名称
     */
    name: string;
}