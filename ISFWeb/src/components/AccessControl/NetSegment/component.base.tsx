import * as React from 'react';
import { every, includes } from 'lodash';
import { dispalyName, IP, IPV6, IPV6Prefix, subNetMask as subNetMaskValidate } from '@/util/validators';
import WebComponent from '../../webcomponent';
import {
    OperateType,
    NetInfo,
    IpVersion,
    ValidateState,
    InputKey,
    ValidateStateList,
} from './helper';

/**
 * Props
 */
export interface NetSegmentProps {
    /**
     * 操作类型
     */
    operateType: OperateType;
    /**
     * 网段信息
     */
    netInfo: NetInfo;

    /**
     * 是否显示弹窗标题
     */
    isShowTitle: boolean;

    /**
     * 是否显示网段名称
     */
    isShowNetName: boolean;

    /**
     * 格式化输出数据
     */
    convererOut: (net: NetInfo) => any;

    /**
     * 取消事件
     */
    onRequestCancel: () => void;

    /**
     * 确定事件
     */
    onRequestConfirm: (net: NetInfo) => void;
}

/**
 * State
 */
export interface NetSegmentState {
    /**
     * ip地址信息
     */
    netInfo: NetInfo;

    /**
     * 网段字段合法性验证
     */
    validateState: ValidateStateList;

    /**
     * 起始IP和终止IP合法性
     */
    ipError: boolean;
}

const ValidateList = [
    InputKey.Name,
    InputKey.OriginIp,
    InputKey.EndIp,
    InputKey.Ip,
    InputKey.Mask,
    InputKey.Prefix,
]

export default class NetSegmentBase extends WebComponent<NetSegmentProps, NetSegmentState>{

    static defaultProps = {
        isShowNetName: false,
    }

    state = {
        netInfo: {
            ...this.props.netInfo,
            single: ((this.props.netInfo.ipVersion === IpVersion.Ipv4 && !!this.props.netInfo.mask) ||
                (this.props.netInfo.ipVersion === IpVersion.Ipv6 && !!this.props.netInfo.prefix)),
        },
        ipError: false,
        validateState: {
            name: ValidateState.OK,
            originIP: ValidateState.OK,
            endIP: ValidateState.OK,
            ip: ValidateState.OK,
            mask: ValidateState.OK,
            prefix: ValidateState.OK,
        },
    }

    /**
     * 数据暂存
     */
    private temporary = {
        originIP: '',
        endIP: '',
        ip: '',
    }

    /**
     * 上次编辑状态
     */
    private lastEdit = {
        [IpVersion.Ipv4]: false,
        [IpVersion.Ipv6]: false,
    }

    componentDidMount() {
        const { netInfo: { ipVersion, single } } = this.state

        this.lastEdit = {
            ...this.lastEdit,
            [ipVersion]: single,
        }
    }

    /**
     * 选择IP地址类型 ipv4 | ipv6
     */
    protected ipVersionSelect = (newIpVersion: IpVersion): void => {
        const { netInfo: { ipVersion, originIP, endIP, ip } } = this.state

        if (ipVersion !== newIpVersion) {
            this.setState({
                netInfo: {
                    ...this.state.netInfo,
                    ipVersion: newIpVersion,
                    originIP: this.temporary.originIP,
                    endIP: this.temporary.endIP,
                    ip: this.temporary.ip,
                    single: this.lastEdit[newIpVersion],
                },
                ipError: false,
            }, () => {
                this.temporary = {
                    originIP,
                    endIP,
                    ip,
                }
                this.updateValidateState()
            })
        }
    }

    /**
     *  选择IP地址范围
     */
    protected ipRangeSelect = (single: boolean): void => {
        this.setState({
            netInfo: {
                ...this.state.netInfo,
                single,
            },
            ipError: false,
        })
    }

    /**
     * 设置IP信息
     */
    protected setIpInfo = (netInfoKey: InputKey, netInfoValue: string): void => {
        const { netInfo: { ipVersion } } = this.state
        const editState = includes([InputKey.Ip, InputKey.Mask, InputKey.Prefix], netInfoKey)

        this.setState({
            netInfo: {
                ...this.state.netInfo,
                [netInfoKey]: netInfoValue,
            },
            ipError: false,
        }, () => {
            this.lastEdit[ipVersion] = editState
            this.updateValidateState({ [netInfoKey]: ValidateState.OK })
        })
    }

