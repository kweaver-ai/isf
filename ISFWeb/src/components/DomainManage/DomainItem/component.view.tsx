import * as React from 'react'
import { Drawer, Button } from '@/sweet-ui'
import SetUpKeywordSync from '../../SetUpKeywordSync/component.view'
import SpareDomain from '../../SpareDomain/component.view'
import DomainExpandItem from '../DomainExpandItem/component.view'
import { ActionType, DomainType } from '../helper'
import DomainInfo from './DomainInfo/component.view'
import SyncSetInfo from './SyncSetInfo/component.view'
import DomainItemBase, { ExpandItem } from './component.base'
import styles from './styles.view';
import __ from './locale'

export default class DomainItem extends DomainItemBase {
    render() {
        const { selection, editDomain, actionType } = this.props;
        const {
            restShow,
            expandItem: {
                infoExpand, syncExpand, keyExpand, spareExpand,
            },
            isDomainInfoEditStatus,
            isSyncSettingEditStatus,
            isKeywordSettingEditStatus,
            isSpareDomainEditStatus,
        } = this.state;
        const { name, id, domainName, password, useSSL, ipAddress, domainPort } = this.domainInfo;

        return (
            <Drawer
                role={'sweetui-drawer'}
                open={true}
                maskClosable={true}
                canOutsideClickClose={false}
                destroyOnClose={true}
                onDrawerClose={this.cancel}
                size={'50%'}
                title={actionType === ActionType.Add && this.domainInfo.id === -1 ? __('新建') : __('编辑')}
                footer={
                    actionType === ActionType.Add && selection.id || actionType === ActionType.Edit && editDomain.type !== DomainType.Primary ?
                        <div className={styles['drawer']}>
                            <Button
                                role={'sweetui-button'}
                                className={styles['button-left']}
                                theme={'oem'}
                                onClick={this.saveDomainInfo}
                            >
                                {__('确定')}
                            </Button>
                            <Button
                                role={'sweetui-button'}
                                className={styles['button-left']}
                                onClick={this.cancel}
                            >
                                {__('取消')}
                            </Button>
                        </div> : null
                }
            >
                <DomainExpandItem
                    title={__('域控信息')}
                    disabled={isSyncSettingEditStatus || isKeywordSettingEditStatus || isSpareDomainEditStatus}
                    isExpand={infoExpand}
                    onExpandItem={(isExpand) => this.expandItem(ExpandItem.DomainInfo, isExpand)}
                    showIcon={actionType === ActionType.Add ? id !== -1 : editDomain.type === DomainType.Primary}
                >
                    <DomainInfo
                        onRef={this.onRef}
                        selection={selection}
                        editDomain={editDomain}
                        actionType={actionType}
                        onSetDomainInfoSuccess={this.getDoaminInfo}
                        onRequestEditStatus={this.changeDomainInfoEditStatus}
                        onRequestDomainInvalid={this.handleDomainError}
                    />
                </DomainExpandItem>
                {
                    restShow ?
                        <div>
                            <DomainExpandItem
                                title={__('定期同步设置')}
                                disabled={isDomainInfoEditStatus || isKeywordSettingEditStatus || isSpareDomainEditStatus}
                                isExpand={syncExpand}
                                onExpandItem={(isExpand) => this.expandItem(ExpandItem.SyncSetting, isExpand)}
                            >
                                {
                                    syncExpand ?
                                        <SyncSetInfo
                                            domainInfo={actionType === ActionType.Add ? this.domainInfo : editDomain}
                                            selection={selection}
                                            actionType={actionType}
                                            onSetSyncSetInfoSuccess={this.getSynSetInfo}
                                            onRequestEditStatus={this.changeSyncInfoEditStatus}
                                            onRequestDomainInvalid={this.handleDomainError}
                                        /> : null
                                }
                            </DomainExpandItem>
                            <DomainExpandItem
                                title={__('同步关键字设置')}
                                disabled={isDomainInfoEditStatus || isSyncSettingEditStatus || isSpareDomainEditStatus}
                                isExpand={keyExpand}
                                onExpandItem={(isExpand) => this.expandItem(ExpandItem.KeySetting, isExpand)}
                            >
                                {
                                    keyExpand ?
                                        <SetUpKeywordSync
                                            domainInfo={actionType === ActionType.Add ? this.domainInfo : editDomain}
                                            onRequestEditStatus={this.changeKeySettingEditStatus}
                                            onRequestDomainInvalid={this.handleDomainError}
                                        /> : null
                                }
                            </DomainExpandItem>
                            <DomainExpandItem
                                title={__('备用域控信息')}
                                disabled={isDomainInfoEditStatus || isSyncSettingEditStatus || isKeywordSettingEditStatus}
                                isExpand={spareExpand}
                                onExpandItem={(isExpand) => this.expandItem(ExpandItem.SpareDomain, isExpand)}
                            >
                                {
                                    spareExpand ? <SpareDomain
                                        mainDomain={actionType === ActionType.Add ? { name, adminName: domainName, password, useSSL, ipAddress, id, port: domainPort } : editDomain}
                                        onRequestEditStatus={this.changeSpareDomainEditStatus}
                                        onRequestDomainInvalid={this.handleDomainError}
                                    /> : null
                                }
                            </DomainExpandItem>
                        </div> : null
                }
            </Drawer >
        )
    }
}