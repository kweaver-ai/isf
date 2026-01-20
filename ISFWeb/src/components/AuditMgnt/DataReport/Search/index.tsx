import * as React from 'react';
import { useState, useEffect, useCallback, useMemo, useRef } from 'react';
import classNames from 'classnames';
import { includes } from 'lodash';
import {
    Button,
    ValidateBox,
    ValidateNumber,
    SweetIcon,
    Tag,
    Trigger,
    ValidateComboArea,
    Menu,
} from '@/sweet-ui';
import { Form, DateBox } from '@/ui/ui.desktop';
import {
    getDataReportFieldValuesList,
} from '@/core/apis/console/auditlog/index';
import { FieldType } from '@/core/apis/console/auditlog/types';
import { NodeType } from '@/core/organization';
import { Selection } from "../../../OrgAndAccountPick/helper"
import OrgAndGroupPicker from '../../../OrgAndGroupPick/component.view'
import { ValidateMessages, ValidateStatus, TabType } from '../../type';
import { verifyRequiredFields, computeConditionCount } from '../method';
import Select from './Select';
import Ellipsis from './Ellipsis';
import { Option } from './Select/type';
import {
    DefaultOrgPickerInfo,
    DefaultRangeStartTime,
    SearchProps,
    OrgTypeMap,
} from './type';
import  styles from './styles.view.css';
import __ from './locale';

const {
    Row: FormRow,
    Label: FormLabel,
    Field: FormField,
} = Form;

