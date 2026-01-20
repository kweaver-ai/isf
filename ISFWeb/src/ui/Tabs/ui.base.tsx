import React from 'react';
import { isFunction } from 'lodash';

export default class TabsBase extends React.Component<any, any> {
    state = {
        activeIndex: 0,
    }

    createChildren(Navigator: React.ReactElement<any>, Main: React.ReactElement<any>): Array<React.ReactElement<any>> {
        const _this = this;

        return [
            React.cloneElement(Navigator, {
                onNavigate(activeIndex) {
                    _this.setState({ activeIndex })
                    isFunction(_this.props.onNavigate) && _this.props.onNavigate(activeIndex)
                },
            }),

            React.cloneElement(Main, {
                activeIndex: this.state.activeIndex,
            }),
        ]
    }
}