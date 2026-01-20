import React, { FunctionComponent, useEffect } from "react";
import { reaction } from "mobx";
import { useLocalStore, useObserver } from "mobx-react-lite";
import classNames from "classnames";
import { pki } from "node-forge";
import openapi, { Auth1SendauthvcodeReqDeviceinfo } from "../../http/index";
import { ErrorCode, getErrorMessage } from "../../core/errorcode";
import Button from "antd/lib/button";
import Input from "antd/lib/input";
import axios from "axios";
import VerificationCode from "@icons/verification-code.svg";
import BackIcon from "@icons/back.svg";
import { getUrlPrefix } from "../../common/getUrlPrefix";

export interface AuthBySMSProps {
  cellphone: string; // 短信验证手机号
  timeInterval: number; // 短信验证时间间隔
  isDuplicateSend: boolean; // 是否重复发送
  account: string; // 账号
  password: string; // 密码
  device: Auth1SendauthvcodeReqDeviceinfo; // 设备信息
  rememberLogin: boolean; // 是否记住登录
  challenge: string; // 认证唯一标识
  t: any; // 全球化函数
  csrftoken: string; // csrf token令牌
  refreshAuthMethod: () => void; // 获取登录方式
  oemConfig?: { [key: string]: any }; // oem配置
}

export interface AuthBySMSState {
  count: number; // 倒计时
  captcha: string; // 验证码值
  isDuplicateSend: boolean; // 是否重复发送
  cellphone: string; // 短信验证手机号
  verifying: boolean; // 是否验证中
  errorStatus: number; // 错误状态码
  errorInfo: any; // 错误附加信息
  verify: () => void; // 登录验证
  sendVcodeAgain: () => void; // 发送验证码
  [key: string]: any;
}

const PublicKey: any = pki.publicKeyFromPem(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDB2fhLla9rMx+6LWTXajnK11Kd
p520s1Q+TfPfIXI/7G9+L2YC4RA3M5rgRi32s5+UFQ/CVqUFqMqVuzaZ4lw/uEdk
1qHcP0g6LB3E9wkl2FclFR0M+/HrWmxPoON+0y/tFQxxfNgsUodFzbdh0XY1rIVU
IbPLvufUBbLKXHDPpwIDAQAB
-----END PUBLIC KEY-----`);

const PublicKeyPlus: any = pki.publicKeyFromPem(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsyOstgbYuubBi2PUqeVj
GKlkwVUY6w1Y8d4k116dI2SkZI8fxcjHALv77kItO4jYLVplk9gO4HAtsisnNE2o
wlYIqdmyEPMwupaeFFFcg751oiTXJiYbtX7ABzU5KQYPjRSEjMq6i5qu/mL67XTk
hvKwrC83zme66qaKApmKupDODPb0RRkutK/zHfd1zL7sciBQ6psnNadh8pE24w8O
2XVy1v2bgSNkGHABgncR7seyIg81JQ3c/Axxd6GsTztjLnlvGAlmT1TphE84mi99
fUaGD2A1u1qdIuNc+XuisFeNcUW6fct0+x97eS2eEGRr/7qxWmO/P20sFVzXc2bF
1QIDAQAB
-----END PUBLIC KEY-----
`);

