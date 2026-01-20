import React from 'react';
import PlainButton from '../PlainButton/ui.desktop';
import styles from './styles.desktop';

const ToolBarButton: UI.ToolBarButton.Component = function ToolBarButton({ children, ...props }) {
    return (
        <PlainButton
            className={ styles['toolbar-button'] }
            {...props}
        >
            {
                children
            }
        </PlainButton>
    )
}

export default ToolBarButton;