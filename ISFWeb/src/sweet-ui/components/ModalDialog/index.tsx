import React from 'react';
import Dialog from '../Dialog';
import Modal from '../Modal';
import Portal from '../Portal';
import Locator from '../Locator';
import SweetIcon from '../SweetIcon';

interface ModalDialogProps {
    /**
     * 把内容挂载到目标节点
     */
    target: HTMLElement;

    /**
     * 对话框标题
     */
    title: string;

    /**
     * 对话框宽度
     */
    width?: number | string;

    /**
     * 对话框是否支持拖拽
     */
    draggable?: boolean;

    /**
     * 请求关闭窗口
     */
    onRequestClose: () => void;
}

const ModalDialog: React.SFC<ModalDialogProps> = function ModalDialog({ target, children, title, width, draggable, onRequestClose }) {
    return (
        <Portal getContainer={() => target}>
            <Modal>
                <Locator draggable={draggable} anchorOrigin={['center', 'center']} alignOrigin={['center', 'center']}>
                    <Dialog
                        width={width}
                        buttons={[
                            {
                                icon: <SweetIcon name="x" size={13} color="#bdbdbd" />,
                                onClick: (event) => onRequestClose(),
                            },
                        ]}
                        {...{ title }}
                    >
                        {children}
                    </Dialog>
                </Locator>
            </Modal>
        </Portal>
    );
};

export default ModalDialog;
