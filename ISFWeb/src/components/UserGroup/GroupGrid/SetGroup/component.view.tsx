import * as React from 'react';
import { Form } from '@/ui/ui.desktop';
import { TextArea, ValidateBox, ModalDialog2, SweetIcon, Button, ComboArea } from '@/sweet-ui';
import OrgAndGroupPicker from '../../../OrgAndGroupPicker';
import { TabType } from '../../../OrgAndGroupPick/helper';
import SetGroupBase from './component.base';
import { ValidateMessages } from './helper'
import __ from './locale';
import styles from './styles.view';

export default class SetGroup extends SetGroupBase {
    render() {
        const { editGroup } = this.props
        const { nameStatus, name, notes, userGroupSources, isShowSelectGroup } = this.state

        return (
            <ModalDialog2
                role={'sweetui-modaldialog2'}
                width={542}
                title={editGroup && editGroup.id ? __('编辑用户组') : __('新建用户组')}
                icons={[
                    {
                        icon: <SweetIcon name="x" size={16} />,
                        onClick: () => this.props.onRequestCancel(),
                    },
                ]}
                buttons={[
                    {
                        text: __('确定'),
                        theme: 'oem',
                        onClick: this.confirm,
                    },
                    {
                        text: __('取消'),
                        theme: 'regular',
                        onClick: () => this.props.onRequestCancel(),
                    },
                ]}
            >
                <Form
                    role={'ui-form'}
                    className={styles['form']}
                >
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'}>{__('用户组名称：')}</Form.Label>
                        <Form.Field
                            role={'ui-form.field'}
                            key={'name'}>
                            <div className={styles['box-wrapper']}>
                                <span className={styles['key-dot']}>*</span>
                                <ValidateBox
                                    role={'sweetui-validatebox'}
                                    width={this.isCreate ? 308 : 378}
                                    value={name}
                                    placeholder={__('1-128个字符，不能包含空格 或 \\ / : * ? " < > | 特殊字符')}
                                    validateMessages={ValidateMessages}
                                    validateState={nameStatus}
                                    onValueChange={({ detail }) => this.changeName(detail)}
                                    onBlur={this.verifyName}
                                />
                            </div>
                        </Form.Field>
                    </Form.Row>
                    {
                        this.isCreate && (
                            <Form.Row role={'ui-form.row'}>
                                <Form.Label role={'ui-form.label'} align={'top'}>{__('用户组成员：')}</Form.Label>
                                <Form.Field
                                    role={'ui-form.field'}
                                    key={'name'}>
                                    <div className={styles['box-wrapper']}>
                                        <ComboArea
                                            width={310}
                                            height={96}
                                            placeholder={__('可从已有用户组中添加成员')}
                                            value={userGroupSources}
                                            formatter={({ name }) => name}
                                            onChange={(value) => this.changeSelectGroups(value)}
                                        />
                                        <Button
                                            className={styles['add-button']}
                                            width={'60px'}
                                            onClick={() => this.changeSelectGroupsState(true)}
                                        >
                                            {__('选择')}
                                        </Button>
                                    </div>
                                </Form.Field>
                            </Form.Row>
                        )
                    }
                    <Form.Row>
                        <Form.Label role={'ui-form.label'} align={'top'}><span className={styles['form-lable']}>{__('备注：')}</span></Form.Label>
                        <Form.Field role={'ui-form.field'} key={'notes'}>
                            <TextArea
                                role={'sweetui-textarea'}
                                width={this.isCreate ? 308 : 378}
                                height={122}
                                maxLength={300}
                                value={notes}
                                onValueChange={({ detail }) => this.changeNote(detail)}
                            />
                        </Form.Field>
                    </Form.Row>
                </Form>

                {
                    isShowSelectGroup && (
                        <OrgAndGroupPicker
                            title={__('选择已有用户组')}
                            tabType={[TabType.Group]}
                            isMult={false}
                            defaultSelections={userGroupSources}
                            onRequestConfirm={(selections) => {
                                this.changeSelectGroups(selections)
                                this.changeSelectGroupsState(false)
                            }}
                            onRequestCancel={() => this.changeSelectGroupsState(false)}
                        />
                    )
                }
            </ModalDialog2>
        )
    }
}