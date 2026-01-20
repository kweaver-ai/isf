import React from 'react';
import { union, without } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import { Locator } from '../BaseTree/TreeNode';
import BaseTree, { BaseTreeNodeData, BaseTreeData, BaseTreeNodeKeyExtractor } from '../BaseTree';
import CheckBox from '../CheckBox';
import styles from './styles';
import SweetIcon from '../SweetIcon';

type Selection = any | ReadonlyArray<any>;

interface TreeProps extends React.ClassAttributes<void> {
    data: BaseTreeData;

    keyExtractor: BaseTreeNodeKeyExtractor;

    renderTreeNode: (node: any) => React.ReactNode;

    onSelectionChange?: (node: any) => void;
}

interface TreeState {
    selection?: Selection;
}

export default class Tree extends React.Component<TreeProps, TreeState> {
    state: TreeState = {};

    dispatchSelectionChangeEvent = createEventDispatcher(this.props.onSelectionChange, (event) => {
        this.changeSelection(event);
    });

    private changeSelection = (event: SweetUIEvent<Selection>) => {
        const { detail: selection } = event;

        this.setState({
            selection,
        });
    };

    private handleSelectionChange = (locator: Locator, checked: boolean) => {
        const { selection = [] } = this.state;
        const [ancestors, node] = this.findNodesByLocator(locator);
        const children = this.findChildrenNodes(node);
        let nextSelection;

        // 选中
        if (checked) {
            nextSelection = union(selection, ancestors);
        } else {
            // 取消选中
            nextSelection = without(selection, ...ancestors);
        }

        this.dispatchSelectionChangeEvent(nextSelection);
    };

    /**
     * 根据locator查找所有节点
     * @param locator 定位数组
     */
    private findNodesByLocator(locator: Locator) {
        const { data } = this.props;

        // root 不用从children中获取，因此作为后续reduce的初始值
        const [root, ...rest] = locator;

        // 假设locator为 [0, 7, 9]
        // 则对应的节点为 arr = [data[0], data[0].children[7], data[0].children[7].children[9]]
        // 除 rest 外其他项转换成下标查找为 arr[0].children[7], arr[1].children[9]
        return rest.reduce(
            (result: ReadonlyArray<any>, nodeIndex: number, index: number) => {
                return [...result, result[index].children[nodeIndex]];
            },
            [data[root]],
        );
    }

    private findChildrenNodes(node: BaseTreeNodeData) { }

    private renderBaseTreeNode = (node: BaseTreeNodeData, locator: Locator) => {
        const { renderTreeNode } = this.props;
        const { selection = [] } = this.state;
        const checked = selection.includes(node);

        return (
            <View>
                <View className={styles['expand']}>
                    <SweetIcon name="x" />
                </View>
                <View className={styles['checkbox']}>
                    <CheckBox
                        checked={checked}
                        onCheckedChange={(event) => this.handleSelectionChange(locator, event.detail)}
                    />
                </View>
                {renderTreeNode(node)}
            </View>
        );
    };

    render() {
        const { data, keyExtractor } = this.props;

        return <BaseTree {...{ data, keyExtractor }} renderBaseTreeNode={this.renderBaseTreeNode} />;
    }
}
