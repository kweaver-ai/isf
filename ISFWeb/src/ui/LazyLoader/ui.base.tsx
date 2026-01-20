import React from 'react';
import { noop } from 'lodash';

export default class LazyLoaderBase extends React.Component<UI.LazyLoader.Props, any> {

    static defaultProps = {
        limit: 200,

        trigger: 0.75,

        onScroll: noop,

        onChange: noop,
    }

    state: UI.LazyLoader.State = {
        page: 1,
    }

    componentDidMount() {
        if (this.props.scroll) {
            this.scrollTo(this.props.scroll);
        }
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.props.scroll !== undefined && this.props.scroll !== this.scrollTop) {
            this.scrollTo(this.props.scroll);
        }
    }

    /**
     * 滚动位置
     */
    scrollTop: number;

    /**
     * 滚动容器的引用
     */
    scrollView: HTMLDivElement;

    /**
     * 实际内容的高度
     */
    viewHeight = 0;

    /**
     * 滚动到顶部
     */
    private scrollTo(scroll: number): void {
        if (scroll !== undefined) {
            this.scrollTop = this.scrollView.scrollTop = scroll;
        }
    }

    /**
     * 计算滚动位置并触发懒加载
     */
    protected handleScroll(event: Event): void {
        const { scrollTop, clientHeight, scrollHeight } = event.target;
        const triggerTop = scrollHeight - (this.props.trigger * scrollHeight);
        this.props.onScroll(scrollTop);

        if ((scrollHeight - clientHeight - scrollTop) < triggerTop && scrollTop !== 0) {
            if (scrollHeight !== this.viewHeight) {
                this.viewHeight = scrollHeight

                this.setState({ page: this.state.page + 1 }, () => {
                    this.props.onChange(this.state.page, this.props.limit)
                })
            }
        }
    }

    reset() {
        this.setState({
            page: 1,
            scroll: 0,
        })
        this.scrollTo(0)

        this.viewHeight = 0
    }
}