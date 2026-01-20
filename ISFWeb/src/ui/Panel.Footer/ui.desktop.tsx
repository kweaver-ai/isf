import React from 'react';
import styles from './styles.desktop';

const PanelFooter: UI.PanelFooter.Component = function PanelFooter({ role, children }) {
    return (
        <div
            role={role}
            className={styles['container']}
        >
            <div className={styles['padding']}>
                {
                    children
                }
            </div>
        </div>
    )
}

export default PanelFooter;