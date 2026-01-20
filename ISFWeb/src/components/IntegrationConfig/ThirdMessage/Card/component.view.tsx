import * as React from 'react';
import * as classnames from 'classnames';
import { Form, Button, UIIcon, Panel, Title, Text, ConfirmDialog, Icon, Overlay } from '@/ui/ui.desktop';
import { Switch, ValidateBox } from '@/sweet-ui';
import Dialog from '@/ui/Dialog2/ui.desktop';
import AdvancedConfig from './AdvancedConfig/component.view';
import ResetInvalidConfig from './ResetInvalidConfig/component.view';
import MessageTypes from './MessageTypes/component.view';
import ConfigItem from './ConfigItem/component.view';
import { ValidateMessage, formatterError } from '../helper';
import CardBase from './component.base';
import * as loading from './assets/loading.gif';
import styles from './styles.view.css';
import __ from './locale';

export default class Card extends CardBase {
    render() {
        const { showDeleteIcon } = this.props
        const {
            enabled,
            thirdPartyName,
            pluginClassName,
            internalConfig,
            showInternalAdvancedConfigDialog,
            edited,
            error,
            invalidConfig,
            validateStatus: {
                thirdPartyNameValidateStatus,
                pluginClassNameValidateStatus,
            },
            messages,
            showDeleteCardDialog,
            plugin,
            errorToast,
        } = this.state;

        // 高级配置config
        const config = {
            thirdparty_name: thirdPartyName,
            enabled,
            class_name: pluginClassName,
            channels: messages,
            config: internalConfig,
        }

        return (
            <div>
                <div className={styles['container']}>
                    {
                        showDeleteIcon ?
                            <UIIcon
                                code={'\uf013'}
                                color={'#505050'}
                                className={styles['delete-btn']}
                                onClick={() => this.handleDelete()}
                            />
                            : null
                    }
                    <Form>
                        <Form.Row>
                            <Form.Label>{__('状态：')}</Form.Label>
                            <Form.Field className={styles['field']}>
                                <div className={styles['status-container']}>
                                    <div className={styles['mark']}></div>
                                    <div className={styles['switch']}>
                                        <Switch
                                            checked={enabled}
                                            onChange={({ detail }) => this.statusChange(detail)}
                                        />
                                    </div>
                                    <Button
                                        className={classnames(styles['button'], styles['Advanced-button'])}
                                        onClick={() => this.openInternalAdvancedConfigDialog()}
                                    >
                                        {__('高级配置')}
                                    </Button>
                                </div>
                            </Form.Field>
                        </Form.Row>
                        <Form.Row className={styles['row']}>
                            <Form.Label>{__('消息服务名称：')}</Form.Label>
                            <Form.Field className={styles['field']}>
                                <div className={styles['mark']}>{'*'}</div>
                                <ValidateBox
                                    width={300}
                                    value={thirdPartyName}
                                    validateState={thirdPartyNameValidateStatus}
                                    validateMessages={ValidateMessage}
                                    onValueChange={({ detail: value }) => this.handleThirdPartyNameChange(value)}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row  className={styles['row']}>
                            <Form.Label>{__('插件类名：')}</Form.Label>
                            <Form.Field className={styles['field']}>
                                <div className={styles['mark']}>{'*'}</div>
                                <ValidateBox
                                    width={300}
                                    value={pluginClassName}
                                    validateState={pluginClassNameValidateStatus}
                                    validateMessages={ValidateMessage}
                                    onValueChange={({ detail: value }) => this.handlePluginClassNameChange(value)}
                                />
                            </Form.Field>
                        </Form.Row>
                    </Form>

                    <div>
                        <span> {__('插件参数配置：')}</span>
                        <Button
                            className={classnames(styles['button'])}
                            onClick={() => this.addInternalConfig()}
                        >
                            {__('添加')}
                        </Button>
                    </div>
                    <div className={styles['item-wrapper']} ref={(ref) => this.internalConfigContainer = ref}>
                        {
                            internalConfig.map((item, configIndex) => {
                                return <ConfigItem
                                    key={configIndex}
                                    configItem={item}
                                    onRequestConfigNameChange={(value) => this.internalConfigNameChange(value, configIndex)}
                                    onRequestConfigTypeChange={(type) => this.internalConfigTypeChange(type, configIndex)}
                                    onRequestConfigValueChange={(value, type) => this.internalConfigValueChange(value, type, configIndex)}
                                    onRequestDelete={() => this.deleteInternalConfigItem(configIndex)}
                                />
                            })
                        }
                    </div>

                    <div>
                        <span className={styles['message-type']}>{__('消息类型：')}</span>
                        <span className={styles['mark']}>{'*'}</span>
                        <Button
                            className={classnames(styles['button'])}
                            onClick={this.addMessageConfig}
                        >
                            {__('添加')}
                        </Button>
                    </div>
                    <div className={styles['item-wrapper']} ref={(ref) => this.messageConfigItemContainer = ref}>
                        {
                            messages.map((item, messageConfigIndex) => (
                                <MessageTypes
                                    key={messageConfigIndex}
                                    messageConfigItem={item}
                                    onRequestConfigValueChange={(value) => this.handleConfigValueChange(value, messageConfigIndex)}
                                    onRequestRef={this.handleMessagesItemRef}
                                    onRequestDelete={() => this.deleteMessageConfigItem(messageConfigIndex)}
                                />
                            ))
                        }
                    </div>

                    {
                        edited ?
                            <div>
                                <Button onClick={() => this.save()} className={styles['save-button']}>{__('保存')}</Button>
                                <Button onClick={() => this.cancel()} className={styles['save-button']}>{__('取消')}</Button>
                                {
                                    error !== null ? <div className={styles['error']}>{formatterError(error.error.errID)}</div> : null
                                }
                            </div>
                            :
                            null
                    }
                    <div className={styles['plugin']}>
                        <div className={styles['row']}>
                            <div className={styles['plugin-tip']}>{__('消息模块插件：')}</div>
                            <div className={styles['mark']}>{'*'}</div>
                            <div className={styles['plugin-content']}>
                                <div className={classnames(styles['file-name'], { [styles['disabled']]: edited || !this.indexId })}>
                                    {plugin && plugin.filename ? <Text>{plugin.filename}</Text> : null}
                                </div>
                                <div className={styles['btn-wraper']}>
                                    <div ref={(ref) => this.select = ref} className={styles['btn-uploader-picker']}>
                                        {__('上传')}
                                    </div>
                                    {
                                        edited || !this.indexId ?
                                            <div className={styles['disabled']}>
                                                {__('上传')}
                                            </div>
                                            : null
                                    }
                                </div>
                                <Title content={__('消息模块插件必须在启用且保存上方配置后，才能上传。')}>
                                    <UIIcon
                                        code={'\uf055'}
                                        size={16}
                                        color={'#69C0FF'}
                                    />
                                    <p className={styles['tips']}>{__('说明')}</p>
                                </Title>
                            </div>
                        </div>
                        <div className={styles['row']}>
                            <div className={styles['plugin-tip-empty']}></div>
                            <div className={styles['mark-empty']}></div>
                            <div className={styles['format-req']}>{__('请从本地选择tar.gz的格式文件上传，大小不超过200MB')}</div>
                        </div>
                    </div>
                </div>
                {
                    this.state.uploading ?
                        <Dialog
                            width={450}
                            title={__('提示')}
                            buttons={[]}
                        >
                            <Panel>
                                <Panel.Main>
                                    <div className={styles['main']}>
                                        <div className={styles['icon']}>
                                            <Icon url={loading} size={44} />
                                        </div>
                                        <div className={styles['message']}>
                                            {__('正在上传消息模块插件，请稍候...')}
                                        </div>
                                    </div>
                                </Panel.Main>
                            </Panel>
                        </Dialog>
                        :
                        null
                }
                {
                    showDeleteCardDialog ? (
                        <ConfirmDialog
                            onConfirm={() => { this.handleDeleteSavedCard() }}
                            onCancel={() => this.handleCancleDelete()}
                        >
                            {__('此操作将删除该第三方消息集成，系统不再向该应用推送消息，您确定要执行此操作吗？')}
                        </ConfirmDialog >
                    ) : null
                }
                {
                    showInternalAdvancedConfigDialog ?
                        <AdvancedConfig
                            originalConfig={config}
                            onRequestClose={() => this.closeInternalAdvancedConfigDialog()}
                            onRequestConfirm={(config) => this.updateConfig(config)}
                        />
                        : null
                }
                {
                    invalidConfig ?
                        <ResetInvalidConfig
                            invalidConfig={invalidConfig}
                            onRequestConfirm={() => this.resetInvalidConfig()}
                            onRequestClose={() => this.resetInvalidConfig()}
                        />
                        : null
                }
                {
                    errorToast ?
                        <Overlay
                            className={styles['toast']}
                            position="top center"
                        >
                            {formatterError(errorToast.errCode)}
                        </Overlay>
                        : null
                }
            </div >
        )
    }
}