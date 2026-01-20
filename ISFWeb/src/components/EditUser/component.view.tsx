import * as React from 'react';
import { Text, Form, UIIcon, Button } from '@/ui/ui.desktop';
import { Select, ModalDialog2, SweetIcon, ValidateBox, ComboArea } from '@/sweet-ui';
import { ValidateMessages, Type, UserInfoType } from '@/core/user';
import IDNumberBox from '../IDNumberBox/component.view';
import ValidityBox2 from '../ValidityBox2/component.view';
import EditUserBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import OrgAndAccountPick from '../OrgAndAccountPick/component.view';
import session from '@/util/session';
import { NodeType } from '@/core/organization';
import { TabType } from '../OrgAndAccountPick/helper';
import { Doclibs } from '@/core/doclibs/doclibs';

export default class EditUser extends EditUserBase {
    /*
    * 密级
    */
    renderCsfLevel() {
        const { userInfo: { csfLevel, csfLevel2, show_csf_level2 }, csfOptions, csfOptions2 } = this.state;

        return (
            <>
                <Form.Row role={'ui-form.row'}>
                    <Form.Label role={'ui-form.label'}>
                        {
                            __('用户密级：')
                        }
                    </Form.Label>
                    <Form.Field role={'ui-form.field'}>
                        <Select
                            role={'sweetui-select'}
                            width={320}
                            value={csfLevel}
                            disabled={this.triSystemStatus && this.isAdmin}
                            onChange={({ detail }) => this.updateCsfLevel('csfLevel', detail)}
                        >
                            {
                                csfOptions.map((secret) => (
                                    <Select.Option
                                        role={'sweetui-select.option'}
                                        value={secret.value}
                                        key={secret.value}
                                        selected={csfLevel === secret.value}
                                    >
                                        {secret.name}
                                    </Select.Option>
                                ))
                            }
                        </Select>
                    </Form.Field>
                </Form.Row>
                {
                    show_csf_level2 ?
                     <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'}>
                            {
                                __('用户密级2：')
                            }
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <Select
                                role={'sweetui-select'}
                                width={320}
                                value={csfLevel2}
                                disabled={this.triSystemStatus && this.isAdmin}
                                onChange={({ detail }) => this.updateCsfLevel('csfLevel2', detail)}
                            >
                                {
                                    csfOptions2.map((secret) => (
                                        <Select.Option
                                            role={'sweetui-select.option'}
                                            value={secret.value}
                                            key={secret.value}
                                            selected={csfLevel2 === secret.value}
                                        >
                                            {secret.name}
                                        </Select.Option>
                                    ))
                                }
                            </Select>
                        </Form.Field>
                    </Form.Row> : null
                }
            </>
        )
    }

