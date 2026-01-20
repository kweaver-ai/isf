import * as React from 'react';
import classnames from 'classnames';
import AutoComplete from '@/ui/AutoComplete/ui.desktop';
import AutoCompleteList from '@/ui/AutoCompleteList/ui.desktop';
import { isBrowser, Browser } from '@/util/browser';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import Title from '@/ui/Title/ui.desktop';
import SearchUserGroupBase from './component.base';
import styles from './styles.desktop';
import __ from './locale';

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class SearchUserGroup extends SearchUserGroupBase {

    render() {
        return (
            <AutoComplete
                role={'ui-autocomplete'}
                ref={(autocomplete) => this.autocomplete = autocomplete}
                disabled={this.props.disabled}
                width={this.props.width}
                autoFocus={this.props.autoFocus}
                value={this.state.searchKey}
                loader={this.getGroupsByKey}
                onChange={this.handelChange}
                onLoad={(data) => { this.getSearchData(data) }}
                lazyLoader={
                    {
                        limit: 10,
                        trigger: 0.99,
                        onChange: this.lazyLoad,
                    }
                }
                placeholder={this.props.placeholder ? this.props.placeholder : __('搜索')}
                missingMessage={__('未找到匹配的结果')}
                onEnter={this.handleEnter}
            >
                {
                    this.state.results && this.state.results.length ?
                        <AutoCompleteList role={'ui-autocompleteList'}>
                            {
                                this.state.results.map((value) => (
                                    <AutoCompleteList.Item
                                        role={'ui-autocompleteList.item'}
                                        key={value.id}
                                    >
                                        <span
                                            className={styles['search-item']}
                                            onClick={() => { this.selectItem(value) }}
                                        >
                                            <span className={styles['selected-data-Icon']}>
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    code={'\uf107'}
                                                    size={16}
                                                />
                                            </span>
                                            <span className={styles['seleted-data']}>
                                                <Title role={'ui-title'} content={value.name}>
                                                    <div className={classnames(
                                                        styles['allname'],
                                                        {
                                                            [styles['safari']]: isSafari,
                                                        },
                                                    )}>
                                                        <span className={styles['dename']}>{value.name}</span>
                                                    </div>
                                                </Title>
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