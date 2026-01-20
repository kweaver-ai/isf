import * as React from 'react';
import classnames from 'classnames';
import { UIIcon, Title } from '@/ui/ui.desktop'
import { AutoComplete } from '@/sweet-ui';
import { isBrowser, Browser } from '@/util/browser';
import { convertPath } from '../DomainTree/helper';
import SearchDomainUserBase from './component.base';
import styles from './styles.desktop';
import __ from './locale';

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class SearchDomainUser extends SearchDomainUserBase {
    render() {
        const { searchKey, domainError, results } = this.state;

        return (
            <AutoComplete
                role={'ui-autocomplete'}
                disabled={this.props.disabled}
                width={this.props.width}
                maxHeight={210}
                autoFocus={false}
                value={searchKey}
                limit={10}
                loader={({ key, start, limit }) => this.getDomainUserOrDepByKey(key, start, limit)}
                onError={({ detail }) => this.loaderFailed(detail)}
                onValueChange={({ detail }) => this.handelChange(detail)}
                onLoad={({ detail }) => this.getSearchData(detail)}
                onBlur={this.handelOnBlur}
                placeholder={this.props.placeholder}
                ListEmptyComponent={!searchKey ? '' : domainError ? '' : __('未找到匹配的结果')}
                onSelect={({ detail }) => this.selectItem(detail)}
                onPressEnter={({ detail }) => this.selectItem(detail)}
                getItemLayout={() => ({ length: 52 })}
                data={results}
                renderItem={(record) => {
                    return (
                        <span className={styles['search-item']}>
                            <span className={styles['selected-data-Icon']}>
                                <UIIcon
                                    role={'ui-uiicon'}
                                    code={record.name ? '\uf009' : '\uf007'}
                                    size={16}
                                />
                            </span>
                            <span className={styles['seleted-data']}>
                                <Title content={record.name || record.displayName}>
                                    <div className={classnames(
                                        styles['allname'],
                                        {
                                            [styles['safari']]: isSafari,
                                        },
                                    )}>
                                        <span className={styles['dename']}>{record.name || record.displayName}</span>
                                        <div className={styles['depaths']}>
                                            {convertPath(record.parentOUPath || record.ouPath)}
                                        </div>
                                    </div>
                                </Title>
                            </span>
                        </span>
                    )
                }}
            >
            </AutoComplete>
        )
    }

}