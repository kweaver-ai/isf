import * as React from 'react';
import { ModalDialog2, SweetIcon, ValidateBox, ComboArea, Switch } from '@/sweet-ui';
import { Button, Form, Title, UIIcon } from '@/ui/ui.desktop';
import { UserInfoType, ValidateMessages } from '@/core/user';
import CreateOrganizationBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import { Doclibs } from '@/core/doclibs';
import OrgAndAccountPick from '../OrgAndAccountPick/component.view';
import session from '@/util/session';
import { NodeType } from '@/core/organization';
import { TabType } from '../OrgAndAccountPick/helper';

export default class CreateOrganization extends CreateOrganizationBase {
    render() {
        const { validateState, managerInfo, showAddDepartmentLeaderDialog,  organizeInfo: { orgName, code, remark, status, email } } = this.state;

        return (
            <div>
                <ModalDialog2
                    role={'sweetui-modaldialog2'}
                    title={__('新建组织')}
                    width={545}
                    zIndex={18}
                    icons={[
                        {
                            icon: <SweetIcon name="x" size={16} role={'sweetui-sweeticon'}/>,
                            onClick: this.props.onRequestCancelCreateOrg,
                        },
                    ]}
                    buttons={[
                        {
                            text: __('确定'),
                            theme: 'oem',
                            onClick: this.createOrganization,
                        },
                        {
                            text: __('取消'),
                            theme: 'regular',
                            onClick: this.props.onRequestCancelCreateOrg,
                        },
                    ]}
                >
                    <Form role={'ui-form'}>
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label className={styles['left-width']} role={'ui-form.label'}>
                                {
                                    __('组织名称：')
                                }
                                <span className={styles['required']}>*</span>
                            </Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <ValidateBox
                                    role={'sweetui-validatebox'}
                                    width={320}
                                    value={orgName}
                                    disabled={false}
                                    validateMessages={ValidateMessages}
                                    validateState={validateState.orgName}
                                    onBlur={() => this.handleOnBlurOrg()}
                                    onValueChange={({ detail }) => this.handleValueChange({ orgName: detail })}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label
                                className={styles['text-label']}
                                role={'ui-form.label'}
                            >
                                {
                                    __('组织编码：')
                                }
                                <UIIcon
                                    role={'ui-uiicon'}
                                    code={'\uf055'}
                                    size={'16px'}
                                    title={
                                        __('组织唯一标识，全局唯一。')
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
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label
                                className={styles['text-label']}
                                role={'ui-form.label'}
                            >
                                {
                                    __('组织负责人：')
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
                        {/* <Form.Row role={'ui-form.row'}>
                            <Form.Label role={'ui-form.label'}>
                                {
                                    __('存储位置：')
                                }
                            </Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <OssInfoSelect
                                    width={320}
                                    type={LocationType.User}
                                    ossInfo={ossInfo}
                                    validateState={validateState.ossInfo.code}
                                    validateMessages={validateState.ossInfo.tipMsg}
                                    onRequestSelectOssInfo={(oss) => this.updateSelectedOss(oss)}
                                />
                            </Form.Field>
                        </Form.Row> */}
                        <Form.Row role={'ui-form.row'}>
                            <Form.Label className={styles['left-width']} role={'ui-form.label'}>
                                {
                                    __('状态：')
                                }
                            </Form.Label>
                            <Form.Field role={'ui-form.field'}>
                                <div className={styles['mark']}></div>
                                <div className={styles['switch']}>
                                    <Title content={__(`点此${status ? '禁用' : '启用'}组织`)} role={'sweetui-title'}>
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