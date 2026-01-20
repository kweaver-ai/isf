import { isEqual, isFunction, isString, isNumber, isUndefined, isArray, pick, assign, without } from 'lodash';
import { useHTTPS as browserUseHTTPS } from '@/util/browser';
import { access, evaluate, isNil } from '@/util/accessor';
import { post, joinURL } from '@/util/http';

/**
 * OpenAPI 配置对象
 */
type Config = {
    /**
     * 访问地址
     */
    host: string;

    /**
     * 访问前缀
     */
    prefix: string;

    port: string | number;

    /**
     * userid 或获取userid的函数
     */
    userid?: string | Function;

    /**
     * tokenid 或获取tokenid的函数
     */
    tokenid?: string | Function;

    /**
     * token超时退出时执行
     */
    onTokenExpire?: Function;

    /**
     * 用户被禁用
     */
    onUserDisabled?: Function;

    /**
     * 网络异常时执行
     */
    onNetworkError?: Function;

    /**
     * 管理员禁止设备登录时执行
     */
    onDeviceForbid?: Function;

    /**
     * 管理员禁止身份证登录时执行
     */
    onIDAuthDisabled?: Function;

    /**
     * 设备绑定
     */
    onDeviceBind?: Function;

    /**
     * 另一地点登录
     */
    onForeignLogin?: Function;

    /**
     * Ip网段限制登录
     */
    onIPSegment?: Function;

    /**
     * 网络切换
     */
    onNetWorkChange?: Function;

    /**
     * 服务器资源不足，无法访问
     */
    onResourcesNotEnough?: Function;

    /**
     * 用户未实名认证
     */
    onUserRealNameRequired?: Function;

    /**
     * 创建者未实名认证
     */
    onCreatorRealNameRequired?: Function;
}
interface OpenAPIOptions {
    /**
     * 请求头
     */
    headers?: {
        [key: string]: any;
    };

    /**
     * 'cors'或'no-cors'
     */
    mode?: string;

    /**
     * 如何读取响应体，‘json'或'blob'
     */
    readAs?: string;

    /**
     * userid,默认从sessionStorage中获取
     */
    userid?: (() => string) | string | null;

    /**
     * tokenid,默认从sessionStorage中获取
     */
    tokenid?: (() => string) | string | null;
}

/**
 * 开放API调用函数
 */
type OpenAPI = (
    /**
     * 资源
     */
    resource: string,

    /**
     * 方法
     */
    method?: string,

    /**
     * 请求体
     */
    params?: any,

    /**
     * 配置项
     */
    Options?,
) => Promise<any>

/**
 * 使用缓存配置
 */
type UseCacheOption = {
    /**
     * 调用OpenAPI时的缓存控制，0 表示总是从服务器获取，Infinity表示总是从缓存获取
     */
    useCache: number;
}

type FactoryOptions = {
    /**
     * 缓存超时时间，超时的缓存将被清除
     */
    expires: number;
}

/**
 * 可缓存的开放API工厂函数
 */
interface CacheableOpenAPIFactory {
    /**
     * 开放API调用函数
     */
    (fn: OpenAPI, source: string, method: string, options?: FactoryOptions): (params: object, options: OpenAPIOptions & UseCacheOption) => Promise<any>;
}

/**
 * OpenAPI工厂函数
 */
type OpenAPIFactory = (module: 'eacp') => OpenAPI

/**
 * OpenAPI配置
 */
const Config: Config = {
    host: `${location.protocol}//${location.hostname}`,

    port: location.port,

    prefix: '',
}

/**
 * 开放API工厂函数
 * @param port 端口号
 */
const OpenAPIFactory: OpenAPIFactory = function (module) {
    // TODO:默认请求超时时间由前端限制为1小时，后续依赖后端优化后修改
    return function (resource, method, body, { readAs = 'json', tokenid, sign, timeout = 60 * 60 * 1000 } = {}) {
        let abort
        const process = new Promise(async (resolve, reject) => {
            // 只判断userid、tokenid为非undefined的情况，以允许传递null避免在外链访问时拼接了userid和tokenid
            // userid = evaluate(userid !== undefined ? userid : Config.userid);
            tokenid = evaluate(tokenid !== undefined ? tokenid : Config.tokenid);

            const url = joinURL(`${Config.host}${Config.port ? `:${Config.port}` : ''}${Config.prefix}/api/${module}/v1/${resource}/${method}`, { tokenid, sign });

            try {
                const request = post(url, JSON.stringify(body), { readAs, sendAs: 'text', timeout })

                // abort 已经在post方法中绑定了xhr对象
                abort = request.abort

                const { status, response } = await request

                if (status >= 400) {
                    // 会话超时
                    if (response) {
                        switch (response.errcode) {
                            case 401001:
                                isFunction(Config.onTokenExpire) && Config.onTokenExpire();
                                break;

                            case 401004:
                                isFunction(Config.onUserDisabled) && Config.onUserDisabled();
                                break;

                            case 401011:
                                isFunction(Config.onDeviceBind) && Config.onDeviceBind();
                                break;

                            case 401025:
                                isFunction(Config.onForeignLogin) && Config.onForeignLogin();
                                break;

                            case 401031:
                                isFunction(Config.onIPSegment) && Config.onIPSegment();
                                break;

                            case 401033:
                                isFunction(Config.onDeviceForbid) && Config.onDeviceForbid();
                                break;

                            case 401036:
                                isFunction(Config.onNetWorkChange) && Config.onNetWorkChange();
                                break;

                            case 401051:
                                isFunction(Config.onIDAuthDisabled) && Config.onIDAuthDisabled();
                                break;

                            case 403179:
                                isFunction(Config.onUserRealNameRequired) && Config.onUserRealNameRequired();
                                break;

                            case 403180:
                                isFunction(Config.onCreatorRealNameRequired) && Config.onCreatorRealNameRequired();
                                break;

                            case 404027:
                                isFunction(Config.onResourcesNotEnough) && Config.onResourcesNotEnough();
                                break;
                        }
                    }

                    return reject(response);
                }
                else {
                    return resolve(response);
                }
            }
            catch (ex) {
                isFunction(Config.onNetworkError) && Config.onNetworkError();

                return reject(ex);
            }
        }) as Promise<any> & { abort: Function }

        process.abort = abort

        return process
    }
}

