import { Request, Response } from "express";
import { getConsentRequest, acceptConsentRequest } from "../services/hydra";
import { getServiceNameFromApi, getErrorCodeFromService } from "../core/config";
import { ErrorCode } from "../core/errorcode";
import { getApplicationbyslug } from "../api/getApplicationbyslug";

const consent = async (req: Request, res: Response) => {
    const { consent_challenge: challenge } = req.query;
    let urlPrefix = "";
    try {
        const prefixPath = req.cookies["X-Forwarded-Prefix"];
        urlPrefix = !prefixPath || prefixPath === "/" ? "" : prefixPath;
    } catch {
        urlPrefix = "";
    }
    if (challenge) {
        try {
            const { requested_access_token_audience, requested_scope, context } = await getConsentRequest(
                challenge as string
            );
            const lastTimestamp = Date.now();

            console.log(
                `[${Date()}] [INFO]  {/api/authentication/v1/config/remember_for get} START  +${
                    Date.now() - lastTimestamp
                }ms`
            );
            const { remember_for } = await getApplicationbyslug("remember_for");

            console.log(
                `[${Date()}] [INFO]  {/api/authentication/v1/config/remember_for get} SUCCESS  +${
                    Date.now() - lastTimestamp
                }ms`
            );

            const { redirect_to } = await acceptConsentRequest(challenge as string, {
                grant_access_token_audience: requested_access_token_audience,
                grant_scope: requested_scope,
                session: {
                    access_token: { ...context },
                },
                remember: false,
                remember_for: remember_for ? remember_for : 0,
            });
            res.redirect(redirect_to);
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
        console.error(`参数不合法，缺少consent_challenge参数`);
        res.statusCode = 400;
        res.redirect(
            `${urlPrefix}/oauth2/error?code=${ErrorCode.INVALID_NO_CONSENT_CHALLENGE}&cause=${encodeURIComponent(
                "参数不合法"
            )}&message=${encodeURIComponent("缺少consent_challenge参数")}`
        );
    }
};

export default consent;
