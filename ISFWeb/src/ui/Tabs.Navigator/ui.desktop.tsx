import React from 'react';
import classnames from 'classnames';
import TabsNavigatorBase from './ui.base';
import styles from './styles.desktop';

export default class TabsNavigator extends TabsNavigatorBase {
    render() {
        return (
            <div
                role={this.props.role}
                className={classnames(styles['navigator'], this.props.className)}
            >
                {
                    React.Children.map(this.props.children, (Tab, index) => {
                        return React.cloneElement(Tab, {
                            active: this.state.activeIndex === index,
                            onActive: this.navigate.bind(this, index),
                        })
                    })
                }
            </div>
        )
    }
}