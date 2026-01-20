import React from 'react';
import { isEqual } from 'lodash';

export default class SelectMenuBase extends React.Component<UI.SelectMenu2.Props, UI.SelectMenu2.State> {
    static defaultProps = {
        label: '',
        candidateItems: [{ name: '' }],
        selectValue: '',
    }

    state = {
        selectValue: this.props.selectValue ? this.props.selectValue : this.props.candidateItems[0],
        hover: true,
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        if (!isEqual(prevState.selectValue, nextProps.selectValue)) {
            return {
                selectValue: nextProps.selectValue,
            }
        }
        return null
    }

    /**
     * 选择候选菜单项
     */
    protected handleClickCandidateItem(e, item) {
        this.setState({
            hover: true,
        }, () => {
            this.props.onSelect(item);
        })

    }

    // 延迟定时器
    timeout: number | null = null;

    /**
     * 鼠标进入任意选项时触发
     */
    protected handleMouseEnter() {
        this.setState({
            hover: false,
        })
    }

    /**
     * 点击弹出框外时触发
     */
    protected handleCloseMenuWhenBlur() {
        if (this.timeout) {
            clearTimeout(this.timeout)
        }

        this.timeout = setTimeout(() => {
            this.setState({
                hover: true,
            })
        }, 150)

        if (typeof (this.props.onBeforePopupClose) === 'function') {
            this.props.onBeforePopupClose()
        }
    }
}