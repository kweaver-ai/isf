import * as React from 'react';
import { includes } from 'lodash';
import classnames from 'classnames';
import Panel from '@/ui/Panel/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { CascadeDirection, SelectType as NodeSelectType } from '@/ui/Tree2/ui.base';
import { NodeType, getDepName } from '@/core/organization';
import { ModalDialog2, SweetIcon, Button, CheckBox } from '@/sweet-ui';
import { isBrowser, Browser } from '@/util/browser';
import SearchDep from '../SearchDep/component.desktop';
import OrganizationTree from '../OrganizationTree/component.view';
import DepartmentTree from '../DepartmentTree/component.view';
import ListTipComponent from '../ListTipComponent/component.view'
import { ListTipStatus } from '../ListTipComponent/helper'
import OrganizationPickerBase from './component.base';
import styles from './styles.desktop';
import __ from './locale';
import * as deleteIcon from './assets/delete.png'

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class OrganizationPicker extends OrganizationPickerBase {

    private onCheckChange = (detail) => {
        if(this.props.userInfo) {
            if(detail) {
                const currentUser = this.state.data.find((item) => this.props.userInfo && item.id === this.props.userInfo.id)
                if(!currentUser) {
                    this.setState({
                        data: [
                            ...this.state.data,
                            {
                                ...this.props.userInfo,
                            },
                        ],
                    })
                }
            }else {
                const selects = this.state.data.filter((item) => item.id !== this.props.userInfo.id)
                this.setState({
                    data: selects,
                })
            }
        }
    }

    render() {
        const data = this.props.isSingleChoice ? this.state.data[0] : this.state.data

        const {
            isSingleChoice,
            isCascadeTree,
            title,
            tip,
            tipStyle,
            userid,
            selectType,
            isShowDisabledUsers,
            isShowUndistributed,
            extraRoots,
            isShowSetLoginUser,
        } = this.props

        const { isGetUsersLoading } = this.state

        const selectionIds = this.state.data.map((item) => item.id)

        return (
            <div className={styles['container']}>
                <ModalDialog2
                    role={'sweetui-modaldialog2'}
                    zIndex={this.props.zIndex || 50}
                    title={title}
                    width={isSingleChoice ? 420 : isCascadeTree ? 630 : 600}
                    icons={[{
                        icon: <SweetIcon role={'sweetui-sweeticon'} name="x" size={16} />,
                        onClick: this.cancelAddDep.bind(this),
                    },
                    ]}
                    buttons={[
                        {
                            text: __('确定'),
                            theme: 'oem',
                            onClick: this.confirmAddDep.bind(this),
                            disabled: !this.state.data.length,
                        },
                        {
                            text: __('取消'),
                            theme: 'regular',
                            onClick: this.cancelAddDep.bind(this),
                        },
                    ]}
                >
                    <Panel role={'ui-panel'}>
                        {
                            isSingleChoice ?
                                <div>
                                    <div className={styles['single-selete']}>
                                        <span className={styles['selected-text']}>{__('已选：')}</span>
                                        {
                                            data ?
                                                <div className={styles['selected']}>
                                                    <div
                                                        role={'ui-title'}
                                                        title={getDepName(data)}
                                                    >
                                                        <div className={classnames(
                                                            styles['dep-name'],
                                                            {
                                                                [styles['safari']]: isSafari,
                                                            },
                                                        )}>
                                                            {data.name}
                                                        </div>
                                                    </div>

                                                </div>
                                                :
                                                <span className={styles['no-content']}>
                                                    {'---'}
                                                </span>
                                        }
                                    </div>
                                    <SearchDep
                                        width={'100%'}
                                        selectType={selectType}
                                        isShowDisabledUsers={isShowDisabledUsers}
                                        isShowUndistributed={isShowUndistributed}
                                        onSelectDep={(value) => { this.selectDep(value) }}
                                    />
                                    <div className={styles['single-organization-tree']}>
                                        <OrganizationTree
                                            userid={userid}
                                            selectType={selectType}
                                            isShowDisabledUsers={isShowDisabledUsers}
                                            getNodeStatus={this.props.getNodeStatus}
                                            onSelectionChange={(value) => { this.selectDep(value) }}
                                        />
                                    </div>
                                </div>
                                :
                                <div>
                                    <div className={classnames(
                                        styles['content-row'],
                                        styles['row-margin'],
                                    )}>
                                        <div>
                                            <SearchDep
                                                placeholder={
                                                    (selectType.length === 1 && selectType[0] === NodeType.USER) // type只有用户
                                                        ? __('搜索用户')
                                                        : (
                                                            includes(selectType, NodeType.USER)
                                                                ? __('搜索用户或部门')
                                                                : __('搜索部门')
                                                        )
                                                }
                                                width={280}
                                                selectType={selectType}
                                                isShowDisabledUsers={isShowDisabledUsers}
                                                isShowUndistributed={isShowUndistributed}
                                                onSelectDep={(value) => { this.selectDep(value) }}
                                            />
                                        </div>
                                        {
                                            isCascadeTree
                                                ? <div className={styles['org-picker-blank']}></div>
                                                : null
                                        }
                                        <div className={styles['org-picker-clear']}>
                                            <span>{__('已选：')}</span>

                                            <Button
                                                role={'sweetui-button'}
                                                className={styles['clear-text']}
                                                theme={'text'}
                                                disabled={!this.state.data.length}
                                                onClick={this.clearSelectDep.bind(this)}
                                            >
                                                {__('清空')}
                                            </Button>
                                        </div>
                                    </div>
                                    <div className={styles['content-row']}>
                                        <div className={classnames(
                                            styles['org-picker-tree'],
                                            {
                                                [styles['cascade-tree']]: isCascadeTree,
                                            },
                                        )}>
                                            {
                                                isCascadeTree
                                                    ? (
                                                        <DepartmentTree
                                                            selectType={selectType}
                                                            nodeSelectType={NodeSelectType.CASCADE_MULTIPLE}
                                                            cascadeDirection={CascadeDirection.DOWN}
                                                            extraRoots={extraRoots}
                                                            isShowDisabledUsers={isShowDisabledUsers}
                                                            ref={(departmentTreeData) => this.departmentTreeData = departmentTreeData}
                                                        />
                                                    ) : (
                                                        <OrganizationTree
                                                            userid={userid}
                                                            selectType={selectType}
                                                            isShowDisabledUsers={isShowDisabledUsers}
                                                            isShowUndistributed={isShowUndistributed}
                                                            onSelectionChange={(value) => { this.selectDep(value) }}
                                                        />
                                                    )
                                            }

                                        </div>

                                        {
                                            isCascadeTree
                                                ? (
                                                    <div className={styles['org-picker-arrow']}>
                                                        <UIIcon
                                                            role={'ui-uiicon'}
                                                            size={28}
                                                            code={'\uf0f5'}
                                                            color={'#757575'}
                                                            onClick={this.addTreeData}
                                                        />
                                                    </div>
                                                )
                                                : <div className={styles['org-pick-gap']}></div>
                                        }

                                        <div className={styles['org-picker-selections']}>
                                            <ul className={styles['selections']}>
                                                {
                                                    this.state.data.map((sharer) => (
                                                        <li
                                                            key={sharer.id}
                                                            style={{ position: 'relative' }}
                                                            className={styles['selection']}
                                                        >
                                                            <div
                                                                role={'ui-title'}
                                                                className={styles['selection-name']}
                                                                title={getDepName(sharer)}
                                                            >
                                                                <span>{sharer.name}</span>
                                                            </div>
                                                            <UIIcon
                                                                role={'ui-uiicon'}
                                                                className={styles['icon-del']}
                                                                size={13}
                                                                code={'\uf014'}
                                                                fallback={deleteIcon}
                                                                onClick={() => { this.deleteSelectDep(sharer) }}
                                                            />
                                                        </li>
                                                    ))
                                                }
                                            </ul>
                                            {
                                                isGetUsersLoading ?
                                                    <div className={styles['loading']}>
                                                        <ListTipComponent
                                                            listTipStatus={ListTipStatus.Loading}
                                                            isInDialog={true}
                                                        />
                                                    </div>
                                                    : null
                                            }
                                        </div>
                                    </div>
                                </div>
                        }
                    </Panel>
                    {
                        isShowSetLoginUser ? (
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
                    {
                        tip
                            ? (
                                <div
                                    className={styles['tip']}
                                    style={tipStyle}
                                >
                                    {tip}
                                </div>
                            )
                            : null
                    }
                </ModalDialog2>
            </div>
        )
    }
}