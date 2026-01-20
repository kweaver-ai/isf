import * as React from 'react';
import { Select } from '@/sweet-ui'
import ConfirmDialog from '@/ui/ConfirmDialog/ui.desktop';
import UserExpireTimeTips from './UserExpireTimeTips/component.view';
import SetUserExpireTime from '../SetUserExpireTime/component.view'
import { Range } from '../helper';
import RemoveUserBase from './component.base';
import { Status } from './helper';
import __ from './locale';
import styles from './styles.view';
import ErrorMessage from './ErrorMessage/component.view';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class RemoveUser extends RemoveUserBase {

    render() {
        const { expireTimeUsers, status } = this.state
        return (
            <div>
                {
                    this.state.status === Status.NORMAL ?
                        (<ConfirmDialog
                            role={'ui-confirmdialog'}
                            onConfirm={this.confirmEnableUsers.bind(this)}
                            onCancel={this.props.onComplete}  >
                            <div className={styles['select-dialog']}>
                                <label>{__('您将启用 ')} </label>
                                <div className={styles['select-item']}>
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
                                                        selected={this.state.selected === value}>
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
                                </div>
                                <label>{__(' 的个人账号，确定要执行此操作吗？')}</label>
                            </div>

                        </ConfirmDialog>)
                        :
                        null
                }

                {

                    status !== Status.NORMAL
                        && status !== Status.LOADING
                        && status !== Status.SET_EXPIRE_TIME ?
                        <ErrorMessage
                            errorType={status}
                            errorUser={this.state.errorUser}
                            onConfirm={this.props.onComplete}
                        /> :
                        null
                }
                {
                    status === Status.LOADING && expireTimeUsers.length === 0 ?
                        <Spin size='large' tip={__('正在启用用户，请稍候……')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />} /> :
                        null
                }

                {
                    expireTimeUsers.length !== 0 && status !== Status.SET_EXPIRE_TIME ?
                        <UserExpireTimeTips
                            expireTimeUsers={expireTimeUsers}
                            completeExpireTimeTips={() => this.completeExpireTimeTips()}
                            cancelExpireTimeTips={() => this.props.onComplete()}
                            userid={this.props.userid}
                        />
                        : null
                }

                {
                    status === Status.SET_EXPIRE_TIME ?
                        <SetUserExpireTime
                            users={expireTimeUsers}
                            userid={this.props.userid}
                            shouldEnableUsers={true}
                            onCancel={() => this.props.onComplete()}
                            onSuccess={() => this.props.onSuccess(this.props.users)}
                        />
                        :
                        null
                }

            </div>)
    }
}