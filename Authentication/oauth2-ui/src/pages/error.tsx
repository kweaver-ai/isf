import { withI18n, useI18n, getRequestLanguage } from "../i18n";
import React from "react";
import { GetServerSideProps } from "next";
import ErrorPage from "./_error";
import { getServerPrefix } from "../common";

interface IErrorProps {
    /**
     * 错误信息
     */
    error: {
        [key: string]: any;
    };
    urlPrefix: string;
}

export default withI18n<IErrorProps>(({ error, urlPrefix }) => {
    const { lang } = useI18n();
    return <ErrorPage {...error} lang={lang} urlPrefix={urlPrefix} />;
});

export const getServerSideProps: GetServerSideProps<IErrorProps> = async ({ req, query }) => {
    const requestLanguage = getRequestLanguage(req);
    const urlPrefix = getServerPrefix(req);

    return {
        props: {
            lang: requestLanguage,
            error: query,
            urlPrefix,
        },
    };
};
