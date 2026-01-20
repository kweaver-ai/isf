import { NodeType } from '@/core/organization';
import { TabType } from '../../type'
import { SearchFieldsItem2 } from '@/core/apis/console/auditlog/types';

export interface SearchProps {
    datasourceId: string;
    searchFields: ReadonlyArray<SearchFieldsItem2>;
    disabled: boolean;
    onRequestSearch: (conditions: Record<string, any>, fieldValueNames: Record<string, string>) => void;
    onRequestReset: (conditions: Record<string, any>) => void;
}

export const DefaultRangeStartTime = new Date('1970.1.1');

export const DefaultOrgPickerInfo = {
    show: false,
    field: '',
    selectType: 1,
    isMultiple: false,
};

export enum SelectType {
    /**
     * 可选人员
     */
    U = 1,

    /**
     * 可选部门
     */
    D = 2,

    /**
     * 可选人员和部门
     */
    UD = 3,

    /**
     * 可选用户组
     */
    Ug = 4,

    /**
     * 可选人员和用户组
     */
    UUg = 5,

    /**
     * 可选部门和用户组
     */
    DUg = 6,

    /**
     * 可选人员、部门和用户组
     */
    UDUg = 7,
}

export const OrgTypeMap = {
    '1': {
        tabType: [
            TabType.Org,
        ],
        selectType: [
            NodeType.USER,
        ],
    },
    '2': {
        tabType: [
            TabType.Org,
        ],
        selectType: [
            NodeType.DEPARTMENT,
            NodeType.ORGANIZATION,
        ],
    },
    '3': {
        tabType: [
            TabType.Org,
        ],
        selectType: [
            NodeType.USER,
            NodeType.DEPARTMENT,
            NodeType.ORGANIZATION,
        ],
    },
    '4': {
        tabType: [
            TabType.Group,
        ],
        selectType: [],
    },
    '5': {
        tabType: [
            TabType.Org,
            TabType.Group,
        ],
        selectType: [
            NodeType.USER,
        ],
    },
    '6': {
        tabType: [
            TabType.Org,
            TabType.Group,
        ],
        selectType: [
            NodeType.DEPARTMENT,
            NodeType.ORGANIZATION,
        ],
    },
    '7': {
        tabType: [
            TabType.Org,
            TabType.Group,
        ],
        selectType: [
            NodeType.USER,
            NodeType.DEPARTMENT,
            NodeType.ORGANIZATION,
        ],
    },
}