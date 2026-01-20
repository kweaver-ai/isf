import * as React from 'react';
import * as classnames from 'classnames';
import { Text, UIIcon, SearchBox } from '@/ui/ui.desktop';
import { DataGrid, Button } from '@/sweet-ui';
import { NodeType } from '@/core/organization';
import ListTipComponent from '../../../ListTipComponent/component.view';
import { ListTipStatus } from '../../../ListTipComponent/helper';
import OrganizationPicker from '../../../OrganizationPicker/component.view';
import { PageSize, DefaultPage } from '../helper';
import VisitorListBase from './component.base';
import styles from './styles.view.css';
import __ from './locale';

export default class VisitorList extends VisitorListBase {
    render() {
        const { isEnabled } = this.props;
        const { searchKey, visitorList, visitorsCount, isAddVisitor, page, listTipStatus } = this.state;
        const DataGridPager = {
            size: PageSize,
            total: visitorsCount,
            page,
            onPageChange: ({ detail: { page } }) => this.handlePageChange(page),
        }

        return (
            <div className={styles['container']}>
                <div className={styles['header']}>
                    <Button
                        icon={'add'}
                        width={'auto'}
                        disabled={!isEnabled}
                        onClick={() => { this.openVisitorsDialog() }}>
                        {__('添加访问者')}
                    </Button>
                    <div className={styles['search-box']}>
                        <SearchBox
                            className={styles['search']}
                            role={'ui-searchbox'}
                            placeholder={__('请输入访问者名称')}
                            value={searchKey}
                            onChange={(value) => this.changeSearchKey(value)}
                            loader={(data) => { this.loadSearchResult(data) }}
                            disabled={!isEnabled}
                        />
                    </div>
                </div >

                <div className={classnames(styles['list'], { [styles['disabled']]: !isEnabled })}>
                    <DataGrid
                        ref={(dataGrid) => { this.dataGrid = dataGrid }}
                        data={visitorList}
                        height={'100%'}
                        showBorder={true}
                        enableSelect={true}
                        enableMultiSelect={false}
                        limit={PageSize}
                        refreshing={listTipStatus !== ListTipStatus.None}
                        DataGridPager={DataGridPager}
                        RefreshingComponent={
                            <ListTipComponent
                                listTipStatus={listTipStatus}
                            />
                        }
                        columns={[
                            {
                                title: __('访问者'),
                                key: 'name',
                                width: '45%',
                                renderCell: (name, record) => (
                                    <Text>{name}</Text>
                                ),
                            },
                            {
                                title: __('操作'),
                                key: 'operate',
                                width: '15%',
                                renderCell: (id, record) => (
                                    <UIIcon
                                        className={styles['operation']}
                                        title={__('删除')}
                                        size={16}
                                        code={'\uf046'}
                                        color={'#999'}
                                        disabled={!isEnabled}
                                        onClick={() => this.deleteVisitor(record)}
                                    />
                                ),
                            },
                        ]}
                    />
                </div>

                {
                    isAddVisitor ?
                        <OrganizationPicker
                            title={__('添加访问者')}
                            isCascadeTree={true}
                            selectType={[
                                NodeType.USER,
                                NodeType.DEPARTMENT,
                                NodeType.ORGANIZATION,
                            ]}
                            onCancel={() => this.closeVisitorsDialog()}
                            onConfirm={(visitorList) => { this.saveVisitors(visitorList) }}
                        />
                        : null
                }

            </div>
        )
    }
}