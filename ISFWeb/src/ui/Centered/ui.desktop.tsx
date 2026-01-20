import React from 'react'
import styles from './styles.desktop'

/**
 * 内容居中组件
 * @param props.children
 */
const Centered = function Centered({ role, children }) {
    return (
        <div role={role} className={styles['centered']}>
            <div className={styles['positioned']}>
                <div className={styles['content']}>
                    {
                        children
                    }
                </div>
            </div>
        </div>
    )
}

export default Centered