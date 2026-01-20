import React from 'react';
import { noop } from 'lodash';

export default class ValidateAreaBase extends React.Component<UI.ValidateArea.Props, UI.ValidateArea.State> implements UI.TextBox.Element {
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

    /**
     * 聚焦文本框
     * @param event 事件对象
     */
    focus(event) {
        this.setState({ focus: true });
        this.props.onFocus(event);
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
}