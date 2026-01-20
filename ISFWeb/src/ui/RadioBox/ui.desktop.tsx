import React from 'react';
import RadioBoxBase from './ui.base';
import classnames from 'classnames';

export default class RadioBox extends RadioBoxBase {
    render() {
        return (
            <input
                type="radio"
                id={this.props.id}
                name={this.props.name}
                value={this.props.value}
                disabled={this.props.disabled}
                defaultChecked={this.state.checked}
                onChange={this.changeHandler.bind(this)}
                className={classnames(this.props.className)}
                checked={this.state.checked}
            />
        )
    }
}