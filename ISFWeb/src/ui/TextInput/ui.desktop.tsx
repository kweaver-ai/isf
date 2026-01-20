import React from 'react';
import classnames from 'classnames';
import { isBrowser, Browser } from '@/util/browser';
import TextInputBase from './ui.base';
import styles from './styles.desktop';

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class TextInput extends TextInputBase {
    render() {
        return (
            <span>
                {
                    this.props.type === 'password' ?
                        <input
                            ref={(input) => this.input = input}
                            id={this.props.id}
                            className={classnames(styles['input'], { [styles['disabled']]: this.props.disabled && !isSafari }, this.props.className)}
                            style={this.props.style}
                            autoComplete="off"
                            type={this.props.type}
                            placeholder={this.props.placeholder}
                            readOnly={this.props.readOnly}
                            disabled={this.props.disabled}
                            maxLength={this.props.maxLength}
                            onChange={this.changeHandler.bind(this)}
                            onFocus={this.focusHandler.bind(this)}
                            onBlur={this.blurHandler.bind(this)}
                            onClick={this.clickHandler.bind(this)}
                            onKeyDown={this.keyDownHandler.bind(this)}
                            onMouseOver={this.mouseoverHandler.bind(this)}
                            onMouseOut={this.mouseoutHandler.bind(this)}
                        />
                        :
                        <input
                            ref={(input) => this.input = input}
                            id={this.props.id}
                            className={classnames(styles['input'], { [styles['disabled']]: this.props.disabled && !isSafari }, this.props.className)}
                            style={this.props.style}
                            autoComplete="off"
                            type={this.props.type}
                            value={this.state.value}
                            placeholder={this.props.placeholder}
                            readOnly={this.props.readOnly}
                            disabled={this.props.disabled}
                            maxLength={this.props.maxLength}
                            onChange={this.changeHandler.bind(this)}
                            onFocus={this.focusHandler.bind(this)}
                            onBlur={this.blurHandler.bind(this)}
                            onClick={this.clickHandler.bind(this)}
                            onKeyDown={this.keyDownHandler.bind(this)}
                            onMouseOver={this.mouseoverHandler.bind(this)}
                            onMouseOut={this.mouseoutHandler.bind(this)}
                        />
                }
            </span>
        )
    }
}