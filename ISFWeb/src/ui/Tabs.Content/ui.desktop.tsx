import React from 'react';
import classnames from 'classnames';
import TabsContentBase from './ui.base';
import styles from './styles.desktop';

export default class TabsContent extends TabsContentBase {
    render() {
        return (
            <div role={this.props.role} className={classnames(styles['content'], this.props.className)}>
                {
                    this.props.children
                }
            </div>
        )
    }
}