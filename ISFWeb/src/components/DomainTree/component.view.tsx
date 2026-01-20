import * as React from 'react'
import classnames from 'classnames'
import Tree from '@/ui/Tree2/ui.desktop'
import { CascadeDirection, SelectType as NodeSelectType } from '@/ui/Tree2/ui.base'
import { Text, UIIcon, Centered, Icon, LinkChip } from '@/ui/ui.desktop'
import SearchDomainUser from '../SearchDomainUser/component.view'
import DomainTreeBase from './component.base'
import { isLeaf, getNodeIcon } from './helper'
import * as loadingImg from './assets/images/loading.gif'
import styles from './styles.desktop'
import __ from './locale'

export default class DomainTree extends DomainTreeBase {
    private renderNode = (node) => {
        return (
            <div className={classnames(
                styles['inline'],
            )} title={node.name || node.displayName}>
                <UIIcon
                    role={'ui-uiicon'}
                    {...getNodeIcon(node)}
                    size={16}
                    className={styles['node']}
                />
                <Text role={'ui-text'}> {node.name || node.displayName} </Text>
            </div >
        )
    }

    private getCheckStatus = (node, item) => {
        if (node.parentNode) {
            const { parentNode, parentNode: { objectGUID, name } } = node;

            if ((parentNode && objectGUID && objectGUID.indexOf(item.objectGUID) !== -1)
                || (parentNode && name && name.indexOf(item.name) !== -1 && item.ipAddress)) {

                return true
            } else {
                return this.getCheckStatus(node.parentNode, item)
            }
        }

        return
    }

    render() {
        const { domainId, doRedirectDomain } = this.props;
        const { treeData } = this.state;

        return (
            <div>
                <div className={styles['search-box']}>
                    <SearchDomainUser
                        width={230}
                        disabled={!treeData.length}
                        domainId={domainId}
                        placeholder={__('搜索用户或部门')}
                        onRequestSelect={this.handleSelectResult}
                    />
                </div>
                <div className={styles['tree-list']}>
                    {!this.state.loading ?
                        <div>
                            {
                                treeData.length ?
                                    <Tree
                                        role={'ui-tree2'}
                                        selectType={NodeSelectType.CASCADE_MULTIPLE}
                                        cascadeDirection={CascadeDirection.DOWN}
                                        checkbox={true}
                                        data={treeData}
                                        isLeaf={(node) => isLeaf(node)}
                                        renderNode={this.renderNode}
                                        getNodeChildren={this.getNodeChildren}
                                        onSelectionChange={this.handleSelect}
                                        ref={(ref) => this.ref = ref}
                                    /> :
                                    <div className={styles['empty']}>
                                        <p>{__('暂无可选择的用户和部门')}</p>
                                        <p className={styles['remarks']}>
                                            <span>{__('请前往')}</span>
                                            <LinkChip
                                                role={'ui-linkchip'}
                                                className={styles['link']}
                                                onClick={doRedirectDomain}
                                            >
                                                <a>{__('【域认证】')}</a>
                                            </LinkChip>
                                            <span>{__('设置')}</span>
                                        </p>
                                    </div>
                            }
                        </div>
                        :
                        < Centered role={'ui-centered'}>
                            <Icon
                                role={'ui-icon'}
                                url={loadingImg}
                                size={16}
                            />
                        </Centered>
                    }
                </div>
            </div>
        )
    }
}