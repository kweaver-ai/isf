import * as React from 'react';
import { ModalDialog2, SweetIcon, ValidateBox, Switch, ComboArea, ValidateComboArea } from '@/sweet-ui';
import { Button, Form, Text, Title, UIIcon } from '@/ui/ui.desktop';
import { UserInfoType, ValidateMessages, ValidateState } from '@/core/user';
import CreateDepartmentBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import { Doclibs } from '@/core/doclibs/doclibs';
import OrgAndAccountPick from '../OrgAndAccountPick/component.view';
import { TabType } from '../OrgAndAccountPick/helper';
import { NodeType } from '@/core/organization';
import session from '@/util/session';
import OrganizationPicker from '../OrganizationPicker/component.view';

export default class CreateDepartment extends CreateDepartmentBase {

    render() {
        const { validateState, showAddDepartmentLeaderDialog, showAddDepartmentDialog, managerInfo,  departmentInfo: { departName, code, status, remark, parentName, parentId, parentType, email, ossInfo } } = this.state;
        return (
            <div>
                <ModalDialog2
                    role={'sweetui-modaldialog2'}
                    title={__('新建部门')}
                    width={545}
                    zIndex={18}
                    icons={[
                        {
                            icon: <SweetIcon name="x" size={16} role={'sweetui-sweeticon'}/>,
                            onClick: this.props.onRequestCancelCreateDep,
                        },
                    ]}
                    buttons={[
                        {
                            text: __('确定'),
                            theme: 'oem',
                            onClick: this.createDepartment,
                        },
                        {
                            text: __('取消'),
                            theme: 'regular',
                            onClick: this.props.onRequestCancelCreateDep,
                        },
                    ]}
                >
                    <Form role={'ui-form'}>
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label className={styles['left-width']} role={'ui-form.label'}>
                                {
                                    __('部门名称：')
                                }
                                <span className={styles['required']}>*</span>
                            </Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <ValidateBox
                                    role={'sweetui-validatebox'}
                                    width={320}
                                    value={departName}
                                    disabled={false}
                                    validateMessages={ValidateMessages}
                                    validateState={validateState.departName}
                                    onBlur={() => this.handleOnBlurDep()}
                                    onValueChange={({ detail }) => this.handleValueChange({ departName: detail })}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label
                                className={styles['text-label']}
                                role={'ui-form.label'}
                            >
                                {
                                    __('部门编码：')
                                }
                                <UIIcon
                                    role={'ui-uiicon'}
                                    code={'\uf055'}
                                    size={'16px'}
                                    title={
                                        __('部门唯一标识，全局唯一。')
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
                                    onBlur={() => this.handleOnBlurCode()}
                                    onValueChange={({ detail }) => this.handleValueChange({ code: detail })}
                                />
                            </Form.Field>
                        </Form.Row>
                        {
                            this.props.sourcePage !== 'default' ?
                                <Form.Row role={'ui-form.row'}>
                                    <Form.Label className={styles['left-width']} role={'ui-form.label'}>
                                        {
                                            __('上级部门：')
                                        }
                                        <span className={styles['required']}>*</span>
                                    </Form.Label>
                                    <Form.Field role={'ui-form.field'}>
                                        <div style={{ display: 'flex' }}>
                                            <ValidateComboArea
                                                role={'sweetui-validatecomboarea'}
                                                width={320}
                                                height={30}
                                                required={false}
                                                uneditable={true}
                                                placeholder={__('选择部门')}
                                                value={ parentId && parentName ? [{ id: parentId, name: parentName, type: parentType }]: []}
                                                formatter={this.formatDP}
                                                onChange={this.changeDp}
                                                validateState={validateState.parentId}
                                                validateMessages={ValidateMessages}
                                                style={{ overflowY: 'hidden' }}
                                                tagClassName={styles['tag-style']}
                                            />
                                            <Button
                                                size={'auto'}
                                                theme={'regular'}
                                                onClick={() => this.setState({ showAddDepartmentDialog: true })}
                                                style={{ marginLeft: 8 }}
                                            >
                                                {__('选择')}
                                            </Button>
                                        </div>
                                    </Form.Field>
                                </Form.Row> : null
                        }
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label
                                className={styles['text-label']}
                                role={'ui-form.label'}
                            >
                                {
                                    __('部门负责人：')
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
                                        onClick={() => this.setState({ showAddDepartmentLeaderDialog: true })}
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
                                    onBlur={() => this.handleOnBlurRemark()}
                                    onValueChange={({ detail }) => this.handleValueChange({ remark: detail })}
                                />
                            </Form.Field>
                        </Form.Row>
                        {
                            this.props.sourcePage === 'default' ?
                                <Form.Row role={'ui-form.row'}>
                                    <Form.Label className={styles['left-width']} role={'ui-form.label'}>
                                        {
                                            __('上级部门：')
                                        }
                                    </Form.Label>
                                    <Form.Field role={'ui-form.field'}>
                                        <Text className={styles['unable-edit']} role={'ui-text'}>
                                            {parentName}
                                        </Text>
                                    </Form.Field>
                                </Form.Row> : null
                        }
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label className={styles['left-width']} role={'ui-form.label'}>
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
                                    onBlur={() => this.handleOnBlurEmail()}
                                    onValueChange={({ detail }) => this.handleValueChange({ email: detail })}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label className={styles['left-width']} role={'ui-form.label'}>
                                {
                                    __('状态：')
                                }
                            </Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <div className={styles['mark']}></div>
                                <div className={styles['switch']}>
                                    <Title content={__(`点此${status ? '禁用' : '启用'}部门`)} role={'sweetui-title'}>
                                        <Switch
                                            checked={status}
                                            onChange={({ detail }) => this.handleValueChange({ status: detail })}
                                        />
                                    </Title>
                                </div>
                            </Form.Field>
                        </Form.Row>
                    </Form>
                </ModalDialog2>
                {
                    showAddDepartmentLeaderDialog ?
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
                                    showAddDepartmentLeaderDialog: false,
                                })
                            }}
                            onRequestConfirm={(data) =>{
                                this.setState({
                                    managerInfo: data as unknown as UserInfoType[],
                                    showAddDepartmentLeaderDialog: false,
                                })
                            }}
                        />
                        : null
                }
                {
                    showAddDepartmentDialog ?
                        <OrganizationPicker
                            title={__('选择部门')}
                            isSingleChoice={true}
                            userid={session.get('isf.userid')}
                            selectType={[NodeType.ORGANIZATION, NodeType.DEPARTMENT]}
                            data={parentId && parentName ? [{ id: parentId, name: parentName, type: parentType }]: []}
                            convererOut={this.convererOutData()}
                            onConfirm={(data) => {
                                this.setState({
                                    departmentInfo: { ...this.state.departmentInfo, parentId: data[0].id, parentName: data[0].name, parentType: data[0].type },
                                    showAddDepartmentDialog: false,
                                    validateState: {
                                        ...validateState,
                                        parentId: ValidateState.Normal,
                                    },
                                })
                            }}
                            onCancel={() => {
                                this.setState({
                                    showAddDepartmentDialog: false,
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

    private formatDP = (dp) =>{
        return dp.name
    }

    private changeDp = (dp) =>{
        this.setState({
            departmentInfo: { ...this.state.departmentInfo, parentId: dp.id, parentName: dp.name, parentType: dp.type },
        })
    }
    /**
     * 格式化所有者函数
     */
    private formatOwner = (user: Doclibs.UserInfo): string => {
        return user.name
    }

    /**
     * 更改直属上级
     */
    private changeOwner = (users): void => {
        this.setState({
            managerInfo: users,
        })
    }
}