import React from 'react';
import classnames from 'classnames';
import { range } from '@/util/validators';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import Trigger from '../Trigger';
import Palette from './Palette';
import styles from './styles';
import disabled from './assets/palette-disabled.png';

/**
 * 色值的16进制表示，包含#
 */
export type value = string;

/**
 * R、G、B单项输入值
 */
type RGBInput = number | null;

/**
 * 色值的RGB表示
 */
export type RGB = { R: RGBInput; G: RGBInput; B: RGBInput };

export interface ColorPickerProps {
    /**
     * 16进制色值，支持以'#'开头
     */
    value?: string;

    /**
     * 取色板禁用状态
     */
    disabled?: boolean;

    /**
     * 色板
     */
    palette?: ReadonlyArray<value>;

    element?: Element | null;

    /**
     * 收起颜色选择面板且选中颜色改变时触发
     */
    onValueChange: (event: SweetUIEvent<value>) => void;
}

interface ColorPickerState {
    /**
     * 颜色选择器是否处于选中状态
     */
    active: boolean;

    /**
     * HEX颜色值
     */
    value: value;

    /**
     * HEX颜色文本框输入值
     */
    hexInputValue: string;

    /**
     * R/G/B输入值
     */
    rgb: RGB;

    /**
     * 保存props value值，用于比较props.value是否发生变化
     */
    propsValue: string;
}

const isChildNode = (el, target) => {
    if (target !== null) {
        return el === target || isChildNode(el, target.parentNode);
    }

    return false;
};

const DefaultColor = '#126EE3'

export default class ColorPicker extends React.Component<ColorPickerProps, ColorPickerState> {
    static defaultProps = {
        value: DefaultColor,
        disabled: false,
    };

    state: ColorPickerState;

    constructor(props: ColorPickerProps & { value: string }, ...args: Array<any>) {
        super(props, ...args);
        const validColor = /^#[0-9A-Fa-f]{6}$/.test(props.value)
            ? props.value.slice(1)
            : /^[0-9A-Fa-f]{6}$/.test(props.value) ? props.value : DefaultColor.slice(1);
        this.state = {
            active: false,
            value: validColor.toLocaleUpperCase(),
            hexInputValue: validColor.toLocaleUpperCase(),
            rgb: {
                R: parseInt(validColor, 16) >> 16,
                G: (parseInt(validColor, 16) >> 8) & 0xff,
                B: parseInt(validColor, 16) & 0xff,
            },
            propsValue: props.value,
        };
    }

    static getDerivedStateFromProps({ value }: ColorPickerProps, prevState: ColorPickerState) {
        if (value && value !== prevState.propsValue) {
            const validColor = /^#[0-9A-Fa-f]{6}$/.test(value)
                ? value.slice(1)
                : /^[0-9A-Fa-f]{6}$/.test(value) ? value : DefaultColor.slice(1);

            return {
                value: validColor.toLocaleUpperCase(),
                hexInputValue: validColor.toLocaleUpperCase(),
                rgb: {
                    R: parseInt(validColor, 16) >> 16,
                    G: (parseInt(validColor, 16) >> 8) & 0xff,
                    B: parseInt(validColor, 16) & 0xff,
                },
                propsValue: value,
            }
        }

        return null
    }

    /**
     * 16进制颜色值输入内容改变时触发
     * @param hexInput HEX文本框输入值
     */
    private handleHexColorChange: (input: string) => void = (hexInput: string) => {
        const hexValue = /^#/.test(hexInput) ? hexInput.slice(1) : hexInput;
        const value = /^[0-9A-Fa-f]{6}$/.test(hexValue) ? hexValue : this.state.value;

        this.setState({
            value,
            hexInputValue: hexValue,
            rgb: {
                R: parseInt(value, 16) >> 16,
                G: (parseInt(value, 16) >> 8) & 0xff,
                B: parseInt(value, 16) & 0xff,
            },
        });
    };

    /**
     * R/G/B颜色输入值改变时触发
     * @param rgb Partial<RGB>
     */
    private handleRGBChange: (rgb: Partial<RGB>) => void = ({
        R = this.state.rgb.R,
        G = this.state.rgb.G,
        B = this.state.rgb.B,
    }: Partial<RGB>) => {
        this.setState({
            rgb: { R, G, B },
        });

        if (range(0, 255)(R) && range(0, 255)(G) && range(0, 255)(B)) {
            this.setState({
                value: this.rgbToHex({ R, G, B }),
                hexInputValue: this.rgbToHex({ R, G, B }),
            });
        }
    };

