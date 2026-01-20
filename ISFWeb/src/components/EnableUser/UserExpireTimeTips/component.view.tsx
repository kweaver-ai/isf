import * as React from 'react';
import { Dialog2 as Dialog, Panel, Text, UIIcon } from '@/ui/ui.desktop';
import { DataGrid } from '@/sweet-ui';
import __ from './locale';
import UserExpireTimeTipsBase from './component.base';
import styles from './styles.view';

export default class UserExpireTimeTips extends UserExpireTimeTipsBase {
    render() {
        return (
            <div>
                <Dialog
                    role={'ui-dialog'}
                    title={__('提示')}
                    onClose={() => this.props.cancelExpireTimeTips()}
                >
                    <Panel role={'ui-panel'}>
                        <Panel.Main role={'ui-panel.main'}>
                            {
                                this.props.expireTimeUsers.length !== 1 ?
                                    <div>
                                        <div className={styles['expire-time-tips-title']}>
                                            {
                                                __('以下 ${expirtTimeUsersCount} 个用户的账号已过期，是否重新设置有效期限？', { expirtTimeUsersCount: this.props.expireTimeUsers.length })
                                            }
                                        </div>
                                        <div className={styles['expire-time-tips-datagrid']}>
                                            <DataGrid
                                                role={'sweetui-datagrid'}
                                                height={400}
                                                data={this.props.expireTimeUsers}
                                                columns={[
                                                    {
                                                        title: __('显示名称'),
                                                        key: 'displayName',
                                                        width: '35%',
                                                        renderCell: (displayName, record) => (
                                                            <Text role={'ui-text'}>{record.user.displayName}</Text>
                                                        ),
                                                    },
                                                    {
                                                        title: __('用户名称'),
                                                        key: 'loginName',
                                                        width: '35%',
                                                        renderCell: (loginName, record) => (
                                                            <Text role={'ui-text'}>{record.user.loginName}</Text>
                                                        ),
                                                    },
                                                    {
                                                        title: __('直属部门'),
                                                        key: 'departmentNames',
                                                        width: '30%',
                                                        renderCell: (departmentNames, record) => (
                                                            <Text role={'ui-text'}>
                                                                {this.getDepartmentName(record.user)}
                                                            </Text>
                                                        ),
                                                    },
                                                ]}
                                            />
                                        </div>
                                    </div>
                                    :
                                    <div>
                                        <div className={styles['icon-warning']}>
                                            <UIIcon
                                                role={'ui-uiicon'}
                                                code={'\uf076'}
                                                color={'#5a8cb4'}
                                                size={40}
                                            />
                                        </div>
                                        <div className={styles['is-reset-expire-time-tips']}>
                                            {__('该用户账号已过期，是否重新设置有效期限？')}
                                        </div>
                                    </div>
                            }
                        </Panel.Main>
                        <Panel.Footer role={'ui-panel.footer'}>
                            <Panel.Button role={'ui-panel.button'} onClick={() => this.props.completeExpireTimeTips()}>{__('确定')}</Panel.Button>
                            <Panel.Button role={'ui-panel.button'} onClick={() => this.props.cancelExpireTimeTips()}>{__('取消')}</Panel.Button>
                        </Panel.Footer>
                    </Panel>
                </Dialog>
            </div>
        )
    }

    private getDepartmentName(user: any): any {
        return user.departmentNames ? user.departmentNames.join(',') : __('未分配组')
    }
}