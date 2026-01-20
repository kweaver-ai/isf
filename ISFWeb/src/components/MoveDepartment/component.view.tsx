import * as React from 'react';
import { shrinkText } from '@/util/formatters';
import { ModalDialog2, SweetIcon } from '@/sweet-ui';
import { Title, Icon } from '@/ui/ui.desktop';
import { NodeType } from '@/core/organization';
import OrganizationTree from '../OrganizationTree/component.view';
import * as loadImg from './assets/loading.gif';
import MoveDepartmentBase from './component.base';
import { Status } from './helper';
import __ from './locale';
import styles from './styles.view';

export default class MoveDepartment extends MoveDepartmentBase {
    render() {
        return (
            <div>{this.getTemplate(this.state.status)}</div>
        )
    }

    private getTemplate(status: Status): React.ReactNode {
        switch (status) {
            case Status.Config:
                return (
                    <ModalDialog2
                        role={'sweetui-modaldialog2'}
                        title={__('移动部门')}
                        zIndex={18}
                        icons={[
                            {
                                icon: <SweetIcon name="x" size={16} role={'sweetui-sweeticon'}/>,
                                onClick: this.props.onRequestCancelMoveDep,
                            },
                        ]}
                        buttons={[
                            {
                                text: __('确定'),
                                theme: 'oem',
                                disabled: !this.state.targetDep,
                                onClick: this.confirmMoveDep,
                            },
                            {
                                text: __('取消'),
                                theme: 'regular',
                                onClick: this.props.onRequestCancelMoveDep,
                            },
                        ]}
                    >
                        <div>
                            {__('您可以将部门 “')}
                            <Title content={this.props.srcDep.name} role={'ui-title'}>{shrinkText(this.props.srcDep.name)}</Title>
                            {__('” 移动至以下选中的部门下面：')}
                        </div>
                        <div className={styles['tree-wrapper']}>
                            <OrganizationTree
                                userid={this.userid}
                                selectType={[NodeType.DEPARTMENT, NodeType.ORGANIZATION]}
                                onSelectionChange={this.selectDep}
                                getNodeStatus={this.getDepartmentStatus}
                                isDisableChildrenByParent={true}
                            />
                        </div>
                    </ModalDialog2>
                )

            case Status.Loading:
                return (
                    <ModalDialog2
                        role={'sweetui-modaldialog2'}
                        title={__('移动部门')}
                    >
                        <div className={styles['loading']}>
                            <div className={styles['loading-icon']}>
                                <Icon
                                    role={'ui-icon'}
                                    url={loadImg}
                                    size={32}
                                />
                            </div>
                            {__('正在替换存储位置，请稍候...')}
                        </div>
                    </ModalDialog2>
                )

            default:
                return null
        }
    }

    /**
     * 禁用当前部门
     * @param node 部门节点
     */
    private getDepartmentStatus = (node: Core.ShareMgnt.ncTDepartmentInfo): { disabled: boolean } => {
        return { disabled: this.props.srcDep.id === node.id }
    }
}