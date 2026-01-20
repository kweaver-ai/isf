import React from 'react';
import classnames from 'classnames';
import PasswordInputBase from './ui.base';
import styles from './styles.desktop';

export default class PasswordInput extends PasswordInputBase {
    render() {
        return (
            <input
                id={this.props.id}
                type="text"
                className={classnames(styles['input'], this.props.className)}
                ref={(input) => this.input = input}
                value={this.state.value}
                readOnly={this.props.readOnly}
                disabled={this.props.disabled}
                placeholder={this.props.placeholder}
                onChange={this.changeHandler.bind(this)}
                onFocus={this.focusHandler.bind(this)}
                onBlur={this.blurHandler.bind(this)}
                onMouseOver={this.mouseoverHandler.bind(this)}
                onMouseOut={this.mouseoutHandler.bind(this)}
                onKeyDown={this.keyDownHandler.bind(this)}
            />
        )
    }
}