    /**
     * 16进制色值输入框失焦时触发
     */
    private handleHexBlur: () => void = () => {
        this.setState({
            value: this.state.value.toLocaleUpperCase(),
            hexInputValue: this.state.value.toLocaleUpperCase(),
        });
    };

    /**
     * R / G / B 数字输入框失焦时触发
     */
    private handleNumberBoxBlur: () => void = () => {
        const { value, rgb: { R, G, B } } = this.state;

        const rgb = {
            R: (typeof R === 'number' && R > 255 || !R) ? parseInt(value, 16) >> 16 : R,
            G: (typeof G === 'number' && G > 255 || !G) ? (parseInt(value, 16) >> 8) & 0xff : G,
            B: (typeof B === 'number' && B > 255 || !B) ? parseInt(value, 16) & 0xff : B,
        };

        this.setState({
            rgb,
        });
    };

    /**
     * rgb转化为16进制色值(大写字母)
     * @param {R: number; G: number; B: number;}
     */
    private rgbToHex: (value: { R: number; G: number; B: number }) => string = ({
        R,
        G,
        B,
    }: {
            R: number;
            G: number;
            B: number;
        }) => {
        return ((1 << 24) + (R << 16) + (G << 8) + B).toString(16).slice(1).toLocaleUpperCase();
    };

    /**
     * 验证16进制颜色值输入
     * @param value HEX颜色输入值
     * 验证规则：最多输入6位，仅支持数字0-9，A-F(不区分大小写)
     */
    private validateColor: (value: string) => boolean = (value: string) => {
        return /^[0-9A-Fa-f]{1,6}$/.test(value);
    };

    /**
     * 点击面板以外的区域触发
     */
    private handleBeforePanelClose: () => void = () => {
        this.setState({
            rgb: {
                R: parseInt(this.state.value, 16) >> 16,
                G: (parseInt(this.state.value, 16) >> 8) & 0xff,
                B: parseInt(this.state.value, 16) & 0xff,
            },
            hexInputValue: this.state.value,
        });
        if (this.state.value !== this.props.value) {
            this.dispatchColorChangeEvent(`#${this.state.value}`);
        }
    };

    /**
     * 向上层抛出`onValueChange`事件
     */
    private dispatchColorChangeEvent = createEventDispatcher(this.props.onValueChange);

    /**
     * 点击调色板里的色块时触发
     * @param value 色块对应的HEX色值
     * @param close 关闭面板的方法
     */
    private handleClickPalette: (value: string, close: () => void) => void = (value: string, close: () => void) => {
        const hexValue = /^#/.test(value) ? value.slice(1) : value;
        if (/^[0-9A-Fa-f]{6}$/.test(hexValue)) {
            this.setState(
                {
                    value: hexValue,
                    hexInputValue: hexValue,
                    rgb: {
                        R: parseInt(hexValue, 16) >> 16,
                        G: (parseInt(hexValue, 16) >> 8) & 0xff,
                        B: parseInt(hexValue, 16) & 0xff,
                    },
                },
                () => {
                    if (this.state.value !== this.props.value) {
                        this.dispatchColorChangeEvent(`#${this.state.value}`);
                    }
                    close();
                },
            );
        }
    };

    render() {
        const { active, value, rgb, hexInputValue, role } = this.state;

        return (
            <Trigger
                role={role}
                renderer={({ setPopupVisibleOnClick }) => (
                    <View
                        key={'color-picker-view'}
                        onClick={this.props.disabled ? undefined : setPopupVisibleOnClick}
                        inline={true}
                        style={{ backgroundColor: this.props.disabled ? '#FFFFFF' : `#${value}` }}
                        className={classnames(
                            { [styles['color-picker']]: !this.props.disabled },
                            { [styles['active']]: active },
                            { [styles['disabled']]: this.props.disabled },
                        )}
                    >
                        {this.props.disabled ? <img width={56} height={24} draggable={false} src={disabled} /> : null}
                    </View>
                )}
                freeze={true}
                onBeforePopupClose={this.handleBeforePanelClose}
                onPopupVisibleChange={({ detail }) => {
                    this.setState({
                        active: detail,
                    });
                }}
                element={this.props.element}
            >
                {({ close }) => (
                    <Palette
                        {...{ hexInputValue, rgb, close, value }}
                        palette={this.props.palette}
                        onHexColorChange={this.handleHexColorChange}
                        onPaletteClick={(value) => this.handleClickPalette(value, close)}
                        onHexBlur={this.handleHexBlur}
                        onNumberBoxBlur={this.handleNumberBoxBlur}
                        onRGBChange={this.handleRGBChange}
                        validateColor={this.validateColor}
                    />
                )}
            </Trigger>
        );
    }
}