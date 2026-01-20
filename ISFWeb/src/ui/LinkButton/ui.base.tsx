import React from 'react'

export default class LinkButtonBase extends React.Component<UI.LinkButton.Props, any> {
    static defaultProps = {
        disabled: false,
    }

    handleClick(e) {
        if (!this.props.disabled && typeof this.props.onClick === 'function') {
            this.props.onClick(e)
        }
    }
}