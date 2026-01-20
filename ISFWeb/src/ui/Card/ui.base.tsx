import React from 'react';
import { isEqual } from 'lodash';

export default class CardBase extends React.Component<UI.Card.Props, UI.Card.State> {

    static defaultProps = {
        width: '100%',
        height: '100%',
    }

    state = {
        width: this.props.width,
        height: this.props.height,
    }

    componentDidUpdate(preProps) {
        if (!isEqual(this.props, preProps)) {
            const { width, height } = this.props;
            this.setState({ width, height })
        }
    }
}