import * as React from 'react';
import classnames from 'classnames';
import { Form, Button, UIIcon, Panel, Title, Text, Icon, Overlay } from '@/ui/ui.desktop';
import Dialog from '@/ui/Dialog2/ui.desktop';
import { Switch, ValidateBox } from '@/sweet-ui';
import AdvancedConfig from './AdvancedConfig/component.view';
import ResetInvalidConfig from './ResetInvalidConfig/component.view';
import ConfigItem from './ConfigItem/component.view';
import { ValidateMessage, formatterError } from './helper';
import ThirdConfigBase from './component.base';
import * as loading from './assets/loading.gif';
import styles from './styles.view';
import __ from './locale';

export default class ThirdConfig extends ThirdConfigBase {
    render() {
        const {
            enabled,
            thirdPartyId,
            thirdPartyName,
            clientConfig,
            internalConfig,
            showClientAdvancedConfigDialog,
            showInternalAdvancedConfigDialog,
            edited,
            error,
            invalidConfig,
            validateStatus,
            plugin,
            uploading,
            errorToast,
        } = this.state

        return (
            <div>
                <div className={styles['header']}>
                    <UIIcon code={'\uf016'} size={16} />
                    <span className={styles['header-title']}>{__('第三方认证通用配置')}</span>
                </div>
                <div className={styles['container']}>
                    <Form>
                        <Form.Row>
                            <Form.Label>{__('状态：')}</Form.Label>
                            <Form.Field>
                                <div className={styles['mark']}></div>
                                <div className={styles['switch']}>
                                    <Switch
                                        checked={enabled}
                                        onChange={({ detail }) => this.statusChange(detail)}
                                    />
                                </div>
                            </Form.Field>
                        </Form.Row>
                        <Form.Row>
                            <Form.Label>{__('认证服务ID：')}</Form.Label>
                            <Form.Field className={styles['field']}>
                                <div className={styles['mark']}>{'*'}</div>
                                <ValidateBox
                                    width={300}
                                    disabled={!enabled}
                                    value={thirdPartyId}
                                    validateState={validateStatus.thirdPartyIdValidateStatus}
                                    validateMessages={ValidateMessage}
                                    onValueChange={({detail}) => this.handleThirdPartyIdChange(detail)}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row>
                            <Form.Label>{__('认证服务名称：')}</Form.Label>
                            <Form.Field className={styles['field']}>
                                <div className={styles['mark']}>{'*'}</div>
                                <ValidateBox
                                    width={300}
                                    disabled={!enabled}
                                    value={thirdPartyName}
                                    validateState={validateStatus.thirdPartyNameValidateStatus}
                                    validateMessages={ValidateMessage}
                                    onValueChange={({detail}) => this.handleThirdPartyNameChange(detail)}
                                />
                            </Form.Field>
                        </Form.Row>
                    </Form>

                    <div>
                        <span> {__('客户端参数配置：')}</span>
                        <Button
                            disabled={!enabled}
                            className={classnames(styles['button'], { [styles['btn-disabled']]: !enabled })}
                            onClick={() => this.addClientConfig()}
                        >
                            {__('添加')}
                        </Button>
                        <Button
                            disabled={!enabled}
                            className={classnames(styles['button'], { [styles['btn-disabled']]: !enabled })}
                            onClick={() => this.openClientAdvancedConfigDialog()}
                        >
                            {__('高级配置')}
                        </Button>
                    </div>
                    <div
                        className={styles['item-wrapper']}
                        ref={(ref) => this.clientConfigContainer = ref}
                    >
                        {
                            clientConfig.map((item, configIndex) => {
                                return <ConfigItem
                                    key={configIndex}
                                    enabled={enabled}
                                    configItem={item}
                                    onRequestConfigNameChange={(value) => this.clientConfigNameChange(value, configIndex)}
                                    onRequestConfigTypeChange={(type) => this.clientConfigTypeChange(type, configIndex)}
                                    onRequestConfigValueChange={(value, type) => this.clientConfigValueChange(value, type, configIndex)}
                                    onRequestDelete={() => this.deleteClientConfigItem(configIndex)}
                                />
                            })
                        }
                    </div>
                    <div>
                        <span> {__('服务端参数配置：')}</span>
                        <Button
                            disabled={!enabled}
                            className={classnames(styles['button'], { [styles['btn-disabled']]: !enabled })}
                            onClick={() => this.addInternalConfig()}
                        >
                            {__('添加')}
                        </Button>
                        <Button
                            disabled={!enabled}
                            className={classnames(styles['button'], { [styles['btn-disabled']]: !enabled })}
                            onClick={() => this.openInternalAdvancedConfigDialog()}
                        >
                            {__('高级配置')}
                        </Button>
                    </div>
                    <div
                        className={styles['item-wrapper']}
                        ref={(ref) => this.internalConfigContainer = ref}
                    >
                        {
                            internalConfig.map((item, configIndex) => {
                                return <ConfigItem
                                    key={configIndex}
                                    enabled={enabled}
                                    configItem={item}
                                    onRequestConfigNameChange={(value) => this.internalConfigNameChange(value, configIndex)}
                                    onRequestConfigTypeChange={(type) => this.internalConfigTypeChange(type, configIndex)}
                                    onRequestConfigValueChange={(value, type) => this.internalConfigValueChange(value, type, configIndex)}
                                    onRequestDelete={() => this.deleteInternalConfigItem(configIndex)}
                                />
                            })
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
                            <div className={styles['plugin-tip']}>{__('认证模块插件：')}</div>
                            <div className={styles['mark']}>{'*'}</div>
                            <div className={styles['plugin-content']}>
                                <div className={classnames(styles['file-name'], { [styles['disabled']]: edited || !this.indexId || !enabled })}>
                                    {plugin && plugin.filename ? <Text>{plugin.filename}</Text> : null}
                                </div>
                                <div className={styles['btn-wraper']}>
                                    <div ref={(ref) => this.select = ref} className={styles['btn-uploader-picker']}>
                                        {__('上传')}
                                    </div>
                                    {
                                        edited || !this.indexId || !enabled ?
                                            <div className={styles['disabled']}>
                                                {__('上传')}
                                            </div>
                                            : null
                                    }
                                </div>
                                <Title content={__('上传此插件，需启用并保存上方配置')}>
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
                    uploading ?
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
                                            {__('正在上传认证模块插件，请稍候...')}
                                        </div>
                                    </div>
                                </Panel.Main>
                            </Panel>
                        </Dialog>
                        :
                        null
                }
                {
                    showClientAdvancedConfigDialog ?
                        <AdvancedConfig
                            title={__('客户端参数配置：')}
                            originalConfig={clientConfig}
                            onRequestClose={() => this.closeClientAdvancedConfigDialog()}
                            onRequestConfirm={(config) => this.updateClientConfig(config)}
                        /> : null
                }
                {
                    showInternalAdvancedConfigDialog ?
                        <AdvancedConfig
                            title={__('服务端参数配置：')}
                            originalConfig={internalConfig}
                            onRequestClose={() => this.closeInternalAdvancedConfigDialog()}
                            onRequestConfirm={(config) => this.updateInternalConfig(config)}
                        />
                        : null
                }
                {
                    invalidConfig && invalidConfig.length > 0 ?
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