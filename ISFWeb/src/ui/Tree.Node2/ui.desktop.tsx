import React from 'react';
import classnames from 'classnames';
import TreeNodeBase from './ui.base';
import LinkChip from '../LinkChip/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import CheckBox from '../CheckBox/ui.desktop';
import styles from './styles.desktop';
import collapsed from './assets/collapsed.desktop.png';
import expanded from './assets/expanded.desktop.png';

export default class TreeNode extends TreeNodeBase {
    render() {

        return (
            <div className={classnames(styles['node'])}>
                <div className={classnames(styles['detail'], { [styles['selected']]: this.state.selected })}>

                    <UIIcon
                        size={13}
                        className={classnames([styles['toggler']], { [styles['leaf']]: this.props.isLeaf })}
                        code={this.state.collapsed ? '\uf04e' : '\uf04c'}
                        fallback={this.state.collapsed ? collapsed : expanded}
                        onClick={() => this.toggleNode(this.props.data)}
                    />

                    <LinkChip
                        className={classnames(styles['name'],
                            {
                                [styles['disabled']]: this.getDisabledStatus(this.props.data),
                                [styles['name-selected']]: this.state.selected,
                            })}
                        onClick={this.selectNode.bind(this)}
                        onDoubleClick={() => this.toggleNode(this.props.data)}
                    >
                        {
                            this.props.checkbox
                                ? <CheckBox checked={this.state.selected} onClick={this.handleCheckBoxChange.bind(this)} />
                                : null
                        }
                        {
                            this.props.formatter(this.props.data)
                        }
                    </LinkChip>
                </div>
                {
                    this.props.children ?
                        <div className={classnames(styles['branch'], { [styles['collapsed']]: this.state.collapsed, [styles['expanded']]: !this.state.collapsed })}>
                            {
                                this.props.children
                            }
                        </div>
                        : null
                }
            </div>
        )
    }
}