import * as React from 'react';
import { trim, noop } from 'lodash';
import { getUseAccountMgnt } from '@/core/apis/console/useaccountmgnt';
import WebComponent from '../webcomponent';
import { ListTipStatus } from '../ListTipComponent/helper';
import { NodeData } from '../OrgAndAccountPick/helper';
import { AppAccount } from './type';

interface Props {
    /**
     * 该组件树是否用多选框
     */
    isMult?: boolean;

    /**
     * 是否禁用
     */
    disabled?: boolean;

    /**
     * 选中节点时触发
     */
    onRequestSelection?: (node: NodeData) => void;
}

interface State {
    /**
     * 应用账户数据列
     */
    data: ReadonlyArray<AppAccount>;

    /**
     * 选中项的id
     */
    selectionId: string;

    /**
     * 搜索关键字
     */
    searchKey: string;

    /**
     * 加载状态
     */
    listTipStatus: ListTipStatus;

    /**
     * isMulti为true时，选中项
     */
    selections: ReadonlyArray<AppAccount>;
}

export default class AppAccountTreeBase extends WebComponent<Props, State>{
    static defaultProps = {
        isMulti: false,
        disabled: false,
        onRequestSelection: noop,
    }

    state = {
        data: [],
        selectionId: '',
        searchKey: '',
        listTipStatus: ListTipStatus.Loading,
        selections: [],
    }

    protected tree = null;

    lazyLoaderRef = null; // 获取懒加载组件

    componentDidMount() {
        this.loader().then(this.handleLoadSuccess, this.handleLoadFailed)
    }

    /**
     * 应用账户数据懒加载
     */
    protected handleLazyLoad = async (page: number, limit: number): Promise<void> => {
        try {
            const { entries } = await this.loader({ limit, offset: (page - 1) * limit })

            this.setState({
                data: [
                    ...this.state.data,
                    ...entries,
                ],
            })

        } catch (ex) { }
    }

    /**
     * 改变搜索关键字
     */
    protected changeSearchKey = (searchKey: string): void => {
        this.setState({
            searchKey,
        })
    }

    /**
     * 点击进行添加
     */
    protected addSelection = (selection: AppAccount): void => {
        if (selection) {
            this.setState({
                selectionId: selection.id,
            })

            this.props.onRequestSelection(selection)
        }
    }

    /**
     * 请求的函数
     */
    protected loader = ({ limit = 150, offset = 0 } = {}): Promise<{ entries: ReadonlyArray<AppAccount> }> => {
        return getUseAccountMgnt({
            limit,
            offset,
            direction: 'desc',
            sort: 'date_created',
            keyword: trim(this.state.searchKey),
        })
    }

    /**
     * 处理列表加载成功
     */
    protected handleLoadSuccess = ({ entries }: { entries: ReadonlyArray<AppAccount> }): void => {
        this.setState({
            data: entries,
            listTipStatus: entries.length
                ? ListTipStatus.None
                : trim(this.state.searchKey)
                    ? ListTipStatus.NoSearchResults
                    : ListTipStatus.OrgEmpty,
        })

        this.lazyLoaderRef && this.lazyLoaderRef.reset() // 搜索之后 必须要执行，否则 滚动条位置不会复位、懒加载不会再次被触发
    }

    /**
     * 处理列表加载失败
     */
    protected handleLoadFailed = (): void => {
        this.setState({
            data: [],
            listTipStatus: ListTipStatus.LoadFailed,
        })
    }

    /**
     * 多选
     */
    protected multiSelect(selection: AppAccount, checked: boolean) {
        const { selections } = this.state

        this.setState({
            selections: checked ? [...selections, selection] : selections.filter(({ id }) => id !== selection.id),
        })
    }

    /**
     * public 获取选中项（isMulti为true时使用）
     */
    public getSelections(): ReadonlyArray<AppAccount> {
        return this.state.selections
    }

    /**
     * public 清空选中项（isMulti为true时使用）
     */
    public cancelSelections() {
        this.setState({
            selections: [],
        })
    }
}