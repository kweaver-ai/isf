import * as React from 'react';
import classnames from 'classnames';
import Tabs from '@/ui/Tabs/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import Title from '@/ui/Title/ui.desktop';
import { UIIcon } from '@/ui/ui.desktop';
import { isBrowser, Browser } from '@/util/browser';
import { ModalDialog2, SweetIcon, Button, CheckBox } from '@/sweet-ui';
import { getDepName } from '@/core/organization'
import SearchDep from '../SearchDep/component.desktop';
import OrganizationTree from '../OrganizationTree/component.view';
import AppAccountTree from '../AppAccountTree/component.view'
import OrgAndAccountPickBase from './component.base';
import { TabType, SelectionType } from './helper';
import __ from './locale';
import styles from './styles.desktop'

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

const tabText = {
    [TabType.Org]: '组织结构',
    [TabType.AppAccount]: '应用账户',
}

export default class OrgAndAccountPick extends OrgAndAccountPickBase {

    private onCheckChange = (detail) => {
        if(this.props.userInfo) {
            if(detail) {
                const currentUser = this.state.selections.find((item) => this.props.userInfo && item.id === this.props.userInfo.id)
                if(!currentUser) {
                    this.setState({
                        selections: [
                            ...this.state.selections,
                            {
                                ...this.props.userInfo,
                            },
                        ],
                    })
                }
            }else {
                const selects = this.state.selections.filter((item) => item.id !== this.props.userInfo.id)
                this.setState({
                    selections: selects,
                })
            }
        }
    }

    render() {
        const {
            title,
            isSingleChoice,
            tabType,
            isShowDisabledUsers,
            isShowSetLoginUser,
        } = this.props

        const { selections } = this.state

        const selectionIds = selections.map((item) => item.id)

        return (
            <div className={styles['container']}>
                <ModalDialog2
                    role={'sweetui-modaldialog2'}
                    zIndex={this.props.zIndex || 50}
                    title={title}
                    width={isSingleChoice ? 420 : 600}
                    icons={[{
                        icon: <SweetIcon role={'sweetui-sweeticon'} name="x" size={16} />,
                        onClick: this.cancelAddSelection,
                    }]}
                    buttons={[
                        {
                            text: __('确定'),
                            theme: 'oem',
                            onClick: this.confirmAddSelection,
                            disabled: !selections.length,
                        },
                        {
                            text: __('取消'),
                            theme: 'regular',
                            onClick: this.cancelAddSelection,
                        },
                    ]}
                >
                    <Panel role={'ui-panel'}>
                        {
                            isSingleChoice
                                ? (
                                    <div className={styles['single-selete']}>
                                        <span className={styles['selected-text']}>{__('已选：')}</span>
                                        {
                                            selections.length
                                                ? (
                                                    <div className={styles['selected']}>
                                                        <Title
                                                            role={'ui-title'}
                                                            content={getDepName(selections[0])}
                                                        >
                                                            <div className={classnames(
                                                                styles['dep-name'],
                                                                {
                                                                    [styles['safari']]: isSafari,
                                                                },
                                                            )}>
                                                                {
                                                                    selections[0].type === SelectionType.AppAccount
                                                                        ? selections[0].name + __('（应用账户）')
                                                                        : selections[0].name
                                                                }
                                                            </div>
                                                        </Title>
                                                    </div>
                                                )
                                                : (
                                                    <span className={styles['no-content']}>
                                                        {'---'}
                                                    </span>
                                                )
                                        }
                                    </div>
                                )
                                :
                                null
                        }
                        <div className={styles['user-wrp']}>
                            {/* 左边Tab栏 */}
                            <div className={classnames(
                                styles['tree-wrp'],
                                { [styles['tree-wrp-flex']]: isSingleChoice },
                            )}>
                                {
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

                                                        case TabType.AppAccount:
                                                            return (
                                                                <Tabs.Content key={'org-cnt'}>
                                                                    {this.renderAccountTree()}
                                                                </Tabs.Content>
                                                            )
                                                    }
                                                })
                                            }
                                        </Tabs.Main>
                                    </Tabs>
                                }
                            </div>
                            {/* 右边 */}
                            {
                                !isSingleChoice
                                    ? (
                                        <div className={styles['list']}>
                                            <div className={styles['clear']}>
                                                <span>{__('已选：')}</span>
                                                <Button
                                                    role={'sweetui-button'}
                                                    className={styles['btn-clear']}
                                                    theme={'text'}
                                                    disabled={!selections.length}
                                                    onClick={this.clearSelections}
                                                >
                                                    {__('清空')}
                                                </Button>
                                            </div>
                                            <div className={styles['list-wrp']}>
                                                <ul>
                                                    {
                                                        selections.length
                                                            ? (
                                                                selections.map((item) => (
                                                                    <li
                                                                        key={item.id}
                                                                        className={styles['selection']}
                                                                    >
                                                                        <Title
                                                                            role={'ui-title'}
                                                                            key={item.id}
                                                                            content={getDepName(item)}
                                                                        >
                                                                            <span className={styles['selection-name']}>
                                                                                {`${item.name}${item.type == SelectionType.AppAccount ? __('（应用账户）') : ''}`}
                                                                            </span>

                                                                        </Title>
                                                                        <UIIcon
                                                                            role={'ui-uiicon'}
                                                                            className={styles['icon']}
                                                                            size={13}
                                                                            code={'\uf014'}
                                                                            onClick={() => { this.deleteSelection(item) }}
                                                                        />
                                                                    </li>
                                                                ))
                                                            )
                                                            : null
                                                    }
                                                </ul>
                                            </div>
                                        </div>
                                    )
                                    : null
                            }
                        </div>
                    </Panel>
                    {
                        isShowSetLoginUser && this.props.userInfo ? (
                            <div className={styles['login-user']}>
                                <div className={styles['login-user-checkbox']}>
                                    <CheckBox
                                        checked={selectionIds.includes(this.props.userInfo.id)}
                                        onCheckedChange={({ detail }) => this.onCheckChange(detail)}
                                    />
                                </div>
                                {__('设置当前登陆用户为文档库所有者')}
                            </div>
                        ) : null
                    }
                </ModalDialog2>
            </div>
        )
    }

    // 应用账户
    private renderAccountTree = (): JSX.Element => {
        return (
            <div className={styles['tree-cnt']}>
                <AppAccountTree
                    onRequestSelection={this.addSelections(TabType.AppAccount)}
                />
            </div>
        )
    }

    // 组织结构
    private renderOrgTree = (): JSX.Element => {
        const { userid, selectType, isShowDisabledUsers } = this.props

        return (
            <div className={styles['tree-cnt']}>
                <SearchDep
                    placeholder={
                        __('搜索用户')
                    }
                    width={'100%'}
                    selectType={selectType}
                    isShowDisabledUsers={isShowDisabledUsers}
                    onSelectDep={this.addSelections(TabType.Org)}
                />
                <div className={styles['tree-main']}>
                    <OrganizationTree
                        userid={userid}
                        selectType={selectType}
                        isShowDisabledUsers={isShowDisabledUsers}
                        onSelectionChange={this.addSelections(TabType.Org)}
                    />
                </div>
            </div>
        )
    }
}