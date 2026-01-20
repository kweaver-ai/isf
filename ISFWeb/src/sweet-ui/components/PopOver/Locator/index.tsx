import React from 'react';
import classnames from 'classnames';
import { throttle, isFunction, isArray } from 'lodash';
import EventListener, { withOptions } from '../../EventListener/index';
import styles from './styles';

export type AlignPosition = number | 'top' | 'left' | 'right' | 'bottom' | 'center';

interface LocatorProps {
    /**
     * 定位锚点元素
     */
    anchor?: HTMLElement;

    /**
     * 定位锚点相对坐标
     */
    anchorOrigin?: [AlignPosition, AlignPosition];

    /**
     * 对齐坐标
     */
    alignOrigin?: [AlignPosition, AlignPosition];

    /**
	 * 定位自适应
	 */
    autoFix?: boolean | 'vertical' | 'horizontal';

    /**
	 * window监听到鼠标按下时触发
	 */
    onMouseDown: (event: MouseEvent) => void;

    /**
	 * 位置改变后触发
	 */
    onLayoutChange?: (point: { x: number; y: number }) => void;

    popContainerZIndex?: number;

    /**
	* 事件监听的目标元素，默认window
	*/
    target?: HTMLElement;

    /**
     * 根容器，非数组时，默认相对于element定位
     * 数组时：element[0]：为弹出层的容器，element[1]：为相对定位容器
     */
    element?: HTMLElement | [HTMLElement, HTMLElement | 'window'];
}

// eslint-disable-next-line
interface LocatorState { }

export default class Locator extends React.Component<LocatorProps, LocatorState> {
    static defaultProps = {
        autoFix: true,
        target: 'window',
    };

    // 包裹定位内容的容器
    container: HTMLElement | null = null;

    componentDidMount() {
        this.locate();
    }

    componentDidUpdate() {
        this.locate();
    }

    handleResize = () => {
        this.handleLocateLazily();
    };

    handleScroll = () => {
        this.handleLocateLazily();
    };

    handleLocateLazily = throttle(this.locate, 1000 / 24, { trailing: true, leading: false });

