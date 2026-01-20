import React from 'react';

export default class DataGridRowBase extends React.Component<UI.DataGridRow.Props, any> {
    static defaultProps = {
        selected: false,
    }

    protected onCheckBoxClick(event) {
        this.props.onCheckChange({ multi: true });

        event.stopPropagation();
    }
}