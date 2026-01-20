
import { Component } from 'react';

export default class Stack extends Component<any, any> {
    static defaultProps = {
        background: '#EFEFEF',
        rate: 0,
    }

    componentDidMount() {
        const { stack } = this.refs;
        if (stack && this.props.rate > 0 && stack.offsetWidth < 5 && stack.offsetWidth >= 0) {
            stack.style.width = '5px';
        }
    }

    componentDidUpdate() {
        const { stack } = this.refs;
        if (stack && this.props.rate > 0 && stack.offsetWidth < 5 && stack.offsetWidth >= 0) {
            stack.style.width = '5px';
        }
    }
}