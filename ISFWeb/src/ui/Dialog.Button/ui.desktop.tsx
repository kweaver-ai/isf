import React from 'react';
import Button from '../Button/ui.desktop';
import styles from './styles.desktop';

export default function DialogButton({ children, ...props }: UI.DialogButton.Props) {
    return (
        <Button className={styles['button']} minWidth={80} {...props}>
            {
                children
            }
        </Button >
    )
}