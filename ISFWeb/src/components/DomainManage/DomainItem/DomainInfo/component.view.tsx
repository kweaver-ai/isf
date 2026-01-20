import * as React from 'react'
import { Form, Text } from '@/ui/ui.desktop'
import { Radio, ValidateBox, CheckBox, ValidateNumber, Button } from '@/sweet-ui'
import { ActionType, DomainType, ValidateMessages } from '../../helper'
import DomainInfoBase, { VerifyType } from './component.base'
import styles from './styles.view';
import __ from './locale'

export default class DomainInfo extends DomainInfoBase {
    render() {
        const { actionType, selection, editDomain } = this.props;
        const {
            domainInfo: {
                domainType, domainName, domainIP, useSSL, domainPort, account, password, id,
            },
            accountError,
            isDomainInfoEditStatus,
            validateStatus: {
                domainNameValidateStatus, domainIPValidateStatus, domainPortValidateStatus, accountValidateStatus, passwordValidateStatus,
            },
            showDomainNameInput,
        } = this.state;

        return (
            <Form role={'ui-form'}>
                {
                    actionType === ActionType.Add && selection.id ?
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label role={'ui-form.label'}>{__('当前选择的域：')}</Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <Text role={'ui-text'} className={styles['name']}>{selection.name}</Text>
                            </Form.Field>
                        </Form.Row> : null
                }
                {
                    actionType === ActionType.Add && showDomainNameInput ?
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label role={'ui-form.label'}>{__('域类型：')}</Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <Radio
                                    role={'sweetui-radio'}
                                    value={DomainType.Primary}
                                    disabled={selection.type}
                                    checked={domainType === DomainType.Primary}
                                    onChange={() => this.changeDomainType(DomainType.Primary)}
                                >
                                    {__('主域')}
                                </Radio>
                                <Radio
                                    role={'sweetui-radio'}
                                    value={DomainType.Sub}
                                    disabled={!selection.type || selection.type === DomainType.Trust}
                                    checked={domainType === DomainType.Sub}
                                    onChange={() => this.changeDomainType(DomainType.Sub)}
                                >
                                    {__('子域')}
                                </Radio>
                                <Radio
                                    role={'sweetui-radio'}
                                    value={DomainType.Trust}
                                    disabled={!selection.type || selection.type !== DomainType.Primary}
                                    checked={domainType === DomainType.Trust}
                                    onChange={() => this.changeDomainType(DomainType.Trust)}
                                >
                                    {__('信任域')}
                                </Radio>
                            </Form.Field>
                        </Form.Row> :
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label role={'ui-form.label'}>{__('域类型：')}</Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                {
                                    domainType === DomainType.Primary ? __('主域') : domainType === DomainType.Sub ? __('子域') : __('信任域')
                                }
                            </Form.Field>
                        </Form.Row>
                }
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>{__('域名：')}</Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        {
                            showDomainNameInput ?
                                <div className={styles['wrapper']}>
                                    <span className={styles['required']}>*</span>
                                    <ValidateBox
                                        role={'sweetui-validatebox'}
                                        width={320}
                                        type={'text'}
                                        disabled={false}
                                        value={domainName}
                                        onBlur={() => this.verifyDomainInfo(VerifyType.DomainName)}
                                        onValueChange={this.changeDomainName}
                                        validateState={domainNameValidateStatus}
                                        validateMessages={ValidateMessages}
                                    />
                                </div> :
                                <div className={styles['wrapper']}>
                                    <Text role={'ui-text'} className={styles['name']}>{domainName}</Text>
                                </div>
                        }
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'} align={'top'}>
                        <span className={styles['form-label']}>{__('域控制器地址：')}</span>
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <div className={styles['wrapper']}>
                            <span className={styles['required']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                type={'text'}
                                value={domainIP}
                                onValueChange={this.changeDomainIP}
                                onBlur={() => this.verifyDomainInfo(VerifyType.DomainIP)}
                                validateState={domainIPValidateStatus}
                                validateMessages={ValidateMessages}
                            />
                            <CheckBox
                                role={'sweetui-checkbox'}
                                className={styles['block']}
                                checked={useSSL}
                                onCheckedChange={this.changeSSL}
                            >
                                {__('使用SSL')}
                            </CheckBox>
                        </div>
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'} align={'top'}>
                        <span className={styles['form-label']}>{__('域控制器端口：')}</span>
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <div className={styles['wrapper']}>
                            <span className={styles['required']}>*</span>
                            <ValidateNumber
                                role={'sweetui-validatenumber'}
                                width={320}
                                precision={0}
                                value={domainPort}
                                onValueChange={this.changeDomainPort}
                                onBlur={() => this.verifyDomainInfo(VerifyType.DomainPort)}
                                validateState={domainPortValidateStatus}
                                validateMessages={ValidateMessages}
                            />
                        </div>
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>{__('域管理员账号：')}</Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <div className={styles['wrapper']}>
                            <span className={styles['required']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                type={'text'}
                                value={account}
                                onValueChange={this.changeAccount}
                                onBlur={() => this.verifyDomainInfo(VerifyType.Account)}
                                validateState={accountValidateStatus}
                                validateMessages={ValidateMessages}
                            />
                        </div>
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>{__('域管理员密码：')}</Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <div className={styles['wrapper']}>
                            <span className={styles['required']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                type={'password'}
                                value={password}
                                onValueChange={this.changePassword}
                                onBlur={() => this.verifyDomainInfo(VerifyType.Password)}
                                onClick={this.clearPassword}
                                validateState={passwordValidateStatus}
                                validateMessages={ValidateMessages}
                            />
                        </div>
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}></Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        {
                            accountError && <div className={styles['err']}>
                                {__('账号或密码不正确。')}
                            </div>
                        }
                        {
                            actionType === ActionType.Add && selection.id || actionType === ActionType.Edit && editDomain.type !== DomainType.Primary ?
                                null :
                                isDomainInfoEditStatus ?
                                    <div className={styles['btns']}>
                                        <Button role={'sweetui-button'} onClick={this.saveDomainInfo}>{__('保存')}</Button>
                                        <Button
                                            role={'sweetui-button'}
                                            className={styles['cancel']}
                                            disabled={!domainName && !domainIP && !account && !password && domainPort === this.originDomainInfo.domainPort}
                                            onClick={this.cancelDomainInfo}
                                        >
                                            {__('取消')}
                                        </Button>
                                    </div> : null
                        }
                    </Form.Field>
                </Form.Row>
            </Form >
        )
    }
}