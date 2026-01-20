import * as React from 'react';
import classnames from 'classnames';
import { Text } from '@/ui/ui.desktop';
import { SweetIcon, Button, Select2, Trigger } from '@/sweet-ui';
import { NodeType } from '@/core/organization';
import OrgAndGroupPick from '../OrgAndGroupPick/component.view';
import { SelectionType } from '../OrgAndGroupPick/helper';
import { DocType, DocTypeText, ScopeType, DataItem, ConfigStatus } from './helper';
import SearchScopePickerBase from './component.base';
import styles from './styles.view';
import __ from './locale';

export default class SearchScopePicker extends SearchScopePickerBase {
    render() {
        const {
            open,
            scopeType,
        } = this.state

        return (
            <Trigger
                triggerEvent={'click'}
                anchorOrigin={['left', 'bottom']}
                alignOrigin={['left', -2]}
                freeze={false}
                onBeforePopupClose={this.closeAway}
                renderer={({ setPopupVisibleOnClick }) =>
                    <div key={'search-scope-trigger'} className={styles['search-scope-picker']}>
                        <div
                            className={classnames(styles['scope'], {
                                [styles['scope-active']]: open,
                            })}
                            onClick={() => { this.toggleShowPicker(); setPopupVisibleOnClick() }}
                        >
                            <span className={styles['text']}><Text role={'ui-text'}>{this.formatScope(scopeType, this.scope)}</Text></span>
                            <SweetIcon
                                role={'sweetui-sweeticon'}
                                className={styles['arrow']}
                                name={open ? 'arrowUp' : 'arrowDown'}
                                size={16}
                            />
                        </div>
                    </div>
                }
            >
                {
                    ({ close }) => this.renderPicker(close)
                }
            </Trigger>
        )
    }

    /**
     * 渲染选择弹框
     */
    private renderPicker = (close: () => void): JSX.Element => {
        const {
            selections,
            isContainSubDep,
        } = this.state

        const { docType } = this.props

        switch (this.state.configStatus) {
            case ConfigStatus.ScopePicker:
                return (
                    <div className={classnames(styles['picker'], styles['scope-picker'])}>
                        <div className={styles['tip-text']}>{__(`选择需查看的${DocTypeText[docType]}：`)}</div>
                        <Select2.Option
                            role={'sweetui-select2.option'}
                            onClick={() => this.changeScopeType(ScopeType.All, close)}
                        >
                            {this.formatScope(ScopeType.All)}
                        </Select2.Option>
                        <Select2.Option
                            role={'sweetui-select2.option'}
                            onClick={() => this.changeScopeType(ScopeType.Custom)}
                        >
                            {__('自定义选择')}
                            <SweetIcon
                                role={'sweetui-sweeticon'}
                                className={styles['arrow-right']}
                                name={'arrowRight'}
                                size={16}
                            />
                        </Select2.Option>
                    </div>
                )

            case ConfigStatus.CustomPicker:
                return (
                    <div className={classnames(styles['picker'], styles['custom-picker'])}>
                        <OrgAndGroupPick
                            isMult={docType === DocType.User}
                            tabType={this.props.tabType}
                            nodeType={docType === DocType.User ? [NodeType.ORGANIZATION, NodeType.DEPARTMENT, NodeType.USER] : [NodeType.ORGANIZATION, NodeType.DEPARTMENT]}
                            selections={selections}
                            isIncludeSubDeps={isContainSubDep}
                            onRequestSelectionsChange={this.selectionsChanged}
                            onRequestSubDepsChange={this.changeContainStatus}
                        />
                        <div className={styles['footer']}>
                            <Button
                                role={'sweetui-button'}
                                className={styles['button']}
                                disabled={!selections.length}
                                theme={'oem'}
                                width={56}
                                onClick={() => {
                                    close();
                                    this.confirm();
                                }}
                            >
                                {__('确定')}
                            </Button>
                            <Button
                                role={'sweetui-button'}
                                width={56}
                                onClick={() => { this.cancel(); close() }}
                            >
                                {__('取消')}
                            </Button>
                        </div>
                    </div>
                )

            default:
                return null
        }
    }

    /**
     * 渲染被选择列表
     */
    private formatScope = (scopeType: ScopeType, selections?: ReadonlyArray<DataItem>): string => {
        if (scopeType === ScopeType.All) {
            return this.props.docType === DocType.User ? __(`全部${DocTypeText[DocType.User]}`) : __(`全部${DocTypeText[DocType.Department]}`)
        } else {
            return selections && selections.map((item) => `${item.name}${item.type === SelectionType.Group ? __('（用户组）') : ''}`).join(', ')
        }
    }
}
