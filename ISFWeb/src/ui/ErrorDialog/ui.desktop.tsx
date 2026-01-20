import React from 'react';
import { noop } from 'lodash';
import UIIcon from '../UIIcon/ui.desktop'
import Panel from '../Panel/ui.desktop'
import Dialog from '../Dialog2/ui.desktop';
import ErrorDialogTitle from '../ErrorDialog.Title/ui.desktop'
import ErrorDialogDetail from '../ErrorDialog.Detail/ui.desktop'
import __ from './locale';
import styles from './styles.desktop'

const ErrorDialog: React.FunctionComponent<UI.ErrorDialog.Props> = function ErrorDialog({
    role,
    iconCode = '\uf075',
    iconColor = '#be0000',
    children,
    onConfirm = noop,
}) {
    return (
        <Dialog
            role={role}
            width={400}
            title={__('提示')}
            onClose={onConfirm}
        >
            <Panel>
                <Panel.Main>
                    <div className={styles['main']}>
                        <div className={styles['icon']}>
                            <UIIcon
                                code={iconCode}
                                color={iconColor}
                                size={40}
                            />
                        </div>
                        <div className={styles['message']}>
                            {
                                children
                            }
                        </div>
                    </div>
                </Panel.Main>
                <Panel.Footer>
                    <Panel.Button
                        type="submit"
                        onClick={onConfirm}
                    >
                        {__('确定')}
                    </Panel.Button>
                </Panel.Footer>
            </Panel>
        </Dialog>
    )
}

ErrorDialog.Title = ErrorDialogTitle;
ErrorDialog.Detail = ErrorDialogDetail;

export default ErrorDialog