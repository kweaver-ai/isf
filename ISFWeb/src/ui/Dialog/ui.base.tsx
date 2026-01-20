import React from 'react';
import { throttle, noop } from 'lodash';
import { bindEvent, unbindEvent } from '@/util/browser';
import { getCenterCoordinate } from '../helper';

export default class DialogBase extends React.Component<UI.Dialog.Props, any> implements UI.Dialog.Element {
    // 鼠标拖拽初始位置
    mouseInitialCord: {
        x: number;

        y: number;
    }

    // 对话框初始位置
    dialogInitialCord: {
        top: number;

        left: number;
    }

    constructor(props, context) {
        super(props, context);
        this.move = throttle(this.move.bind(this), 20);
    }

    state: UI.Dialog.State = {
    }

    static defaultProps = {
        draggable: true, // 是否支持拖拽

        onResize: noop, // 对话框尺寸变化时触发
    }

    componentDidMount() {
        const { clientWidth, clientHeight } = this.refs.container;
        this.center();
        this.updateSize({
            width: clientWidth,
            height: clientHeight,
        })
    }

    componentDidUpdate() {
        const { clientWidth, clientHeight } = this.refs.container;

        if (clientWidth !== this.width || clientHeight !== this.height) {
            this.updateSize({
                width: clientWidth,
                height: clientHeight,
            })
        }
    }

    center() {
        this.setState(getCenterCoordinate(this.refs.container))
    }

    /**
     * 更新对话框尺寸并触发onResize事件
     * @param param0 width 宽度，height 高度
     */
    updateSize({ width, height }) {
        this.width = width;
        this.height = height;
        this.fireResizeEvent({ width, height })
    }

    protected startDrag(event) {
        if (event.target.getAttribute('role') === 'drag-area') {
            const el = this.refs.container;
            const { top, left } = el.getBoundingClientRect();

            this.mouseInitialCord = {
                x: event.clientX,
                y: event.clientY,
            }

            this.dialogInitialCord = {
                top,
                left,
            }

            bindEvent(document, 'mousemove', this.move);
        }
    }

    /**
     * 随鼠标移动对话框
     * @param event 鼠标移动事件对象
     */
    private move(event) {
        const el = this.refs.container;
        const { top, left } = el.getBoundingClientRect();

        el.style.top = `${event.clientY - this.mouseInitialCord.y + this.dialogInitialCord.top}px`;
        el.style.left = `${event.clientX - this.mouseInitialCord.x + this.dialogInitialCord.left}px`;
    }

    /**
     * 放开鼠标注销移动对话框
     * @param event 放开鼠标事件
     */
    protected endDrag(event) {
        if (event.target.getAttribute('role') === 'drag-area') {
            unbindEvent(document, 'mousemove', this.move);
        }
    }

    private fireResizeEvent({ width, height }) {
        this.props.onResize({ width, height });
    }

}