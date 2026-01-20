import React from 'react';
import classnames from 'classnames';
import styles from './styles';

/**
 * 分割字符串
 */
function splitStr(name: string) {
    const sliceLen = 7
    const nameLen = name.length
    const firstEnd = nameLen - sliceLen

    const strBefore = name.slice(0, firstEnd)
    const strAter = name.slice(firstEnd, nameLen)
    return [strBefore, strAter]
}

interface MidEllipsisProps {
    children: React.ReactElement;
    classNames?: string;
}

const MidEllipsis = React.memo<MidEllipsisProps>(({
    children,
    classNames,
    ...otherProps
}) => {
    if (typeof children === 'string') {
        const [nameBefore, nameAfter] = splitStr(children)

        return (
            <div className={classnames(styles['name'], classNames)} {...otherProps}>
                <span className={styles['name-before']}>{nameBefore}</span>
                <span className={styles['name-after']}>{nameAfter}</span>
            </div>
        )
    }

    return children
})

export default MidEllipsis;