import * as React from 'react';
import { Form, ProgressBar, FlexBox, Text } from '@/ui/ui.desktop';
import { ValidateMessages } from '@/core/user';
import { shrinkText } from '@/util/formatters';
import { Select2, ModalDialog2, SweetIcon, CheckBox, ValidateSelect } from '@/sweet-ui';
import ValidityBox2 from '../ValidityBox2/component.view';
import { Range } from '../helper'
import BatchEditUserBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import { Status } from './type'

const selectMapType = (value: Range, depName: string) => {
    const rangeMapType = {
        [Range.USERS]: __('当前选中用户'),
        [Range.DEPARTMENT]: __('${name} 部门成员', { name: depName }),
        [Range.DEPARTMENT_DEEP]: __('${name} 及其子部门成员', { name: depName }),
    }

    return rangeMapType[value]
}

export default class BatchEditUser extends BatchEditUserBase {
    /*
    * 密级
    */
    renderCsfLevel() {
        const { csfLevel, csfLevel2, csfOptions, csfOptions2, show_csf_level2, csfIsChecked, csf2IsChecked, csfValidateState, csf2ValidateState } = this.state;

        return (
            <>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>
                        {
                            this.isSecurit ?
                                <span className={styles['label']}>{__('用户密级：')}</span>
                                :
                                <CheckBox
                                    role={'sweetui-checkbox'}
                                    checked={csfIsChecked}
                                    onCheckedChange={({ detail }) => this.updateCsfIsChecked('csfIsChecked', detail)}
                                >
                                    {__('用户密级：')}
                                </CheckBox>
                        }
                    </Form.Label>
                    <Form.Field
                        role={'ui-form.field'}
                        key={'storageSite'}
                        className={styles['box-wrapper']}>
                        {
                            csfIsChecked ?
                                <span className={styles['key-dot']}>*</span>
                                : null
                        }
                        <div className={styles['select-height']}>
                            <ValidateSelect
                                role={'sweetui-select'}
                                selectorStyle={{ height: '98%' }}
                                width={320}
                                selectorWidth={320}
                                menuWidth={320}
                                disabled={!csfIsChecked}
                                validateState={csfValidateState}
                                validateMessages={ValidateMessages}
                                placeholder={__('请选择密级')}
                                value={csfLevel}
                                onBlur={() => this.handleOnBlur('csfLevel')}
                                onChange={({ detail }) => this.updateCsfLevel('csfLevel', detail)}
                            >
                                {
                                    csfOptions.map((secret) => (
                                        <ValidateSelect.Option
                                            role={'sweetui-select.option'}
                                            value={secret.value}
                                            key={secret.value}
                                            selected={csfLevel === secret.value}
                                        >
                                            {secret.name}
                                        </ValidateSelect.Option>
                                    ))
                                }
                            </ValidateSelect>
                        </div>
                    </Form.Field>
                </Form.Row>
                {
                   show_csf_level2 ?
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'}>
                            {
                                this.isSecurit ?
                                    <span className={styles['label']}>{__('用户密级2：')}</span>
                                    :
                                    <CheckBox
                                        role={'sweetui-checkbox'}
                                        checked={csf2IsChecked}
                                        onCheckedChange={({ detail }) => this.updateCsfIsChecked('csf2IsChecked', detail)}
                                    >
                                        {__('用户密级2：')}
                                    </CheckBox>
                            }
                        </Form.Label>
                        <Form.Field
                            role={'ui-form.field'}
                            key={'storageSite'}
                            className={styles['box-wrapper']}>
                            {
                                csf2IsChecked ?
                                    <span className={styles['key-dot']}>*</span>
                                    : null
                            }
                            <div className={styles['select-height']}>
                                <ValidateSelect
                                    role={'sweetui-select'}
                                    selectorStyle={{ height: '98%' }}
                                    width={320}
                                    selectorWidth={320}
                                    menuWidth={320}
                                    disabled={!csf2IsChecked}
                                    validateState={csf2ValidateState}
                                    validateMessages={ValidateMessages}
                                    placeholder={__('请选择密级')}
                                    value={csfLevel2}
                                    onBlur={() => this.handleOnBlur('csfLevel2')}
                                    onChange={({ detail }) => this.updateCsfLevel('csfLevel2', detail)}
                                >
                                    {
                                        csfOptions2.map((secret) => (
                                            <ValidateSelect.Option
                                                role={'sweetui-select.option'}
                                                value={secret.value}
                                                key={secret.value}
                                                selected={csfLevel2 === secret.value}
                                            >
                                                {secret.name}
                                            </ValidateSelect.Option>
                                        ))
                                    }
                                </ValidateSelect>
                            </div>
                        </Form.Field>
                    </Form.Row> : null
                }
            </>
        )
    }

    /*
    * 有效期限
    */
    renderExpTime() {
        const { expireTime, expIsChecked } = this.state;

        return (
            <Form.Row role={'ui-form.row'}>
                <Form.Label role={'ui-form.label'}>
                    {
                        this.isAdmin ?
                            <span className={styles['label']}>{__('有效期限：')}</span>
                            :
                            <CheckBox
                                role={'sweetui-checkbox'}
                                checked={expIsChecked}
                                onCheckedChange={({ detail }) => this.updateExpIsChecked(detail)}
                            >
                                {__('有效期限：')}
                            </CheckBox>
                    }
                </Form.Label>
                <Form.Field role={'ui-form.field'}>
                    <div className={styles['box-wrapper']}>
                        {
                            expIsChecked ?
                                <span className={styles['key-dot']}>*</span>
                                : null
                        }
                        {
                            expireTime ?
                                (
                                    <ValidityBox2
                                        width={320}
                                        disabled={!expIsChecked}
                                        className={styles['expire-date']}
                                        allowPermanent={true}
                                        value={expireTime}
                                        selectRange={[new Date()]}
                                        onChange={this.changeExpireTime}
                                    />
                                ) : null
                        }
                    </div>
                </Form.Field>
            </Form.Row>
        )
    }

