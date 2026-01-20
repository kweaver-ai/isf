import { LoginAuthType, LoginWay } from '@/core/logincertification/logincertification';
import { Message } from '@/sweet-ui';
import __ from './locale';
/**
 * 自动禁用时间
 */
export enum DisableTime {

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
}

/**
 * 账号禁止登录客户端
 */
export enum OstypeInfo {

    /**
     * 默认值，所有端都允许登录
     */
    Default = 0,

    /**
     * 禁止iOS客户端登录
     */
    iOS = 1,

    /**
     * 禁止Android客户端登录
     */
    Android = 2,

    /**
     * 禁止wibndows客户端登录
     */
    Windows = 4,

    /**
     * 禁止Mac客户端登录
     */
    Mac = 5,

    /**
     * 禁止Web客户端登录
     */
    Web = 6,

    /**
     * 禁止移动客户端登录
     */
    Mobileweb = 7,

    /**
     * 禁止Linux客户端登录
     */
    Linux = 8,
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
     * 不合法的密码输错次数
     */
    InvalidPasswdErrCnt,

}

export const ValidateMessages = {
    [ValidateState.Empty]: __('此项不允许为空。'),
    [ValidateState.InvalidPasswdErrCnt]: __('密码错误次数范围为0~99。'),
}

export function rederErrorMsg(error) {
    if (error.message) {
        Message.alert({ message: error.message });
    } else if (error.error) {
        Message.alert({ message: error.error.errMsg });
    }
}

export const LoginWays = [
    LoginWay.Account,
    LoginWay.AccountAndImgCaptcha,
    LoginWay.AccountAndSmsCaptcha,
    LoginWay.AccountAndDynamicPassword,
]

export const LoginAuthTypeMapper = {
    [LoginAuthType.account]: LoginWay.Account,
    [LoginAuthType.accountAndImageCaptcha]: LoginWay.AccountAndImgCaptcha,
    [LoginAuthType.accountAndSMSCaptcha]: LoginWay.AccountAndSmsCaptcha,
    [LoginAuthType.accountAndDynamicPassword]: LoginWay.AccountAndDynamicPassword,
}