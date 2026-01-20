import React from 'react';
import classnames from 'classnames';
import styles from './styles.desktop';

export default ({ className, align = 'right', children, borderColor = '#bdbdbd', backgroundColor = '#fafafa' }) => (
    <div className={classnames(className, styles['tip'])}>
        <div
            className={align === 'right' ?
                classnames(styles['pointer'], styles['pointer-right'])
                : classnames(styles['pointer'], styles['pointer-top'])}
            style={align === 'right' ?
                { borderColor: `transparent ${borderColor} transparent transparent` }
                : { borderColor: `${borderColor} transparent transparent transparent` }}>
            <div
                className={align === 'right' ?
                    classnames(styles['pointer-inner'], styles['pointer-inner-right'])
                    : classnames(styles['pointer-inner'], styles['pointer-inner-top'])}
                style={align === 'right' ?
                    { borderColor: `transparent ${backgroundColor} transparent transparent` }
                    : { borderColor: `${backgroundColor} transparent transparent transparent` }}>
            </div>
        </div>
        <span className={align === 'right' ?
            classnames(styles['content'])
            : classnames(styles['content'], styles['nowarp'])}
        style={{ borderColor, backgroundColor }}>
            {
                children
            }
        </span>
    </div>
)
