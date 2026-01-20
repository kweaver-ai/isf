import * as React from 'react';
import { Text } from '@/ui/ui.desktop';
import { getRoleName } from '@/core/role/role';
import Dialog from '@/ui/Dialog2/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import { NodeType } from '@/core/organization';
import OrganizationPick from '../../../OrganizationPick/component.view';
import SetOrgManagerBase from './component.base';
import styles from './styles.desktop';
import __ from './locale';

export default class SetOrgManager extends SetOrgManagerBase {
    render() {
        const {
            selectDeps,
            selectState,
        } = this.state;
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
                                                <span className={styles['user-name']}><Text role={'ui-text'}>{this.props.userInfo.user.displayName}</Text></span>
                                                <span>{__(' 的角色：')}</span>
                                            </div>
                                        ) :
                                        (
                                            <div className={styles['setRole-tit']}>
                                                <span>{__('为用户 “')}</span>
                                                <span className={styles['user-name']}><Text role={'ui-text'}>{this.props.userInfo.user.displayName}</Text></span>
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
                                    height={'290px'}
                                    autoFocus={true}
                                    userid={this.props.userid}
                                    selectType={[NodeType.DEPARTMENT, NodeType.ORGANIZATION]}
                                    selectTit={__('已选部门：')}
                                    placeholder={__('查找')}
                                    data={selectDeps}
                                    converterIn={this.convertData}
                                    convererOut={this.convertDataOut}
                                    onSelectionChange={(value) => { this.selectDeparment(value) }}
                                />
                                {
                                    <span className={styles['errMessage']}>
                                        {
                                            selectState ? __('请至少添加一个部门。') : ''
                                        }
                                    </span>
                                }
                            </div>
                        </div>
                    </Panel.Main>
                    <Panel.Footer role={'ui-panel.footer'}>
                        <Panel.Button
                            theme='oem'
                            role={'ui-panel.button'}
                            onClick={this.validateRole.bind(this)}
                        >
                            {
                                this.props.editRateInfo ?
                                    __('确定') : __('添加')
                            }
                        </Panel.Button>
                        <Panel.Button
                            role={'ui-panel.button'}
                            onClick={this.cancelSetRoleConfig.bind(this)}
                        >
                            {__('取消')}
                        </Panel.Button>
                    </Panel.Footer>
                </Panel>
            </Dialog>
        )
    }
}