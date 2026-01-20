import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash'
import UIIcon from '../UIIcon/ui.desktop';
import PaneBase from './ui.base';
import PopMenu from '../PopMenu/ui.desktop'
import styles from './styles.desktop';

export default class Pane extends PaneBase {

    render() {
        const { icon, color, type = 'local', fallback, label, disabled, menuItems, msgNum, triggerEvent = 'click' } = this.props;

        const pane = (
            <div
                className={classnames(
                    styles['container'],
                    styles[type],
                    {
                        [styles['disabled']]: disabled,
                    },
                )}
                onClick={this.handleClick.bind(this)}
                onBlur={this.fireBlurEvent.bind(this)}
            >
                <div
                    className={styles['pane']}
                >
                    <UIIcon
                        code={icon}
                        color={color}
                        fallback={fallback}
                        size={24}
                        disabled={disabled}
                        onClick={noop /* 防止鼠标手型变为指针 */}
                    />
                    <div className={styles['label']}>
                        {label}
                    </div>
                </div>
                {
                    msgNum !== 0 ?
                        <div className={classnames(styles['badge'], { [styles['circles']]: msgNum < 10, [styles['oval']]: msgNum > 99 })}>
                            {msgNum > 99 ? '99+' : msgNum}
                        </div>
                        : null
                }
                {
                    (triggerEvent !== 'mouseover' && menuItems && menuItems.length) ?
                        (
                            <PopMenu
                                anchor={this.state.anchor}
                                anchorOrigin={['right', 'bottom']}
                                targetOrigin={['right', 'top']}
                                autoFix={false}
                                freezable={true}
                                open={this.state.open}
                                onRequestCloseWhenClick={this.closePopMenu.bind(this)}
                                onRequestCloseWhenBlur={this.closePopMenu.bind(this)}
                                triggerEvent={'click'}
                            >
                                {
                                    menuItems
                                }
                            </PopMenu>
                        ) :
                        null
                }
            </div>
        )

        if (triggerEvent === 'mouseover' && menuItems && menuItems.length) {
            return (
                <PopMenu
                    anchorOrigin={['right', 'bottom']}
                    targetOrigin={['right', 'top']}
                    autoFix={false}
                    freezable={false}
                    onRequestCloseWhenClick={(close) => close()}
                    closeWhenMouseLeave={true}
                    triggerEvent={triggerEvent}
                    trigger={pane}
                >
                    {
                        menuItems
                    }
                </PopMenu>
            )
        } else {
            return pane
        }
    }

}