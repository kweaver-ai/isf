import * as React from 'react';
import { Select } from '@/sweet-ui';
import { getErrorMessage } from '@/core/exception';
import { Range, Status } from './component.base';
import MessageDialog from '@/ui/MessageDialog/ui.desktop';
import ConfirmDialog from '@/ui/ConfirmDialog/ui.desktop';
import DeleteUserBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class DeleteUser extends DeleteUserBase {

    render() {
        return (
            <div className={styles['container']}>
                {
                    this.state.status === Status.NORMAL ?
                        <ConfirmDialog
                            role={'ui-confirmdialog'}
                            onConfirm={this.confirmDeleteUsers.bind(this)}
                            onCancel={this.props.onComplete}
                        >
                            <div className={styles['select-dialog']}>
                                <div className={styles['title']}>
                                    <label>{__('您将彻底删除 ')} </label>
                                    <span style={{ verticalAlign: 'middle', display: 'inline-block' }}>
                                        <Select
                                            role={'sweetui-select'}
                                            value={this.state.selected}
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
                                                            selected={this.state.selected === value}
                                                        >
                                                            {
                                                                {
                                                                    [Range.USERS]: __('当前选中用户'),
                                                                    [Range.DEPARTMENT]: __('${name} 部门成员', { name: this.props.dep.name }),
                                                                    [Range.DEPARTMENT_DEEP]: __('${name} 及其子部门成员', { name: this.props.dep.name }),
                                                                }[value]
                                                            }
                                                        </Select.Option>
                                                    )
                                                })
                                            }

                                        </Select>
                                    </span>
                                    <label>{__(' 的个人账号，确定执行此操作吗？')}</label>
                                </div>
                            </div>
                        </ConfirmDialog> :
                        null
                }

                {

                    this.state.status !== Status.NORMAL && this.state.status !== Status.LOADING ?
                        <MessageDialog
                            role={'ui-messagedialog'}
                            onConfirm={() => { this.props.onComplete() }}
                        >
                            {
                                this.getErrorMessage(this.state.status)
                            }
                        </MessageDialog> :
                        null
                }
                {
                    this.state.status === Status.LOADING ?
                        <Spin size='large' tip={__('正在删除用户，请稍候...')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/> :
                        null
                }

            </div>)
    }

    getErrorMessage(error) {
        switch (error) {
            case Status.CURRENT_USER_INCLUDED:
                return __('您无法删除自身账号。');

            default:
                return getErrorMessage(error);
        }
    }
}