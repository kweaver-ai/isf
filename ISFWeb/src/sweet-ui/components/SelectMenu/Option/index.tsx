import React from 'react';
import MenuItem from '../../Menu/MenuItem';

export interface SelectMenuOptionProps {
    className?: string;

    icon: string;

    disabled: boolean;

    selected?: boolean;

    value: any;

    title?: string;

    onSelect: (selected: any) => void;
}

const SelectMenuOption: React.FunctionComponent<SelectMenuOptionProps> = function SelectMenuOption({
    selected,
    disabled,
    value,
    ...otherProps
}) {
    return (
        <MenuItem
            icon={typeof value !== 'undefined' && !disabled && selected ? 'selected' : 'empty'}
            {...otherProps}
        />
    );
};

export default SelectMenuOption;
