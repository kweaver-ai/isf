import React from 'react';
import { noop } from 'lodash';
import Panel from '../Panel/ui.desktop'
import Dialog from '../Dialog2/ui.desktop';
import __ from './locale';
import styles from './styles.desktop'

const SimpleDialog: React.FunctionComponent<UI.SimpleDialog.Props> = function SimpleDialog({
    onConfirm = noop,
    onClose,
    children,
    title,
    footless = false,
    hide = false,
    headless = false,
    className = '',
}) {
    return (
        <Dialog
            width={400}
            title={title ? title : __('提示')}
            onClose={onClose || onConfirm}
            {...{ hide, headless, className }}
        >
            <Panel>
                <Panel.Main>
                    <div className={styles['main']}>
                        {
                            children
                        }
                    </div>
                </Panel.Main>
                {
                    !footless ?
                        <Panel.Footer>
                            <Panel.Button
                                type="submit"
                                onClick={onConfirm}
                            >
                                {__('确定')}
                            </Panel.Button>
                        </Panel.Footer>
                        : null
                }
            </Panel>
        </Dialog>
    )
}

export default SimpleDialog