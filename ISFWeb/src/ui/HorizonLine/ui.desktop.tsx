import React from 'react';
import classnames from 'classnames';
import { ClassName } from '../helper';
import styles from './styles.desktop';

export default function HorizonLine({ height = 1 } = {}): UI.HorizonLine.Element {
    return (
        <hr
            style={ { borderBottomWidth: height } }
            className={ classnames(styles['hr'], ClassName.BorderColor) }
        />
    )
}