import axios, { AxiosInstance } from "axios";
import { Agent } from "https";
import axiosRetry from "axios-retry";
import { urlFormat } from "../utils/format";

const hydraPublicApiBase = urlFormat(process.env.HYDRA_PUBLIC_API_BASE) || "http://127.0.0.1";
const hydraPrivateApiBase =
    urlFormat(process.env.HYDRA_PRIVATE_API_BASE) || "http://127.0.0.1" + process.env.HYDRA_PRIVATE_API_BASE_PREFIX;
const deployPrivateApiBase = urlFormat(process.env.DEPLOY_WEB_SERVICE_API_BASE) || "http://127.0.0.1:18080";
const eacpPublicApiBase = urlFormat(process.env.EACP_PUBLIC_API_BASE) || "http://127.0.0.1";
const eacpPrivateApiBase = urlFormat(process.env.EACP_PRIVATE_API_BASE) || "http://127.0.0.1";
const authenticationPublicApiBase = urlFormat(process.env.AUTHENTICATION_PUBLIC_API_BASE) || "http://127.0.0.1";
const authenticationPrivateApiBase = urlFormat(process.env.AUTHENTICATION_PRIVATE_API_BASE) || "http://127.0.0.1";
const usermanagementPrivateApiBase = urlFormat(process.env.USER_MANAGEMENT_PRIVATE_API_BASE) || "http://127.0.0.1";
const deployWebServicePrivateApiBase = urlFormat(process.env.DEPLOY_WEB_SERVICE) || "http://127.0.0.1";
const timeoutLabel = "timeout";
const retryCount = 2;
const EACP_API_TIMEOUT = Number(process.env.EACP_API_TIMEOUT);
const eacpApiTimeout = isNaN(EACP_API_TIMEOUT) || !EACP_API_TIMEOUT ? 3000 : EACP_API_TIMEOUT;
const deployApiTimeout = 5000;
export const hydraPrivateApiBasePrefix = process.env.HYDRA_PRIVATE_API_BASE_PREFIX || "admin";
//isThirdRememberLogin为undefined或者enable时，第三方登录记住登录状态
export const isThirdRememberLogin = process.env.IS_THIRD_REMEMBER_LOGIN !== "disable";

export const hydraPublicApi = axios.create({
    baseURL: hydraPublicApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: 5000,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});
export const hydraPrivateApi = axios.create({
    baseURL: hydraPrivateApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: 5000,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});
export const deployPublicApi = axios.create({
    baseURL: deployPrivateApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: deployApiTimeout,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});

export const deployServicePrivateApi = axios.create({
    baseURL: deployWebServicePrivateApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: eacpApiTimeout,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});

export const eacpPublicApi = axios.create({
    baseURL: eacpPublicApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: eacpApiTimeout,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});
export const eacpPrivateApi = axios.create({
    baseURL: eacpPrivateApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: eacpApiTimeout,
});

export const authenticationPublicApi = axios.create({
    baseURL: authenticationPublicApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: 5000,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});
export const authenticationPrivateApi = axios.create({
    baseURL: authenticationPrivateApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: 5000,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});
export const usermanagementPrivateApi = axios.create({
    baseURL: usermanagementPrivateApiBase,
    headers: {
        "Content-Type": "application/json",
    },
    timeout: 5000,
    httpsAgent: new Agent({ rejectUnauthorized: false }),
});
enum Service {
    UNKNOWN = "unknown",
    EACP = "eacp",
    HYDRA = "hydra",
    AUTHENTICATION = "authentication",
    DEPLOY = "deploy",
}
export function getServiceNameFromApi(api: string = ""): Service {
    switch (true) {
        case api.includes(`/${Service.EACP}/`):
            return Service.EACP;
        case api.includes(`/${Service.HYDRA}/`):
            return Service.HYDRA;
        case api.includes(`/${Service.AUTHENTICATION}/`):
            return Service.AUTHENTICATION;
        case api.includes(`/${Service.DEPLOY}/`):
            return Service.DEPLOY;
        default:
            return Service.UNKNOWN;
    }
}
export function getErrorCodeFromService(service: Service): number {
    switch (service) {
        case Service.EACP:
            return 500041001;
        case Service.HYDRA:
            return 500041002;
        case Service.AUTHENTICATION:
            return 500041003;
        case Service.DEPLOY:
            return 500041003;
        default:
            return 500000000;
    }
}

function retry(service: AxiosInstance) {
    axiosRetry(service, {
        // 设置自动发送请求次数
        retries: retryCount,
        // 重复请求延迟（毫秒）
        retryDelay: () => 0,
        //  重置超时时间
        shouldResetTimeout: true,
        //控制是否应该重试请求的回调，true为打开自动发送请求，false为关闭自动发送请求
        retryCondition: (error) => error.message.includes(timeoutLabel),
    });
}

retry(eacpPublicApi);
retry(eacpPrivateApi);
