import * as React from 'react';
import ToolBar from '@/ui/ToolBar/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { generatePwdWithSpecialChar } from '@/core/password';
import { Select, CheckBox, Button, ValidateNumber, ValidateBox } from '@/sweet-ui';
import { PwdValidity, PwdStrength, ValidateState, hiddenPwd } from './helper';
import __ from './locale';
import PasswordPolicyBase from './component.base';
import styles from './styles.view';

export default class PasswordPolicy extends PasswordPolicyBase {
    render() {
        const {
            configInfo,
            pwdPolicy: {
                max_err_count,
                min_strong_pwd_length,
                weak_pwd_disabled,
            },
            lockStatusDisabled,
            isSMSReset,
            isEmailReset,
            isSMSConfig,
            isEmailConfig,
            isTriSystem,
            timeCtrl,
            isChanged,
            validateState,
            resetViaSMSDisabled,
            loading,
        } = this.state;

        return !loading && (
            <div className={styles['container']}>
                <ToolBar className={styles['too-bar']}>
                    <span className={styles['tool-bar']}>
                        <UIIcon className={styles['toolbar-icon']} size={18} code={'\uf016'} />
                        {__('密码策略')}
                    </span>
                </ToolBar>
                <div className={styles['main']}>
                    <div className={styles['title']}>
                        {__('用户初始密码设置：')}
                    </div>
                    <div className={styles['title']}>
                        <span className={styles['text-middle']}>
                            {__('初始密码：')}
                        </span>
                        <div className={styles['password']}>
                            <ValidateBox
                            role={'sweetui-validatebox'}
                            width={200}
                            value={configInfo.initPwd}
                            disabled={configInfo.initPwd === hiddenPwd}
                            onValueChange={({ detail }) => this.updateInitPwd(detail)}
                            validateState={validateState.initPwdState}
                            validateMessages={this.getValidateMessages()}
                            onBlur={this.checkInitPwd}
                        />
                        </div>
                        {
                            configInfo.initPwd === hiddenPwd ?
                                <Button
                                    role={'ui-button'}
                                    className={styles['initpwd-btn']}
                                    size={'auto'}
                                    onClick={() => this.updateInitPwd('')}
                                >
                                    {__('重置密码')}
                                </Button>
                                :
                                <Button
                                    role={'ui-button'}
                                    className={styles['initpwd-btn']}
                                    size={'auto'}
                                    onClick={() => this.updateInitPwd(generatePwdWithSpecialChar(10))}
                                >
                                    {__('随机密码')}
                                </Button>
                        }
                    </div>
                    <div className={styles['explain']}>{__('初始密码生成后无法再次查看，请妥善保管')}</div>
                    <div className={styles['title']}>
                        <span className={styles['text-middle']}>
                            {__('密码有效期：')}
                        </span>
                        <Select
                            className={styles['text-middle']}
                            width={200}
                            value={configInfo.expireTime}
                            onChange={({ detail }) => this.selectPwdValidity({ detail })}
                        >
                            {
                                this.expireTimeSelected(timeCtrl).map((item) => {
                                    return (
                                        <Select.Option
                                            key={item}
                                            value={item}
                                            selected={configInfo.expireTime === item}>
                                            {this.formateExpireTime(item)}
                                        </Select.Option>
                                    )
                                })
                            }
                        </Select>
                    </div>
                    <div className={styles['explain']}>
                        <span>{__('密码仅在指定时间段内有效，若超过该有效期，则需要修改密码，否则无法登录')}</span>
                    </div>
                    <div className={styles['title']}>
                        <span className={styles['text-middle']}>{__('密码强度：')}</span>
                        {
                            weak_pwd_disabled ? (
                                <Select
                                    className={styles['text-middle']}
                                    width={200}
                                    value={configInfo.strongStatus}
                                    disabled={configInfo.pwdLockStatus}
                                    onChange={(value) => this.selectPwdStrength(value)}
                                >
                                    <Select.Option
                                        value={PwdStrength.StrongPwd}
                                        selected={configInfo.strongStatus === PwdStrength.StrongPwd}>
                                        {__('强密码')}
                                    </Select.Option>
                                </Select>
                            ) :
                                (
                                    <Select
                                        className={styles['text-middle']}
                                        width={200}
                                        value={configInfo.strongStatus}
                                        disabled={configInfo.pwdLockStatus}
                                        onChange={(value) => this.selectPwdStrength(value)}
                                    >
                                        <Select.Option
                                            value={PwdStrength.StrongPwd}
                                            selected={configInfo.strongStatus === PwdStrength.StrongPwd}>
                                            {__('强密码')}
                                        </Select.Option>
                                        <Select.Option
                                            value={PwdStrength.WeakPwd}
                                            selected={configInfo.strongStatus === PwdStrength.WeakPwd}
                                        >{__('弱密码')}
                                        </Select.Option>
                                    </Select>
                                )
                        }
                    </div>
                    {
                        configInfo.strongStatus ?
                            <div className={styles['explain']}>
                                <span>{__('强密码格式：密码长度至少为')}</span>
                                <div className={styles['number-input']}>
                                    <ValidateNumber
                                        width={160}
                                        step={1}
                                        precision={0}
                                        min={0}
                                        maxLength={2}
                                        max={configInfo.strongPwdLength ? 99 : null}
                                        placeholder={__('请输入${min}-${max}的数值', { min: min_strong_pwd_length, max: 99 })}
                                        value={configInfo.strongPwdLength}
                                        validateState={validateState.strongPwdLength}
                                        validateMessages={this.getValidateMessages()}
                                        disabled={configInfo.pwdLockStatus}
                                        onBlur={(event, value) => this.handleBlurValidate('strongPwdLength', value)}
                                        onValueChange={({ detail }) => this.handleChange({ strongPwdLength: detail })}
                                    />
                                </div>

                                <span>{__('个字符，需同时包含 大小写英文字母、数字与特殊字符')}</span>
                            </div>
                            :
                            <div className={styles['explain']}>
                                {__('弱密码格式：密码长度至少为6个字符')}
                            </div>
                    }
                    <div className={styles['title']}>
                        <label>
                            <CheckBox
                                checked={configInfo.lockStatus}
                                onClick={(event) => event.stopPropagation()}
                                onChange={() => this.openPwdLock()}
                                disabled={lockStatusDisabled}
                            />
                            <div className={styles['checked-text']}>{__('启用密码错误锁定：')}</div>
                        </label>
                    </div>
                    <div className={styles['explain']}>
                        <span>
                            {__('用户在任意情况下，密码连续输错')}
                        </span>
                        <div className={styles['number-input']}>
                            <ValidateNumber
                                width={160}
                                min={0}
                                step={1}
                                precision={0}
                                maxLength={max_err_count.toString().length}
                                placeholder={__('请输入${min}-${max}的数值', { min: 1, max: max_err_count })}
                                value={configInfo.passwdErrCnt}
                                disabled={!configInfo.lockStatus || lockStatusDisabled}
                                validateState={validateState.passwdErrCnt}
                                validateMessages={this.getValidateMessages()}
                                onBlur={(event, value) => this.handleBlurValidate('passwdErrCnt', value)}
                                onValueChange={({ detail }) => this.handleChange({ passwdErrCnt: detail })}
                            />
                        </div>
                        <span>{__('次，则账号被锁定')}</span>
                    </div>
                    <div className={styles['explain']}>
                        <span>
                            {__('账号被锁定后，')}
                        </span>
                        <div className={styles['number-input']}>
                            <ValidateNumber
                                width={170}
                                max={configInfo.passwdLockTime ? 999 : null}
                                min={0}
                                step={1}
                                precision={0}
                                maxLength={3}
                                placeholder={__('请输入${min}-${max}的数值', { min: 10, max: 180 })}
                                disabled={!configInfo.lockStatus || lockStatusDisabled}
                                value={configInfo.passwdLockTime}
                                validateState={validateState.passwdLockTime}
                                validateMessages={this.getValidateMessages()}
                                onBlur={(event, value) => this.handleBlurValidate('passwdLockTime', value)}
                                onValueChange={({ detail }) => this.handleChange({ passwdLockTime: detail })}
                            />
                        </div>
                        <span>{__('分钟后自动解锁或由管理员解锁')}</span>
                    </div>
                    <div>
                        <div className={styles['title']}>
                            <span>{__('忘记密码重置：')}</span>
                        </div>
                        {
                            !resetViaSMSDisabled && (
                                <div className={styles['reset-pwd']}>
                                    <label>
                                        <CheckBox
                                            checked={isSMSReset}
                                            onClick={(event) => event.stopPropagation()}
                                            onChange={() => this.selectSMSReset()}
                                        />
                                        <div className={styles['checked-text']}>{__(('通过短信验证'))}</div>
                                    </label>
                                    {
                                        isSMSConfig ?
                                            <span className={styles['explain']}>
                                                {__('（短信验证需要配置对应的短信服务器插件才能生效，如果您还没有配置，可以在')}
                                                {
                                                    !isTriSystem ?
                                                        <a
                                                            className={styles['link']}
                                                            onClick={() => this.openThirdMessage()}
                                                        >{__(' 第三方消息插件 ')}</a>
                                                        :
                                                        <span>{__(' 第三方消息插件 ')}</span>
                                                }
                                                {__('页面进行操作）')}
                                            </span>
                                            : null
                                    }
                                </div>
                            )
                        }
                        <div className={styles['reset-pwd']}>
                            <label>
                                <CheckBox
                                    checked={isEmailReset}
                                    onClick={(event) => event.stopPropagation()}
                                    onChange={() => this.selectEmailReset()}
                                />
                                <div className={styles['checked-text']}>{__(('通过邮箱验证'))}</div>
                            </label>
                            {
                                isEmailConfig ?
                                    <span className={styles['explain']}>
                                        {__('（邮箱验证需要配置对应的邮箱服务器插件才能生效，如果您还没有配置，可以在')}
                                        {
                                            !isTriSystem ?
                                                <a
                                                    className={styles['link']}
                                                    onClick={() => this.openThirdPartyServer()}
                                                >
                                                    {__(' 邮件服务 ')}
                                                </a>
                                                :
                                                <span>{__(' 邮件服务 ')}</span>
                                        }
                                        {__('页面进行操作）')}
                                    </span>
                                    : null
                            }
                        </div>
                        <div className={styles['explain']}>
                            <span>{__('用户忘记密码时，可以通过绑定的${phone}邮箱发送验证码验证身份，重新设置密码（管控密码的用户除外）', {
                                phone: resetViaSMSDisabled ? '' : __('手机或'),
                            })}</span>
                        </div>
                    </div>

                    {
                        isChanged ?
                            <div className={styles['change-btn']}>
                                <Button  
                                    className={styles['btn']}
                                    onClick={this.saveConfigInfo.bind(this)}
                                >
                                    {__('保存')}
                                </Button>
                                <Button
                                    className={styles['btn']}
                                    onClick={this.cancalConfigInfo.bind(this)}
                                >
                                    {__('取消')}
                                </Button>
                            </div>
                            : null
                    }
                </div>
            </div >
        )
    }

