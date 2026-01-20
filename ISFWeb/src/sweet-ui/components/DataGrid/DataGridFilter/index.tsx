import React from 'react';
import { isEqual } from 'lodash';
import CheckBox from '../../CheckBox';
import View from '../../View';
import SweetIcon from '../../SweetIcon';
import Trigger from '../../Trigger';
import styles from './styles';
import __ from './locale';

interface DataGridFilterProps {
    /**
     * 过滤列关键字
     */
    filterKey: string;
    /**
     * 过滤选项
     */
    filters: ReadonlyArray<any>;

    /**
     * 过滤项改变时触发
     */
    onFilterChange: (filter: { key: string; filters: Array<string> }) => void;

    /**
     * 获取过滤菜单项key值
     */
    getFilterKey: (filter: any) => string;
}

interface DataGridFilterState {
    selectedKeys: Array<string>;
    filtered: boolean;
}

export default class DataGridFilter extends React.Component<DataGridFilterProps, DataGridFilterState> {
    state = {
        selectedKeys: [],
        filtered: false,
    };

    /**
     * 上一次的过滤项选中结果
     */
    prevFilterKeys = []

    filterOptions = this.props.filters.map(({ text, value }) => value)

    componentDidMount() {
        const selectedKeys = this.props.filters.reduce((prev, cur) => cur.checked ? [...prev, cur.value] : [...prev], [])

        this.setState({
            selectedKeys,
            filtered: !!selectedKeys.length,
        })

        this.prevFilterKeys = selectedKeys
    }

    changeSelectedKeys = (selectedKeys: Array<string>) => {
        this.setState({
            selectedKeys,
        })
    }

    /**
     * 点击全选复选框
     */
    handleCheckAllFilter = (checked: boolean) => {
        this.changeSelectedKeys(checked ? this.filterOptions : [])
    };

    /**
     * 点击全选菜单项
     */
    handleClickFilterAll = () => {
        this.changeSelectedKeys(this.state.selectedKeys.length === this.props.filters.length ? [] : this.filterOptions)
    };

    /**
     * 点击过滤项复选框
     */
    handleFilterChange = (checked: boolean, item: any) => {
        this.changeSelectedKeys(checked ? [...this.state.selectedKeys, item.value] : this.state.selectedKeys.filter((key) => key !== item.value))
    };

    /**
     * 点击非全选菜单项
     */
    handleClickFilterItem = (item: any) => {

        if (this.state.selectedKeys.length > 0 && this.state.selectedKeys.find((key) => key === item.value)) {
            this.changeSelectedKeys(this.state.selectedKeys.filter((key) => key !== item.value))
        } else {
            this.changeSelectedKeys([...this.state.selectedKeys, item.value])
        }
    };

    /**
     * 关闭菜单时触发
     */
    handleBeforeFilterClose = () => {
        this.setState({
            filtered: this.state.selectedKeys.length > 0 && this.state.selectedKeys.length <= this.filterOptions.length,
        });
        if (!isEqual(this.prevFilterKeys, this.state.selectedKeys)) {
            this.props.onFilterChange({ key: this.props.filterKey, filters: this.state.selectedKeys });
            this.prevFilterKeys = this.state.selectedKeys
        }
    };

    /**
     * 处理CheckBox点击事件
     */
    handleCheckBoxClicked = (event: React.MouseEvent<HTMLLabelElement>) => {
        event.stopPropagation();
    }

    render() {
        const { selectedKeys, filtered } = this.state;
        const { filters, element } = this.props;

        return (
            <Trigger
                anchorOrigin={['right', 'bottom']}
                alignOrigin={['right', 'top']}
                freeze={true}
                element={element}
                renderer={({ setPopupVisibleOnClick }) => (
                    <SweetIcon
                        key={'filter'}
                        size={12}
                        title={__('筛选')}
                        color={filtered ? '#40A9FF' : '#505050'}
                        name={'filter'}
                        onClick={setPopupVisibleOnClick}
                        className={styles['filter-icon']}
                    />
                )}
                onBeforePopupClose={this.handleBeforeFilterClose}
            >
                <ul className={styles['drop-menu']}>
                    <li
                        key={'filter-all'}
                        className={styles['drop-menu-item']}
                        onClick={this.handleClickFilterAll}
                    >
                        <CheckBox
                            checked={selectedKeys.length === this.filterOptions.length}
                            onChange={(event) => this.handleCheckAllFilter(event.target.checked)}
                            onClick={(event) => this.handleCheckBoxClicked(event)}
                            indeterminate={selectedKeys.length > 0 && selectedKeys.length < this.filterOptions.length}
                        />
                        <View inline={true} className={styles['item-text']}>
                            {__('全选')}
                        </View>
                    </li>
                    {filters.map((item) => {
                        return (
                            <li
                                key={this.props.getFilterKey(item)}
                                className={styles['drop-menu-item']}
                                onClick={() => this.handleClickFilterItem(item)}
                            >
                                <CheckBox
                                    checked={selectedKeys && selectedKeys.indexOf(item.value) >= 0}
                                    onChange={(event) => this.handleFilterChange(event.target.checked, item)}
                                    onClick={(event) => this.handleCheckBoxClicked(event)}
                                />
                                <View inline={true} className={styles['item-text']}>
                                    {item.text}
                                </View>
                            </li>
                        );
                    })}
                </ul>
            </Trigger>
        );
    }
}
