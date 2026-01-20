import React from 'react';
import { noop } from 'lodash';

interface Props {
    /**
     * 角色
     */
    role?: string;

    /**
     * 禁用（灰化，不可操作）
     */
    disabled: boolean;

    /**
     * 选中节点时触发
     * @param node 节点数据
     */
    onSelectionChange?: (node: Record<string, any> | Array<Record<string, any>>) => any;

    /**
     * 选择模式: 'single'--单选, 'multi'--多选, 'cascade'--级联
     */
    selectMode: string;
}

export default class TreeBase extends React.Component<Props, any> {
    static defaultProps = {
        disabled: false,

        onSelectionChange: noop,

        selectMode: 'single',
    }

    selected = this.props.selectMode === 'single' ? null : []

    props: Props;

    /**
     * 遍历当前所有节点
     * @param callback 对每个节点执行函数
     */
    private forEachNode(callback: Function = noop): void {
        const iterateNodes = function (nodes: React.ReactNode, callback: Function = noop) {
            React.Children.forEach(nodes, (node) => {
                callback(node);

                if (node.props.children) {
                    iterateNodes(node.props.children, callback)
                }
            })
        }

        iterateNodes(this.props.children, callback);
    }

    /**
     * 选择节点
     * @param target 选择目标
     * @param checkbox： true -- 勾选或取消复选框导致 ; false -- 直接选中TreeNode导致
     */
    private selectNode(target: React.ReactNode, checkbox: boolean) {
        switch (this.props.selectMode) {
            case 'single': {
                // 单选
                if (this.selected !== target) {
                    // 未选中 --> 选中， 且原来选中的要取消选中
                    if (this.selected) {
                        this.selected.setState({ selected: false })
                    }
                    target.setState({ selected: true })
                    this.selected = target
                }
                this.props.onSelectionChange(target.props.data)
                break
            }
            case 'multi': {
                // 多选
                if (checkbox) {
                    // 通过复选框触发的变化（1）如果已选中，要取消选中（2）如果未选中，要选中
                    if (this.selected.some((item) => item === target)) {
                        // 已选中 -->  取消
                        target.setState({ selected: false })
                        this.selected = this.selected.filter((item) => item !== target)
                    } else {
                        // 未选中 --> 选中
                        target.setState({ selected: true })
                        this.selected = [...this.selected, target]
                    }
                } else {
                    // 通过选中整条数据触发的变化,
                    // (1)当这条数据有checkbox才触发变化，否则不变。
                    // (2)取消所有已选中状态。
                    // (3)这条数据如果已被选中，则取消；如果没有被选中，则选中
                    if (target.props.checkbox) {
                        this.selected.forEach((item) => item.setState({ selected: false }))
                        if (this.selected.includes(target)) {
                            // 已选中 -->  取消
                            this.selected = []
                        } else {
                            // 未选中 --> 选中
                            target.setState({ selected: true })
                            this.selected = [target]
                        }
                    }
                }
                this.props.onSelectionChange(this.selected.map((item) => item.props.data))
                break
            }
            case 'cascade': {
                // 级联选择
            }
        }
    }

    /**
     * 对Tree.Node进行扩展
     * @param children 某一层级下的所有Tree.Node
     */
    protected extendsChildren(children: Array<React.ReactNode>) {
        return React.Children.map(children, (node) => React.cloneElement(node, {
            _onSelect: this.selectNode.bind(this),
            children: !node.props.data ? node.props.children : (node.props.children ? this.extendsChildren(node.props.children) : null),
        }));
    }

    /**
     * 防止事件传播
     */
    protected stopPropagation = (e): void => {
        if (this.props.disabled) {
            e.stopPropagation()
        }
    }
}