import * as React from 'react';
import classnames from 'classnames';
import AutoComplete from '@/ui/AutoComplete/ui.desktop';
import AutoCompleteList from '@/ui/AutoCompleteList/ui.desktop';
import { isBrowser, Browser } from '@/util/browser';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { getNodeIcon, FormatedNodeInfo, NodeType } from '@/core/organization';
import SearchDepBase from './component.base';
import styles from './styles.desktop';
import __ from './locale';

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class SearchDep extends SearchDepBase {

    // 部门用户路径信息
    getDepName(value: FormatedNodeInfo): string {
        if (value.type === NodeType.DEPARTMENT && value.parent_path) {
            return (
                `${value.name}(${__('部门')})
${value.parent_path}`
            )
        }

        if (value.type === NodeType.USER) {
            return (
                `${value.name}(${value.account})
${value.parent_path}`
            )
        }

        return value.name
    }

    render() {
        return (
            <AutoComplete
                role={'ui-autocomplete'}
                ref={(autocomplete) => this.autocomplete = autocomplete}
                disabled={this.props.canInput !== undefined ? !this.state.canInputValue : false}
                width={this.props.width}
                autoFocus={this.props.autoFocus}
                value={this.state.searchKey}
                loader={this.getDepsByKey.bind(this)}
                onChange={this.handelChange.bind(this)}
                onLoad={(data) => { this.getSearchData(data) }}
                lazyLoader={
                    {
                        limit: 10,
                        trigger: 0.99,
                        onChange: this.lazyLoade,
                    }
                }
                placeholder={this.props.placeholder ? this.props.placeholder : __('搜索')}
                missingMessage={__('未找到匹配的结果')}
                onEnter={this.handleEnter.bind(this)}
            >
                {
                    this.state.results && this.state.results.length ?
                        <AutoCompleteList role={'ui-autocompletelist'}>
                            {
                                this.state.results.map((value) => (
                                    <AutoCompleteList.Item
                                        role={'ui-autocompletelist.item'}
                                        key={value.id}
                                    >
                                        <span className={styles['search-item']} onClick={() => { this.selectItem(value) }}>
                                            <span className={styles['selected-data-Icon']}>
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    {...getNodeIcon(value)}
                                                    size={16}
                                                />
                                            </span>
                                            <span className={styles['seleted-data']}>
                                                <div
                                                    role={'ui-title'}
                                                    title={this.getDepName(value)}
                                                >
                                                    <div className={classnames(
                                                        styles['allname'],
                                                        {
                                                            [styles['safari']]: isSafari,
                                                        },
                                                    )}>
                                                        <span className={styles['dename']}>{value.name}</span>
                                                        <div className={styles['depaths']}>
                                                            {value.parent_path || ''}
                                                        </div>
                                                    </div>
                                                </div>
                                            </span>
                                        </span>
                                    </AutoCompleteList.Item>
                                ))
                            }
                        </AutoCompleteList>
                        : null
                }
            </AutoComplete>
        )
    }

}