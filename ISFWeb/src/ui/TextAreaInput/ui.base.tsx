import React from 'react';
import { noop } from 'lodash';

export default class TextAreaInput extends React.PureComponent<UI.TextAreaInput.Props, UI.TextAreaInput.State> {

    textarea: HTMLTextAreaElement;

    static defaultProps = {

        validator: () => true,

        onChange: noop,

        onFocus: noop,

        onBlur: noop,

        onMouseout: noop,

        onMouseover: noop,
    }

    state = {
        value: '',

        focus: false,
    }

    componentDidMount() {
        this.updateValue(this.props.value)
    }

    static getDerivedStateFromProps({ value }, prevState) {
        if(value !== prevState.value) {
            return {
                value,
            }
        }
        return null
    }

    /**
     * props的值更新
     * @param value
     */
    private updateValue(value) {
        this.setState({
            value,
        })
    }

    /**
     * 处理onChange事件
     * @param event
     */
    protected changeHandler(event) {
        const value = event.target.value;

        if ((!this.props.required && value === '') || (this.props.validator && this.props.validator(value))) {
            this.updateValue(value);
            this.props.onChange && this.props.onChange(value);
        } else {
            event.preventDefault();
        }
    }

    /**
     * 处理onFocus事件
     */
    protected focusHandler() {
        this.setState({ focus: true })
        this.props.onFocus && this.props.onFocus();
    }

    /**
     * 处理onBlur事件
     */
    protected blurHandler() {
        this.setState({ focus: false })
        this.props.onBlur && this.props.onBlur();
    }

    /**
     * 处理onClick事件
     */
    protected clickHandler() {
        if (!this.props.disabled) {
            this.setState({ focus: true })
            this.textarea.focus()
        }
    }

    /**
     * 处理onMouseover事件
     * @param event
     */
    protected mouseoverHandler(event) {
        if (!this.props.disabled) {
            this.props.onMouseover && this.props.onMouseover(event);
        }
    }

    /**
     * 处理onMouseout事件
     * @param event
     */
    protected mouseoutHandler(event) {
        if (!this.props.disabled) {
            this.props.onMouseout && this.props.onMouseout(event);
        }
    }
}