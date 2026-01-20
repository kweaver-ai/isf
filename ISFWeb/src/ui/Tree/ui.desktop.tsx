import React from 'react';
import classnames from 'classnames';
import TreeBase from './ui.base';
import TreeNode from '../Tree.Node/ui.desktop';
import styles from './styles.desktop';

export default class Tree extends TreeBase {
    static Node = TreeNode;

    render() {
        return (
            <div
                role={this.props.role}
                className={classnames(styles['root'], { [styles['disabled']]: this.props.disabled })}
                onClickCapture={(e) => this.stopPropagation(e)}
                onDoubleClickCapture={(e) => this.stopPropagation(e)}
            >
                {
                    this.extendsChildren(this.props.children)
                }
            </div>
        )
    }
}