import React from 'react';
import classnames from 'classnames';
import { map, noop, isFunction } from 'lodash';
import { render } from 'react-dom';
import { bindEvent } from '@/util/browser'
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import SweetIcon from '../../components/SweetIcon';
import ModalDialog from '../../components/ModalDialog2';
import View from '../../components/View';
import __ from './locale';
import styles from './styles';
import { Portal } from '../../components';

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 */
type Alert = (params: { element?: HTMLElement; title?: string | HTMLElement; message?: string; detail?: string | HTMLElement; showCancelIcon?: boolean; zIndex?: number; transparentRegion?: { width: number; height: number }; restrictedDragRange?: { top?: number; left?: number } }) => Promise<boolean>

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 */
type Info = (params: { element?: HTMLElement; title?: string | HTMLElement; message?: string; detail?: string | HTMLElement; showCancelIcon?: boolean; zIndex?: number; transparentRegion?: { width: number; height: number }; restrictedDragRange?: { top?: number; left?: number } }) => Promise<boolean>

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 */
type Error = (params: { element?: HTMLElement; title?: string | HTMLElement; message?: string; detail?: string | HTMLElement; showCancelIcon?: boolean; zIndex?: number; transparentRegion?: { width: number; height: number }; restrictedDragRange?: { top?: number; left?: number } }) => Promise<boolean>

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 */
type Success = (params: { element?: HTMLElement; title?: string | HTMLElement; message?: string; zIndex?: number; transparentRegion?: { width: number; height: number }; restrictedDragRange?: { top?: number; left?: number } }) => Promise<boolean>

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 */
type Custom = (params: {
    title?: string | HTMLElement;
    message?: string;
    detail?: string | HTMLElement;
    showCancelIcon?: boolean;
    showCloseIcon?: boolean;
    zIndex?: number;
    element?: HTMLElement;
    iconName?: string;
    iconColor?: string;
    confirmText?: string;
    cancelText?: string;
    confirmProps?: Record<string, any>;
    transparentRegion?: { width: number; height: number };
    restrictedDragRange?: { top?: number; left?: number };
}) => Promise<boolean>

/**
 * 点击确定/取消后返回的结果
 */
type ActionResult = (resolve: boolean) => Promise<boolean>

/**
 * 按钮
 */
type MessageBoxButton = { text: string; theme: Theme; onClick: (event: SweetUIEvent<ActionResult>) => void }

/**
 * 图标按钮
 */
type MessageBoxIcon = { icon: React.ReactNode; onClick: (event: SweetUIEvent<ActionResult>) => void }

/**
 * 消息对话框
 * @param params.type 图标
 * @param params.title 对话框标题
 * @param params.message 消息
 * @param params.detail 详情
 * @param params.icons 顶部的图标按钮
 * @param params.buttons 底部的按钮
 * @param params.zIndex 对话框zIndex值
 */
type MessageBox = (params: { type: React.ReactNode; title?: string | HTMLElement; message?: string; detail?: string | HTMLElement; icons?: ReadonlyArray<MessageBoxIcon>; buttons?: ReadonlyArray<MessageBoxButton>; zIndex?: number; role?: string; element?: HTMLElement; transparentRegion?: { width: number; height: number }; restrictedDragRange?: { top?: number; left?: number } }) => Promise<boolean> & { destroy: Function }

/**
 * 按钮主题
 */
type Theme = 'regular' | 'oem' | 'gray' | 'dark' | 'text' | 'noborder';

const keyDownEvent = (event: KeyboardEvent) => {
    event.stopPropagation()
    event.preventDefault()
}

/**
 * 消息对话框
 * @param params
 */
export const messageBox: MessageBox = ({ type, title, message, detail, icons, buttons, zIndex = 20, role, element, transparentRegion, restrictedDragRange }) => {
    let container: HTMLDivElement | null = null;
    let unbindEvent: () => void = noop;
    let rootElement: HTMLElement | null = null; 

    const box = new Promise((resolve, reject) => {
        const MessageBoxComponent = () => {
            const contextElement = document.querySelector('#isf-web-plugins')
            rootElement = element || (contextElement as HTMLElement) || document.body;

            container = document.createElement('div');
            container.style.position = 'fixed';
            container.style.top = '0';
            container.style.right = '0';
            container.style.bottom = '0';
            container.style.left = '0';
            container.style.zIndex = String(zIndex);
            if (transparentRegion) {
                container.style.pointerEvents = 'none';
            }
            container.setAttribute('role', 'sweetui-message2');

            unbindEvent = bindEvent(window, 'keydown', keyDownEvent);

            const destroy = () => {
                // When destroying, do not execute the componentWillUnmount event of ModalDialog, manually remove the overflow of the body
                if (container && rootElement) {
                    rootElement.style.overflow = '';
                    rootElement.removeChild(container);
                    container = null;
                }
                isFunction(unbindEvent) && unbindEvent();
            };

            rootElement.appendChild(container);

            return (
                <Portal getContainer={() => container}>
                    <ModalDialog
                        transparentRegion={transparentRegion}
                        restrictedDragRange={restrictedDragRange}
                        target={container}
                        width={400}
                        icons={map(icons, (icon) => {
                            return {
                                ...icon,
                                onClick: () => {
                                    createEventDispatcher(icon.onClick, destroy)(resolve);
                                },
                            };
                        })}
                        buttons={map(buttons, (button) => {
                            return {
                                ...button,
                                onClick: () => {
                                    createEventDispatcher(button.onClick, destroy)(resolve);
                                },
                            };
                        })}
                    >
                        <View className={styles['content']}>
                            <View>
                                <View className={styles['icon']}>{type}</View>
                                <View className={styles['article']}>
                                    {title ? (
                                        <View className={styles['title']}>{title}</View>
                                    ) : null}
                                    {message ? (
                                        <View className={classnames(styles['message'], { [styles['small-padding']]: title })}>{message}</View>
                                    ) : null}
                                    {detail ? (
                                        <View className={styles['message-with-scroll']}>{detail}</View>
                                    ) : null}
                                </View>
                            </View>
                        </View>
                    </ModalDialog>
                </Portal>
            );
        };

        render(<MessageBoxComponent />, document.createElement('div'));
    }) as Promise<boolean> & { destroy: Function };

    box.destroy = () => {
        if (container) {
            (rootElement || document.body).style.overflow = '';
            (rootElement || document.body).removeChild(container);
            container = null;
        }
        isFunction(unbindEvent) && unbindEvent();
    };

    return box;
};

