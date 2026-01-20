import React from 'react';
import { noop } from 'lodash';

export default class TextAreaBase extends React.PureComponent<UI.TextArea.Props, UI.TextArea.State> {
    static defaultProps = {

        validator: () => true,

        onChange: noop,

        onFocus: noop,

        onBlur: noop,

        showCounter: false,
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

    updateValue(value) {
        this.setState({
            value,
        })
    }

    changeHandler(event) {
        const value = event.target.value;

        if ((!this.props.required && value === '') || (this.props.validator && this.props.validator(value))) {
            this.updateValue(value);
            this.props.onChange && this.props.onChange(value);
        } else {
            event.preventDefault();
        }
    }

    focusHandler() {
        this.setState({ focus: true })
        this.props.onFocus && this.props.onFocus();
    }

    blurHandler() {
        this.setState({ focus: false })
        this.props.onBlur && this.props.onBlur();
    }
}