import React from 'react';
import classnames from 'classnames';
import { render } from 'react-dom';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import SweetIcon from '../../components/SweetIcon';
import Button from '../../components/Button';
import ModalDialog from '../../components/ModalDialog';
import View from '../../components/View';
import { Portal } from '../../components';
import { noop, isFunction } from 'lodash'
import __ from './locale';
import styles from './styles';

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 */
type Alert = (params: { title?: string; message?: string; zIndex?: number }) => Promise<boolean>

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 */
type Info = (params: { title?: string; message?: string; zIndex?: number }) => Promise<boolean>

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.subtitle 内容标题
 * @param params.message 消息
 */
type Confirm = (params: { title?: string; message?: string; zIndex?: number }) => Promise<boolean>

/**
 * @param params 对话框参数
 * @param params.title 对话框标题
 * @param params.message 消息
 * @param params.overflow 是否滚动显示消息
 */
type Error = (params: { title?: string; message?: string; overflow?: boolean; zIndex?: number }) => Promise<boolean>

/**
 * 点击确定/取消后返回的结果
 */
type ActionResult = (resolve: boolean) => Promise<boolean>

type MessageBoxButton = { text: string; onClick: (event: SweetUIEvent<ActionResult>) => void; theme: 'regular' | 'oem' | 'gray' | 'dark' | 'text' | 'noborder' }

/**
 * 消息对话框
 * @param params.icon 图标
 * @param params.title 对话框标题
 * @param params.subtitle 内容标题
 * @param params.message 消息
 * @param params.overflow 是否滚动显示消息
 * @param params.buttons 按钮
 * @param params.onClose 关闭对话框时执行
 * @param params.zIndex 对话框zIndex值
 */
type MessageBox = (params: { icon: React.ReactElement; title: string; message?: string; overflow: boolean; buttons: ReadonlyArray<MessageBoxButton>; onClose: (event: SweetUIEvent<ActionResult>) => void; zIndex?: number }) => Promise<boolean>

/**
 * 消息对话框
 * @param params
 */
const messageBox: MessageBox = ({ icon, title, message, overflow = false, buttons = [], onClose, zIndex = 20 }) => {
    let container: HTMLDivElement | null = null;
    let unbindEvent: () => void = noop;
    let rootElement: HTMLElement | null = null; 

    const box = new Promise((resolve, reject) => {
        const MessageBoxComponent = () => {
            const contextElement = document.querySelector('#isf-web-plugins')
            rootElement = (contextElement as HTMLElement) || document.body;

            const container = document.createElement('div');
            container.style.position = 'fixed'
            container.style.top = '0'
            container.style.right = '0'
            container.style.bottom = '0'
            container.style.left = '0'
            container.style.zIndex = String(zIndex) // 适应ShareWebUI/Mask的临时处理

            const destroy = () => {
                if(container && rootElement) {
                    rootElement.removeChild(container);
                }
            }
            document.activeElement && document.activeElement.blur();
            rootElement.appendChild(container);

            return (
                <Portal getContainer={() => container}>
                    <ModalDialog
                        target={container}
                        width={400}
                        title={__('提示')}
                        draggable={false}
                        onRequestClose={() => {
                            createEventDispatcher(onClose, destroy)(resolve)
                        }}
                    >
                        <View>
                            <View className={styles['content']}>
                                <View>
                                    <View className={styles['icon']}>{icon}</View>
                                    <View className={styles['article']}>
                                        {
                                            title ?
                                                (
                                                    <View className={styles['title']}>{title}</View>
                                                )
                                                : null
                                        }
                                        {
                                            message ?
                                                (
                                                    <View className={classnames(styles['message'], { [styles['message-with-scroll']]: overflow })}>{message}</View>
                                                )
                                                : null
                                        }
                                    </View>
                                </View>
                            </View>
                            <View className={styles['footer']}>
                                {
                                    buttons.map(({text, theme = 'regular', onClick}) => (
                                        <View key={text} className={styles['button']} >
                                            <Button
                                                theme={theme}
                                                style={{ minWidth: 80 }}
                                                onClick={() => {
                                                    createEventDispatcher(onClick, destroy)(resolve)
                                                }}
                                            >
                                                {text}
                                            </Button>
                                        </View>
                                    ))
                                }
                            </View>
                        </View>
                    </ModalDialog>
                </Portal>
            )
        }

        render(
            <MessageBoxComponent />,
            document.createElement('div'),
        )
    }) as Promise<boolean> & { destroy: Function };

    box.destroy = () => {
        if (container) {
            (rootElement || document.body).style.overflow = '';
            (rootElement || document.body).removeChild(container);
            container = null;
        }
        isFunction(unbindEvent) && unbindEvent();
    };

    return box
}

/**
 * 警告对话框
 * @returns Promise<boolean>
 * @example
   ```ts
   if (await Message.alert({title: '执行XX操作失败'})) {
       // doSomething
   }
   ```
 */
export const alert: Alert = ({ title, message, zIndex }) => {
    return messageBox({
        title,
        message,
        icon: <SweetIcon name="alert" size={40} color="#be0000" />,
        buttons: [
            {
                text: __('确定'),
                onClick: ({ detail: resolve }) => resolve(true),
            },
        ],
        onClose: ({ detail: resolve }) => resolve(true),
        zIndex,
    })
}

/**
 * 一般消息对话框
 * @returns Promise<boolean>
 * @example
   ```ts
   if (await Message.info({title: '执行XX操作成功'})) {
       // doSomething
   }
   ```
 */
export const info: Info = ({ title, message, zIndex }) => {
    return messageBox({
        title,
        message,
        icon: <SweetIcon name="info" size={40} color="#5a8cb4" />,
        buttons: [
            {
                text: __('确定'),
                onClick: ({ detail: resolve }) => resolve(true),
            },
        ],
        onClose: ({ detail: resolve }) => resolve(true),
        zIndex,
    })
}

/**
 * 确认对话框
 * @returns Promise<boolean>
 * @example
   ```ts
   if (await Message.confirm({title: '确认执行XX操作吗？'})) {
       // doSomething
   }
   ```
 */
export const confirm: Confirm = ({ title, message, zIndex }) => {
    return messageBox({
        title,
        message,
        icon: <SweetIcon key={'notice'} name="notice" size={40} color="#f5a415" />,
        buttons: [
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
        ],
        onClose: ({ detail: resolve }) => resolve(false),
        zIndex,
    })
}

/**
 * 错误提示框
 * @returns Promise<boolean>
 * @example
   ```ts
   if (await Message.error({title: '执行XX操作失败'})) {
       // doSomething
   }
   ```
 */
export const error: Error = ({ title, message, zIndex }) => {
    return messageBox({
        title,
        message,
        overflow: true,
        icon: <SweetIcon name="alert" size={40} color="#be0000" />,
        buttons: [
            {
                text: __('确定'),
                onClick: ({ detail: resolve }) => resolve(true),
            },
        ],
        onClose: ({ detail: resolve }) => resolve(true),
        zIndex,
    })
}