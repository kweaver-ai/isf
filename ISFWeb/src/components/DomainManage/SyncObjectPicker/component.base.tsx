import * as React from 'react'
import { noop } from 'lodash';
import { NodeType } from '@/core/organization';
import WebComponent from '../../webcomponent'
import __ from './locale'

export default class SyncObjectPickerBase extends WebComponent<any, any> {

    static defaultProps = {
        onConfirm: noop,
        onCancel: noop,
        domainId: -1,
        userid: '',
        selectType: [NodeType.ORGANIZATION, NodeType.DEPARTMENT],
        data: [],
        converterIn: (x) => x,
    }

    state = {
        data: this.props.data,
    }

    /**
     * 树结构选中的节点
     */
    selectTree: ReadonlyArray<any> = this.props.data

    /**
     * 选择部门
     */
    protected selectDep = (data: ReadonlyArray<any>, search?: boolean): void => {
        if (search) {
            this.setState({
                data,
            })
        } else {
            this.selectTree = data
        }
    }

    /**
     * 添加
     */
    protected addTreeData = (): void => {
        if (this.selectTree.length) {
            this.setState({
                data: this.selectTree,
            })
            this.ref.cancelSelections()
            this.selectTree = []
        }
    }

    /**
     * 删除已选部门
     * @param dep 部门
     */
    deleteSelectDep = (dep: Node) => {
        this.setState({
            data: this.state.data.filter((value) => value.pathName !== dep.pathName),
        })
    }

    /**
     * 清空已选择部门
     */
    clearSelectDep = () => {
        this.setState({
            data: [],
        })
    }

    /**
     * 取消本次操作
     */
    cancelAddDep = () => {
        this.clearSelectDep()
        this.props.onCancel()
    }

    /**
     * 确定本次操作
     */
    confirmAddDep = () => {
        this.props.onConfirm(this.state.data.map(this.props.convererOut))
    }

    /**
     * 禁用
     */
    protected getNodeStatus = (node: Core.ShareMgnt.ncTDepartmentInfo): { disabled: boolean } => {
        return { disabled: this.props.data.some((item) => item.id === node.id) }
    }
}