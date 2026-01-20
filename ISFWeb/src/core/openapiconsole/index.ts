import { isEqual, isFunction, isNumber, without } from 'lodash';
import { post, get, del, put, patch, joinURL } from '@/util/http';
import { access, isNil } from '@/util/accessor';
import { ErrorCode, PublicErrorCode } from '../apis/openapiconsole/errorcode';
import { getHeaders, autoRefreshToken } from '../token'

/**
 * OpenAPI 配置对象
 */
type Config = {
    /**
     * 访问地址
     */
    host: string;

    port: string | number;

    /**
     * 访问前缀
     */
    prefix: string;

    /**
     * 当前语言
     */
    lang?: string;

    /**
     * 网络异常时执行
     */
    onNetworkError?: Function;

    /**
     * AS服务器无法连接时执行
     */
    onServerError?: Function;

    /**
     * token无效或过期
     */
    onTokenError?: Function;

    /**
     * 请求超时
     */
    onTimeout?: Function;

    /**
     * 获取token
     */
    getToken: () => string;
}

const Method = {
    post,
    get,
    put,
    patch,
    'delete': del,
}

const TokenExpireArr = [
    ErrorCode.TokenExpire,
    ErrorCode.TokenEmpty,
    PublicErrorCode.Unauthorized,
]

type CacheableOpenAPIFactoryResult = (
    otherPathParams: ReadonlyArray<string>,
    body: object,
    queryParam: object,
    options: {
        readAs?: string;
        timeout?: number;
        resHeader?: boolean;
        locationPort?: boolean;
    }) => Promise<any>

/**
 * 开放API调用函数
 */
type OpenAPI = (
    httpMethod: 'post' | 'get' | 'put' | 'patch' | 'delete',
    /**
     * 资源
     * 比如 ['doc-domain', 'sync-plan']
     */
    pathParams: ReadonlyArray<string>,

    /**
     * 请求体
     */
    body: any,

    /**
     * query
     * 比如 {start: 0, limit: 20}
     */
    queryParam: {
        [key: string]: any | ReadonlyArray<string | number>;
    },

    /**
     * 配置项
     */
    options?: any,
) => Promise<any>

/**
 * OpenAPI工厂函数
 */
type OpenAPIFactory = () => OpenAPI

/**
 * OpenAPI配置
 */
const Config: Config = {
    host: `${location.protocol}//${location.hostname}`,

    port: location.port,

    prefix: '',

    getToken: () => '',

}

/**
 * 设置参数
 */
export function setup(...config) {
    access(Config, ...config)
}

interface DoRequestParams {
    /**
     * http请求函数
     */
    request: () => Promise<any>;

    /**
     * 是否需要错误处理
     */
    errHandling: boolean;

    /**
     * 请求地址
     */
    url: string;

    /**
     * 是否需要返回响应头
     */
    resHeader: boolean;
}

/**
 * 发起请求
 */
export const doRequest = ({ request, errHandling, url, resHeader }: DoRequestParams) => {
    return new Promise(async (resolve, reject) => {
        try {
            const { status, response, getResponseHeader } = await request()

            if (status >= 400) {
                if (errHandling) {
                    // 编目 503，底层不处理，在编目界面处理
                    // 登录安全策略页面第三方消息服务503时，不在底层处理
                    if (status === 503 && (
                        /metadata/.test(url) ||
                        (/thirdparty-message-plugin/.test(url) && /loginvisit/.test(location.hash))
                    )) {
                        return reject(503)
                    }

                    // 503错误，表现无法连接服务器
                    if (status === 503 || ((response && (response.code === ErrorCode.InternalError  || response.code === PublicErrorCode.InternalServerError) && !/metadata/.test(url)))) {
                        // 已对 涉密接口(confidential) 做容错处理，不用提示
                        if (!/confidential/.test(url)) {
                            isFunction(Config.onServerError) && Config.onServerError();
                        }

                        return reject(status);
                    }
                }

                // token无效或过期
                if (response && TokenExpireArr.includes(response.code)) {
                    try {
                        const res = await autoRefreshToken(() => doRequest({ request, errHandling, url, resHeader }))
                        resolve(res)
                    } catch {
                        isFunction(Config.onTokenError) && Config.onTokenError()
                    }
                } else {
                    return reject(response);
                }
            } else {
                if (resHeader) {
                    return resolve({ response, getResponseHeader });
                }

                return resolve(response);
            }
        } catch (ex) {
            if (errHandling) {
                if (!navigator.onLine) {
                    isFunction(Config.onNetworkError) && Config.onNetworkError();
                } else if (ex.status === 0 || ex.code === ErrorCode.InternalError || ex.code === PublicErrorCode.InternalServerError) {
                    // 无法连接服务器，服务器没有响应
                    isFunction(Config.onServerError) && Config.onServerError();
                } else if (ex.isTimeout) {
                    // 请求超时
                    isFunction(Config.onTimeout) && Config.onTimeout();
                } 
            }

            return reject(ex);
        }
    })
}
/**
 * 开放API工厂函数
 * @param port 端口号
 */
