import React from 'react';
import classnames from 'classnames'
import { bindEvent } from '@/util/browser';
import styles from './styles';

type AlignPosition = number | 'top' | 'left' | 'right' | 'bottom' | 'center';

interface LocatorProps {
    /**
     * 样式
     */
    className: string;

    /**
     * 定位锚点元素
     */
    anchor?: HTMLElement;

    /**
     * 定位锚点相对坐标
     */
    anchorOrigin?: [number | string, number | string];

    /**
     * 对齐坐标
     */
    alignOrigin?: [AlignPosition, AlignPosition];

    /**
     * 是否支持拖拽
     */
    draggable?: boolean;
}

export default class Locator extends React.Component<LocatorProps, any> {
    // 包裹定位内容的容器
    container: HTMLElement | null = null;

    // 用来标记是否拖拽过，如果为true则不再自动定位
    dragged: boolean = false;

    // 鼠标拖拽初始位置
    mouseInitialCord?: {
        x: number;

        y: number;
    };

    // 对话框初始位置
    containerInitialCord?: {
        top: number;

        left: number;
    };

    componentDidMount() {
        this.locate();
    }

    componentDidUpdate() {
        this.locate();
    }

    /**
     * 自动定位内容
     */
    private locate() {
        const { anchor = document.documentElement, anchorOrigin, alignOrigin } = this.props;
        const [anchorOriginX, anchorOriginY] = anchorOrigin;
        const [alignOriginX, alignOriginY] = alignOrigin;
        const { width, height, top, left, bottom, right } = anchor.getBoundingClientRect();
        let offsetY: number, offsetX: number;

        /**
         * 计算Y坐标
         * @param anchorOriginY 锚点的Y原点
         * @param alignOriginY 目标的Y原点
         * @param param2
         */
        const calcContentOffsetY = (
            anchorOriginY: AlignPosition,
            alignOriginY: AlignPosition,
            { autoFit = false } = {},
        ) => {
            let offsetY: number;

            if (this.container) {
                switch (anchorOriginY) {
                    case 'top':
                        switch (alignOriginY) {
                            case 'top':
                                offsetY = top;

                                if (offsetY + this.container.clientHeight > window.innerHeight && autoFit) {
                                    return calcContentOffsetY('top', 'bottom', { autoFit: false });
                                } else {
                                    return offsetY;
                                }

                            case 'bottom':
                                offsetY = top - this.container.clientHeight;

                                if (offsetY < 0 && autoFit) {
                                    return calcContentOffsetY('bottom', 'top', { autoFit: false });
                                } else {
                                    return offsetY;
                                }

                            case 'center':
                                offsetY = top - this.container.clientHeight / 2;

                                if (offsetY < 0 && autoFit) {
                                    return calcContentOffsetY('top', 'top', { autoFit: false });
                                } else {
                                    return offsetY;
                                }
                        }
                        break

                    case 'bottom':
                        switch (alignOriginY) {
                            case 'top':
                                offsetY = top + anchor.clientHeight;

                                if (offsetY + this.container.clientHeight > window.innerHeight && autoFit) {
                                    return calcContentOffsetY('top', 'bottom', { autoFit: false });
                                } else {
                                    return offsetY;
                                }

                            case 'bottom':
                                offsetY = top + anchor.clientHeight - this.container.clientHeight;

                                if (offsetY < 0 && autoFit) {
                                    return calcContentOffsetY('bottom', 'top', { autoFit: false });
                                } else {
                                    return offsetY;
                                }

                            case 'center':
                                offsetY = top + anchor.clientHeight - this.container.clientHeight / 2;

                                if (offsetY + this.container.clientHeight > window.innerHeight && autoFit) {
                                    return calcContentOffsetY('bottom', 'bottom', { autoFit: false });
                                } else {
                                    return offsetY;
                                }
                        }
                        break

                    case 'center':
                        switch (alignOriginY) {
                            case 'top':
                                offsetY = (top + anchor.clientHeight) / 2;

                                if (offsetY + this.container.clientHeight > window.innerHeight && autoFit) {
                                    return calcContentOffsetY('center', 'bottom', { autoFit: false });
                                } else {
                                    return offsetY;
                                }

                            case 'bottom':
                                offsetY = (top + anchor.clientHeight) / 2 - this.container.clientHeight;

                                if (offsetY < 0 && autoFit) {
                                    return calcContentOffsetY('center', 'top', { autoFit: false });
                                } else {
                                    return offsetY;
                                }

                            case 'center':
                                offsetY = (top + anchor.clientHeight) / 2 - this.container.clientHeight / 2;

                                return offsetY;
                        }
                        break
                    default:
                        break
                }
            }
        };

        /**
         * 计算X坐标
         * @param anchorOriginX 锚点的X原点
         * @param alignOriginX 目标的X原点
         * @param param2
         */
        const calcContentOffsetX = (
            anchorOriginX: AlignPosition,
            alignOriginX: AlignPosition,
            { autoFit = false } = {},
        ) => {
            let offsetX: number;

            if (this.container) {
                switch (anchorOriginX) {
                    case 'left':
                        switch (alignOriginX) {
                            case 'left':
                                offsetX = left;

                                if (offsetX + this.container.clientWidth > window.innerWidth && autoFit) {
                                    return calcContentOffsetY('right', 'right', { autoFit: false });
                                } else {
                                    return offsetX;
                                }

                            case 'right':
                                offsetX = left - this.container.clientWidth;

                                if (offsetX < 0 && autoFit) {
                                    return calcContentOffsetY('left', 'left', { autoFit: false });
                                } else {
                                    return offsetX;
                                }

                            case 'center':
                                offsetX = left - this.container.clientWidth / 2;

                                if (offsetX < 0 && autoFit) {
                                    return calcContentOffsetY('left', 'left', { autoFit: false });
                                } else {
                                    return offsetX;
                                }
                        }
                        break

                    case 'right':
                        switch (alignOriginX) {
                            case 'left':
                                offsetX = left + anchor.clientWidth;

                                if (offsetX + this.container.clientWidth > window.innerWidth && autoFit) {
                                    return calcContentOffsetY('right', 'right', { autoFit: false });
                                } else {
                                    return offsetX;
                                }

                            case 'right':
                                offsetX = left + anchor.clientWidth - this.container.clientWidth;

                                if (offsetX < 0 && autoFit) {
                                    return calcContentOffsetY('left', 'left', { autoFit: false });
                                } else {
                                    return offsetX;
                                }

                            case 'center':
                                offsetX = left + anchor.clientWidth - this.container.clientWidth / 2;

                                if (offsetX + this.container.clientWidth > window.innerWidth && autoFit) {
                                    return calcContentOffsetY('right', 'right', { autoFit: false });
                                } else {
                                    return offsetX;
                                }
                        }
                        break

                    case 'center':
                        switch (alignOriginX) {
                            case 'left':
                                offsetX = (left + anchor.clientWidth) / 2;

                                if (offsetX + this.container.clientWidth > window.innerWidth && autoFit) {
                                    return calcContentOffsetY('center', 'right', { autoFit: false });
                                } else {
                                    return offsetX;
                                }

                            case 'right':
                                offsetX = (left + anchor.clientWidth) / 2 - this.container.clientWidth;

                                if (offsetX < 0 && autoFit) {
                                    return calcContentOffsetY('center', 'left', { autoFit: false });
                                } else {
                                    return offsetX;
                                }

                            case 'center':
                                offsetX = (left + anchor.clientWidth) / 2 - this.container.clientWidth / 2;

                                return offsetX;
                        }
                        break
                    default:
                        break
                }
            }
        };

        if (this.container) {
            this.container.style.top = `${Math.round(calcContentOffsetY(anchorOriginY, alignOriginY, { autoFit: true }))}px`;
            this.container.style.left = `${Math.round(calcContentOffsetX(anchorOriginX, alignOriginX, { autoFit: true }))}px`;
        }
    }