/**
 * 设置参数
 */
export function setup(...config) {
    access(Config, ...config)
}

/**
 * 获取OpenAPI配置
 * @param options 要获取的配置项，为空返回全部配置，传递字符串返回单个配置，传递数组返回配置对象
 * @description 如果options传递userid／tokenid，则返回evaluate的结果，其他配置项返回原值
 */
export function getOpenAPIConfig(options?: string | Array<string>): Config | any {
    // const evaluateSpec = (value: any, key: string): any => /^(tokenid)$/.test(key) ? evaluate(value) : value;
    const evaluateSpec = (value: any, key: string): any => /^(userid|tokenid)$/.test(key) ? evaluate(value) : value;

    // 不传递参数，返回所有配置
    if (!options) {
        return assign({}, Config, (_objectValue, sourceValue, key) => evaluateSpec(sourceValue, key));
    }
    // 传递数组返回配置对象
    else if (isArray(options)) {
        return assign({}, pick(Config, options), (_objectValue, sourceValue, key) => evaluateSpec(sourceValue, key))
    }
    // options 传递单个字符串，则返回单个配置
    else if (isString(options)) {
        return evaluateSpec(Config[options], options);
    }
}

/**
 * 可缓存的开放API工厂函数
 * @param mod 开放API模块，eachttp或efshttp
 * @param source 开放API资源，如：dir
 * @param method 开放API方法，如：list
 * @param options
 * @param options.expires 缓存有效期，超出有效期缓存将被删除
 * @example
 *
 * ```ts
 * // 定义一个缓存有效期为 60s 的开放API接口
 * const list = CacheableOpenAPIFactory(efshttp, 'dir', 'list', {expires: 60 * 1000})
 *
 * // 调用一次开放API
 * list({docid: 'gns://B17A113F14124939BE25D07BB9B0B41E'})
 *
 * // 再次调用，如果时间在 60s内并且参数无变化，则直接使用缓存结果
 * list({docid: 'gns://B17A113F14124939BE25D07BB9B0B41E'})
 *
 * // 再次调用，并且只期望获取最近1s内的缓存结果，如果最近 1s 内没有调用过，则从服务器获取数据
 * setTimeout(() => {
 *  list({docid: 'gns://B17A113F14124939BE25D07BB9B0B41E'}, { useCache: 1000 })
 * }, 1000 * 10)
 * ```
 */
export const CacheableOpenAPIFactory: CacheableOpenAPIFactory = function (mod, source, method, { expires = 1000 }: FactoryOptions = { expires: 1000 }) {
    // 存放缓存的请求，分别代表[参数对象, 鉴权信息， 请求Promise，时间戳]
    type Cache = [object, object, Promise<any>, number]

    let caches: ReadonlyArray<Cache> = []

    return (params, { useCache = expires, tokenid, ...options } = { useCache: expires }) => {
        let nextAuthArgs

        if (tokenid === undefined) {
            nextAuthArgs = getOpenAPIConfig(['tokenid'])
        } else {
            nextAuthArgs = { tokenid }
        }

        const match = caches.find(([cacheParams, cacheAuthArgs]) => {
            return isEqual(cacheAuthArgs, nextAuthArgs) &&
                (
                    (isNil(params) && isNil(cacheParams)) ||
                    isEqual(params, cacheParams)
                )
        })
        // 检查是否需要更新缓存
        // 如果有命中项，检查缓存时间与useCache是否匹配，匹配则不需要重现缓存，否则需要重现缓存
        // 如果没有命中，则一定需要重新缓存
        if (match) {
            const [, , cacheResult, cacheTimestamp] = match

            if (Date.now() < cacheTimestamp + useCache) {
                return cacheResult
            }
        }
        // 如果需要重现缓存，则检查是否有旧缓存，如果有则删除旧缓存，没有则直接新增缓存
        if (match) {
            caches = without(caches, match)
        }

        const nextResult = mod(source, method, params, { ...nextAuthArgs, ...options })
        const cache: Cache = [params, nextAuthArgs, nextResult, Date.now()]

        caches = [...caches, cache]

        if (isNumber(expires) && isFinite(expires)) {
            setTimeout(() => {
                caches = without(caches, cache)
            }, expires)
        }
        return nextResult
    }
}

/**
 * EACHttp协议
 */
export const eachttp = OpenAPIFactory('eacp');

/**
 * 当前运行环境是否使用https
 */
export function useHTTPS(): boolean {
    return (Config.host || '').startsWith('https://') || browserUseHTTPS()
}