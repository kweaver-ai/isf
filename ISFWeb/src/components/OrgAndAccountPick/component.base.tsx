import * as React from 'react';
import { uniqBy } from 'lodash';
import { Doclibs } from '@/core/doclibs/doclibs';
import { NodeType, nodeTypeMaptoMixinType } from '@/core/organization';
import WebComponent from '../webcomponent';
import { TabType, Selection, SelectionType } from './helper';

interface Props {

    /**
     * zIndex
     */
    zIndex?: number;

    /**
     * dialog 头信息
     */
    title: string;

    /**
     * 是否为单选
     */
    isSingleChoice: boolean;

    /**
     * 已选中
     */
    selected?: ReadonlyArray<Doclibs.UserInfo>;

    /**
     * 用户类型(组织结构搜索)
     */
    selectType?: Array<NodeType>;

    /**
     * 是否只显示被禁用的用户
     */
    isShowDisabledUsers?: boolean;

    /**
     * 当前管理员id
     */
    userid: string;

    /**
     * tab栏选项
     */
    tabType: ReadonlyArray<TabType>;

    /**
     * 点击取消事件
     */
    onRequestCancel?: () => void;

    /**
     * 点击确定事件
     */
    onRequestConfirm?: (data: ReadonlyArray<Doclibs.UserInfo>) => void;

    /**
     * 数据转换内部数据结构
     */
    converterIn?: (x) => Selection;

    /**
     * 数据转换外部数据结构
     */
    convererOut?: (Node) => any;

    /**
     * 是否显示设置当前登陆用户为文档库所有者
     */
    isShowSetLoginUser?: boolean;

    /**
     * 用户信息
     */
    userInfo?: {name: string; type: SelectionType; id: string};
}

interface State {

    /**
     * 选择添加的库所有者
     */
    selections: ReadonlyArray<Selection>;
}

export default class OrgAndAccountPickBase extends WebComponent<Props, State>{

    static defaultProps = {
        tabType: [TabType.Org, TabType.AppAccount],
        isSingleChoice: false,
        selected: [],
    }

    state = {
        selections: [],
    }

    componentDidMount() {
        // 打开时向state中添加已选的数据
        const { selected: [selection] } = this.props

        if (selection) {
            this.setState({
                selections: this.state.selections.concat(selection),
            })
        }
    }

    /**
     * 将选项加入已选列表
     * 接收两种数据： 应用账户/组织结构用户
     */
    protected addSelections = (selectType: TabType) => {
        return async (selection): Promise<void> => {
            const { isSingleChoice } = this.props
            let additions: Selection

            // 应用账户数据处理
            if (selectType === TabType.AppAccount) {
                additions = {
                    ...selection,
                    type: SelectionType.AppAccount,
                    original: selection,
                }
            } else {
                additions = {
                    ...selection,
                    type: nodeTypeMaptoMixinType(selection.type),
                }
            }

            // 数据处理完毕，将选项添加进state   情境：单选/多选
            this.setState(({ selections }) => ({
                selections: isSingleChoice ? [additions] : uniqBy([...selections, additions], 'id'),
            }))
        }
    }

    /**
     * 移除选项
     */
    protected deleteSelection = (deleteItem: Selection): void => {

        this.setState(({ selections }) => ({
            selections: selections.filter((item) => item.id !== deleteItem.id),
        }))
    }

    /**
     * 清空所有选项
     */
    protected clearSelections = (): void => {
        const { selections } = this.state;

        if (selections.length) {
            this.setState({
                selections: [],
            });
        }
    }

    /**
     * 取消本次操作
     */
    protected cancelAddSelection = (): void => {
        this.clearSelections()
        this.props.onRequestCancel()
    }

    /**
     * 确定本次操作
     */
    protected confirmAddSelection = (): void => {
        this.props.onRequestConfirm(this.state.selections.map(this.props.convererOut))
    }
}