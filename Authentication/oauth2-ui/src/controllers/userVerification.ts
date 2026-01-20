import { Request, Response } from "express";
import { usermanagementPrivateApi } from "../core/config";
const userVerification = async (req: Request, res: Response) => {
    const { account } = req.query;
    try {
        if (account && typeof account === "string") {
            console.log(`[${Date()}] [INFO]  {/api/user-management/v1/pwd-retrieval-method} GET}  START`);
            const { data } = await usermanagementPrivateApi.get(
                `/api/user-management/v1/pwd-retrieval-method?account=${encodeURIComponent(account)}`
            );
            console.log(`[${Date()}] [INFO]  {/api/user-management/v1/pwd-retrieval-method} GET}  SUCCESS`);
            const formatTelephoneFn = (telephone: string) => {
                if (telephone.length <= 6) {
                    if (telephone.length === 1) {
                        return "*";
                    } else {
                        return telephone[0] + "*".repeat(telephone.length - 1);
                    }
                } else {
                    const displayNum = 3;
                    return (
                        telephone.substring(0, displayNum) +
                        "*****" +
                        telephone.substring(telephone.length - displayNum)
                    );
                }
            };
            const formatEmailFn = (email: string) => {
                let atIndex = email.indexOf("@");
                let accountBefore = email.slice(0, atIndex);
                if (accountBefore.length <= 5) {
                    const suffix = email.slice(atIndex);
                    const hiddenPart = "*".repeat(accountBefore.length - 1);
                    return accountBefore.length === 1 ? "*" + suffix : accountBefore[0] + hiddenPart + suffix;
                } else {
                    let firstThree = email.substring(0, 3);
                    let secondTwo = email.substring(atIndex - 2, atIndex);
                    let stars = "";
                    for (let i = 0; i < atIndex - 5; i++) {
                        stars += "*";
                    }
                    return `${firstThree}${stars}${secondTwo}@${email.substring(atIndex + 1)}`;
                }
            };
            res.json({
                telephone: data?.telephone ? formatTelephoneFn(data.telephone) : data?.telephone,
                email: data?.email ? formatEmailFn(data.email) : data?.email,
                status: data.status,
            });
        } else {
            console.log(
                `[${Date()}] [INFO]  {/api/user-management/v1/pwd-retrieval-method} GET} 账户名不能为空 `,
                account
            );
        }
    } catch (error: any) {
        console.log(`[${Date()}] [INFO]  {/api/user-management/v1/pwd-retrieval-method} GET}  FAILD `, error);
        if (error && error?.response && error?.response?.status && error?.response?.data) {
            const { status, data } = error.response;
            res.statusCode = status;
            res.json({ ...data });
        } else {
            res.statusCode = 500;
            res.json({
                code: 500,
                cause: "oauth2-ui内部错误",
                message: "内部错误",
            });
        }
    }
};
export default userVerification;
