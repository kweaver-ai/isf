import React from 'react';
import classnames from 'classnames';
import { getAlignCls } from './helper';
import styles from './styles.desktop';

export default function FlexBoxItem({ role, children, width, align = 'left middle', className }: UI.FlexBoxItem.Props) {
    return (
        <div
            role={role}
            className={classnames(styles['item'],
                getAlignCls(styles, align), className)}
            style={{ width }}
        >
            {
                children
            }
        </div>
    )
}