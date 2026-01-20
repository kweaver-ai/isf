import * as React from 'react'
import { Button, DataGrid, Switch } from '@/sweet-ui'
import { Text, Title, InlineButton, UIIcon } from '@/ui/ui.desktop'
import { ListTipStatus } from '../ListTipComponent/helper'
import ListTipComponent from '../ListTipComponent/component.view'
import DomainItem from './DomainItem/component.view'
import { ActionType, DomainType } from './helper'
import DomainManageBase from './component.base'
import styles from './styles.view';
import __ from './locale'

export default class DomainManage extends DomainManageBase {

    render() {
        const { listTipStatus, domainList, actionType, selection } = this.state;

        return (
            <div className={styles['container']}>
                <div className={styles['header']}>
                    <Button
                        role={'sweetui-button'}
                        disabled={selection.type === DomainType.Trust}
                        width={'auto'}
                        icon={'add'}
                        size={'auto'}
                        theme={'oem'}
                        onClick={() => this.changeActionType(ActionType.Add)}
                    >
                        {__('新建')}
                    </Button>
                </div>

                <div className={styles['main']}>
                    <DataGrid
                        role={'sweetui-datagrid'}
                        height={'100%'}
                        enableSelect={true}
                        enableMultiSelect={false}
                        onSelectionChange={this.changeSelection}
                        EmptyComponent={true}
                        refreshing={listTipStatus !== ListTipStatus.None}
                        RefreshingComponent={
                            <ListTipComponent
                                listTipStatus={listTipStatus}
                            />
                        }
                        expandable={true}
                        data={domainList}
                        rowKeyName={'id'}
                        expandedKeys={this.state.expandedKeys}
                        onExpand={({ detail }) => this.expand(detail)}
                        columns={[
                            {
                                title: __('域名称'),
                                key: 'name',
                                width: '30%',
                                minWidth: 90,
                                renderCell: (value, record) => (
                                    <div className={styles['domain-name']}>
                                        <UIIcon
                                            role={'ui-uiicon'}
                                            code={record.type === DomainType.Primary ? '\uf119' : '\uf11a'}
                                            size={record.type === DomainType.Primary ? 24 : 18}
                                        />
                                        <div className={styles['text']}>
                                            <Text role={'ui-text'}>{record.name}</Text>
                                        </div>
                                    </div>
                                ),
                            },
                            {
                                title: __('类型'),
                                key: 'type',
                                width: '15%',
                                renderCell: (value, record) => (
                                    <Text role={'ui-text'}>{record.type === DomainType.Primary ? __('主域') : record.type === DomainType.Sub ? __('子域') : __('信任域')}</Text>
                                ),
                            },
                            {
                                title: __('同步周期'),
                                key: 'operation',
                                width: '20%',
                                renderCell: (value, record) => (
                                    <Text role={'ui-text'}>{record.type === DomainType.Primary ? record.config.syncInterval >= 1440 ? record.config.syncInterval / 24 / 60 + __('天') : record.config.syncInterval >= 60 ? record.config.syncInterval / 60 + __('小时') : record.config.syncInterval + __('分钟') : '---'}</Text>
                                ),
                            },
                            {
                                title: __('定期同步'),
                                key: 'interval',
                                width: '15%',
                                renderCell: (value, record) => (
                                    record.type === DomainType.Primary ?
                                        <Title
                                            role={'ui-title'}
                                            content={record.syncStatus === 0 ? __('关闭定期同步') : __('开启定期同步')}
                                        >
                                            <Switch
                                                role={'sweetui-switch'}
                                                checked={record.syncStatus === 0}
                                                disabled={!record.status}
                                                onChange={({ detail }) => this.changeSyncStatus(detail, record.id, record.name)}
                                            />
                                        </Title> : '---'
                                ),
                            },
                            {
                                title: __('操作'),
                                key: 'oprate',
                                width: '20%',
                                minWidth: 150,
                                renderCell: (value, record) => (
                                    <div>
                                        <Title
                                            role={'ui-title'}
                                            content={record.status ? __('禁用域控制器') : __('启用域控制器')}
                                        >
                                            <Switch
                                                role={'sweetui-switch'}
                                                checked={record.status}
                                                onChange={({ detail }) => this.changeDomainStatus(detail, record.id, record.name)}
                                            />
                                        </Title>
                                        <InlineButton
                                            role={'ui-inlinebutton'}
                                            className={styles['inline-button']}
                                            code={'\uf111'}
                                            title={__('立即同步')}
                                            disabled={!record.status || record.type !== DomainType.Primary || record.syncStatus !== 0}
                                            onClick={() => this.syncNow(record.id, record.status, record.name)}
                                        />
                                        <InlineButton
                                            role={'ui-inlinebutton'}
                                            className={styles['inline-button']}
                                            code={'\uf085'}
                                            title={__('编辑')}
                                            onClick={(e) => this.setDomain(e, record)}
                                        />
                                        <InlineButton
                                            role={'ui-inlinebutton'}
                                            className={styles['inline-button']}
                                            code={'\uf000'}
                                            title={__('删除')}
                                            onClick={() => this.deleteDomain(record.id, record.name)}
                                        />
                                    </div>
                                ),
                            },
                        ]}
                    />
                </div>
                {
                    actionType !== ActionType.None ?
                        <DomainItem
                            selection={selection}
                            editDomain={this.editDomain}
                            actionType={actionType}
                            onEditDomainInfoSuccess={(id) => this.setDomainInfoFinish(id)}
                            onRequestCancel={() => this.setState({ actionType: ActionType.None })}
                            onDomainError={this.handleDomainError}
                        /> : null
                }
            </div>
        )
    }
}