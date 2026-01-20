import React from 'react';
import ToolBarButton from '../ToolBar.Button/ui.desktop';
import styles from './styles.desktop';

const ToolBar: UI.ToolBar.Component = function ToolBar({ children }) {
    return (
        <div className={ styles['tool-bar'] } role={'ui-toolbar'}>
            {
                children
            }
        </div>
    )
} as UI.ToolBar.Component

ToolBar.Button = ToolBarButton;

export default ToolBar