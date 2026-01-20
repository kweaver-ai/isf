import * as React from 'react';
import Dialog from '@/ui/Dialog2/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import { Select } from '@/sweet-ui';
import { getErrorMessage } from '@/core/exception';
import { NodeType } from '@/core/organization';
import OrganizationTree from '../OrganizationTree/component.view';
import { Status } from './component.base';
import MessageDialog from '@/ui/MessageDialog/ui.desktop';
import DepartmentsOfUserSelector from '../DepartmentsOfUserSelector';
import { Range } from '../helper';
import MoveUserBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class MoveUser extends MoveUserBase {

    render() {
        return (
            <div className={styles['container']}>
                {
                    this.getTemplate(this.state.status)
                }
            </div>
        )
    }

    getTemplate(status) {
        const { users } = this.props
        const { selected, selectedDep } = this.state

        switch (status) {
            case Status.NORMAL:
                return (
                    <Dialog
                        ref={(ref) => this.dialogRef = ref}
                        role={'ui-dialog'}
                        className={styles['dialog']}
                        title={__('移动用户')}
                        onClose={() => { this.props.onComplete() }}
                    >
                        <Panel role={'ui-panel'}>
                            <Panel.Main role={'ui-panel.main'}>
                                <div>
                                    <span> {__('您可以将 ')}</span>
                                    <span className={styles['select-panel']}>
                                        <Select
                                            role={'sweetui-select'}
                                            value={selected}
                                            onChange={({ detail }) => { this.onSelectedType(detail) }}
                                            width={200}
                                        >
                                            {
                                                [Range.USERS, Range.DEPARTMENT, Range.DEPARTMENT_DEEP].filter((value) => {
                                                    if ((this.props.dep.id === '-2' || this.props.dep.id === '-1') && value === Range.DEPARTMENT_DEEP) {
                                                        return false
                                                    } else if (value === Range.USERS && !this.props.users.length) {
                                                        return false
                                                    } else {
                                                        return true
                                                    }
                                                }).map((value) => {

                                                    return (
                                                        <Select.Option
                                                            role={'sweetui-select.option'}
                                                            key={value}
                                                            value={value}
                                                        >
                                                            {
                                                                ({
                                                                    [Range.USERS]: __('当前选中用户'),
                                                                    [Range.DEPARTMENT]: __('${name} 部门成员', { name: this.props.dep.name }),
                                                                    [Range.DEPARTMENT_DEEP]: __('${name} 及其子部门成员', { name: this.props.dep.name }),
                                                                })[value]
                                                            }
                                                        </Select.Option>
                                                    )
                                                })
                                            }

                                        </Select>
                                    </span>
                                    {
                                        this.singleUser
                                            ? <>
                                                <span className={styles['dialog-text-space']}>
                                                    {__('从')}
                                                </span>
                                                <span className={styles['select-panel']}>
                                                    <DepartmentsOfUserSelector
                                                        userInfo={users[0]}
                                                        dep={this.props.dep}
                                                        onSelectionChange={this.onDepChange}
                                                    />
                                                </span>
                                            </>
                                            : null
                                    }
                                    <span className={styles['select-text']}>{__(' 移动至以下指定的部门下面：')}</span>
                                </div>
                                <div className={styles['dep-tree']}>
                                    <OrganizationTree
                                        ref={(ref) => this.organizationTreeRef = ref}
                                        userid={this.props.userid}
                                        selectType={[NodeType.DEPARTMENT, NodeType.ORGANIZATION]}
                                        onSelectionChange={(value) => { this.selectDep(value) }}
                                        getNodeStatus={this.getDepartmentStatus.bind(this)}
                                    />
                                </div>
                                {
                                    // 是否有选中的部门，但是该部门不存在了
                                    selectedDep && selectedDep.notExist
                                        ? (
                                            <MessageDialog role={'ui-messagedialog'} onConfirm={this.depNotExist}>
                                                {
                                                    __('无法移动用户，您选中的目标部门 “${depName}” 已不存在，请重新选择。', { depName: this.state.selectedDep.name })
                                                }
                                            </MessageDialog>
                                        )
                                        : null
                                }
                            </Panel.Main>
                            <Panel.Footer role={'ui-panel.footer'}>
                                <Panel.Button theme='oem' role={'ui-panel.button'} onClick={() => { this.confirmMoveUsers() }} disabled={!this.state.selectedDep}>
                                    {__('确定')}
                                </Panel.Button>
                                <Panel.Button role={'ui-panel.button'} onClick={() => { this.props.onComplete() }}>
                                    {__('取消')}
                                </Panel.Button>
                            </Panel.Footer>
                        </Panel>
                    </Dialog >
                )
            case Status.LOADING:
                return (
                    <Spin size='large' tip={__('正在移动用户，请稍候……')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/>
                )

            case Status.CURRENT_USER_INCLUDED:
                return (
                    <MessageDialog role={'ui-messagedialog'} onConfirm={() => { this.props.onComplete() }}>
                        {__('您无法移动自身账号。')}
                    </MessageDialog>
                )
            case Status.DESTDEPARTMENT_NOT_EXIST:
                return (
                    <MessageDialog role={'ui-messagedialog'} onConfirm={() => { this.props.onComplete(this.state.selectedDep) }}>
                        {
                            __('无法移动用户，您选中的目标部门 “${depName}” 已不存在，请重新选择。', { depName: this.state.selectedDep.name })
                        }
                    </MessageDialog>
                )

            case Status.SRRDEP_NOT_EXIST:
                return (
                    <MessageDialog role={'ui-messagedialog'} onConfirm={() => { this.props.onComplete(this.state.dep) }}>
                        {
                            __('无法移动用户，部门 “${depName}” 已不存在，请重新选择。', { depName: this.state.dep.name })
                        }
                    </MessageDialog>
                )

            default:
                return (
                    <MessageDialog role={'ui-messagedialog'} onConfirm={() => { this.props.onComplete() }}>
                        {
                            getErrorMessage(status)
                        }
                    </MessageDialog>
                )
        }
    }
}