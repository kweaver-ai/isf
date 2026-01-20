import * as React from 'react';
import classnames from 'classnames';
import { Text } from '@/ui/ui.desktop';
import { getRoleName } from '@/core/role/role';
import Dialog from '@/ui/Dialog2/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import { NodeType } from '@/core/organization';
import OrganizationPick from '../../../OrganizationPick/component.view';
import SetOrgAuditBase from './component.base';
import styles from './styles.desktop';
import __ from './locale';

export default class SetOrgAudit extends SetOrgAuditBase {
    render() {
        return (
            <Dialog
                role={'ui-dialog'}
                title={__('设置角色')}
                onClose={() => this.cancelSetRoleConfig()}
            >
                <Panel role={'ui-panel'}>
                    <Panel.Main role={'ui-panel.main'}>
                        <div className={styles['layout']}>
                            <div className={styles['setRole-head']}>
                                {
                                    this.props.editRateInfo ?
                                        (
                                            <div className={styles['setRole-tit']}>
                                                <span className={styles['user-name']}><Text>{this.props.userInfo.user.displayName}</Text></span>
                                                <span>{__(' 的角色：')}</span>
                                            </div>
                                        ) :
                                        (
                                            <div className={styles['setRole-tit']}>
                                                <span>{__('为用户 “')}</span>
                                                <span className={styles['user-name']}><Text>{this.props.userInfo.user.displayName}</Text></span>
                                                <span>{__('” 添加角色：')}</span>
                                            </div>
                                        )
                                }
                                <span className={styles['roleName']}>
                                    {getRoleName(this.props.roleInfo)}
                                </span>
                            </div>
                            <span className={styles['manageRange']}>
                                {__('管辖范围：')}
                            </span>
                            <div className={styles['managerSet']}>
                                <OrganizationPick
                                    width={'256px'}
                                    height={'340px'}
                                    autoFocus={true}
                                    userid={this.props.userid}
                                    selectType={[NodeType.DEPARTMENT, NodeType.ORGANIZATION]}
                                    selectTit={__('已选部门：')}
                                    placeholder={__('查找')}
                                    data={this.state.selectDeps}
                                    converterIn={this.convertData}
                                    convererOut={this.convertDataOut}
                                    onSelectionChange={(value) => { this.selectDeparment(value) }}
                                />
                            </div>
                        </div>
                    </Panel.Main>
                    <Panel.Footer role={'ui-panel.footer'}>
                        <Panel.Button
                            theme='oem'
                            role={'ui-panel.button'}
                            onClick={() => this.confirmSetRoleConfig()}
                            disabled={!this.state.selectDeps.length}
                        >
                            {
                                this.props.editRateInfo ?
                                    __('确定') : __('添加')
                            }
                        </Panel.Button>
                        <Panel.Button
                            role={'ui-panel.button'}
                            onClick={() => this.cancelSetRoleConfig()}
                        >
                            {__('取消')}
                        </Panel.Button>
                    </Panel.Footer>
                </Panel>
            </Dialog>
        )
    }
}