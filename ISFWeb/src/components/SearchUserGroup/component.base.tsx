import * as React from 'react';
import { noop } from 'lodash';
import { searchInGroup } from '@/core/apis/console/usergroup';
import WebComponent from '../webcomponent';

/**
 * 用户组信息
 */
interface UserGroup {
    /**
     * 用户组id
     */
    id: string;

    /**
     * 用户组名称
     */
    name: string;
}

interface Props {
    /**
     * 选择事件
     */
    onRequestSelect: (userGroup: UserGroup) => void;

    /**
     * 宽度
     */
    width: string | number;

    /**
     * 检测 输入框的值是否发生改变
     */
    onValueChange: (value: string) => boolean;

    /**
     * 是否搜索框能够输入
     */
    disabled?: boolean;
}

interface State {
    /**
     * 搜索结果
     */
    results: Array<{ type: string; id: string; name: string; origin: any }>;

    /**
     * 搜索关键字
     */
    searchKey: string;
}

export default class SearchUserGroupBase extends WebComponent<any, any> {
    static defaultProps = {
        onRequestSelect: noop,
    }

    state = {
        results: [],

        searchKey: '',
    }

    autocomplete

    componentDidMount() {
        this.setState({
            searchKey: this.props.value,
        })
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        if (nextProps.canInput !== prevState.canInputValue) {
            return {
                canInputValue: nextProps.canInput,
            }
        }
        return null
    }

    /**
     * 根据key获取用户组
     * @return 用户组数组
     */
    protected getGroupsByKey = async (key: string, offset: number = 0, limit = 10): Promise<ReadonlyArray<Record<string, any>>> => {
        if (key) {
            return (await searchInGroup({ keyword: key, type: 'group', offset, limit })).groups.entries;
        } else {
            return [];
        }

    }

    /**
     * 获取搜索到的结果
     * @param results 用户组数组
     */
    protected getSearchData = (results: ReadonlyArray<Record<string, any>>): void => {
        this.setState({
            results,
        })
    }

    /**
     * 选择搜索到的单个用户组
     */
    protected selectItem = (userGroup: UserGroup): void => {
        this.props.onRequestSelect(userGroup);

        if (this.props.completeSearchKey && userGroup.name) {
            this.setState({
                searchKey: userGroup.name,
            })
        } else if (this.props.completeSearchKey && userGroup.name) {
            this.setState({
                searchKey: userGroup.name,
            })
        }
        this.autocomplete.toggleActive(false);
    }

    protected handelChange = (searchKey): void => {
        this.setState({
            searchKey,
        })
        this.props.onValueChange ? this.props.onValueChange(searchKey) : null
    }

    /**
     * 按下enter
     */
    protected handleEnter = (e, selectIndex: number): void => {
        if (selectIndex >= 0) {
            this.selectItem(this.state.results[selectIndex])
        }
    }

    /**
     * 懒加载
     */
    protected lazyLoad = async (page: number, limit: number): Promise<void> => {
        this.setState({
            results: [
                ...this.state.results,
                ...(await this.getGroupsByKey(this.state.searchKey, (page - 1) * limit)),
            ],
        })
    }
}
