import React, { FunctionComponent, useEffect } from "react";
import { useLocalStore, useObserver } from "mobx-react-lite";
import { reaction } from "mobx";
import classNames from "classnames";
import { pki } from "node-forge";
import { ErrorCode, getErrorMessage } from "../../core/errorcode";
import openapi, { Auth1ModifypasswordReq } from "../../http/index";
import Button from "antd/lib/button";
import Input from "antd/lib/input";
import message from "antd/lib/message";
import { Password } from "../../controls";
import { useVerification } from "./verification-context";
import {
  VerificationType,
  IForgetPasswordState,
  IResetPasswordProps,
  IResetPasswordState,
} from "./type";
import VerificationCodeIcon from "@icons/verification-code.svg";
import BackIcon from "@icons/back.svg";
import NewPasswordIcon from "@icons/new-password.svg";
import ConfirmPasswordIcon from "@icons/confirm-password.svg";
import openApi from "../../http/index";
import { getUrlPrefix } from "../../common/getUrlPrefix";

const PublicKey: any = pki.publicKeyFromPem(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDB2fhLla9rMx+6LWTXajnK11Kd
p520s1Q+TfPfIXI/7G9+L2YC4RA3M5rgRi32s5+UFQ/CVqUFqMqVuzaZ4lw/uEdk
1qHcP0g6LB3E9wkl2FclFR0M+/HrWmxPoON+0y/tFQxxfNgsUodFzbdh0XY1rIVU
IbPLvufUBbLKXHDPpwIDAQAB
-----END PUBLIC KEY-----`);

export const ResetPassword: FunctionComponent<IResetPasswordProps> = ({
  redirect,
  t,
}) => {
  const verification = useVerification() as IForgetPasswordState;
  const urlPrefix = getUrlPrefix();
  const store = useLocalStore<IResetPasswordState>(() => {
    return {
      captcha: "",
      newPassword: "",
      confirmPassword: "",
      count: 60,
      errorStatus: ErrorCode.Normal,
      errorInfo: null,
      get getError(): string {
        switch (store.errorStatus) {
          case ErrorCode.PasswordInvalidLocked:
            return getErrorMessage(store.errorStatus, t, {
              time: store.errorInfo.detail.remainlockTime,
            });
          case ErrorCode.PasswordWeak:
            return getErrorMessage(store.errorStatus, t, {
              length: verification.strongPasswordLength.toString(),
            });
          default:
            return getErrorMessage(store.errorStatus, t);
        }
      },
      checkInput(): boolean {
        switch (true) {
          case !store.captcha:
            store.errorStatus = ErrorCode.NoCaptcha;
            return false;
          case !store.newPassword:
            store.errorStatus = ErrorCode.NoNewPassword;
            return false;
          case !store.confirmPassword:
            store.errorStatus = ErrorCode.NoConfirmPassword;
            return false;
          case store.newPassword !== store.confirmPassword:
            store.errorStatus = ErrorCode.NewConfirmInconsitent;
            return false;
          default:
            return true;
        }
      },
      checkSendVcodeStatus(): boolean {
        (async () => {
          await verification.updatePasswordConfig();
        })();
        if (
          !verification.sendVcodeType.sendVcodeByEmail &&
          !verification.sendVcodeType.sendVcodeBySMS
        ) {
          store.errorStatus = ErrorCode.CloseForgetPasswordResetBySend;
          return false;
        } else if (
          verification.verificationType === VerificationType.PHONE &&
          !verification.sendVcodeType.sendVcodeBySMS
        ) {
          store.errorStatus = ErrorCode.SMSClose;
          return false;
        } else if (
          verification.verificationType === VerificationType.EMAIL &&
          !verification.sendVcodeType.sendVcodeByEmail
        ) {
          store.errorStatus = ErrorCode.EmailClose;
          return false;
        } else {
          return true;
        }
      },
      redirectTime: 1000,
      async reset() {
        if (store.checkInput() && store.checkSendVcodeStatus()) {
          try {
            const body = {
              account: verification.account,
              oldpwd: "",
              newpwd: btoa(
                PublicKey.encrypt(store.newPassword, "RSAES-PKCS1-V1_5")
              ),
              isforgetpwd: true,
              [verification.verificationType === VerificationType.PHONE
                ? "telnumber"
                : "emailaddress"]:
                verification.verificationType === VerificationType.PHONE
                  ? verification.verificationValue?.telephone
                  : verification.verificationValue?.email,
              vcodeinfo: {
                uuid: verification.verificationId,
                vcode: store.captcha,
              },
            };
            await openapi.post(
              "/eacp/v1/auth1/modifypassword",
              body as Auth1ModifypasswordReq
            );

            message.success(t("reset-password-success"));
            // 重置成功，登录界面
            setTimeout(() => {
              window.location.href =
                redirect ||
                `${location.protocol}//${location.hostname}:${location.port}${urlPrefix}`;
            }, this.redirectTime);
          } catch (e: any) {
            if (e.response) {
              const {
                response: { data: err, status },
              } = e;
              store.errorInfo = err;
              if (err.code === ErrorCode.NewPasswordIsInitial) {
                store.errorStatus = ErrorCode.NewIsInitial;
              } else {
                store.errorStatus = err.code || status;
              }
            } else {
              store.errorInfo = null;
              store.errorStatus = ErrorCode.NoNetwork;
            }
          }
        }
      },
      async sendVcodeAgain() {
        if (store.checkSendVcodeStatus()) {
          try {
            const {
              data: { uuid },
            } = await openApi.post("/eacp/v1/auth1/pwd-retrieval-vcode", {
              account: verification.account,
              type: verification.verificationType as "telephone" | "email",
            });
            verification.updateVerificationId(uuid);
            store.count = 60;
          } catch (e: any) {
            if (e.response) {
              const {
                response: { data: err, status },
              } = e;
              store.errorInfo = err;
              switch (err.code) {
                case ErrorCode.PasswordChangeNotSupported:
                  store.errorStatus =
                    verification.verificationType === VerificationType.PHONE
                      ? ErrorCode.SMSUserNoLocal
                      : ErrorCode.EmailUserNoLocal;
                  break;
                case ErrorCode.PasswordRestricted:
                  store.errorStatus =
                    verification.verificationType === VerificationType.PHONE
                      ? ErrorCode.SMSUserControlled
                      : ErrorCode.EmailUserControlled;
                  break;
                case ErrorCode.SendVcodeServerUnavailable:
                  store.errorStatus =
                    verification.verificationType === VerificationType.PHONE
                      ? ErrorCode.SMSClose
                      : ErrorCode.EmailClose;
                  break;
                default:
                  store.errorStatus = err.code || status;
              }
            } else {
              store.errorInfo = null;
              store.errorStatus = ErrorCode.NoNetwork;
            }
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
      <div className="content">
        <span
          className="back back-pass"
          onClick={async () => {
            await verification.updatePasswordConfig();
            verification.returnSendVcode();
          }}
        >
          <BackIcon />
        </span>
        <div className="reset-password-tip">
          <div className="verification-message">
            {verification.verificationType === VerificationType.PHONE
              ? t("sms-captcha-tip")
              : t("email-captcha-tip")}
            {verification.verificationType === VerificationType.PHONE
              ? verification.verificationValue?.telephone
              : verification.verificationValue?.email}
          </div>
        </div>
        <div className="vcode-wrapper">
          <Input
            className="input-item reset-password-item reset-password-captcha-item"
            type="text"
            autoComplete="off"
            prefix={
              <span className="icon">
                <VerificationCodeIcon />
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
          <div
            className={classNames(
              "vcode vcode-captcha",
              store.count > 0 && "vcode-disabled vcode-captcha"
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
          </div>
        </div>
        <Password
          className="input-item reset-password-item"
          type="password"
          autoComplete="new-password"
          prefix={
            <span className="icon">
              <NewPasswordIcon />
            </span>
          }
          placeholder={t("new-password")}
          value={store.newPassword}
          onChange={(e) => {
            store.newPassword = e.target.value;
            store.errorStatus = ErrorCode.Normal;
          }}
          onDrop={(e) => {
            e.preventDefault();
          }}
        />
        <Password
          className="input-item reset-password-item"
          type="password"
          autoComplete="new-password"
          prefix={
            <span className="icon">
              <ConfirmPasswordIcon />
            </span>
          }
          placeholder={t("confirm-password")}
          value={store.confirmPassword}
          onChange={(e) => {
            store.confirmPassword = e.target.value;
            store.errorStatus = ErrorCode.Normal;
          }}
          onDrop={(e) => {
            e.preventDefault();
          }}
        />
        <div className="tip">
          {verification.strongPasswordStatus
            ? t("strong-password-tip", {
                length: verification.strongPasswordLength,
              })
            : t("weak-password-tip")}
        </div>
        <Button
          className="oem-button reset-password-button as-components-oem-background-color"
          type="primary"
          onClick={store.reset}
        >
          {t("reset-password")}
        </Button>
        {store.errorStatus !== ErrorCode.Normal ? (
          <div className="error-message-text">{store.getError}</div>
        ) : null}
      </div>
    );
  });
};
