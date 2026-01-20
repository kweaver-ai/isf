import React from 'react';
import ReactDOM from 'react-dom';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import { AlignPosition } from '../PopOver/Locator';
import PopOver from '../PopOver';
import AppConfigContext from '@/core/context/AppConfigContext';

export type TriggerEvent = 'click' | 'hover' | 'contextmenu' | 'focus';
interface TriggerProps {
    /**
     * 触发行为
     */
    triggerEvent: TriggerEvent | ReadonlyArray<TriggerEvent>;

    /**
     * 定义如何渲染触发元素
     */
    renderer?: (
        props: Partial<{
            setPopupVisibleOnMouseEnter: () => void;
            setPopupVisibleOnMouseLeave: () => void;
            setPopupVisibleOnClick: () => void;
            setPopupVisibleOnFocus: () => void;
            setPopupVisibleOnBlur: () => void;
            role?: string;
        }>,
    ) => React.ReactNode;

    /**
     * 鼠标移出时的延迟，单位：s
     */
    mouseLeaveDelay?: number;

    /**
     * 弹出层展开时是否冻结滚动条
     */
    freeze?: boolean;

    /**
     * 定位自适应
     */
    autoFix?: boolean;

    anchor?: HTMLElement;

    /**
      * 触发元素定位原点
      */
    anchorOrigin?: [AlignPosition, AlignPosition];

    /**
     * 弹出元素定位原点
     */
    alignOrigin?: [AlignPosition, AlignPosition];

    open?: boolean;

    popupZIndex?: number;

    /**
     * 弹出层关闭之前触发
     */
    onBeforePopupClose?: (event: SweetUIEvent<any>) => void;

    /**
     * 弹出层展开状态变化时触发
     */
    onPopupVisibleChange?: (event: SweetUIEvent<boolean>) => void;

    /**
     * 弹出层对齐后触发
     */
    onPopupAlign?: (alignPoint: SweetUIEvent<{ x: number; y: number }>) => void;

    role?: string;

    element?: any;
}

interface TriggerState {
    /**
     * popup展开状态
     */
    open: boolean;

    /**
     * 弹出层窗口的父NW窗口(pc)
     */
    opener?: any;
}

export default class Trigger extends React.Component<TriggerProps, TriggerState> {
    static contextType = AppConfigContext
    static defaultProps = {
        open: false,
        triggerEvent: 'click',
        mouseLeaveDelay: 0.1,
        freeze: true,
        renderer: () => null,
        anchorOrigin: ['left', 'bottom'],
        alignOrigin: ['left', 'top'],
    };

    constructor(props: TriggerProps, ...args: any[]) {
        super(props);
        this.state = {
            open: !!props.open,
            opener: null,
        };
    }

    /**
     * 弹出层实例
     */
    popupComponent: Element | null = null;

    triggerElement: HTMLElement | null = document.body;

    delayTimer = null;

    componentDidMount() {
        this.triggerElement = ReactDOM.findDOMNode(this);
    }

    public getRootDomNode = () => {
        return ReactDOM.findDOMNode(this);
    };

    componentWillUnmount() {
        // 在组件销毁前设置state，防止内存泄漏
        this.setState = (state, callback) => {
            return;
        }
        this.clearDelayTimer();
    }

    /**
     * 关闭弹出层的方法
     */
    close = () => {
        this.setPopupVisible(false);
    };

    /**
     * 触发元素触发click事件
     */
    setPopupVisibleOnClick = () => {
        if (this.props.triggerEvent.indexOf('click') !== -1 || this.props.triggerEvent === 'click') {
            this.setPopupVisible(!this.state.open);
        }
    };

    setPopupVisibleOnMouseEnter = () => {
        if (this.props.triggerEvent.indexOf('hover') !== -1 || this.props.triggerEvent === 'hover') {
            this.delaySetPopupVisible(true);
        }
    };

