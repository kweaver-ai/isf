import React from 'react';
import classnames from 'classnames';
import LocatorBase from './ui.base';
import styles from './styles.desktop';

export default class Locator extends LocatorBase {
    render() {
        return (
            <div
                className={classnames(styles['locator'], this.props.className)}
                ref={(el) => this.el = el}
            >
                {
                    this.props.children
                }
            </div>
        )
    }
}