    /**
     * 自动定位内容
     */
    private locate() {
        if (this.container) {
            const { anchor = document.documentElement, anchorOrigin, alignOrigin, element } = this.props;
            const [anchorOriginX, anchorOriginY] = anchorOrigin;
            const [alignOriginX, alignOriginY] = alignOrigin;
            const { offsetHeight, offsetWidth }: { offsetHeight: number; offsetWidth: number } = this.container;

            let top = 0, left = 0, bottom = 0, right = 0

            if (element && !isArray(element)) {
                const eleRect = element.getBoundingClientRect()
                const anchorRect = anchor.getBoundingClientRect()
                top = anchorRect.top - eleRect.top
                left = anchorRect.left - eleRect.left
                bottom = anchorRect.bottom - eleRect.bottom
                right = anchorRect.right - eleRect.top
            } else {
                ({ top, left, bottom, right } = anchor.getBoundingClientRect())
            }

            let alignOffsetX: number = 0,
                alignOffsetY: number = 0;

            let windowInnerWidth = window.innerWidth, windowInnerHeight = window.innerHeight

            if (element) {
                try {
                    windowInnerWidth = isArray(element) ? element[1] === 'window' ? window.innerWidth : element[1].clientWidth : element.clientWidth
                    windowInnerHeight = isArray(element) ? element[1] === 'window' ? window.innerHeight : element[1].clientHeight : element.clientHeight
                } catch{
                    windowInnerWidth = window.innerWidth, windowInnerHeight = window.innerHeight
                }
            }
            /**
         * 计算Y坐标
         * @param anchorOriginY 锚点的Y原点
         * @param alignOriginY 目标的Y原点
         * @param param2
         */
            const calcContentOffsetY: (
                anchorOriginY: AlignPosition,
                alignOriginY: AlignPosition,
                { autoFit }: { autoFit?: boolean | 'vertical' | 'horizontal' }
            ) => number = (anchorOriginY: AlignPosition, alignOriginY: AlignPosition, { autoFit = false } = {}) => {
                let offsetY: number;

                switch (anchorOriginY) {
                    case 'top':
                        switch (alignOriginY) {
                            case 'top':
                                offsetY = top;

                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('top', 'bottom', { autoFit: false });
                                }
                                alignOffsetY = 0;
                                return offsetY;

                            case 'bottom':
                                offsetY = top - offsetHeight;

                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('bottom', 'top', { autoFit: false });
                                }
                                alignOffsetY = offsetHeight;
                                return Math.max(offsetY, 0);

                            case 'center':
                                offsetY = top - offsetHeight / 2;

                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('top', 'top', { autoFit: false });
                                } else if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('bottom', 'bottom', { autoFit: false });
                                }

                                alignOffsetY = offsetHeight / 2;
                                return offsetY;

                            default:
                                offsetY = top - Number(alignOriginY);
                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('top', 'top', { autoFit: false });
                                } else if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('bottom', 'bottom', { autoFit: false });
                                }

                                alignOffsetY = Number(alignOriginY);
                                return offsetY;
                        }

                    case 'bottom':
                        switch (alignOriginY) {
                            case 'top':
                                offsetY = top + anchor.offsetHeight;

                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('top', 'bottom', { autoFit: false });
                                }
                                alignOffsetY = 0;
                                return offsetY;

                            case 'bottom':
                                offsetY = top + anchor.offsetHeight - offsetHeight;

                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('bottom', 'top', { autoFit: false });
                                }
                                alignOffsetY = offsetHeight;
                                return offsetY;

                            case 'center':
                                offsetY = top + anchor.offsetHeight - offsetHeight / 2;

                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('bottom', 'bottom', { autoFit: false });
                                } else if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('top', 'top', { autoFit: false });
                                }
                                alignOffsetY = offsetHeight / 2;
                                return offsetY;

                            default:
                                offsetY = bottom - Number(alignOriginY);
                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('top', 'bottom', { autoFit: false });
                                } else if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('bottom', 'top', { autoFit: false });
                                }
                                alignOffsetY = Number(alignOriginY);
                                return offsetY;
                        }

                    case 'center':
                        switch (alignOriginY) {
                            case 'top':
                                offsetY = top + anchor.offsetHeight / 2;

                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('center', 'bottom', { autoFit: false });
                                }
                                alignOffsetY = 0;
                                return offsetY;

                            case 'bottom':
                                offsetY = top + anchor.offsetHeight / 2 - offsetHeight;

                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('center', 'top', { autoFit: false });
                                }
                                alignOffsetY = offsetHeight;
                                return offsetY;

                            case 'center':
                                offsetY = top + anchor.offsetHeight / 2 - offsetHeight / 2;

                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('center', 'top', { autoFit: false });
                                } else if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('center', 'bottom', { autoFit: false });
                                }
                                alignOffsetY = offsetHeight / 2;
                                return offsetY;

                            default:
                                offsetY = bottom - anchor.offsetHeight / 2 - Number(alignOriginY);

                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY('center', 'top', { autoFit: false });
                                } else if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY('center', 'bottom', { autoFit: false });
                                }
                                alignOffsetY = Number(alignOriginY);
                                return offsetY;
                        }
                    default:
                        switch (alignOriginY) {
                            case 'top':
                                offsetY = top + Number(anchorOriginY);

                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY(anchorOriginY, 'bottom', { autoFit: false });
                                }
                                alignOffsetY = 0;
                                return offsetY;

                            case 'bottom':
                                offsetY = top + Number(anchorOriginY) - offsetHeight;

                                if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY(anchorOriginY, 'top', { autoFit: false });
                                }
                                alignOffsetY = offsetHeight;
                                return offsetY;

                            case 'center':
                                offsetY = top + Number(anchorOriginY) - offsetHeight / 2;

                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY(anchorOriginY, 'bottom', { autoFit: false });
                                } else if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY(anchorOriginY, 'top', { autoFit: false });
                                }

                                alignOffsetY = offsetHeight / 2;
                                return offsetY;

                            default:
                                offsetY = top + Number(anchorOriginY) - Number(alignOriginY);
                                if (offsetY + offsetHeight > windowInnerHeight && autoFit) {
                                    alignOffsetY = offsetHeight;
                                    return calcContentOffsetY(anchorOriginY, 'bottom', { autoFit: false });
                                } else if (offsetY < 0 && autoFit) {
                                    alignOffsetY = 0;
                                    return calcContentOffsetY(anchorOriginY, 'top', { autoFit: false });
                                }
                                alignOffsetY = Number(alignOriginY);
                                return offsetY;
                        }
                }
            };

            /**
         * 计算X坐标
         * @param anchorOriginX 锚点的X原点
         * @param alignOriginX 目标的X原点
         * @param param2
         */
            const calcContentOffsetX: (
                anchorOriginX: AlignPosition,
                alignOriginX: AlignPosition,
                { autoFit }: { autoFit?: boolean | 'vertical' | 'horizontal' }
            ) => number = (anchorOriginX: AlignPosition, alignOriginX: AlignPosition, { autoFit = false } = {}) => {
                let offsetX: number;

                switch (anchorOriginX) {
                    case 'left':
                        switch (alignOriginX) {
                            case 'left':
                                offsetX = left;

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX('right', 'right', { autoFit: false });
                                }
                                alignOffsetX = 0;
                                return offsetX;

                            case 'right':
                                offsetX = left - offsetWidth;

                                if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    if (right + offsetWidth > windowInnerWidth) {
                                        return calcContentOffsetX('right', 'left', { autoFit: false });
                                    }
                                    return calcContentOffsetX('left', 'left', { autoFit: false });
                                }
                                alignOffsetX = offsetWidth;
                                return offsetX;

                            case 'center':
                                offsetX = left - offsetWidth / 2;

                                if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    return calcContentOffsetX('left', 'left', { autoFit: false });
                                }
                                alignOffsetX = offsetWidth / 2;
                                return offsetX;

                            default:
                                offsetX = left - Number(alignOriginX);

                                if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    return calcContentOffsetX('left', 'left', { autoFit: false });
                                } else if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX('right', 'right', { autoFit: false });
                                }
                                alignOffsetX = Number(alignOriginX);
                                return offsetX;
                        }

                    case 'right':
                        switch (alignOriginX) {
                            case 'left':
                                offsetX = left + anchor.offsetWidth;

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX('right', 'right', { autoFit: false });
                                }
                                alignOffsetX = 0;
                                return offsetX;

                            case 'right':
                                offsetX = left + anchor.offsetWidth - offsetWidth;

                                if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    return calcContentOffsetX('left', 'left', { autoFit: false });
                                }
                                alignOffsetX = offsetWidth;
                                return Math.max(offsetX, 0);

                            case 'center':
                                offsetX = left + anchor.offsetWidth - offsetWidth / 2;

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    if (right - offsetWidth < 0) {
                                        alignOffsetX = 0;
                                        return calcContentOffsetX('left', 'left', { autoFit: false });
                                    }
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX('right', 'right', { autoFit: false });
                                }
                                alignOffsetX = offsetWidth / 2;
                                return offsetX;

                            default:
                                offsetX = right - Number(alignOriginX);

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX('right', 'right', { autoFit: false });
                                } else if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    return calcContentOffsetX('right', 'left', { autoFit: false });
                                }
                                alignOffsetX = Number(alignOriginX);
                                return offsetX;
                        }

                    case 'center':
                        switch (alignOriginX) {
                            case 'left':
                                offsetX = left + anchor.offsetWidth / 2;

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    if (offsetX < offsetWidth) {
                                        alignOffsetX = 0;
                                        return calcContentOffsetX('left', 'left', { autoFit: false });
                                    }
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX('center', 'right', { autoFit: false });
                                }
                                alignOffsetX = 0;
                                return offsetX;

                            case 'right':
                                offsetX = left + anchor.offsetWidth / 2 - offsetWidth;

                                if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    if (left + anchor.offsetWidth / 2 + offsetWidth > windowInnerWidth) {
                                        return calcContentOffsetX('left', 'left', { autoFit: false });
                                    }
                                    return calcContentOffsetX('center', 'left', { autoFit: false });
                                }
                                alignOffsetX = offsetWidth;
                                return offsetX;

                            case 'center':
                                offsetX = left + anchor.offsetWidth / 2 - offsetWidth / 2;

                                if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    if (left + anchor.offsetWidth / 2 + offsetWidth > windowInnerWidth) {
                                        return calcContentOffsetX('left', 'left', { autoFit: false });
                                    }
                                    return calcContentOffsetX('center', 'left', { autoFit: false });
                                } else if (offsetX + offsetWidth > windowInnerWidth) {
                                    //
                                    alignOffsetX = offsetWidth;
                                    if (left + anchor.offsetWidth / 2 < offsetWidth) {
                                        return calcContentOffsetX('right', 'right', { autoFit: false });
                                    }
                                    return calcContentOffsetX('center', 'right', { autoFit: false });
                                }

                                alignOffsetX = offsetWidth / 2;
                                return offsetX;

                            default:
                                offsetX = right - anchor.offsetWidth / 2 - Number(alignOriginX);
                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX('center', 'right', { autoFit: false });
                                } else if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    if (left + anchor.offsetWidth / 2 + offsetWidth > windowInnerWidth) {
                                        return calcContentOffsetX('left', 'left', { autoFit: false });
                                    }
                                    return calcContentOffsetX('center', 'left', { autoFit: false });
                                }
                                alignOffsetX = Number(alignOriginX);
                                return offsetX;
                        }

                    default:
                        switch (alignOriginX) {
                            case 'left':
                                offsetX = left + Number(anchorOriginX);

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX(anchorOriginX, 'right', { autoFit: false });
                                }
                                alignOffsetX = 0;
                                return offsetX;

                            case 'right':
                                offsetX = left + Number(anchorOriginX) - offsetWidth;

                                if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    return calcContentOffsetX(anchorOriginX, 'left', { autoFit: false });
                                }
                                alignOffsetX = offsetWidth;
                                return offsetX;

                            case 'center':
                                offsetX = left + Number(anchorOriginX) - offsetWidth / 2;

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX(anchorOriginX, 'right', { autoFit: false });
                                } else if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    return calcContentOffsetX(anchorOriginX, 'left', { autoFit: false });
                                }
                                alignOffsetX = offsetWidth / 2;
                                return offsetX;

                            default:
                                offsetX = left + Number(anchorOriginX) - Number(alignOriginX);

                                if (offsetX + offsetWidth > windowInnerWidth && autoFit) {
                                    alignOffsetX = offsetWidth;
                                    return calcContentOffsetX(anchorOriginX, 'right', { autoFit: false });
                                } else if (offsetX < 0 && autoFit) {
                                    alignOffsetX = 0;
                                    return calcContentOffsetX(anchorOriginX, 'left', { autoFit: false });
                                }
                                alignOffsetX = Number(alignOriginX);
                                return offsetX;
                        }
                }
            };

            if (this.container) {
                let top = calcContentOffsetY(anchorOriginY, alignOriginY, { autoFit: this.props.autoFix });
                let left = calcContentOffsetX(anchorOriginX, alignOriginX, { autoFit: this.props.autoFix });

                if (this.props.autoFix) {
                    if (
                        (left + offsetWidth > windowInnerWidth && alignOffsetX < offsetWidth / 2) ||
                        (left < 0 && alignOffsetX > offsetWidth / 2)
                    ) {
                        left = right - Number(anchorOriginX) - offsetWidth + alignOffsetX;
                        if (left < 0) {
                            left = 0;
                        }
                    }
                    if (
                        (top + offsetHeight > windowInnerHeight && alignOffsetY < offsetHeight / 2) ||
                        (top < 0 && alignOffsetY > offsetHeight / 2)
                    ) {
                        top = bottom - Number(anchorOriginY) - offsetHeight + alignOffsetY;
                        if (top < 0) {
                            top = 0;
                        }
                    }
                }

                this.container.style.top = `${Math.round(top)}px`;
                this.container.style.left = `${Math.round(left)}px`;
                if (isFunction(this.props.onLayoutChange)) {
                    this.props.onLayoutChange({ x: alignOffsetX, y: alignOffsetY });
                }
            }
        }
    }

    handleMouseDown = (e: MouseEvent) => {
        this.props.onMouseDown(e);
    };

    render() {
        const { element } = this.props

        return (
            <EventListener
                target={this.props.target}
                onResize={this.handleResize}
                onMouseDown={this.handleMouseDown}
                onScroll={withOptions(this.handleScroll, { passive: true, capture: false })}
            >
                <div
                    ref={(node) => (this.container = node)}
                    className={classnames(styles['wrapper'], { [styles['pop-absolute']]: (element && !isArray(element)) || (isArray(element) && element[1] !== 'window') })}
                    style={{ zIndex: this.props.popContainerZIndex }}
                >
                    {this.props.children}
                </div>
            </EventListener>
        );
    }
}
