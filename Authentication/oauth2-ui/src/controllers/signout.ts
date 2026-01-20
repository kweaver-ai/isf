import { Request, Response } from "express";
import { acceptLogoutRequest, getLogoutRequest } from "../services/hydra";
import {
    authenticationPrivateApi,
    getServiceNameFromApi,
    getErrorCodeFromService,
    eacpPublicApi,
} from "../core/config";
import { ErrorCode } from "../core/errorcode";
import { SignoutLogin3rdPartyStatusKey, SignoutLogin3rdPartyStatus } from "../const";

const signout = async (req: Request, res: Response) => {
    let { logout_challenge } = req.query;
    logout_challenge = req.query.logout_challenge || req.cookies.logout_challenge;
    let is_previous_login_3rd_party = req.cookies.is_previous_login_3rd_party;
    let urlPrefix = "";
    try {
        const prefixPath = req.cookies["X-Forwarded-Prefix"];
        urlPrefix = !prefixPath || prefixPath === "/" ? "" : prefixPath;
    } catch {
        urlPrefix = "";
    }
    if (logout_challenge) {
        (res as Response).cookie("logout_challenge", logout_challenge);
        try {
            if (req.query.logout_challenge) {
                // 正常退出的提前准备 + 第三方退出的第一次退出的提前准备
                const lastTimestamp = Date.now();
                const { sid = "" } = await getLogoutRequest(logout_challenge as string);
                console.log(`[${Date()}] [INFO]  {/api/authentication/v1/session/${sid} DELETE} START`);
                await authenticationPrivateApi.delete(`/api/authentication/v1/session/${sid}`);
                console.log(
                    `[${Date()}] [INFO]  {/api/authentication/v1/session/${sid} DELETE} SUCCESS +${
                        Date.now() - lastTimestamp
                    }ms`
                );
            } else {
                //  第三方退出的第二次退出的提前准备
                (res as Response).clearCookie("is_previous_login_3rd_party");
                (res as Response).cookie(SignoutLogin3rdPartyStatusKey, SignoutLogin3rdPartyStatus.SecoundSignout);
                is_previous_login_3rd_party = false;
            }

            if (!is_previous_login_3rd_party) {
                // 正常退出 / 第三方退出的第二次退出
                const { redirect_to } = await acceptLogoutRequest(logout_challenge as string);
                res.redirect(redirect_to);
            } else {
                // 第三方退出的第一次退出
                const { data } = await eacpPublicApi.get("/api/eacp/v1/auth1/login-configs");
                const logoutUrl = data.thirdauth.config.logoutUrl;
                const removeCookies = data.thirdauth.config.removeCookies as string[];
                if (logoutUrl) {
                    (res as Response).cookie(SignoutLogin3rdPartyStatusKey, SignoutLogin3rdPartyStatus.firstSignout);
                    (res as Response).clearCookie("expires_in");
                    (res as Response).clearCookie("token_type");
                    (res as Response).clearCookie("oauth2_authentication_session");
                    removeCookies?.forEach((cookie) => {
                        (res as Response).clearCookie(cookie);
                    });
                    res.redirect(logoutUrl);
                    return;
                } else {
                    const { redirect_to } = await acceptLogoutRequest(logout_challenge as string);
                    res.redirect(redirect_to);
                }
            }
        } catch (e: any) {
            const path = e && e.request && (e.request.path || (e.request._options && e.request._options.path));
            console.error(`[${Date()}] [ERROR]  ${path}  ERROR ${JSON.stringify(e && e.response && e.response.data)}`);
            console.error(`[${Date()}] [ERROR]  ${path}  ERROR ${JSON.stringify(e)}`);
            if (e && e.response && e.response.status !== 503) {
                const { status, data } = e.response;
                res.statusCode = status;
                res.redirect(
                    `${urlPrefix}/oauth2/error?${
                        (
                            data &&
                            Object.keys(data).map(
                                (key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`
                            )
                        ).join("&") || ""
                    }`
                );
            } else {
                const service = getServiceNameFromApi(path);
                console.error(`内部错误，连接${service}服务失败`);
                res.statusCode = 500;
                res.redirect(
                    `${urlPrefix}/oauth2/error?code=${getErrorCodeFromService(service)}&cause=${encodeURIComponent(
                        "内部错误"
                    )}&message=${encodeURIComponent(`连接${service}服务失败`)}`
                );
            }
        }
    } else {
        console.error(`参数不合法，缺少logout_challenge参数`);
        res.statusCode = 400;
        res.redirect(
            `${urlPrefix}/oauth2/error?code=${ErrorCode.INVALID_NO_LOGOUT_CHALLENGE}&cause=${encodeURIComponent(
                "参数不合法"
            )}&message=${encodeURIComponent("缺少logout_challenge参数")}`
        );
    }
};

export default signout;
