import React from 'react';
import { createRef } from 'react';
import { assign } from 'lodash';

export default class OverlayBase extends React.Component<any, any> {
    static defaultProps = {
        position: '',
    }

    state = {
        align: {},
        position: this.props.position,
    }

    overlayRef = createRef()

    componentDidMount() {
        if (this.props.position) {
            this.setState({
                align: this.align(this.state.position),
            });
        }
    }

    static getDerivedStateFromProps({ position }, prevState) {
        if (position && position !== prevState.position) {
            return {
                position,
            }
        }
        return null
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.state.position !== prevState.position) {
            this.setState({
                align: this.align(this.state.position),
            })
        }
    }

    align(position) {
        return position.split(/\s+/).reduce((align, key, i) => {
            switch (key) {
                case 'top':
                    return assign(align, { top: 0 });

                case 'right':
                    return assign(align, { right: 0 });

                case 'bottom':
                    return assign(align, { bottom: 0 });

                case 'left':
                    return assign(align, { left: 0 });

                case 'middle':
                    return assign(align, { top: (document.documentElement.clientHeight - this.overlayRef.current.clientHeight) / 2 });

                case 'center':
                    return assign(align, { left: (document.documentElement.clientWidth - this.overlayRef.current.clientWidth) / 2 });

                default:
                    if (i === 0) {
                        return assign(align, { left: key });
                    } else if (i === 1) {
                        return assign(align, { top: key });
                    }
            }
        }, {})
    }
}