import React from 'react';
import classnames from 'classnames';
import { isFunction } from 'lodash';
import { Text } from '@/ui/ui.desktop';
import View from '../View';
import SweetIcon from '../SweetIcon';
import styles from './styles';

interface TagProps {
    /**
     * 标签样式
     */
    className?: string;

    /**
     * CSS样式
     */
    style?: React.CSSProperties;

    /**
     * 标签是否可以关闭，配合 onClose 使用，处理关闭事件
     */
    closable?: boolean;

    /**
     * 点击关闭按钮时的回调函数
     */
    onClose?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 标签是否支持选择，配合 checked 和 onChange 使用，处理选择事件
     */
    checkable?: boolean;

    /**
     * 标签的选中状态，配合 onChange 一起使用
     */
    checked?: boolean;

    /**
     * 标签的选中状态发生变化时的回调函数，配合 checkable 属性一起使用
     */
    onChange?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 禁用标签
     */
    disabled?: boolean;

    /**
    * title组件的className
    */
    titleClassName?: string;
}

const Tag: React.FunctionComponent<TagProps> = function Tag({
    className,
    style,
    closable = false,
    onClose,
    checkable = false,
    checked = false,
    onChange,
    disabled = false,
    children,
    titleClassName,
}) {
    const handleClose = (e: React.MouseEvent<HTMLElement>) => {
        if (disabled) {
            return
        }

        e.stopPropagation();

        onClose && isFunction(onClose) && onClose(e)
    };

    const handleChange = (e: React.MouseEvent<HTMLElement>) => {
        if (disabled || !checkable) {
            return
        }

        onChange && isFunction(onChange) && onChange(e)
    }

    return (
        <View
            inline={true}
            className={classnames(
                styles['tag-wrapper'],
                {
                    [styles['closable']]: closable,
                },
                {
                    [styles['checkable']]: checkable,
                },
                {
                    [styles['checked']]: checked,
                },
                {
                    [styles['tag-disabled']]: disabled,
                },
                className,
            )}
            style={style}
            onClick={handleChange}
        >
            <Text titleClassName={titleClassName}>{children}</Text>
            {
                closable ?
                    <View inline={true} onClick={handleClose} className={styles['close-wrapper']}>
                        <SweetIcon name={'x'} size={12} className={styles['close-icon']} />
                    </View> :
                    null
            }
        </View >
    )
}

export default Tag