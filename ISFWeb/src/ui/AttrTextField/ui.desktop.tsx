import React from 'react';
import __ from './locale';
import AttrTextFieldBase from './ui.base';
import styles from './styles.desktop';
import { ExpandStatus } from './helper';

export default class AttrTextField extends AttrTextFieldBase {

    render() {
        return (
            <div className={styles['meta']}>
                <label className={styles['name']} title={this.attribute.name}>{this.attribute.formatName}：</label>
                {
                    this.valueField()
                }
            </div>
        )
    }

    valueField() {
        if (this.state.status === ExpandStatus.EXPAND) {
            return (
                <span className={styles['val']}>
                    <p>{this.attribute.formatValue}</p>
                    <span className={styles['elastic']} href="#" onClick={(e) => { e.stopPropagation();  this.setState({ status: ExpandStatus.COLLAPSE }) } }>{__('展开')}▼</span>
                </span>
            )
        } else if (this.state.status === ExpandStatus.COLLAPSE) {
            return (
                <span className={styles['val']}>
                    <p>{this.attribute.value}</p>
                    <span className={styles['elastic']} href="#" onClick={(e) => { e.stopPropagation(); this.setState({ status: ExpandStatus.EXPAND }) } }>{__('收起')}▲</span>
                </span>
            )
        } else {
            return (
                <span className={styles['val']}>
                    <p title={this.attribute.value}>{this.attribute.formatValue}</p>
                </span>
            )

        }
    }
}