import * as React from 'react';
import { Select } from '@/sweet-ui';
import ConfirmDialog from '@/ui/ConfirmDialog/ui.desktop';
import { Range } from '../helper';
import DisableUserBase from './component.base';
import { Status } from './helper';
import __ from './locale';
import styles from './styles.view';
import ErrorMessage from './ErrorMessage/component.view';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class DisableUser extends DisableUserBase {

    render() {
        return (
            <div>
                {
                    this.state.status === Status.NORMAL ?
                        <ConfirmDialog
                            role={'ui-confirmdialog'}
                            onConfirm={this.confirmDisableUsers.bind(this)}
                            onCancel={this.props.onComplete} >
                            <div className={styles['select-dialog']}>
                                <label>{__('您将禁用 ')} </label>
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
                        </ConfirmDialog> :
                        null
                }

                {

                    this.state.status !== Status.NORMAL && this.state.status !== Status.LOADING ?
                        <ErrorMessage errorType={this.state.status} onConfirm={this.props.onComplete} /> :
                        null
                }
                {
                    this.state.status === Status.LOADING ?
                        <Spin size='large' tip={__('正在禁用用户，请稍候……')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/> :
                        null
                }

            </div>)
    }
}