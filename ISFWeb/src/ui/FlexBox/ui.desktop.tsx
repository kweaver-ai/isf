import React from 'react';
import FlexBoxItem from '../FlexBox.Item/ui.desktop';
import styles from './styles.desktop';

const FlexBox: UI.FlexBox.Component = function FlexBox({ role, children }) {
    return (
        <div role={role} className={styles['flex']}>
            {
                children
            }
        </div>
    )
} as UI.FlexBox.Component

FlexBox.Item = FlexBoxItem;

export default FlexBox