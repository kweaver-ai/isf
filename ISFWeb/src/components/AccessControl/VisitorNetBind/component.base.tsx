import * as React from 'react';
import { isEqual } from 'lodash';
import { Toast } from '@/sweet-ui';
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import { getState, setState } from '@/core/apis/console/networkRestriction'
import WebComponent from '../../webcomponent';
import { getNetBindErrorMessage, NetInfo } from './helper'
import __ from './locale';

interface VisitorNetBindState {
    /**
     * 是否开启访问者网段策略
     */
    networkRestrictionIsEnabled: boolean;

    /**
     * 是否开启未配置网段策略访问者可以访问
     */
    noNetworkPolicyAccessor: boolean;

    /**
     * 访问者网段限制状态是否变更
     */
    hasChanged: boolean;

    /**
    * 访问者列表是否可用
    */
    visitorListIsEnabled: boolean;

    /**
     * 选中的网段
     */
    selectedNet: NetInfo;

    /**
     * 服务是否正常
     */
    isServerNormal: boolean;
}

export default class VisitorNetBindBase extends WebComponent<any, VisitorNetBindState> {

    state: VisitorNetBindState = {
        networkRestrictionIsEnabled: false,
        noNetworkPolicyAccessor: false,
        hasChanged: false,
        visitorListIsEnabled: false,
        selectedNet: null,
        isServerNormal: false,
    }

    // 网段配置初始值
    previousEnabledState =  {
        networkRestrictionIsEnabled: false,
        noNetworkPolicyAccessor: false,
    }

    async componentDidMount() {
        try {
            const policy = (await getState({ name: 'network_restriction,no_network_policy_accessor' })).data;

            if (policy && policy.length !== 0) {
                const { value: { is_enabled } } = policy.find(({ name }) => name === 'network_restriction')
                const { value: { is_enabled: accessorIsEnabled } } = policy.find(({ name }) => name === 'no_network_policy_accessor')

                this.setState({
                    networkRestrictionIsEnabled: is_enabled,
                    noNetworkPolicyAccessor: accessorIsEnabled,
                    visitorListIsEnabled: is_enabled,
                    isServerNormal: true,
                })

                this.previousEnabledState = {
                    networkRestrictionIsEnabled: is_enabled,
                    noNetworkPolicyAccessor: accessorIsEnabled,
                }
            } else {
                Toast.open(__('该策略不存在'))
            }
        } catch (error) {
            if (error.status === 0) {
                this.setState({
                    isServerNormal: false,
                })
            } else {
                getNetBindErrorMessage(error)
            }
        }
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
     * 打开/关闭访问者网段绑定策略
     */
    protected async toggleChecked(value) {
        this.setState({
            networkRestrictionIsEnabled: value,
            hasChanged: !isEqual(value, this.previousEnabledState.networkRestrictionIsEnabled),
            noNetworkPolicyAccessor: value
                ? this.previousEnabledState.noNetworkPolicyAccessor
                : this.state.noNetworkPolicyAccessor,
        })
    }

    /**
     * 改变网段限制选中状态
     */
    protected changeNetworkRestriction = () => {
        this.setState({
            noNetworkPolicyAccessor: false,
            hasChanged: this.previousEnabledState.noNetworkPolicyAccessor
            || this.previousEnabledState.networkRestrictionIsEnabled !== this.state.networkRestrictionIsEnabled,
        })
    }

    /**
     * 改变网段限制选中状态
     */
    protected changeNoNetworkPolicy = () => {
        this.setState({
            noNetworkPolicyAccessor: true,
            hasChanged: !this.previousEnabledState.noNetworkPolicyAccessor
                || this.previousEnabledState.networkRestrictionIsEnabled !== this.state.networkRestrictionIsEnabled,
        })
    }

    /**
     * 保存按钮
     */
    protected confirm = async () => {
        try {
            const { networkRestrictionIsEnabled, noNetworkPolicyAccessor } = this.state;
            const payload = [
                { name: 'network_restriction', value: { is_enabled: networkRestrictionIsEnabled } },
                { name: 'no_network_policy_accessor', value: { is_enabled: noNetworkPolicyAccessor } },
            ]

            await setState({ name: 'network_restriction,no_network_policy_accessor', payload });

            this.previousEnabledState = {
                networkRestrictionIsEnabled,
                noNetworkPolicyAccessor,
            }

            this.setState({
                hasChanged: false,
                isServerNormal: true,
            })

            manageLog(
                ManagementOps.SET,
                networkRestrictionIsEnabled ? __('启用 访问者网段限制 成功') : __('关闭 访问者网段限制 成功'),
                networkRestrictionIsEnabled ?
                    __('开启后，被绑定的访问者只能在指定的网段内登录客户端，未绑定的访问者${canLogin}', { canLogin: noNetworkPolicyAccessor ?  __('登录客户端不受限制') : __('无法登录客户端') })
                    : null,
                Level.WARN,
            )
        } catch (error) {
            getNetBindErrorMessage(error);
        }
    }

    /**
     * 取消按钮
     */
    protected cancel = (): void => {
        this.setState({
            networkRestrictionIsEnabled: this.previousEnabledState.networkRestrictionIsEnabled,
            noNetworkPolicyAccessor: this.previousEnabledState.noNetworkPolicyAccessor,
            hasChanged: false,
        })
    }

    /**
     * 选中某一个网段
     * @param netInfo 选择的网段
     */
    protected handleNetSelect(netInfo: NetInfo) {
        this.setState({
            selectedNet: netInfo,
        })

        if (!netInfo) {
            this.setState({
                visitorListIsEnabled: false,
            })
        }
    }

    /**
     * 禁用访问者列表
     */
    protected handleDisableVisitorList(status: boolean) {
        this.setState({
            visitorListIsEnabled: !status,
        })
    }

}