import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash';
import styles from './styles.desktop';

export default function Control({ role, className, style, width, height, maxHeight, minHeight, focus, disabled, onClick = noop, onBlur = noop, onMouseLeave= noop, onMouseOver= noop, children }: UI.Control.Props): UI.Control.Element {
    return (
        <div
            role={role}
            className={
                classnames(
                    styles['control'],
                    className,
                    {
                        [styles['disabled']]: disabled,
                        [styles['focus']]: focus,
                        [styles['box-sizing-border-box']]: !!width || !!height || !!maxHeight || !!minHeight, // 如果传递了宽或高，则将宽高视为盒模型整体宽高，计算时将包含padding/border
                    },
                )
            }
            style={{ ...style, width, height, minHeight, maxHeight }}
            onClick={(event) => onClick(event)}
            onBlur={(event) => onBlur(event)}
            onMouseLeave={(event) => onMouseLeave(event)}
            onMouseOver={(event) => onMouseOver(event)}
        >
            {
                children
            }
        </div>
    )
}