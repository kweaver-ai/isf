import React from 'react';
import classnames from 'classnames';
import ListBase from './ui.base';
import styles from './styles.desktop';

export default class List extends ListBase {

    render() {
        const Template = this.props.template;

        return (
            <ul
                className={classnames(styles.container, this.props.className)}
                selectId={this.props.selectIndex}
                ref="list"
                onMouseMove={this.MoveByMouse.bind(this)}
            >
                {
                    this.props.data.map((item, index) => {
                        return (
                            <li
                                key={index}
                                className={classnames(index === this.state.selectIndex ? styles['result-li-selected'] : styles['result-li'], this.props.resultClass)}
                                onMouseOver={() => this.handleMouseOver(index)}
                                onMouseLeave={() => this.handleMouseLeave()}
                                onMouseDown={this.props.onMouseDown.bind(this)}
                            >
                                <Template data={item} />
                            </li>
                        )
                    })
                }
            </ul>
        )
    }
}