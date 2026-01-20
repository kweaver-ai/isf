import React from 'react';
import { SweetUIEvent, createEventDispatcher } from '../../utils/event';
import Trigger, { TriggerEvent } from '../Trigger';
import Menu from '../Menu';
import SelectMenuOption from './Option';

interface SelectMenuProps {
    triggerEvent: TriggerEvent;

    freeze: boolean;
    /**
      * 触发元素定位原点
      */
    anchorOrigin?: [number | string, number | string];

    /**
     * 弹出元素定位原点
     */
    alignOrigin?: [number | string, number | string];

    label: (
        props: Partial<{
            setPopupVisibleOnMouseEnter: () => void;
            setPopupVisibleOnMouseLeave: () => void;
            setPopupVisibleOnClick: () => void;
        }>,
    ) => React.ReactNode;

    value?: any;

    className?: string;

    onChange?: (event: SweetUIEvent<any>) => void;

    onRequestCloseWhenBlur: () => void;

    onRequestCloseWhenClick: (close: () => void) => void;

    onPopMenuMouseDown: () => void;
}

interface SelectMenuState {
    value: any;
}

class SelectMenu extends React.Component<SelectMenuProps, SelectMenuState> {
    static defaultProps = {
        triggerEvent: 'hover',
        freeze: false,
        anchorOrigin: ['left', 'bottom'],
        alignOrigin: ['left', 'top'],
        onRequestCloseWhenClick: (close: () => void) => close(),
    };

    state = {
        value: this.props.value,
    };

    static getDerivedStateFromProps({ value }: SelectMenuProps, prevState: SelectMenuState) {
        if (typeof value !== 'undefined' && value !== prevState.value) {
            return {
                value,
            };
        }

        return null;
    }

    handleClick(e, item) {
        if (!item.props.disabled) {
            if (typeof item.props.onClick === 'function') {
                item.props.onClick(e);
            }
            if (!e.defaultPrevented) {
                this.setState({ value: item.props.value });
                this.dispatchChangeEvent(item.props.value);
            }
        }
    }

    dispatchChangeEvent = createEventDispatcher(this.props.onChange);

    render() {
        const {
            children,
            label,
            className,
            onRequestCloseWhenBlur,
            onRequestCloseWhenClick,
            onPopMenuMouseDown,
            ...otherProps
        } = this.props;

        return (
            <Trigger renderer={label} onBeforePopupClose={onRequestCloseWhenBlur} {...otherProps}>
                {({ close }) => (
                    <Menu
                        className={className}
                        onMouseDown={onPopMenuMouseDown}
                        onClick={() =>
                            typeof onRequestCloseWhenClick === 'function' ? onRequestCloseWhenClick(close) : undefined}
                    >
                        {React.Children.map(children, (option) =>
                            React.cloneElement(option, {
                                selected: option.props.value === this.state.value,
                                disabled: option.props.disabled,
                                onClick: (e) => this.handleClick(e, option),
                            }),
                        )}
                    </Menu>
                )}
            </Trigger>
        );
    }
}

SelectMenu.Option = SelectMenuOption;

export default SelectMenu;
