import React from 'react';
import styles from './styles.desktop';

export default function DataGridCell({ children }: UI.DataGridCell.Props) {
    return (
        <td className={ styles['cell'] }>
            {
                children
            }
        </td>
    )
}