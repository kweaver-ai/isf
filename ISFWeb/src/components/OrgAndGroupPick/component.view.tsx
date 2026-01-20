import * as React from 'react';
import classnames from 'classnames';
import Tabs from '@/ui/Tabs/ui.desktop';
import { UIIcon } from '@/ui/ui.desktop';
import { CheckBox, Button } from '@/sweet-ui';
import { CascadeDirection, SelectType as NodeSelectType } from '@/ui/Tree2/ui.base';
import { getDepName } from '@/core/organization'
import SearchDep from '../SearchDep/component.desktop';
import OrganizationTree from '../OrganizationTree/component.view';
import DepartmentTree from '../DepartmentTree/component.view';
import UserGroupTree from '../UserGroupTree/component.view';
import AppAccountTree from '../AppAccountTree/component.view'
import AnonymousPick from './AnonymousPick/component.view';
import OrgAndGroupPickBase from './component.base';
import { TabType, SelectionType } from './helper';
import styles from './styles.desktop';
import __ from './locale';

const tabText = {
    [TabType.Org]: '组织结构',
    [TabType.Group]: '用户组',
    [TabType.Anonymous]: '匿名用户',
    [TabType.App]: '应用账户',
}

export default class OrgAndGroupPick extends OrgAndGroupPickBase {
    render() {
        const { isMult, disabled, isShowCheckBox, isShowPadding } = this.props;
        const { selections, isIncludeSubDeps, tabType } = this.state;

        return (
            <div
                className={classnames(
                    styles['user-wrp'],
                    {
                        [styles['show-padding']]: isShowPadding,
                    },
                )}
            >
                <div className={styles['tree-wrp']}>
                    {
                        tabType.length
                            ? (
                                <Tabs role={'ui-tabs'}>
                                    <Tabs.Navigator role={'ui-tabs.navigator'} key={'navigator'}>
                                        {
                                            tabType.map((tab) =>
                                                <Tabs.Tab role={'ui-tabs.tab'} key={tab}>{__(tabText[tab])}</Tabs.Tab>,
                                            )
                                        }
                                    </Tabs.Navigator>
                                    <Tabs.Main role={'ui-tabs.main'} key={'main'}>
                                        {
                                            tabType.map((tab) => {
                                                switch (tab) {
                                                    case TabType.Org:
                                                        return (
                                                            <Tabs.Content key={'org-cnt'}>
                                                                {this.renderOrgTree()}
                                                            </Tabs.Content>
                                                        )

                                                    case TabType.Group:
                                                        return (
                                                            <Tabs.Content key={'org-cnt'}>
                                                                {this.renderGrpTree()}
                                                            </Tabs.Content>
                                                        )

                                                    case TabType.Anonymous:
                                                        return (
                                                            <Tabs.Content key={'org-cnt'}>
                                                                {this.renderAnonymousTree()}
                                                            </Tabs.Content>
                                                        )

                                                    case TabType.App:
                                                        return (
                                                            <Tabs.Content key={'org-cnt'}>
                                                                {this.renderAppTree()}
                                                            </Tabs.Content>
                                                        )
                                                }
                                            })
                                        }
                                    </Tabs.Main>
                                </Tabs>
                            )
                            : null
                    }
                </div>

                {
                    this.props.isMult
                        ? (
                            <div
                                className={classnames(
                                    styles['add'],
                                )}
                            >
                                <UIIcon
                                    role={'ui-uiicon'}
                                    className={styles['btn']}
                                    size={28}
                                    code={'\uf0f5'}
                                    color={disabled ? '#c0c0c0' : ''}
                                    onClick={this.addTreeDataToSelections}
                                />
                            </div>
                        ) : null
                }
                <div className={styles['list']}>
                    <div className={styles['clear']}>
                        <span>{__('已选：')}</span>
                        <Button
                            role={'sweetui-button'}
                            className={styles['btn-clear']}
                            theme={'text'}
                            disabled={disabled || !selections.length}
                            onClick={this.clearSelections}
                        >
                            {__('清空')}
                        </Button>
                    </div>

                    <div
                        className={classnames(
                            styles['list-wrp'],
                            {
                                [styles['disabled']]: disabled,
                            },
                        )}
                    >
                        <div
                            className={classnames(
                                styles['itms'],
                                {
                                    [styles['dep-itms']]: !isMult,
                                },
                            )}
                        >
                            {selections.length
                                ? selections.map((item, index) => (
                                    <div className={styles['selection']} key={item.id}>
                                        <div
                                            role={'ui-title'}
                                            title={getDepName(item)}
                                        >
                                            <div className={classnames(
                                                styles['selection-name'],
                                            )}>
                                                <span>
                                                    {item.name}
                                                    {
                                                        [SelectionType.Group, SelectionType.App].includes(item.type)
                                                            ? `${__('（${text}）', { text: tabText[item.type] })}`
                                                            : ''
                                                    }
                                                </span>
                                            </div>
                                        </div>
                                        <UIIcon
                                            role={'ui-uiicon'}
                                            size={13}
                                            code={'\uf014'}
                                            onClick={() => { this.deleteSelection(item) }}
                                        />
                                    </div>
                                ))
                                : null
                            }
                        </div>
                        {
                            isShowCheckBox && !isMult
                                ? (
                                    <div className={styles['cfg']}>
                                        <CheckBox
                                            role={'sweetui-checkbox'}
                                            disabled={disabled}
                                            checked={isIncludeSubDeps}
                                            onCheckedChange={({ detail }) => {
                                                this.checkSubDeps(detail)
                                            }}
                                        >
                                            <span
                                                className={classnames(
                                                    {
                                                        [styles['disabled']]: disabled,
                                                    },
                                                )}
                                            >
                                                {__('包含子部门')}
                                            </span>
                                        </CheckBox>
                                        <UIIcon
                                            role={'ui-uiicon'}
                                            className={classnames(
                                                {
                                                    [styles['disabled']]: disabled,
                                                },
                                            )}
                                            code={'\uf055'}
                                            size={'16px'}
                                            title={
                                                <div className={styles['text']} >
                                                    {__('若勾选此项，则包含已选部门的子部门，不包含用户组部门中的子部门。')}
                                                </div>
                                            }
                                        />
                                    </div>
                                )
                                : null
                        }
                    </div>
                </div>
            </div>
        )
    }

