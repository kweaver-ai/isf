import React from 'react';
import styles from './styles.desktop';

const MenuItem: UI.MenuItem.Component = function MenuItem({ children, onClick }) {
    return (
        <div
            className={ styles['container'] }
            onClick={ onClick }
        >
            { children }
        </div>
    )
}

export default MenuItem;