import * as React from 'react';
import classnames from 'classnames'
import { CheckBox } from '@/sweet-ui'
import { SearchBox, Icon, LazyLoader } from '@/ui/ui.desktop';
import ListTipComponent from '../ListTipComponent/component.view';
import { ListTipStatus, ListTipMessage } from '../ListTipComponent/helper';
import AppAccountTreeBase from './component.base';
import * as AppAccountIcon from './assets/appAccount.png'
import styles from './styles.view';
import __ from './locale';

const listTipMessage = {
    ...ListTipMessage,
    [ListTipStatus.OrgEmpty]: __('暂无可选的账户名称'),

}

export default class AppAccountTree extends AppAccountTreeBase {

    render() {
        const { isMult, disabled } = this.props
        const { data,
            listTipStatus,
            searchKey,
            selectionId,
            selections,
        } = this.state

        return (
            <div className={classnames(
                styles['tree-box'],
                {
                    [styles['disabled']]: disabled,
                },
            )}>
                <SearchBox
                    role={'ui-searchbox'}
                    width={'100%'}
                    disabled={disabled}
                    placeholder={__('搜索账户名称')}
                    value={searchKey}
                    onChange={this.changeSearchKey}
                    loader={this.loader}
                    onLoad={this.handleLoadSuccess}
                    onLoadFailed={this.handleLoadFailed}

                />
                <div className={classnames(
                    styles['tree-wrp'],
                    {
                        [styles['multi-tree-wrp']]: isMult,
                    },
                )}>
                    {
                        listTipStatus === ListTipStatus.None
                            ? (
                                <LazyLoader
                                    limit={150}
                                    trigger={0.9}
                                    onChange={this.handleLazyLoad}
                                    ref={(lazyLoaderRef) => this.lazyLoaderRef = lazyLoaderRef}
                                >
                                    {
                                        data.map((selection) => {
                                            const checked = !!selections.find(({ id }) => id === selection.id)
                                            return (
                                                <div
                                                    key={selection.id}
                                                    className={classnames(
                                                        styles['list'],
                                                        {
                                                            [styles['multi-list']]: isMult,
                                                        },
                                                    )}
                                                >
                                                    {
                                                        isMult
                                                            ? (
                                                                <span className={styles['multi-checkbox']}>
                                                                    <CheckBox
                                                                        checked={checked}
                                                                        onCheckedChange={({ detail }) => this.multiSelect(selection, detail)}
                                                                    />
                                                                </span>
                                                            )
                                                            : null
                                                    }
                                                    <div
                                                        className={classnames(
                                                            styles['list-title'],
                                                            {
                                                                [styles['selected']]: data && selectionId === selection.id,
                                                            },
                                                            {
                                                                [styles['multi-line']]: isMult,
                                                            },
                                                        )}
                                                        title={selection.name}
                                                        onClick={() => isMult ? this.multiSelect(selection, !checked) : this.addSelection(selection)}
                                                    >
                                                        <Icon
                                                            role={'ui-uiicon'}
                                                            url={AppAccountIcon}
                                                            size={16}
                                                        />
                                                        <div className={styles['list-name']}>
                                                            {selection.name}
                                                        </div>
                                                    </div>
                                                </div>
                                            )
                                        })
                                    }
                                </LazyLoader>
                            )
                            : (
                                <ListTipComponent
                                    listTipStatus={listTipStatus}
                                    listTipMessage={listTipMessage}
                                    isInDialog={true}
                                />
                            )
                    }
                </div>
            </div >
        )
    }
}