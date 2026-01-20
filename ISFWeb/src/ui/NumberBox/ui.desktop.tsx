import React from 'react';
import classnames from 'classnames';
import UIIcon from '../UIIcon/ui.desktop';
import Control from '../Control/ui.desktop';
import NumberBoxBase from './ui.base';
import TextInput from '../TextInput/ui.desktop';
import ValidateTip from '../ValidateTip/ui.desktop';
import Locator from '../Locator/ui.desktop';
import styles from './styles.desktop';

export default class NumberBox extends NumberBoxBase {

    render() {
        return (
            <Control
                className={classnames({ [styles['validate-fail']]: this.state.validateState in this.state.validateMessages }, this.props.className)}
                width={this.props.width}
                style={this.props.style}
                focus={this.state.focus}
                disabled={this.props.disabled}
            >
                <div className={styles['box-display']}>
                    <div className={styles['input-size']}>
                        <TextInput
                            value={this.state.value}
                            disabled={this.props.disabled}
                            autoFocus={this.props.autoFocus}
                            validator={this.props.validator}
                            onChange={(value) => { this.handleChange(value) }}
                            onFocus={this.props.onFocus}
                            onBlur={this.props.onBlur}
                            onEnter={this.props.onEnter}
                            onKeyDown={this.props.onKeyDown && this.props.onKeyDown.bind(this)}
                        />
                    </div>
                    <div className={styles['set-btn']}>
                        <div key={'add'}>
                            <UIIcon
                                size={13}
                                code={'\uf019'}
                                onMouseDown={this.addValue.bind(this)}
                                // onMouseUp={this.clearTimer.bind(this)}
                                // onMouseOut={this.clearTimer.bind(this)}
                                disabled={this.props.disabled}
                                className={styles['mouse-status']}
                            />
                        </div>
                        <div key={'sub'}>
                            <UIIcon
                                size={13}
                                code={'\uf01a'}
                                onMouseDown={this.subValue.bind(this)}
                                // onMouseUp={this.clearTimer.bind(this)}
                                // onMouseOut={this.clearTimer.bind(this)}
                                disabled={this.props.disabled}
                                className={styles['mouse-status']}
                            />
                        </div>
                    </div>

                    {
                        this.state.validateState in this.state.validateMessages ?
                            <div className={classnames(styles['tip-wrap'], styles['right'])} >
                                {/* {
                                    <UIIcon code={'\uf033'} size={12} className={styles['warning-icon']} />
                                } */}
                                <Locator className={styles['locator']}>
                                    <div className={styles['validate-message']}>
                                        <ValidateTip align={'right'}>
                                            {
                                                this.state.validateMessages[this.state.validateState]
                                            }
                                        </ValidateTip>
                                    </div>
                                </Locator>
                            </div>
                            : null
                    }
                </div>

            </Control>
        )
    }
}