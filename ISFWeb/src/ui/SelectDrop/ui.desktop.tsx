import React from 'react'
import SelectDropBase from './ui.base'
import PopMenu from '../PopMenu/ui.desktop'
import PopMenuItem from '../PopMenu.Item/ui.desktop'
import UIIcon from '../UIIcon/ui.desktop'
import { isEqual } from 'lodash';
import styles from './styles.desktop';
import classnames from 'classnames';

export default class SelectDrop extends SelectDropBase {

    render() {
        const { icon, align, iconFallback, labelFormatter, textSize, textColor,
            defaultOption, options, onChange, ...otherProps } = this.props
        return (
            <PopMenu
                trigger={
                    <span className={classnames(styles['current'])}>
                        {
                            icon ?
                                <span className={classnames(styles['icon'])}>
                                    <UIIcon
                                        code={icon}
                                        textSize={textSize}
                                        textColor={textColor}
                                    />
                                </span > : null
                        }
                        <span className={classnames(styles['text'])} >
                            {
                                this.props.labelFormatter(this.state.option)
                            }
                        </span>

                        <span className={classnames(styles['icon'], styles['drop-icon'])} >
                            <UIIcon
                                code={'\uf04c'}
                                textColor={textColor}
                                textSize={textSize}
                                fallback={iconFallback}
                            />
                        </span >
                    </span>
                }
                watch={true}
                freezable={false}
                triggerEvent={'mouseover'}
                closeWhenMouseLeave={true}
                targetOrigin={[align, 0]}
                anchorOrigin={[align, 25]}
                onRequestCloseWhenClick={(close) => close()}
                {...otherProps}
            >
                {
                    options.map((option, index) => {
                        return (
                            <PopMenuItem
                                key={index}
                                label={this.props.labelFormatter(option)}
                                icon={isEqual(option, this.state.option) ? '\uf068' : '\u0000'}
                                onClick={this.handleClick.bind(this, option)}
                            />
                        )
                    })
                }

            </PopMenu>
        )
    }
}