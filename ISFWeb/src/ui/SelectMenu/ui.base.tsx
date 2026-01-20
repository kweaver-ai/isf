import React from 'react'

export default class SelectMenu extends React.Component<UI.SelectMenu.Props> {

    static defaultProps: UI.SelectMenu.Props = {
        triggerEvent: 'mouseover',
        freezable: false,
        anchorOrigin: ['left', 'bottom'],
        targetOrigin: ['left', 'top'],
        onRequestCloseWhenClick: (close) => close(),
    }

    state = {
        value: 'value' in this.props ? this.props.value : this.props.defaultValue,
    }

    static getDerivedStateFromProps({ value }, prevState) {
        if (typeof value !== 'undefined' && value !== prevState.value) {
            return {
                value,
            }
        }
        return null;
    }

    handleClick(e, item) {
        if (!item.props.disabled) {
            if (typeof item.props.onClick === 'function') {
                item.props.onClick(e)
            }

            if (!e.defaultPrevented) {
                this.setState({
                    value: item.props.value,
                })

                if (typeof this.props.onChange === 'function') {
                    this.props.onChange(item.props.value)
                }
            }
        }
    }
}