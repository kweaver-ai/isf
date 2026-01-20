import React from 'react';
import classnames from 'classnames';
import { isBrowser, Browser } from '@/util/browser';
import Title from '../Title/ui.desktop'
import styles from './styles.desktop';

// IE8/IE9 在HTTPS下不支持 @font-face，使用fallback图片代替
const FALLBACKED = isBrowser({ app: Browser.MSIE, version: 8 }) || (isBrowser({ app: Browser.MSIE, version: 9 }));

const FontIcon: React.FunctionComponent<UI.FontIcon.Props> = function FontIcon({
    code,
    fallback,
    font,
    title,
    element,
    size,
    color,
    cursor,
    onClick,
    disabled,
    className,
    titleClassName,
    onMouseOver,
    onMouseOut,
    onMouseDown,
    role,
    ...otherProps
}) {
    const icon = (role?) => (
        <span
            role={role}
            className={classnames(styles['icon'], {
                [styles['link']]: onClick && typeof onClick === 'function',
                [styles['disabled']]: disabled,
            }, className)}
            onMouseDown={(event) => !disabled && typeof onMouseDown === 'function' && onMouseDown(event)}
            onClick={(event) => !disabled && typeof onClick === 'function' && onClick(event)}
            onMouseOver={(event) => !disabled && typeof onMouseOver === 'function' && onMouseOver(event)}
            onMouseOut={(event) => !disabled && typeof onMouseOut === 'function' && onMouseOut(event)}
            style={{ fontFamily: font, fontSize: size, color, cursor }}
            {...otherProps}
        >
            {
                FALLBACKED || code === '\u0000' ?
                    <img
                        src={fallback}
                        className={styles['fallback-icon']}
                        width={size}
                        height={size}
                    /> :
                    code
            }
        </span>
    )

    return title ?
        <Title
            role={role}
            timeout={0}
            content={title}
            className={titleClassName}
            element={element}
        >
            {icon()}
        </Title>
        :
        icon(role)
}

export default FontIcon