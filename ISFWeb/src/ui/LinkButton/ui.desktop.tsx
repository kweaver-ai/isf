import React from 'react';
import classnames from 'classnames'
import LinkButtonBase from './ui.base';
import styles from './styles.desktop';

export default class LinkButton extends LinkButtonBase {
    render() {
        return (
            <span
                ref="linkButton"
                href="#"
                draggable={false}
                className={classnames(styles['link-button'], { [styles['disabled']]: this.props.disabled }, this.props.className)}
                onClick={this.handleClick.bind(this)}
            >
                {
                    this.props.children
                }
            </span>
        )
    }
}