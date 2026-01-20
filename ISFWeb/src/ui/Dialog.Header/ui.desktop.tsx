import React from 'react';
import classnames from 'classnames';
import { isString, noop } from 'lodash';
import Text from '../Text/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import FlexBox from '../FlexBox/ui.desktop';
import { ClassName } from '../helper';
import styles from './styles.desktop';
import close from './assets/close.desktop.png';

export default function DialogHeader({ children, closable = false, onClose = noop, HeaderButtons = [] }: UI.DialogHeader.Props) {
    return (
        <div className={styles['container']}>
            <div className={classnames(styles['padding'], ClassName.BorderTopColor)}>
                <FlexBox>
                    <FlexBox.Item align="left middle">
                        {
                            isString(children) ?
                                <Text className={styles['title']}>
                                    {
                                        children
                                    }
                                </Text> :
                                children
                        }
                    </FlexBox.Item>
                    <FlexBox.Item align="right middle">
                        {
                            HeaderButtons.map((HeaderButton) => (
                                React.cloneElement(HeaderButton, { className: styles['button-icon'], size: 13 })
                            ))
                        }
                        {
                            closable ?
                                <UIIcon
                                    code={'\uf014'}
                                    size={13}
                                    fallback={close}
                                    className={styles['button-icon']}
                                    onClick={onClose}
                                /> :
                                null
                        }
                    </FlexBox.Item>
                </FlexBox>
            </div>
        </div>
    )
}