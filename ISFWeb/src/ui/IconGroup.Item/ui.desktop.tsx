import React from 'react';
import classnames from 'classnames';
import InlineButton from '../InlineButton/ui.desktop';
import styles from './styles.desktop';

const IconGroupItem: React.FunctionComponent<UI.IconGroupItem.Props> = function IconGroupItem({
    className,
    code,
    disabled,
    ...otherProps
}) {
    return (
        <InlineButton
            code={code}
            disabled={disabled}
            className={classnames(styles['icon'], { [styles['enabled']]: !disabled }, className)}
            {...otherProps}
        />
    );
};

export default IconGroupItem;
