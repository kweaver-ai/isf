import React from 'react'
export default class SelectDrop extends React.Component<UI.SelectDrop.Props> {

    state = {
        option: this.props.defaultOption,
    }

    handleClick(option) {
        this.setState({
            option: option,
        }, () => {
            if (typeof this.props.onChange === 'function') {
                this.props.onChange(option)
            }
        })
    }

}