    /**
     * 检查表单合法性
     */
    protected checkFrom = (validateList: ReadonlyArray<InputKey> = ValidateList): void => {
        const { isShowNetName } = this.props
        const { netInfo } = this.state
        const { name, mask, prefix, single, ipVersion } = this.state.netInfo
        let validateState = {}

        validateList.map((netInfoItem: InputKey) => {
            validateState = {
                ...validateState,
                [netInfoItem]: (
                    () => {
                        switch (netInfoItem) {
                            case InputKey.Name:
                                return isShowNetName && name && (!dispalyName(name) || name.length > 128)
                                    ? ValidateState.InvalidName
                                    : ValidateState.OK

                            case InputKey.OriginIp:
                            case InputKey.EndIp:
                            case InputKey.Ip: {
                                const netInfoValue = netInfo[netInfoItem]
                                const Invalid = ipVersion === IpVersion.Ipv4
                                    ? ValidateState.InvalidIPv4
                                    : ValidateState.InvalidIPv6
                                const verifyIP = ipVersion === IpVersion.Ipv4
                                    ? IP
                                    : IPV6

                                if ((netInfoItem !== InputKey.Ip && !single) || (netInfoItem === InputKey.Ip && single)) {
                                    return !netInfoValue
                                        ? ValidateState.Empty
                                        : !verifyIP(netInfoValue)
                                            ? Invalid
                                            : ValidateState.OK
                                }
                                return ValidateState.OK
                            }

                            case InputKey.Mask:
                                if (single) {

                                    return ipVersion === IpVersion.Ipv4
                                        ? !mask
                                            ? ValidateState.Empty
                                            : IP(mask) && subNetMaskValidate(mask)
                                                ? ValidateState.OK
                                                : IP(mask)
                                                    ? ValidateState.InvalidMask
                                                    : ValidateState.InvalidMaskForm
                                        : ValidateState.OK
                                }
                                return ValidateState.OK

                            case InputKey.Prefix:
                                if (single) {
                                    return ipVersion === IpVersion.Ipv6
                                        ? !prefix
                                            ? ValidateState.Empty
                                            : IPV6Prefix(prefix)
                                                ? ValidateState.OK
                                                : ValidateState.InvalidPrefix
                                        : ValidateState.OK
                                }
                                return ValidateState.OK
                        }
                    }
                )(),
            }
        })
        this.updateValidateState(validateState)
    }

    /**
     * 确定
     */
    protected confirm = async (): Promise<void> => {
        // 输入框检查
        await this.checkFrom()

        const {
            netInfo,
            netInfo: {
                originIP,
                endIP,
                ipVersion,
            },
            validateState,
            ipError,
        } = this.state

        if (
            every(validateState, (item) => item === ValidateState.OK && !ipError) &&
            this.validateRange(originIP.toUpperCase(), endIP.toUpperCase(), ipVersion)
        ) {
            this.props.onRequestConfirm(this.props.convererOut(netInfo))
        }
    }

    /**
     * 更新 | 清除 表单验证结果
     */
    private updateValidateState = (validateStateList?: ValidateStateList): void => {
        this.setState({
            validateState: {
                ...this.state.validateState,
                ...validateStateList
                    ? validateStateList
                    : {
                        name: ValidateState.OK,
                        originIP: ValidateState.OK,
                        endIP: ValidateState.OK,
                        ip: ValidateState.OK,
                        mask: ValidateState.OK,
                        prefix: ValidateState.OK,
                    },
            },
        })
    }

    /**
     * 校验IP和子网掩码的输入合法性
     * @param input 输入的值
     */
    protected validateInput = (input: string, isIpv4: boolean): boolean => {

        if (isIpv4) {
            const inputChar = input.charAt(input.length - 1);
            const inputArr = input.split('.');

            if (inputChar !== '.') {
                if (inputArr.length === 1) {
                    if (!Number.isNaN(parseInt(inputChar)) && /^[1-9]\d{0,2}/.test(inputArr[inputArr.length - 1])) {
                        return true;
                    }
                    return false;

                } else {
                    if (!Number.isNaN(parseInt(inputChar)) && /^0(?!\d)|(^[1-9]\d{0,2})/.test(inputArr[inputArr.length - 1])) {
                        return true;
                    }
                    return false;

                }
            }
            return true;

        } else {
            if (/[0-9a-fA-F]|\:|\./.test(input[input.length - 1])) {
                return true
            }
            return false
        }
    }

    /**
     * 验证起始IP和终止IP范围是否合法
    */
    private validateRange(originIP: string, endIP: string, ipVersion: IpVersion): boolean {
        const { single } = this.state.netInfo

        if (single) {
            return true
        }

        const isIPV4 = ipVersion === IpVersion.Ipv4
        const separator = isIPV4 ? '.' : ':'
        const originIPArr = originIP.split(separator);
        const endIPArr = endIP.split(separator);

        if (isIPV4) {
            for (let i = 0; i < originIPArr.length; i++) {
                if (parseInt(originIPArr[i]) < parseInt(endIPArr[i])) {

                    return true;
                } else if (parseInt(originIPArr[i]) > parseInt(endIPArr[i])) {
                    this.setState({
                        ipError: true,
                    })

                    return false;
                }
            }
            return true;
        } else {
            const fillOriginIPArr = this.fillIPV6(originIPArr)
            const fillEndIPArr = this.fillIPV6(endIPArr)

            for (let i = 0; i < fillOriginIPArr.length; i++) {
                if (fillOriginIPArr[i] < fillEndIPArr[i]) {

                    return true;
                } else if (fillOriginIPArr[i] > fillEndIPArr[i]) {
                    this.setState({
                        ipError: true,
                    })

                    return false;
                }
            }
        }
        return true;
    }

    /**
     * IPV6 位数补齐
     */
    private fillIPV6 = (ipv6Arr: ReadonlyArray<string>): ReadonlyArray<string> => {
        const len = 8 - ipv6Arr.length
        const index = ipv6Arr.indexOf('')
        let newIpv6Arr = []
        let before = ipv6Arr.slice(0, index)
        let after = ipv6Arr.slice(index + 1, ipv6Arr.length)

        for (let i = 0; i < len + 1; i++) {
            newIpv6Arr = newIpv6Arr.concat('')
        }
        ipv6Arr = before.concat(newIpv6Arr).concat(after)

        return ipv6Arr.map((item) => {
            if (!item) {
                return '0000';
            }

            const len = 4 - item.length;
            for (let i = 0; i < len; i++) {
                item = '0' + item;
            }

            return item;
        });
    }
}