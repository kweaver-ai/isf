import * as React from 'react';
import { useState, useEffect, useCallback, useRef, forwardRef, useImperativeHandle } from 'react';
import * as moment from 'moment';
import ListTipComponent from "../../ListTipComponent/component.view"
import {
    ListTipStatus,
    ListTipMessage,
    getTipStatus
} from '../../ListTipComponent/helper'
import { isEqual } from 'lodash';
import { Text, UIIcon } from '@/ui/ui.desktop';
import {
    getDataReportDataList,
} from '@/core/apis/console/auditlog/index';
import {
    SearchFieldsItem2,
    DefaultSortDirection,
    ShowType,
} from '@/core/apis/console/auditlog/types';
import Search from './Search';
import {
    DataReportProps,
    DataReportRef,
    HistoryLog,
    HistoryLogs,
    Limit,
} from './type';
import { formatConditions, historyOperationConfig, operationConfig, managementConfig, historyManagementConfig, loginConfig, historyLoginConfig, verifyRequiredFields } from './method';
import { BigNumber } from 'bignumber.js'
import  styles from './styles.view.css';
import __ from './locale';
import { PublicErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { Table } from 'antd';
import EmptyIcon from "../../../icons/empty.png";
import SearchEmptyIcon from "../../../icons/searchEmpty.png"
import loadFailedIcon from "../../../icons/loadFailed.png";
import intl from 'react-intl-universal';
import ExportLog from '../ExportLog';

const DataReport: React.ForwardRefRenderFunction<DataReportRef, DataReportProps> = ({
    dataReportInfo,
}, ref) => {
    useImperativeHandle(ref, () => ({
        reloadPage: initializePage,
    }))

    const { id, name } = dataReportInfo;

    const searchConditionsJson = useRef<string>('');
    const dataSourceConfig = useRef<any>(null);
    const stopLoad = useRef<boolean>(false);
    const selectFieldValueNames = useRef<Record<string, string>>({});

    const [datasourceId, setDatasourceId] = useState<string>('operation');
    const [searchFields, setSearchFields] = useState<ReadonlyArray<SearchFieldsItem2>>([]);
    const [sortInfo, setSortInfo] = useState<{
        sortField: string;
        sortDirection: DefaultSortDirection;
    }>({
        sortField: 'default',
        sortDirection: DefaultSortDirection.Desc,
    });
    const [columns, setColumns] = useState<ReadonlyArray<Record<string, any>>>([]);
    const [dataList, setDataList] = useState<ReadonlyArray<Record<string, string | number>>>([]);
    const [searchConditions, setSearchConditions] = useState<Record<string, any>>({});
    const [listTipStatus, setListTipStatus] = useState<ListTipStatus>(ListTipStatus.None);
    const [initialPageListTipStatus, setInitialPageListTipStatus] = useState<ListTipStatus>(ListTipStatus.None);
    const [searchConditionVerified, setSearchConditionVerified] = useState<boolean>(false);
    const [showExportLogDialog, setShowExportLogDialog] = useState<boolean>(false);
    const [targetLogId, setTargetLogId] = useState<string>('')
    const [permissionDeniedTip, setPermissionDeniedTip] = useState<string>('');
    const [pageSize, setPageSize] = useState<number>(Limit);
    const [total, setTotal] = useState<number>(0);
    const [curPage, setCurPage] = useState<number>(1);
    const [tableKey, setTableKey] = useState(0);
    const [defaultSort, setDefaultSort] = useState<{
        sortField: string;
        sortDirection: DefaultSortDirection;
    }>(null);

    const validate = useCallback((): boolean => {
        const requiredSearchFields = searchFields.filter(({ search_is_required }) => search_is_required);

        return verifyRequiredFields(searchConditions, requiredSearchFields);
    }, [searchFields, searchConditions]);

    const fromatId = (id: string) => {
        switch(id) {
            case 'operation':
            case 'history-operation':
                return 'operation'
            case 'management':
            case 'history-management':
                return 'management'
            case 'login':
            case 'history-login':
                return 'login'
        }
    }

    const updateDataList = useCallback(async (offset = 0, limit = Limit) => {
        const searchConditionVerified = validate();
        setSearchConditionVerified(searchConditionVerified);

        if (searchConditionVerified && sortInfo && datasourceId) {
            try {
                setListTipStatus(ListTipStatus.LightLoading);

                let params: any = {
                    id: datasourceId,
                    offset,
                    limit,
                }

                const { sortField, sortDirection } = sortInfo;
                const { unique_incremental_field: specifyField } = dataSourceConfig.current;

                let order_by;

                if (sortField) {
                    order_by = [
                        {
                            field: sortField,
                            direction: sortDirection,
                            last_field_value: '',
                        },
                    ]
                }

                if (specifyField) {
                    order_by = [
                        ...order_by,
                        {
                            field: specifyField,
                            direction: DefaultSortDirection.Desc,
                            last_field_value: '',
                        },
                    ]
                }

                if (order_by) {
                    params = {
                        ...params,
                        order_by,
                    }
                }

                const condition = formatConditions(searchConditions, searchFields);

                if (condition) {
                    params = {
                        ...params,
                        condition,
                    }
                }

                searchConditionsJson.current = JSON.stringify(params);
                
                const id = fromatId(datasourceId)
                const { entries, total_count } = await getDataReportDataList({ ...params, id }, datasourceId.includes('history'));
                setTotal(total_count);

                stopLoad.current = entries.length < Limit;

                setDataList(entries);

                setListTipStatus(getTipStatus(entries, condition ? 'hasSearchKey' : ''));
            } catch (error) {
                setDataList([]);
                setListTipStatus(ListTipStatus.LoadFailed);
                if (error.code === PublicErrorCode.Forbidden && error.description) {
                    setPermissionDeniedTip(error.description)
                }
            }
        } else {
            setListTipStatus(ListTipStatus.Empty);
            setDataList([]);
        }
    }, [datasourceId, dataList, validate, sortInfo, searchConditions, searchFields]);

    const handleSort = useCallback((sort: { key: string; type: string }) => {
        setSortInfo({
            sortField: sort.key,
            sortDirection: sort.type
        });
        setCurPage(1)
        setDataList([]);
        stopLoad.current = false;
    }, []);

    const updateSearchConditions = useCallback((conditions: Record<string, any>, fieldValueNames: Record<string, string> = {}) => {
        setCurPage(1)
        setDataList([]);
        if (!isEqual(conditions, searchConditions)) {
            setSearchConditions(conditions);
            stopLoad.current = false;

            selectFieldValueNames.current = fieldValueNames;
        } else {
            // 再次点击搜索按钮，更新数据
            updateDataList();
        }
    }, [searchConditions, updateDataList]);

    const displayExportLogDialog = async(id: string) => {
        setShowExportLogDialog(true)
        setTargetLogId(id)
    }

    const exportComplete = (url: string) => {
        location.replace(url)
        setShowExportLogDialog(false)
    }

    const getRenderCellContent = (showType: ShowType, text: string | BigNumber, record: Record<string, any>, rcLabel: HistoryLog) => {
        switch (showType) {
            case ShowType.Text:
            case ShowType.Link:
            case ShowType.Img:
                return <div className={styles['text']} title={(text instanceof BigNumber) ? text.toString() : text}>{(text instanceof BigNumber) ? text.toString() : text}</div>
            case ShowType.Date:
            case ShowType.Minute:
            case ShowType.Time:
                return text ? moment(text as string).format(
                    showType === ShowType.Date ?
                        'YYYY-MM-DD' :
                        showType === ShowType.Minute ?
                            'YYYY-MM-DD HH:mm' :
                            'YYYY-MM-DD HH:mm:ss',
                ) : '--';
            case ShowType.Download:
                return HistoryLogs.includes(rcLabel) ? (
                    <UIIcon
                        className={styles['download-log']}
                        title={__('下载')}
                        size={16}
                        code={'\uf02a'}
                        color={'#999'}
                        onClick={(e) => {e.stopPropagation(); displayExportLogDialog(record.id)}}
                    />
                ) : (<div></div>)

            default:
                return '';
        }
    }

    const getDataReportConfig = (id) => {
        switch(id) {
            case 'operation':
                return operationConfig
            case 'history-operation':
                return historyOperationConfig
            case 'management':
                return managementConfig
            case 'history-management':
                return historyManagementConfig
            case 'login':
                return loginConfig
            case 'history-login':
                return historyLoginConfig
        }
    }

    const clearSort = () => {
        setTableKey(prev => prev + 1);
    };

    const initializePage = useCallback(async () => {
        try {
            setColumns([]);
            setSearchFields([]);
            setSearchConditions({});
            setDataList([]);
            setCurPage(1);
            setTotal(0);
            clearSort();

            setInitialPageListTipStatus(ListTipStatus.LightLoading);

            const {
                rc_datasource_id,
                show_fields,
                search_fields,
                datasource_config,
                datasource_config: {
                    default_sort_field,
                    default_sort_direction,
                },
                rc_label,
            } = getDataReportConfig(id);

            dataSourceConfig.current = datasource_config;
            setDatasourceId(rc_datasource_id);
            setSortInfo({
                sortField: default_sort_field,
                sortDirection: default_sort_direction,
            });
            setDefaultSort({
                sortField: default_sort_field,
                sortDirection: default_sort_direction,
            })
            setSearchFields(search_fields);

            if (show_fields.length) {
                setColumns(show_fields.map(({ field, field_title_custom, is_can_sort, show_type }) => {
                    let column: Record<string, any> = {
                        title: field_title_custom,
                        key: field,
                        dataIndex: field,
                        render: (text, record) => getRenderCellContent(show_type, text, record, rc_label),
                    }

                    if (is_can_sort) {
                        column = {
                            ...column,
                            sorter: true,
                        }
                    }

                    return column;
                }));
                stopLoad.current = false;
                setInitialPageListTipStatus(ListTipStatus.None);
            } else {
                setInitialPageListTipStatus(ListTipStatus.Empty);
            }
        } catch (error) {
            setInitialPageListTipStatus(ListTipStatus.LoadFailed);
        }
    }, [id]);

    useEffect(() => {
        if (id) {
            initializePage();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id]);

    useEffect(() => {
        setTotal(0);
        if (columns.length) {
            updateDataList(0, pageSize);
        }

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [columns, searchConditions, sortInfo]);

    return (
        <div className={styles['container']}>
            {
                initialPageListTipStatus === ListTipStatus.None && columns.length > 0 ? (
                    <div className={styles['content']}>
                        <div className={styles['header']}>
                            <Text className={styles['name']}><span>{name}</span></Text>
                            <Search
                                datasourceId={datasourceId}
                                searchFields={searchFields}
                                disabled={listTipStatus === ListTipStatus.LoadFailed}
                                onRequestSearch={updateSearchConditions}
                                onRequestReset={updateSearchConditions}
                            />
                        </div>
                        <div className={styles['body']}>
                            <Table
                                key={tableKey}
                                size="small"
                                tableLayout="fixed"
                                scroll={{ x: dataList.length ? id.includes("history") ? "200px": "1100px" : null, y: "calc(100vh - 242px)" }}
                                loading={listTipStatus === ListTipStatus.LightLoading}
                                columns={columns}
                                dataSource={dataList}
                                rowKey={(record) => record.id}
                                onChange={(pagination, _, sorter) => {
                                    if(pagination.current !== curPage) {
                                        setCurPage(pagination.current);
                                        updateDataList((pagination.current - 1) * pagination.pageSize, pagination.pageSize);
                                        return
                                    }

                                    if(pagination.pageSize !== pageSize) {
                                        setCurPage(1);
                                        setPageSize(pagination.pageSize);
                                        updateDataList(0, pagination.pageSize);
                                        return
                                    }
                                    if (sorter) {
                                        if(sorter.field) {
                                            const order = sorter.order === 'ascend' ? 'asc' : 'desc';
                                            handleSort({
                                                key: sorter.field,
                                                type: order
                                            });
                                        }else {
                                            handleSort({
                                                key: defaultSort.sortField,
                                                type: defaultSort.sortDirection,
                                            });
                                        }
                                    }         
                                }}
                                pagination={
                                    {
                                        current: curPage,
                                        pageSize,
                                        total,
                                        showSizeChanger: true,
                                        showQuickJumper: false,
                                        showTotal: (total) => {
                                            return intl.get("list.total.tip", { total });
                                        },
                                    }
                                }
                                locale={{
                                    emptyText: (
                                        <div className={styles["empty"]}>
                                            <img
                                                src={listTipStatus === ListTipStatus.LoadFailed ? loadFailedIcon : listTipStatus === ListTipStatus.NoSearchResults ? SearchEmptyIcon : EmptyIcon }
                                                alt=""
                                                width={128}
                                                height={128}
                                            />
                                            <span>
                                                {
                                                    listTipStatus === ListTipStatus.LoadFailed 
                                                        ? !permissionDeniedTip ?
                                                            ListTipMessage[ListTipStatus.LoadFailed] :
                                                            permissionDeniedTip 
                                                        : 
                                                        searchConditionVerified ? 
                                                            listTipStatus === ListTipStatus.NoSearchResults ? 
                                                                intl.get("no.search.result") :
                                                                ListTipMessage[ListTipStatus.Empty] 
                                                            : (
                                                                __('当前报表已设置必填搜索项，请先填写搜索条件')
                                                            )
                                                }
                                            </span>
                                        </div>
                                    ),
                                }}
                            />
                        </div>
                        {
                            showExportLogDialog ? (
                                <ExportLog 
                                    id={targetLogId} 
                                    onExportComplete={(url) => exportComplete(url)}
                                    onRequestCancel={() => setShowExportLogDialog(false)}
                                />)
                                : null
                        }
                    </div>
                ) : (
                    <ListTipComponent
                        listTipStatus={initialPageListTipStatus}
                        listTipMessage={{
                            ...ListTipMessage,
                            [ListTipStatus.Empty]: (
                                <div className={styles['empty']}>
                                    {__('当前报表对应的数据源无数据，请添加数据')}
                                </div>
                            ),
                        }}
                    />
                )
            }
        </div>
    )
}

export default React.memo(forwardRef(DataReport));