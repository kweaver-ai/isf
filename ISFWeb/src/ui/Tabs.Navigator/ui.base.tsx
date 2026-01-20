import React from 'react';

export default class TabsNavigatorBase extends React.Component<any, any> {
    state = {
        activeIndex: 0,
    }

    componentDidMount() {
        React.Children.forEach(this.props.children, (Tab: any, index: number): void => {
            if (Tab.props.active) {
                this.navigate(index)
            }
        })
    }

    componentDidUpdate(prevProps, prevState) {
        if (!(prevProps.children && this.props.children && this.props.children.length === prevProps.children.length)) {
            this.navigate(0)
        }
    }

    public navigate(activeIndex) {
        this.setState({ activeIndex });
        this.props.onNavigate(activeIndex);
    }
}