import React from 'react';
import View from '../View';
import TreeNode, { Locator } from './TreeNode';

export type BaseTreeNodeData = {
    children?: ReadonlyArray<BaseTreeNodeData>;

    [key: string]: any;
};

export type BaseTreeData = ReadonlyArray<BaseTreeNodeData>;

export type BaseTreeNodeKeyExtractor = (node: any) => string;

interface BaseTreeProps {
    data: BaseTreeData;

    renderBaseTreeNode: (node: any, locator: Locator) => React.ReactNode;

    keyExtractor: BaseTreeNodeKeyExtractor;
}

const BaseTree: React.FunctionComponent<BaseTreeProps> = function BaseTree({
    data,
    keyExtractor,
    renderBaseTreeNode,
}) {
    const renderBaseTreeNodes = ({
        nodes,
        parentLocator = [],
    }: {
        nodes: ReadonlyArray<any>;
        parentLocator: Locator;
    }): React.ReactNode => {
        return nodes.map((node, i) => {
            const locator = [...parentLocator, i];

            return (
                <TreeNode key={keyExtractor(node)} data={node} locator={locator} renderNode={renderBaseTreeNode}>
                    {node.children ? renderBaseTreeNodes({ nodes: node.children, parentLocator: locator }) : null}
                </TreeNode>
            );
        });
    };

    return <View>{renderBaseTreeNodes({ nodes: data })}</View>;
};

export default BaseTree;
