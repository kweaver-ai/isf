import React from 'react';
import classnames from 'classnames';
import View from '../../View';
import TextBox from '../../TextBox';
import NumberBox from '../../NumberBox';
import { RGB, value } from '../index';
import styles from '../styles';

interface PaletteProps {
    /**
     * HEX颜色文本框输入值
     */
    hexInputValue: string;

    /**
     * 16进制颜色文本框输入值验证函数
     */
    validateColor: (value: string) => boolean;

    /**
     * 16进制颜色值输入内容改变时触发
     */
    onHexColorChange: (value: string) => void;

    /**
     * 16进制色值输入框失焦时触发
     */
    onHexBlur: () => void;

    /**
     * R/G/B输入值
     */
    rgb: RGB;

    /**
     * R/G/B颜色输入值改变时触发
     */
    onRGBChange: (rgb: Partial<RGB>) => void;

    /**
     * R / G / B 数字输入框失焦时触发
     */
    onNumberBoxBlur: () => void;

    /**
     * 色板
     */
    palette?: ReadonlyArray<value>;

    /**
     * 点击调色板里的色块时触发
     */
    onPaletteClick: (value: string) => void;

    /**
     * 取色板当前颜色值，包含#
     */
    value: string;
}

const Palette: React.SFC<PaletteProps> = function Palette({
    hexInputValue,
    validateColor,
    onHexColorChange,
    onHexBlur,
    rgb,
    onRGBChange,
    onNumberBoxBlur,
    palette,
    onPaletteClick,
    value,
}) {
    return (
        <View className={styles['panel']}>
            <View inline={true}>
                <TextBox
                    width={74}
                    value={hexInputValue}
                    validator={validateColor}
                    onValueChange={({ detail }) => onHexColorChange(detail)}
                    autoFocus={true}
                    selectOnFocus={true}
                    onBlur={onHexBlur}
                />
                <View className={styles['text-layout']}>#</View>
            </View>
            <View inline={true} className={styles['number-box-layout']}>
                <NumberBox
                    value={rgb.R}
                    min={0}
                    max={255}
                    width={42}
                    selectOnFocus={true}
                    precision={0}
                    onValueChange={({ detail }) => onRGBChange({ R: detail })}
                    onBlur={onNumberBoxBlur}
                />
                <View className={styles['text-layout']}>R</View>
            </View>
            <View inline={true} className={styles['number-box-layout']}>
                <NumberBox
                    min={0}
                    max={255}
                    value={rgb.G}
                    width={42}
                    selectOnFocus={true}
                    precision={0}
                    onValueChange={({ detail }) => onRGBChange({ G: detail })}
                    onBlur={onNumberBoxBlur}
                />
                <View className={styles['text-layout']}>G</View>
            </View>
            <View inline={true} className={styles['number-box-layout']}>
                <NumberBox
                    min={0}
                    max={255}
                    value={rgb.B}
                    width={42}
                    selectOnFocus={true}
                    precision={0}
                    onValueChange={({ detail }) => onRGBChange({ B: detail })}
                    onBlur={onNumberBoxBlur}
                />
                <View className={styles['text-layout']}>B</View>
            </View>
            <View className={styles['dividing-line']} />
            <View className={styles['palette-area']}>
                {
                    Array.isArray(palette) &&
                    palette.map((color) => (
                        <View
                            className={classnames(styles['palette'], { [styles['selected']]: color === `#${value}` })}
                            inline={true}
                            key={color}
                            style={{ backgroundColor: color }}
                            onClick={() => onPaletteClick(color.slice(1))}
                        />
                    ))
                }
            </View>
        </View>
    );
};

Palette.defaultProps = {
    palette: [
        '#DD3333',
        '#FF4D4F',
        '#FA8C16',
        '#FFC62C',
        '#126EE3',
        '#597EF7',
        '#40A9FF',
        '#69C0FF',
        '#13C2C2',
        '#3FC380',
        '#000000',
        '#7D8791',
    ],
};

export default Palette;