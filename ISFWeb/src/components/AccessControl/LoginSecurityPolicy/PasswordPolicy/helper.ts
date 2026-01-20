import { Message } from '@/sweet-ui';
import __ from './locale';
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
 * 密码强度
 */
export enum PwdStrength {
    /**
     * 弱密码
     */
    WeakPwd,

    /**
     * 强密码
     */
    StrongPwd,

    /**
     * 初始值
     */
    None,

}

/**
 * 输入不合法状态
 */
export enum ValidateState {
    /**
     * 正常
     */
    Normal,

    /**
     * 空值
     */
    Empty,

    /**
     * 不合法的强密码字符长度最小值
     */
    InvalidStrongPwdLengthMin,

    /**
    * 不合法的强密码字符长度最大值
    */
    InvalidStrongPwdLengthMax,

    /**
     * 非涉密模式不合法的密码输错次数
     */
    InvalidPasswdErrCnt,

    /**
     * 不合法的密码锁定时间
     */
    InvalidPasswdLockTime,

    /**
     * 涉密模式不合法的密码输错次数
     */
    InvalidSecretPasswdErrCnt,

    /**
     * 初始化密码错误
     */
    InvalidInitPwd,
}

export function rederErrorMsg(error) {
    if (error.message) {
        Message.alert({ message: error.message });
    } else if (error.error) {
        Message.alert({ message: error.error.errMsg });
    }
}

/**
 * 第三方消息插件中的消息类型
 */
export enum MessageTypes {
    /**
     * 短信验证码
     */
    ResetPWDVerificationCode = 'authentication/v1/reset-pwd-verification-code',
}

/**
 * 第三方消息插件中的插件参数配置
 */
export const pluginParamConf = {
    // 短信验证码
    [MessageTypes.ResetPWDVerificationCode]: {
        key: 'ThirdPartyType',
        value: 'SMS',
    },
}

/**
 * 隐藏密码
 */
export const hiddenPwd = '**********'