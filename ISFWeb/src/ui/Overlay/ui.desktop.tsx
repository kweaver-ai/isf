import React from 'react';
import OverlayBase from './ui.base';
import classnames from 'classnames';
import styles from './styles.desktop';

export default class Overlay extends OverlayBase {
    render() {
        return (
            <div
                className={classnames(
                    styles['overlay'],
                    { [styles['fixed']]: this.state.position },
                    this.props.className,
                )}
                ref={this.overlayRef}
                style={this.state.align}
            >
                {
                    this.props.children
                }
            </div>
        )
    }
}