import React from 'react';
import classnames from 'classnames';
import TextAreaInputBase from './ui.base';
import styles from './styles.desktop';

export default class TextAreaInput extends TextAreaInputBase {
    render() {
        return (
            <div className={styles['textarea']}>
                <textarea
                    ref={(textarea) => this.textarea = textarea}
                    className={styles['input']}
                    value={this.state.value}
                    maxLength={this.props.maxlength}
                    readOnly={this.props.readOnly}
                    disabled={this.props.disabled}
                    onChange={this.changeHandler.bind(this)}
                    onFocus={this.focusHandler.bind(this)}
                    onBlur={this.blurHandler.bind(this)}
                    onMouseOver={this.mouseoverHandler.bind(this)}
                    onMouseOut={this.mouseoutHandler.bind(this)}
                />
                <span
                    onClick={this.clickHandler.bind(this)}
                    className={classnames(styles['label-placeholder'], { [styles['hide-lable']]: !!this.state.value })}
                >{this.props.placeholder}</span>
            </div >
        )
    }
}