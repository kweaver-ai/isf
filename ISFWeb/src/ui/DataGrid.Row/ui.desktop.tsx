import React from 'react';
import classnames from 'classnames';
import CheckBox from '../CheckBox/ui.desktop';
import DataGridRowBase from './ui.base';
import styles from './styles.desktop';

export default class DataGridRow extends DataGridRowBase {
    render() {
        return (
            <tr
                className={ classnames(styles['record'], { [styles['selected']]: this.props.selected, [styles['strap']]: this.props.strap }) }
                onClick={ this.props.onClick }
                onDoubleClick={ this.props.onDoubleClick }
            >
                {
                    this.props.checkbox ?
                        <td className={ styles['checkbox-cell'] }>
                            <CheckBox
                                disabled={ this.props.disabled === true }
                                checked={ this.props.selected === true }
                                onClick={ this.onCheckBoxClick.bind(this) } />
                        </td> :
                        null
                }
                {
                    this.props.children
                }
            </tr>
        )
    }
}