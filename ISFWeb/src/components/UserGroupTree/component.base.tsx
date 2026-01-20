import { noop } from 'lodash';
import { getGroups } from '@/core/apis/console/usergroup';
import { ListTipStatus } from '../ListTipComponent/helper'
import { UserGroup, Result, formatUserGroup, Limit } from './helper';
import WebComponent from '../webcomponent';

interface UserGroupTreeProps {
    /**
     * 是否支持复选框
     */
    isMultSelect: boolean;

    /**
     * 是否禁用
     */
    disabled: boolean;

    /**
     * 选择搜索结果回调，传出点击的搜索结果数据
     */
    onRequestSelectSearchResult: (result: Result) => void;

    /**
     * 选择树节点回调，传出选中的树节点数据
     */
    onRequestSelectionsChange: (selections: ReadonlyArray<UserGroup>) => void;

}

interface UserGroupTreeState {
    /**
     * 用户组数据
     */
    userGroups: ReadonlyArray<UserGroup>;

    /**
     * 选择项
     */
    selections: ReadonlyArray<UserGroup>;

    /**
     * 加载状态
     */
    listTipStatus: ListTipStatus;

}

export default class UserGroupBase extends WebComponent<UserGroupTreeProps, UserGroupTreeState> {
    static defaultProps = {
        isMultSelect: true,
        disabled: false,
        onRequestSelectSearchResult: noop,
        onRequestSelectionsChange: noop,
    }

    state = {
        userGroups: [],
        selections: [],
        listTipStatus: ListTipStatus.Loading,
    }

    tree = null;

    async componentDidMount() {
        try {
            const { entries } = await getGroups({ offset: 0, limit: Limit })
            this.setState({
                userGroups: entries,
                listTipStatus: entries.length ? ListTipStatus.None : ListTipStatus.OrgEmpty,
            })
        } catch (ex) {
            this.setState({
                userGroups: [],
                listTipStatus: ListTipStatus.LoadFailed,
            })
        }
    }

    protected handleLazyLoad = async (page: number, limit: number): Promise<void> => {
        this.setState({
            userGroups: [
                ...this.state.userGroups,
                ...(await getGroups({ offset: (page - 1) * limit, limit })).entries,
            ],
        })
    }

    /**
     * 选择搜索结果回调
     */
    protected handleSelectResult = (result: UserGroup): void => {
        this.props.onRequestSelectSearchResult(formatUserGroup(result));
    }

    /**
     * 用户组树选择时回调
     */
    protected handleSelectionsChange = (selections: ReadonlyArray<UserGroup>): void => {
        const result = (selections as ReadonlyArray<UserGroup>).map((selection) => formatUserGroup(selection));
        this.props.onRequestSelectionsChange(result);
    }
}