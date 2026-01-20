import * as React from 'react';
import { noop, uniqBy, isEqual } from 'lodash';
import { getSystemProtectionLevel, ProtLevel } from '@/core/systemprotectionlevel'
import { NodeType } from '@/core/organization';
import WebComponent from '../webcomponent';
import { TabType, SelectionType, Selection, NodeData, filterTabType } from './helper';

export interface OrgAndGroupPickProps {
    /**
     * 该组件树是否用多选框
     */
    isMult: boolean;

    /**
     * 是否单选（该值需和isMult相反）
     */
    isSingleChoice?: boolean;

    /**
     * 显示哪些tab
     */
    tabType: ReadonlyArray<TabType>;

    /**
     * 节点类型,是否包含用户
     */
    nodeType: ReadonlyArray<NodeType>;

    /**
     * 是否禁用
     */
    disabled: boolean;

    /**
     * 占位符
     */
    placeholder?: string;

    /**
     * 是否只显示被禁用的用户
     */
    isShowDisabledUsers?: boolean;

    /**
     * 父组件传入的已选列表数据
     */
    selections: ReadonlyArray<Selection>;

    /**
     * 父组件传入的是否包含子部门
     */
    isIncludeSubDeps: boolean;

    /**
     * 是否显示(包含子部门)复选框
     */
    isShowCheckBox: boolean;

    /**
     * 是否显示padding，默认显示
     */
    isShowPadding: boolean;

     /**
     * 是否为知识库管理员
     */
     isKnowledger?: boolean;

     /**
      * 是否请求普通用户的组织架构
      */
     isRequestNormal?: boolean;

    /**
     * 传出selcetions
     */
    onRequestSelectionsChange: (selections: ReadonlyArray<Selection>) => void;

    /**
     * 传出isIncludeSubDeps
     */
    onRequestSubDepsChange: (isIncludeSubDeps: boolean) => void;

    /**
     * didMount事件回调
     */
    onDidMount?: () => void;
}

interface OrgAndGroupPickState {
    /**
     * 已选数据
     */
    selections: ReadonlyArray<Selection>;

    /**
     * 是否包含子部门
     */
    isIncludeSubDeps: boolean;

    /**
     * 显示哪些tab(需要过滤)
     */
    tabType: ReadonlyArray<TabType>;
}

export default class OrgAndGroupPickBase extends WebComponent<OrgAndGroupPickProps, OrgAndGroupPickState> {

    static defaultProps = {
        isMult: true,
        isSingleChoice: false,
        tabType: [TabType.Org, TabType.Group],
        nodeType: [NodeType.DEPARTMENT, NodeType.ORGANIZATION, NodeType.USER],
        disabled: false,
        selections: [],
        isIncludeSubDeps: false,
        isShowCheckBox: true,
        isShowPadding: true,
        onRequestSelectionsChange: noop,
        onRequestSubDepsChange: noop,
        onDidMount: noop,
    }

    state = {
        selections: this.props.selections || [],
        isIncludeSubDeps: this.props.isIncludeSubDeps || false,
        protLevel: ProtLevel.Common,
        tabType: [],
    }

    /**
     * departmentTree 的 ref
     */
    depTree = null;

    /**
     * userGroupTree 的 ref
     */
    grpTree = null;

    /**
     * AnonymousPick 的 ref
     */
    anonymousTree = null;

    /**
     * AppAcountTree 的 ref
     */
    AppAccountTree = null;

    async componentDidMount() {
        const level = await getSystemProtectionLevel()

        this.setState({
            tabType: this.props.isKnowledger ? [TabType.Org, TabType.Group] : filterTabType(level, this.props.tabType),
        }, () => {
            this.props.onDidMount()
        })
    }

    componentDidUpdate(prevProps, prevState) {
        if (
            this.props.isIncludeSubDeps !== this.state.isIncludeSubDeps
            && prevState.isIncludeSubDeps === this.state.isIncludeSubDeps
        ) {
            this.setState({
                isIncludeSubDeps: this.props.isIncludeSubDeps,
            })
        }

        if (
            !isEqual(this.props.selections, this.state.selections)
            && isEqual(prevState.selections, this.state.selections)
        ) {
            this.setState({
                selections: this.props.selections,
            })
        }

    }

    /**
     * 树无复选框时，将选择项加入选择列表(组织树)
     */
    protected addOrgSelection = (data: NodeData): void => {
        const item = {
            ...data,
            type: data.type === NodeType.USER ? SelectionType.User : SelectionType.Department,
        };

        this.setState(({ selections }) => ({
            selections: this.props.isSingleChoice ? [item] : uniqBy([...selections, item], 'id'),
        }), () => {
            this.props.onRequestSelectionsChange(this.state.selections);
        });
    };

    /**
     * 树无复选框时，将选择项加入选择列表(用户组树、匿名用户树、应用账户树)
     */
    protected addCommonSelections = (data: ReadonlyArray<Selection>): void => {
        this.setState(({ selections }) => ({
            selections: this.props.isSingleChoice ? data : uniqBy([...selections, ...data], 'id'),
        }), () => {
            this.props.onRequestSelectionsChange(this.state.selections);
        });
    }

    /**
     * 树有复选框时，点击添加箭头回调，将组件树的勾选项加入已选列表
     */
    protected addTreeDataToSelections = async (): void => {
        let type
        // 获取当前页面展示的tree ref，以便调用其public方法获取selections
        const { tree: treeRef } = [
            { tree: this.depTree },
            { tree: this.grpTree && this.grpTree.tree, selectionType: SelectionType.Group },
            { tree: this.anonymousTree, selectionType: SelectionType.Anonymous },
            { tree: this.AppAccountTree, selectionType: SelectionType.App },
        ].find(({ tree, selectionType }) => {
            if (tree) {
                type = selectionType
            }

            return tree
        }) || {}

        if (this.props.disabled || !treeRef) {
            return
        }

        const additions = (await treeRef.getSelections()).map(
            (item) => ({
                ...item,
                type: type || (item.type === NodeType.USER ? SelectionType.User : SelectionType.Department),
                id: item.id,
                name: item.name,
                original: this.depTree ? item.original : item,
            }),
        )

        this.setState(({ selections }) => ({
            selections: this.props.isSingleChoice ? additions : uniqBy([...selections, ...additions], 'id'),
        }), () => {
            treeRef.cancelSelections();

            this.props.onRequestSelectionsChange(this.state.selections);
        });
    }

    /**
     * 清空所有选择项
     */
    protected clearSelections = (): void => {
        const { selections } = this.state;

        if (selections.length) {
            this.setState({
                selections: [],
            }, () => {
                this.props.onRequestSelectionsChange(this.state.selections);
            });
        }
    };

    /**
     * 删除已选列表中的选中项
     */
    protected deleteSelection = (deleteItem: Selection): void => {
        if (this.props.disabled) {
            return
        }

        this.setState(({ selections }) => ({
            selections: selections.filter((item) => item.id !== deleteItem.id),
        }), () => {
            this.props.onRequestSelectionsChange(this.state.selections);
        });
    };

    /**
     * 当组件类型为部门时，选择是否包含子部门的勾选框 回调
     */
    protected checkSubDeps = (status: boolean): void => {
        this.setState({
            isIncludeSubDeps: status,
        }, () => {
            this.props.onRequestSubDepsChange(this.state.isIncludeSubDeps);
        })
    }
}