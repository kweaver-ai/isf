import React from 'react';
import classnames from 'classnames';
import FlexTextBoxBase from './ui.base';
import styles from './styles.desktop';

export default class FlexTextBox extends FlexTextBoxBase {
    render() {
        const { className, disabled, placeholder, readOnly, maxLength } = this.props;
        const { value, width } = this.state;

        const style = {
            ...this.props.style,
            width: this.props.width || `${width}px`,
            border: 'none',
        };
        return (
            <span style={{ display: 'inline-block' }}>
                <input
                    type="text"
                    ref={this.saveInput}
                    className={classnames(styles['base-input'], className)}
                    onChange={this.handleValueChange}
                    {...{ value, disabled, placeholder, readOnly, maxLength, style }}
                    onClick={this.handleClick}
                    onKeyDown={this.handleKeyDown}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onPaste={this.handlePaste}
                />
            </span>
        )
    }
}