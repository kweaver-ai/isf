/**
 * 开放API调用函数
 * @param P 请求参数
 * @param R 响应内容
 */
export type OpenAPI<P, R> = (params: P, ...others) => Promise<R>
