import React from 'react';
import { noop } from 'lodash';

const itemHeight = 32

export default class ListBase extends React.PureComponent<UI.List.Props, UI.List.State> {

    static defaultProps = {
        /**
         * 数据列表
         */
        data: [],

        onMouseDown: noop,

        template: null,

        selectIndex: null,

        onSelectionChange: noop,

    }

    state = {
        selectIndex: -1,
    }

    refs: {
        list: HTMLInputElement;
    }

    selectByMouse: boolean = true;

    componentDidUpdate(prevProps, prevState) {
        if (this.props.selectIndex !== prevState.selectIndex) {
            this.setState({
                selectIndex: this.props.selectIndex,
            }, () => {
                if(this.state.selectIndex !== prevState.selectIndex) {
                    const { viewHeight } = this.props;

                    this.selectByMouse = false;
                    // 每次滑动，增加或者减少scrollTop值
                    const height = itemHeight * (this.state.selectIndex + 1);

                    if ((height - this.refs.list.scrollTop) > viewHeight) {
                        this.refs.list.scrollTop = height - viewHeight
                    }
                    if ((height - this.refs.list.scrollTop) < itemHeight) {
                        this.refs.list.scrollTop = height - itemHeight
                    }
                }
            })
        }
    }

    /**
     * 鼠标移动事件
     */
    protected MoveByMouse() {
        this.selectByMouse = true;
    }

    /**
     * 悬浮选中搜索结果
     */
    protected handleMouseOver(index) {
        if (this.selectByMouse) {
            this.setState({ selectIndex: index }, () => {
                this.props.onSelectionChange(index);
            });
        }
    }

    /**
     * 悬浮离开取消选中搜索结果
     */
    protected handleMouseLeave() {
        if (this.selectByMouse) {
            this.setState({ selectIndex: -1 }, () => {
                this.props.onSelectionChange(-1);
            })
        }
    }

}