import React from 'react';
import { noop } from 'lodash';

export default class DropBoxBase extends React.PureComponent<any, any> {
    static defaultProps = {
        fontIcon: '\uf04c',

        dropAlign: 'bottom left',

        formatter: (val) => val,

        onBlur: noop,

        onChange: noop,

        onBeforeDeactivate: () => true,

        disabled: false,

        onActive: noop,
    }

    state = {
        active: false,
    }

    deactivePrevented = false;

    public preventDeactivate(e) {
        this.deactivePrevented = true;
    }

    onSelectBlur(e) {
        if (this.deactivePrevented) {
            this.refs.select.focus();
        } else {
            this.deactivate();
            this.fireBlurEvent();
        }

        this.deactivePrevented = false;
    }

    /**
     * 取消激活
     */
    deactivate() {
        if (!this.props.disabled) {
            this.setState({
                active: false,
            })
        }
    }

    /**
     * 切换选中状态
     */
    toggleActive() {
        if (!this.props.disabled) {
            this.setState({
                active: !this.state.active,
            })
        }
    }

    fireBlurEvent() {
        this.props.onBlur();
    }

    /**
     * 打开时触发
     */
    onOpen() {
        if (typeof this.props.onActive === 'function') {
            this.props.onActive(this.state.active);
        }

    }

    /**
     * 关闭时触发
     */
    onClose(close) {
        close();
        this.setState({
            active: false,
        })

        if (typeof this.props.onActive === 'function') {
            this.props.onActive(false);
        }
    }

    /**
     * 点击下拉选项
     */
    protected toggleSelect(close) {
        this.setState({
            active: false,
        })
        close()
    }

}