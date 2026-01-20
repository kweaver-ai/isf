import React from 'react';
import classnames from 'classnames';
import { ClassName } from '@/ui/helper';
import BaseButton from '../BaseButton';
import View from '../View';
import styles from './styles';

type DialogButton = {
    icon: React.ReactNode;

    onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
};

interface DialogProps {
    /**
     * 对话框标题
     */
    title: React.ReactNode;

    /**
     * 对话框宽度
     */
    width?: number | string;

    /**
     * 对话框按钮
     */
    buttons: ReadonlyArray<DialogButton>;
}

const Dialog: React.SFC<DialogProps> = function Dialog({ title, width, buttons = [], children }) {
    return (
        <View
            className={styles['dialog']}
            style={{ width }}
        >
            <View className={classnames(styles['header'], ClassName.BorderTopColor)}>
                <View className={styles['title']}>{title}</View>
                <View className={styles['buttons']}>
                    {
                        buttons.map(({ icon, onClick }, index) => (
                            <BaseButton
                                key={index}
                                onClick={onClick}
                                className={styles['button']}
                            >
                                {icon}
                            </BaseButton>
                        ))
                    }
                </View>
            </View>
            <View className={styles['content']}>{children}</View>
        </View>
    );
};

export default Dialog;