    render() {
        const { showEditDialog, showAddDirectSupervisorDialog, managerInfo, userInfo, userInfo: { loginName, displayName, code, position, remark, certification, email, telNum, idCard, expireTime }, isIDNumEdit, validateState } = this.state;

        return (
            <div>
                {
                    showEditDialog ?
                        <ModalDialog2
                            role={'sweetui-modaldialog2'}
                            title={__('编辑用户')}
                            width={558}
                            zIndex={18}
                            icons={[
                                {
                                    icon: <SweetIcon name="x" size={16} role={'sweetui-sweeticon'} />,
                                    onClick: this.props.onRequestCancel,
                                },
                            ]}
                            buttons={[
                                {
                                    text: __('确定'),
                                    theme: 'oem',
                                    onClick: this.editUser,
                                },
                                {
                                    text: __('取消'),
                                    theme: 'regular',
                                    onClick: this.props.onRequestCancel,
                                },
                            ]}
                        >
                            {
                                this.isSecurit ?
                                    <Form role={'ui-form'}>
                                        {this.renderCsfLevel()}
                                    </Form>
                                    :
                                    <Form role={'ui-form'}>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('用户名：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <Text
                                                    className={styles['unable-edit']}
                                                    role={'ui-text'}
                                                >
                                                    {loginName}
                                                </Text>
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('显示名：')
                                                }
                                                <span className={styles['required']}>*</span>
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <ValidateBox
                                                    role={'sweetui-validatebox'}
                                                    width={320}
                                                    value={displayName}
                                                    disabled={false}
                                                    validateMessages={ValidateMessages}
                                                    validateState={validateState.displayName}
                                                    onBlur={() => this.handleOnBlur(Type.DisplayName)}
                                                    onValueChange={({ detail }) => this.handleValueChange({ displayName: detail })}
                                                />
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('用户编码：')
                                                }
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    code={'\uf055'}
                                                    size={'16px'}
                                                    title={
                                                        __('用户唯一标识，全局唯一，如：工号。')
                                                    }
                                                    color={'#999'}
                                                />
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <ValidateBox
                                                    role={'sweetui-validatebox'}
                                                    width={320}
                                                    value={code}
                                                    disabled={false}
                                                    validateMessages={ValidateMessages}
                                                    validateState={validateState.code}
                                                    onBlur={() => this.handleOnBlur(Type.Code)}
                                                    onValueChange={({ detail }) => this.handleValueChange({ code: detail })}
                                                />
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('直属上级：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <div style={{ display: 'flex' }}>
                                                    <ComboArea
                                                        role={'sweetui-validatecomboarea'}
                                                        width={320}
                                                        height={30}
                                                        required={false}
                                                        uneditable={true}
                                                        placeholder={__('请添加用户')}
                                                        value={managerInfo}
                                                        formatter={this.formatOwner}
                                                        onChange={this.changeOwner}
                                                        style={{ overflowY: 'hidden' }}
                                                        tagClassName={styles['tag-style']}
                                                    />
                                                    <Button
                                                        size={'auto'}
                                                        theme={'regular'}
                                                        onClick={() => this.setState({ showAddDirectSupervisorDialog: true })}
                                                        style={{ marginLeft: 8 }}
                                                    >
                                                        {__('选择')}
                                                    </Button>
                                                </div>
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('岗位：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <ValidateBox
                                                    role={'sweetui-validatebox'}
                                                    width={320}
                                                    value={position}
                                                    disabled={false}
                                                    validateMessages={ValidateMessages}
                                                    validateState={validateState.position}
                                                    onBlur={() => this.handleOnBlur(Type.Position)}
                                                    onValueChange={({ detail }) => this.handleValueChange({ position: detail })}
                                                />
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('备注：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <ValidateBox
                                                    role={'sweetui-validatebox'}
                                                    width={320}
                                                    value={remark}
                                                    disabled={false}
                                                    validateMessages={ValidateMessages}
                                                    validateState={validateState.remark}
                                                    onBlur={() => this.handleOnBlur(Type.Remark)}
                                                    onValueChange={({ detail }) => this.handleValueChange({ remark: detail })}
                                                />
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('直属部门：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <Text
                                                    className={styles['unable-edit']}
                                                    role={'ui-text'}
                                                >
                                                    {
                                                        this.editInfo && this.editInfo.user &&
                                                        this.editInfo.user.departmentNames.join(', ')
                                                    }
                                                </Text>
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('认证类型：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <Text
                                                    className={styles['unable-edit']}
                                                    role={'ui-text'}
                                                >
                                                    {certification}
                                                </Text>
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('邮箱地址：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <ValidateBox
                                                    role={'sweetui-validatebox'}
                                                    width={320}
                                                    value={email}
                                                    disabled={false}
                                                    validateMessages={ValidateMessages}
                                                    validateState={validateState.email}
                                                    onBlur={() => this.handleOnBlur(Type.Email)}
                                                    onValueChange={({ detail }) => this.handleValueChange({ email: detail })}
                                                />
                                            </Form.Field>
                                        </Form.Row>
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label
                                                className={styles['text-label']}
                                                role={'ui-form.label'}
                                            >
                                                {
                                                    __('手机号：')
                                                }
                                            </Form.Label>
                                            <Form.Field role={'ui-form.field'}>
                                                <ValidateBox
                                                    role={'sweetui-validatebox'}
                                                    width={320}
                                                    value={telNum}
                                                    disabled={false}
                                                    validateMessages={ValidateMessages}
                                                    validateState={validateState.telNum}
                                                    onBlur={() => this.handleOnBlur(Type.TelNumber)}
                                                    onValueChange={({ detail }) => this.handleValueChange({ telNum: detail })}
                                                />
                                            </Form.Field>
                                        </Form.Row>
                                        <div>
                                            <IDNumberBox
                                                width={320}
                                                idcardNumber={idCard ? idCard : ''}
                                                isClickBtn={true}
                                                isIDNumEdit={isIDNumEdit}
                                                defaultidcardNumber={this.editInfo ? this.editInfo.user.idcardNumber : ''}
                                                onChange={(value) => {
                                                    this.setState({
                                                        userInfo: {
                                                            ...userInfo,
                                                            idCard: value,
                                                        },
                                                        isIDNumEdit: false,
                                                        isEditUserInfo: true,
                                                    })
                                                }}
                                            />
                                        </div>
                                        {this.renderCsfLevel()}
                                        <Form.Row role={'ui-form.row'}>
                                            <Form.Label role={'ui-form.label'}>
                                                {
                                                    __('有效期限：')
                                                }
                                            </Form.Label>
                                            {
                                                expireTime ?
                                                    (
                                                        <Form.Field role={'ui-form.field'}>
                                                            <ValidityBox2
                                                                width={320}
                                                                className={styles['expire-date']}
                                                                allowPermanent={true}
                                                                value={expireTime}
                                                                selectRange={[new Date()]}
                                                                onChange={(value) => { this.changeExpireTime(value) }}
                                                            />
                                                        </Form.Field>
                                                    ) : null
                                            }
                                        </Form.Row>
                                       
                                    </Form>
                            }
                        </ModalDialog2>
                        : null
                }
                {
                    showAddDirectSupervisorDialog ?
                        <OrgAndAccountPick
                            title={__('选择用户')}
                            userid={session.get('isf.userid')}
                            selectType={[NodeType.USER]}
                            tabType={[TabType.Org]}
                            isShowDisabledUsers={false}
                            isSingleChoice={true}
                            convererOut={this.convererOutData()}
                            selected={managerInfo}
                            onRequestCancel={() => {
                                this.setState({
                                    showAddDirectSupervisorDialog: false,
                                })
                            }}
                            onRequestConfirm={(data) =>{
                                this.setState({
                                    managerInfo: data as unknown as UserInfoType[],
                                    showAddDirectSupervisorDialog: false,
                                    isEditUserInfo: true,
                                })
                            }}
                        />
                        : null
                }
            </div>
        )
    }

    /**
     * 转出数据时转换数据格式
     */
    private convererOutData = (newType: string = ''): (value: Doclibs.UserInfo) => Doclibs.UserInfo => {
        return ({ id, name, type }): Doclibs.UserInfo => {
            return {
                id,
                name,
                type: newType ? newType : type,
            }
        }
    }
    /**
     * 格式化所有者函数
     */
    private formatOwner = (users: Doclibs.UserInfo): string => {
        return users.name
    }

    /**
     * 更改直属上级
     */
    private changeOwner = (users): void => {
        this.setState({
            managerInfo: users,
            isEditUserInfo: true,
        })
    }
}