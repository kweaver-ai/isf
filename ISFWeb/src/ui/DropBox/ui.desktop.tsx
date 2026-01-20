import React from 'react';
import classnames from 'classnames';
import DropBoxBase from './ui.base';
import Control from '../Control/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import Text from '../Text/ui.desktop';
import PopOver from '../PopOver/ui.desktop'
import styles from './styles.desktop';
import arrowDown from './assets/arrow-down.png';

export default class DropBox extends DropBoxBase {
    render() {
        const { width, height, maxHeight, icon, role } = this.props;

        const { active } = this.state;
        return (
            <PopOver
                role={role}
                trigger={
                    <div
                        ref="container"
                        className={styles['dropbox']}
                        style={{ width, height, maxHeight }}
                    >
                        <Control
                            focus={active}
                            className={classnames(styles['control'], this.props.className)}
                            width={width}
                            height={height}
                            maxHeight={maxHeight}
                            disabled={this.props.disabled}
                        >
                            <span
                                ref="select"
                                href="#"
                                className={classnames(styles['select'], { [styles['disabled']]: this.props.disabled })}
                                onMouseDown={(e) => { e.stopPropagation();   this.toggleActive.bind(e) }}
                                onBlur={this.onSelectBlur.bind(this)}
                            >
                                {
                                    icon ? <span className={styles['icon']} >
                                        <UIIcon
                                            size={16}
                                            code={icon}
                                        />
                                    </span > :
                                        null
                                }
                                <div className={classnames(styles['text'], { [styles['text-left']]: icon })}>

                                    <Text>
                                        {
                                            this.props.formatter(this.props.value)
                                        }
                                    </Text>
                                </div>
                                <span className={styles['drop-icon']}>
                                    <UIIcon
                                        disabled={this.props.disabled}
                                        size={16}
                                        code={this.props.fontIcon}
                                        fallback={this.props.fallbackIcon || arrowDown}
                                    />
                                </span>
                            </span>
                        </Control>
                    </ div >
                }
                /** 兼容angular弹框中的下拉框，防止弹窗覆盖下拉框 */
                style={{ zIndex: 10000 }}
                triggerEvent="click"
                anchorOrigin={['left', 'bottom']}
                targetOrigin={['left', 'top']}
                freezable={false}
                watch={true}
                onRequestCloseWhenBlur={(close) => this.onClose(close)}
                onRequestCloseWhenClick={this.toggleSelect.bind(this)}
                open={active}
                onOpen={() => this.onOpen()}
            >
                <div
                    className={classnames(styles['drop'], { [styles['active']]: active })}
                    ref="drop"
                    onMouseDown={this.preventDeactivate.bind(this)}
                >
                    {
                        this.props.children
                    }
                </div>

            </PopOver>

        )

    }

}