import React from 'react';
import classnames from 'classnames';
import FlexTextBoxBase from './ui.base';
import styles from './styles.desktop';

export default class FlexTextBox extends FlexTextBoxBase {
    render() {
        return (
            <div
                ref={(node) => this.textBox = node}
                className={classnames(
                    styles['textbox'],
                    {
                        [styles['placeholder']]: this.state.placeholder !== '',
                        [styles['disabled']]: this.props.disabled,
                    },
                    this.props.className,
                )}
                data-placeholder={this.state.placeholder}
                onKeyDown={this.keyDownHandler.bind(this)}
                onPaste={this.pasteHandler.bind(this)}
                onBlur={this.props.onBlur}
                contentEditable={!this.props.disabled}
            />
        )
    }
}