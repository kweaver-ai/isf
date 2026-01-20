import React from 'react';
import { createRoot } from 'react-dom/client';
import { isFunction } from 'lodash';
import { isDom } from '@/util/validators'
import ToastChunk, { ToastMode, NoticeType } from '../ToastChunk/ToastChunk';
import styles from './styles';

let seed = 0

export function getUuid() {
    return `toast_${Date.now()}_${seed++}`;
}

/**
 * toastChunk的信息
 */
export interface ToastInfo {
    /**
     * 标签内容
     */
    content: JSX.Element | string;

    /**
     * toast持续时间
     */
    duration?: number;

    /**
     * toastChunk唯一标识
     */
    key?: string | number;

    /**
     * 样式
     */
    className?: string;

    /**
     * 模式
     */
    mode?: ToastMode | string;

    /**
     * 提示类型
     */
    type?: NoticeType | string;

    /**
     * 关闭函数
     */
    onClose?: () => void;
}

interface ToasterProps {
    /**
     * 列表最多toastCuhnk个数
     */
    maxCount?: number;

    /**
     * 容器
     */
    holder?: any;
}

interface ToasterState {
    /**
     * toastChunk列表
     */
    toasts: ReadonlyArray<ToastInfo>;
}

class Toaster extends React.Component<ToasterProps, ToasterState> {
    static newInstance

    static defaultProps = {
        maxCount: 10,
    }

    state = {
        toasts: [],
    }

    /**
     * 添加toastChunk
     */
    private add = (toast: ToastInfo): void => {
        toast.key = toast.key || getUuid()

        const { key } = toast

        const { maxCount } = this.props

        this.setState((prevState) => {
            const toasts = prevState.toasts

            if (maxCount && maxCount > toasts.length) {
                if (!toasts.filter((t) => t.key === key).length) {
                    return {
                        toasts: toasts.concat(toast),
                    }
                }
            }
        })
    }

    /**
     * 移除指定的ToastChunk
     */
    private remove = (key: string | number): void => {
        this.setState((prevState) => ({
            toasts: prevState.toasts.filter((toast) => toast.key !== key),
        }))
    }

    /**
     * 持续时间结束自动关闭ToastChunk
     */
    private close = (toast: ToastInfo): void => {
        if (toast && toast.key) {
            this.remove(toast.key)

            isFunction(toast.onClose) && toast.onClose()
        }
    }

    render() {
        const { toasts } = this.state

        const toastNodes = toasts.map((toast, index) => {
            return (
                <ToastChunk
                    {...toast}
                    key={toast.key}
                    onClose={() => this.close(toast)}
                />
            )
        })

        return (
            <div className={styles['toaster']}>
                {toastNodes}
            </div>
        )
    }
}

/**
 * Toaster实例化函数
 */
Toaster.newInstance = function newToasterInstance(properties: ToasterProps, callback: (methods: any) => void) {
    const props = properties || {}

    const div = window.document.createElement('div')

    const parent = props.holder && isDom(props.holder) ? props.holder : window.document.body

    parent.appendChild(div)

    // 创建根节点
    const root = createRoot(div);

    let called = false

    function ref(toaster) {
        if (called) {
            return
        }

        called = true

        callback({
            addToast(props: ToastInfo) {
                toaster.add(props)
            },
            removeToast(key: string | number) {
                toaster.remove(key)
            },
            destroy() {
                root.unmount();
                div.parentNode.removeChild(div);
            },
        })
    }

    root.render(<Toaster {...props} ref={ref} />);
}

export default Toaster