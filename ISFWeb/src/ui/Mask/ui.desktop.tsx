import React from 'react';
import styles from './styles.desktop';

const Mask = function Mask({ children }) {
    return (
        <div className={styles['mask']}></div>
    )
}

export default Mask;