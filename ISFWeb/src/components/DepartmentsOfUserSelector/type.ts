export interface Props {
    /**
     * 用户信息
     */
    userInfo: any;

    /**
     * 宽度
     */
    width?: number;

    /**
     * 当前部门
     */
    dep: DepInfo;

    /**
     * 选中项
     */
    onSelectionChange: (dep: DepInfo) => void;
}

export interface DepInfo {
    id: string;
    name: string;
    path?: string;
    is_root?: boolean;
}