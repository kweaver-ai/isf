import * as moment from 'moment';
import { trim } from 'lodash';
import { SelectionType } from '../../OrgAndAccountPick/helper'
import { FieldType, SearchFieldsItem2, DatasourceConfig, DefaultSortDirection } from '@/core/apis/console/auditlog/types';
import __ from './locale';

export const formatConditions = (
    conditions: Record<string, any>,
    searchFields: ReadonlyArray<SearchFieldsItem2>,
): Record<string, string | number | [string | number, string | number] | (string | number)[]> => {
    let conditionsTemp;

    if (Object.keys(conditions).length) {
        searchFields.forEach(({ field, field_type, org_structure_field_config }) => {
            switch (field_type) {
                case FieldType.Text: {
                    const trimValue = trim(conditions[field]);

                    if (trimValue !== '') {
                        conditionsTemp = {
                            ...conditionsTemp,
                            [field]: trimValue,
                        }
                    }
                    break;
                }
                case FieldType.Select: {
                    const value = conditions[field];

                    if (value !== '') {
                        conditionsTemp = {
                            ...conditionsTemp,
                            [field]: value,
                        };
                    }
                    break;
                }
                case FieldType.TextRange: {
                    const { min, max } = conditions[field];
                    const trimMin = trim(min);
                    const trimMax = trim(max);

                    if (trimMin !== '' && trimMax !== '') {
                        conditionsTemp = {
                            ...conditionsTemp,
                            [field]: [trimMin, trimMax],
                        }
                    }
                    break;
                }
                case FieldType.DateRange: {
                    const { start, end } = conditions[field];

                    if (start !== '' && end !== '') {
                        conditionsTemp = {
                            ...conditionsTemp,
                            [field]: [
                                new Date(`${start.toLocaleDateString()} 00:00:00`).getTime(),
                                end.getTime(),
                            ],
                        }
                    }
                    break;
                }
                case FieldType.Org: {
                    const value = conditions[field];

                    if (value.length > 0) {
                        conditionsTemp = {
                            ...conditionsTemp,
                            [field]: value.reduce(({
                                user_ids,
                                dep_ids,
                                group_ids,
                            }, { id, type }) => {
                                if (type === SelectionType.User) {
                                    user_ids = [
                                        ...user_ids,
                                        id,
                                    ]
                                }

                                if (type === SelectionType.Department) {
                                    dep_ids = [
                                        ...dep_ids,
                                        id,
                                    ]
                                }

                                if (type === SelectionType.Group) {
                                    group_ids = [
                                        ...group_ids,
                                        id,
                                    ]
                                }

                                return {
                                    user_ids,
                                    dep_ids,
                                    group_ids,
                                }
                            }, {
                                user_ids: [],
                                dep_ids: [],
                                group_ids: [],
                            }),
                        }
                    }
                    break;
                }
                default:
                    break;
            }
        });
    }

    return conditionsTemp;
}

export const verifyRequiredFields = (conditions: Record<string, any>, searchFields: ReadonlyArray<SearchFieldsItem2>): boolean => {
    let pass = true;

    if (searchFields.length) {
        if (Object.keys(conditions).length) {
            for (const { field, field_type, search_is_required } of searchFields) {
                switch (field_type) {
                    case FieldType.Text:
                    case FieldType.Select:
                        if (search_is_required) {
                            pass = trim(conditions[field]) !== '';
                        }
                        break;
                    case FieldType.TextRange: {
                        const { min, max } = conditions[field];
                        const trimMin = trim(min), trimMax = trim(max);

                        if (search_is_required) {
                            pass = trimMin !== '' && trimMax !== '';
                        } else {
                            pass = trimMin === '' && trimMax === '' || trimMin !== '' && trimMax !== '';
                        }
                        break;
                    }
                    case FieldType.DateRange: {
                        const date = conditions[field];

                        if (search_is_required) {
                            pass = date?.start !== '' && date?.end !== '';
                        } else {
                            pass = date?.start === '' && date?.end === '' || date?.start !== '' && date?.end !== '';
                        }
                        break;
                    }
                    case FieldType.Org:
                        if (search_is_required) {
                            pass = conditions[field].length > 0;
                        }
                        break;
                    default:
                        break;
                }

                if (!pass) {
                    break;
                }
            }
        } else {
            pass = false;
        }
    }

    return pass;
}

export const computeConditionCount = (conditions: Record<string, any>, searchFields: ReadonlyArray<SearchFieldsItem2>): number => {
    let count = 0;

    if (searchFields.length) {
        if (Object.keys(conditions).length) {
            for (const { field, field_type } of searchFields) {
                switch (field_type) {
                    case FieldType.Text:
                    case FieldType.Select:
                        trim(conditions[field]) !== '' && count++;
                        break;
                    case FieldType.TextRange: {
                        const { min, max } = conditions[field];

                        trim(min) !== '' && trim(max) !== '' && count++;
                        break;
                    }
                    case FieldType.DateRange: {
                        const { start, end } = conditions[field];

                        start !== '' && end !== '' && count++;
                        break;
                    }
                    case FieldType.Org:
                        conditions[field].length > 0 && count++;
                        break;
                    default:
                        break;
                }
            }
        }
    }

    return count;
}

