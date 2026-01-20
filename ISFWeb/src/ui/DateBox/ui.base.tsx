import React from 'react';
import { noop, isFunction } from 'lodash';

export default class DateBoxBase extends React.PureComponent<UI.DateBox.Props, UI.DateBox.State> implements UI.DateBox.Component {
    static defaultProps = {
        active: false,

        format: 'yyyy-MM-dd',

        dropAlign: 'bottom left',

        onChange: noop,

        onActive: noop,

        selectRange: [],

        shouldShowblankStatus: false,

        startsFromZero: false,

        onDatePickerClick: noop,

        disabled: false,

        placeholder: '---',
    }

    state = {
        value: this.props.value || new Date(),

        active: this.props.active,
    }

    static getDerivedStateFromProps({ value }, prevState) {
        if (value !== prevState.value) {
            return {
                value,
                active: false,
            }
        }
        return null
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.state.value !== prevState.value) {
            this.fireChangeEvent(this.state.value);
        }
    }

    /**
     * 更新选中日期
     * @param value 日期对象
     */
    private updateDate(value: Date) {
        this.setState({
            value,
            active: false,
        });

        this.fireChangeEvent(value);
    }

    /**
     * 选中日期时触发
     * @param value 日期对象
     */
    private fireChangeEvent(value: Date) {
        isFunction(this.props.onChange) && this.props.onChange(value);
    }

    /**
     * 选中日期时触发
     * @param value 日期对象
     */
    protected select(value, close?) {
        this.updateDate(value);
        isFunction(close) && close()
    }

    protected handleBeforePopupClose = () => {
        this.props.onDatePickerClick(false)
    }
}