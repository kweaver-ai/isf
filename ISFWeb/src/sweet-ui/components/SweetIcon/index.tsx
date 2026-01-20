import React from 'react';
import '@/core/fonts/sweetui/font.css';
import Icon from '../Icon';
import { iconSet } from './helper';

interface SweetIconProps {
    /**
     * 图标名字
     */
    name: string;

    /**
     * 图标尺寸
     */
    size?: number;

    /**
     * 图标颜色
     */
    color?: string;
}

/**
 * 根据环境支持程度获取图标code或src
 * @param name 图标名
 */
const getIconSource = (name: string) => {
    const code = iconSet[name];

    return { code };
};

/**
 * 图标
 */
const SweetIcon: React.FC<SweetIconProps> = function SweetIcon({ role, name, size = 16, color, ...otherProps }) {
    return (
        <Icon
            role={role}
            size={size}
            color={color}
            fontFamily={'SweetUI'}
            {...getIconSource(name)}
            {...otherProps}
        />
    )
};

export default SweetIcon;
