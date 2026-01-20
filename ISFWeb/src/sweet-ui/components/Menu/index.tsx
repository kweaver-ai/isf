import React from 'react';
import classnames from 'classnames';
import MenuItem from './MenuItem';
import styles from './styles';

export interface MenuProps {
    width?: number;

    maxHeight?: number;

    close?: () => void;

    className?: string;
}

const Menu: React.FunctionComponent<MenuProps> = function Menu({ children, width, maxHeight, className, close, ...otherProps }) {
    return (
        <ul
            className={classnames(styles['drop-menu'], { [styles['box-sizing-border-box']]: !!width }, className)}
            style={{ width: `${width}px`, maxHeight: `${maxHeight}px` }}
            {...otherProps}
        >
            {children}
        </ul>
    );
};

Menu.Item = MenuItem;

export default Menu;
