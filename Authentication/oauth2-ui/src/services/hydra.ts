/* eslint-disable @typescript-eslint/no-var-requires */
import querystring from "query-string";
import { hydraPrivateApi, hydraPrivateApiBasePrefix } from "../core/config";

/**
 * hydra服务中认证、授权、退出请求前缀
 */
const hydraUrlPrefix = hydraPrivateApiBasePrefix + "/oauth2/auth/requests/";

/**
 * @param flow - 请求类型（login consent logout）
 * @param challenge - 认证或授权唯一标识
 */
function get(flow: string, challenge: string) {
    const url = new URL(hydraUrlPrefix + flow, hydraPrivateApi.defaults.baseURL);
    url.search = querystring.stringify({ [flow + "_challenge"]: challenge });
    console.log(`[${Date()}] [INFO]  {${url.toString()} GET} START`);
    return hydraPrivateApi
        .request({
            url: url.toString(),
            method: "GET",
        })
        .then((res) => {
            console.log(`[${Date()}] [INFO]  {${url.toString()} GET} SUCCESS`);
            // 请求成功
            const { data } = res;
            return data;
        })
        .catch((err) => {
            console.error(
                `[${Date()}] [ERROR]  {${url.toString()} GET} ERROR ${JSON.stringify(
                    err && err.response && err.response.data
                )}`
            );
            // 请求失败
            return Promise.reject(err);
        });
}

/**
 * @param flow - 请求类型（login consent logout）
 * @param action - 请求行为（accept reject）
 * @param challenge - 认证或授权唯一标识
 * @param body - 请求体
 */
function put(flow: string, action: string, challenge: string, body: any) {
    const url = new URL(hydraUrlPrefix + flow + "/" + action, hydraPrivateApi.defaults.baseURL);
    url.search = querystring.stringify({ [flow + "_challenge"]: challenge });

    console.log(`[${Date()}] [INFO]  {${url.toString()} GET} START`);
    return hydraPrivateApi
        .request({
            url: url.toString(),
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            data: JSON.stringify(body),
        })
        .then(function (res) {
            console.log(`[${Date()}] [INFO]  {${url.toString()} GET} SUCCESS`);
            // 请求成功
            const { data } = res;
            return data;
        })
        .catch((err) => {
            console.error(
                `[${Date()}] [ERROR]  {${url.toString()} GET} ERROR ${JSON.stringify(
                    err && err.response && err.response.data
                )}`
            );
            // 请求失败
            return Promise.reject(err);
        });
}

/**
 * 认证、授权、退出服务接入hydra的请求
 */

// 获取login请求
export function getLoginRequest(challenge: string) {
    return get("login", challenge);
}

// 接受login请求
export function acceptLoginRequest(challenge: string, body: any) {
    return put("login", "accept", challenge, body);
}

// 拒绝login请求
export function rejectLoginRequest(challenge: string, body: any) {
    return put("login", "reject", challenge, body);
}

// 获取consent请求
export function getConsentRequest(challenge: string) {
    return get("consent", challenge);
}

// 接受consent请求
export function acceptConsentRequest(challenge: string, body: any) {
    return put("consent", "accept", challenge, body);
}

// 拒绝consent请求
export function rejectConsentRequest(challenge: string, body: any) {
    return put("consent", "reject", challenge, body);
}

// 获取logout请求
export function getLogoutRequest(challenge: string) {
    return get("logout", challenge);
}

// 接受logout请求
export function acceptLogoutRequest(challenge: string) {
    return put("logout", "accept", challenge, {});
}

// 拒绝logout请求
export function rejectLogoutRequest(challenge: string) {
    return put("logout", "reject", challenge, {});
}