export const AuthBySMS: FunctionComponent<AuthBySMSProps> = ({
  cellphone,
  timeInterval = 60,
  account,
  password,
  isDuplicateSend,
  device,
  challenge,
  rememberLogin,
  t,
  csrftoken,
  refreshAuthMethod,
  oemConfig,
}) => {
  const urlPrefix = getUrlPrefix();

  const store = useLocalStore<AuthBySMSState>(() => {
    return {
      count: timeInterval,
      cellphone: cellphone,
      captcha: "",
      isDuplicateSend: isDuplicateSend,
      verifying: false,
      errorStatus: ErrorCode.Normal,
      errorInfo: null,
      get getError(): string {
        switch (store.errorStatus) {
          case ErrorCode.PasswordInvalidLocked:
            return getErrorMessage(store.errorStatus, t, {
              time: store.errorInfo.detail.remainlockTime,
            });
          default:
            return getErrorMessage(store.errorStatus, t);
        }
      },
      async verify() {
        store.verifying = true;
        if (store.captcha) {
          try {
            const {
              data: { redirect },
            } = await axios.post(
              `${location.protocol}//${location.hostname}:${location.port}${urlPrefix}/oauth2/signin`,
              {
                _csrf: csrftoken,
                challenge,
                account: account,
                password: btoa(
                  PublicKeyPlus.encrypt(password, "RSAES-PKCS1-V1_5")
                ),
                vcode: { id: "", content: "" },
                dualfactorauthinfo: {
                  validcode: { vcode: store.captcha },
                  OTP: { OTP: "" },
                  UKey: {},
                },
                remember: rememberLogin,
                device: device as Auth1SendauthvcodeReqDeviceinfo,
              }
            );
            store.verifying = false;
            window.location.href = redirect;
          } catch (e: any) {
            if (e.response) {
              const {
                response: { data: err },
              } = e;
              // 获取最新登录方式
              await refreshAuthMethod();
              store.errorInfo = err;
              switch (err.code) {
                case ErrorCode.AuthFailed:
                  store.errorStatus = ErrorCode.PasswordChange;
                  break;
                default:
                  store.errorStatus = err.code || err.status_code;
              }
              store.verifying = false;
            } else {
              store.errorInfo = null;
              store.errorStatus = ErrorCode.NoNetwork;
              store.verifying = false;
            }
          }
        } else {
          store.errorStatus = ErrorCode.NoCaptcha;
          store.verifying = false;
        }
      },
      async sendVcodeAgain() {
        try {
          const {
            data: { authway, isduplicatesended, sendinterval },
          } = await openapi.post("/eacp/v1/auth1/sendauthvcode", {
            account,
            password: btoa(PublicKey.encrypt(password, "RSAES-PKCS1-V1_5")),
            oldtelnum: "",
            device,
          });
          store.isDuplicateSend = isduplicatesended;
          store.count = sendinterval;
          store.cellphone = authway;
        } catch (e: any) {
          if (e.response) {
            const {
              response: { data: err },
            } = e;
            store.errorInfo = err;
            store.errorStatus = err.code;
          } else {
            store.errorInfo = null;
            store.errorStatus = ErrorCode.NoNetwork;
          }
        }
      },
    };
  });

  useEffect(
    () =>
      reaction(
        () => [store.count],
        () => {
          if (store.count > 0) {
            const timer = setTimeout(() => {
              store.count = store.count - 1;
            }, 1000);

            return () => {
              clearTimeout(timer);
            };
          }
        },
        { fireImmediately: true }
      ),
    []
  );

  return useObserver(() => {
    return (
      <div className="sms-verification-wrapper">
        <span
          className="back back-sms"
          onClick={() =>
            (location.href = `${location.protocol}//${location.hostname}:${location.port}${urlPrefix}/oauth2/signin?login_challenge=${challenge}`)
          }
        >
          <BackIcon />
        </span>
        <div className="title">{t("verify")}</div>
        <div className="content">
          <div className="tip">{`${t("sms-captcha-tip")}${
            store.cellphone
          }`}</div>
          <Input
            className="input-item sms-verification-item"
            type="text"
            prefix={
              <span className="icon">
                <VerificationCode />
              </span>
            }
            placeholder={t("captcha")}
            value={store.captcha}
            onChange={(e) => {
              store.captcha = e.target.value;
              store.errorStatus = ErrorCode.Normal;
            }}
            onDrop={(e) => {
              e.preventDefault();
            }}
          />
          <span
            className={classNames(
              "vcode vcode-captcha vcode-captcha-sms",
              store.count > 0 &&
                "vcode-disabled vcode-captcha vcode-captcha-sms"
            )}
            onClick={() => {
              if (store.count > 0) {
                return false;
              } else {
                store.sendVcodeAgain();
              }
            }}
          >
            {store.count > 0
              ? t("send-captcha-again", { count: store.count })
              : t("send-captcha")}
          </span>
          <Button
            className="oem-button as-components-oem-background-color"
            type="primary"
            onClick={store.verify}
          >
            {store.verifying ? t("verifying") : t("verify-now")}
          </Button>
          <div>
            {store.errorStatus === ErrorCode.Normal ? null : (
              <div className="error-message-text">{store.getError}</div>
            )}
            {store.errorStatus === ErrorCode.PasswordFailure ||
            store.errorStatus === ErrorCode.PasswordInsecure ||
            store.errorStatus === ErrorCode.PasswordIsInitial ? (
              <div className="change-password-text">
                <a
                  href={`${urlPrefix}/oauth2/change?${
                    oemConfig?.productId
                      ? `product=${oemConfig?.productId}&`
                      : ""
                  }redirect=${location.protocol}//${location.hostname}:${
                    location.port
                  }${urlPrefix}/oauth2/signin?login_challenge=${challenge}&account=${account}`}
                >
                  {t("change")}
                </a>
              </div>
            ) : null}
          </div>
        </div>
      </div>
    );
  });
};
