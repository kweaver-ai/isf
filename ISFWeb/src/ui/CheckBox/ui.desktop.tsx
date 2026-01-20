import React from 'react';
import classnames from 'classnames';
import CheckBoxBase from './ui.base';
import styles from './styles.desktop';

export default class CheckBox extends CheckBoxBase {
    render() {
        return (
            <input
                role={this.props.role}
                type="checkbox"
                id={this.props.id}
                className={classnames(
                    styles.checkbox,
                    { [styles['disabled']]: this.props.disabled },
                    this.props.className,
                )}
                checked={this.state.checked}
                disabled={this.props.disabled}
                onChange={this.changeHandler.bind(this)}
                onClick={this.handleClick.bind(this)}
                ref={(checkbox) => this.checkbox = checkbox}
            />
        )
    }
}