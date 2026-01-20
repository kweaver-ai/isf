import { Component } from 'react';

export default class DataGridFieldBase extends Component<UI.DataGrid.Field, any> {
    static defaultProps = {
        width: 100,

        field: '',

        label: '',

        // 默认不传是居左显示
        align: 'left',

        className: '',

        formatter: (val, record) => val,
    }
}