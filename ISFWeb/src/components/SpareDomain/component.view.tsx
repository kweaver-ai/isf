import * as React from 'react'
import { noop } from 'lodash'
import classnames from 'classnames'
import { Form, Button, InlineButton, UIIcon } from '@/ui/ui.desktop'
import { ValidateBox, ValidateNumber, CheckBox } from '@/sweet-ui'
import { VerifyType } from '../DomainManage/DomainItem/DomainInfo/component.base'
import { ValidateMessages, ValidateStatus } from '../DomainManage/helper'
import SpareDomainBase from './component.base'
import styles from './styles.view';
import __ from './locale'

const KeyValue = [__('一'), __('二'), __('三'), __('四'), __('五')]

const DomainComponent = ({
    domain,
    index,
    editingIndex,
    validateState,
    onDeleteDomain = noop,
    onClearPassword = noop,
    onChangeSSL = noop,
    onEditDomain = noop,
    onTestDomain = noop,
    onCancelEdit = noop,
    onVerifyParams = noop,
}) => {
    const { editable, address, useSSL, port, account, password, accountError } = domain
    // 如果有正在编辑的数据，那么其他的数据都变成disabled
    const disabled = editingIndex !== -1 && editingIndex !== index

    return (
        <div className={styles['domain-area']}>
            <div className={styles['text-title']}>
                {__('备用域控制器') + __('（') + KeyValue[index] + __('）')}
            </div>
            <span className={styles['text-link']}>
                <InlineButton
                    role={'ui-inlinebutton'}
                    disabled={disabled}
                    code={'\uf000'}
                    title={__('删除')}
                    onClick={() => onDeleteDomain(index)}
                />
            </span>
            <Form role={'ui-form'}>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>
                        {__('域控制器地址：')}
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <span className={styles['warning-style']}>*</span>
                        <ValidateBox
                            role={'sweetui-validatebox'}
                            width={290}
                            disabled={disabled}
                            value={address}
                            validateState={editable ? validateState.address : ValidateStatus.Normal}
                            validateMessages={ValidateMessages}
                            onValueChange={({ detail }) => onEditDomain(index, 'address', detail)}
                            onBlur={() => onVerifyParams(index, VerifyType.DomainIP)}
                        />
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <span className={styles['placeholder']}></span>
                        <CheckBox
                            role={'sweetui-checkbox'}
                            className={styles['block']}
                            checked={useSSL}
                            disabled={disabled}
                            onCheckedChange={({ detail }) => onChangeSSL(detail, index)}
                        >
                            {__('使用SSL')}
                        </CheckBox>
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>
                        {__('域控制器端口：')}
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <span className={styles['warning-style']}>*</span>
                        <ValidateNumber
                            role={'sweetui-validatenumber'}
                            width={290}
                            precision={0}
                            disabled={disabled}
                            value={port}
                            validateState={editable ? validateState.port : ValidateStatus.Normal}
                            validateMessages={ValidateMessages}
                            onValueChange={({ detail }) => onEditDomain(index, 'port', detail)}
                            onBlur={() => onVerifyParams(index, VerifyType.DomainPort)}
                        />
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>
                        {__('域管理员账号：')}
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <span className={styles['warning-style']}>*</span>
                        <ValidateBox
                            role={'sweetui-validatebox'}
                            width={290}
                            disabled={disabled}
                            value={account}
                            validateState={editable ? validateState.account : ValidateStatus.Normal}
                            validateMessages={ValidateMessages}
                            onValueChange={({ detail }) => onEditDomain(index, 'account', detail)}
                            onBlur={() => onVerifyParams(index, VerifyType.Account)}
                        />
                    </Form.Field>
                </Form.Row>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>
                        {__('域管理员密码：')}
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <span className={styles['warning-style']}>*</span>
                        <input className={styles['hidden-input']} type="text" />
                        <ValidateBox
                            role={'sweetui-validatebox'}
                            className={password ? styles['secret-password'] : null}
                            type={'password'}
                            width={290}
                            disabled={disabled}
                            value={password}
                            validateState={editable ? validateState.password : ValidateStatus.Normal}
                            validateMessages={ValidateMessages}
                            onValueChange={({ detail }) => onEditDomain(index, 'password', detail)}
                            onBlur={() => onVerifyParams(index, VerifyType.Password)}
                            onClick={() => onClearPassword(index)}
                        />
                    </Form.Field>
                </Form.Row>
                {
                    accountError && (
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label role={'ui-form.label'}>
                            </Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <span className={styles['placeholder']}></span>
                                <span className={styles['err']}>
                                    {__('账号或密码不正确。')}
                                </span>
                            </Form.Field>
                        </Form.Row>
                    )
                }
                {
                    editable && (
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label role={'ui-form.label'}>
                            </Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <span className={styles['placeholder']}></span>
                                <Button
                                    role={'ui-button'}
                                    key={'confirm-btn'}
                                    className={styles['confirm-btn']}
                                    width={80}
                                    onClick={() => onTestDomain(index)}
                                >
                                    {__('保存')}
                                </Button>
                                <Button
                                    role={'ui-button'}
                                    key={'cancel-btn'}
                                    className={styles['cancel-btn']}
                                    width={80}
                                    onClick={() => onCancelEdit(index)}
                                >
                                    {__('取消')}
                                </Button>
                            </Form.Field>
                        </Form.Row>
                    )
                }
            </Form>
        </div>
    )
}

export default class SpareDomain extends SpareDomainBase {
    render() {
        const { validateState, domains, editingIndex } = this.state

        return (
            <div>
                <div className={styles['tip-text']}>{__('在主域控制器无法正常运行时，备用域控制器可进行用户登录验证、域导入、域同步和反向同步。（最多可添加5个备用域控制器）')}</div>
                {
                    !!domains.length && (
                        domains.map((domain, index) => (
                            <DomainComponent
                                key={index}
                                validateState={validateState}
                                domain={domain}
                                index={index}
                                editingIndex={editingIndex}
                                onDeleteDomain={this.deleteDomain.bind(this)}
                                onClearPassword={this.clearPassword.bind(this)}
                                onChangeSSL={this.changeSSL.bind(this)}
                                onTestDomain={this.clickTestDomain.bind(this)}
                                onEditDomain={this.editDomain.bind(this)}
                                onCancelEdit={this.cancelEdit.bind(this)}
                                onVerifyParams={this.verifyDomainInfo.bind(this)}
                            />
                        ))
                    )
                }
                {
                    domains.length !== 5 && (
                        <div
                            className={classnames(styles['add-domain'], { [styles['disable']]: editingIndex !== -1 })}
                            onClick={editingIndex === -1 ? this.addDomain.bind(this) : null}
                        >
                            <UIIcon
                                role={'ui-uiicon'}
                                className={styles['plus']}
                                code={'\uf089'}
                                size={16}
                            />
                            <span>{__('添加备用域控制器')}</span>
                        </div>
                    )
                }
            </div>
        )
    }
}