import React from 'react';
import ReactDOM from 'react-dom';

interface PortalProps {
    getContainer: () => HTMLElement;

    children: React.ReactNode;
}

export default class Portal extends React.Component<PortalProps, any> {

    container: HTMLElement | null = null;

    componentDidMount() {
        this.createContainer();
    }

    componentWillUnmount() {
        this.removeContainer();
    }

    createContainer() {
        this.container = this.props.getContainer();
        this.forceUpdate();
    }

    public removeContainer() {
        this.container && this.container.parentNode && this.container.parentNode.removeChild(this.container);
    }

    render() {
        if (this.container) {
            return ReactDOM.createPortal(this.props.children, this.container);
        }

        return null;
    }
}
