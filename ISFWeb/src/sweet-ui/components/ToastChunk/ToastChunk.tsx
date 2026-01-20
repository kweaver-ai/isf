import React from 'react';
import ReactDOM from 'react-dom';
import classnames from 'classnames'
import { noop, isFunction } from 'lodash';
import SweetIcon from '../../components/SweetIcon/index';
import styles from './styles';

interface ToastChunkProps {
    /**
     * 标签里的内容
     */
    content: HTMLElement | string;

    /**
     * toastChunk持续时间(ms)
     */
    duration?: number;

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
     * 挂载ToastChunk的容器
     */
    container?: HTMLElement;

    /**
     * 关闭
     */
    onClose?: () => void;

    /**
     * 其他
     */
    [x: string]: any;
}

interface ToastChunkState {
    /**
     * 淡入淡出样式
     */
    gradualStyle: GradualStyle;
}

/**
 * 渐入渐出样式
 */
export const enum GradualStyle {
    /**
     * 渐入
     */
    ToastIn,

    /**
     * 渐出
    */
    ToastOut,
}

export const enum ToastMode {
    /**
     * 白色背景
     */
    Light = 'light',

    /**
     * 黑色背景
     */
    Black = 'black',
}

/**
 * 提示类型
 */
export const enum NoticeType {
    /**
     * 无图标
     */
    None = 'none',

    /**
     * 成功
     */
    Success = 'success',

    /**
     * 警告
     */
    Warning = 'warning',

    /**
     * 失败
     */
    Error = 'error',

    /**
     * 消息
     */
    Info = 'info',
}

/**
 * 提示图标
 */
const TypeIcon = {
    [NoticeType.Success]: {
        name: 'successTip',
        color: '#52c41a',
    },
    [NoticeType.Error]: {
        name: 'errorTip',
        color: '#ff4d4f',
    },
    [NoticeType.Warning]: {
        name: 'warningTip',
        color: '#faad14',
    },
    [NoticeType.Info]: {
        name: 'infoTip',
        color: '#1890ff',
    },
}

export default class ToastChunk extends React.Component<ToastChunkProps, ToastChunkState> {
    static defaultProps = {
        duration: 2000,
        onClose: noop,
        mode: ToastMode.Black,
        type: NoticeType.None,
    }

    state = {
        gradualStyle: GradualStyle.ToastIn,
    }

    closeTimer: number | null = null

    componentDidMount() {
        this.startCloseTimer()
    }

    /**
     * 开始定时器
     */
    private startCloseTimer: () => void = () => {
        if (this.props.duration) {
            this.closeTimer = window.setTimeout(() => {
                this.close()
            }, this.props.duration)
        }
    }

    /**
     * 清除定时器
     */
    private clearCloseTimer: () => void = () => {
        if (this.closeTimer) {
            clearTimeout(this.closeTimer);
            this.closeTimer = null
        }
    }

    /**
     * 关闭
     */
    private close: () => void = () => {
        this.setState(() => ({
            gradualStyle: GradualStyle.ToastOut,
        }))

        setTimeout(() => {
            this.clearCloseTimer()

            isFunction(this.props.onClose) && this.props.onClose()
        }, 100)
    }

    render() {
        const {
            className,
            content,
            container,
            mode,
            type,
            ...otherProps
        } = this.props

        const node = (
            <div>
                <div
                    className={classnames(
                        mode === ToastMode.Black ? styles['black-chunk'] : styles['light-chunk'],
                        this.state.gradualStyle === GradualStyle.ToastIn ? styles['toast-in'] : styles['toast-out'],
                        className,
                    )}
                    {...otherProps}
                >
                    {
                        type && type !== NoticeType.None ?
                            [
                                (
                                    <SweetIcon
                                        key={'toast-icon'}
                                        {...TypeIcon[type]}
                                        className={styles['icon']}
                                    />
                                ),
                                content,
                            ]
                            : content
                    }
                </div>
            </div>
        )

        if (container) {
            return ReactDOM.createPortal(node, container)
        }

        return node
    }
}