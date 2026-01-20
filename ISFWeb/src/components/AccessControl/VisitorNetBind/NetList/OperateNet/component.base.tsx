import * as React from 'react';
import { dispalyName, IP, subNetMask as subNetMaskValidate } from '@/util/validators';
import WebComponent from '../../../../webcomponent';
import { NetType, ValidateState, NetInfo } from '../../helper';

interface OperateNetProps extends React.Props<void> {
    /**
     * 编辑或添加的网段信息
     */
    net: NetInfo;

    /**
     * 编辑成功
     */
    onRequestEditSuccess: (data: NetInfo) => void;

    /**
     * 取消编辑或添加
     */
    onRequestEditCancel: () => void;
}
interface OperateNetState {
    /**
     * 网段信息
     */
    net: NetInfo;

    /**
     * 网段字段合法性状态
     */
    validateState: {
        /**
         * 名称合法性状态
         */
        name: ValidateState;

        /**
         * 起始IP合法性状态
         */
        originIP: ValidateState;

        /**
         * 终止IP合法性状态
         */
        endIP: ValidateState;

        /**
         * ip地址合法性状态
         */
        ip: ValidateState;

        /**
         * 子网掩码合法性状态
         */
        mask: ValidateState;
    };

    /**
     * 终止IP小于起始IP时错误信息
     */
    error: ValidateState;
}

export default class OperateNetBase extends WebComponent<OperateNetProps, OperateNetState> {

    state: OperateNetState = {
        net: {
            ...this.props.net,
            originIP: this.props.net.originIP ? this.props.net.originIP : '',
            endIP: this.props.net.endIP ? this.props.net.endIP : '',
            ip: this.props.net.ip ? this.props.net.ip : '',
            mask: this.props.net.mask ? this.props.net.mask : '',
        },
        validateState: {
            name: ValidateState.OK,
            originIP: ValidateState.OK,
            endIP: ValidateState.OK,
            ip: ValidateState.OK,
            mask: ValidateState.OK,
        },
        error: ValidateState.OK,
    }

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return;
        };
    }

    /**
     * 校验IP和子网掩码的输入合法性
     * @param input 输入的值
     */
    protected validateInput(input) {
        let inputChar = input.charAt(input.length - 1);
        let inputArr = input.split('.');
        if (inputChar !== '.') {
            if (inputArr.length === 1) {
                if (!Number.isNaN(parseInt(inputChar)) && /^[1-9]\d{0,2}/.test(inputArr[inputArr.length - 1])) {
                    return true;
                } else {
                    return false;
                }
            } else {
                if (!Number.isNaN(parseInt(inputChar)) && /^0(?!\d)|(^[1-9]\d{0,2})/.test(inputArr[inputArr.length - 1])) {
                    return true;
                } else {
                    return false;
                }
            }
        } else {
            return true;
        }
    }

    /**
     * 设置state中表单的值
     * @param key 对应的输入框键值
     * @param value 填写的字段值
     */
    protected setValue(key: string, value: string) {
        const { validateState, net } = this.state;
        this.setState({
            net: {
                ...net,
                [key]: value,
            },
            validateState: {
                name: key === 'name' ? ValidateState.OK : validateState.name,
                originIP: key === 'originIP' ? ValidateState.OK : validateState.originIP,
                endIP: key === 'endIP' ? ValidateState.OK : validateState.endIP,
                ip: key === 'ip' ? ValidateState.OK : validateState.ip,
                mask: key === 'mask' ? ValidateState.OK : validateState.mask,
            },
            error: ValidateState.OK,
        })
    }

    /**
     * 验证输入框输入的合法性
     * @return  boolean 验证的结果
     */
    protected validateNetInfo(): boolean {
        const { net: { name, originIP, endIP, ip, mask, netType } } = this.state;
        const nameValidity = !name || (name && dispalyName(name) && name.length <= 128);
        const originIPValidity = (originIP && IP(originIP)) || netType !== NetType.Range;
        const endIPValidity = (endIP && IP(endIP)) || netType !== NetType.Range;
        const ipValidity = (ip && IP(ip)) || netType !== NetType.Mask;
        const subNetMaskValidity = (mask && IP(mask) && subNetMaskValidate(mask)) || netType !== NetType.Mask;

        if (nameValidity && originIPValidity && endIPValidity && ipValidity && subNetMaskValidity) {

            if (netType === NetType.Mask) {
                return true;
            } else if (this.validateRange(originIP, endIP)) {
                return true;
            } else {
                return false;
            }

        } else {
            this.setState({
                validateState: {
                    name: nameValidity ? ValidateState.OK : ValidateState.InvalidName,
                    originIP: (netType !== NetType.Range || originIP) ? (originIPValidity ? ValidateState.OK : ValidateState.InvalidIP) : ValidateState.Empty,
                    endIP: (netType !== NetType.Range || endIP) ? (endIPValidity ? ValidateState.OK : ValidateState.InvalidIP) : ValidateState.Empty,
                    ip: (netType !== NetType.Mask || ip) ? (ipValidity ? ValidateState.OK : ValidateState.InvalidIP) : ValidateState.Empty,
                    mask: (netType !== NetType.Mask || mask) ? (subNetMaskValidity ? ValidateState.OK : (!IP(mask) ? ValidateState.InvalidMaskForm : ValidateState.InvalidMask)) : ValidateState.Empty,
                },
            })

            return false;
        }
    }

    /**
     * 验证起始IP和终止IP范围是否合法
     * @param originIP 起始IP
     * @param endIP 终止IP
     * @return  boolean 验证的结果
     */
    protected validateRange(originIP: string, endIP: string): boolean {
        const originIPArr = originIP.split('.');
        const endIPArr = endIP.split('.');

        for (let i = 0; i < originIPArr.length; i++) {
            if (parseInt(originIPArr[i]) < parseInt(endIPArr[i])) {
                return true;
            } else if (parseInt(originIPArr[i]) > parseInt(endIPArr[i])) {
                this.setState({
                    error: ValidateState.InvalidRange,
                })
                return false;
            }
        }

        return true;
    }

    /**
     * 选择的网段设置形式
     * @param moduleOption 选择的形式
     */
    protected changeNetType({ detail: { value } }) {
        const { net, validateState } = this.state;
        this.setState({
            net: {
                ...net,
                netType: value,
            },
            error: ValidateState.OK,
        })

        if (value === NetType.Range) {
            this.setState({
                validateState: {
                    ...validateState,
                    ip: ValidateState.OK,
                    mask: ValidateState.OK,
                },
            })
        } else {
            this.setState({
                validateState: {
                    ...validateState,
                    originIP: ValidateState.OK,
                    endIP: ValidateState.OK,
                },
            })
        }

    }

    /**
     * 保存网段信息
     */
    protected saveNet() {
        if (this.validateNetInfo()) {
            this.props.onRequestEditSuccess(this.state.net);
        }
    }

}