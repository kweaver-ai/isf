import { merge, isFunction } from 'lodash';
import { post, get } from '@/util/http';
import { currify } from '@/util/currify';
import { evaluate } from '@/util/accessor';
import { getHeaders, autoRefreshToken } from '../token'

/**
 * 可配置参数
 */
interface Settings {
    /**
     * 访问前缀
     */
    prefix?: string;

    CSRFToken?: () => string | string;

    onNetworkError?: Function;

    onServerError?: Function;

    onTokenError?: Function;

    onAsServerUnusable?: Function;

    getToken: () => string;

    protocol: string;

    host: string;

    port: string;
}

/**
 * 全局配置项
 */
const Config: Settings = {
    CSRFToken: undefined,

    prefix: '',

    getToken: () => '',

    protocol: 'https',

    host: '127.0.0.1',

    port: '443'

}

/**
 * toast提示是否存在
 */
let isToast: boolean = false

/**
 * 获取配置参数
 * @param option 要获取的配置参数
 */
export function getConfig(option?: string): any {
    return option === void (0) ? Config : evaluate(Config[option]);
}

/**
 * 配置Thrift
 * @param param0 配置参数
 */
export function setup(...config) {
    merge(Config, ...config)

}

const Method = {
    post,
    get,
}

interface DoRequestParams {
    /**
     * http请求函数
     */
    request: () => Promise<any>;

    /**
     * 请求地址
     */
    url: string;
}

/**
 * 发起请求
 */
export const doRequest = async ({ request, url }: DoRequestParams) => {
    return new Promise(async (resolve, reject) => {
        try {
            const { status, response } = await request()

            if (status >= 400) {
                // 并发量过大，导致服务器暂时不可用，返回503错误
                // log服务 501，底层处理，不显示服务返回的错误信息
                if (status === 503 || status === 502 || (status === 501 && (/GetLogCount/.test(url))) || status === 500) {
                    !isToast && isFunction(Config.onServerError) && Config.onServerError(() => { isToast = false });
                    isToast = true;
                    return reject(status);
                }

                // token无效或过期，刷新token
                if (status === 403) {
                    try {
                        const res = await autoRefreshToken(() => doRequest({ request, url }))
                        resolve(res)
                    } catch (ex) {
                        isFunction(Config.onTokenError) && Config.onTokenError()
                    }
                } else {
                    return reject({ ...response, status });
                }
            } else {
                return resolve(response);
            }
        } catch (ex) {
            if (!navigator.onLine) {
                isFunction(Config.onNetworkError) && Config.onNetworkError();
            } else {
                // 提示无法连接文档域
                isFunction(Config.onServerError) && Config.onServerError();
            }

            return reject(ex);
        }
    })
}

/**
 * thrift协议代理
 * @param module 模块名
 * @param method 方法名
 * @param params 参数，按顺序传递
 */
function thrift(module: string, method: string, params: Array<any> = [], { ip = '127.0.0.1', timeout = 60 * 10000, httpMethod = 'post', sendAs = 'json' } = {}, host = ''): Promise<any> {
    let abort
    const process = new Promise(async (resolve, reject) => {
        const url = `${Config.prefix}/isfweb/api${module ? `/${module}` : ''}/${method}`;
        const port = Config.port ? `:${Config.port}` : ''

        const request = () => {
            const res = Method[httpMethod](
                Config.host ? `${Config.protocol}//${Config.host}${port}${url}` : url,
                params,
                { sendAs, readAs: 'json', timeout, ...getHeaders(Config?.getToken?.(), { 'x-tclient-addr': ip }) },
            )

            // abort 已经在post方法中绑定了xhr对象
            abort = res.abort

            return res
        }

        return doRequest({ request, url }).then(resolve).catch(reject)

    }) as Promise<any> & { abort: Function }

    process.abort = abort
    return process
}

export const ShareMgnt = currify(thrift, 'ShareMgnt');
export const ShareMgntSingle = currify(thrift, 'ShareMgntSingle');
export const EACP = currify(thrift, 'EACP');
export const ConsoleInterface = currify(thrift, '');
export const PrivateAPI = currify(thrift, '');