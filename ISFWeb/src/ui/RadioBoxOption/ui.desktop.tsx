import React from 'react';
import classnames from 'classnames';
import RadioBox from '../RadioBox/ui.desktop';
import styles from './styles.desktop';

export default function RadioBoxOption({ children, className, ...props }: UI.RadioBoxOption.Props) {
    return (
        <label className={classnames(styles['container'], className)}>
            <RadioBox
                className={styles['radio-box']}
                {...props}
            />
            <span className={classnames(styles['text'], { [styles['disabled']]: props.disabled })}>
                {
                    children
                }
            </span>
        </label>
    )
}