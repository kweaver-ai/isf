import React from 'react';
import classnames from 'classnames';
import View from '../../View';
import IconButton from '../../IconButton';
import { DialogIcon, titleTextAlign } from '../index';
import styles from './styles';

interface DialogHeaderProps {
    /**
     * 对话框标题
     */
    title?: React.ReactNode;

    /**
     * 标题位置
     */
    titleTextAlign?: titleTextAlign;

    /**
     * 对话框顶部的图标按钮
     */
    icons?: ReadonlyArray<DialogIcon>;

    /**
	 * 鼠标按下时的回调
	 */
    onMouseDown?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标松开时的回调
     */
    onMouseUp?: (event: React.MouseEvent<HTMLElement>) => void;
}

const DialogHeader: React.FC<DialogHeaderProps> = function DialogHeader({ onMouseDown, onMouseUp, title, titleTextAlign = 'left', icons = [] }) {
    const higherHeader = !title || titleTextAlign === 'center'// 无标题和标题位置居中时，标题栏高度为64px

    return (
        <View
            className={classnames(
                styles['header'],
                { [styles['high-header']]: higherHeader },
            )}
            onMouseDown={onMouseDown}
            onMouseUp={onMouseUp}
        >
            {title}
            <View className={styles['icons']}>
                {
                    icons.map(({ icon, onClick }, index) => (
                        <IconButton
                            key={index}
                            icon={icon}
                            size={higherHeader ? 32 : 40}
                            onClick={onClick}
                            className={styles['icon']}
                        />
                    ))
                }
            </View>
        </View>
    );
};

export default DialogHeader;