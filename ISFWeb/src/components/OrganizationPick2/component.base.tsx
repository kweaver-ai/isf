import * as React from 'react';
import { noop, uniqBy, isEqual } from 'lodash';
import { NodeType } from '@/core/organization';
import WebComponent from '../webcomponent';
import { OrganizationPick2Props, OrganizationPick2State, Node, SearchDepSelectDepInfo, SearchDepSelectUserInfo } from './helper';
import __ from './locale';

export default class OrganizationPick2Base extends WebComponent<OrganizationPick2Props, OrganizationPick2State> {
    static defaultProps = {
        disabled: false,
        selectType: [NodeType.ORGANIZATION, NodeType.DEPARTMENT],
        describeTip: __('请选择范围：'),
        selections: [],
        extraRoots: [],
        onRequestChangeSelections: noop,
        convererOut: (x) => x,
    }

    state: OrganizationPick2State = {
        selections: typeof this.props.converterIn === 'function' ?
            this.props.selections.map(this.props.converterIn)
            : this.props.selections,
    }

    departmentTree = {
        getSelections: () => [],
        cancelSelections: noop,
    };

    componentDidUpdate(prevProps, prevState) {
        if (!isEqual(this.props.selections, prevProps.selections)) {
            const { selections, converterIn } = this.props

            this.setState({
                selections: typeof converterIn === 'function' ? selections.map(converterIn) : selections,
            })
        }
    }

    /**
     * 添加到已选列表
     */
    protected addToList = async (): void => {
        let extras = []

        const addition = (await this.departmentTree.getSelections()).reduce((prev, cur) => {
            if (this.props.extraRoots.length && this.props.extraRoots.some((item) => item.id === cur.id)) {
                extras = [...extras, cur]
                return prev
            } else {
                return [...prev, cur]
            }
        }, [])

        const selections = uniqBy([...extras, ...this.state.selections, ...addition], 'id')

        this.setState({
            selections,
        })

        this.departmentTree.cancelSelections()

        this.props.onRequestChangeSelections(selections.map(this.props.convererOut))
    }

    /**
     * 点击搜索下拉框选择
     * @param value 选中项
     */
    protected select = async (value: SearchDepSelectDepInfo | SearchDepSelectUserInfo): Promise<void> => {
        const selections = uniqBy(this.state.selections.concat(value), 'id')

        this.setState({
            selections,
        })

        this.props.onRequestChangeSelections(selections.map(this.props.convererOut))
    }

    /**
     * 删除已选
     */
    protected deleteSelected = (node: Node): void => {
        const selections = this.state.selections.filter((value) => value.id !== node.id)

        this.setState({
            selections,
        })

        this.props.onRequestChangeSelections(selections.map(this.props.convererOut))
    }

    /**
     * 清空已选
     */
    protected clearSelections = (): void => {
        this.setState({
            selections: [],
        })

        this.props.onRequestChangeSelections([])
    }
}