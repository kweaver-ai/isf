import * as React from 'react';
import { noop, isFunction } from 'lodash';
import { NodeType } from '@/core/organization';
import WebComponent from '../webcomponent';

interface OrganizationTreeBoxProps extends React.Props<void> {
    /**
     * 选择展示的类型
     */
    selectType: ReadonlyArray<NodeType>;

    /**
     * 是否展示搜索框
     */
    isShowSearch?: boolean;

    /**
     * 搜索框宽度
     */
    searchWidth?: number | string;

    /**
     * 是否展示为分配组
     */
    isShowUndistributed?: boolean;

    /**
     * 选中方法
     */
    onSelectionChange?: (selections: any) => void;

    /**
     * 搜索框选中方法
     */
    selectSearch?: (selections: any) => void;
}

export default class OrganizationTreeBoxBase extends WebComponent<OrganizationTreeBoxProps, any> {
    static defaultProps = {
        isShowSearch: false,
        isShowUndistributed: false,
        onSelectionChange: noop,
    }

    state = {
        // 选择的部门/文档
        selected: [],
    }

    /**
    * 选择访问者
    * @param value 访问者
    */
    async selectDep(value: any) {
        isFunction(this.props.selectSearch) && this.props.selectSearch(value)
    }

}