/**
 * 警告对话框
 * @returns Promise<boolean>
 */
export const alert: Alert = ({ title, message, detail, showCancelIcon = false, zIndex, element, transparentRegion, restrictedDragRange }) => {
    return messageBox({
        role: 'sweetui-message2.alert',
        transparentRegion,
        restrictedDragRange,
        element,
        title,
        message,
        detail,
        type: <SweetIcon name={'notice'} size={32} color={'#FAAD14'} />,
        icons: showCancelIcon ? [
            {
                icon: <SweetIcon name={'x'} size={16} />,
                onClick: ({ detail: resolve }) => resolve(false),
            },
        ] : [],
        buttons: showCancelIcon ?
            [
                {
                    text: __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                },
                {
                    text: __('取消'),
                    theme: 'regular',
                    onClick: ({ detail: resolve }) => resolve(false),
                },
            ] :
            [
                {
                    text: __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                },
            ],
        zIndex,
    })
}

/**
 * 一般消息对话框
 * @returns Promise<boolean>
 */
export const info: Info = ({ title, message, detail, showCancelIcon, zIndex, element, transparentRegion, restrictedDragRange }) => {
    return messageBox({
        role: 'sweetui-message2.info',
        transparentRegion,
        restrictedDragRange,
        element,
        title,
        message,
        detail,
        type: <SweetIcon name={'info'} size={32} color={'#1890FF'} />,
        icons: showCancelIcon ? [
            {
                icon: <SweetIcon name={'x'} size={16} />,
                onClick: ({ detail: resolve }) => resolve(false),
            },
        ] : [],
        buttons: showCancelIcon ?
            [
                {
                    text: __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                },
                {
                    text: __('取消'),
                    theme: 'regular',
                    onClick: ({ detail: resolve }) => resolve(false),
                },
            ] :
            [
                {
                    text: __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                },
            ],
        zIndex,
    })
}

/**
 * 错误提示框
 * @returns Promise<boolean>
 */
export const error: Error = ({ title, message, detail, showCancelIcon, zIndex, element, transparentRegion, restrictedDragRange }) => {
    return messageBox({
        role: 'sweetui-message2.error',
        transparentRegion,
        restrictedDragRange,
        element,
        title,
        message,
        detail,
        type: <SweetIcon name={'alert'} size={32} color={'#FF4D4F'} />,
        icons: showCancelIcon ? [
            {
                icon: <SweetIcon name={'x'} size={16} />,
                onClick: ({ detail: resolve }) => resolve(false),
            },
        ] : [],
        buttons: showCancelIcon ?
            [
                {
                    text: __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                },
                {
                    text: __('取消'),
                    theme: 'regular',
                    onClick: ({ detail: resolve }) => resolve(false),
                },
            ] :
            [
                {
                    text: __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                },
            ],
        zIndex,
    })
}

/**
 * 成功提示框
 * @returns Promise<boolean>
 */
export const success: Success = ({ title, message, zIndex, element, transparentRegion, restrictedDragRange }) => {
    return messageBox({
        role: 'sweetui-message2.success',
        transparentRegion,
        restrictedDragRange,
        element,
        title: <span style={{ color: '#505050', fontWeight: 600 }}>{title}</span>,
        message,
        type: <SweetIcon name={'success'} size={32} color={'#52C41B'} />,
        buttons: [
            {
                text: __('确定'),
                theme: 'oem',
                onClick: ({ detail: resolve }) => resolve(true),
            },
        ],
        zIndex,
    })
}

/**
 * 自定义对话框
 * @returns Promise<boolean>
 */
export const custom: Custom = ({
    title,
    transparentRegion,
    restrictedDragRange,
    message,
    detail,
    showCancelIcon,
    showCloseIcon,
    zIndex,
    element,
    iconName,
    iconColor,
    confirmText,
    cancelText,
    confirmProps,
}) => {
    return messageBox({
        role: 'sweetui-message2.info',
        transparentRegion,
        restrictedDragRange,
        element,
        title,
        message,
        detail,
        type: (<SweetIcon name={iconName || 'info'} size={32} color={iconColor || '#1890FF'} />),
        icons: showCloseIcon ? [
            {
                icon: <SweetIcon name={'x'} size={16} />,
                onClick: ({ detail: resolve }) => resolve(false),
            },
        ] : [],
        buttons: showCancelIcon ?
            [
                {
                    text: confirmText || __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                    ...confirmProps,
                },
                {
                    text: cancelText || __('取消'),
                    theme: 'regular',
                    onClick: ({ detail: resolve }) => resolve(false),
                },
            ] :
            [
                {
                    text: confirmText || __('确定'),
                    theme: 'oem',
                    onClick: ({ detail: resolve }) => resolve(true),
                },
            ],
        zIndex,
    })
}