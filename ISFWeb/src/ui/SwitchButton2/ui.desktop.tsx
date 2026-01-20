
import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash';
import styles from './styles.desktop';

const SwitchButton2: React.FunctionComponent<UI.SwitchButton2.Props> = function SwitchButton2({ active, onChange = noop, disabled = false, value }) {

    return (
        <div
            className={
                classnames(
                    styles['button-style'],
                    { [styles['disabled']]: disabled },
                )
            }
        >
            <span onClick={() => { !disabled && onChange(value, !active) }}>
                <div
                    className={
                        classnames(
                            styles['slide'],
                            { [styles['slide-on']]: active },
                            { [styles['slide-close']]: !active },
                        )
                    }
                >

                </div>
                <div
                    className={
                        classnames(
                            styles['btn'],
                            { [styles['btn-on']]: active },
                            { [styles['btn-close']]: !active },

                        )
                    }
                >

                </div>
            </span>
        </div>
    )
}

export default SwitchButton2;