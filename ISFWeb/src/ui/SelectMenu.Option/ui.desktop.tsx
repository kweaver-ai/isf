import React from 'react'
import PopMenuItem from '../PopMenu.Item/ui.desktop'

const SelectMenuOption: React.FunctionComponent<any> = function SelectMenuOption({ selected, disabled, value, ...otherProps }) {
    return (
        <PopMenuItem
            icon={typeof value !== 'undefined' && !disabled && selected ? '\uf068' : '\u0000'}
            {...otherProps}
        />
    )
}

export default SelectMenuOption