import React from 'react';
import { noop } from 'lodash';
import UIIcon from '../UIIcon/ui.desktop'
import Panel from '../Panel/ui.desktop'
import Dialog from '../Dialog2/ui.desktop';
import __ from './locale';
import styles from './styles.desktop'

const MessageDialog: React.FunctionComponent<UI.MessageDialog.Props> = function MessageDialog({
    role,
    onConfirm = noop,
    children,
    hide,
    headless,
    className,
}) {
    return (
        <Dialog
            role={role}
            width={400}
            title={__('提示')}
            onClose={onConfirm}
            {...{ hide, headless, className }}
        >
            <Panel>
                <Panel.Main>
                    <div className={styles['main']}>
                        <div className={styles['icon']}>
                            <UIIcon
                                code={'\uf076'}
                                color={'#5a8cb4'}
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

export default MessageDialog