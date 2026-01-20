import React from 'react';
import classnames from 'classnames';
import { ClassName } from '../helper';
import styles from './styles.desktop';

const TabsTab: React.FunctionComponent<UI.TabsTab.Props> = function TabsTab({ role, children, active, style, onActive, className }) {
    return (
        <div
            role={role}
            className={classnames(
                styles['tab'],
                [ClassName.Color__Hover],
                {
                    [styles['active']]: active,
                    [ClassName.BorderBottomColor]: active,
                    [ClassName.Color]: active,
                },
                className,
            )}
            style={{ ...style }}
            onClick={onActive}
        >
            {
                children
            }
        </div>
    )
}

export default TabsTab