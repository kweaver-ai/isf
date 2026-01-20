import React from 'react';
import classnames from 'classnames';
import { ClassName } from '@/ui/helper';
import Text from '@/ui/Text/ui.desktop';
import styles from './styles';

export interface SelectOptionProps {
    className?: string;

    close: () => void;

    disabled: boolean;

    selected?: boolean;

    title?: string;

    value: any;

    onClick?: () => void;
}

const SelectOption: React.FunctionComponent<SelectOptionProps> = function SelectOption({
    children,
    className,
    title,
    value,
    disabled = false,
    selected = false,
    onClick,
    ...otherProps
}) {
    return (
        <li
            className={classnames(
                styles['select-option'],
                {
                    [styles['disabled']]: disabled,
                },
                {
                    [styles['select-option-active']]: selected,
                },
                { [ClassName.Color__Hover]: !disabled },
                { [ClassName.Color]: selected },
                className,
            )}
            onClick={onClick}
            {...otherProps}
        >
            <Text>{children}</Text>
        </li>
    );
};

export default SelectOption;
