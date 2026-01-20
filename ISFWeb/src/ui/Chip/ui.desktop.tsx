import React from 'react';
import classnames from 'classnames';
import { noop, isFunction } from 'lodash';
import Text from '../Text/ui.desktop';
import styles from './styles.desktop';

export default function Chip({ actionClassName, className, maxWidth, removeHandler, readOnly = false, disabled = false, children, onClick = noop }: UI.Chip.Props) {
    return (
        <div className={classnames(styles['chip'], className, { [styles['disabled']]: disabled })} onClick={onClick}>
            <div className={styles['text']} style={{ maxWidth }}>
                <Text>
                    {
                        children
                    }
                </Text>
            </div>
            {
                !readOnly && isFunction(removeHandler) ?
                    <span href="#" className={classnames(styles['action'], actionClassName)} onClick={(e) => { e.stopPropagation(); !disabled && removeHandler() }} >x</span>
                    : null
            }
        </div>
    )
}