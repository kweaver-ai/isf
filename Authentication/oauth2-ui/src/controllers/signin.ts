import { Request, Response } from "express";
import { acceptLoginRequest, getLoginRequest } from "../services/hydra";
import {
    eacpPrivateApi,
    authenticationPrivateApi,
    getServiceNameFromApi,
    getErrorCodeFromService,
} from "../core/config";
import { ErrorCode } from "../core/errorcode";
import { getApplicationbyslug } from "../api/getApplicationbyslug";

const signin = async (req: Request, res: Response) => {
    const { challenge, account, password, vcode, dualfactorauthinfo, remember, device } = req.body;
    const DefaultValidityPeriod = 30 * 24 * 60 * 60;
    if (challenge && !challenge.includes("/") && typeof remember === "boolean") {
        try {
            let user_id: string, context: any;
            let lastTimestamp = Date.now();
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

            if (device.client_type === "console_web" || device.client_type === "deploy_web") {
                console.log(`[${Date()}] [INFO]  {/api/eacp/v1/auth1/consolelogin POST} START`);
                // 管理控制台登录
                ({
                    data: { user_id, context },
                } = await eacpPrivateApi.post("/api/eacp/v1/auth1/consolelogin", {
                    credential: {
                        type: "account",
                        account,
                        password,
                        vcode,
                    },
                    device,
                    ip: req.ip || "",
                }));
                console.log(
                    `[${Date()}] [INFO]  {/api/eacp/v1/auth1/consolelogin POST} SUCCESS +${
                        Date.now() - lastTimestamp
                    }ms`
                );
            } else {
                lastTimestamp = Date.now();
                console.log(`[${Date()}] [INFO]  {/api/eacp/v1/auth1/getnew POST} START`);
                // 客户端登录
                ({
                    data: { user_id, context },
                } = await eacpPrivateApi.post("/api/eacp/v1/auth1/getnew", {
                    account,
                    password,
                    vcode,
                    dualfactorauthinfo,
                    device,
                    ip: req.ip || "",
                }));
                console.log(
                    `[${Date()}] [INFO]  {/api/eacp/v1/auth1/getnew POST} SUCCESS +${Date.now() - lastTimestamp}ms`
                );
                if (remember) {
                    const {
                        session_id = "",
                        client: { client_id = "" },
                    } = await getLoginRequest(challenge);
                    lastTimestamp = Date.now();
                    console.log(`[${Date()}] [INFO]  {/api/authentication/v1/session/${session_id} PUT} START`);
                    await authenticationPrivateApi.put(`/api/authentication/v1/session/${session_id}`, {
                        subject: user_id,
                        client_id,
                        remember_for: remember_for ? remember_for : DefaultValidityPeriod,
                        context,
                    });
                    console.log(
                        `[${Date()}] [INFO]  {/api/authentication/v1/session/${session_id} PUT} SUCCESS  +${
                            Date.now() - lastTimestamp
                        }ms`
                    );
                }
            }

            const { redirect_to } = await acceptLoginRequest(challenge, {
                subject: user_id,
                context,
                remember,
                remember_for: remember_for ? remember_for : remember ? DefaultValidityPeriod : 0,
            });
            res.json({
                redirect: redirect_to,
            });
        } catch (e: any) {
            const path = e && e.request && (e.request.path || (e.request._options && e.request._options.path));
            console.error(`[${Date()}] [ERROR]  ${path}  ERROR ${JSON.stringify(e && e.response && e.response.data)}`);
            console.error(`[${Date()}] [ERROR]  ${path}  ERROR ${JSON.stringify(e)}`);
            if (e && e.response && e.response.status !== 503) {
                const { status, data } = e.response;
                res.statusCode = status;
                res.json({ ...data });
            } else {
                const service = getServiceNameFromApi(path);
                console.error(`内部错误，连接${service}服务失败`);
                res.statusCode = 500;
                res.json({
                    code: getErrorCodeFromService(service),
                    cause: "内部错误",
                    message: `连接${service}服务失败`,
                });
            }
        }
    } else {
        console.error(`参数不合法，challenge或remember参数验证失败`);
        res.statusCode = 400;
        res.json({
            code: ErrorCode.INVALID_CHALLENGE_OR_REMEMBER,
            cause: "参数不合法",
            message: "challenge或remember参数验证失败",
        });
    }
};

export default signin;