const OpenAPIFactory: OpenAPIFactory = function () {
    // TODO:默认请求超时时间由前端限制为1小时，后续依赖后端优化后修改
    return function (httpMethod, pathParams, body, queryParam = {}, { readAs = 'json', timeout = 60 * 1000, resHeader = false, locationPort = false, errHandling = true, token = '' } = {}) {

        let abort
        const process = new Promise(async (resolve, reject) => {
            const url = joinURL(`${Config.prefix}/api/${pathParams.join('/')}`, queryParam);
            const headers = {
                'Cache-Control': 'no-cache',
                Pragma: 'no-cache',
                'x-language': (Config.lang || 'zh-CN').replace(/-.*$/, (s) => s.toUpperCase()),
                'x-error-code': 'string',
                ...(token ? { Authorization: 'Bearer ' + token } : {}),
            }

            const request = () => {
                const res = Method[httpMethod](url, body, { readAs, sendAs: 'json', timeout, ...getHeaders(Config?.getToken?.(), headers) })
                // abort 已经在post方法中绑定了xhr对象
                abort = res.abort

                return res
            }

            return doRequest({ request, errHandling, url, resHeader }).then(resolve).catch(reject)
        }) as Promise<any> & { abort: Function }

        process.abort = abort

        return process
    }
}

/**
 * console协议
 */
export const consolehttp = OpenAPIFactory();

/**
 * 获取OpenAPI配置
 * @param options 要获取的配置项，传递字符串返回单个配置
 */
export function getOpenAPIConfig(options: string): any {
    return Config[options];
}

/**
 * 带缓存的工厂请求函数
 */
export const CacheableConsoleAPIFactory = (httpMethod: string, pathParams: ReadonlyArray<string>, expires: number = 1000): CacheableOpenAPIFactoryResult => {
    // 存放缓存的请求，分别代表[httpMethod, otherPathParmas, body, query, 请求Promise，时间戳]
    type Cache = [string, ReadonlyArray<string>, object, object, Promise<any>, number]
    let caches: ReadonlyArray<Cache> = []

    return (otherPathParams, body, queryParam = {}, { readAs = 'json', timeout = 60 * 1000, resHeader = false, locationPort = false } = {}) => {
        const match = caches.find(([cacheHttpMethod, cacheOtherParams, cacheBody, cacheQuery]) => {
            return (
                httpMethod === cacheHttpMethod
                && ((isNil(body) && isNil(cacheBody)) || isEqual(body, cacheBody))
                && ((isNil(queryParam) && isNil(cacheQuery)) || isEqual(queryParam, cacheQuery)))
                && isEqual(otherPathParams, cacheOtherParams)
        })
        // 检查是否需要更新缓存
        // 如果有命中项，检查缓存时间与useCache是否匹配，匹配则不需要重现缓存，否则需要重现缓存
        // 如果没有命中，则一定需要重新缓存
        if (match) {
            const [, , , , cacheResult, cacheTimestamp] = match

            if (Date.now() < cacheTimestamp + expires) {
                return cacheResult
            }

            caches = without(caches, match)
        }

        const nextResult = consolehttp(httpMethod, pathParams, body, queryParam, { readAs, timeout, resHeader, locationPort })
        const cache: Cache = [httpMethod, otherPathParams, body, queryParam, nextResult, Date.now()]

        caches = [...caches, cache]

        if (isNumber(expires) && isFinite(expires)) {
            setTimeout(() => {
                caches = without(caches, cache)
            }, expires)
        }

        return nextResult
    }
}
