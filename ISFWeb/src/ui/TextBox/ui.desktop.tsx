import React from 'react';
import Control from '../Control/ui.desktop';
import TextInput from '../TextInput/ui.desktop';
import TextBoxBase from './ui.base';

export default class TextBox extends TextBoxBase {

    render() {
        const { style, className, width, disabled, role, ...props } = this.props;

        return (
            <Control
                role={role}
                className={ className }
                style={ style }
                width={width}
                disabled={ disabled }
                focus={ this.state.focus }
            >
                <TextInput
                    { ...props }
                    disabled={ disabled }
                    onFocus={ this.focus.bind(this) }
                    onBlur={ this.blur.bind(this) }
                />
            </Control >
        )
    }
}