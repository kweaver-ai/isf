import React from 'react';
import classnames from 'classnames';
import Icon from '../Icon/ui.desktop';
import LinkIconBase from './ui.base';
import styles from './styles.desktop';

export default class LinkIcon extends LinkIconBase {
    render() {
        return (
            <span
                href="#"
                className={classnames(styles.linkIcon, this.props.className, { [styles['disabled']]: this.props.disabled })}
                onClick={this.clickHandler.bind(this)}
                style={{ width: this.props.size, height: this.props.size }}
            >
                <Icon
                    url={this.props.url}
                    size={this.props.size}
                />
            </span>
        )
    }
}