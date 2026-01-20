import React from 'react';
import classnames from 'classnames';
import { omit } from 'lodash';
import UIIcon from '../UIIcon/ui.desktop';
import TextInput from '../TextInput/ui.desktop';
import ValidateTip from '../ValidateTip/ui.desktop';
import Locator from '../Locator/ui.desktop';
import Control from '../Control/ui.desktop';
import ValidateBoxBase from './ui.base';
import styles from './styles.desktop';

export default class ValidateBox extends ValidateBoxBase {
    render() {
        let { role, width, align, validateState, validateMessages, className, disabled, type, isHidePwdEyes, readOnly = false } = this.props
        return (
            <Control
                role={role}
                className={classnames(styles['validate-box'], { [styles['validate-fail']]: validateState in validateMessages }, className)}
                disabled={disabled}
                focus={this.state.focus}
                width={width}
            >
                <TextInput
                    {...omit(this.props, 'className')}
                    className={classnames({ [styles['secret-pwd']]: isHidePwdEyes })}
                    onBlur={this.blur.bind(this)}
                    onFocus={this.focus.bind(this)}
                    onMouseover={this.mouseOver.bind(this)}
                    onMouseout={this.mouseOut.bind(this)}
                    onEnter={this.onEnter.bind(this)}
                    readOnly={readOnly}
                    type={type}
                    ref={(ref) => this.textInput = ref}
                />
                {
                    validateState in validateMessages ?
                        <div className={classnames(styles['tip-wrap'], styles[align])} >
                            {
                                align === 'right' ?
                                    <UIIcon
                                        code={'\uf033'} size={16}
                                        className={styles['warning-icon']}
                                        onMouseOver={this.mouseOver.bind(this)}
                                        onMouseLeave={this.mouseOut.bind(this)}
                                    />
                                    : null
                            }
                            {
                                this.state.focus || this.state.hover ?
                                    (
                                        <Locator className={styles['locator']}>
                                            <div className={styles['validate-message']}>
                                                <ValidateTip align={align}>
                                                    {
                                                        validateMessages[validateState]
                                                    }
                                                </ValidateTip>
                                            </div>
                                        </Locator>
                                    ) :
                                    null
                            }
                        </div>
                        : null
                }
            </Control>
        )
    }
}