import React from 'react';
import classnames from 'classnames';
import Button, { ButtonProps } from '../Button/Button';
import SweetIcon from '../SweetIcon';
import View from '../View';
import styles from './styles';

/**
 * 传递到Button上的props
 */
type ButtonReceivedProps = 'disabled' | 'style' | 'className' | 'onClick' | 'theme' | 'size' | 'width';

interface DropButtonProps extends Pick<ButtonProps, ButtonReceivedProps> {
    /**
     * 有无背景（下拉文字按钮）
     */
    background?: boolean;
}

const DropButton: React.FunctionComponent<DropButtonProps> = function DropButton({
    disabled,
    theme,
    background = true,
    children,
    ...restProps
}) {
    return (
        <Button
            className={classnames(
                { [styles['text-button']]: theme === 'text' },
                { [styles['no-background']]: !background },
                { [styles['disabled']]: disabled },
            )}
            {...{ ...restProps, disabled, theme }}
        >
            <View inline={true} className={classnames(styles['text'], { [styles['disabled']]: disabled })}>
                {children}
            </View>
            <SweetIcon
                size={16}
                name={'arrowDown'}
                color={theme === 'gray' ? '#fff' : '#000'}
                className={classnames(
                    { [styles['drop-icon']]: theme !== 'gray' },
                    { [styles['disabled']]: disabled && theme !== 'gray' },
                )}
            />
        </Button>
    );
};

export default DropButton;
