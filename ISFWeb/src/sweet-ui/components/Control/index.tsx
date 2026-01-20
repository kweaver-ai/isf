import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import styles from './styles';

interface ControlProps {
    /**
     * className
     */
    className?: string;

    /**
     * style
     */
    style?: React.CSSProperties;

    /**
     * width，包含盒模型的padding和border
     */
    width?: number | string;

    /**
     * height，包含盒模型的padding和border
     */
    height?: number | string;

    /**
     * minHeight，包含盒模型的padding和border
     */
    minHeight?: number | string;

    /**
     * maxHeight，包含盒模型的padding和border
     */
    maxHeight?: number | string;

    /**
     * 是否默认聚焦
     */
    focus?: boolean;

    /**
     * 是否禁用
     */
    disabled?: boolean;

    /**
     * 点击触发
     */
    onClick?: (event: SweetUIEvent<MouseEvent>) => any;

    /**
     * 失焦触发
     */
    onBlur?: (event: SweetUIEvent<MouseEvent>) => any;
}

const Control: React.SFC<ControlProps> = function Control({
    className,
    style,
    width,
    height,
    maxHeight,
    minHeight,
    focus,
    disabled,
    onClick = noop,
    onBlur = noop,
    children,
}) {
    const dispatchClickEvent = createEventDispatcher(onClick);
    const dispatchBlurEvent = createEventDispatcher(onBlur);

    return (
        <View
            className={
                classnames(
                    styles['control'],
                    {
                        [styles['disabled']]: disabled,
                        [styles['focus']]: focus,
                        // 如果传递了宽或高，则将宽高视为盒模型整体宽高，计算时将包含padding/border
                        [styles['box-sizing-border-box']]: !!width || !!height || !!maxHeight || !!minHeight,
                    },
                    className,
                )
            }
            style={{ ...style, width, height, minHeight, maxHeight }}
            onClick={dispatchClickEvent}
            onBlur={dispatchBlurEvent}
        >
            {
                children
            }
        </View>
    )
}

export default Control;
