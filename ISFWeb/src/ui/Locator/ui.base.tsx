import React from 'react';
import { timer } from '@/util/timer'

interface Props {
    className?: string;
}

export default class LocatorBase extends React.Component<Props, any> {
    el: HTMLElement;

    container: HTMLElement

    parent: Node;

    timer: number;

    componentDidMount() {
        this.watch();
    }

    watch() {
        this.stopWatch = timer(() => {
            if (this.el) {

                // 首先得到el绝对定位时的位置
                this.el.style.position = 'absolute';
                this.el.style.top = '';
                this.el.style.left = '';

                const { top, left } = this.el.getBoundingClientRect();
                // 加1防止遮住父元素边框或阴影
                const nextTop = `${top + 1}px`
                const nextLeft = `${left}px`

                // 获得相对文档的位置
                this.el.style.position = 'fixed';
                this.el.style.top = nextTop;
                this.el.style.left = nextLeft;
            }
        }, 1000 / 24)
    }

    componentWillUnmount() {
        if (this.stopWatch) {
            this.stopWatch()
        }
    }
}