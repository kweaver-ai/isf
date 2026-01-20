import { ProtLevel } from '@/core/systemprotectionlevel';
import __ from './locale';

export enum Status {
    // 验证合格
    OK,

    // 密码错误锁定次数超出范围
    COUNT_RANGE_ERROR,

    /**
     * 锁定时间超出范围
     */
    PASSWD_LOCK_TIME_ERROR,

    // 密码长度超出范围
    COUNT_PWD_RANGE_ERROR,

    // 自定义密级名称含有特殊字符
    FORBIDDEN_SPECIAL_CHARACTER,

    //  密级个数超出限制
    SECU_OUT_SUM,

    // 密级名重名
    DUPLICATE_NAMES_ERROR,

    // 密级名为空
    EMPTY_NAME,
}

export interface PwdPolicy {
    /**
     * 是否屏蔽弱密码选项
     */
    weak_pwd_disabled: boolean;

    /**
     * 强密码最小长度
     */
    min_strong_pwd_length: number;

    /**
     * 密码错误最大次数
     */
    max_err_count: number;
}

/**
 * 密码有效期
 */
export enum PwdValidity {
    /**
     * 初始值
     */
    None = 0,

    /**
     * 1天
     */
    OneDay = 1,

    /**
     * 3天
     */
    TreeDays = 3,

    /**
     * 7天
     */
    SevenDays = 7,

    /**
     * 1个月
     */
    OneMonth = 30,

    /**
     * 3个月
     */
    TreeMonths = 90,

    /**
     * 6个月
     */
    SixMonths = 180,

    /**
     * 12个月
     */
    TwelveMonths = 365,

    /**
     * 永久
     */
    Permanent = -1,
}

/**
 * 验证自定义密级等级名称
 * @params input 输入值
 * @return 是否输入不合法
 */
export function customedSecuName(input: any): boolean {
    return /[\*\:\/\?\"\<\>\|]/.test(input);
}

/**
 * 系统保护等级文字
 */
export const ProtLevelText = {
    [ProtLevel.Classified]: __('秘密级'),
    [ProtLevel.Confidential]: __('机密级一般'),
    [ProtLevel.MoreConfidential]: __('机密级增强'),
}

/**
 * 系统保护等级列表
 */
export const getProtLevels = (protLevel: ProtLevel): ReadonlyArray<ProtLevel> => {
    return [ProtLevel.Classified, ProtLevel.Confidential, ProtLevel.MoreConfidential].filter((item) => item >= protLevel)
}

/**
 * 用户密级分类
 */
export enum UserCsfLevelType {
    /**
     * 用户密级
     */
    UserLevel = 'csf_level_enum',

    /**
     * 用户密级2
     */
    UserLevel2 = 'csf_level2_enum',
}
