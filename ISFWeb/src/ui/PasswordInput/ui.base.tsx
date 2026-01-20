import React from 'react';
import { noop } from 'lodash';
import { isBrowser, Browser } from '@/util/browser';

// 判断是否为Chromei浏览器，是时，添加解决Chrome浏览器下弹出“是否保存密码框”
const isChrome = isBrowser({ app: Browser.Chrome });

export default class PasswordInputBase extends React.PureComponent<UI.PasswordInput.Props, UI.PasswordInput.State> {

    static defaultProps = {

        validator: () => true,

        onChange: noop,

        onFocus: noop,

        onBlur: noop,

        onMouseout: noop,

        onMouseover: noop,

        onKeyDown: noop,
    }

    state = {
        value: this.props.value,

        focus: false,
    }

    input: HTMLInputElement;

    componentDidMount() {
        // 避免value有值时，再次渲染出来组件时，密码明文显示
        this.setType(this.state.value);
    }

    componentDidUpdate(prevProps, prevState) {
        const { value } = this.props;
        if(prevProps.value !== value) {
            this.setState({
                value,
            })
            this.setType(value);
        }
    }

    /**
     * props的值更新
     * @param value
     */
    private updateValue(value) {
        this.setState({
            value,
        })
    }

    /**
     * 改变输入框type
     * @param value 文本值
     */
    protected setType(value) {
        if (value.length === 0) {
            // 避免输入密码撤销为空时,出现密码框填充选项
            this.input.type = 'text';
        }
        else if (!isChrome) {
            // 兼容其他非chrome浏览器,避免显示输入密码
            this.input.type = 'password';
        }
    }

    /**
     * 处理onChange事件
     * @param event
     */
    protected changeHandler(event) {
        const value = event.target.value;

        if ((!this.props.required && value === '') || (this.props.validator && this.props.validator(value))) {
            this.setType(value);
            this.updateValue(value);
            this.props.onChange && this.props.onChange(value);
        } else {
            event.preventDefault();
        }
    }

    /**
     * 处理onFocus事件
     */
    protected focusHandler() {
        this.setType(this.state.value);
        this.setState({ focus: true })
        this.props.onFocus && this.props.onFocus();
    }

    /**
     * 处理onBlur事件
     */
    protected blurHandler() {
        this.setState({ focus: false })
        this.props.onBlur && this.props.onBlur();
    }

    /**
     * 处理onClick事件
     */
    protected clickHandler() {
        if (!this.props.disabled) {
            this.setState({ focus: true })
            this.input.focus()
        }
    }

    /**
     * 处理onMouseover事件
     * @param event
     */
    protected mouseoverHandler(event) {
        if (!this.props.disabled) {
            this.props.onMouseover && this.props.onMouseover(event);
        }
    }

    /**
     * 处理onMouseout事件
     * @param event
     */
    protected mouseoutHandler(event) {
        if (!this.props.disabled) {
            this.props.onMouseout && this.props.onMouseout(event);
        }
    }

    /**
     * 处理onKeyDown事件
     * @param event
     */
    protected keyDownHandler(event) {
        if (!this.props.disabled) {
            this.props.onKeyDown && this.props.onKeyDown(event);
        }
    }

    /**
    * 输入框聚焦
    */
    public focus() {
        this.input.focus();
    }

    /**
     * 输入框失焦
     */
    public blur() {
        this.input.blur();
    }
}