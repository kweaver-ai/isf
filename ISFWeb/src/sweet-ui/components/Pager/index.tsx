import React from 'react';
import { isFunction } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import IconButton from '../IconButton';
import BaseNumberInput from '../BaseNumberInput';
import styles from './styles';
import SweetIcon from '../SweetIcon';
import __ from './locale';

/**
 * 页码变化事件
 * @param detail 页码变化信息
 * @param detail.page 页码
 * @param detail.size 每页大小
 */
export type PageChangeEvent = SweetUIEvent<{ page: number; size: number }>;

/**
 * 当前页码发生变化时触发
 */
export type PageChangeHandler = (event: PageChangeEvent) => void;

export interface PagerProps extends React.ClassAttributes<void> {
    /**
     * 总页数
     */
    total: number;

    /**
     * 当前页码
     */
    page?: number;

    /**
     * 每页条目数
     */
    size?: number;

    /**
     * 当前页码发生变化时触发
     */
    onPageChange?: PageChangeHandler;
}

interface PagerState {
    page: number;
    pageInput: number | string;
}

/**
 * 默认开始页
 */
const DEFAULT_PAGE = 1;

/**
 * 默认分页大小
 */
const DEFAULT_PAGE_SIZE = 20;

export default class Pager extends React.PureComponent<PagerProps, PagerState> {
    static defaultProps = {
        page: DEFAULT_PAGE,
        size: DEFAULT_PAGE_SIZE,
    };

    state = {
        page: this.props.page,
        pageInput: this.props.page,
    };

    // 上一次的有效页码
    lastVaildPage: number = DEFAULT_PAGE

    static getDerivedStateFromProps({ page }: { page: number }, prevState: PagerState) {
        if (page !== prevState.page) {
            return {
                page,
                pageInput: page,
            }
        }

        return null
    }

    dispatchPageChangeEvent = createEventDispatcher(this.props.onPageChange, (event: PageChangeEvent) => {
        this.setState({
            page: event.detail.page,
            pageInput: event.detail.page,
        });
    });

    /**
     * 更新页码
     * @param nextPage 页码
     */
    private updatePage(nextPage: number) {
        const { onPageChange } = this.props;

        this.setState({ page: nextPage, pageInput: nextPage }, () => {
            const { page } = this.state;
            const { size } = this.props;

            if (isFunction(onPageChange)) {
                onPageChange({ page, size });
            }
        });
    }

    private handlePageChange(nextPage: number) {
        const { total, size = DEFAULT_PAGE_SIZE } = this.props;

        if (nextPage <= 0) {
            return;
        }

        if (nextPage > Math.ceil(total / size)) {
            return;
        }

        this.lastVaildPage = nextPage
        this.dispatchPageChangeEvent({ page: nextPage, size });
    }

    handleEmptyInput = () => {
        this.setState({
            pageInput: '',
        });
    }

    // 失焦时，如果页码为空则还原为上一次的有效页码
    handleBlur = () => {
        if (this.state.pageInput === '') {
            this.setState({
                pageInput: this.lastVaildPage,
            });
        }
    }

    pageBox = ({ page, totalPage }: { page: number; totalPage: number }) => {
        return (
            <BaseNumberInput
                value={page}
                className={styles['page-box']}
                min={1}
                max={totalPage}
                onValueChange={(event: any) => this.handlePageChange(event.detail)}
                onEmptyValue={this.handleEmptyInput}
                onBlur={this.handleBlur}
            />
        );
    }

    render() {
        const { page = DEFAULT_PAGE, pageInput } = this.state;
        const { total, size = DEFAULT_PAGE_SIZE } = this.props;
        const totalPage = Math.max(Math.ceil(total / size), 1);

        return (
            <View className={styles['pager']}>
                <View className={styles['pages']}>
                    <IconButton
                        icon={<SweetIcon name="first" />}
                        className={styles['icon']}
                        onClick={() => this.handlePageChange(1)}
                        disabled={page === 1}
                    />
                    <IconButton
                        icon={<SweetIcon name="prev" />}
                        className={styles['icon']}
                        onClick={() => this.handlePageChange(page - 1)}
                        disabled={page === 1}
                    />
                    <View className={styles['page-info']}>
                        {__('第 ')}
                        {this.pageBox({ page: pageInput, totalPage })}
                        {__(' 页，共 ')}
                        {totalPage}
                        {__(' 页')}
                    </View>
                    <IconButton
                        icon={<SweetIcon name="next" />}
                        className={styles['icon']}
                        onClick={() => this.handlePageChange(page + 1)}
                        disabled={page === totalPage}
                    />
                    <IconButton
                        icon={<SweetIcon name="last" />}
                        className={styles['icon']}
                        onClick={() => this.handlePageChange(totalPage)}
                        disabled={page === totalPage}
                    />
                </View>
                <View className={styles['count']}>
                    {__('显示 ')}
                    {total > 0 ? (page - 1) * size + 1 : 0} - {Math.min(page * size, total)}
                    {__(' 条，共 ')}
                    {total}
                    {__(' 条')}
                </View>
            </View>
        );
    }
}
