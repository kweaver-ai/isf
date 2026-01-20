import { ValidateStatus } from '../../../type';

export interface SelectProps {
    selectedKeyOfProps: string | number;
    inputValueOfProps: string;
    optionsOfProps: ReadonlyArray<Option>;
    validateStatus: ValidateStatus;
    showSearch: boolean;
    allowClear?: boolean;
    loader: ({
        limit,
        offset,
        searchKey,
    }: {
        limit: number;
        offset: number;
        searchKey?: string;
    }) => Promise<{
        entries: ReadonlyArray<any>;
        total_count: number;
    }>;
    onChange: (option: Option, options: ReadonlyArray<Option>) => void;
}

export interface Option {
    value_code: string | number;
    value_name: string;
}

export enum BlurType {
    Selected,

    Other,
}