import React from 'react';
import DataGridFieldBase from './ui.base';

export default class DataGridField extends DataGridFieldBase {
    render() {
        return (
            <col
                width={this.props.width}
            />
        )
    }
}