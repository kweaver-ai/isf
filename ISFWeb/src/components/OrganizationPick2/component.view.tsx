import * as React from 'react';
import { isFunction } from 'lodash';
import { FlexBox, UIIcon } from '@/ui/ui.desktop';
import { SweetIcon, Button } from '@/sweet-ui';
import { CascadeDirection, SelectType as NodeSelectType } from '@/ui/Tree2/ui.base';
import { isBrowser, Browser } from '@/util/browser';
import { getDepName } from '@/core/organization'
import SearchDep from '../SearchDep/component.desktop';
import DepartmentTree from '../DepartmentTree/component.view';
import OrganizationPick2Base from './component.base';
import styles from './styles.desktop';
import __ from './locale';

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class OrganizationPick2 extends OrganizationPick2Base {

    render() {
        const {
            userid,
            searchBoxwWidth,
            autoFocus,
            placeholder,
            disabled,
            selectType,
            extraRoots,
            describeTip,
            formatSelectedItem,
        } = this.props

        const { selections } = this.state

        return (
            <div className={styles['flex']}>
                <div className={styles['flex-item']}>
                    <div className={styles['tree-container']}>
                        <div className={styles['describe-tip']}>
                            {typeof describeTip === 'function' ? describeTip() : describeTip}
                        </div>
                        <div className={styles['organization-tree']}>
                            <div className={styles['search-box']}>
                                <SearchDep
                                    onSelectDep={this.select}
                                    userid={userid}
                                    width={searchBoxwWidth || '100%'}
                                    placeholder={placeholder}
                                    autoFocus={autoFocus}
                                    selectType={selectType}
                                />
                            </div>
                            <div className={styles['tree-main']}>
                                <DepartmentTree
                                    extraRoots={extraRoots}
                                    selectType={selectType}
                                    disabled={disabled}
                                    nodeSelectType={NodeSelectType.CASCADE_MULTIPLE}
                                    cascadeDirection={CascadeDirection.DOWN}
                                    ref={(departmentTree) => this.departmentTree = departmentTree}
                                />
                            </div>
                        </div>
                    </div>
                </div>
                <div className={styles['flex-arrow']}>
                    <UIIcon
                        role={'ui-uiicon'}
                        className={styles['arrow-right']}
                        code={'\uf0f5'}
                        size={28}
                        onClick={this.addToList}
                    />
                </div>
                <div className={styles['flex-item']}>
                    <div className={styles['select-container']}>
                        <div className={styles['select-header']}>
                            <FlexBox role={'ui-flexbox'}>
                                <FlexBox.Item
                                    role={'ui-flexBox.item'}
                                    align={'left middle'}
                                >
                                    <label>
                                        {__('已选：')}
                                    </label>
                                </FlexBox.Item>
                                <FlexBox.Item
                                    role={'ui-flexbox.item'}
                                    align={'right middle'}
                                >
                                    <Button
                                        role={'sweetui-button'}
                                        theme={'text'}
                                        onClick={this.clearSelections}
                                        disabled={!selections.length}
                                    >
                                        {__('清空')}
                                    </Button>
                                </FlexBox.Item>
                            </FlexBox>
                        </div>
                        <div className={styles['select-content']}>
                            <ul>
                                {
                                    selections.map((item) => (
                                        <li
                                            key={item.id || ''}
                                            className={styles['dep-item']}
                                        >
                                            <div className={styles['seleted-data']} >
                                                <div className={styles['selection']} key={item.id}>
                                                    {
                                                        isFunction(formatSelectedItem) ?
                                                            formatSelectedItem(item)
                                                            :
                                                            <div
                                                                role={'ui-title'}
                                                                className={styles['selection-name']}
                                                                title={getDepName(item)}
                                                            >
                                                                <span>{item.name}</span>
                                                            </div>
                                                    }
                                                    <SweetIcon
                                                        role={'sweetui-sweeticon'}
                                                        name={'x'}
                                                        size={13}
                                                        className={styles['delete-icon']}
                                                        onClick={() => { this.deleteSelected(item) }}
                                                    />
                                                </div>
                                            </div>
                                        </li>
                                    ))
                                }
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}