import React from 'react';
import { throttle, isFunction, noop } from 'lodash';
import { bindEvent, unbindEvent } from '@/util/browser';
import { getCenterCoordinate } from '../helper';

export default class DialogBase extends React.Component<UI.Dialog2.Props, UI.Dialog2.State> {

    static defaultProps = {
        draggable: true, // 是否支持拖拽

        onResize: noop, // 对话框尺寸变化时触发

        buttons: ['close'],

        width: 'auto',

        headless: false,
    }

    constructor(props, context) {
        super(props, context);
        this.move = throttle(this.move.bind(this), 20);
    }

    state = {
    }

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

    /**
     * 是否拖拽移动过对话框
     */
    moved: boolean;

    /**
     * 最外层容器
     */
    container: HTMLElement;

    /**
     * 宽度
     */
    width: number;

    /**
     * 高度
     */
    height: number;

    componentDidMount() {
        const { clientWidth, clientHeight } = this.container;

        this.updateSize({
            width: clientWidth,
            height: clientHeight,
        })
    }

    componentDidUpdate() {
        const { clientWidth, clientHeight } = this.container;

        if (clientWidth !== this.width || clientHeight !== this.height) {
            this.updateSize({
                width: clientWidth,
                height: clientHeight,
            })
        }
    }

    center() {
        const el = this.container;
        const { top, left } = getCenterCoordinate(this.container);

        el.style.top = `${Math.max(top, 0)}px`;
        el.style.left = `${left}px`;
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

    /**
     * 鼠标按下开始移动
     * @param event 鼠标事件
     */
    protected startDrag(event: MouseEvent) {
        const el = this.container;
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

    /**
     * 随鼠标移动对话框
     * @param event 鼠标移动事件对象
     */
    private move(event: MouseEvent) {
        const el = this.container;

        el.style.top = `${event.clientY - this.mouseInitialCord.y + this.dialogInitialCord.top}px`;
        el.style.left = `${event.clientX - this.mouseInitialCord.x + this.dialogInitialCord.left}px`;

        this.moved = true;
    }

    /**
     * 放开鼠标注销移动对话框
     * @param event 放开鼠标事件
     */
    protected endDrag(_event: MouseEvent) {
        unbindEvent(document, 'mousemove', this.move);
    }

    /**
     * 触发尺寸变化事件
     * @param size.width
     * @param size.height
     */
    private fireResizeEvent({ width, height }) {
        isFunction(this.props.onResize) && this.props.onResize({ width, height });
    }

    /**
     * 触发关闭对话框事件
     */
    protected fireCloseEvent() {
        isFunction(this.props.onClose) && this.props.onClose();
    }

}