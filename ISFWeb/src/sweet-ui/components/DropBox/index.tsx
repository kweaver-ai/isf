import React from 'react';
import { noop } from 'lodash';
import classnames from 'classnames';
import Text from '@/ui/Text/ui.desktop';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import Trigger from '../Trigger';
import View from '../View';
import SweetIcon from '../SweetIcon';
import styles from './styles';

interface DropBoxProps {
    /**
     * 下拉框宽度
     */
    width: string;

    /**
     * 下拉框图标名
     */
    iconName?: string;

    /**
     * 输入框末尾的图标名（默认为下拉箭头）
     */
    iconLabel?: string;

    /**
     * 下拉框禁用状态
     */
    disabled?: boolean;

    /**
     * 下拉框的值
     */
    value: any;

    placeholder?: string;

    /**
     * 下拉框样式
     */
    className?: string;

    freeze?: boolean;

    /**
     * 格式化下拉框显示内容
     */
    formatter?: (val: any) => string;

    onActive?: (event: SweetUIEvent<boolean>) => void;

    onBeforePopupClose?: (event: SweetUIEvent<any>) => void;

    element?: HTMLElement;
}

interface DropBoxState {
    active: boolean;
}

export default class DropBox extends React.PureComponent<DropBoxProps, DropBoxState> {
    static defaultProps = {
        formatter: (val: any) => val,

        iconLabel: 'arrowDown',

        disabled: false,

        freeze: false,

        onActive: noop,
    };

    state = {
        active: false,
    };

    handlePopupVisiableChange = (event: SweetUIEvent<boolean>) => {
        this.setState({ active: event.detail });

        this.dispatchActiveEvent(event.detail);

    };

    dispatchActiveEvent = createEventDispatcher(this.props.onActive);

    render() {
        const { width, iconName, iconLabel, disabled, value, placeholder, className, formatter, freeze } = this.props;
        const { active } = this.state;

        return (
            <Trigger
                renderer={({ setPopupVisibleOnClick }) => (
                    <View
                        key={'dropBoxTrigger'}
                        inline={true}
                        className={classnames(
                            styles['dropbox'],
                            { [styles['disabled']]: disabled },
                            { [styles['focus']]: active },
                            className,
                        )}
                        style={{ width }}
                        onClick={disabled ? undefined : setPopupVisibleOnClick}
                    >
                        {iconName ? <SweetIcon name={iconName} /> : null}
                        <View
                            inline={true}
                            className={classnames(styles['layout'], { [styles['disabled']]: disabled })}
                            placeholder={placeholder}
                        >
                            <Text>{formatter(value)}</Text>
                        </View>
                        <SweetIcon
                            name={iconLabel}
                            size={16}
                            className={classnames(styles['arrow-down'], { [styles['disabled-icon']]: disabled })}
                        />
                    </View>
                )}
                anchorOrigin={['left', 'bottom']}
                alignOrigin={['left', 'top']}
                freeze={freeze}
                popupZIndex={10000}
                open={active}
                onPopupVisibleChange={this.handlePopupVisiableChange}
                onBeforePopupClose={this.props.onBeforePopupClose}
                element={this.props.element}
            >
                {({ close, open }) => this.props.children({ close, open })}
            </Trigger>
        );
    }
}
