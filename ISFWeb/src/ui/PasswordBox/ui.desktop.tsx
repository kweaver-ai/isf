import React from 'react';
import Control from '../Control/ui.desktop';
import PasswordInput from '../PasswordInput/ui.desktop';
import PasswordBoxBase from './ui.base';

export default class PasswordBox extends PasswordBoxBase {

    render() {
        const { style, className, width, disabled, ...props } = this.props;

        return (
            <Control
                className={ className }
                style={ style }
                width={width}
                disabled={ disabled }
                focus={ this.state.focus }
            >
                <PasswordInput
                    { ...props }
                    disabled={ disabled }
                    onFocus={ this.focus.bind(this) }
                    onBlur={ this.blur.bind(this) }
                />
            </Control >
        )
    }
}