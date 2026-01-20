import { SelectionType, Selection, TabType } from '../helper';
import __ from './locale';

export interface AnonymousPickProps {
    /**
     * 是否用复选框
     */
    isMult: boolean;

    /**
     * 是否禁用
     */
    disabled: boolean;

    /**
     * 传递选中项
     */
    onRequsetSelection: (selection: ReadonlyArray<Selection>) => void;
}

export interface AnonymousPickState {
    /**
     * 复选框状态
     */
    checkStatus: boolean;

    /**
     * 选中匿名用户
     */
    selections: ReadonlyArray<Selection>;
}

/**
 * 选中项默认值
 */
export const DefaultSelection = {
    id: TabType.Anonymous,
    type: SelectionType.Anonymous,
    name: __('匿名用户'),
    origin: TabType.Anonymous,
};