import * as React from 'react';
import classnames from 'classnames';
import Tree from '@/ui/Tree2/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { Text } from '@/ui/ui.desktop';
import LazyLoader from '@/ui/LazyLoader/ui.desktop';
import { SelectType } from '@/ui/Tree2/ui.base';
import ListTipComponent from '../ListTipComponent/component.view';
import { ListTipStatus, ListTipMessage } from '../ListTipComponent/helper';
import UserGroupTreeBase from './component.base';
import SearchUserGroup from '../SearchUserGroup/component.view';
import __ from './locale';
import styles from './styles.view';

const listTipMessage = {
    ...ListTipMessage,
    [ListTipStatus.OrgEmpty]: __('暂无可选的用户组'),
}

export default class UserGroupTree extends UserGroupTreeBase {
    renderNode(node) {
        return (
            <div title={node.name} className={styles['node-title']}>
                <UIIcon role={'ui-uiicon'} code={'\uf107'} size={16} className={styles['node-icon']} />
                <Text role={'ui-text'} className={styles['node-name']}> {node.name} </Text>
            </div>
        )
    }

    render() {

        const { userGroups, listTipStatus } = this.state;

        const { isMultSelect, disabled } = this.props;

        return (
            <div className={styles['wrp']}>
                <div className={styles['search-box']}>
                    <SearchUserGroup
                        width={'100%'}
                        disabled={disabled}
                        placeholder={__('搜索用户组')}
                        onRequestSelect={this.handleSelectResult}
                    />
                </div>
                <div
                    className={classnames(
                        styles['tree'],
                        {
                            [styles['disabled-tree']]: disabled,
                        },
                    )}
                >
                    {
                        listTipStatus === ListTipStatus.None ?
                            <LazyLoader
                                limit={50}
                                trigger={0.999}
                                onChange={this.handleLazyLoad}
                            >
                                <div className={styles['left-tree']}>
                                    <Tree
                                        disabled={disabled}
                                        selectType={isMultSelect ? SelectType.MULTIPLE : SelectType.SINGLE}
                                        checkbox={isMultSelect}
                                        data={userGroups}
                                        isLeaf={() => true}
                                        renderNode={this.renderNode}
                                        ref={(tree) => { this.tree = tree }}
                                        onSelectionChange={this.handleSelectionsChange}
                                    />
                                </div>
                            </LazyLoader>
                            :
                            <div className={styles['list-tip']}>
                                <ListTipComponent
                                    listTipStatus={listTipStatus}
                                    listTipMessage={listTipMessage}
                                    isInDialog={true}
                                />
                            </div>
                    }
                </div>
            </div>
        )
    }
}