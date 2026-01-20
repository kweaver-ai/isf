import * as React from 'react';
import classnames from 'classnames';
import { ValidateBox, Select } from '@/sweet-ui'
import SwitchButton2 from '@/ui/SwitchButton2/ui.desktop';
import Title from '@/ui/Title/ui.desktop';
import ToolBar from '@/ui/ToolBar/ui.desktop';
import Button from '@/ui/Button/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import * as descriptionImg from './assets/description.png';
import SMTPConfigBase from './component.base';
import { ValidateState, SafeMode, TestStatus } from './helper';
import __ from './locale';
import styles from './styles.desktop.css';

export default class SMTPConfig extends SMTPConfigBase {
    render() {
        const { configInfo, isFormChanged, testStatus, validateState, isTestSuccess, testError, saveError } = this.state;
        const ValidateMessages = {
            [ValidateState.Empty]: __('此输入项不允许为空。'),
            [ValidateState.ServerError]: __('SMTP服务器名只能包含 英文、数字 及 @-_. 字符，长度范围 3~100 个字符，请重新输入。'),
            [ValidateState.PortError]: __('请输入1~65535范围内的整数。'),
            [ValidateState.EmailError]: __('邮箱地址只能包含 英文、数字 及 @-_. 字符，格式形如 XXX@XXX.XXX，长度范围 5~100 个字符，请重新输入。'),
        };

        return (
            <div className={styles['container']}>
                <ToolBar>
                    <span className={styles['tool-bar']}>
                        <UIIcon className={styles['toolbar-icon']} size={18} code={'\uf016'} />
                        {__('SMTP服务器设置')}
                    </span>
                </ToolBar>
                <div className={styles['main']}>
                    <div className={styles['table']}>
                        <div className={styles['row']}>
                            <div htmlFor="server" className={styles['label']}>
                                <div className={styles['form-font']}>
                                    {__('邮件服务器（SMTP）:')}
                                    <span className={styles['required']}>*</span>
                                </div>
                            </div>
                            <div className={styles['field']}>
                                <ValidateBox
                                    width={220}
                                    id="server"
                                    value={configInfo.server}
                                    onValueChange={({ detail }) => this.handleChange('server', { server: detail })}
                                    validateMessages={ValidateMessages}
                                    validateState={validateState.server}
                                />
                            </div>
                        </div>
                        <div className={styles['row']}>
                            <div className={styles['label']} htmlFor="safeMode">
                                <div className={styles['form-font']}>
                                    {__('安全连接:')}
                                    <span className={styles['required']}>*</span>
                                </div>
                            </div>
                            <div className={styles['field']} id="safeMode">
                                <Select
                                    value={configInfo.safeMode}
                                    onChange={this.selectChangeHandler.bind(this)}
                                >
                                    <Select.Option
                                        value={SafeMode.Default}
                                        selected={configInfo.safeMode === SafeMode.Default}
                                    >
                                        {__('无')}
                                    </Select.Option>
                                    <Select.Option
                                        value={SafeMode.SslOrTsl}
                                        selected={configInfo.safeMode === SafeMode.SslOrTsl}
                                    >
                                        {'SSL/TLS'}
                                    </Select.Option>
                                    <Select.Option
                                        value={SafeMode.Starttls}
                                        selected={configInfo.safeMode === SafeMode.Starttls}
                                    >
                                        {'STARTTLS'}
                                    </Select.Option>
                                </Select>
                            </div>
                        </div>
                        <div className={styles['row']}>
                            <div className={styles['label']} htmlFor="port">
                                <div className={styles['form-font']}>
                                    {__('端口:')}
                                    <span className={styles['required']}>*</span>
                                </div>
                            </div>
                            <div className={styles['field']}>
                                <ValidateBox
                                    width={220}
                                    id="port"
                                    value={configInfo.port}
                                    onValueChange={({ detail }) => this.handleChange('port', { port: detail })}
                                    validateMessages={ValidateMessages}
                                    validateState={validateState.port}
                                />
                            </div>
                        </div>
                        <div className={styles['row']}>
                            <div className={styles['label']}>
                                <div className={styles['form-font']}>
                                    {__('Open Relay:')}
                                </div>
                            </div>
                            <div className={classnames(styles['field'], styles['switch'])}>
                                <div className={styles['open-relay']}>
                                    <SwitchButton2
                                        active={configInfo.openRelay}
                                        onChange={() => this.switchOpenRelay()}
                                    />
                                </div>
                                <div className={styles['explain']}>
                                    <Title
                                        content={
                                            <div className={styles['explain-content']}>
                                                {__('开启Open Relay，需要邮件服务器已支持Open Relay方能操作成功，开启后，邮箱验证不需要输入密码；关闭Open Relay，邮箱验证则需要输入密码。')}
                                            </div>
                                        }
                                    >
                                        <UIIcon
                                            code={'\uf055'}
                                            size={16}
                                            color={'#69c0ff'}
                                            fallback={descriptionImg}
                                        />
                                        <span className={styles['explain-text']}>{__('说明')}</span>
                                    </Title>
                                </div>
                            </div>
                        </div>
                        <div className={styles['row']}>
                            <div className={styles['label']} htmlFor="email">
                                <div className={styles['form-font']}>
                                    {__('邮箱地址:')}
                                    <span className={styles['required']}>*</span>
                                </div>
                            </div>
                            <div className={styles['field']}>
                                <ValidateBox
                                    width={220}
                                    id="email"
                                    ref={(pwdInput) => this.pwdInput = pwdInput}
                                    value={configInfo.email}
                                    onValueChange={({ detail }) => this.handleChange('email', { email: detail })}
                                    validateMessages={ValidateMessages}
                                    validateState={validateState.email}
                                />
                            </div>
                        </div>
                        <input className={styles['hidden-input']} type="text" />
                        {

                            !configInfo.openRelay ?
                                <div className={styles['row']}>
                                    <div className={styles['label']} htmlFor="password">
                                        <div className={styles['form-font']}>
                                            {__('邮箱密码:')}
                                            <span className={styles['required']}>*</span>
                                        </div>
                                    </div>
                                    <div className={styles['field']}>
                                        <ValidateBox
                                            width={220}
                                            id="password"
                                            type={'password'}
                                            value={configInfo.password}
                                            onValueChange={({ detail }) => this.handleChange('password', { password: detail })}
                                            onClick={() => this.cleanPassword()}
                                            validateMessages={ValidateMessages}
                                            validateState={validateState.password}
                                        />
                                    </div>
                                </div>
                                : null
                        }
                    </div>
                    <div className={styles['row']}>
                        <Button
                            className={styles['btn']}
                            onClick={this.testHandler.bind(this)}
                            width='auto'
                        >
                            {__('测试')}
                        </Button>
                        {
                            isFormChanged ?
                                <div className={styles['form-change-btn']}>
                                    <Button
                                        className={styles['btn']}
                                        onClick={this.saveHandler.bind(this)}
                                    >
                                        {__('保存')}
                                    </Button>
                                    <Button
                                        className={styles['btn']}
                                        onClick={this.cancalHandler.bind(this)}
                                    >
                                        {__('取消')}
                                    </Button>
                                </div>
                                : null
                        }
                        {
                            testStatus === TestStatus.Tested ?
                                <div className={classnames(styles['test'], {
                                    [styles['test-success']]: isTestSuccess,
                                    [styles['test-failed']]: !isTestSuccess,
                                })}>
                                    {
                                        isTestSuccess ?
                                            __('测试连接成功，指定的服务器可用，您可以进入邮箱查看测试邮件。') : (
                                                testError ? (this.getErrorMessage(testError, configInfo.openRelay)) : null
                                            )
                                    }
                                </div>
                                :
                                (
                                    testStatus === TestStatus.Testing ?
                                        <div className={classnames(styles['test'], styles['testing'])}>
                                            {__('测试中...')}
                                        </div>
                                        : null
                                )
                        }
                        {
                            saveError ?
                                <div className={classnames(styles['test'], styles['test-failed'])}>
                                    {this.getErrorMessage(saveError, configInfo.openRelay)}
                                </div>
                                : null
                        }
                    </div>
                </div>
            </div>
        )
    }

    getErrorMessage(err: any, openRelay: boolean): string {
        if (err.error) {
            switch (err.error.errID) {
                case 20807:
                case 20808:
                case 20811:
                case 20812:
                case 20813:
                case 20814:
                    return __('SMTP服务器不可用，请检查服务器地址、安全连接或端口是否正确。');
                case 20103:
                case 20801:
                case 20809:
                case 20810:
                case 20815:
                    return openRelay ? __('测试连接失败，邮箱地址不正确，请重新输入。') : __('测试连接失败，邮箱地址或密码不正确，请重新输入。')
                default:
                    return err.error.errMsg;
            }
        } else {
            return '';
        }
    }
}