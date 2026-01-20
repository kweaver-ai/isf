import React from 'react';
import classnames from 'classnames';
import { Title } from '@/ui/ui.desktop';
import BaseInput from '../../BaseInput';
import SweetIcon from '../../SweetIcon';
import styles from './styles';

interface SingleSelectorProps {
    /**
     * 文本框中显示的内容
     */
    label?: string;

    /**
     * 宽度
     */
    width?: number;

    /**
     * className
     */
    className?: string;

    /**
     * css样式，传入TextInput，控制TextInput样式
     */
    style?: React.CSSProperties;

    /**
     * placeholder
     */
    placeholder?: string;

    /**
     * 是否聚焦状态
     */
    active?: boolean;

    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * 文本框状态
     */
    status?: 'normal' | 'error';

    /**
     * 点击时触发
     */
    onClick?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 渲染完成后触发
     */
    onMounted?: (ref: HTMLInputElement) => void;

    /**
     * 输入框聚焦时触发
     */
    onFocus?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 输入框失焦时触发
     */
    onBlur?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 鼠标进入文本域时触发
     */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;
}

export default class SingleSelector extends React.Component<SingleSelectorProps, any> {
    static defaultProps = {
        status: 'normal',
    }

    render() {
        const { label, width, className, style, placeholder, status, active, disabled, onClick, onMounted, onFocus, onBlur, onMouseEnter, onMouseLeave, role } = this.props;

        return (
            <Title
                content={label}
                timeout={1500}
                role={role}
            >
                <div
                    className={classnames(
                        styles['wrapper'],
                        {
                            [styles['disabled']]: disabled,
                            [styles[`${status}`]]: status,
                            [styles[`${status}-active`]]: active,
                        },
                        className,
                    )}
                    style={{ width }}
                    onClick={onClick}
                    onMouseEnter={onMouseEnter}
                    onMouseLeave={onMouseLeave}
                >
                    {/*
                        Edge禁用状态title提示不消失临时解决方案
                        禁用时，Edge中使用input高频率偶现不触发Title的mouseLeave事件
                        替换Title后检查是否可以去掉这边的判断
                    */}
                    {
                        disabled ?
                            <div
                                className={classnames(
                                    styles['disable-input'],
                                    {
                                        [styles['empty']]: !(label || placeholder),
                                    },
                                )}
                            >
                                {label || placeholder || ''}
                            </div>
                            :
                            <BaseInput
                                type="text"
                                className={classnames(
                                    styles['input'],
                                    {
                                        [styles['placeholder']]: !label && placeholder,
                                    },
                                    {
                                        [styles['disabled']]: disabled,
                                    },
                                )}
                                value={label}
                                readOnly={true}
                                style={style}
                                {...{ disabled, placeholder }}
                                onMounted={onMounted}
                                onFocus={onFocus}
                                onBlur={onBlur}
                            />
                    }
                    <SweetIcon
                        name={'arrowDown'}
                        size={16}
                        className={classnames(
                            styles['arrow-down'],
                            {
                                [styles['disabled']]: disabled,
                            },
                        )}
                    />
                    {
                        status === 'error' && !disabled ?
                            <SweetIcon
                                name={'caution'}
                                size={16}
                                color={'#e60012'}
                                className={styles['caution']}
                            /> : null
                    }
                </div>
            </Title>
        );
    }
}