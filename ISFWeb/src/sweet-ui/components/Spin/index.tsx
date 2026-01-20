import React from 'react';
import { range } from 'lodash'
import classnames from 'classnames'
import styles from './styles';

type Size = 'small' | 'middle' | 'large'

interface Props {
    /**
     * 样式名称
     */
    className?: string;

    /**
     * 大小
     */
    size?: Size;
}

/**
 * Spin
 */
const Spin = ({
    className,
    size = 'middle',
}: Props) => {
    return (
        <div className={classnames(styles['container'], className)}>
            <span className={classnames(styles['dot'], styles[`dot-${size}`])}>
                {
                    range(0, 4).map((item) => (
                        <i
                            key={item}
                            className={styles['dot-item']}
                        ></i>
                    ))
                }
            </span>
        </div>
    )
}

export default Spin