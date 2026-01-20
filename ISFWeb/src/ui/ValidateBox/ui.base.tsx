import React from 'react';
import { noop } from 'lodash';

export default class ValidateBoxBase extends React.Component<UI.ValidateBox.Props, UI.ValidateBox.State> implements UI.TextBox.Element {
    static defaultProps = {
        onFocus: noop,

        onBlur: noop,

        validateMessages: [],

        align: 'right',
    }

    state = {
        focus: false,
        hover: false,
    }

    textInput = null

    /**
     * 聚焦文本框
     * @param event 事件对象
     */
    public focus(event) {
        if (this.textInput) {
            this.textInput.focus()
        }
        this.setState({ focus: true });
        this.props.onFocus && this.props.onFocus(event);
    }

    /**
     * 失焦文本框
     * @param event 事件对象
     */
    blur(event) {
        this.setState({ focus: false });
        this.props.onBlur(event);
    }

    /**
     * 鼠标悬浮
     */
    mouseOver() {
        this.setState({
            hover: true,
        })
    }

    /**
     * 鼠标移除
     */
    mouseOut() {
        this.setState({
            hover: false,
        })
    }

    /**
     * 鼠标回车事件
     */
    onEnter(e) {
        if (this.props.onEnter) {
            this.props.onEnter(e)
        }
    }
}