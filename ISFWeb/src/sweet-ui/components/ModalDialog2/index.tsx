import React from 'react';
import Dialog, { DialogIcon, DialogButton, titleTextAlign } from '../Dialog2';
import Modal from '../Modal';
import Portal from '../Portal';
import Locator from '../Locator';
import AppConfigContext from '@/core/context/AppConfigContext';

interface ModalDialogProps {
    /**
     * 把内容挂载到目标节点
     */
    target?: HTMLElement;

    /**
     * 对话框标题
     */
    title?: string;

    /**
     * 对话框标题的位置
     */
    titleTextAlign?: titleTextAlign;

    /**
     * 对话框宽度
     */
    width?: number | string;

    /**
    * 对话框顶部的图标按钮
    */
    icons?: ReadonlyArray<DialogIcon>;

    /**
     * 对话框底部按钮
     */
    buttons?: ReadonlyArray<DialogButton>;

    /**
     * 对话框是否支持拖拽
     */
    draggable?: boolean;

    /**
     * className
     */
    className?: string;

    /**
     * style
     */
    style?: React.CSSProperties;

    /**
     * 对话框zIndex值
     */
    zIndex: number;

    /**
     * role
     */
    role?: string;

    /**
     * 是否透视，目前仅支持右上角透视
     */
    transparentRegion?: {
        width: number;
        height: number;
    };

    /**
     * 限制的拖拽范围（目前只用到了top，后续还可补充left）
     */
    restrictedDragRange?: {
        top?: number;
        left?: number;
    };
}

interface ModalDialogState {
    /**
     * 把内容挂载到目标节点
     */
    target: HTMLElement;
}

export default class ModalDialog extends React.Component<ModalDialogProps, ModalDialogState> {
    static contextType = AppConfigContext
    static defaultProps = {
        draggable: true, // 弹窗默认支持拖拽
        zIndex: 20,
    };

    // 缓存容器元素
    private container: HTMLElement | null = null;

    constructor(props: ModalDialogProps) {
        super(props)
        this.state = {
            target: this.props.target,
        }
    }

    componentDidMount() {
        // 弹窗弹出时候，底部页面不允许滚镀
        document.body.style.overflow = 'hidden'
    }

    componentWillUnmount() {
        // 弹窗销毁时，取消overflow，hidden 设置
        document.body.style.overflow = '';
        // 移除缓存的容器元素
        if (this.container) {
            const parent = this.container.parentElement;
            if (parent) {
                parent.removeChild(this.container);
            }
            this.container = null;
        }
    }

    private getDefaultTarget = (transparentRegion?: { width: number; height: number }): HTMLElement => {
        if (this.container) {
            return this.container;
        }
        const { target } = this.props
        // 使用可选链操作符避免报错
        const rootElement = target || this.context?.element || document.body
        const container = document.createElement('div');
        container.style.position = 'fixed'
        container.style.top = '0'
        container.style.right = '0'
        container.style.bottom = '0'
        container.style.left = '0'
        container.style.zIndex = String(this.props.zIndex)
        if (transparentRegion) {
            container.style.pointerEvents = 'none'
        }
        container.setAttribute('role', this.props.role);

        rootElement.appendChild(container);
        this.container = container;

        return container
    }

    render() {
        const { children, title, titleTextAlign, width, icons, buttons, draggable, style, className, transparentRegion, restrictedDragRange } = this.props
        const { target } = this.state
        const rootElement = this.getDefaultTarget(this.props.transparentRegion)

        return (
            <Portal getContainer={() => rootElement }>
                <Modal transparentRegion={transparentRegion}>
                    <Locator
                        anchor={rootElement}
                        anchorOrigin={['center', 'center']}
                        alignOrigin={['center', 'center']}
                    >
                        <Dialog
                            title={title}
                            titleTextAlign={titleTextAlign}
                            width={width}
                            icons={icons}
                            buttons={buttons}
                            draggable={draggable}
                            restrictedDragRange={restrictedDragRange}
                            style={style}
                            className={className}
                        >
                            {children}
                        </Dialog>
                    </Locator>
                </Modal>
            </Portal>
        );
    }
}