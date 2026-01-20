import * as React from 'react';
import { ProtLevel } from '@/core/systemprotectionlevel';
import Dialog from '@/ui/Dialog2/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import Form from '@/ui/Form/ui.desktop';
import Button from '@/ui/Button/ui.desktop';
import { Select } from '@/sweet-ui';
import CheckBoxOption from '@/ui/CheckBoxOption/ui.desktop';
import TextBox from '@/ui/TextBox/ui.desktop';
import ConfirmDialog from '@/ui/ConfirmDialog/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { positiveIntegerAndMaxLength } from '@/util/validators';
import SecurityIntegrationBase from './component.base';
import { PwdValidity, Status, getProtLevels, ProtLevelText } from './helper';
import __ from './locale';
import styles from './styles';
import CustomSecurity from './CustomSecurity';

export default class SecurityIntegration extends SecurityIntegrationBase {

    render() {
        const {
            pwdPolicy: {
                weak_pwd_disabled,
            },
            lockStatus,
            passwdErrCnt,
            passwdLockTime,
            validateLockTimeStatus,
            validatePwdStatus,
            timeCtrl,
            strongPassword,
            isShowProtLevel,
            protLevel,
            loginStrategyStatus,
        } = this.state

        return (
            <div className={styles.container}>
                {
                    this.state.init ?
                        <Dialog
                            role={'ui-dialog2'}
                            width={850}
                            title={__('初始化配置')}
                            buttons={[]}
                        >
                            {
                                this.state.openCustomSecurity ?
                                    this.CustomSecurityClassification()
                                    : null
                            }
                            {
                                this.state.isConfirmCustSecu ?
                                    this.ConfirmCustomSecu()
                                    : null
                            }
                            <Panel role={'ui-panel'}>
                                <Panel.Main role={'ui-panel.main'}>
                                    {
                                        this.state.supportCsf ?
                                            <div>
                                                <div className={styles['info-content-header']}>
                                                    <UIIcon
                                                        role={'ui-uiicon'}
                                                        size={16}
                                                        code={'\uf016'}
                                                    />
                                                    <label className={styles['info-header-title']}>{__('用户密级策略')}</label>
                                                </div>
                                                <Form role={'ui-form'}>
                                                    <Form.Row role={'ui-form.row'}>
                                                        <Form.Label role={'ui-form.label'}>
                                                            <div className={styles['security-label']}>{__('密级列表：')}</div>
                                                        </Form.Label>
                                                        <Form.Field role={'ui-form.field'}>
                                                            <Button
                                                                role={'ui-button'}
                                                                onClick={this.customSecurityClassification.bind(this)}
                                                                width={120}
                                                            >
                                                                {__('设置')}
                                                            </Button>
                                                        </Form.Field>
                                                    </Form.Row>
                                                </Form>
                                                {
                                                    !this.state.isCustomed ?
                                                        <div className={styles['security-tip']}>{__('在保存初始化配置之前，请先完成“自定义密级”的设置')}</div> : null
                                                }
                                            </div>
                                            : null
                                    }
                                    {/* 系统保护等级 */}
                                    {
                                        isShowProtLevel
                                            ? (
                                                <div>
                                                    <div className={styles['info-content-header']}>
                                                        <UIIcon
                                                            role={'ui-uiicon'}
                                                            size={16}
                                                            code={'\uf016'}
                                                        />
                                                        <label className={styles['info-header-title']}>{__('系统保护等级')}</label>
                                                    </div>
                                                    <Form role={'ui-form'} className={styles['prot-level-select']}>
                                                        <Form.Row role={'ui-form.row'}>
                                                            <Form.Label role={'ui-form.label'}>
                                                                <div className={styles['security-label']}>{__('系统保护等级：')}</div>
                                                            </Form.Label >
                                                            <Form.Field role={'ui-form.field'}>
                                                                <div className={styles['dropbox-first']}>
                                                                    {
                                                                        <Select
                                                                            role={'sweetui-select'}
                                                                            onChange={({ detail }) => this.selectProLevel(detail)}
                                                                            value={protLevel}
                                                                        >
                                                                            {
                                                                                getProtLevels(ProtLevel.Classified).map((leve) => {
                                                                                    return (
                                                                                        <Select.Option
                                                                                            role={'sweetui-select.option'}
                                                                                            key={leve}
                                                                                            value={leve}
                                                                                            selected={leve === protLevel}
                                                                                        >
                                                                                            {ProtLevelText[leve]}
                                                                                        </Select.Option>
                                                                                    )
                                                                                })
                                                                            }
                                                                        </Select>
                                                                    }
                                                                </div>
                                                            </Form.Field>
                                                        </Form.Row>
                                                    </Form>
                                                </div>
                                            )
                                            : null
                                    }

                                    <div className={styles['info-content-header']}>
                                        <UIIcon
                                            role={'ui-uiicon'}
                                            size={16}
                                            code={'\uf016'}
                                        />
                                        <label className={styles['info-header-title']}>{__('密码策略')}</label>
                                    </div>
                                    <Form role={'ui-form'}>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label role={'ui-form.label'}>
                                                <div className={styles['security-label']}>{__('密码有效期：')}</div>
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <div className={styles['dropbox-second']}>
                                                    {
                                                        <Select
                                                            role={'sweetui-select'}
                                                            onChange={this.selectPasswordExpiration.bind(this)}
                                                            value={this.state.expireTime}
                                                        >
                                                            {
                                                                this.expireTimeSelected(timeCtrl).map((item) => {
                                                                    return (
                                                                        <Select.Option
                                                                            key={item}
                                                                            value={item}
                                                                            selected={this.state.expireTime === item}>
                                                                            {this.formateExpireTime(item)}
                                                                        </Select.Option>
                                                                    )
                                                                })
                                                            }
                                                        </Select>
                                                    }
                                                </div>
                                            </Form.Field>
                                        </Form.Row>
                                    </Form>
                                    <div className={styles['password-comment']}>{__('密码仅在指定时间段内有效，若超过该有效期，则需要修改密码，否则无法登录')}</div>
                                    <Form role={'ui-form'}>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label role={'ui-form.label'}>
                                                <div className={styles['security-label']}>{__('密码强度：')}</div>
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                {
                                                    weak_pwd_disabled ? (
                                                        <Select
                                                            role={'sweetui-select'}
                                                            value={strongPassword ? 1 : 0}
                                                            onChange={this.selectPasswordStrength.bind(this)}
                                                        >
                                                            <Select.Option value={1} selected={strongPassword} role={'sweetui-select.option'}>{__('强密码')}</Select.Option>
                                                        </Select>
                                                    ) : (
                                                        <Select
                                                            role={'sweetui-select'}
                                                            value={strongPassword ? 1 : 0}
                                                            onChange={this.selectPasswordStrength.bind(this)}
                                                        >
                                                            <Select.Option value={1} selected={strongPassword} role={'sweetui-select.option'}>{__('强密码')}</Select.Option>
                                                            <Select.Option value={0} selected={!strongPassword} role={'sweetui-select.option'}>{__('弱密码')}</Select.Option>
                                                        </Select>
                                                    )
                                                }
                                            </Form.Field>
                                        </Form.Row>
                                    </Form>
                                    {

                                        strongPassword ?
                                            <div>
                                                <div className={styles['comment-start']}>{__('强密码格式：密码长度至少为')}</div>
                                                <div className={styles['inline-layout']}>
                                                    <TextBox
                                                        role={'ui-textbox'}
                                                        width={50}
                                                        value={this.state.passwordLength}
                                                        onChange={(pass) => this.handlePwdCntChange(pass)}
                                                        validator={positiveIntegerAndMaxLength(2)}
                                                        disabled={!strongPassword}
                                                    />
                                                    <div className={styles['operate-btn']}>
                                                        <UIIcon
                                                            role={'ui-uiicon'}
                                                            className={styles['triangle-up']}
                                                            size={16}
                                                            code={'\uf019'}
                                                            onClick={this.passwordLengthAdd.bind(this)}
                                                        />
                                                        <UIIcon
                                                            role={'ui-uiicon'}
                                                            className={styles['triangle-down']}
                                                            size={16}
                                                            code={'\uf01A'}
                                                            onClick={this.passwordLengthSub.bind(this)}
                                                        />
                                                    </div>
                                                </div>
                                                <label className={styles['comment-after']}>{__('个字符，需同时包含 大小写英文字母、数字与特殊字符')}</label>
                                                {
                                                    this.renderValidateError(this.state.validatePwdLenStatus)
                                                }
                                            </div>
                                            :
                                            <div className={styles['password-comment']}>{__('弱密码格式：密码长度至少为6个字符')}</div>
                                    }
                                    <div className={styles['pass-lock']}>
                                        <CheckBoxOption
                                            role={'ui-checkboxoption'}
                                            onChange={this.setPasswordLock.bind(this)}
                                            checked={this.state.lockStatus}
                                        >
                                            {__('启用密码错误锁定：')}
                                        </CheckBoxOption>
                                    </div>
                                    <div>
                                        <div className={styles['comment-start']}>{__('用户在任一情况下，密码连续输错')}</div>
                                        <div className={styles['inline-layout']}>
                                            <TextBox
                                                role={'ui-textbox'}
                                                width={50}
                                                value={passwdErrCnt}
                                                onChange={(cnt) => this.handleErrCntChange(cnt)}
                                                validator={positiveIntegerAndMaxLength(2)}
                                                disabled={!lockStatus}
                                            />
                                            <div className={styles['operate-btn']}>
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    className={styles['triangle-up']}
                                                    size={16}
                                                    code={'\uf019'}
                                                    onClick={this.errCntAdd.bind(this)}
                                                    disabled={!lockStatus}
                                                />
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    className={styles['triangle-down']}
                                                    size={16}
                                                    code={'\uf01A'}
                                                    onClick={this.errCntSub.bind(this)}
                                                    disabled={!lockStatus}
                                                />
                                            </div>
                                        </div>
                                        <label className={styles['comment-after']}>{__('次，则账号将被锁定')}</label>
                                        {
                                            this.renderValidateError(validatePwdStatus)
                                        }
                                    </div>
                                    <div>
                                        <div className={styles['comment-start']}>{__('账号被锁定后，')}</div>
                                        <div className={styles['inline-layout']}>
                                            <TextBox
                                                role={'ui-textbox'}
                                                width={50}
                                                value={passwdLockTime}
                                                onChange={this.changePasswdLockTime}
                                                validator={positiveIntegerAndMaxLength(3)}
                                                disabled={!lockStatus}
                                            />
                                            <div className={styles['operate-btn']}>
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    className={styles['triangle-up']}
                                                    size={16}
                                                    code={'\uf019'}
                                                    onClick={this.addPasswdLockTime}
                                                    disabled={!lockStatus}
                                                />
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    className={styles['triangle-down']}
                                                    size={16}
                                                    code={'\uf01A'}
                                                    onClick={this.subPasswdLockTime}
                                                    disabled={!lockStatus}
                                                />
                                            </div>
                                        </div>
                                        <label className={styles['comment-after']}>{__('分钟后自动解锁或由管理员解锁')}</label>
                                        {
                                            this.renderValidateError(validateLockTimeStatus)
                                        }
                                    </div>
                                </Panel.Main>
                                <Panel.Footer role={'ui-panel.footer'}>
                                    <Panel.Button
                                        theme='oem'
                                        role={'ui-panel.button'}
                                        onClick={this.triggerConfirmInit.bind(this)}
                                        disabled={
                                            !this.state.isCustomed ||
                                            this.state.validatePwdStatus === Status.COUNT_RANGE_ERROR ||
                                            this.state.validatePwdLenStatus === Status.COUNT_PWD_RANGE_ERROR ||
                                            this.state.validateLockTimeStatus === Status.PASSWD_LOCK_TIME_ERROR
                                        }
                                    >
                                        {__('确定')}
                                    </Panel.Button>
                                </Panel.Footer>
                            </Panel>
                        </Dialog>
                        : null
                }
            </div>
        )
    }

