import React from 'react';
import classnames from 'classnames';
import SwitchButtonBase from './ui.base';
import styles from './styles.desktop';

export default class SwitchButton extends SwitchButtonBase {
    render() {
        return (
            <div className={classnames(styles['container'], { [styles['switch-left']]: this.state.active===false }, { [styles['switch-right']]: this.state.active ===true })}
                onClick={() => this.toggleStatus(this.props.value, !this.state.active)} >
                <div className={styles['child']}></div>
            </div>
        )
    }
}