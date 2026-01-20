import React from 'react';
import classnames from 'classnames';
import { isFunction } from 'lodash'
import { SweetUIEvent } from '../../utils/event';
import TextInput from '../TextInput';
import SweetIcon from '../SweetIcon';
import View from '../View';
import styles from './styles';

interface TextBoxProps {
    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * width，包含盒模型的padding和border
     */
    width?: string;

    /**
     * 文本框默认内容(建议非受控场景下使用)
     */
    defaultValue?: string;

    /**
     * 文本框内容
     */
    value: string;

    /**
     * className，控制外层view样式
     */
    className?: string;

    /**
     * 自动聚焦
     */
    autoFocus?: boolean;

    /**
     * 聚焦时选中
     */
    selectOnFocus?: [number] | [number, number] | boolean;

    /**
     * 文本框状态
     */
    status?: 'normal' | 'error';

    /**
     * css样式，传入TextInput，控制TextInput样式
     */
    style?: React.CSSProperties;

    /**
     * 输入内容发生变化时触发
     */
    onValueChange?: (event: SweetUIEvent<string>) => void;

    /**
     * 按下enter键事触发
     */
    onPressEnter?: (event: React.KeyboardEvent<HTMLInputElement>) => void;

    /**
     * 键盘按下时触发
     */
    onKeyDown?: (event: React.KeyboardEvent<HTMLInputElement>) => void;

    /**
     * 聚焦时触发
     */
    onFocus?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 失焦时触发
     */
    onBlur?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 输入限制函数
     */
    validator?: (value: string) => boolean;

    /**
     * 鼠标进入文本域时触发（外层View）
     */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发（外层View）
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;

    role?: string;
}

interface TextBoxState {
    active: boolean;
}

export default class TextBox extends React.Component<TextBoxProps, TextBoxState> {
    static defaultProps = {
        disabled: false,
        status: 'normal',
    };

    state = {
        active: false,
    };

    input: TextInput | null = null;

    /**
     * 触发onFocus
     */
    private handleFocus = (event: React.FocusEvent<HTMLInputElement>) => {
        this.setState({ active: true });
        isFunction(this.props.onFocus) && this.props.onFocus(event)
    };

    /**
     * 触发onBlur
     */
    private handleBlur = (event: React.FocusEvent<HTMLInputElement>) => {
        this.setState({ active: false });
        isFunction(this.props.onBlur) && this.props.onBlur(event)
    };

    render() {
        const { className, style, disabled, width, onFocus, onBlur, onMouseEnter, onMouseLeave, status, role, ...otherProps } = this.props;

        return (
            <View
                role={role}
                className={classnames(styles['text-box'], className, {
                    [styles['disabled']]: disabled,
                    [styles[`${status}`]]: status,
                    [styles[`${status}-active`]]: !!this.state.active,
                    [styles[`${status}-active`]]: !!this.state.active,
                })}
                style={{ width }}
                inline={true}
                {...{ onMouseEnter, onMouseLeave }}
            >
                <TextInput
                    style={style}
                    disabled={disabled}
                    width={'100%'}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    {...otherProps}
                />
                {
                    status === 'error' ?
                        <SweetIcon
                            name={'caution'}
                            size={16}
                            color={'#e60012'}
                            className={styles['caution']}
                        /> : null
                }
            </View>
        );
    }
}
