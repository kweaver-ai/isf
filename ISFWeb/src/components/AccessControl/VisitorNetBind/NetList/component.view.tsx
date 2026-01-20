import * as React from 'react';
import * as classnames from 'classnames';
import { Text, UIIcon, SearchBox } from '@/ui/ui.desktop';
import { DataGrid, Button } from '@/sweet-ui';
import ListTipComponent from '../../../ListTipComponent/component.view';
import { ListTipStatus } from '../../../ListTipComponent/helper';
import { PageSize, NetType } from '../helper';
import NetSegment from '../../NetSegment/component.view';
import { OperateType } from '../../NetSegment/helper'
import NetListBase from './component.base';
import { PubliceNet, setNetSegment, getNetSegment } from './helper';
import styles from './styles.view.css';
import __ from './locale';

export default class NetList extends NetListBase {
    render() {
        const { isEnabled } = this.props;
        const { searchKey, page, netList, isEditNet, netsCount, editingNet, listTipStatus, selection, operateType } = this.state;
        const DataGridPager = {
            size: PageSize,
            total: netsCount,
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
                        onClick={() => { this.updateNet(OperateType.Add) }}
                        >
                        {__('添加网段')}
                    </Button>
                    <div className={styles['search-box']}>
                        <SearchBox
                            className={styles['search']}
                            role={'ui-searchbox'}
                            width={210}
                            placeholder={__('请输入名称或网段')}
                            value={searchKey}
                            onChange={(value) => this.changeSearchKey(value)}
                            loader={() => { this.loadSearchResult() }}
                            disabled={!isEnabled}
                        />
                    </div>

                </div >

                <div className={classnames(styles['list'], { [styles['disabled']]: !isEnabled })}>
                    <DataGrid
                        ref={(dataGrid) => { this.dataGrid = dataGrid }}
                        data={netList}
                        height={'100%'}
                        showBorder={true}
                        enableSelect={true}
                        enableMultiSelect={false}
                        limit={PageSize}
                        refreshing={listTipStatus !== ListTipStatus.None}
                        selection={selection}
                        DataGridPager={DataGridPager}
                        RefreshingComponent={
                            <ListTipComponent
                                listTipStatus={listTipStatus}
                            />
                        }
                        onSelectionChange={this.selectedNet.bind(this)}
                        columns={[
                            {
                                title: __('名称'),
                                key: 'name',
                                width: '35%',
                                renderCell: (name, record) => (
                                    <div className={styles['item']}>
                                        <Text>
                                            {
                                                record.id === 'public-net'
                                                    ? PubliceNet.name
                                                    : (name ? name : '---')
                                            }
                                        </Text>
                                    </div>
                                ),
                            },
                            {
                                title: __('网段'),
                                key: 'net',
                                width: '50%',
                                renderCell: (net, record) => (
                                    <Text className={styles['net-segment']}>
                                        {
                                            record.id === 'public-net'
                                                ? PubliceNet.netInfo
                                                : (record.netType === NetType.Range
                                                    ? (record.originIP + '-' + record.endIP)
                                                    : (record.ip + '/' + record.mask)
                                                )
                                        }
                                    </Text>
                                ),
                            },
                            {
                                title: __('操作'),
                                key: 'operate',
                                width: '15%',
                                minWidth: 80,
                                renderCell: (id, record) => (
                                    record.id !== 'public-net' ?
                                        <div>
                                            <UIIcon
                                                className={styles['operation']}
                                                title={__('编辑')}
                                                size={16}
                                                code={'\uf085'}
                                                color={'#999'}
                                                disabled={!isEnabled}
                                                onClick={() => this.updateNet(OperateType.Edit, record)}
                                            />
                                            <UIIcon
                                                className={styles['operation']}
                                                title={__('删除')}
                                                size={16}
                                                code={'\uf046'}
                                                color={'#999'}
                                                disabled={!isEnabled}
                                                onClick={(e) => this.deleteNet(e, record)}
                                            />
                                        </div>
                                        : null
                                ),
                            },
                        ]}
                    />
                </div >
                {
                    isEditNet ?
                        <NetSegment
                            operateType={operateType}
                            netInfo={setNetSegment(editingNet)}
                            isShowNetName={true}
                            isShowTitle={true}
                            convererOut={getNetSegment}
                            onRequestConfirm={(netInfo) => this.handleRequestEditSuccess(netInfo)}
                            onRequestCancel={() => this.handleRequestEditCancel()}
                        />
                        : null
                }
            </div >
        )
    }
}