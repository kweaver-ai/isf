import __ from './locale';

/**
 * 输入框状态验证结果
 */
export enum ValidateStatusEnum {
    /**
     * 正常
     */
    Normal,

    /**
     * 输入项为空
     */
    Empty,
}

/**
 * 输入框输入字段
 */
export interface KeywordInput {
    /**
     * 部门名对应的域key字段
     */
    departNameKeys: string;

    /**
     * 部门ID对应的域key字段
     */
    departThirdIdKeys: string;

    /**
     * 登录名对应的域key字段
     */
    loginNameKeys: string;

    /**
     * 显示名对应的域key字段
     */
    displayNameKeys: string;

    /**
     * 用户邮箱对应的域key字段
     */
    emailKeys: string;

    /**
     * 用户Id对应的域key字段
     */
    userThirdIdKeys: string;

    /**
     * 安全组信息的key字段
     */
    groupKeys: string;

    /**
     * 搜索子部门的Filter
     */
    subOuFilter: string;

    /**
     * 搜索子用户的Filter
     */
    subUserFilter: string;

    /**
     * 具体某个部门或用户信息的filter
     */
    baseFilter: string;

    /**
     * 用户状态对应的域key字段
     */
    statusKeys: string;

    /**
     * 用户身份证号对应的域key字段
     */
    idcardNumberKeys: string;

    /**
     * 用户电话号对应的域key字段
     */
    telNumberKeys: string;
}

/**
 * 验证输入框内容是否正确
 */
interface ValidateStatus {
    /**
     * 部门名对应的域key字段验证状态
     */
    departNameKeysValidateStatus: ValidateStatusEnum;

    /**
     * 部门ID对应的域key字段验证状态
     */
    departThirdIdKeysValidateStatus: ValidateStatusEnum;

    /**
     * 登录名对应的域key字段验证状态
     */
    loginNameKeysValidateStatus: ValidateStatusEnum;

    /**
     * 显示名对应的域key字段验证状态
     */
    displayNameKeysValidateStatus: ValidateStatusEnum;

    /**
     * 用户邮箱对应的域key字段验证状态
     */
    emailKeysValidateStatus: ValidateStatusEnum;

    /**
     * 用户Id对应的域key字段验证状态
     */
    userThirdIdKeysValidateStatus: ValidateStatusEnum;

    /**
     * 搜索子部门的Filter验证状态
     */
    subOuFilterValidateStatus: ValidateStatusEnum;

    /**
     * 搜索子用户的Filter验证状态
     */
    subUserFilterValidateStatus: ValidateStatusEnum;

    /**
     * 具体某个部门或用户信息的filter验证状态
     */
    baseFilterValidateStatus: ValidateStatusEnum;
}

/**
 * 验证后气泡提示信息
 */
export const ValidateMessage = {
    [ValidateStatusEnum.Empty]: __('此项不允许为空。'),
}

export interface SetUpKeywordSyncState {
    /**
     * 同步输入关键字状态
     */
    keywordInput: KeywordInput;

    /**
     * 输入框是否为编辑状态，用来决定是否展示保存/取消按钮以及父组件的下一步按钮是否灰化
     */
    isEditStatus: boolean;

    /**
     * 验证输入框内容是否正确
     */
    validateStatus: ValidateStatus;
}

/**
 * 父组件传入的域控制器相关信息
 */
interface DomainInfo {
    /**
     * 添加或者编辑域控时传入的id
     */
    id: number;

    /**
     * 主控域名
     */
    name: string;
}

export interface SetUpKeywordSyncProps {
    /**
     * 父组件传入的域控制器相关信息
     */
    domainInfo: DomainInfo;

    /**
     * 用来传递isEditStatus的回调函数，用来通知父组件下一步按钮是否灰化
     */
    onRequestEditStatus: (isEditStatus: boolean) => void;

    /**
     * 主域id失效时的回调
     */
    onRequestDomainInvalid: (name: string) => void;
}