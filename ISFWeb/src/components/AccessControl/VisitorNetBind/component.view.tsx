import * as React from 'react';
import { CheckBox, Radio, Button } from '@/sweet-ui'
import ToastProvider from '@/ui/ToastProvider/ui.desktop';
import NetList from './NetList/component.view';
import VisitorList from './VisitorList/component.view';
import VisitorNetBindBase from './component.base';
import styles from './styles.view.css';
import __ from './locale';

export default class VisitorNetBind extends VisitorNetBindBase {

    render() {
        const { networkRestrictionIsEnabled, noNetworkPolicyAccessor, hasChanged, selectedNet, visitorListIsEnabled, isServerNormal } = this.state;
        return (
            <div className={styles['container']}>
                <div className={styles['header']}>
                    <div>
                        <CheckBox
                            className={styles['open-checkbox']}
                            checked={networkRestrictionIsEnabled}
                            onClick={(event) => event.stopPropagation()}
                            onCheckedChange={({ detail }) => this.toggleChecked(detail)}
                        >
                            <span className={styles['checkbox']} >{__('启用访问者网段限制')}</span>
                        </CheckBox>
                    </div>
                    {
                        networkRestrictionIsEnabled ? (
                            <div className={styles['limit-options']}>
                                <Radio
                                    checked={!noNetworkPolicyAccessor}
                                    onChange={this.changeNetworkRestriction}
                                >
                                    {__('开启后，被绑定的访问者只能在指定的网段内登录客户端，未绑定的访问者无法登录客户端。')}
                                </Radio>

                                <Radio
                                    checked={noNetworkPolicyAccessor}
                                    onChange={this.changeNoNetworkPolicy}
                                >
                                    {__('开启后，被绑定的访问者只能在指定的网段内登录客户端，未绑定的访问者登录客户端不受限制。')}
                                </Radio>
                            </div>
                        ) : null
                    }
                    {
                        hasChanged ? (
                            <div className={styles['config-buttons']}>
                                <Button
                                    theme="oem"
                                    className={styles['button']}
                                    onClick={this.confirm}
                                >
                                    {__('保存')}
                                </Button>
                                <Button
                                    onClick={this.cancel}
                                >
                                    {__('取消')}
                                </Button>
                            </div>
                        ) : null
                    }
                </div>
                <div className={styles['list-container']}>
                    <div className={styles['net-list']}>
                        <NetList
                            isEnabled={this.previousEnabledState.networkRestrictionIsEnabled}
                            serverStatus={isServerNormal}
                            onDisableVisitorList={(status) => this.handleDisableVisitorList(status)}
                            onSelectNet={(netInfo) => this.handleNetSelect(netInfo)}
                        />
                    </div>
                    <ToastProvider className={styles['visitor-list']}>
                        <VisitorList
                            isEnabled={this.previousEnabledState.networkRestrictionIsEnabled && visitorListIsEnabled}
                            selectedNet={selectedNet}
                        />
                    </ToastProvider>
                </div>
            </div>
        )
    }
}