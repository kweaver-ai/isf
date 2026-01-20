import * as React from 'react';
import { LoginWay, LoginAuthType } from '@/core/logincertification/logincertification';
import { EndpointType } from '@/core/endpointtype/endpointtype';
import ToolBar from '@/ui/ToolBar/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { Select, CheckBox, Button, ValidateNumber } from '@/sweet-ui';
import { OstypeInfo, DisableTime, ValidateMessages } from './helper';
import __ from './locale';
import LoginPolicyBase from './component.base';
import styles from './styles.view';

export default class LoginPolicy extends LoginPolicyBase {
    render() {
        const {
            idCardStatus,
            forbidOstypeInfo,
            autoDisableDays,
            autoDisableStatus,
            isChanged,
            loginWay,
            passwdErrCnt,
            passwdErrCntValidity,
            clientLockStatus,
            multiAutnLocked,
            loginStrategyStatus,
            domainUserExemptLogin,
            disabledEndpointTypes,
            disabledLoginAuth,
            loading,
            loginStrategyDisabled,
        } = this.state;

        return !loading && (
            <div className={styles['container']}>
                <ToolBar className={styles['too-bar']}>
                    <span className={styles['tool-bar']}>
                        <UIIcon
                            className={styles['toolbar-icon']}
                            size={18}
                            code={'\uf016'} />
                        {__('登录策略')}
                    </span>
                </ToolBar>
                <div className={styles['main']}>
                    {/* <div className={styles['title']}>
                        <label className={styles['item-label']}>
                            <CheckBox
                                checked={loginStrategyStatus}
                                onClick={(event) => event.stopPropagation()}
                                onChange={() => this.setloginStrategyStatus()}
                                disabled={loginStrategyDisabled}
                            />
                            <div className={styles['checked-text']}>{__('禁止同一帐号多地同时登录（仅对使用Windows客户端登录进行限制，使用其它终端登录不受限制）')}</div>
                        </label>
                    </div> */}
                    <div className={styles['title']}>
                        <label className={styles['item-label']}>
                            <CheckBox
                                checked={domainUserExemptLogin}
                                onClick={(event) => event.stopPropagation()}
                                onChange={() => this.setDomainUserExemptLogin()}
                            />
                            <div className={styles['checked-text']}>{__('域用户免登录')}</div>
                        </label>
                    </div>
                    <div className={styles['title']}>
                        <label>
                            <CheckBox
                                checked={autoDisableStatus}
                                onClick={(event) => event.stopPropagation()}
                                onChange={() => this.autoDisableLogin()}
                            />
                            <span className={styles['disabled-text']}>{__('自动禁用')}</span>
                        </label>
                        <div className={styles['auto-checkbox']}>
                            <Select
                                className={styles['text-middle']}
                                width={200}
                                value={autoDisableDays}
                                disabled={!autoDisableStatus}
                                onChange={(value) => this.selectAutoDisableDays(value)}
                            >
                                <Select.Option value={DisableTime.OneMonth} selected={DisableTime.OneMonth === autoDisableDays}>{1 + __('个月')}</Select.Option>
                                <Select.Option value={DisableTime.TreeMonths} selected={DisableTime.TreeMonths === autoDisableDays}>{3 + __('个月')}</Select.Option>
                                <Select.Option value={DisableTime.SixMonths} selected={DisableTime.SixMonths === autoDisableDays}>{6 + __('个月')}</Select.Option>
                                <Select.Option value={DisableTime.TwelveMonths} selected={DisableTime.TwelveMonths === autoDisableDays}>{12 + __('个月')}</Select.Option>
                            </Select>
                            <span>{__('内未登录的用户账号')}</span>
                        </div>
                    </div>
                    {/* <div className={styles['platform-selections']}>
                        <span className={styles['selection-label']}>{__('客户端登录选项：')}</span>
                        <div className={styles['selections']}>
                            <div className={styles['selections-item']}>
                                <label className={styles['item-label']}>
                                    <CheckBox
                                        checked={forbidOstypeInfo[OstypeInfo.Web] || false}
                                        onClick={(event) => event.stopPropagation()}
                                        disabled={clientLockStatus}
                                        onChange={() => this.forbidTerminalLogin(OstypeInfo.Web)}
                                    />
                                    <div className={styles['checked-text']}>{__('禁止网页端登录')}</div>
                                </label>
                                <label className={styles['item-label']}>
                                    <CheckBox
                                        checked={forbidOstypeInfo[OstypeInfo.Windows] || false}
                                        onClick={(event) => event.stopPropagation()}
                                        disabled={clientLockStatus}
                                        onChange={() => this.forbidTerminalLogin(OstypeInfo.Windows)}
                                    />
                                    <div className={styles['checked-text']}>{__('禁止Windows客户端登录')}</div>
                                </label>
                            </div>
                            <div className={styles['selections-item']}>
                                {
                                    !disabledEndpointTypes.includes(EndpointType.IOS) && (
                                        <label className={styles['item-label']}>
                                            <CheckBox
                                                checked={forbidOstypeInfo[OstypeInfo.iOS] || false}
                                                onClick={(event) => event.stopPropagation()}
                                                disabled={clientLockStatus}
                                                onChange={() => this.forbidTerminalLogin(OstypeInfo.iOS)}
                                            />
                                            <div className={styles['checked-text']}>{__('禁止iOS客户端登录')}</div>
                                        </label>
                                    )
                                }
                                <label className={styles['item-label']}>
                                    <CheckBox
                                        checked={forbidOstypeInfo[OstypeInfo.Mac] || false}
                                        onClick={(event) => event.stopPropagation()}
                                        disabled={clientLockStatus}
                                        onChange={() => this.forbidTerminalLogin(OstypeInfo.Mac)}
                                    />
                                    <div className={styles['checked-text']}>{__('禁止Mac客户端登录')}</div>
                                </label>
                            </div>
                            <div className={styles['selections-item']}>
                                {
                                    !disabledEndpointTypes.includes(EndpointType.Android) && (
                                        <label className={styles['item-label']}>
                                            <CheckBox
                                                checked={forbidOstypeInfo[OstypeInfo.Android] || false}
                                                onClick={(event) => event.stopPropagation()}
                                                disabled={clientLockStatus}
                                                onChange={() => this.forbidTerminalLogin(OstypeInfo.Android)}
                                            />
                                            <div className={styles['checked-text']}>{__('禁止Android客户端登录')}</div>
                                        </label>
                                    )
                                }
                                <label className={styles['item-label']}>
                                    <CheckBox
                                        checked={forbidOstypeInfo[OstypeInfo.Linux] || false}
                                        onClick={(event) => event.stopPropagation()}
                                        disabled={clientLockStatus}
                                        onChange={() => this.forbidTerminalLogin(OstypeInfo.Linux)}
                                    />
                                    <div className={styles['checked-text']}>{__('禁止Linux客户端登录')}</div>
                                </label>
                            </div>
                        </div>
                    </div> */}
                    <div className={styles['account-selection']}>
                        <span className={styles['account-selection-label']}>{__('账号选项：')}</span>
                        <div className={styles['account-selection-content']}>
                            <label>
                                <CheckBox
                                    checked={idCardStatus}
                                    onClick={(event) => event.stopPropagation()}
                                    onChange={() => this.allowIdCardAccount()}
                                />
                                <div className={styles['checked-text']}>{__('允许身份证号作为账号登录')}</div>
                            </label>
                            <div className={styles['explation']}>
                                <span>{__('勾选后，允许用户使用绑定的身份证号作为账号登录')}</span>
                            </div>
                        </div>
                    </div>
                    <div className={styles['title']}>
                        <span>{__('登录认证：')}</span>
                        <span>{__('设置用户必须通过')}</span>
                        <Select
                            className={styles['text-middle']}
                            width={200}
                            value={loginWay}
                            disabled={multiAutnLocked}
                            onChange={({ detail }) => this.selectLoginWay({ detail })}
                        >
                            {
                                !disabledLoginAuth.includes(LoginAuthType.account) ? (
                                    <Select.Option
                                        value={LoginWay.Account}
                                        selected={loginWay === LoginWay.Account}
                                    >
                                        {__('账号密码')}
                                    </Select.Option>
                                ) : <></>
                            }
                            {
                                !disabledLoginAuth.includes(LoginAuthType.accountAndImageCaptcha) ? (
                                    <Select.Option
                                        value={LoginWay.AccountAndImgCaptcha}
                                        selected={loginWay === LoginWay.AccountAndImgCaptcha}
                                    >
                                        {__('账号密码 + 图形验证码')}
                                    </Select.Option>
                                ) : <></>
                            }
                            {
                                !disabledLoginAuth.includes(LoginAuthType.accountAndSMSCaptcha) ? (
                                    <Select.Option
                                        value={LoginWay.AccountAndSmsCaptcha}
                                        selected={loginWay === LoginWay.AccountAndSmsCaptcha}
                                    >
                                        {__('账号密码 + 短信验证码')}
                                    </Select.Option>
                                ) : <></>
                            }
                            {
                                !disabledLoginAuth.includes(LoginAuthType.accountAndDynamicPassword) ? (
                                    <Select.Option
                                        value={LoginWay.AccountAndDynamicPassword}
                                        selected={loginWay === LoginWay.AccountAndDynamicPassword}
                                    >
                                        {__('账号密码 + 动态密码')}
                                    </Select.Option>
                                ) : <></>
                            }
                        </Select>
                        <span >{__('登录')}</span>
                    </div>
                    {
                        loginWay === LoginWay.AccountAndImgCaptcha ? (
                            <div className={styles['explation']}>
                                <span >
                                    {
                                        __('用户登录时，连续输错密码 ')
                                    }
                                </span>
                                <div className={styles['number-input']}>
                                    <ValidateNumber
                                        width={160}
                                        step={1}
                                        precision={0}
                                        max={passwdErrCnt ? 99 : null}
                                        min={0}
                                        maxLength={2}
                                        placeholder={__('请输入0-99的数值')}
                                        value={passwdErrCnt}
                                        validateState={passwdErrCntValidity}
                                        validateMessages={ValidateMessages}
                                        disabled={multiAutnLocked}
                                        onBlur={(event, value) => this.handleBlurPwdErrCntValidate(value)}
                                        onValueChange={({ detail }) => this.changePwdErrCnt(detail)}
                                    />
                                </div>

                                <span>{__(' 次，则开启登录验证码')}
                                </span>
                            </div>
                        ) : null
                    }
                    {
                        loginWay === LoginWay.AccountAndSmsCaptcha ? (
                            <div className={styles['explation']}>
                                {__('(短信验证码需要配置对应的服务器插件才能获取，如果您还没有配置，可以在 ')}
                                {
                                    this.triSystemStatus ?
                                        <span>{__('第三方认证')}</span> :
                                        <span className={styles['link']} onClick={() => this.handleJump()}>{__('第三方认证')}</span>
                                }
                                {__(' 页面进行操作)')}
                            </div>
                        ) : null
                    }
                    {
                        loginWay === LoginWay.AccountAndDynamicPassword ? (
                            <div className={styles['explation']}>
                                {__('(动态密码需要配置对应的应用插件才能获取，如果您还没有配置，可以在 ')}
                                {
                                    this.triSystemStatus ?
                                        <span>{__('第三方认证')}</span> :
                                        <span className={styles['link']} onClick={() => this.handleJump()}>{__('第三方认证')}</span>
                                }
                                {__(' 页面进行操作)')}
                            </div>
                        ) : null
                    }
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
}