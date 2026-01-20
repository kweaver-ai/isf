import React from 'react'
import { noop } from 'lodash'

export default class Fold extends React.Component<UI.Fold.Props, UI.Fold.State> {

    static defaultProps = {
        labelProps: {
        },
        open: true,
        onToggle: noop,
    }

    state = {
        open: this.props.open,
    }

    static getDerivedStateFromProps({ open }, prevState) {
        if(open !== prevState.open ) {
            return {
                open,
            }
        }
        return null
    }

    toggle() {
        this.setState({
            open: !this.state.open,
        })
        this.props.onToggle(!this.state.open)
    }
}