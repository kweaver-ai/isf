import React from 'react';
import classnames from 'classnames';
import TextAreaBase from './ui.base';
import Control from '../Control/ui.desktop';
import styles from './styles.desktop';

export default class TextArea extends TextAreaBase {
    render() {
        const { showCounter, maxLength, role } = this.props
        const { value } = this.state

        return (
            <Control
                role={role}
                className={classnames(
                    styles['textarea'],
                    {
                        [styles['count-limit']]: showCounter,
                    },
                    this.props.className,
                )}
                disabled={this.props.disabled}
                focus={this.state.focus}
                width={this.props.width}
                height={this.props.height}
                maxHeight={this.props.maxHeight}
                minHeight={this.props.minHeight}
            >
                <textarea
                    className={styles['input']}
                    value={this.state.value}
                    maxLength={this.props.maxLength}
                    placeholder={this.props.placeholder}
                    readOnly={this.props.readOnly}
                    disabled={this.props.disabled}
                    onChange={this.changeHandler.bind(this)}
                    onFocus={this.focusHandler.bind(this)}
                    onBlur={this.blurHandler.bind(this)}
                />

                {
                    showCounter ?
                        maxLength ?
                            <span
                                className={classnames(
                                    styles['count'],
                                    {
                                        [styles['overflow']]: value.length > maxLength,
                                    },
                                )}
                            >
                                {`${value.length}/${maxLength}`}
                            </span>
                            :
                            <span
                                className={styles['count']}
                            >
                                {value.length}
                            </span>
                        :
                        null
                }
            </Control>
        )
    }
}