    /*
    * 设置创建
    */
    renderConfig() {
        const { csfIsChecked, expIsChecked, selected } = this.state;
        const { dep, users, onRequestCancel } = this.props;

        return (
            <ModalDialog2
                role={'sweetui-modaldialog2'}
                title={__('批量编辑')}
                icons={[{
                    icon: <SweetIcon role={'sweetui-sweeticon'} name="x" size={16} />,
                    onClick: onRequestCancel,
                },
                ]}
                buttons={[
                    {
                        text: __('确定'),
                        theme: 'oem',
                        disabled: (!expIsChecked && !csfIsChecked) ? true : false,
                        onClick: this.confirmBatchEditUser,
                    },
                    {
                        text: __('取消'),
                        theme: 'regular',
                        onClick: onRequestCancel,
                    },
                ]}
            >
                <div className={styles['container-box']}>
                    {
                        [
                            __(' 您可以为 '),
                            <div className={styles['select-range']} key={'selectRange'}>
                                <Select2
                                    role={'sweetui-select'}
                                    value={selected}
                                    onChange={({ detail }) => this.selectedType(detail)}
                                    selectorWidth={280}
                                    menuWidth={280}
                                >
                                    {
                                        [Range.USERS, Range.DEPARTMENT, Range.DEPARTMENT_DEEP].filter((value) =>
                                            !((dep.id === '-2' || dep.id === '-1') &&
                                                value === Range.DEPARTMENT_DEEP ||
                                                (value === Range.USERS && !users.length)),
                                        ).map((value) => {
                                            return (
                                                <Select2.Option
                                                    role={'sweetui-select.option'}
                                                    key={value}
                                                    value={value}
                                                    selected={selected === value}
                                                >
                                                    {
                                                        selectMapType(value, dep.name)
                                                    }
                                                </Select2.Option>
                                            )
                                        })
                                    }
                                </Select2>
                            </div>,
                            __(' 做以下设置：'),
                        ]
                    }
                    {
                        this.isSecurit ?
                            (
                                <Form role={'ui-form'}>
                                    {this.renderCsfLevel()}
                                    {this.updateExpIsChecked(false)}
                                </Form>
                            )
                            :
                            this.isAdmin ?
                                (
                                    <Form role={'ui-form'}>
                                        {this.renderExpTime()}
                                        {this.updateCsfIsChecked(false)}
                                    </Form>
                                )
                                :
                                (
                                    <Form role={'ui-form'}>
                                        {this.renderCsfLevel()}
                                        {this.renderExpTime()}
                                    </Form>
                                )
                    }
                </div>
            </ModalDialog2>
        )
    }

    /*
    * 设置过程
    */
    renderProgress() {
        const { progress, currentUserName: { loginName, displayName } } = this.state;

        return (
            <div>
                {
                    (displayName && loginName) ?
                        <ModalDialog2
                            role={'sweetui-modaldialog2'}
                            className={styles['sweetui-modaldialog2']}
                            title={__('批量编辑')}
                            icons={[{
                                icon: <SweetIcon role={'sweetui-sweeticon'} name="x" size={16} />,
                                onClick: this.changeStatus,
                            },
                            ]}
                            buttons={[]}
                        >
                            <div className={styles['progress']}>
                                <div className={styles['edit-tip']}>
                                    <Text className={styles['edit-text']}>
                                        {__(`正在设置用户`)}
                                    </Text>
                                    {`“${shrinkText(displayName, { limit: 16, indicator: '...' })}(${shrinkText(loginName, { limit: 16, indicator: '...' })})”：`}
                                </div>
                                <ProgressBar
                                    role={'ui-progressbar'}
                                    value={progress / this.currentUserLength}
                                    width={350}
                                    height={16}
                                    progressBackground={'#9abbef'}
                                />
                            </div>
                        </ModalDialog2>
                        : null
                }
            </div>
        )
    }

    /*
    * 取消确认
    */
    renderConfirm() {
        return (
            <ModalDialog2
                role={'sweetui-modaldialog2'}
                className={styles['sweetui-modaldialog2']}
                icons={[]}
                buttons={[
                    {
                        text: __('确定'),
                        theme: 'oem',
                        onClick: this.confirmCancel,
                    },
                    {
                        text: __('取消'),
                        theme: 'regular',
                        onClick: this.changeStatus,
                    },
                ]}
            >
                <FlexBox role={'ui-flexbox'}>
                    <FlexBox.Item align="left top">
                        <SweetIcon
                            role={'sweetui-sweeticon'}
                            name="notice"
                            size={40}
                            color="#F39422"
                        />
                    </FlexBox.Item>
                    <FlexBox.Item className={styles['container']}>
                        <div className={styles['conform-tip']}>
                            {__('提示')}
                        </div>
                        <div className={styles['conform-container']}>
                            {__('关闭后，将会停止批量编辑用户密级和有效期限，确认要关闭吗？')}
                        </div>
                    </FlexBox.Item>
                </FlexBox>
            </ModalDialog2>
        )
    }

    render() {
        switch (this.state.status) {
            case Status.Config:
                return this.renderConfig()

            case Status.Progress:
                return this.renderProgress()

            case Status.Confirm:
                return this.renderConfirm()

            default:
                return null
        }
    }
}