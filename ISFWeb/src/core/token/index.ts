import { getOpenAPIConfig } from '../openapiconsole';

export const getHeaders = (token, options = {}) => {
    return { headers: { Authorization: 'Bearer ' + token, ...options } }
}

export const concatWithToken = (prefix, token, url: string): string => {
    // 如果 accessPrefix 不为空且 url 不包含 accessPrefix，则拼接上当前的 accessPrefix

    return `${prefix}${url}&token=${token}`
}

let requests: ReadonlyArray<() => void> = []

/**
 * 请求自动刷新token
 */
export const autoRefreshToken = async (doRequest: () => Promise<any>) => {
    return new Promise(async (resolve, reject) => {
        try {
            const refreshToken = getOpenAPIConfig('refreshToken')
            await refreshToken()
            const allRequests = [() => resolve(doRequest()), ...requests]
            allRequests.forEach((cb) => cb())
            requests = []
        } catch (err) {
            requests = []
            reject(err)
        }
    })
}