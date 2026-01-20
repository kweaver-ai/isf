import React from 'react';
import AutoCompleteListItem from '../AutoCompleteList.Item/ui.desktop';
import AutoCompleteListBase from './ui.base';
import styles from './styles.desktop';

export default class AutoCompleteList extends AutoCompleteListBase {

    static Item = AutoCompleteListItem;

    render() {
        return (
            <div
                role={this.props.role}
                ref={(list) => this.list = list}
                className={styles['autocomplte-list']}
                onMouseMove={this.setSelectByMouseMove.bind(this)}
            >
                <ul>
                    {
                        React.Children.map(this.props.children, (child, index) => {
                            return React.cloneElement(child, {
                                selected: this.props.selectIndex === index,
                                onMouseOver: (e) => this.handleMouseOver(e, index),
                                onMount: (listItem) => this.getItem(listItem, index),
                            })
                        })
                    }
                </ul>
            </div>
        )
    }
}