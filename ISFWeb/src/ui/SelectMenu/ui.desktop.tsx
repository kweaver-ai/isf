import React from 'react'
import SelectMenuBase from './ui.base'
import PopMenu from '../PopMenu/ui.desktop'
import Option from '../SelectMenu.Option/ui.desktop'

export default class SelectMenu extends SelectMenuBase {

    static Option = Option

    render() {
        const { children, label, ...otherProps } = this.props
        return (
            <PopMenu
                trigger={label}
                {...otherProps}
            >
                {
                    React.Children.map(children, (option) => React.cloneElement(
                        option,
                        {
                            selected: option.props.value === this.state.value,
                            disabled: option.props.disabled,
                            onClick: (e) => this.handleClick(e, option),
                        },
                    ))
                }
            </PopMenu>
        )
    }
}