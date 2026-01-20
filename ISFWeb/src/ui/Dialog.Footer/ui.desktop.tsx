import React from 'react';
import styles from './styles.desktop';

const DialogFooter: UI.DialogFooter.StatelessComponent = function DialogFooter({ children }) {
    return (
        <div className={ styles['container'] }>
            <div className={ styles['padding'] } role="drag-area">
                {
                    children
                }
            </div>
        </div>
    )
}

export default DialogFooter;