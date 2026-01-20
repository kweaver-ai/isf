import * as React from 'react';
import { useState, useRef, useEffect } from 'react';
import classNames from 'classnames';
import { EllipsisProps } from './type';
import  styles from './styles.view.css';

const Ellipsis: React.FC<EllipsisProps> = ({
    className = '',
    children,
}) => {
    const ref = useRef(null);

    const [showTitle, setShowTitle] = useState<boolean>(false);

    useEffect(() => {
        if (ref.current && ref.current!.scrollWidth > ref.current!.clientWidth) {
            setShowTitle(true);
        }
    }, []);

    return (
        <span
            ref={ref}
            title={typeof children === 'string' && showTitle ? children : ''}
            className={classNames(styles['container'], className)}
        >
            {children}
        </span>
    )
}

export default React.memo(Ellipsis);