    setPopupVisibleOnMouseLeave = () => {
        if (this.props.triggerEvent.indexOf('hover') !== -1 || this.props.triggerEvent === 'hover') {
            this.delaySetPopupVisible(false, this.props.mouseLeaveDelay);
        }
    };

    setPopupVisibleOnFocus = () => {
        if (this.props.triggerEvent.indexOf('focus') !== -1 || this.props.triggerEvent === 'focus') {
            this.setPopupVisible(true);
        }
    };

    setPopupVisibleOnBlur = () => {
        if (this.props.triggerEvent.indexOf('focus') !== -1 || this.props.triggerEvent === 'focus') {
            this.setPopupVisible(false);
        }
    };

    handlePopupMouseEnter = () => {
        if (this.props.triggerEvent.indexOf('hover') !== -1 || this.props.triggerEvent === 'hover') {
            this.clearDelayTimer();
        }
    };

    handlePopupMouseLeave = () => {
        if (this.props.triggerEvent.indexOf('hover') !== -1 || this.props.triggerEvent === 'hover') {
            this.delaySetPopupVisible(false, this.props.mouseLeaveDelay);
        }
    };

    delaySetPopupVisible(visible: boolean, delayS: number = 0) {
        const delay = delayS * 1000;

        this.clearDelayTimer();
        if (delay) {
            this.delayTimer = setTimeout(() => {
                this.setPopupVisible(visible);
            }, delay);
        } else {
            this.setPopupVisible(visible);
        }
    }

    clearDelayTimer() {
        if (this.delayTimer) {
            clearTimeout(this.delayTimer);
            this.delayTimer = null;
        }
    }

    /**
     * 设置弹出层是否可见
     */
    setPopupVisible = (open: boolean) => {
        this.clearDelayTimer();
        if (this.state.open !== open) {
            this.setState({ open });
            this.dispatchPopupVisibleChangeEvent(open);
        }
    };

    /**
     * 弹出层展开状态改变时触发
     */
    dispatchPopupVisibleChangeEvent = createEventDispatcher(this.props.onPopupVisibleChange);

    /**
     * TODO 应用场景待定
     */
    savePopup = (node: Element) => {
        this.popupComponent = node;
    };

    handleClickAway = () => {
        this.dispatchClickAway();
    };

    dispatchClickAway = createEventDispatcher(this.props.onBeforePopupClose, () => {
        this.close();
    });

    renderPopOver() {
        return (
            <PopOver
                role={this.props.role}
                element={this.props.element || this.context?.element}
                key={'triggerPopOver'}
                ref={this.savePopup}
                open={this.props.anchor ? this.props.open : this.state.open}
                anchor={this.props.anchor || this.triggerElement}
                onClickAway={this.handleClickAway}
                close={this.close}
                popup={this.props.children}
                freeze={this.props.freeze}
                autoFix={this.props.autoFix}
                onMouseEnter={this.handlePopupMouseEnter}
                onMouseLeave={this.handlePopupMouseLeave}
                anchorOrigin={this.props.anchorOrigin}
                alignOrigin={this.props.alignOrigin}
                onPopupAlign={this.props.onPopupAlign}
                popupZIndex={this.props.popupZIndex}
            />
        );
    }

    render() {
        return (
            [
                typeof this.props.renderer === 'function' ? (
                    this.props.renderer({
                        setPopupVisibleOnClick: this.setPopupVisibleOnClick,
                        setPopupVisibleOnMouseEnter: this.setPopupVisibleOnMouseEnter,
                        setPopupVisibleOnMouseLeave: this.setPopupVisibleOnMouseLeave,
                        setPopupVisibleOnFocus: this.setPopupVisibleOnFocus,
                        setPopupVisibleOnBlur: this.setPopupVisibleOnBlur,
                        role: this.props.role,
                    })
                ) : null,
                this.renderPopOver(),
            ]
        );
    }
}
