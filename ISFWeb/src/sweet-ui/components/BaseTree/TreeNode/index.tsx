import React from 'react';
import View from '../../View';
import styles from './styles';

export type Locator = ReadonlyArray<number>;

interface TreeNodeProps {
    data: any;

    locator: Locator;

    renderNode: (node: any, locator: Locator) => React.ReactNode;

    children?: React.ReactNode;
}

type TreeNode = React.FunctionComponent<TreeNodeProps>;

const TreeNode: TreeNode = function TreeNode({ children, data, locator, renderNode }) {
    return (
        <View>
            <View>{renderNode(data, locator)}</View>
            {React.Children.count(children) > 0 ? <View className={styles['branch']}>{children}</View> : null}
        </View>
    );
};

export default TreeNode;
