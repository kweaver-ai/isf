import * as React from 'react';
import { Select } from '@/sweet-ui'
import { getErrorMessage } from '@/core/exception';
import { Status } from './component.base';
import MessageDialog from '@/ui/MessageDialog/ui.desktop';
import ConfirmDialog from '@/ui/ConfirmDialog/ui.desktop';
import DepartmentsOfUserSelector from '../DepartmentsOfUserSelector';
import { Range } from '../helper';
import RemoveUserBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class RemoveUser extends RemoveUserBase {

    render() {
        const { users } = this.props
        const { selected } = this.state

        return (
            <div className={styles['container']}>
                {
                    this.state.status === Status.NORMAL
                        ? (
                            <ConfirmDialog
                                role={'ui-confirmdialog'}
                                onConfirm={this.confirmRemoveUsers.bind(this)}
                                onCancel={this.props.onComplete}
                            >
                                <div className={styles['select-dialog']}>
                                    <label>{__('此操作会将 ')} </label>
                                    <span className={styles['container-select']}>
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

                                                    return <Select.Option
                                                        role={'sweetui-select.option'}
                                                        key={value}
                                                        value={value}
                                                    >
                                                        {
                                                            {
                                                                [Range.USERS]: __('当前选中用户'),
                                                                [Range.DEPARTMENT]: __('${name} 部门成员', { name: this.props.dep.name }),
                                                                [Range.DEPARTMENT_DEEP]: __('${name} 及其子部门成员', { name: this.props.dep.name }),
                                                            }[value]
                                                        }
                                                    </Select.Option>
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
                                                <span className={styles['container-select']}>
                                                    <DepartmentsOfUserSelector
                                                        userInfo={users[0]}
                                                        dep={this.props.dep}
                                                        onSelectionChange={this.onDepChange}
                                                    />
                                                </span>
                                                <span >
                                                    {__(' 中移除，但不会删除用户账户，您确定要执行吗？')}
                                                </span>
                                            </>
                                            : (
                                                <span>
                                                    {__('从“${depName}”中移除，但不会删除用户账户，您确定要执行吗？', { depName: this.props.dep.name })}
                                                </span>
                                            )
                                    }
                                </div>
                            </ConfirmDialog>
                        )
                        : null
                }

                {

                    this.state.status !== Status.NORMAL && this.state.status !== Status.LOADING ?
                        <MessageDialog role={'ui-messagedialog'} onConfirm={() => { this.props.onComplete() }}>
                            {
                                this.getErrorMessage(this.state.status)
                            }
                        </MessageDialog> :
                        null
                }
                {
                    this.state.status === Status.LOADING ?
                        <Spin size='large' tip={__('正在移除用户，请稍候……')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/> :
                        null
                }

            </div>)
    }

    getErrorMessage(error) {
        switch (error) {
            case Status.CURRENT_USER_INCLUDED:
                return __('您无法移除自身账号。');
            default:
                return getErrorMessage(error);
        }
    }
}