export const conditionConvertJsonOfHuman = (
    conditions: Record<string, any>,
    sortInfo: {
        sortField: string;
        sortDirection: string;
    },
    datasourceConfig: DatasourceConfig,
    searchFields: ReadonlyArray<SearchFieldsItem2>,
    selectFieldValueNames: Record<string, string>,
): string => {
    const { sortField, sortDirection } = sortInfo;
    const { unique_incremental_field: specifyField } = datasourceConfig;

    const field = sortField === specifyField ? {
        'zh-cn': '默认排序',
        'zh-tw': '預設排序',
        'en-us': 'Default order',
    } : {
        'zh-cn': `“${sortField}”排序`,
        'zh-tw': `“${sortField}”排序`,
        'en-us': `Sort by ${sortField}`,
    };

    const direction = sortDirection === DefaultSortDirection.Asc ? {
        'zh-cn': '升序',
        'zh-tw': '升序',
        'en-us': 'ascending',
    } : {
        'zh-cn': '降序',
        'zh-tw': '降序',
        'en-us': 'descending',
    };

    const orderBy = {
        'zh-cn': `${field['zh-cn']} ${direction['zh-cn']}`,
        'zh-tw': `${field['zh-tw']} ${direction['zh-tw']}`,
        'en-us': `${field['en-us']} ${direction['en-us']}`,
    };

    let condition = '';

    if (Object.keys(conditions).length) {
        searchFields.forEach(({ field, field_type, field_title_custom, org_structure_field_config }) => {
            switch (field_type) {
                case FieldType.Text:
                case FieldType.Select: {
                    const trimName = trim(selectFieldValueNames[field]);

                    if (trimName !== '') {
                        condition = `${condition}${field_title_custom || field}: ${trimName}; `
                    }
                    break;
                }
                case FieldType.TextRange: {
                    const { min, max } = conditions[field];
                    const trimMin = trim(min);
                    const trimMax = trim(max);

                    if (trimMin !== '' && trimMax !== '') {
                        condition = `${condition}${field_title_custom || field}: ${trimMin}-${trimMax}; `
                    }
                    break;
                }
                case FieldType.DateRange: {
                    const { start, end } = conditions[field];

                    if (start !== '' && end !== '') {
                        condition = `${condition}${field_title_custom || field}: ${moment(start).format('YYYY/MM/DD')}-${moment(end).format('YYYY/MM/DD')}; `
                    }
                    break;
                }
                case FieldType.Org: {
                    const value = conditions[field];

                    if (conditions[field].length > 0) {
                        const { is_multiple } = org_structure_field_config || { is_multiple: true };
                        const names = is_multiple ? value.map(({ name }) => name).join('、') : value[0].name;

                        condition = `${condition}${field_title_custom || field}: ${names}; `
                    }
                    break;
                }
                default:
                    break;
            }
        });
    }

    return JSON.stringify({ condition, orderBy });
}

