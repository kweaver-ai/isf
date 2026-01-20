import React from 'react';
import classnames from 'classnames';
import LazyLoaderBase from './ui.base';
import styles from './styles.desktop';

export default class LazyLoader extends LazyLoaderBase {
    render() {
        return (
            <div
                role={this.props.role}
                className={classnames(styles['container'], this.props.className)}
                onScroll={this.handleScroll.bind(this)}
                ref={(scrollView) => this.scrollView = scrollView}
            >
                {
                    this.props.children
                }
            </div>
        )
    }
}