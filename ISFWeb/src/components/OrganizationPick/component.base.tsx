import * as React from 'react';
import { noop, uniqBy, isEqual } from 'lodash';
import { NodeType } from '@/core/organization';
import WebComponent from '../webcomponent';

export default class OrganizationPickBase extends WebComponent<any, any> {

    static defaultProps = {
        disabled: false,
        onConfirm: noop,
        onCancel: noop,
        userid: '',
        selectType: [NodeType.ORGANIZATION, NodeType.DEPARTMENT],
        data: [],
        isShowUndistributed: false,
        onSelectionChange: noop,
        converterIn: (x) => x,
    }

    state = {
        data: [],
    }

    componentDidMount() {
        this.setState({
            data: this.props.data.map(this.props.converterIn),
        })
    }

    componentDidUpdate(prevProps, prevState) {
        if (!isEqual(this.props.data, prevProps.data)) {
            this.setState({
                data: this.props.data.map(this.props.converterIn),
            })
        }
    }

    /**
     * 选择共享者
     * @param value 共享者
     */
    async selectDep(value) {
        this.setState({
            data: uniqBy(this.state.data.concat(value), 'id'),
        }, () => {
            this.props.onSelectionChange(this.state.data.map(this.props.convererOut))
        })
    }

    /**
     * 删除已选部门
     * @param dep 部门
     */
    deleteSelectDep(sharer: Node) {
        this.setState({
            data: this.state.data.filter((value) => value.id !== sharer.id),
        }, () => {
            this.props.onSelectionChange(this.state.data.map(this.props.convererOut))
        })
    }

    /**
     * 清空已选择部门
     */
    clearSelectDep() {
        this.setState({
            data: [],
        })
        this.props.onSelectionChange([])
    }
}