import React from 'react';
import classnames from 'classnames';
import { ClassName } from '@/ui/helper';
import View from '../../View';
import SweetIcon from '../../SweetIcon';
import CheckBox from '../../CheckBox';
import styles from './styles';

export interface MenuItemProps {
    className?: string;

    checkbox: object;

    icon: any;

    close: () => void;

    disabled: boolean;

    selected?: boolean;

    title?: string;

    value: any;

    onClick?: () => void;

    label?: string;
}

const MenuItem: React.FunctionComponent<MenuItemProps> = function MenuItem({
    children,
    className,
    title,
    value,
    checkbox,
    icon,
    disabled = false,
    selected = false,
    onClick,
    label,
    ...otherProps
}) {
    return (
        <li
            className={classnames(
                styles['drop-menu-item'],
                {
                    [styles['disabled']]: disabled,
                },
                {
                    [styles['drop-menu-item-active']]: selected,
                },
                { [ClassName.Color__Hover]: !disabled },
                { [ClassName.Color]: selected },
                className,
            )}
            key={value}
            onClick={onClick}
            {...otherProps}
        >
            {checkbox ? (
                <CheckBox
                    disabled={!!checkbox.disabled}
                    checked={!!checkbox.checked}
                    onChange={checkbox.onChange}
                />
            ) : null}
            {
                typeof icon !== 'undefined' ?
                    typeof icon === 'string' ?
                        <SweetIcon color={disabled ? '#c8c8c8' : '#757575'} name={icon} /> :
                        icon :
                    null
            }
            {
                // 兼容PopMenu.Item
                label ?
                    <View
                        inline={true}
                        className={classnames({ [styles['label']]: checkbox || icon })}
                    >
                        {label}
                    </View>
                    : null
            }
            {
                children ?
                    <View
                        inline={true}
                        className={classnames({ [styles['label']]: checkbox || icon })}
                    >
                        {children}
                    </View>
                    : null
            }
        </li>
    );
};

export default MenuItem;
