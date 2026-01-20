import React from 'react';
import classnames from 'classnames';
import { isFunction } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import SweetIcon from '../SweetIcon';
import styles from './styles';

interface TextAreaProps {
    /**
     * 是否禁用状态
     */
    disabled?: boolean;

    /**
     * 是否控制光标位置
     */
    cursorControl?: boolean;

    /**
     * 输入内容的最大长度(正整数)
     */
    maxLength?: number;

    /**
     * 文本域是否支持编辑
     */
    readOnly?: boolean;

    /**
     * 占位提示
     */
    placeholder?: string;

    /**
     * 是否必填
     */
    required?: boolean;

    /**
     * 文本域默认内容
     */
    defaultValue?: string;

    /**
     * 文本域内容
     */
    value: string;

    /**
     * 文本域宽度
     */
    width?: number;

    /**
     * 文本域高度
     */
    height?: number;

    /**
     * 文本框padding
     */
    paddingRight?: number;

    /**
     * 文本域状态
     */
    status?: 'normal' | 'error';

    /**
     * 文本框的换行属性
     */
    wrap?: 'off' | 'hard' | 'soft' ;

    /**
     * 输入验证
     * @param value 文本值
     */
    validator?: (input: string) => boolean;

    /**
     * 文本域内容变化时的回调
     */
    onValueChange?: (e: SweetUIEvent<string>) => void;

    /**
     * 按下回车的回调
     */
    onPressEnter?: (e: React.KeyboardEvent) => void;

    /**
     * 聚焦事件回调
     */
    onFocus?: (e: React.FocusEvent<HTMLElement>) => void;

    /**
     * 失焦事件回调
     */
    onBlur?: (e: React.FocusEvent<HTMLElement>) => void;

    /**
     * 鼠标进入文本域时触发（外层View）
     */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发（外层View）
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;
}

interface TextAreaState {
    /**
     * 文本域的输入内容
     */
    value: string;

    /**
     * 文本域选中状态
     */
    active: boolean;
}

export default class TextArea extends React.Component<TextAreaProps, TextAreaState> {
    static defaultProps = {
        disabled: false,
        cursorControl: false,
        validator: () => true,
        width: 200,
        height: 96,
        status: 'normal',
    };

    constructor(props: TextAreaProps, ...args: any[]) {
        super(props);
        this.state = {
            value: (typeof props.value === 'undefined' ? props.defaultValue : props.value) || '',
            active: false,
        };
    }

    textareaRef = React.createRef();
    cursorPosition = 0;

    static getDerivedStateFromProps(nextProps: TextAreaProps, prevState: TextAreaState) {
        if ('value' in nextProps && nextProps.value !== prevState.value) {
            return {
                value: nextProps.value,
            };
        }

        return null;
    }

    handleInputValueChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        const value = event.target.value;

        if (isFunction(this.props.validator) && this.props.validator(value)) {
            this.setState(() => {
                return {
                    value,
                };
            });
            if (this.props.cursorControl) {
                this.cursorPosition = event.target.selectionStart;
            }

            this.dispatchValueChangeEvent(value);
        }
    };

    componentDidUpdate() {
        if (this.props.cursorControl) {
            this.textareaRef.current.selectionStart = this.cursorPosition;
            this.textareaRef.current.selectionEnd = this.cursorPosition;
        }
    }

    dispatchValueChangeEvent = createEventDispatcher(this.props.onValueChange);

    handleFocus = (e: React.FocusEvent<HTMLElement>) => {
        this.setState({
            active: true,
        });
        isFunction(this.props.onFocus) && this.props.onFocus(e);
    };

    handleBlur = (e: React.FocusEvent<HTMLElement>) => {
        this.setState({
            active: false,
        });
        isFunction(this.props.onBlur) && this.props.onBlur(e);
    };

    handleKeyDown = (e: React.KeyboardEvent) => { // TODO
        if (e.keyCode === 13) {
            if (this.props.cursorControl) {
                const { value } = this.state;
                if (value[value.length - 1] !== '\n') {
                    this.setState({
                        value: value + '\n',
                    }, () => {
                        this.textareaRef.current.selectionStart = this.state.value.length;
                        this.textareaRef.current.selectionEnd = this.state.value.length;
                    });
                } else {
                    e.preventDefault();
                    this.textareaRef.current.selectionStart = this.state.value.length;
                    this.textareaRef.current.selectionEnd = this.state.value.length;
                }
            } else {
                isFunction(this.props.onPressEnter) && this.props.onPressEnter(e)
            }
        }
    };

    render() {
        const { maxLength, disabled, placeholder, required, readOnly, width, height, paddingRight, status, wrap, onMouseEnter, onMouseLeave } = this.props;
        const { value, active } = this.state;

        return (
            <View
                style={{ width, height, ...(paddingRight && { paddingRight }) }}
                className={classnames(styles['text-area'], {
                    [styles['disabled']]: disabled,
                    [styles['limit-length']]: maxLength,
                    [styles[`${status}`]]: status,
                    [styles[`${status}-active`]]: !!active,
                    [styles['read-only']]: readOnly,
                })}
                {...{ onMouseEnter, onMouseLeave }}
            >
                <textarea
                    ref={this.textareaRef}
                    className={classnames(styles['text-input'], { [styles['disabled']]: disabled })}
                    {...{ value, maxLength, disabled, placeholder, required, readOnly, wrap }}
                    onChange={this.handleInputValueChange}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onKeyDown={this.handleKeyDown}
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
                {
                    maxLength ? (
                        <View
                            inline={true}
                            className={classnames(styles['text-count'], {
                                [styles['over-length']]: value.length > maxLength,
                            })}
                        >
                            {`${value.length}/${maxLength}`}
                        </View>
                    ) : null}
            </View>
        );
    }
}
