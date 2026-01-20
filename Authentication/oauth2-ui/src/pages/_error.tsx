import React, { useEffect, FunctionComponent } from "react";
import { useCookie } from "react-use";
import { getErrorMessage } from "../core/errorcode";
import { useI18n, withI18n } from "../i18n";
import { WebPortalUrlBasePathName, getUrlPrefix } from "../common";
import LoginError from "../assets/login-error.png";
interface ErrorProps {
  code?: number;
  cause?: string;
  message?: string;
  status_code?: number;
  error?: string;
  error_description?: string;
  error_hint?: string;
  request_id?: string;
  [key: string]: any;
}

const getDefaultErrorMessage = (errorProps: ErrorProps, t: any) => {
  const message =
    errorProps["cause"] || errorProps["message"]
      ? errorProps["cause"]
        ? errorProps["cause"] + "：" + errorProps["message"]
        : errorProps["message"]
      : errorProps["error_description"]
      ? errorProps["error_description"]
      : errorProps["error_hint"];
  const code = errorProps["code"]
    ? errorProps["code"]
    : errorProps["status_code"]
    ? errorProps["status_code"]
    : null;
  return code ? message + `（${t("error-code")}${code}）` : message;
};

export function redirectOrigin(
  topOrigin: string,
  redirects: {
    [key: string]: string | any;
  }
) {
  const url = new URL(topOrigin);
  const urlPrefix = getUrlPrefix();
  switch (true) {
    case /^\/console\//.test(url.pathname):
      if (redirects["console"]) {
        window.location.href = redirects["console"];
      } else {
        top!.window.location.href = `${url.origin}${urlPrefix}/console/`;
      }
      break;
    case /^\/deploy\//.test(url.pathname):
      if (redirects["deploy"]) {
        window.location.href = redirects["deploy"];
      } else {
        top!.window.location.href = `${url.origin}${urlPrefix}/deploy/`;
      }
      break;
    default:
      if (redirects["client"]) {
        window.location.href = redirects["client"];
      } else {
        top!.window.location.href = `${url.origin}${urlPrefix}${WebPortalUrlBasePathName}/`;
      }
      break;
  }
}

export const Redirect: FunctionComponent = ({}) => {
  const [clientOriginUri] = useCookie("client.origin_uri");
  const [consoleOriginUri] = useCookie("console.origin_uri");
  const [deployOriginUri] = useCookie("deploy.origin_uri");
  const redirects = {
    client: clientOriginUri,
    console: consoleOriginUri,
    deploy: deployOriginUri,
  };
  useEffect(() => {
    redirectOrigin(top!.window.location.href, redirects);
  }, []);
  return <></>;
};

export default withI18n<ErrorProps>((errorProps) => {
  const { t } = useI18n();
  const errorStatus = errorProps["code"] || errorProps["status_code"] || 500;
  if (
    [400041000, 400041001, 400041002, 400041003, 400041004, 401, 409].some(
      (error) => error == errorStatus
    )
  ) {
    return <Redirect />;
  }

  return (
    <div className="top-signin">
      <div className={"oauth2-ui-wrapper"}>
        <div className="signin-wrapper">
          <div className="loginerror-logo-wrapper">
            <img className="loginerror-logo-image" src={LoginError}></img>
          </div>
          <div className="loginerror-tips-wrapper">
            {t("signin-error-title")}
          </div>
          <div className="loginerror-message-wrapper">
            {typeof errorProps === "object"
              ? errorStatus
                ? getErrorMessage(
                    errorStatus,
                    t,
                    void 0,
                    getDefaultErrorMessage(errorProps, t)
                  )
                : getDefaultErrorMessage(errorProps, t)
              : null}
          </div>
        </div>
      </div>
    </div>
  );
});
