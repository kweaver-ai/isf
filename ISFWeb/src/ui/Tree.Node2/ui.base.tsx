import React from 'react';
import { noop } from 'lodash';
import { promisify } from '@/util/accessor';

enum LoadStatus {
    PENDING, // 未加载子节点

    LOADING, // 正在加载子节点

    LOADED, // 已加载子节点
}

interface Status {
    disabled: boolean;
}

interface Props {

    onExpand?: (data) => void;

    _onSelect?: (data, select?) => void;

    onCollapse?: (data) => void;

    loader?: (data) => void;

    selected?: boolean;

    formatter?: (data) => void;

    getStatus?: (data) => Status;

    collapse?: boolean;

    childless?: boolean;
}

interface State {
    loadStatus: LoadStatus;

    collapsed: boolean;

    selected: boolean;

    status: Status;

}

export default class TreeNodeBase extends React.Component<any, any> {
    static defaultProps = {
        onExpand: noop,

        _onSelect: noop,

        onCollapse: noop,

        loader: noop,

        formatter: (data) => data.name,

        checkbox: false,    // 是否有复选框
        getStatus: (data) => ({
            disabled: false,
        }),
    };

    state: State = {
        loadStatus: LoadStatus.PENDING,

        collapsed: this.props.collapse, // false为展开

        selected: this.props.selected,

        status: this.props.getStatus(this.props.data),
    };

    componentDidUpdate(prevProps, prevState) {
        const { collapse, selected } = this.props
        if (collapse !== prevState.collapsed) {
            this.setState({
                collapsed: collapse,
            })
        }

        if (selected !== prevState.selected) {
            this.setState({ selected })
        }
    }
    /**
     * 选中节点
     */
    protected selectNode() {
        if (!this.state.status.disabled) {
            this.props._onSelect(this);
        }
        this.setState({
            selected: true,
        });
        this.props.onSelect();
    }

    /**
     * 展开／收起节点
     */
    protected toggleNode(node) {
        if (this.state.loadStatus === LoadStatus.PENDING) {
            this.setState({ loadStatus: LoadStatus.LOADING });
            promisify(this.props.loader(node)).then(() => {
                this.setState({ loadStatus: LoadStatus.LOADED })
            });
        }

        if (this.state.collapsed) {
            this.props.onExpand(node);
        } else {
            this.props.onCollapse(node);
        }

        this.setState({ collapsed: !this.state.collapsed });
    }

    /**
     * 复选框触发的变化
     */
    handleCheckBoxChange(target) {
        target.stopPropagation();
        this.props._onSelect(this, true)
    }
    protected getDisabledStatus(data): boolean {
        return this.props.getStatus(data).disabled
    }
}