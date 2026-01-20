import React from 'react';
import classnames from 'classnames';
import { ClassName } from '../helper';
import { noop } from 'lodash';
import Text from '../Text/ui.desktop';
import styles from './styles.desktop';

const SelectOption: UI.SelectOption.Component = function SelectOption({ className, selected, value, disabled, onSelect = noop, children, role }) {
    return (
        <div
            role={role}
            className={
                classnames(
                    styles['option'],
                    {
                        [styles['selected']]: selected,
                    },
                    {
                        [ClassName.Color]: selected,
                    },
                    className,
                )
            }
            onMouseDown={!disabled && onSelect.bind(null, { value, text: children })}
        >
            <Text>
                {
                    children
                }
            </Text>
        </div>
    )
}

export default SelectOption