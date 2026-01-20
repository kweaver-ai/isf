import React from 'react'
import { noop, filter, trim, isEmpty, isEqual } from 'lodash'
import __ from './locale'

let index = 0;
export default class ComboSearchBoxBase extends React.PureComponent<UI.ComboSearchBox.Props, UI.ComboSearchBox.State> {
    static defaultProps = {
        keys: [],
        placeholder: __('搜索'),
        onComboChange: noop,
        renderOption: noop,
        renderComboItem: noop,
    }

    state = {
        value: '',
        searchAnchor: null,
        searchValue: [],
        isSearchFocus: false,
        isSearchMenu: false,
        filterInputValue: '',
    }

    timeout = null;

    firstFilterInput: HTMLInputElement;

    searchInput: HTMLInputElement;

    /**
     * 输入块前置输入框实例集合
     */
    searchFilterInput = [];

    /**
     * 当接收到的searchValue发生改变的时，更新显示
     */
    componentDidUpdate({ searchValue, searchKey }, prevState) {
        if (!isEqual(searchValue, this.props.searchValue) || searchKey !== this.props.searchKey) {
            if (!isEqual(this.props.searchValue, prevState.searchValue)) {
                this.setState({
                    searchValue: this.props.searchValue,
                })
            }
            if (this.props.searchKey !== prevState.value) {
                this.setState({
                    value: this.props.searchKey,
                })
            }
        }
    }
    /**
     * 搜素框获得焦点时触发,获取聚焦元素并设置聚焦状态为true
     */
    protected handleSearchBoxFocus(e) {
        if (this.timeout) {
            clearTimeout(this.timeout);
        }
        this.setState({
            searchAnchor: e.currentTarget,
            isSearchMenu: true,
            isSearchFocus: true,
        })
    }

    /**
     * 搜索框失去焦点时触发，聚焦状态设置为false
     */
    protected handleSearchBoxBlur() {
        // 如果当面面板是打开状态
        if (this.timeout) {
            clearTimeout(this.timeout);
        }

        this.timeout = setTimeout(() => {
            if (this.state.value) {
                this.setState({
                    isSearchFocus: false,
                    isSearchMenu: false,
                })
            } else {
                this.setState({
                    isSearchFocus: false,
                })
            }

        }, 200)
    }

    /**
     * 点击删除每一项搜索关键字时触发
     * @param {any} key 删除项的信息
     */
    protected handleItemDelete(e, key, index) {
        // 筛选出非删除项
        const newValue = filter([...this.state.searchValue], (item) => {
            return item['index'] !== key['index']
        })

        this.setState({
            searchValue: newValue,
        }, () => {
            this.props.onComboChange(this.state.searchValue, this.state.value)
            this.searchFilterInput = [...this.searchFilterInput.slice(0, index), ...this.searchFilterInput.slice(index + 1, this.searchFilterInput.length)]
        })
    }

    /**
     * 监听搜索框变化
     */
    protected handleSearchInputChange(value) {
        this.setState({
            value,
        })
    }

    /**
     * 当输入框中有键按下时
     */
    protected handleInputKeyDown(e) {
        const { keys } = this.props
        const { value, searchValue } = this.state
        // 当按下删除键，搜索框的值（value）为空，且searchValue的长度不为0，从后向前依次删除生成的每一项搜索关键字
        if (e.keyCode === 8 && (!value || this.searchInput.selectionStart === 0) && searchValue.length !== 0) {
            this.setState({
                searchValue: searchValue.slice(0, searchValue.length - 1),
            }, () => {
                this.props.onComboChange(this.state.searchValue, this.state.value)
                this.searchFilterInput = this.searchFilterInput.slice(1)
            });
        }

        // 当按下enter键且value值不为空时且没有重复条件时
        if (e.keyCode === 13 && trim(value) && !searchValue.find((item) => item && item.key === keys[0] && item.value === trim(value))) {
            const newValue = keys && keys.length !== 0
                ?
                {
                    index: index++,
                    key: keys[0],
                    value: trim(value),
                }
                :
                {
                    index: index++,
                    value: trim(value),
                }
            this.setState({
                searchValue: [...searchValue, newValue],
                value: '',
            }, () => {
                this.searchInput.focus()
                this.props.onComboChange(this.state.searchValue, this.state.value)
            })
        }

        // 如果输入框为空 && 按键为左方向键
        if (e.keyCode === 37 && (!value || this.searchInput.selectionStart === 0) && this.searchFilterInput.length > 0) {
            this.searchFilterInput[0].focus();
        }
    }

    /**
     * 清空全部搜索关键词
     */
    protected handleTotalDelete() {
        this.setState({
            value: '',
            searchValue: [],
        }, () => {
            this.props.onComboChange(this.state.searchValue, this.state.value)
        })
    }

    /**
     * 选择下拉框中相应的搜索条件时触发
     * @param key
     * @param value
     */
    protected handleSearchItemClick(key, value) {
        if (this.state.searchValue.find((item) => item && item.key === key && item.value === value)) {
            return
        }
        const newValue = {
            index: index++,
            key: key,
            value: value,
        }
        this.setState({
            searchValue: [...this.state.searchValue, newValue],
            value: '',
        }, () => {
            this.searchInput.focus()
            this.props.onComboChange(this.state.searchValue, this.state.value)
        })
    }

    /**
     * 禁止输入任何字符，只触发聚焦光标
     */
    protected handleStopInput(e) {
        this.setState({
            filterInputValue: '',
        })
    }

    /**
     * 前置输入框按键监听
     */
    protected handleSearchFilterDelete(e, index) {
        let { searchValue } = this.state;
        if (e.keyCode === 8 && index + 1 !== this.searchFilterInput.length) {
            this.setState({
                searchValue: [...searchValue.slice(0, this.searchFilterInput.length - index - 2), ...searchValue.slice(this.searchFilterInput.length - index - 1, searchValue.length)],
            }, () => {
                this.props.onComboChange(this.state.searchValue, this.state.value)
                this.searchFilterInput = [...this.searchFilterInput.slice(0, index + 1), ...this.searchFilterInput.slice(index + 2, this.searchFilterInput.length)]
            });
        }
        // 如果输入框为空 && 按键为左方向键
        if (e.keyCode === 37 && index + 1 !== this.searchFilterInput.length) {
            this.searchFilterInput[index + 1].focus();
        }
        // 如果输入框为空 && 按键为右方向键
        if (e.keyCode === 39) {
            if (index === 0) {
                this.firstFilterInput.focus();
                this.searchInput.focus();
            } else {
                this.searchFilterInput[index - 1].focus();
            }

        }
    }
}