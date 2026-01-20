import React from 'react'
import { noop } from 'lodash'
import { KeyDown } from '../AutoComplete/ui.base'

export default class AutoCompleteListBase extends React.Component<UI.AutoCompleteList.Props, UI.AutoCompleteList.State> {
    static defaultProps = {
        maxHeight: 200,

        selectIndex: -1,

        keyDown: KeyDown.NONE,

        onSelectionChange: noop,
    }

    list: HTMLElement;

    state = {
        selectIndex: this.props.selectIndex,
    }

    selectedByMouseOver: boolean = false;

    /**
     * 列表项dom集合
     */
    items = [];

    /**
     * 列表项高度值
     */
    itemHeight;

    componentDidMount() {
        this.updateListHeight()
    }

    componentDidUpdate(prevProps, prevState) {
        // children数目发生变化 或 Item高度没有全部更新时再去执行更新高度操作
        this.selectedByMouseOver = false;
        if (this.items.some((item) => item.style.height !== `${this.itemHeight}px`) ||
            React.Children.count(this.props.children) !== React.Children.count(prevProps.children)
        ) {
            this.updateListHeight()
        }

        const { selectIndex, keyDown } = this.props

        if(selectIndex !== -1 && selectIndex !== prevProps.selectIndex) {
            const selectIndex = this.state.selectIndex
            if (selectIndex === 0) {
                // 如果选中第一条，scrollTop为0
                this.list.scrollTop = 0
            } else if (selectIndex === React.Children.count(this.props.children) - 1) {
                // 选中最后一条，scrollTop为超出的高度
                this.list.scrollTop = this.itemHeight * (selectIndex + 1) - this.props.maxHeight
            } else {
                // 每次滑动，增加或者减少scrollTop
                const height = this.itemHeight * (selectIndex + 1);

                if ((height - this.list.scrollTop) > this.props.maxHeight) {
                    this.list.scrollTop = height - this.props.maxHeight
                }
                if ((height - this.list.scrollTop) < this.itemHeight) {
                    this.list.scrollTop = height - this.itemHeight
                }
            }
        }

        if (keyDown !== prevProps.keyDown) {
            switch (keyDown) {
                case KeyDown.DOWNARROW: {
                    // 按下向下键
                    this.handleDownArrow()
                    break;
                }
                case KeyDown.UPARROW: {
                    // 按下向上键
                    this.handleUpArrow()
                    break;
                }
            }
        }
    }

    /**
     * 更新列表显示高度
     */
    updateListHeight() {
        this.itemHeight = this.items && this.items.length && Math.max(...this.items.map((item) => item.clientHeight))
        const contentHeight = React.Children.count(this.props.children) * this.itemHeight

        this.list.style.height = (this.props.maxHeight > contentHeight) ?
            `${(contentHeight + 3)}px`
            : `${this.props.maxHeight}px`
        this.items.forEach((item) => {
            item.style.height = `${this.itemHeight}px`
        })
        this.items = []
    }

    /**
     * 获取子元素
     */
    getItem(item, index) {
        if (item) {
            this.items[index] = item
        }
    }

    /**
     * 处理鼠标移动到上面事件
     */
    handleMouseOver(e, selectIndex: number) {
        if (this.selectedByMouseOver) {
            this.setState({
                selectIndex,
            }, () => this.props.onSelectionChange(selectIndex))
        }
    }

    setSelectByMouseMove() {
        this.selectedByMouseOver = true;
    }

    /**
     * 按向下键触发
     */
    handleDownArrow() {
        const count = React.Children.count(this.props.children)

        if (count) {
            if (this.props.selectIndex === -1) {
                // 没有选中项选择第一个
                this.setState({
                    selectIndex: 0,
                }, () => this.props.onSelectionChange(0))
            } else if ((this.props.selectIndex + 1) < count) {
                // 有选中项且不是最后一个，选择下一个
                this.setState({
                    selectIndex: this.props.selectIndex + 1,
                }, () => this.props.onSelectionChange(this.state.selectIndex))
            } else {
                // 有选择项且为最后一个，选择第一
                this.setState({
                    selectIndex: 0,
                }, () => this.props.onSelectionChange(0))
            }
        }
    }

    /**
     * 按向上键触发
     */
    handleUpArrow() {
        const count = React.Children.count(this.props.children)

        if (count) {
            if (this.props.selectIndex === -1) {
                // 没有任何选中项选择最后一个
                this.setState({
                    selectIndex: count - 1,
                }, () => this.props.onSelectionChange(count - 1))
            } else if (this.props.selectIndex > 0) {
                // 有选中项且不是第一项，选择其上一个
                this.setState({
                    selectIndex: this.props.selectIndex - 1,
                }, () => this.props.onSelectionChange(this.state.selectIndex))
            } else {
                // 选择项是第一个的时候
                this.setState({
                    selectIndex: count - 1,
                }, () => this.props.onSelectionChange(count - 1))
            }
        }
    }
}