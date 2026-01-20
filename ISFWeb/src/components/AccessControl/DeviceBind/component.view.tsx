import * as React from 'react';
import * as classnames from 'classnames';
import DeviceBindBase, { SearchField, OsType, ValidateStates } from './component.base';
import ToolBar from '@/ui/ToolBar/ui.desktop';
import SearchBox from '@/ui/SearchBox/ui.desktop';
import { Select, ValidateBox } from '@/sweet-ui';
import SwitchButton from '@/ui/SwitchButton/ui.desktop';
import Dialog from '@/ui/Dialog2/ui.desktop';
import ConfirmDialog from '@/ui/ConfirmDialog/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import MessageDialog from '@/ui/MessageDialog/ui.desktop';
import ErrorDialog from '@/ui/ErrorDialog/ui.desktop';
import { DataGrid } from '@/sweet-ui';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import Form from '@/ui/Form/ui.desktop';
import Title from '@/ui/Title/ui.desktop';
import { formatTime } from '@/util/formatters'
import { EmptyResult, Centered, Icon, Text } from '@/ui/ui.desktop';
import * as loadImg from './assets/loading.gif'
import __ from './locale';
import styles from './styles.view.css';

const deviceType = {
    '0': __('未知设备'),
    '1': __('iOS'),
    '2': __('Android'),
    '3': __('Windows Phone'),
    '4': __('Windows'),
    '5': __('Mac OS X'),
    '6': __('Web'),
}

const ValidateMessages = {
    [ValidateStates.InvalidMac]: __('MAC地址格式错误'),
}

/**
 * 内链共享列表为空
 */
const EmptyComponent = () => (
    <EmptyResult
        details={'内链共享列表为空'}
    />

)

/**
 * 正在加载列表数据
 */
const RefreshingComponent = () => (
    <Centered>
        <Icon
            url={loadImg}
            size={48}
        />
        <p
            className={styles['loading-text']}
        >
            {__('正在加载，请稍候......')}
        </p>
    </Centered>
)

