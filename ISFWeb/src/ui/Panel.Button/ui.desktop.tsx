import React from 'react';
import Button from '../Button/ui.desktop';
import styles from './styles.desktop';

const PanelButton: UI.PanelButton.Component = function PanelButton({ children, ...props }) {
    return (
        <Button
            className={styles['button']}
            minWidth={80}
            {...props}
        >
            {
                children
            }
        </Button >
    )
}

export default PanelButton