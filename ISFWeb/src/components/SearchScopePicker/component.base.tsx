import * as React from 'react';
import { noop, isFunction } from 'lodash';
import { TabType } from '../OrgAndGroupPick/helper';
import { ConfigStatus, DocType, ScopeType, DataItem } from './helper';

interface SearchScopePickerProps {
    /**
     * 操作类型（个人/部门）
     */
    docType: DocType;

    /**
     * tab类型
     */
    tabType: ReadonlyArray<TabType>;

    /**
     * 是否为重置状态
     */
    isReseted: boolean;

    /**
     * 抛出选择范围
     */
    onRequestConfirm: (data: {
        is_all_checked: boolean;
        is_contain_child?: boolean;
        scope: ReadonlyArray<DataItem>;
    }) => void;
}

interface SearchScopePickerState {
    /**
     * 展开状态
     */
    open: boolean;

    /**
     * 显示内容
     */
    configStatus: ConfigStatus;

    /**
     * 筛选范围类型
     */
    scopeType: ScopeType;

    /**
     * 被选择列表
     */
    selections: ReadonlyArray<DataItem>;

    /**
     * 是否包含子部门
     */
    isContainSubDep: boolean;
}

export default class SearchScopePickerBase extends React.PureComponent<SearchScopePickerProps, SearchScopePickerState> {
    static defaultProps = {
        tabType: [TabType.Org, TabType.Group],
        isReseted: false,
        docType: DocType.User,
        onRequestConfirm: noop,
    }

    state = {
        open: false,
        configStatus: ConfigStatus.None,
        scopeType: ScopeType.All,
        selections: [],
        isContainSubDep: false,
    }

    /**
     * 点击确定之后的范围
     */
    scope = []

    /**
     * 点击确定之后是否包含子部门
     */
    containSubDepChecked: boolean = false

    componentDidUpdate({ isReseted }, { scopeType }) {
        if (isReseted !== this.props.isReseted && !!this.props.isReseted && scopeType !== ScopeType.All) {
            this.reset()
        }
    }

    /**
     * 展开/收起下拉框
     */
    protected toggleShowPicker = (): void => {
        const { configStatus, open } = this.state

        this.setState({
            open: !open,
            configStatus: configStatus === ConfigStatus.None ? ConfigStatus.ScopePicker : configStatus,
        })
    }

    /**
     * 点击空白关闭
     */
    protected closeAway = (): void => {
        this.setState({
            open: false,
        })
    }

    /**
     * 更改筛选范围类型
     */
    protected changeScopeType = (scopeType: ScopeType, close?: () => void): void => {
        if (scopeType === ScopeType.All) {
            this.reset()

            this.props.onRequestConfirm({ is_all_checked: true, scope: [] })

            isFunction(close) && close()
        } else {
            this.setState({
                configStatus: ConfigStatus.CustomPicker,
            })
        }
    }

    /**
     * 已选列表变化
     */
    protected selectionsChanged = (selections: ReadonlyArray<DataItem>): void => {
        this.setState({
            selections,
        })
    }

    /**
     * 是否包含子部门变化
     */
    protected changeContainStatus = (isContainSubDep: boolean): void => {
        this.setState({
            isContainSubDep,
        })
    }

    /**
     * 确定按钮回调
     */
    protected confirm = (): void => {
        const { selections, isContainSubDep } = this.state

        this.scope = selections

        this.containSubDepChecked = isContainSubDep

        this.setState({
            open: false,
            configStatus: ConfigStatus.None,
            scopeType: ScopeType.Custom,
        })

        this.props.onRequestConfirm({ is_all_checked: false, is_contain_child: isContainSubDep, scope: selections })
    }

    /**
     * 取消
     */
    protected cancel = () => {
        this.setState({
            open: false,
            configStatus: ConfigStatus.None,
            selections: this.scope,
            isContainSubDep: this.containSubDepChecked,
        })
    }

    /**
     * 重置
     */
    private reset = (): void => {
        this.setState({
            open: false,
            configStatus: ConfigStatus.None,
            scopeType: ScopeType.All,
            selections: [],
            isContainSubDep: false,
        })

        this.scope = []

        this.containSubDepChecked = false
    }
}
