import { getErrorMessage } from '@/core/errcode';
import { ErrorCode, PublicErrorCode } from '@/core/apis/openapiconsole/errorcode'
import { Message, Toast } from '@/sweet-ui';
import { NodeType } from '@/core/organization';
import { IpVersion } from '../NetSegment/helper';
import __ from './locale';

/**
 * 网段信息
 */
export interface NetInfo {
    /**
     * 网段id
     */
    id: string;

    /**
     * 网段名称
     */
    name: string;

    /**
     * 网段类型
     */
    netType: NetType;

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
     * ip类型
     */
    ipVersion: IpVersion;
}

/**
 * 访问者
 */
export interface Visitor {
    /**
     * 访问者id
     */
    id: string;

    /**
     * 访问者名称
     */
    name: string;

    /**
     * 访问者类型
     */
    type: NodeType;
}

/**
 * 网段形式
 */
export enum NetType {

    /**
     * 选择起始IP和终止IP
     */
    Range = 'ip_segment',

    /**
     *选择IP地址和子网掩码
     */
    Mask = 'ip_mask',

}

/**
 * 访问者类型
 */
export enum UserType {

    /**
     * 用户
     */
    User = 'user',

    /**
     * 部门
     */
    Department = 'department',
}

/**
 * 选择类型
 */
export enum NodeSelectType {
    /**
     * 禁止选中
     */
    None = 0,

    /**
     * 同级单选
     */
    Single = 1,

    /**
     * 同级多选
     */
    Multiple = 2,

    /**
     * 级联单选
     */
    CascadeSingle = 3,

    /**
     * 级联多选
     */
    CascadeMultiple = 4,

    /**
     * 无限制
     */
    Unrestricted = 5,
}

/**
 * 每页显示数据条数
 */
export const PageSize = 20;

/**
 * 默认开始页
 */
export const DefaultPage = 1;

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
     * ip输入错误
     */
    InvalidIP,

    /**
     * 子网掩码输入错误
     */
    InvalidMaskForm,

    /**
     * 子网掩码不合法
     */
    InvalidMask,

    /**
     * 终止IP小于起始IP
     */
    InvalidRange,

}

/**
 * 输入框标识
 */
export enum NetInputKey {

    /**
     * 网段名输入框
     */
    Name = 'name',

    /**
     * 起始IP输入框
     */
    OriginIp = 'originIp',

    /**
     * 终止IP输入框
     */
    EndIp = 'endIp',

    /**
     * IP地址输入框
     */
    Ip = 'ip',

    /**
     * 子网掩码输入框
     */
    Mask = 'mask',
}

export const ValidateMessages = {
    [ValidateState.Empty]: __('此输入项不允许为空。'),
    [ValidateState.InvalidName]: __('网段名称不能包含 \\ / : * ? " < > | 特殊字符，长度不能超过128个字符。'),
    [ValidateState.InvalidIP]: __('IP地址格式形如xxx.xxx.xxx.xxx，每段必须是 0~255 之间的整数。'),
    [ValidateState.InvalidMaskForm]: __('子网掩码格式形如xxx.xxx.xxx.xxx，每段必须是 0~255 之间的整数。'),
    [ValidateState.InvalidMask]: __('非法的网段掩码参数。'),
}

/**
 * 错误信息
 */
export function getNetBindErrorMessage(error: any, net?: NetInfo) {
    if (error.status !== 0) {
        switch (error.code) {
            case ErrorCode.InvalidRequest:
            case PublicErrorCode.BadRequest:
                Message.alert({ message: __('请求参数错误。') });
                break;

            case ErrorCode.ResourceInaccessibleByPolicy:
            case PublicErrorCode.NotFound:
                error.detail.notfound_params.indexOf('network_restriction') >= 0 ?
                    Toast.open(__('该策略不存在'))
                    :
                    Message.alert({
                        message: error.detail.notfound_params.indexOf('id') >= 0 ?
                            __('该网段已不存在。') : __('请求出现错误。'),
                    })
                break;

            case ErrorCode.TooManyRequestsByPolicy:
                Message.alert({ message: __('请求过多。') });
                break;

            case ErrorCode.NoPermissionToOperateByPolicy:
            case PublicErrorCode.Forbidden:
                Message.alert({ message: __('不允许此操作。') });
                break;

            case ErrorCode.ResourceConflictByPolicy:
            case PublicErrorCode.Conflict:
                Message.alert({
                    message: error.detail.conflict_params.indexOf('name') >= 0 ?
                        __('该网段名称已存在，请重新输入。')
                        :
                        __('网段“${network}”已存在，无法重复添加。',
                            {
                                network: net.netType === NetType.Range ?
                                    net.originIP + '-' + net.endIP
                                    : net.ip + '(' + net.mask + ')',
                            },
                        ),
                })
                break;

            default: {
                const message = error.code ? getErrorMessage(error.code) : (error.message || '')

                if (message) {
                    Message.alert({ message });
                }
            }
        }
    }
}

/**
 * 从Location获取id
 */
export function getIdFromLocaltion(location) {
    return location.split('/').slice(-1)[0];
}
