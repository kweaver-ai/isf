import React from 'react';
import { noop, isFunction } from 'lodash';

interface Props extends React.Props<any> {
    /**
     * 样式
     */
    className?: string;

    /**
     * 禁用
     */
    readOnly?: boolean;

    /**
     * 只读
     */
    disabled?: boolean;

    /**
     * 占位文本
     */
    placeholder?: string;

    /**
     * 键盘keyDown事件
     */
    onKeyDown?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 粘贴事件
     */
    onPaste?: (event: React.ClipboardEvent<HTMLElement>) => void;

    /**
     * 失去焦点
     */
    onBlur?: (event: React.FocusEvent<HTMLElement>) => void;

    /**
     * 输入值变化时触发
     */
    onValueChange: (value: string) => void;

}

interface State {
    value: string;

    placeholder?: string;
}

export default class FlexTextBoxBase extends React.Component<Props, State> {
    static defaultProps = {
        readOnly: false,

        disabled: false,

        onKeyDown: noop,

        onPaste: noop,

        placeholder: '',
    }

    state: State = {
        value: '',
        placeholder: this.props.placeholder,
    }

    textBox: HTMLAnchorElement | null;

    static getDerivedStateFromProps({ placeholder }, prevState) {
        if(placeholder !== prevState.placeholder) {
            return {
                placeholder,
            }
        }
        return null
    }

    keyDownHandler(e) {
        e.keyCode === 13 && e.preventDefault();

        if (e.keyCode === 8 && this.value().length === 1) {
            // 修复EDGE & IE下无法删除最后一个字符的问题
            this.value('')
        }

        // 将setState放到setTimeout内是为了：修复FlexTextBox中state与实际输入不同步的问题
        // this.value()获取到的是输入框中的值，所以this.value()在系统事件使输入框中的值更新之后才能获取到最新的值
        // 如果setState在keyDownHandler函数中被同步执行，则this.value()会在系统更新输入框内容之前被调用
        // 导致state中的value一直比实际输入慢一步
        // 将setState放到setTimeout中，因为setTimeout是异步函数，所以会在系统更新输入框内容之后被调用
        // 此时this.value就是最新的实际输入的内容

        isFunction(this.props.onKeyDown) && this.props.onKeyDown(e);

        setTimeout(() => {
            this.setState({
                value: this.value(),
            }, () => this.fireValueChangeEvent(this.state.value));
        });

    }

    pasteHandler(e) {
        setTimeout(() => {
            this.setState({
                value: this.value(),
            }, () => {
                this.fireValueChangeEvent(this.state.value);
                isFunction(this.props.onPaste) && this.props.onPaste(this.state.value);
            });
        });
    }

    /**
     * 清空文本框输入
     */
    clear() {
        this.value('');
    }

    /**
     * 触发文本框输入变化事件
     * @param value 文本框内的输入值
     */
    private fireValueChangeEvent(value) {
        if (isFunction(this.props.onValueChange)) {
            this.props.onValueChange(value)
        }
    }

    /**
     * 获取或设定输入框的值
     * @param [text] 设定值
     */
    public value(text?: string): string {
        if (text !== undefined) {
            this.setState({ value: text }, () => this.fireValueChangeEvent(this.state.value))

            return 'textContent' in this.textBox ?
                (this.textBox.textContent = text.trim()) :
                (this.textBox.innerText = text.trim());
        } else {
            // innerText可以保留多行文本的换行，因此优先支持
            return 'innerText' in this.textBox ?
                this.textBox.innerText :
                this.textBox.textContent;
        }
    }

    focus() {
        this.textBox && this.textBox.focus()
    }
}