import React, { useContext } from 'react';
import classnames from 'classnames';
import { SweetUIEvent } from '../../utils/event';
import Trigger from '../Trigger';
import Item from '../Menu/MenuItem';
import styles from './styles';
import AppConfigContext from '@/core/context/AppConfigContext';

interface PopMenuProps {
    /**
     * 下拉根元素的类名称
     */
    className?: string;

    /**
     * 弹出层展开时是否冻结滚动条
     */
    freeze?: boolean;

    /**
     * 锚点元素
     */
    anchor?: HTMLElement;

    /**
     * 触发元素定位原点
     */
    anchorOrigin?: [number | string, number | string];

    /**
     * 弹出元素定位原点
     */
    alignOrigin?: [number | string, number | string];

    /**
     * 定义如何渲染触发元素
     */
    trigger?: (
        props: Partial<{
            setPopupVisibleOnMouseEnter: () => void;
            setPopupVisibleOnMouseLeave: () => void;
            setPopupVisibleOnClick: () => void;
            setPopupVisibleOnFocus: () => void;
            setPopupVisibleOnBlur: () => void;
        }>,
    ) => React.ReactNode;

    open?: boolean;

    onRequestCloseWhenClick: (close: () => void, event?: Event) => void;
    onRequestCloseWhenBlur: (event: SweetUIEvent<any>) => void;
    onPopMenuMouseDown: () => void;
}

const PopMenu: React.FunctionComponent<PopMenuProps> = function PopMenu({
    children,
    className,
    anchor,
    open,
    alignOrigin,
    anchorOrigin,
    freeze,
    trigger,
    onRequestCloseWhenClick,
    onRequestCloseWhenBlur,
    onPopMenuMouseDown,
    ...otherProps
}) {

    const { element } = useContext(AppConfigContext)

    return (
        <Trigger
            element={element}
            renderer={trigger}
            onBeforePopupClose={onRequestCloseWhenBlur}
            {...{ anchor, anchorOrigin, alignOrigin, freeze, open }}
            {...otherProps}
        >
            {({ close }) => (
                <ul
                    className={classnames(styles['popmenu'], className)}
                    onMouseDown={onPopMenuMouseDown}
                    onClick={(e) =>
                        typeof onRequestCloseWhenClick === 'function' ? onRequestCloseWhenClick(close, e) : undefined}
                >
                    {children}
                </ul>
            )}
        </Trigger>
    );
};

PopMenu.Item = Item;

export default PopMenu;
