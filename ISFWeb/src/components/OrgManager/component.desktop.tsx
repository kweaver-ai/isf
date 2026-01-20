import * as React from 'react';
import classnames from 'classnames';
import ComboArea from '@/ui/ComboArea/ui.desktop';
import FlexBox from '@/ui/FlexBox/ui.desktop';
import Button from '@/ui/Button/ui.desktop';
import { NodeType } from '@/core/organization';
import OrganizationPicker from '../OrganizationPicker/component.view';
import OrgManagerBase from './component.base';
import __ from './locale';
import styles from './styles.desktop';

export default class OrgManager extends OrgManagerBase {
    render() {
        return (
            <div className={styles['container']}>
                <FlexBox>
                    <FlexBox.Item align="left top">
                        <div className={styles['dep-label']}>{__('管辖部门：')}</div>
                    </FlexBox.Item>
                    <FlexBox.Item align="left top">
                        <div className={styles['dep-box']}>
                            <ComboArea
                                minHeight={80}
                                width={380}
                                uneditable={true}
                                disabled={!this.props.isManager}
                                value={this.state.deps}
                                formatter={this.depDataFormatter}
                                onChange={(data) => { this.setDepsData(data) }}
                            />
                        </div>
                    </FlexBox.Item>
                    <FlexBox.Item align="right top">
                        <div className={styles['dep-btn']}>
                            <Button disabled={!this.props.isManager} onClick={this.showOrganizationPicker.bind(this)}>
                                {__('添加')}
                            </Button>
                        </div>
                    </FlexBox.Item>
                </FlexBox>
                {
                    this.state.showOrganizationPicker ?
                        <OrganizationPicker
                            onCancel={this.cancelAddDeps.bind(this)}
                            onConfirm={(deps) => { this.confirmAddDeps(deps) }}
                            userid={this.props.userid}
                            convererOut={(value) => this.convererData(value)}
                            selectType={[NodeType.DEPARTMENT, NodeType.ORGANIZATION]}
                        /> :
                        null
                }
            </div>
        )
    }
}