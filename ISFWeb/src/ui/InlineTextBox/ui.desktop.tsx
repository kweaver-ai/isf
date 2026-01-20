import React from 'react';
import classnames from 'classnames';
import TextBox from '../TextBox/ui.desktop';
import styles from './styles.desktop';

export default function InlineTextBox({ className, ...props }: UI.InlineTextBox.Props) {
    return (
        <TextBox className={classnames(styles['inline-textbox'], className)} {...props} />
    )
}