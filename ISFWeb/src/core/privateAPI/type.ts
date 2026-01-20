/**
 * 获取导出文件参数headers
 */
interface Headers {
    /**
     * 身份鉴权
     */
    Authorization: string;

    /**
     * 时间
     */
    'x-amz-date': string;
}

/**
 * 获取导出文件参数
 */
export interface GetExportFile {
    headers: Headers;
    method: string;
    url: string;
}