    /**
     * 随鼠标移动对话框
     * @param event 鼠标事件
     */
    private move(event: React.MouseEvent<HTMLElement>) {
        const container = this.container;

        if (container && this.mouseInitialCord && this.containerInitialCord) {
            container.style.top = `${event.clientY - this.mouseInitialCord.y + this.containerInitialCord.top}px`;
            container.style.left = `${event.clientX - this.mouseInitialCord.x + this.containerInitialCord.left}px`;

            this.dragged = true;
        }
    }

    /**
     * 鼠标按下开始移动
     * @param event 鼠标事件
     */
    protected startDrag(event: React.MouseEvent<HTMLElement>) {
        if (this.container) {
            const { top, left } = this.container.getBoundingClientRect();

            this.mouseInitialCord = {
                x: event.clientX,
                y: event.clientY,
            };

            this.containerInitialCord = {
                top,
                left,
            };

            bindEvent(document, 'mousemove', this.move);
        }
    }

    /**
     * TODO
     * 如果鼠标按下超过一定时间，则允许拖拽
     */
    private waitToDrag = (event: React.MouseEvent<HTMLDivElement>) => {
        const { draggable } = this.props;

        if (draggable) {
            if (this.waitToRegisterDragEvent) {
                clearTimeout(this.waitToRegisterDragEvent);
            }

            this.waitToRegisterDragEvent = setTimeout(() => {
                const handleMousemove = (event: MouseEvent) => {
                    // console.log(event);
                }

                window.addEventListener('mousemove', handleMousemove);

                window.addEventListener('mouseup', (event: MouseEvent) => {
                    window.removeEventListener('mousemove', handleMousemove);
                });
            }, 1000);
        }
    };

    render() {
        const { children } = this.props;

        return (
            <div ref={(node) => this.container = node} onMouseDown={this.waitToDrag} className={classnames(styles['wrapper'], this.props.className)}>
                {children}
            </div>
        );
    }
}
