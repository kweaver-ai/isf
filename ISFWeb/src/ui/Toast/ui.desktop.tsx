import React from 'react'
import { noop } from 'lodash'
import UIIcon from '../UIIcon/ui.desktop'
import styles from './styles.desktop'

const Toast: React.FunctionComponent<UI.Toast.Props> = ({ children, closable, code, onClose, ...otherProps }) => (
    <div className={styles['toast']} data-test-scope="ui/Toast">
        <div className={styles['toast-bg']}></div>
        {
            code ? <UIIcon code={code} {...otherProps} className={styles['icon']} /> : null
        }
        <span className={styles['text']}>
            {
                children
            }
        </span>
        {
            closable ? <UIIcon className={styles['close']} onClick={onClose} code={'\uf046'} size={14} /> : null
        }
    </div>
)

Toast.defaultProps = {
    closable: false,
    onClose: noop,
}

export default Toast