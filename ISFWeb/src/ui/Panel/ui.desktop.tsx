import React from 'react';
import PanelMain from '../Panel.Main/ui.desktop';
import PanelFooter from '../Panel.Footer/ui.desktop';
import PanelButton from '../Panel.Button/ui.desktop';

const Panel: UI.Panel.Component = function Panel({ role, children }) {
    return (
        <div role={role}>
            {
                children
            }
        </div>
    )
} as UI.Panel.Component

Panel.Main = PanelMain;

Panel.Footer = PanelFooter;

Panel.Button = PanelButton;

export default Panel;