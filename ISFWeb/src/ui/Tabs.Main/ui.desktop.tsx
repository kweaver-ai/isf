import React from 'react';
import classnames from 'classnames';
import TabsMainBase from './ui.base';
import styles from './styles.desktop';

export default class TabsMain extends TabsMainBase {
    render() {
        return (
            <div role={this.props.role} className={classnames(styles['main'], this.props.className)}>
                {
                    React.Children.map(this.props.children, (Content, i) => this.props.activeIndex === i ? Content : null)
                }
            </div>
        )
    }
}