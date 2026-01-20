
import * as React from 'react';
import { MessageDialog, Dialog2 as Dialog, Panel, ErrorDialog } from '@/ui/ui.desktop'
import { Select } from '@/sweet-ui'
import ValidityBox2 from '../ValidityBox2/component.view';
import { Range } from '../helper'
import SetUserExpireTimeBase, { Status } from './component.base';
import styles from './styles.view';
import __ from './locale';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class SetUserExpireTime extends SetUserExpireTimeBase {
    render() {
        return (
            <div>
                {
                    this.state.status === Status.Normal ?
                        <Dialog
                            role={'ui-dialog'}
                            width={425}
                            title={__('设置用户有效期限')}
                            onClose={this.props.onCancel}
                        >
                            <Panel role={'ui-panel'}>
                                <Panel.Main role={'ui-panel.main'}>
                                    <div className={styles['container']}>
                                        {
                                            !this.props.shouldEnableUsers ?
                                                [
                                                    __('您可以设置'),
                                                    <div className={styles['select-range']} key={'selectRange'}>
                                                        <Select
                                                            role={'sweetui-select'}
                                                            value={this.state.selected}
                                                            onChange={({ detail }) => this.onSelectedType(detail)}
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
                                                                        <Select.Option role={'sweetui-select.option'} key={value} value={value} selected={this.state.selected === value}>
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
                                                    </div>,
                                                    __(' 的有效期限为：'),
                                                ]
                                                :
                                                null
                                        }
                                        < div className={styles['select-date']}>
                                            <ValidityBox2
                                                width="100%"
                                                allowPermanent={true}
                                                value={this.state.expireTime}
                                                selectRange={[new Date()]}
                                                onChange={(value) => { this.changeExpireTime(value) }}
                                            />
                                        </div>
                                    </div>
                                </Panel.Main>
                                <Panel.Footer role={'ui-panel.footer'}>
                                    <Panel.Button
                                        role={'ui-panel.button'}
                                        type="submit"
                                        onClick={this.confirmSetUserExpireTime.bind(this)}
                                    >
                                        {__('确定')}
                                    </Panel.Button>
                                    <Panel.Button
                                        role={'ui-panel.button'}
                                        onClick={this.props.onCancel}
                                    >
                                        {__('取消')}
                                    </Panel.Button>
                                </Panel.Footer>
                            </Panel>
                        </Dialog>
                        :
                        null
                }

                {
                    this.state.status === Status.CurrentUser ?
                        <MessageDialog role={'ui-messagedialog'} onConfirm={() => { this.props.onCancel() }}>
                            {
                                __('您无法为自身账号设置有效期。')
                            }
                        </MessageDialog>
                        : null
                }

                {
                    this.state.invalidExpireTime ?
                        <MessageDialog role={'ui-messagedialog'} onConfirm={() => this.setState({ invalidExpireTime: false })}>
                            {
                                __('该日期已过期，请重新选择。')
                            }
                        </MessageDialog>
                        : null
                }

                {
                    this.state.status === Status.Error ?
                        <ErrorDialog role={'ui-errordialog'} onConfirm={() => { this.props.onCancel() }}>
                            {this.state.errors.map((error) => {
                                <div key={error.errID}>{error.errMsg}</div>
                            })}
                        </ErrorDialog>
                        : null
                }

                {
                    this.state.status === Status.Loading ?
                        <Spin size='large' tip={this.props.shouldEnableUsers ? __('正在启用，请稍候...') : __('正在设置，请稍候...')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/> :
                        null
                }
            </div>
        )
    }
}