    formateExpireTime(expireTime) {
        switch (expireTime) {
            case PwdValidity.OneDay:
                return 1 + __('天');
            case PwdValidity.TreeDays:
                return 3 + __('天');
            case PwdValidity.SevenDays:
                return 7 + __('天');
            case PwdValidity.OneMonth:
                return 1 + __('个月');
            case PwdValidity.TreeMonths:
                return 3 + __('个月');
            case PwdValidity.SixMonths:
                return 6 + __('个月');
            case PwdValidity.TwelveMonths:
                return 12 + __('个月');
            case PwdValidity.Permanent:
                return __('永久有效');
        }
    }

    getValidateMessages() {
        const { pwdPolicy: { min_strong_pwd_length, max_err_count } } = this.state

        return {
            [ValidateState.Empty]: __('此项不允许为空。'),
            [ValidateState.InvalidStrongPwdLengthMin]: __('强密码长度至少为${min}个字符。', { min: min_strong_pwd_length }),
            [ValidateState.InvalidStrongPwdLengthMax]: __('强密码长度最多为99个字符。'),
            [ValidateState.InvalidPasswdErrCnt]: __('密码错误次数范围为1~${max}。', { max: max_err_count }),
            [ValidateState.InvalidPasswdLockTime]: __('锁定时间范围为10~180分钟。'),
            ...this.errMsg,
        }
    }
}