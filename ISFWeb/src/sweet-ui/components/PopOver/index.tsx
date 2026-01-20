import React from 'react';
import { isArray, isFunction } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import Portal from '../Portal';
import Locator, { AlignPosition } from './Locator';
import styles from './styles';
import AppConfigContext from '@/core/context/AppConfigContext';

interface PopOverProps {
    open?: boolean;

    popup: React.ReactNode | ((
        props: {
            open: boolean;
            triggerElement?: Element;
            close: () => void;
        },
    ) => React.ReactNode);

    freeze: boolean;

    /**
     * 定位自适应
     */
    autoFix?: boolean;

    anchor?: HTMLElement;

    close: () => void;

    onClickAway: (e: SweetUIEvent<MouseEvent>) => void;

    onMouseEnter: () => void;

    onMouseLeave: () => void;

    anchorOrigin: [AlignPosition, AlignPosition];

    alignOrigin: [AlignPosition, AlignPosition];

    onPopupAlign?: (event: SweetUIEvent<{ x: number; y: number }>) => void;

    popupZIndex?: number;

    /**
    * 事件监听的目标元素，默认window
    */
    target?: HTMLElement | string;

    role?: string;

    /**
     * 根容器
     */
    element?: () => HTMLElement | HTMLElement;
}

export default class PopOver extends React.Component<PopOverProps, any> {
    static contextType = AppConfigContext;
    static defaultProps = {
        freeze: true,
    };

    constructor(props: PopOverProps, ...args: any[]) {
        super(props);
    }

    /**
     * 弹出层元素
     */
    popupInstance: Element | null = null;

    getContainer = () => {
        const element = isFunction(this.props.element) ? this.props.element() : this.props.element || this.context && this.context.element 
        const popupContainer = document.createElement('div');
        (((element && isArray(element) ? element[0] : element) || window.document.querySelector('body')) as HTMLBodyElement).appendChild(popupContainer as HTMLDivElement);
        if (this.props.freeze && !this.props.element) {
            (popupContainer as HTMLDivElement).setAttribute('class', styles['layer']);
        }

        popupContainer.setAttribute('role', this.props.role)

        return popupContainer
    }

    savePopupRef = (node: Element) => {
        this.popupInstance = node;
    };

    handleClickAway = (e: MouseEvent) => {
        if (
            this.props.open &&
            this.popupInstance &&
            !this.popupInstance.contains(e.target) &&
            ((this.props.anchor && !this.props.anchor.contains(e.target)) || !this.props.anchor)
        ) {
            this.dispatchClickAway(e);
        }
    };

    dispatchClickAway = createEventDispatcher(this.props.onClickAway);

    /**
     * popup对齐后触发
     */
    handleLayoutChange = ({ x, y }: { x: number; y: number }) => {
        this.dispatchPopupAlignEvent({ x, y });
    };

    dispatchPopupAlignEvent = createEventDispatcher(this.props.onPopupAlign);

    render() {
        return this.props.open ? (
            <Portal getContainer={this.getContainer}>
                <Locator
                    target={this.props.target}
                    anchor={this.props.anchor}
                    anchorOrigin={this.props.anchorOrigin}
                    alignOrigin={this.props.alignOrigin}
                    onMouseDown={this.handleClickAway}
                    onLayoutChange={this.handleLayoutChange}
                    autoFix={this.props.autoFix}
                    popContainerZIndex={this.props.popupZIndex}
                    element={isFunction(this.props.element) ? this.props.element() : this.props.element || this.context && this.context.element}
                >
                    <View
                        onMouseEnter={this.props.onMouseEnter}
                        onMouseLeave={this.props.onMouseLeave}
                        onMounted={this.savePopupRef}
                    >
                        {typeof (this.props.popup) === 'function' ?
                            this.props.popup({
                                close: this.props.close,
                                triggerElement: this.props.anchor,
                                open: this.props.open,
                            })
                            : this.props.popup
                        }
                    </View>
                </Locator>
            </Portal>
        ) : null;
    }
}
