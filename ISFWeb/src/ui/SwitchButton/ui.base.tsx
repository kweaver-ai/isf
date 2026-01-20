import React from 'react';
import { noop, isBoolean } from 'lodash';

export default class SwitchButtonBase extends React.Component<any, any> {

    static defaultProps = {
        active: false,
        onChange: noop,
        onCheck: noop,
        onUncheck: noop,
    }

    state = {
        active: this.props.active,
    }

    static getDerivedStateFromProps({ active }, prevState) {
        if(isBoolean(active) && active !== prevState.active) {
            return {
                active,
            }
        }
        return null
    }

    toggleStatus(value, active) {
        if (active) {
            this.props.onCheck(value);
        } else {
            this.props.onUncheck(value);
        }
        this.props.onChange(active, value)
        this.setState({
            active,
        })
    }
}