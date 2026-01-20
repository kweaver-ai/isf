import React from 'react';
import styles from './styles.desktop';

export default function Icon({ role, url, size }: UI.Icon.Props) {
    return (
        <img role={role} className={styles['icon']} src={url} style={{ width: size, height: size }} />
    )
}