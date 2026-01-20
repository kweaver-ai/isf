import * as React from 'react';
import classnames from 'classnames';
import { Control, UIIcon, Text, CheckBox, DatePicker, Menu } from '@/ui/ui.desktop';
import { Trigger } from '@/sweet-ui';
import ValidityBox2Base from './component.base';
import styles from './styles.view';
import __ from './locale';

export default class ValidityBox2 extends ValidityBox2Base {
    render() {
        const { active, value } = this.state;
        const DateBox = ({ setPopupVisibleOnClick, role }) => (
            <div
                role={role}
                key={'ValidityBoxDateBox'}
                className={styles['dropbox']}
                style={{ width: this.props.width }}
            >
                <Control
                    role={'ui-control'}
                    focus={active}
                    className={classnames(styles['control'])}
                    width={this.props.width}
                    disabled={this.props.disabled}
                >
                    <span
                        ref={(select) => (this.select = select)}
                        href="#"
                        className={classnames(styles['select'], { [styles['disabled']]: this.props.disabled }, this.props.className)}
                        onClick={(e) => {
                            e.preventDefault();
                            if (!this.props.disabled) {
                                setPopupVisibleOnClick(e);
                            }
                        }}
                        onBlur={this.onSelectBlur.bind(this)}
                    >
                        <div className={styles['text']}>
                            <Text role={'ui-text'}>{this.validityFormatter(value === -1 ? value : value / 1000)}</Text>
                        </div>
                        <span className={styles['drop-icon']}>
                            <UIIcon
                                role={'ui-uiicon'}
                                size={16}
                                code={'\uf00e'}
                                color={this.props.disabled ? '#c8c8c8' : '#505050'}
                            />
                        </span>
                    </span>
                </Control>
            </div>
        );

        if (this.props.disabled) {
            return (
                <div className={styles['dropbox']} style={{ width: this.props.width }}>
                    <Control
                        role={'ui-control'}
                        focus={active}
                        className={classnames(styles['control'])}
                        width={this.props.width}
                        disabled={true}
                    >
                        <a
                            ref={(select) => (this.select = select)}
                            href="#"
                            className={classnames(styles['select'], styles['disabled'])}
                        >
                            <div className={styles['text']}>
                                <Text role={'ui-text'}>{this.validityFormatter(value === -1 ? value : value / 1000)}</Text>
                            </div>
                            <span className={styles['drop-icon']}>
                                <UIIcon
                                    role={'ui-uiicon'}
                                    size={16}
                                    code={'\uf00e'}
                                    color={'#c8c8c8'}
                                />
                            </span>
                        </a>
                    </Control>
                </div>
            );
        }

        return (
            <Trigger
                role={'sweetui-trigger'}
                renderer={DateBox}
                triggerEvent="click"
                anchorOrigin={['left', 'bottom']}
                alignOrigin={['left', 'top']}
                freeze={false}
                /** 兼容angular弹框中的下拉框，防止弹窗覆盖下拉框 */
                popupZIndex={10000}
                onPopupVisibleChange={(event) => this.setState({ active: event.detail })}
            >
                {({ close }) => (
                    <Menu role={'ui-menu'}>
                        <DatePicker
                            role={'ui-datepicker'}
                            value={value === -1 ? null : new Date(value / 1000)}
                            selectRange={this.props.selectRange}
                            onChange={(value) => this.setValidity(value, close)}
                            disabled={this.state.value === -1}
                        />
                        <div className={styles['options']}>
                            {this.props.allowPermanent ? (
                                <div>
                                    <CheckBox
                                        role={'ui-checkBox'}
                                        value={-1}
                                        checked={this.state.value === -1}
                                        onChange={this.switchPermanent.bind(this, close)}
                                    />
                                    <label className={styles['option-label']}>{__('永久有效')}</label>
                                </div>
                            ) : null}
                        </div>
                    </Menu>
                )}
            </Trigger>
        )
    }
}
