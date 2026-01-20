import * as React from 'react';
import classnames from 'classnames';
import { FlexBox, Title, UIIcon } from '@/ui/ui.desktop';
import Button from '@/ui/Button/ui.desktop';
import { isBrowser, Browser } from '@/util/browser';
import { getDepName } from '@/core/organization';
import SearchDep from '../SearchDep/component.desktop';
import OrganizationTree from '../OrganizationTree/component.view';
import OrganizationPickBase from './component.base';
import styles from './styles.desktop';
import __ from './locale';
import * as deleteIcon from './assets/delete.png'

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class OrganizationPick extends OrganizationPickBase {

    render() {
        return (
            <div>
                <FlexBox role={'ui-flexbox'}>
                    <FlexBox.Item
                        role={'ui-flexBox.item'}
                        align="left middle"
                    >
                        <div className={styles['search-box']}>
                            <SearchDep
                                onSelectDep={(value) => { this.selectDep(value) }}
                                userid={this.props.userid}
                                width={this.props.width ? this.props.width : 202}
                                placeholder={this.props.placeholder}
                                autoFocus={this.props.autoFocus}
                                selectType={this.props.selectType}
                                isShowUndistributed={this.props.isShowUndistributed}
                                isShowDisabledUsers={this.props.isShowDisabledUsers}
                            />
                        </div>
                        <div
                            style={{ height: this.props.height, width: this.props.width }}
                            className={styles['organization-tree']}
                        >
                            <div className={styles['tree-main']}>
                                <OrganizationTree
                                    isShowUndistributed={this.props.isShowUndistributed}
                                    isShowDisabledUsers={this.props.isShowDisabledUsers}
                                    disabled={this.props.disabled}
                                    userid={this.props.userid}
                                    selectType={this.props.selectType}
                                    onSelectionChange={(value) => { this.selectDep(value) }}
                                />
                            </div>
                        </div>
                    </FlexBox.Item>
                    <FlexBox.Item
                        role={'ui-flexbox.item'}
                        align="right middle"
                    >
                        <div className={styles['select-content']}>
                            <FlexBox role={'ui-flexbox'}>
                                <FlexBox.Item role={'ui-flexbox.item'} align="left middle">
                                    <label>
                                        {__('已选：')}
                                    </label>
                                </FlexBox.Item>
                                <FlexBox.Item role={'ui-flexbox.item'} align="right middle">
                                    <div>
                                        <Button
                                            role={'ui-button'}
                                            onClick={this.clearSelectDep.bind(this)}
                                            disabled={!this.state.data.length}
                                        >
                                            {__('清空')}
                                        </Button>
                                    </div>
                                </FlexBox.Item>
                            </FlexBox>
                        </div>
                        <div
                            style={{ height: this.props.height, width: this.props.width }}
                            className={classnames(styles['organization-tree'], styles['select-content'])}
                        >
                            <ul>
                                {
                                    this.state.data.map((sharer) => (
                                        <li
                                            key={sharer.id}
                                            style={{ position: 'relative' }}
                                            className={styles['dep-item']}>
                                            <div className={styles['seleted-data']}>
                                                <Title role={'ui-title'} content={getDepName(sharer)}>
                                                    <div className={classnames(
                                                        styles['dep-name'],
                                                        {
                                                            [styles['safari']]: isSafari,
                                                        },
                                                    )}>
                                                        {sharer.name}
                                                    </div>
                                                </Title>
                                            </div>
                                            <div className={styles['selected-data-del']}>
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    size={13}
                                                    code={'\uf014'}
                                                    fallback={deleteIcon}
                                                    onClick={() => { this.deleteSelectDep(sharer) }}
                                                />
                                            </div>
                                        </li>
                                    ))
                                }
                            </ul>
                        </div>
                    </FlexBox.Item>
                </FlexBox>
            </div>
        )
    }
}