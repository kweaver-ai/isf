import React from 'react';
import TextInput from '../TextInput/ui.desktop';
import SearchInputBase from './ui.base';

export default class SearchInput extends SearchInputBase {
    render() {
        return (
            <TextInput
                ref={(textInput) => this.textInput = textInput}
                className={this.props.className}
                value={this.state.value}
                disabled={this.props.disabled}
                placeholder={this.props.placeholder}
                autoFocus={this.props.autoFocus}
                validator={this.props.validator.bind(this)}
                maxLength={this.props.maxLength}
                onChange={this.handleChange.bind(this)}
                onClick={this.handleClick.bind(this)}
                onFocus={this.handleFocus.bind(this)}
                onBlur={this.handleBlur.bind(this)}
                onEnter={this.props.onEnter.bind(this)}
                onKeyDown={this.props.onKeyDown.bind(this)}
            />
        )
    }
}