export default class DeviceBind extends DeviceBindBase {
    render() {
        let { scope, searchResults, deviceInfos, addingDevice, deviceIsExist, addBox: { mac, osType }, validateState, currentSelectUser, uploading, loadingUsers, loadingDevices, errorlist, searchudidKey, clearing } = this.state;
        return (
            <div className={styles['container']}>
                <div className={styles['text-style']}>
                    <label className={styles['tip']}>{__('请将用户与设备绑定，绑定后用户在非绑定的设备上将无法登录')}</label>
                </div>
                <div className={styles['wrapper']}>
                    <div className={styles['table-left']}>
                        <div className={classnames(styles['fl'], styles['user-container'])}>
                            <div className={styles['user-wrapper']}>
                                <div className={styles['info-header']}>
                                    <ToolBar>
                                        <div className={styles['fl']}>
                                            <Select
                                                onChange={this.handleChangeSearchField.bind(this)}
                                                value={scope}
                                            >
                                                <Select.Option selected={scope === SearchField.AllUser} value={SearchField.AllUser}>{__('全部用户')}</Select.Option>
                                                <Select.Option selected={scope === SearchField.BindUser} value={SearchField.BindUser}>{__('已绑定的用户')}</Select.Option>
                                                <Select.Option selected={scope === SearchField.NoBindUser} value={SearchField.NoBindUser}>{__('未绑定的用户')}</Select.Option>
                                            </Select>
                                        </div>
                                        <div className={classnames(styles['fr'], styles['search-box'])}>
                                            <SearchBox
                                                className={styles['search-box-ui']}
                                                role={'ui-searchbox'}
                                                disabled={false}
                                                width={200}
                                                placeholder={__('请输入用户名称')}
                                                value={this.state.searchKey}
                                                onChange={this.updateUserSearchKey}
                                                loader={this.handleSearchKeyChange.bind(this)}
                                            />
                                        </div>
                                    </ToolBar>
                                </div>
                                <div className={styles['dategrid-wrapper']}>
                                    <DataGrid
                                        limit={this.UserListLimit}
                                        height={'100%'}
                                        data={searchResults}
                                        enableSelect={true}
                                        onSelectionChange={this.handleSelectUser.bind(this)}
                                        selection={currentSelectUser}
                                        refreshing={loadingUsers}
                                        RefreshingComponent={RefreshingComponent}
                                        onRequestLazyLoad={this.handleRequestLoadUsers.bind(this)}
                                        columns={[
                                            {
                                                title: __('用户'),
                                                key: 'displayName',
                                                width: 70,
                                                renderCell: (displayName, users) =>
                                                    (
                                                        <div className={styles['displayName']}>
                                                            <Title inline={true} content={displayName}>{displayName}</Title>
                                                        </div>
                                                    ),
                                            },
                                            {
                                                title: __('绑定状态'),
                                                key: 'bindStatus',
                                                width: 30,
                                                minWidth: 80,
                                                renderCell: (bindStatus, users) => (
                                                    bindStatus ?
                                                        <span className={styles['bind-status']}>{__('已绑定')}</span> :
                                                        <span className={styles['nobind-status']}>{__('未绑定')}</span>
                                                ),
                                            },
                                        ]}
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className={styles['table-right']}>
                        <div className={classnames(styles['fr'], styles['device-container'], { [styles['opacity']]: !currentSelectUser })}>
                            <div className={styles['device-wrapper']}>
                                <div className={styles['info-header']}>
                                    <ToolBar>
                                        <ToolBar.Button
                                            icon={'\uf018'}
                                            className={styles['margin']}
                                            disabled={!currentSelectUser}
                                            onClick={this.handleAddDevice.bind(this)}
                                        >
                                            {__('添加')}
                                        </ToolBar.Button>

                                        <div className={styles['file-select']}>
                                            <div
                                                ref="batchImport"
                                                className={
                                                    classnames(
                                                        styles['plainbutton'],
                                                        styles['box-sizing-border-box'],
                                                        {
                                                            [styles['disabled']]: !currentSelectUser,
                                                        },
                                                    )
                                                }
                                                type="button"
                                                disabled={!currentSelectUser}
                                            >
                                                <div 
                                                    className={styles['batch-import']}
                                                >
                                                    <span className={styles['icon']}>
                                                        <UIIcon
                                                            size={'13px'}
                                                            code={'\uf018'}
                                                        />
                                                    </span>
                                                    {
                                                        __('批量导入')
                                                    }
                                                </div>
                                               
                                            </div>
                                            {
                                                !currentSelectUser ? (
                                                    <div
                                                        className={
                                                            classnames(
                                                                styles['plainbutton'],
                                                                styles['box-sizing-border-box'],
                                                                styles['disabled'],
                                                                styles['btn-cover'],
                                                            )
                                                        }
                                                        type="button"
                                                        disabled={true}
                                                    >
                                                        <div 
                                                            className={styles['batch-import']}
                                                            style={{
                                                                opacity: 0.5,
                                                            }}
                                                        >
                                                            <span className={styles['icon']}>
                                                                <UIIcon
                                                                    size={'13px'}
                                                                    code={'\uf018'}
                                                                />
                                                            </span>
                                                            {
                                                                __('批量导入')
                                                            }
                                                        </div>
                                                    </div>
                                                ) : null
                                            }
                                        </div>

                                        <ToolBar.Button
                                            icon={'\uf000'}
                                            className={styles['margin']}
                                            disabled={!currentSelectUser || !deviceInfos.length}
                                            onClick={() => {
                                                this.setState({
                                                    clearing: true,
                                                })
                                            }}
                                        >
                                            {__('清空')}
                                        </ToolBar.Button>
                                        <div className={classnames(styles['fr'], styles['search-box'])}>
                                            <SearchBox
                                                className={styles['search-box-ui']}
                                                role={'ui-searchbox'}
                                                disabled={!currentSelectUser}
                                                width={176}
                                                value={searchudidKey}
                                                placeholder={__('请输入设备识别码')}
                                                onChange={this.updateUDIDSearchKey}
                                                loader={this.handleSearchUdidChange.bind(this)}
                                            />
                                        </div>
                                    </ToolBar>
                                </div>
                                <div className={styles['dategrid-wrapper']}>
                                    <DataGrid
                                        height={'100%'}
                                        data={deviceInfos}
                                        limit={this.DeviceListLimit}
                                        refreshing={loadingDevices}
                                        RefreshingComponent={RefreshingComponent}
                                        onRequestLazyLoad={this.handleRequestLoadDevice.bind(this)}
                                        columns={[
                                            {
                                                title: __('设备识别码'),
                                                key: 'udid',
                                                width: 30,
                                                renderCell: (udid, deviceInfo) => (<Text>{deviceInfo.baseInfo.udid}</Text>),
                                            },
                                            {
                                                title: __('设备类型'),
                                                key: 'osType',
                                                width: 20,
                                                renderCell: (osType, deviceInfo) => (<Text>{deviceType[deviceInfo.baseInfo.osType]}</Text>),
                                            },
                                            {
                                                title: __('最后登录时间'),
                                                key: 'lastLoginTime',
                                                width: 25,
                                                renderCell: (lastLoginTime, deviceInfo) => (
                                                    <Text>{deviceInfo.baseInfo.lastLoginTime === -1 ? '--' : formatTime(deviceInfo.baseInfo.lastLoginTime / 1000, 'yyyy/MM/dd')}</Text>
                                                ),
                                            },
                                            {
                                                title: __('操作'),
                                                key: 'bindFlag',
                                                width: 25,
                                                minWidth: 100,
                                                renderCell: (bindFlag, deviceInfo) => (
                                                    <div>
                                                        <SwitchButton
                                                            disabled={!currentSelectUser}
                                                            active={bindFlag ? true : false}
                                                            onChange={() => this.switchDeviceBind(deviceInfo)}
                                                        />
                                                        <UIIcon
                                                            disabled={!currentSelectUser}
                                                            size={15}
                                                            className={styles['delete-icon']}
                                                            code={'\uf013'}
                                                            color={'#9a9a9a'}
                                                            onClick={() => this.deleteDeviceBind(deviceInfo)}
                                                        />
                                                    </div>
                                                ),
                                            },
                                        ]}
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                {
                    addingDevice ?
                        <Dialog
                            title={__('添加设备')}
                            onClose={this.handleCancelAddDevice.bind(this)}
                        >
                            <Panel>
                                <Panel.Main>
                                    <Form>
                                        <Form.Row>
                                            <Form.Label>
                                                <label className={styles['device-label']}>
                                                    {__('设备类型')}
                                                </label>
                                            </Form.Label>
                                            <Form.Field>
                                                <div className={styles['device-field']}>
                                                    <Select value={osType} onChange={({ detail: osType }) => this.handleAddBoxChange({ osType })}>
                                                        <Select.Option selected={osType === OsType.Windows} value={OsType.Windows}>{__('Windows')}</Select.Option>
                                                        <Select.Option selected={osType === OsType.MacOS} value={OsType.MacOS}>{__('Mac OS X')}</Select.Option>
                                                    </Select>
                                                </div>
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row>
                                            <Form.Label>
                                                <label className={styles['device-label']}>
                                                    {__('设备识别码')}
                                                </label>
                                            </Form.Label>
                                            <Form.Field>
                                                <div className={styles['device-field']}>
                                                    <ValidateBox
                                                        width={220}
                                                        validateMessages={ValidateMessages}
                                                        validateState={validateState}
                                                        value={mac}
                                                        placeholder={__('MAC地址,如 00-00-00-00-00-E0')}
                                                        onValueChange={({detail: mac}) => this.handleAddBoxChange({ mac })}
                                                    />
                                                </div>
                                            </Form.Field>
                                        </Form.Row>
                                    </Form>
                                </Panel.Main>
                                <Panel.Footer>
                                    <Panel.Button theme='oem' onClick={this.handleSubmitAddDevice.bind(this)} width='auto'>{__('确定')}</Panel.Button>
                                    <Panel.Button onClick={this.handleCancelAddDevice.bind(this)} width='auto'>{__('取消')}</Panel.Button>
                                </Panel.Footer>
                            </Panel>
                        </Dialog> :
                        null
                }
                {
                    deviceIsExist ?
                        <MessageDialog onConfirm={this.handleCancelTip.bind(this)}>
                            {__('该设备已经存在，无法再添加')}
                        </MessageDialog> :
                        null
                }
                {
                    errorlist ? (
                        <ErrorDialog onConfirm={() => this.setState({
                            errorlist: null,
                        })}>
                            <ErrorDialog.Title>
                                {__('导入失败，错误信息如下:')}
                            </ErrorDialog.Title>
                            <ErrorDialog.Detail>
                                <p>{__('以下MAC地址格式错误：')}</p>
                                {errorlist.map((error, index) => (
                                    <p key={index}>{error}</p>
                                ))}
                            </ErrorDialog.Detail>
                        </ErrorDialog>
                    ) : null
                }
                {
                    uploading ?
                        <div className={styles['uploading-wrap']}>
                            <Centered>
                                <Icon
                                    url={loadImg}
                                    size={48}
                                />
                                <p
                                    className={styles['loading-text']}
                                >
                                    {__('正在导入，请稍候......')}
                                </p>
                            </Centered>
                        </div>
                        : null
                }
                {
                    clearing ?
                        <ConfirmDialog
                            onConfirm={this.handleClearDevices.bind(this)}
                            onCancel={() => {
                                this.setState({
                                    clearing: false,
                                })
                            }}
                        >
                            {__('此操作将清空该用户绑定的所有设备识别码，您确定要执行此操作吗？')}
                        </ConfirmDialog> :
                        null
                }
            </div>
        );
    }
}
