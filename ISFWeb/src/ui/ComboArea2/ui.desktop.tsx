import React from 'react';
import classnames from 'classnames';
import FlexTextBox from '../FlexTextBox/ui.desktop';
import Control from '../Control/ui.desktop';
import ComboAreaBase from './ui.base';
import Item from '../ComboArea2.Item/ui.desktop';
import styles from './styles.desktop';

export default class ComboArea extends ComboAreaBase {

    static Item = Item;
    render() {
        const { focus } = this.state;

        return (
            <Control
                className={classnames(
                    styles['comboarea'],
                    { [styles['disabled']]: this.props.disabled },
                    this.props.className,
                )}
                width={this.props.width}
                height={this.props.height}
                minHeight={this.props.minHeight}
                maxHeight={this.props.maxHeight}
                onClick={this.focusInput.bind(this)}
                onBlur={this.blurInput.bind(this)}
                focus={focus}
                disabled={this.props.disabled}
            >
                {
                    React.Children.toArray(this.props.children).map((item: React.ReactElement<UI.ComboAreaItem.Props>) => {
                        const { removeChip, ...otherProps } = item.props
                        return React.cloneElement(item, {
                            removeChip: (data) => { this.props.removeChip(data) },
                            ...otherProps,
                        })
                    })
                }

                {
                    this.props.uneditable === false ?
                        <div className={styles['chip-wrap']}>
                            <FlexTextBox
                                className={styles['flextextbox']}
                                ref="input"
                                disabled={this.props.disabled || this.props.readOnly}
                                placeholder={this.state.placeholder}
                                onKeyDown={this.keyDownHandler.bind(this)}
                            />
                        </div> :
                        null
                }
            </Control>
        )
    }
}