import React, { FunctionComponent } from "react";
import { useLocalStore, useObserver } from "mobx-react-lite";
import { ErrorCode, getErrorMessage } from "../../core/errorcode";
import Button from "antd/lib/button";
import Input from "antd/lib/input";
import { useVerification } from "./verification-context";
import {
  IForgetPasswordState,
  IUserVerificationProps,
  UserverificationStatusEnum,
} from "./type";
import BackIcon from "@icons/back.svg";
import AccountIcon from "@icons/account.svg";
import axios from "axios";
import { getUrlPrefix } from "../../common/getUrlPrefix";

enum VerifymethodsEnum {
  email = "email",
  telephone = "telephone",
}

export const UserVerification: FunctionComponent<IUserVerificationProps> = ({
  t,
  redirect,
}) => {
  const verification = useVerification() as IForgetPasswordState;
  const urlPrefix = getUrlPrefix();

  const store = useLocalStore(() => {
    return {
      account: "",
      errorStatus: ErrorCode.Normal,
      checkAccountSuccess: false,
      sendVerificationCodeSuccess: false,
      verificatioCodeType: VerifymethodsEnum.email,
      telephone: "",
      email: "",
      async getUserverification() {
        if (!this.account) {
          this.errorStatus = ErrorCode.NOACCOUNT;
          return;
        }
        try {
          const {
            data: { telephone, email, status },
          } = await axios.get(
            `${location.protocol}//${location.hostname}:${location.port}${urlPrefix}/oauth2/userverification`,
            {
              params: {
                account: this.account,
              },
            }
          );
          this.telephone = telephone;
          this.email = email;
          switch (status) {
            case UserverificationStatusEnum.available:
              if (telephone || email) {
                this.checkAccountSuccess = true;
              } else if (telephone === "" && email === "") {
                this.errorStatus = ErrorCode.UNBOUND_PHONENUMBER_AND_EMAIL;
              } else if (email === "") {
                this.errorStatus = ErrorCode.UNBOUND_EMAIL;
              } else if (telephone === "") {
                this.errorStatus = ErrorCode.UNBOUND_PHONENUMBER;
              }
              break;
            case UserverificationStatusEnum.invalid_account:
              this.errorStatus = ErrorCode.invalid_account;
              break;
            case UserverificationStatusEnum.disable_user:
              this.errorStatus = ErrorCode.disable_user;
              break;
            case UserverificationStatusEnum.enable_pwd_control:
              this.errorStatus = ErrorCode.enable_pwd_control;
              break;
            case UserverificationStatusEnum.non_local_user:
              this.errorStatus = ErrorCode.non_local_user;
              break;
            case UserverificationStatusEnum.unable_pwd_retrieval:
              this.errorStatus = ErrorCode.unable_pwd_retrieval;
              break;
            default:
              break;
          }
          if (!this.checkAccountSuccess) return;
          verification.updateVerificationValue({
            telephone,
            email,
          });
          verification.updateAccount(this.account);
          verification.checkAccountSuccess();
        } catch (e: any) {
          if (e?.response?.data) {
            const {
              response: { data: err },
            } = e;
            this.errorStatus = err.code;
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
          onClick={() => (location.href = redirect)}
        >
          <BackIcon />
        </span>
        <Input
          className="input-item verification-item"
          type="text"
          prefix={
            <span className="icon">
              <AccountIcon />
            </span>
          }
          placeholder={t("reset-pass-placeholder")}
          value={store.account}
          onChange={(e) => {
            store.account = e.target.value;
            store.errorStatus = ErrorCode.Normal;
          }}
          onDrop={(e) => {
            e.preventDefault();
          }}
        />
        <Button
          className="oem-button as-components-oem-background-color"
          type="primary"
          onClick={() => {
            store.getUserverification();
          }}
        >
          {t("reset-next")}
        </Button>
        {store.errorStatus !== ErrorCode.Normal ? (
          <div className="error-message-text">
            {getErrorMessage(store.errorStatus, t)}
          </div>
        ) : null}
      </div>
    );
  });
};
