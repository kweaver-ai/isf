import React from 'react';
import classnames from 'classnames';
import ValidateBox from '../ValidateBox/ui.desktop';
import styles from './styles.desktop';

export default function InlineValidateBox({ className, width, ...props }: UI.InlineValidateBox.Props) {
    return (
        <ValidateBox width={width ? width : 48} className={classnames(styles['inline-validatebox'], className)} {...props} />
    )
}