    /**
     * 自定义密级对话框
     */
    CustomSecurityClassification() {
        return <CustomSecurity tempCsfLevel={this.state.isClickCustomSecur ? this.state.userCsfLevel : this.state.tempCsfLevel} triggerConfirmCust={this.triggerConfirmCust.bind(this)} closeCustomSecurity={this.closeCustomSecurity.bind(this)}/>
    }

    /**
     * 确认自定义密级
     */
    ConfirmCustomSecu() {
        return (
            <ConfirmDialog onConfirm={this.setCustomedSecurity.bind(this)} onCancel={this.cancelCustom.bind(this)} role={'ui-confirmdialog'}>
                <div className={styles['confirm-msg']}>{__('初始化后将无法更改已设置的用户密级，请确认你的操作。')}</div>
            </ConfirmDialog>
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

    renderValidateError(status) {
        const { pwdPolicy: { max_err_count, min_strong_pwd_length } } = this.state

        switch (status) {
            case Status.OK:
                return null

            case Status.COUNT_PWD_RANGE_ERROR:
                return <div className={styles['passerrcnt']}>{__('强密码长度至少为${num}个字符，请重新输入。', { num: min_strong_pwd_length })}</div>

            case Status.COUNT_RANGE_ERROR:
                return <div className={styles['errcnt']}>{__('密码错误次数范围为1~${num}，请重新输入。', { num: max_err_count })}</div>

            case Status.PASSWD_LOCK_TIME_ERROR:
                return <div className={styles['errcnt']}>{__('账号锁定时间范围为10~180分钟，请重新输入。')}</div>

            case Status.FORBIDDEN_SPECIAL_CHARACTER:
                return <div className={styles['errmsg']}>{__('不能包含 / : * ? " < > | 特殊字符，请重新输入。')}</div>
        }
    }
}