/**
 * 历史日志
 */
export enum HistoryLog {
    /**
     * 历史操作日志
     */
    HistoryLogOperation = 'history_log_operation',

    /**
     * 历史管理日志
     */
    HistoryLogManagement = 'history_log_management',

    /**
     * 历史访问日志
     */
    HistoryLogLogin = 'history_log_login',
}

export enum ExportStatus {
    /**
   * 开关关闭
   */
    SWITCH_CLOSE,

    /**
   * 开关开启
   */
    SWITCH_OPEN,

    /**
   * 转圈圈组件正在加载中
   */
    LOADING,
}

export enum ValidateState {
    /**
   * 验证合法
   */
    Normal,

    /**
   * 验证不合法
   */
    Diff,
}

export interface DataReportProps {
    dataReportInfo: {
        id: string;
        name: string;
    };
}

export interface DataReportRef {
    reloadPage: () => Promise<void>;
}

export const Limit = 50;

export const HistoryLogs = [HistoryLog.HistoryLogOperation, HistoryLog.HistoryLogManagement, HistoryLog.HistoryLogLogin];
