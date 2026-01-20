import { Auth1SendauthvcodeReqDeviceinfo } from "../../http/index";
import Button from "antd/lib/button";
import classNames from "classnames";
import React, { FC } from "react";
import { useI18n } from "../../i18n";

interface PropsType {
    onThirdAuthClick: () => void;
    onAccountLoginClick: () => void;
    device?: Auth1SendauthvcodeReqDeviceinfo;
    className?: string;
    authconfig?: { [key: string]: any };
}

const ConsoleWebdDvice = "console_web";

export const ThirdAuthFirst: FC<PropsType> = ({
    onThirdAuthClick,
    onAccountLoginClick,
    authconfig,
    device,
    className,
}) => {
    const { t } = useI18n();

    const thirdAuthBtnText = authconfig!.thirdauth?.config?.loginButtonText || t("third");

    const getAccountLoginBtnText = () => {
        let text = "";
        if (device?.client_type === ConsoleWebdDvice) {
            text = authconfig!.thirdauth?.config?.consoleAccountLoginButtonText;
        } else {
            text = authconfig!.thirdauth?.config?.clientAccountLoginButtonText;
        }

        return text || t("account-login");
    };

    return (
        <div className={classNames("signin-third-auth-first", "signin-content", className)}>
            <Button
                onClick={onThirdAuthClick}
                className={classNames("login-button", "as-components-oem-background-color")}
                type="primary"
            >
                {thirdAuthBtnText}
            </Button>
            <Button onClick={onAccountLoginClick} className="as-controls-normal-button login-button last-btn">
                {getAccountLoginBtnText()}
            </Button>
        </div>
    );
};
