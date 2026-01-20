import * as React from 'react';
import { DataGrid, Button } from '@/sweet-ui';
import { Text, InlineButton, SearchBox } from '@/ui/ui.desktop';
import { ListTipStatus, ListTipMessage } from '../../ListTipComponent/helper';
import ListTipComponent from '../../ListTipComponent/component.view';
import { Limit } from '../helper';
import SetGroup from './SetGroup/component.view';
import GroupGridBase from './component.base';
import __ from './locale';
import styles from './styles.view';

const listTipMessage = {
    ...ListTipMessage,
    [ListTipStatus.Empty]: __('请点击左上角【新建用户组】进行添加'),
}

export default class GroupGrid extends GroupGridBase {
    render() {
        const {
            operatedGroup,
            searchKey,
            data: { groups, total, page },
            listTipStatus,
            selection,
        } = this.state

        return (
            <div className={styles['group-grid']}>
                <div className={styles['header']}>
                    <Button
                        role={'sweetui-button'}
                        theme={'oem'}
                        icon={'add'}
                        size={'auto'}
                        onClick={this.addGroup}
                    >
                        {__('新建用户组')}
                    </Button>
                    <SearchBox
                        role={'ui-searchbox'}
                        ref={(searchBox) => this.searchBox = searchBox}
                        className={styles['search-box']}
                        width={220}
                        placeholder={__('搜索用户组名称')}
                        value={searchKey}
                        onChange={this.changeSearchKey}
                        loader={this.getGroup}
                        onLoad={this.loadGroup}
                        onLoadFailed={this.loadFailed}
                    />
                </div>
                <div className={styles['grid-wrapper']}>
                    <DataGrid
                        role={'sweetui-datagrid'}
                        ref={(dataGrid) => this.dataGrid = dataGrid}
                        height={'100%'}
                        enableSelect={true}
                        selection={selection}
                        showBorder={true}
                        data={groups}
                        DataGridPager={{
                            size: Limit,
                            total,
                            page,
                            onPageChange: ({ detail: { page } }) => this.handlePageChange(page),
                        }}
                        refreshing={listTipStatus !== ListTipStatus.None}
                        RefreshingComponent={
                            <ListTipComponent
                                listTipStatus={listTipStatus}
                                listTipMessage={listTipMessage}
                            />
                        }
                        onSelectionChange={this.selectGroup}
                        columns={[
                            {
                                title: __('名称'),
                                key: 'name',
                                width: 40,
                                renderCell: (name, record) => (
                                    <Text role={'ui-text'}>{name}</Text>
                                ),
                            },
                            {
                                title: __('备注'),
                                key: 'notes',
                                width: 40,
                                renderCell: (notes, record) => (
                                    <Text role={'ui-text'}>{notes || '---'}</Text>
                                ),
                            },
                            {
                                title: __('操作'),
                                key: 'operation',
                                width: 20,
                                minWidth: 80,
                                renderCell: (id, record, index) => (
                                    <div className={styles['table-btn']}>
                                        <div className={styles['edit-btn']}>
                                            <InlineButton
                                                role={'ui-inlinebutton'}
                                                code={'\uf085'}
                                                size={24}
                                                title={__('编辑')}
                                                onClick={(event) => this.editGroup(event, record)}
                                            />
                                        </div>
                                        <InlineButton
                                            role={'ui-inlinebutton'}
                                            code={'\uf000'}
                                            size={24}
                                            title={__('删除')}
                                            onClick={(event) => this.deletegroup(event, record)}
                                        />
                                    </div>
                                ),
                            },
                        ]}
                    />
                </div>
                {
                    operatedGroup ?
                        <SetGroup
                            editGroup={operatedGroup}
                            onRequestCancel={this.cancelSetGroup}
                            onRequestSuccess={this.setSuccess}
                        />
                        : null
                }
            </div>
        )
    }
}