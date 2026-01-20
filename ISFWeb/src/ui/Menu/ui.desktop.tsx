import React from 'react';
import { isFunction } from 'lodash'
import classnames from 'classnames';
import MenuItem from '../Menu.Item/ui.desktop';
import styles from './styles.desktop';

const Menu: UI.Menu.Component = function Menu({ width, maxHeight, children, onMouseDown, role }) {
    return (
        <div
            role={role}
            className={ classnames(styles['menu'], { [styles['box-sizing-border-box']]: !!width }) }
            onMouseDown={ (e) => isFunction(onMouseDown) && onMouseDown(e) }
            style={ { width, maxHeight } }
        >
            {
                children
            }
        </div>
    )
} as UI.Menu.Component

Menu.Item = MenuItem;

export default Menu