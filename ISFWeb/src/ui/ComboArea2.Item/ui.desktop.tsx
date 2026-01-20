import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash'
import Chip from '../Chip/ui.desktop';
import styles from './styles.desktop';

const ComboAreaItem: React.FunctionComponent<UI.ComboAreaItem.Props> = function ComboAreaItem({
    readOnly,
    disabled,
    data,
    children,
    className,
    chipClassName,
    removeChip = noop,
    ...otherProps
}) {

    return (
        <div
            className={classnames(styles['chip-wrap'], className)}
            {...otherProps}
        >
            <Chip
                readOnly={readOnly}
                disabled={disabled}
                removeHandler={() => { removeChip(data) }}
                className={chipClassName}
                actionClassName={styles['chip-action']}>
                {
                    children
                }
            </Chip>
        </div>
    );
}

ComboAreaItem.defaultProps = {
    readOnly: false,
    disabled: false,
    removeChip: noop,
}

export default ComboAreaItem;