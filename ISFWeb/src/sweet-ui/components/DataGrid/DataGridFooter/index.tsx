import React from 'react';
import Pager, { PagerProps, PageChangeHandler as PageChangeHandlerType } from '../../Pager';

export type PageChangeHandler = PageChangeHandlerType;

interface DataGridFooterProps {
    DataGridPager: PagerProps;
}

const DataGridFooter: React.FC<DataGridFooterProps> = function DataGridFooter({ DataGridPager }) {
    const { page, total, size, onPageChange } = DataGridPager;

    return DataGridPager ? <Pager {...{ page, total, size, onPageChange }} /> : null;
};

export default DataGridFooter;
