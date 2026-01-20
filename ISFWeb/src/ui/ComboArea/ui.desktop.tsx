import React from 'react';
import classnames from 'classnames';
import { Tag } from '@/sweet-ui'
import FlexTextBox2 from '../FlexTextBox2/ui.desktop';
import FlexTextBox from '../FlexTextBox/ui.desktop';
import Control from '../Control/ui.desktop';
import ValidateTip from '../ValidateTip/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import Locator from '../Locator/ui.desktop';
import ComboAreaBase from './ui.base';
import styles from './styles.desktop';

export default class ComboArea extends ComboAreaBase {
    render() {
        const {
            validateState,
            validateMessages,
        } = this.props

        return (
            <Control
                role={this.props.role}
                className={classnames(
                    styles['comboarea'],
                    { [styles['validate-fail']]: validateState in validateMessages },
                    this.props.className,
                )}
                width={this.props.width}
                height={this.props.height}
                minHeight={this.props.minHeight}
                maxHeight={this.props.maxHeight}
                onBlur={this.blur.bind(this)}
                onClick={this.focusInput.bind(this)}
                onMouseOver={this.mouseOver}
                onMouseLeave={this.mouseOut}
            >
                {
                    !this.state.value.length && this.props.uneditable && this.props.placeholder ?
                        <span className={styles['placeholder']}>{this.props.placeholder}</span>
                        : null
                }
                {
                    this.state.value.map((o, index) => {
                        return (
                            <div
                                className={styles['chip-wrap']}
                                key={index}
                            >
                                <Tag
                                    disabled={this.props.disabled}
                                    style={{ maxWidth: this.props.maxWidth }}
                                    closable={true}
                                    onClose={() => this.removeChip(o)}
                                >
                                    {
                                        this.props.formatter(o)
                                    }
                                </Tag>
                            </div>
                        )
                    })
                }
                {
                    this.props.uneditable === false ?
                        <div className={styles['chip-wrap']}>
                            {
                                this.props.useNewFlextInput ?
                                    <FlexTextBox2
                                        ref={this.saveFlexInput}
                                        disabled={this.props.disabled || this.props.readOnly}
                                        placeholder={this.state.value.length > 0 ? '' : this.props.placeholder}
                                        onKeyDown={this.keyDownHandler.bind(this)}
                                        onPressEnter={this.keyDownHandler.bind(this)}
                                        onPaste={this.pasteHandler.bind(this)}
                                        onBlur={this.blurHandler.bind(this)}
                                        onFocus={this.focusInput.bind(this)}
                                        onValueChange={this.handleValueChange}
                                        value={this.state.inputValue}
                                        maxWidth={this.props.maxWidth}
                                    />
                                    :
                                    <FlexTextBox
                                        ref={this.saveFlexInput}
                                        disabled={this.props.disabled || this.props.readOnly}
                                        placeholder={this.state.value.length > 0 ? '' : this.props.placeholder}
                                        onKeyDown={this.keyDownHandler.bind(this)}
                                        onPaste={this.pasteHandler.bind(this)}
                                        onBlur={this.blurHandler.bind(this)}
                                        onValueChange={this.handleValueChange}
                                    />
                            }

                        </div> :
                        null
                }
                {
                    validateState in validateMessages ?
                        <div className={classnames(styles['tip-wrap'])}>
                            <UIIcon
                                code={'\uf033'}
                                size={16}
                                className={styles['warning-icon']}
                                onMouseOver={this.mouseOver}
                                onMouseLeave={this.mouseOut}
                            />
                            {
                                this.state.isHoverWarning ?
                                    <Locator className={styles['locator']}>
                                        <div className={styles['validate-message']}>
                                            <ValidateTip align={'right'}>
                                                {validateMessages[validateState]}
                                            </ValidateTip>
                                        </div>
                                    </Locator>
                                    : null
                            }
                        </div>
                        : null
                }
            </Control>
        )
    }
}