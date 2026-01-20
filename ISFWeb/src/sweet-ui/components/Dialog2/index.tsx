import React from 'react';
import classnames from 'classnames';
import { bindEvent, unbindEvent } from '@/util/browser';
import View from '../View';
import DialogHeader from './DialogHeader';
import Button from '../Button';
import styles from './styles';

/**
 * 对话框顶部的icon
 */
export type DialogIcon = {
    /**
     * 图标
     */
    icon: React.ReactNode;

    /**
     * 图标的点击事件
     */
    onClick: (event: React.MouseEvent<HTMLButtonElement>) => any;
}

/**
 * 按钮主题
 */
type Theme = 'regular' | 'oem' | 'gray' | 'dark' | 'text' | 'noborder';

/**
 * 按钮尺寸
 */
type Size = 'normal' | 'auto';

/**
 * 对话框底部的button
 */
export type DialogButton = {
    /**
     * 按钮文字
     */
    text: string;

    /**
     * 按钮主题
     */
    theme: Theme;

    /**
     * 按钮尺寸
     */
    size: Size;

    /**
     * 是否禁用
     */
    disabled: boolean;

    /**
     * 按钮点击事件
     */
    onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
}

/**
 * 标题位置（默认的标题位置居左）
 */
export type titleTextAlign = 'center' | 'left';

interface DialogProps {
    /**
     * 对话框宽度
     */
    width?: number | string;

    /**
     * 对话框标题
     */
    title?: React.ReactNode;

    /**
     * 标题位置
     */
    titleTextAlign?: titleTextAlign;

    /**
     * 对话框顶部的图标按钮
     */
    icons?: ReadonlyArray<DialogIcon>;

    /**
     * 对话框底部的操作按钮
     */
    buttons?: ReadonlyArray<DialogButton>;

    /**
     * 是否支持拖拽
     */
    draggable?: boolean;

    /**
     * 限制拖拽的范围（目前只用到了top，后续还可补充left）
     */
    restrictedDragRange?: {
        top?: number;
        left?: number;
    };

    /**
     * className
     */
    className?: string;

    /**
     * style
     */
    style?: React.CSSProperties;
}

export default class Dialog extends React.Component<DialogProps, any> {
    static defaultProps = {
        titleTextAlign: 'left',
        icons: [],
        buttons: [],
    }

    constructor(props, context) {
        super(props, context);
        this.startDrag = this.startDrag.bind(this)
        this.endDrag = this.endDrag.bind(this)
        this.move = this.move.bind(this);
    }

    /**
     * 外层容器
     */
    container: null | HTMLElement;

    /**
     * 鼠标拖拽初始位置
     */
    mouseInitialCord: {
        x: number;

        y: number;
    }

    /**
     * 对话框初始位置
     */
    dialogInitialCord: {
        top: number;

        left: number;
    }

    /**
     * 开始拖拽
     */
    startDrag = (event: React.MouseEvent<HTMLElement>) => {
        if (!this.props.draggable) {
            return
        }

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
    };

    /**
    * 随鼠标移动对话框
    * @param event 鼠标移动事件对象
    */
    private move(event: React.MouseEvent<HTMLElement>) {
        event.stopPropagation();
        const el = this.container;

        const { restrictedDragRange: { top: restrictedTop } = {} } = this.props

        el.style.position = `fixed`

        const top = event.clientY - this.mouseInitialCord.y + this.dialogInitialCord.top
        el.style.top = `${restrictedTop ? Math.max(restrictedTop, top) : top}px`
        el.style.left = `${event.clientX - this.mouseInitialCord.x + this.dialogInitialCord.left}px`
    }

    /**
     * 拖拽事件结束
     */
    endDrag = () => {
        if (!this.props.draggable) {
            return
        }

        unbindEvent(document, 'mousemove', this.move);
    }

    render() {
        const { title, titleTextAlign, width, style, className, icons, buttons, draggable, children } = this.props
        return (
            <div
                className={classnames(styles['dialog'], className)}
                style={{ ...style, width }}
                ref={(container) => this.container = container}
            >
                <DialogHeader
                    title={title}
                    titleTextAlign={titleTextAlign}
                    icons={icons}
                    onMouseDown={this.startDrag}
                    onMouseUp={this.endDrag}
                />
                <View
                    className={classnames(
                        styles['content'],
                        { [styles['content-padding-top']]: title && titleTextAlign === 'left' },
                    )}
                >
                    {children}
                </View>
                {
                    buttons && buttons.length !== 0 && (
                        <View className={styles['footer']}>
                            {
                                buttons.map(({ text, theme, size, disabled, onClick, ...others }, index) => (
                                    <Button
                                        role={'sweetui-button'}
                                        key={index}
                                        theme={theme}
                                        size={size}
                                        disabled={disabled}
                                        onClick={onClick}
                                        className={styles['button']}
                                        {...others}
                                    >
                                        {text}
                                    </Button>
                                ))
                            }
                        </View>
                    )
                }
            </div>
        );
    }
}