import React from 'react';

export default class StarBase extends React.Component<any, any> {

    static defaultProps = {
        score: 0,
    }

    state = {
        score: this.props.score,
    }

    static getDerivedStateFromProps({ score }, prevState) {
        if (score !== prevState.score) {
            return {
                score,
            }
        }
        return null
    }

    handleMouseEnter(score) {
        this.setState({
            score,
        })
    }

    handleMouseLeave() {
        this.setState({
            score: this.props.score,
        })
    }
}