const Search: React.FC<SearchProps> = ({
    datasourceId,
    searchFields,
    disabled,
    onRequestSearch,
    onRequestReset,
}) => {
    const requiredSearchFields = useMemo(() => searchFields.filter(({ search_is_required }) => !!search_is_required), [searchFields]);
    const optionalSearchFields = useMemo(() => searchFields.filter(({ search_is_required }) => !search_is_required), [searchFields]);

    const curSelectableRange = useRef<{
        startRange: [Date] | [Date, Date];
        endRange: [Date] | [Date, Date];
    }>({
        startRange: [DefaultRangeStartTime],
        endRange: [DefaultRangeStartTime],
    });
    const defaultInfo = useRef<{
        condition: Record<string, any>;
        fieldValueNames: Record<string, string>;
        fieldValueOptions: Record<string, ReadonlyArray<Option>>;
        validateStatus: Record<string, ValidateStatus | Record<string, ValidateStatus>>;
        orgSelections: Record<string, ReadonlyArray<Selection>>;
    }>({
        condition: {},
        fieldValueNames: {},
        fieldValueOptions: {},
        validateStatus: {},
        orgSelections: {},
    });
    const lastInfo = useRef<{
        condition: Record<string, any>;
        fieldValueNames: Record<string, string>;
        fieldValueOptions: Record<string, ReadonlyArray<Option>>;
        validateStatus: Record<string, ValidateStatus | Record<string, ValidateStatus>>;
        orgSelections: Record<string, ReadonlyArray<Selection>>;
        selectableRange: {
            startRange: [Date] | [Date, Date];
            endRange: [Date] | [Date, Date];
        };
    }>({
        condition: {},
        fieldValueNames: {},
        fieldValueOptions: {},
        validateStatus: {},
        orgSelections: {},
        selectableRange: {
            startRange: [DefaultRangeStartTime],
            endRange: [DefaultRangeStartTime],
        },
    });
    const curFieldValueNames = useRef<Record<string, string>>({});
    const curFieldValueOptions = useRef<Record<string, ReadonlyArray<Option>>>({});
    const triggerElement = useRef(null);

    const [open, setOpen] = useState<boolean>(false);

    const [conditionCount, setConditionCount] = useState<number>(0);

    const [condition, setCondition] = useState<Record<string, any>>({});

    const [fieldValues, setFieldValues] = useState<Record<string, any>>({});

    const [validateStatus, setValidateStatus] = useState<Record<string, ValidateStatus | Record<string, ValidateStatus>>>({});

    const [orgSelections, setOrgSelections] = useState<Record<string, ReadonlyArray<Selection>>>({});

    const [disbleSearchBtn, setDisableSearchBtn] = useState<boolean>(true);

    const [{
        show,
        field,
        selectType,
        isMultiple,
    }, setOrgPickerInfo] = useState<{
        show: boolean;
        field: string;
        selectType: number;
        isMultiple: boolean;
    }>(DefaultOrgPickerInfo);

    const confirmAddOrg = useCallback((field: string, getValue: (fieldValue: any) => any) => {
        setOrgPickerInfo(DefaultOrgPickerInfo);

        updateCondition(field, getValue);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const updateCondition = useCallback((field: string, getValue: (fieldValue: any) => any) => {
        setCondition((condition) => ({
            ...condition,
            [field]: getValue(condition[field]),
        }));
        setValidateStatus((validateStatus) => ({
            ...validateStatus,
            [field]: ValidateStatus.Normal,
        }));
    }, []);

    const reset = useCallback(() => {
        const { validateStatus, condition, fieldValueNames, fieldValueOptions, orgSelections } = defaultInfo.current;

        setValidateStatus(validateStatus);
        setCondition(condition);
        setOrgSelections(orgSelections);
        setConditionCount(0);

        curSelectableRange.current = {
            startRange: [DefaultRangeStartTime],
            endRange: [DefaultRangeStartTime],
        };
        curFieldValueNames.current = fieldValueNames;
        curFieldValueOptions.current = fieldValueOptions;

        lastInfo.current = {
            condition,
            fieldValueNames,
            fieldValueOptions,
            validateStatus,
            orgSelections,
            selectableRange: {
                startRange: [DefaultRangeStartTime],
                endRange: [DefaultRangeStartTime],
            },
        }

        onRequestReset(condition);
    }, [onRequestReset]);

    const handleDeleteTag = useCallback(() => {
        setConditionCount(0);
        reset();
    }, [reset]);

    const renderField = useCallback((field_type, field, search_field_config, org_structure_field_config) => {
        switch (field_type) {
            case FieldType.Text:
                return (
                    <ValidateBox
                        value={condition[field]}
                        width={244}
                        placeholder={__('请输入')}
                        validateState={validateStatus[field]}
                        validateMessages={ValidateMessages}
                        onValueChange={({ detail }) => {
                            updateCondition(field, () => detail);
                        }}
                    />
                )
            case FieldType.Select: {
                return (
                    <div>
                        <Select
                            allowClear={true}
                            selectedKeyOfProps={condition[field]}
                            inputValueOfProps={curFieldValueNames.current[field]}
                            optionsOfProps={curFieldValueOptions.current[field]}
                            validateStatus={validateStatus[field] as ValidateStatus}
                            showSearch={search_field_config.is_can_search_by_api}
                            loader={async ({ limit, offset, searchKey }) => {
                                const { entries, total_count } = await getDataReportFieldValuesList({
                                    id: datasourceId,
                                    offset,
                                    limit,
                                    field,
                                    keyword: searchKey,
                                })

                                return {
                                    entries,
                                    total_count,
                                }
                            }}
                            onChange={({ value_code, value_name }, options) => {
                                curFieldValueNames.current = {
                                    ...curFieldValueNames.current,
                                    [field]: value_name,
                                }
                                curFieldValueOptions.current = {
                                    ...curFieldValueOptions.current,
                                    [field]: options,
                                }

                                updateCondition(field, () => value_code);
                            }}
                        />
                    </div>
                )
            }
            case FieldType.TextRange: {
                const { min, max } = condition[field];
                const fieldTipStatus = validateStatus[field] as Record<string, ValidateStatus>;

                return (
                    <div className={styles['text-range-container']}>
                        <ValidateNumber
                            role={'sweetui-validatenumber'}
                            width={112}
                            value={min}
                            precision={0}
                            placeholder={__('请输入')}
                            validateMessages={ValidateMessages}
                            validateState={fieldTipStatus.min}
                            onValueChange={({ detail }) => {
                                updateCondition(field, (fieldValue) => ({
                                    ...fieldValue,
                                    min: detail,
                                }));
                            }}
                        />
                        &nbsp; - &nbsp;
                        <ValidateNumber
                            role={'sweetui-validatenumber'}
                            width={112}
                            value={max}
                            precision={0}
                            placeholder={__('请输入')}
                            validateMessages={ValidateMessages}
                            validateState={fieldTipStatus.max}
                            onValueChange={({ detail }) => {
                                updateCondition(field, (fieldValue) => ({
                                    ...fieldValue,
                                    max: detail,
                                }));
                            }}
                        />
                    </div>
                )
            }
            case FieldType.DateRange: {
                const { start, end } = condition[field];

                return (
                    <div className={styles['date-container']}>
                        <DateBox
                            width={112}
                            placeholder={__('开始日期')}
                            shouldShowblankStatus={!start}
                            value={start}
                            element={[triggerElement.current, 'window']}
                            selectRange={curSelectableRange.current.startRange}
                            onChange={(date) => {
                                updateCondition(field, (fieldValue) => ({
                                    ...fieldValue,
                                    start: date,
                                }));

                                curSelectableRange.current = {
                                    ...curSelectableRange.current,
                                    endRange: [date],
                                }
                            }}
                        />
                        &nbsp; - &nbsp;
                        <DateBox
                            width={112}
                            placeholder={__('结束日期')}
                            shouldShowblankStatus={!end}
                            value={end}
                            element={[triggerElement.current, 'window']}
                            selectRange={curSelectableRange.current.endRange}
                            onChange={(date) => {
                                updateCondition(field, (fieldValue) => ({
                                    ...fieldValue,
                                    end: date,
                                }));

                                curSelectableRange.current = {
                                    ...curSelectableRange.current,
                                    startRange: [DefaultRangeStartTime, date],
                                }
                            }}
                        />
                    </div>
                )
            }
            case FieldType.Org:
                return (
                    <div className={styles['org-container']}>
                        <ValidateComboArea
                            width={192}
                            placeholder={__('请选择')}
                            validateMessages={ValidateMessages}
                            validateState={validateStatus[field]}
                            value={condition[field]}
                            formatter={(item) => item.name}
                            style={{ whiteSpace: 'normal' }}
                            onChange={(value) => {
                                updateCondition(field, () => value);
                                setOrgSelections((orgSelections) => ({
                                    ...orgSelections,
                                    [field]: value,
                                }));
                            }}
                        />
                        <Button
                            size='auto'
                            className={styles['org-select-btn']}
                            onClick={() => {
                                setOrgPickerInfo({
                                    show: true,
                                    field,
                                    selectType: org_structure_field_config ?
                                        org_structure_field_config.select_type : DefaultOrgPickerInfo.selectType,
                                    isMultiple: org_structure_field_config ?
                                        !!org_structure_field_config.is_multiple : DefaultOrgPickerInfo.isMultiple,
                                });
                            }}
                        >
                            {__('选择')}
                        </Button>
                    </div>
                )
            default:
                return null;
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [condition, validateStatus, updateCondition, fieldValues])

    useEffect(() => {
        const initializePage = () => {
            const condition = {};
            const fieldValueNames = {};
            const fieldValueOptions = {};
            const fieldValues = {};
            const validateStatus = {};
            const orgSelections = {};
            setConditionCount(0);

            for (let {
                field,
                field_type,
                search_field_config,
                org_structure_field_config,
            } of searchFields) {
                switch (field_type) {
                    case FieldType.Text:
                        condition[field] = '';
                        validateStatus[field] = ValidateStatus.Normal;
                        break;
                    case FieldType.Select: {
                        fieldValues[field] = [];

                        condition[field] = '';
                        fieldValueNames[field] = '';
                        fieldValueOptions[field] = [];
                        validateStatus[field] = ValidateStatus.Normal;

                        break;
                    }
                    case FieldType.TextRange:
                        condition[field] = {
                            min: '',
                            max: '',
                        };
                        validateStatus[field] = {
                            min: ValidateStatus.Normal,
                            max: ValidateStatus.Normal,
                        };
                        break;
                    case FieldType.DateRange:
                        condition[field] = {
                            start: '',
                            end: '',
                        };
                        validateStatus[field] = {
                            start: ValidateStatus.Normal,
                            end: ValidateStatus.Normal,
                        };
                        break;
                    case FieldType.Org:
                        condition[field] = [];
                        validateStatus[field] = ValidateStatus.Normal;
                        orgSelections[field] = [];
                        break;
                    default:
                        break;
                }
            }

            setFieldValues(fieldValues);
            setValidateStatus(validateStatus);
            setCondition(condition);
            setOrgSelections(orgSelections);

            curFieldValueNames.current = fieldValueNames;
            curFieldValueOptions.current = fieldValueOptions;
            curSelectableRange.current = {
                startRange: [DefaultRangeStartTime],
                endRange: [DefaultRangeStartTime],
            };

            defaultInfo.current = {
                condition,
                fieldValueNames,
                fieldValueOptions,
                validateStatus,
                orgSelections,
            };
            lastInfo.current = {
                condition,
                fieldValueNames,
                fieldValueOptions,
                validateStatus,
                orgSelections,
                selectableRange: {
                    startRange: [DefaultRangeStartTime],
                    endRange: [DefaultRangeStartTime],
                },
            };
        };

        initializePage();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [searchFields, datasourceId]);

    useEffect(() => {
        setDisableSearchBtn(!verifyRequiredFields(condition, searchFields));
    }, [condition, searchFields]);

    const orgType = OrgTypeMap[selectType];

    const tabType =
        orgType && orgType.tabType ?
            orgType.tabType :
            [TabType.Org, TabType.Group];
    const nodeType =
        orgType && orgType.selectType ?
            orgType.selectType :
            [
                NodeType.USER,
                NodeType.DEPARTMENT,
                NodeType.ORGANIZATION,
            ];

    const placeholder = (nodeType.length === 1 && nodeType[0] === NodeType.USER)
        ? __('搜索用户')
        : (
            includes(nodeType, NodeType.USER)
                ? __('搜索用户或部门')
                : __('搜索部门')
        )

    const conditionListIsNotEmpty = !!(searchFields.length && Object.keys(condition).length);
    const disabledSearch = !conditionListIsNotEmpty || disabled;

    return (
        <Trigger
            triggerEvent={'click'}
            anchorOrigin={['right', 'bottom']}
            alignOrigin={['right', -2]}
            freeze={false}
            onBeforePopupClose={() => {
                setOpen(false);

                const { validateStatus, condition, fieldValueNames, fieldValueOptions, orgSelections, selectableRange } = lastInfo.current;

                setValidateStatus(validateStatus);
                setCondition(condition);
                setOrgSelections(orgSelections);
                curSelectableRange.current = selectableRange;
                curFieldValueNames.current = fieldValueNames;
                curFieldValueOptions.current = fieldValueOptions;
            }}
            element={[triggerElement.current, 'window']}
            renderer={({ setPopupVisibleOnClick }) =>
                <div
                    ref={triggerElement}
                    key={'search-trigger'}
                    className={classNames(
                        styles['search'],
                        {
                            [styles['search-active']]: open,
                            [styles['search-disabled']]: disabledSearch,
                        },
                    )}
                    onClick={disabledSearch ? undefined : () => { setOpen(!open); setPopupVisibleOnClick() }}
                >
                    <SweetIcon
                        role={'sweetui-sweeticon'}
                        className={styles['search-icon']}
                        name={'search'}
                        size={14}
                    />
                    {
                        conditionCount ? (
                            <Tag
                                className={styles['tag']}
                                closable={true}
                                onClose={handleDeleteTag}
                            >
                                {__('${conditionCount} 个搜索项', { conditionCount })}
                            </Tag>
                        ) : (
                            <span className={styles['placeholder']}>{__('搜索')}</span>
                        )
                    }
                    <SweetIcon
                        role={'sweetui-sweeticon'}
                        className={styles['arrow']}
                        name={open ? 'arrowUp' : 'arrowDown'}
                        size={16}
                    />
                </div>
            }
        >
            {
                ({ close }) => {
                    return (
                        conditionListIsNotEmpty ? (
                            <Menu>
                                <div className={styles['container']}>
                                    <div className={styles['form-container']}>
                                        <Form>
                                            {
                                                [
                                                    ...requiredSearchFields,
                                                    ...optionalSearchFields,
                                                ].map(({
                                                    id,
                                                    field,
                                                    field_title_custom,
                                                    field_type,
                                                    search_is_required,
                                                    search_field_config,
                                                    search_field_config: { search_label },
                                                    org_structure_field_config,
                                                }) => {
                                                    return (
                                                        <FormRow key={id}>
                                                            <FormLabel
                                                                className={styles['label']}
                                                                colon
                                                                align={includes([FieldType.Org], field_type) ? 'top' : ''}
                                                            >
                                                                <Ellipsis>{search_label || field_title_custom}</Ellipsis>
                                                            </FormLabel>
                                                            <FormField isRequired={!!search_is_required}>
                                                                {
                                                                    renderField(field_type, field, search_field_config, org_structure_field_config)
                                                                }
                                                            </FormField>
                                                        </FormRow>
                                                    )
                                                })
                                            }
                                        </Form>
                                    </div>
                                    <div className={styles['footer']}>
                                        <Button
                                            className={styles['button-left']}
                                            theme={'oem'}
                                            disabled={disbleSearchBtn}
                                            onClick={() => {
                                                onRequestSearch(condition, curFieldValueNames.current);
                                                lastInfo.current = {
                                                    condition,
                                                    fieldValueNames: curFieldValueNames.current,
                                                    fieldValueOptions: curFieldValueOptions.current,
                                                    validateStatus,
                                                    orgSelections,
                                                    selectableRange: curSelectableRange.current,
                                                };
                                                setConditionCount(computeConditionCount(condition, searchFields))
                                                setOpen(false);
                                                close();
                                            }}
                                        >
                                            {__('搜索')}
                                        </Button>
                                        <Button
                                            className={styles['button-left']}
                                            onClick={() => {
                                                reset();
                                                setOpen(false);
                                                close();
                                            }}
                                        >
                                            {__('重置')}
                                        </Button>
                                    </div>
                                    {
                                        show
                                            ? (
                                                <OrgAndGroupPicker
                                                    title={__('选择')}
                                                    element={triggerElement.current}
                                                    tabType={tabType}
                                                    defaultSelections={orgSelections[field]}
                                                    isMult={false}
                                                    isSingleChoice={!isMultiple}
                                                    isShowCheckBox={false}
                                                    nodeType={nodeType}
                                                    placeholder={placeholder}
                                                    onRequestConfirm={(selections) => {
                                                        setOrgSelections((orgSelections) => ({
                                                            ...orgSelections,
                                                            [field]: selections,
                                                        }));

                                                        confirmAddOrg(field, () => selections)
                                                    }}
                                                    onRequestCancel={() => {
                                                        setOrgPickerInfo(DefaultOrgPickerInfo);

                                                        setOrgSelections((orgSelections) => ({
                                                            ...orgSelections,
                                                            [field]: condition[field],
                                                        }));
                                                    }}
                                                />
                                            )
                                            : null
                                    }
                                </div>
                            </Menu>
                        ) : null
                    )
                }
            }
        </Trigger>
    )
}

export default React.memo(Search);