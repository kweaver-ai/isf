import React from 'react';
import PaginatorBase from './ui.base';
import UIIcon from '../UIIcon/ui.desktop';
import TextInput from '../TextInput/ui.desktop';
import FlexBox from '../FlexBox/ui.desktop';
import __ from './locale';
import styles from './styles.desktop';

function PageInput({ onChange, value, onBlur }) {
    return (
        <div className={styles['page-input']}>
            <TextInput {...{ value, onChange, onBlur }} />
        </div>
    )
}

export default class Paginator extends PaginatorBase {
    render() {
        return (
            <div className={styles['paginator']}>
                <FlexBox>
                    <FlexBox.Item align="middle left">
                        <div className={styles['page']}>
                            <div className={styles['action']}>
                                <UIIcon
                                    disabled={this.state.pages === 0 || this.state.page === 1}
                                    code={'\uf010'}
                                    size={16}
                                    onClick={this.first.bind(this)}
                                />
                            </div>
                            <div className={styles['action']}>
                                <UIIcon
                                    disabled={this.state.pages === 0 || this.state.page === 1}
                                    code={'\uf012'}
                                    size={16}
                                    onClick={this.navigate.bind(this, -1)}
                                />
                            </div>
                            {__('第')}
                            <div className={styles['page-info']}>
                                <PageInput
                                    value={this.state.pageInput}
                                    onChange={this.goto.bind(this)}
                                    onBlur={this.handlePageInputBlur.bind(this)}
                                />
                            </div>
                            {__('页，共')}
                            <div className={styles['page-info']}>{Math.max(this.state.pages, 1)}</div>
                            {__('页')}
                            <div className={styles['action']}>
                                <UIIcon
                                    disabled={this.state.pages === 0 || this.state.page === this.state.pages}
                                    code={'\uf011'}
                                    size={16}
                                    onClick={this.navigate.bind(this, 1)}
                                />
                            </div>
                            <div className={styles['action']}>
                                <UIIcon
                                    disabled={this.state.pages === 0 || this.state.page === this.state.pages}
                                    code={'\uf00f'}
                                    size={16}
                                    onClick={this.last.bind(this)}
                                />
                            </div>
                        </div>
                    </FlexBox.Item>
                    <FlexBox.Item align="middle right">
                        <div className={styles['total']}>
                            {__('显示')}
                            <div className={styles['page-info']}>
                                {
                                    `${this.props.total === 0 ? 0 : this.props.limit * (this.state.page - 1) + 1} - ${Math.min(this.props.limit * this.state.page, this.props.total)}`
                                }
                            </div>
                            {__('条，共')}
                            <div className={styles['page-info']}>
                                {
                                    this.props.total
                                }
                            </div>
                            {__('条')}
                        </div>
                    </FlexBox.Item>
                </FlexBox>
            </div>
        )
    }
}