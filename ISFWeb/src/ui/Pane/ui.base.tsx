import React from 'react';
import { noop } from 'lodash';

export default class PaneBase extends React.PureComponent<any, any> {

    static defaultProps = {

        /**
         * 图标
         */
        icon: '',

        /**
         * 颜色
         */
        color: '',

        /**
         * 图标引用资源
         */
        fallback: '',

        /**
         * 点击事件
         */
        onClick: noop,

        /**
         * 文字
         */
        label: '',

        /**
         * 样式
         */
        className: '',

        /**
         * 禁用
         */
        disabled: false,

        /**
         * 消息数目
         */
        msgNum: 0,

        /**
         * 浮动菜单
         */
        menuItems: [],

        /**
         * 展开浮动菜单的方式
         */
        triggerEvent: 'click',
    }

    state = {

        /**
         * 浮动菜单锚点
         */
        anchor: null,

        /**
         * 是否打开
         */
        open: false,
    }

    closePopMenu() {
        this.setState({
            open: false,
        })
        this.fireBlurEvent()
    }

    handleClick(event) {
        if (!this.props.disabled) {
            this.setState({
                anchor: event.currentTarget,
                open: true,
            });

            this.props.onClick();
        }

        event.stopPropagation();
    }

    fireBlurEvent() {
        this.props.onBlur();
    }
}