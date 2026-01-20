import * as React from 'react';
import { ModalDialog2, SweetIcon } from '@/sweet-ui';
import { Title, Icon } from '@/ui/ui.desktop';
import { shrinkText } from '@/util/formatters';
import { NodeType } from '@/core/organization';
import OrganizationPick from '../OrganizationPick/component.view';
import * as loadImg from './assets/loading.gif';
import AddUsersToDepBase from './component.base';
import { Status } from './helper';
import __ from './locale';
import styles from './styles.view';

export default class AddUsersToDep extends AddUsersToDepBase {
    render() {
        return <div>{this.getTemplate(this.state.renderStatus)}</div>
    }

    /**
     * 获取显示界面
    */
    private getTemplate(status: Status): React.ReactNode {
        switch (status) {
            case Status.Config:
                return (
                    <ModalDialog2
                        title={__('添加用户至部门')}
                        zIndex={18}
                        icons={[
                            {
                                icon: <SweetIcon name="x" size={16} />,
                                onClick: this.props.onRequestCancel,
                            },
                        ]}
                        buttons={[
                            {
                                text: __('确定'),
                                theme: 'oem',
                                disabled: !this.state.users.length,
                                onClick: this.confrimAddUsers,
                            },
                            {
                                text: __('取消'),
                                theme: 'regular',
                                onClick: this.props.onRequestCancel,
                            },
                        ]}
                    >
                        <div>
                            {__('选择用户添加至部门')}
                            <Title content={this.props.targetDep.name}>{`“${shrinkText(this.props.targetDep.name)}”`}</Title>
                        </div>
                        <OrganizationPick
                            userid={this.userid}
                            isShowUndistributed={this.isShowUndistributed}
                            selectType={[NodeType.USER]}
                            data={this.state.users}
                            convererOut={(user) => user}
                            onSelectionChange={this.addUsers}
                        />
                    </ModalDialog2>
                )

            case Status.Adding:
            case Status.ChangeOSSing:
                return (
                    <ModalDialog2
                        title={__('添加用户至部门')}
                    >
                        <div className={styles['loading']}>
                            <div className={styles['loading-icon']}>
                                <Icon
                                    url={loadImg}
                                    size={32}
                                />
                            </div>
                            {status === Status.Adding ? __('正在添加用户，请稍候...') : __('正在替换存储位置，请稍候...')}
                        </div>
                    </ModalDialog2>
                )

            default:
                return null
        }
    }
}