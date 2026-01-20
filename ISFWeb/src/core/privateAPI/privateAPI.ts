import { PrivateAPI } from '../thrift'
import __ from './locale'

/**
 * privateAPI接口代理
 * @param module 模块名
 * @param method 方法名
 * @param params 参数，按顺序传递
 */
export function privateAPI(service: string, module: string, httpMethod: string, method: string, params: any, { ip = '127.0.0.1', timeout = 60 * 1000 } = {}): Promise<any> {
    let url = method ? `${service}/${module}/${method}` : `${service}/${module}`

    return PrivateAPI(url, params, { ip, timeout, httpMethod })
}