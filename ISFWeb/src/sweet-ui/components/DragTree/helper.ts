import React from 'react'

/**
 * 节点props
 */
export interface TreeNodeProps {
    eventKey?: Key;
    className?: string;
    style?: React.CSSProperties;

    // By parent
    expanded?: boolean;
    selected?: boolean;
    checked?: boolean;
    loaded?: boolean;
    loading?: boolean;
    halfChecked?: boolean;
    title?: React.ReactNode | ((data: DataNode) => React.ReactNode);
    dragOver?: boolean;
    dragOverGapTop?: boolean;
    dragOverGapBottom?: boolean;
    pos?: string;
    domRef?: React.Ref<HTMLDivElement>;
    data?: DataNode;
    isStart?: ReadonlyArray<boolean>;
    isEnd?: ReadonlyArray<boolean>;
    active?: boolean;
    onMouseMove?: React.MouseEventHandler<HTMLDivElement>;

    // By user
    isLeaf?: boolean;
    checkable?: boolean;
    selectable?: boolean;
    disabled?: boolean;
    disableCheckbox?: boolean;
    icon?: IconType;
    switcherIcon?: IconType;
    children?: React.ReactNode;
}

export type Key = string | number;

export interface DataNode {
    checkable?: boolean;
    children?: ReadonlyArray<DataNode>;
    disabled?: boolean;
    disableCheckbox?: boolean;
    icon?: IconType;
    isLeaf?: boolean;
    key?: Key;
    title?: React.ReactNode;
    selectable?: boolean;
    switcherIcon?: IconType;
    className?: string;
    style?: React.CSSProperties;
    data: any;
    parent: any;
}

export type IconType = React.ReactNode | ((props: TreeNodeProps) => React.ReactNode);

export type NodeInstance = React.Component<TreeNodeProps> & {
    selectHandle?: HTMLSpanElement;
};

export interface EventDataNode extends DataNode {
    expanded: boolean;
    selected: boolean;
    checked: boolean;
    loaded: boolean;
    loading: boolean;
    halfChecked: boolean;
    dragOver: boolean;
    dragOverGapTop: boolean;
    dragOverGapBottom: boolean;
    pos: string;
    active: boolean;
}

export type NodeMouseEventHandler<T = HTMLSpanElement> = (e: React.MouseEvent<T>, node: EventDataNode) => void;

export type NodeDragEventHandler<T = HTMLDivElement> = (e: React.MouseEvent<T>, node: NodeInstance) => void;

export interface DragTreeProps {
    selectable: boolean;
    showIcon: boolean;
    icon: IconType;
    switcherIcon: IconType;
    draggable: boolean;
    checkable: boolean | React.ReactNode;
    checkStrictly: boolean;
    disabled: boolean;
    loadData: (treeNode: EventDataNode) => Promise<void>;
    titleRender?: (node: DataNode) => React.ReactNode;

    onNodeClick: NodeMouseEventHandler;
    onNodeDoubleClick: NodeMouseEventHandler;
    onNodeExpand: NodeMouseEventHandler;
    onNodeSelect: NodeMouseEventHandler;
    onNodeCheck: (e: React.MouseEvent<HTMLSpanElement>, treeNode: EventDataNode, checked: boolean) => void;
    onNodeLoad: (treeNode: EventDataNode) => void;
    onNodeMouseEnter: NodeMouseEventHandler;
    onNodeMouseLeave: NodeMouseEventHandler;
    onNodeContextMenu: NodeMouseEventHandler;
    onNodeDragStart: NodeDragEventHandler;
    onNodeDragEnter: NodeDragEventHandler;
    onNodeDragOver: NodeDragEventHandler;
    onNodeDragLeave: NodeDragEventHandler;
    onNodeDragEnd: NodeDragEventHandler;
    onNodeDrop: NodeDragEventHandler;
}