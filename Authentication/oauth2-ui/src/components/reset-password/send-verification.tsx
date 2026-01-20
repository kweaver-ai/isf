import React, { FunctionComponent } from "react";
import { useLocalStore, useObserver } from "mobx-react-lite";
import { ErrorCode, getErrorMessage } from "../../core/errorcode";
import Button from "antd/lib/button";
import { useVerification } from "./verification-context";
import {
  VerificationType,
  IForgetPasswordState,
  ISendVerificationProps,
  ISendVerificationState,
} from "./type";
import BackIcon from "@icons/back.svg";
import Radio from "antd/lib/radio";
import Space from "antd/lib/space";
import openApi from "../../http/index";

export const SendVerification: FunctionComponent<ISendVerificationProps> = ({
  t,
}) => {
  const verification = useVerification() as IForgetPasswordState;
  const store = useLocalStore<ISendVerificationState>(() => {
    return {
      errorStatus: ErrorCode.Normal,
      verificationType: verification?.verificationType,
      async sendVcode() {
        try {
          const {
            data: { uuid },
          } = await openApi.post("/eacp/v1/auth1/pwd-retrieval-vcode", {
            account: verification.account,
            type: store.verificationType,
          });
          verification.updateVerificationType(store.verificationType);
          verification.sendVcodeSuccess(uuid);
        } catch (e: any) {
          if (e.response) {
            const {
              response: { data: err, status },
            } = e;
            this.errorStatus = err.code || status;
          } else {
            store.errorStatus = ErrorCode.NoNetwork;
          }
        }
      },
    };
  });

  return useObserver(() => {
    return (
      <div className="content">
        <span
          className="back back-pass"
          onClick={() => verification.returnUserVerification()}
        >
          <BackIcon />
        </span>
        <div className="verifymethods-box">
          <p className="verifymethods-tips">{t("reset-verifymethods")}</p>
          {verification.verificationValue?.email &&
          verification.verificationValue?.telephone ? (
            <Radio.Group
              onChange={(v) => {
                store.verificationType = v.target.value;
              }}
              value={store.verificationType}
              className="verifymethods-radio"
            >
              <Space direction="vertical">
                <Radio value={VerificationType.EMAIL}>
                  <span className="verifymethods-tip">
                    {t("reset-verifymethods-email", {
                      number: verification?.verificationValue?.email,
                    })}
                  </span>
                </Radio>
                <Radio value={VerificationType.PHONE}>
                  <span className="verifymethods-tip">
                    {t("reset-verifymethods-telephone", {
                      number: verification?.verificationValue?.telephone,
                    })}
                  </span>
                </Radio>
              </Space>
            </Radio.Group>
          ) : verification.verificationValue?.email ? (
            <Radio
              value={VerificationType.EMAIL}
              checked
              className="verifymethods-radio"
            >
              <span className="verifymethods-tip">
                {t("reset-verifymethods-email", {
                  number: verification?.verificationValue?.email,
                })}
              </span>
            </Radio>
          ) : (
            <Radio
              value={VerificationType.PHONE}
              checked
              className="verifymethods-radio"
            >
              <span className="verifymethods-tip">
                {t("reset-verifymethods-telephone", {
                  number: verification?.verificationValue?.telephone,
                })}
              </span>
            </Radio>
          )}
          <Button
            className="oem-button as-components-oem-background-color"
            type="primary"
            onClick={() => store.sendVcode()}
          >
            {t("send-captcha")}
          </Button>
          {store.errorStatus !== ErrorCode.Normal ? (
            <div className="error-message-text">
              {getErrorMessage(store.errorStatus, t)}
            </div>
          ) : null}
        </div>
      </div>
    );
  });
};
