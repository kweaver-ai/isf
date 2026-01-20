import * as React from 'react';
import ToastProvider from '@/ui/ToastProvider/ui.desktop';
import SetRoleComponent from '../SetRole/SetRole/component.view';
import SetRoleBase from './component.base';

export default class SetRole extends SetRoleBase {
    render() {
        return (
            <ToastProvider role={'ui-toastprovider'}>
                <SetRoleComponent
                    users={this.props.users}
                    dep={this.props.dep}
                    userid={this.props.userid}
                    onComplete={this.props.onComplete}
                />
            </ToastProvider>
        )
    }
}