    renderOrgTree() {
        const { isMult, disabled, nodeType, placeholder, isShowDisabledUsers, isRequestNormal } = this.props;

        return (
            <div className={styles['tree-cnt']}>
                <SearchDep
                    canInput={!disabled}
                    placeholder={
                        placeholder ? placeholder : !isMult ? __('搜索部门') : __('搜索用户或部门')
                    }
                    onSelectDep={(value) => { this.addOrgSelection(value) }}
                    width={'100%'}
                    selectType={nodeType}
                    isShowDisabledUsers={isShowDisabledUsers}
                />

                <div
                    className={classnames(
                        styles['tree'],
                        {
                            [styles['mult-tree']]: isMult,
                            [styles['disabled']]: disabled,
                        },
                    )}
                >
                    {
                        !isMult
                            ? (
                                <OrganizationTree
                                    disabled={disabled}
                                    selectType={nodeType}
                                    isShowDisabledUsers={isShowDisabledUsers}
                                    onSelectionChange={(value) => { this.addOrgSelection(value) }}
                                    isRequestNormal={isRequestNormal}
                                />
                            )
                            : (
                                <DepartmentTree
                                    disabled={disabled}
                                    selectType={nodeType}
                                    nodeSelectType={NodeSelectType.CASCADE_MULTIPLE}
                                    cascadeDirection={CascadeDirection.DOWN}
                                    isShowDisabledUsers={isShowDisabledUsers}
                                    ref={(depTree) => this.depTree = depTree}
                                    isRequestNormal={isRequestNormal}
                                />
                            )
                    }
                </div>
            </div>
        )
    }

    renderGrpTree() {
        const { disabled, isMult } = this.props;

        return (
            <div className={styles['tree-cnt']}>
                <UserGroupTree
                    disabled={disabled}
                    isMultSelect={isMult}
                    ref={(grpTree) => this.grpTree = grpTree}
                    onRequestSelectionsChange={(value) => { isMult ? null : this.addCommonSelections(value) }}
                    onRequestSelectSearchResult={(value) => { this.addCommonSelections(value) }}
                />
            </div>
        )
    }

    renderAnonymousTree() {
        const { disabled, isMult } = this.props;
        return (
            <div className={styles['tree-cnt']}>
                <AnonymousPick
                    ref={(ref) => { this.anonymousTree = ref }}
                    isMult={isMult}
                    disabled={disabled}
                    onRequsetSelection={(selection) => isMult ? null : this.addCommonSelections(selection)}
                />
            </div>
        )
    }

    renderAppTree() {
        const { disabled, isMult } = this.props;

        return (
            <div className={styles['tree-cnt']}>
                <AppAccountTree
                    ref={(ref) => { this.AppAccountTree = ref }}
                    isMult={isMult}
                    disabled={disabled}
                    onRequestSelection={(selection) => isMult ? null : this.addCommonSelections([{ ...selection, type: SelectionType.App }])}
                />
            </div>
        )
    }
}