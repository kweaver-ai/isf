import React from 'react';
import FontIcon from '../FontIcon/ui.desktop';
import '@/core/fonts/ui/font.css';

const UIIcon: React.FunctionComponent<UI.UIIcon.Props> = function UIIcon({
    role,
    code,
    fallback,
    size = 16,
    cursor,
    ...otherProps
}) {
    return (
        <FontIcon
            role={role}
            font="AnyShare"
            code={code}
            size={size}
            cursor={cursor}
            fallback={fallback} {...otherProps}
        />
    )
}

export default UIIcon