export const operationConfig = {
    "rc_datasource_id": 'operation',
    "show_fields": [
        {
            "id": 37,
            "field": "level",
            "field_title_custom": __("级别"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 38,
            "field": "date",
            "field_title_custom": __("时间"),
            "is_can_sort": 1,
            "show_type": 6
        },
        {
            "id": 39,
            "field": "mac",
            "field_title_custom": __("设备地址"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 40,
            "field": "ip",
            "field_title_custom": __("IP地址"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 41,
            "field": "user_name",
            "field_title_custom": __("用户"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 42,
            "field": "user_paths",
            "field_title_custom": __("部门"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 43,
            "field": "op_type",
            "field_title_custom": __("操作"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 44,
            "field": "obj_type",
            "field_title_custom": __("操作对象"),
            "is_can_sort": 0,
            "show_type": 1
        },
        {
            "id": 45,
            "field": "obj_name",
            "field_title_custom": __("对象名称"),
            "is_can_sort": 0,
            "show_type": 1
        },
        {
            "id": 46,
            "field": "msg",
            "field_title_custom": __("日志描述"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 47,
            "field": "exmsg",
            "field_title_custom": __("附加信息"),
            "is_can_sort": 1,
            "show_type": 1
        }
    ],
    "search_fields": [
        {
            "id": 37,
            "field": "level",
            "field_title_custom": __("级别"),
            "field_type": 2,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    2
                ],
                "is_can_search_by_api": true,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 38,
            "field": "date",
            "field_title_custom": __("时间"),
            "field_type": 4,
            "search_is_required": 1,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    4
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 39,
            "field": "mac",
            "field_title_custom": __("设备地址"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 40,
            "field": "ip",
            "field_title_custom": __("IP地址"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 41,
            "field": "user_name",
            "field_title_custom": __("用户"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 42,
            "field": "user_paths",
            "field_title_custom": __("部门"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 43,
            "field": "op_type",
            "field_title_custom": __("操作"),
            "field_type": 2,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    2
                ],
                "is_can_search_by_api": true,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 44,
            "field": "msg",
            "field_title_custom": __("日志描述"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 45,
            "field": "exmsg",
            "field_title_custom": __("附加信息"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        }
    ],
    "datasource_config": {
        "default_sort_field": "date",
        "default_sort_direction": "desc",
        "unique_incremental_field": "",
        "id_field": "log_id"
    },
    "rc_label": "active_log_operation",
    "is_system_config": 1,
    "is_exportable": 1,
    "is_encrypted": 0,
    "file_type": 1
}

export const historyOperationConfig ={
    "rc_datasource_id": 'history-operation',
    "show_fields": [
        {
            "id": 32,
            "field": "name",
            "field_title_custom": __("文件名称"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 33,
            "field": "dump_date",
            "field_title_custom": __("转存时间"),
            "is_can_sort": 1,
            "show_type": 6
        },
        {
            "id": 34,
            "field": "size",
            "field_title_custom": __("大小"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 35,
            "field": "operation",
            "field_title_custom": __("操作"),
            "is_can_sort": 0,
            "show_type": 7
        }
    ],
    "search_fields": [
        {
            "id": 32,
            "field": "name",
            "field_title_custom": __("文件名称"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        }
    ],
    "datasource_config": {
        "default_sort_field": "dump_date",
        "default_sort_direction": "desc",
        "unique_incremental_field": "",
        "id_field": "id"
    },
    "rc_label": "history_log_operation",
    "is_system_config": 1,
    "is_exportable": 0,
    "is_encrypted": 0,
    "file_type": 1
}

export const managementConfig = {
    "rc_datasource_id": 'management',
    "show_fields": [
        {
            "id": 22,
            "field": "level",
            "field_title_custom": __("级别"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 23,
            "field": "date",
            "field_title_custom": __("时间"),
            "is_can_sort": 1,
            "show_type": 6
        },
        {
            "id": 24,
            "field": "mac",
            "field_title_custom": __("设备地址"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 25,
            "field": "ip",
            "field_title_custom": __("IP地址"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 26,
            "field": "user_name",
            "field_title_custom": __("用户"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 27,
            "field": "user_paths",
            "field_title_custom": __("部门"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 28,
            "field": "op_type",
            "field_title_custom": __("操作"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 29,
            "field": "obj_type",
            "field_title_custom": __("操作对象"),
            "is_can_sort": 0,
            "show_type": 1
        },
        {
            "id": 30,
            "field": "obj_name",
            "field_title_custom": __("对象名称"),
            "is_can_sort": 0,
            "show_type": 1
        },
        {
            "id": 31,
            "field": "msg",
            "field_title_custom": __("日志描述"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 32,
            "field": "exmsg",
            "field_title_custom": __("附加信息"),
            "is_can_sort": 1,
            "show_type": 1
        }
    ],
    "search_fields": [
        {
            "id": 22,
            "field": "level",
            "field_title_custom": __("级别"),
            "field_type": 2,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    2
                ],
                "is_can_search_by_api": true,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 23,
            "field": "date",
            "field_title_custom": __("时间"),
            "field_type": 4,
            "search_is_required": 1,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    4
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 24,
            "field": "mac",
            "field_title_custom": __("设备地址"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 25,
            "field": "ip",
            "field_title_custom": __("IP地址"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 26,
            "field": "user_name",
            "field_title_custom": __("用户"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 27,
            "field": "user_paths",
            "field_title_custom": __("部门"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 28,
            "field": "op_type",
            "field_title_custom": __("操作"),
            "field_type": 2,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    2
                ],
                "is_can_search_by_api": true,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 29,
            "field": "msg",
            "field_title_custom": __("日志描述"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 30,
            "field": "exmsg",
            "field_title_custom": __("附加信息"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        }
    ],
    "datasource_config": {
        "default_sort_field": "date",
        "default_sort_direction": "desc",
        "unique_incremental_field": "",
        "id_field": "log_id"
    },
    "rc_label": "active_log_management",
    "is_system_config": 1,
    "is_exportable": 1,
    "is_encrypted": 0,
    "file_type": 1
}

export const historyManagementConfig = {
    "rc_datasource_id": 'history-management',
    "show_fields": [
        {
            "id": 17,
            "field": "name",
            "field_title_custom": __("文件名称"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 18,
            "field": "dump_date",
            "field_title_custom": __("转存时间"),
            "is_can_sort": 1,
            "show_type": 6
        },
        {
            "id": 19,
            "field": "size",
            "field_title_custom": __("大小"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 20,
            "field": "operation",
            "field_title_custom": __("操作"),
            "is_can_sort": 0,
            "show_type": 7
        }
    ],
    "search_fields": [
        {
            "id": 17,
            "field": "name",
            "field_title_custom": __("文件名称"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        }
    ],
    "datasource_config": {
        "default_sort_field": "dump_date",
        "default_sort_direction": "desc",
        "unique_incremental_field": "",
        "id_field": "id"
    },
    "rc_label": "history_log_management",
    "is_system_config": 1,
    "is_exportable": 0,
    "is_encrypted": 0,
    "file_type": 1
}

export const loginConfig = {
    "rc_datasource_id": 'login',
    "show_fields": [
        {
            "id": 7,
            "field": "level",
            "field_title_custom": __("级别"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 8,
            "field": "date",
            "field_title_custom": __("时间"),
            "is_can_sort": 1,
            "show_type": 6
        },
        {
            "id": 9,
            "field": "mac",
            "field_title_custom": __("设备地址"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 10,
            "field": "ip",
            "field_title_custom": __("IP地址"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 11,
            "field": "user_name",
            "field_title_custom": __("用户"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 12,
            "field": "user_paths",
            "field_title_custom": __("部门"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 13,
            "field": "op_type",
            "field_title_custom": __("操作"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 14,
            "field": "obj_type",
            "field_title_custom": __("操作对象"),
            "is_can_sort": 0,
            "show_type": 1
        },
        {
            "id": 15,
            "field": "obj_name",
            "field_title_custom": __("对象名称"),
            "is_can_sort": 0,
            "show_type": 1
        },
        {
            "id": 16,
            "field": "msg",
            "field_title_custom": __("日志描述"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 17,
            "field": "exmsg",
            "field_title_custom": __("附加信息"),
            "is_can_sort": 1,
            "show_type": 1
        }
    ],
    "search_fields": [
        {
            "id": 7,
            "field": "level",
            "field_title_custom": __("级别"),
            "field_type": 2,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    2
                ],
                "is_can_search_by_api": true,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 8,
            "field": "date",
            "field_title_custom": __("时间"),
            "field_type": 4,
            "search_is_required": 1,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    4
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 9,
            "field": "mac",
            "field_title_custom": __("设备地址"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 10,
            "field": "ip",
            "field_title_custom": __("IP地址"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 11,
            "field": "user_name",
            "field_title_custom": __("用户"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 12,
            "field": "user_paths",
            "field_title_custom": __("部门"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 13,
            "field": "op_type",
            "field_title_custom": __("操作"),
            "field_type": 2,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    2
                ],
                "is_can_search_by_api": true,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 14,
            "field": "msg",
            "field_title_custom": __("日志描述"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        },
        {
            "id": 15,
            "field": "exmsg",
            "field_title_custom": __("附加信息"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        }
    ],
    "datasource_config": {
        "default_sort_field": "date",
        "default_sort_direction": "desc",
        "unique_incremental_field": "",
        "id_field": "log_id"
    },
    "rc_label": "active_log_login",
    "is_system_config": 1,
    "is_exportable": 1,
    "is_encrypted": 0,
    "file_type": 1
}

export const historyLoginConfig = {
    "rc_datasource_id": 'history-login',
    "show_fields": [
        {
            "id": 2,
            "field": "name",
            "field_title_custom": __("文件名称"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 3,
            "field": "dump_date",
            "field_title_custom": __("转存时间"),
            "is_can_sort": 1,
            "show_type": 6
        },
        {
            "id": 4,
            "field": "size",
            "field_title_custom": __("大小"),
            "is_can_sort": 1,
            "show_type": 1
        },
        {
            "id": 5,
            "field": "operation",
            "field_title_custom": __("操作"),
            "is_can_sort": 0,
            "show_type": 7
        }
    ],
    "search_fields": [
        {
            "id": 2,
            "field": "name",
            "field_title_custom": __("文件名称"),
            "field_type": 1,
            "search_is_required": 0,
            "search_field_config": {
                "dependent_fields": null,
                "support_types": [
                    1
                ],
                "is_can_search_by_api": false,
                "search_label": ""
            },
            "org_structure_field_config": null
        }
    ],
    "datasource_config": {
        "default_sort_field": "dump_date",
        "default_sort_direction": "desc",
        "unique_incremental_field": "",
        "id_field": "id"
    },
    "rc_label": "history_log_login",
    "is_system_config": 1,
    "is_exportable": 0,
    "is_encrypted": 0,
    "file_type": 1
}