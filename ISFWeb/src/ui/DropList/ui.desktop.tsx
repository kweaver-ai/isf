import React from 'react';
import styles from './styles.desktop'
import __ from './locale'

export default class DropList extends React.Component<any, any> {

    render() {
        return (
            <div className={styles['container']}>
                {
                    this.props.keyWords === '' || this.props.keyWords === null ? null : this.renderDropList()
                }
            </div>
        )
    }

    renderDropList() {
        return (

            <ul className={styles['ul']}>
                {

                    this.props.data.length > 0 ?
                        this.props.data.map((value, index) => {
                            return <li className={styles['li']} value={value} key={index} onClick={this.handleClick.bind(this, value)} title={value}>{value}</li>
                        })
                        :
                        <li className={styles['empty']} onClick={this.hideDropList.bind(this)}>
                            {__('未找到匹配的结果')}
                        </li>
                }
            </ul>

        )
    }

    handleClick(value) {
        this.setState({
            value,
        }, () => { this.props.onClick(value) })
    }

    hideDropList() {
        this.setState(() => { this.props.onClick(null) })
    }

}
