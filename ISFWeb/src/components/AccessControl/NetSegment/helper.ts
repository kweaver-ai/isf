import __ from './locale';

/**
 * ip版本
 */
export enum IpVersion {
    /**
     * ipv4
     */
    Ipv4 = 'ipv4',

    /**
     * ipv6
     */
    Ipv6 = 'ipv6',
}

/**
 * 操作类型
 */
export enum OperateType {
    /**
     * 添加
     */
    Add = 'add',
    /**
     * 编辑
     */
    Edit = 'edit',
}

/**
 * 网段信息
 */
export interface NetInfo {
    /**
     * 网段名称
     */
    name?: string;

    /**
     * ip版本
     */
    ipVersion: IpVersion;

    /**
     * 起始IP
     */
    originIP: string;

    /**
     * 终止IP
     */
    endIP: string;

    /**
     * ip地址
     */
    ip: string;

    /**
     * 子网掩码
     */
    mask: string;

    /**
     * 前缀
     */
    prefix: string;

    /**
     * ip范围
     */
    single: boolean;

    /**
     * ID
     */
    id?: string;
}

/**
 * 表单验证
 */
export interface ValidateStateList {
    /**
     * 名称合法性状态
    */
    name?: ValidateState;

    /**
     * 起始IP合法性状态
     */
    originIP?: ValidateState;

    /**
     * 终止IP合法性状态
     */
    endIP?: ValidateState;

    /**
     * ip地址合法性状态
     */
    ip?: ValidateState;

    /**
     * 子网掩码合法性状态
     */
    mask?: ValidateState;

    /**
     * 子网掩码合法性状态
     */
    prefix?: ValidateState;
}

/**
 * 输入框标识
 */
export enum InputKey {

    /**
     * 网段名输入框
     */
    Name = 'name',

    /**
     * 起始IP输入框
     */
    OriginIp = 'originIP',

    /**
     * 终止IP输入框
     */
    EndIp = 'endIP',

    /**
     * IP地址输入框
     */
    Ip = 'ip',

    /**
     * 子网掩码输入框
     */
    Mask = 'mask',

    /**
     * 前缀长度输入框
     */
    Prefix = 'prefix',
}

/**
 * 输入框默认提示
 */
export const Placeholder = {
    /**
     * 网段名称
     */
    name: __('可选填'),

    /**
     * ipv4
     */
    ipv4: __('例如：') + ' 192.168.10.10',

    /**
     * ipv6
     */
    ipv6: __('例如：') + ' 2002:50::44',

    /**
     * 子网掩码
     */
    mask: __('例如：') + ' 255.255.255.0',

    /**
     * 前缀
     */
    prefix: __('例如：') + ' 64',
};

/**
 * 输入不合法状态
 */
export enum ValidateState {
    /**
     * 正常
     */
    OK,

    /**
     * 空值
     */
    Empty,

    /**
     * 网段输入错误
     */
    InvalidName,

    /**
     * ipv4输入错误
     */
    InvalidIPv4,

    /**
     * ipv6输入错误
     */
    InvalidIPv6,

    /**
     * 子网掩码输入错误
     */
    InvalidMaskForm,

    /**
     * 子网掩码不合法
     */
    InvalidMask,

    /**
     * 前缀长度不合法
     */
    InvalidPrefix,

    /**
     * 终止IP小于起始IP
     */
    InvalidRange
}

/**
 * 不合法提示
 */
export const ValidateMessages = {
    [ValidateState.Empty]: __('此输入项不允许为空。'),
    [ValidateState.InvalidName]: __('网段名称不能包含 \\ / : * ? " < > | 特殊字符，长度不能超过128个字符。'),
    [ValidateState.InvalidIPv4]: __('IPv4地址格式形如 XXX.XXX.XXX.XXX，每段必须是 0~255 之间的整数。'),
    [ValidateState.InvalidIPv6]: __('IPv6地址格式形如 XXXX:XXXX:XXXX:XXXX:XXXX:XXXX:XXXX:XXXX，其中每个X都为十六进制数。'),
    [ValidateState.InvalidMaskForm]: __('子网掩码格式形如 XXX.XXX.XXX.XXX，每段必须是0~255之间的整数。'),
    [ValidateState.InvalidMask]: __('非法的网段掩码参数。'),
    [ValidateState.InvalidPrefix]: __('前缀长度必须是0~128之间的整数。'),
}