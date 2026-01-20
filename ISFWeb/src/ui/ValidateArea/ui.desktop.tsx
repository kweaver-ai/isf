import React from 'react';
import classnames from 'classnames';
import { omit } from 'lodash';
import UIIcon from '../UIIcon/ui.desktop';
import TextAreaInput from '../TextAreaInput/ui.desktop';
import ValidateTip from '../ValidateTip/ui.desktop';
import Locator from '../Locator/ui.desktop';
import Control from '../Control/ui.desktop';
import ValidateAreaBase from './ui.base';
import styles from './styles.desktop';

export default class ValidateArea extends ValidateAreaBase {
    render() {
        let { height, width, align, validateState, validateMessages, className, disabled } = this.props
        return (
            <Control
                className={classnames(styles['validate-box'], { [styles['validate-fail']]: validateState in validateMessages }, className)}
                disabled={disabled}
                focus={this.state.focus}
                width={width}
                height={height}
            >
                <TextAreaInput
                    {...omit(this.props, 'className') }
                    onBlur={this.blur.bind(this)}
                    onFocus={this.focus.bind(this)}
                    onMouseover={this.mouseOver.bind(this)}
                    onMouseout={this.mouseOut.bind(this)}
                />
                {
                    validateState in validateMessages ?
                        <div className={classnames(styles['tip-wrap'], styles[align])} >
                            {
                                align === 'right' ?
                                    <UIIcon code={'\uf033'}
                                        size={16}
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