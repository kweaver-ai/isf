import React from 'react';
import classnames from 'classnames';
import Text from '@/ui/Text/ui.desktop';
import SweetIcon from '../../SweetIcon'
import styles from './styles';

export interface SelectOptionProps {
    /**
	 * className
	 */
    className?: string;

    /**
	 * 选项是否禁用
	 */
    disabled?: boolean;

    /**
	 * 根据此属性值进行筛选
	 */
    value?: any;

    /**
	 * children为显示的内容
	 */
    children: any;

    /**
	 * 选项被选中
	 */
    selected?: boolean;

    /**
	 * 标签名称，默认没有标签，string类型的标签名称匹配SweetIcon中的名称
	 */
    iconName?: string | null;

    /**
	 * 选项的点击事件
	 */
    onClick?: () => void;
}

const SelectOption: React.FunctionComponent<SelectOptionProps> = function SelectOption({
    children,
    className,
    disabled = false,
    iconName = null,
    selected = false,
    value,
    onClick,
    ...otherProps
}) {
    return (
        <li
            className={classnames(
                styles['select-option'],
                {
                    [styles['disabled']]: disabled,
                },
                {
                    [styles['select']]: selected,
                },
                {
                    [styles['icon']]: !!iconName,
                },
                className,
            )}
            onClick={onClick}
            {...otherProps}
        >
            {
                iconName && <SweetIcon
                    name={iconName}
                    size={16}
                    className={classnames(
                        styles['sweet-icon'],
                        {
                            [styles['disabled']]: disabled,
                        },
                    )}
                />
            }
            <Text>{children}</Text>
        </li>
    );
};

export default SelectOption;