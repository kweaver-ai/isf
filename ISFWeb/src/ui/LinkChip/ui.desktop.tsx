import React from 'react';
import classnames from 'classnames';
import LinkChipBase from './ui.base';
import styles from './styles.desktop';

export default class LinkChip extends LinkChipBase {
    render() {
        return (
            <span
                href="#"
                className={classnames(styles.linkChip, this.props.className)}
                onClick={this.clickHandler.bind(this)}
                disabled={this.props.disabled}
                title={this.props.title}
                onDoubleClick={this.doubleClickHandler.bind(this)}
            >
                {
                    this.props.children
                }
            </span>
        )
    }
}