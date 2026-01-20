import React from 'react';
import Dialog from '../Dialog2/ui.desktop';
import ProgressDialogBase from './ui.base';
import ProgressDialogView from './ui.view';

export default class ProgressDialog extends ProgressDialogBase {
    render() {
        const { progress, item } = this.state;
        return (
            <Dialog
                title={this.props.title}
                width={440}
                onClose={this.handleCancel.bind(this)}
            >
                <ProgressDialogView
                    detailTemplate={this.props.detailTemplate}
                    item={item}
                    progress={progress}
                    prohandleCancel={this.handleCancel.bind(this)}
                />
            </Dialog>
        )
    }
}