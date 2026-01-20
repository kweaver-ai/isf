import { reduce, isEmpty, isNumber, noop, map, filter, forEach, isArray } from 'lodash';
import * as JSONbig from 'json-bigint';
import { safeEvalPrimitive } from '../accessor';
import { isBrowser, Browser } from '../browser';
import { currify } from '../currify';
import { HTTPConfig } from './types'

declare const XDomainRequest

const Config: HTTPConfig = {
    protocol: '',
    host: '',
    port: '',
    prefix: '',
    getToken: () => '',
}

export const setupOpenApi = (config: Record<string, string | (() => void) | number>): void => {
    for (const key in config) {
        Config[key] = config[key]
    }
}

// ------ XMLHttpRequest.onload ------
// monkey patch XMLHttpRequest to make IE8 call onload when readyState === 4
if (navigator.userAgent.match(/MSIE\s8/)) {
    const sendFn = XMLHttpRequest.prototype.send;
    XMLHttpRequest.prototype.send = function () {
        // only if onreadystatechange has not already been set
        // to avoid breaking anything outside of angular
        if (!this.onreadystatechange) {
            this.onreadystatechange = function () {
                if (this.readyState === 4 && this.onload) {
                    this.onload();
                }
            };
        }
        // apply this & args to original send
        sendFn.apply(this, arguments);
    };
}

/**
 * 通过对象构建queryString
 * ```
 * queryString({id: 'aaa', keyword: ['w', 'b', 'c']})  // id=aaa&keyword=w&keyword=b&keyword=c
 * ```
 */
export function queryString(data: Record<string, any>): string {
    if (!data) {
        return '';
    }

    return map(
        filter(
            Object.entries(data),
            ([, value]) => value !== undefined && value !== null,
        ),
        ([key, value]) => (
            isArray(value) && value.length > 0 ?
                value.map((item) => [key, encodeURIComponent(item)].join('=')).join('&')
                : [key, encodeURIComponent(value)].join('=')
        ),
    ).join('&')
}

/**
 * 连接URL和参数
 */
export function joinURL(url: string, query?: Record<string, any> | string) {
    if (isEmpty(query)) {
        return url;
    }

    const [base, qs = ''] = url.split('?');
    const args = typeof query === 'object' ? queryString(query) : query;
    const fullQuery = qs ? [qs, args].join('&') : args;

    return fullQuery ? [base, fullQuery].join('?') : base;
}

/**
 * http工厂函数
 */
function http(method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'HEAD' = 'POST', url: string, data: Record<string, any> | ReadonlyArray<Record<string, any>>, { sendAs = 'form', readAs = 'form', headers = {}, beforeSend = noop, timeout = 1000 * 60 } = { timeout: 1000 * 60 }): PromiseLike<any> {

    let body;

    // 处理data，GET/HEAD请求无法在请求体中传递参数，因此拼接到URL上
    if (method === 'GET' || method === 'HEAD') {
        url = joinURL(url, data)
    }

    // 非GET/HEAD请求，根据sendAs，处理请求参数
    else {
        switch (sendAs) {
            case 'json':
                body = JSONbig.stringify(data);
                break;

            case 'text':
                body = String(data);
                break;

            case 'form':
            default:
                body = queryString(data);
                break;
        }
    }

    const useXDomain = (isBrowser({ app: Browser.MSIE, version: 8 }) || isBrowser({ app: Browser.MSIE, version: 9 })) && /^(https|http):\/\//.test(url) && url.indexOf(location.hostname) === -1;

    const xhr = useXDomain ? new XDomainRequest() : new XMLHttpRequest();

    // xhr.status置为1
    xhr.open(method, url);

    // 必须设置timeout，否则从休眠状态恢复Chrome会触发异常"Network IO Suspended"
    xhr.timeout = isNumber(timeout) ? timeout : 1000 * 60;

    // 设置请求头
    // 必须在open()之后执行
    if (xhr.setRequestHeader) {
        switch (sendAs) {
            case 'text':
                xhr.setRequestHeader('Content-Type', 'text/plain;charset=utf-8'); // 使用简单请求处理CORS
                break;

            case 'json':
                xhr.setRequestHeader('Content-Type', 'application/json;charset=utf-8'); // 使用简单请求处理CORS
                break;

            case 'form':
            default:
                xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8'); // 使用简单请求处理CORS
                break;
        }

        // 设定自定义头
        forEach(headers, (value, name) => value && xhr.setRequestHeader(name, value));
    }

    const HTTPPromise: any = new Promise((resolve, reject) => {

        xhr.onload = function () {
            let response;

            if (readAs === 'json') {
                try {
                    response = JSONbig.parse(this.responseText || null);
                } catch {
                    response = this.responseText || null;
                }
            } else if (readAs === 'xml') {
                response = this.responseXML
            } else {
                response = this.responseText
            }

            /**
             * ie8 xhr.getResponseHeader 没有 bind
             */
            resolve({ status: xhr.status, response, getResponseHeader: (name) => xhr.getResponseHeader(name) });
        }

        xhr.onerror = function () {
            reject(xhr);
        }
        xhr.ontimeout = function () {
            reject({ ...xhr, isTimeout: true })
        }

        // 增加beforeSend允许发送请求前对XHR对象做最后的配置
        beforeSend(xhr);

        body ? xhr.send(body) : xhr.send();

        xhr.onreadystatechange = function () {
            // 服务不可用
            if (xhr.status === 502) {
                reject(xhr)
            }
        }
    });

    // 允许HTTP请求中断
    // xhr.abort在IE8下没有bind方法，因此需要构造函数而不是直接xhr.abort.bind(xhr)
    HTTPPromise.abort = () => xhr.abort()

    return HTTPPromise;
}

/**
 * 将querystring中的值转换为JS数据格式
 * @param query 请求对象
 * @return object
 * @example
 * ```
 * evalQuery({k1: 'true', k2: '1', k3: 'foo'}) // { k1: true, k2: 1, k3: 'foo' }
 * ```
 */
export function evalQuery(query: { [key: string]: any }): { [key: string]: any } {

    return reduce(query, (ret, value, key) => {
        try {
            return {
                ...ret,
                [key]: safeEvalPrimitive(value),
            }
        } catch (ex) {
            return {
                ...ret,
                [key]: value,
            }
        }
    }, {})
}

// GET
export const get = currify(http, 'GET');

// POST
export const post = currify(http, 'POST');

// PUT
export const put = currify(http, 'PUT');

// DELETE
export const del = currify(http, 'DELETE');

// PATCH
export const patch = currify(http, 'PATCH');

// HEAD
export const head = currify(http, 'HEAD');