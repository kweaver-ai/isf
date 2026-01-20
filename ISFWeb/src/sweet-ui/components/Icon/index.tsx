import React from 'react';
import styles from './styles';
import View from '../View';

interface IconProps {
    /**
     * 字体图标字族
     */
    fontFamily?: string;

    /**
     * 字体图标Code
     */
    code?: string;

    /**
     * 图标大小
     */
    size?: number;

    /**
     * 字体图标颜色
     */
    color?: string;

    /**
     * 图片URL或Base64
     */
    src?: string;
}

const Icon: React.FunctionComponent<IconProps> = function Icon({ size = 14, code, fontFamily, color, src, ...otherProps }) {
    return code ?
        (
            <View className={styles['fonticon']} style={{ fontSize: size, fontFamily, color }} inline={true} {...otherProps}>
                {code}
            </View>
        ) : (
            <View inline={true} {...otherProps}>
                <img width={size} height={size} src={src} />
            </View>
        );
};

export default Icon;
