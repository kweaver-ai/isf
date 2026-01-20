import { createContext, FunctionComponent, useState, useRef, useEffect, useContext } from "react";
import rosetta from "rosetta";
import acceptLanguageParser from "accept-language-parser";
import zhHans from "./translations/zh-hans.json";
import zhHant from "./translations/zh-hant.json";
import en from "./translations/en.json";
import viVN from "./translations/vi-vn.json";
import { NextPage, GetServerSideProps } from "next";
import { ParsedUrlQuery } from "querystring";
import { IncomingMessage } from "http";

const i18n = rosetta();

i18n.set("zh", zhHans);
i18n.set("zh-cn", zhHans);
i18n.set("zh-tw", zhHant);
i18n.set("zh-hk", zhHant);
i18n.set("en", en);
i18n.set("en-us", en);
i18n.set("vi-vn", viVN);
i18n.set("vi", viVN);

export interface II18nContextType {
    lang: string;
    t: typeof i18n.t;
    setLang: typeof i18n.locale;
}

export const defaultLanguage = "en";

export const supportLanguages = ["zh", "zh-cn", "zh-tw", "zh-hk", "en", "en-us", "vi", "vi-vn"];

export const I18nContext = createContext<II18nContextType>({
    lang: defaultLanguage,
    t: (...args) => i18n.t(...args),
    setLang: (lang: string) => i18n.locale(lang),
});

export interface II18nProviderProps {
    lang?: string;
}

export const I18nProvider: FunctionComponent<II18nProviderProps> = ({ lang = defaultLanguage, children }) => {
    const firstRender = useRef(true);
    const currentLangRef = useRef(lang);
    const [, setTick] = useState(0);

    if (lang && firstRender.current) {
        firstRender.current = false;
        i18n.locale(lang);
    }

    useEffect(() => {
        if (lang) {
            i18n.locale(lang);
            currentLangRef.current = lang;
            setTick((tick) => tick + 1);
        }
    }, [lang]);

    const i18nContextValue: II18nContextType = {
        lang: currentLangRef.current,
        t: (...args) => i18n.t(...args),
        setLang: ((lang: string) => {
            i18n.locale(lang);
            currentLangRef.current = lang;
            setTick((tick) => tick + 1);
        }) as typeof i18n.locale,
    };

    return <I18nContext.Provider value={i18nContextValue}>{children}</I18nContext.Provider>;
};

export const useI18n: () => II18nContextType = () => useContext(I18nContext);

export function withI18n<P>(Page: NextPage<P>): NextPage<P & { lang: string }> {
    const WrapperedPage: NextPage<P & { lang: string }> = ({ lang, ...otherProps }) => {
        return (
            <I18nProvider lang={lang}>
                <Page {...(otherProps as any)} />
            </I18nProvider>
        );
    };

    WrapperedPage.displayName = `I18nWrappered(${Page.displayName || Page.name})`;

    return WrapperedPage;
}

export type GetI18nServerSideProps<
    P extends { [key: string]: any } = { [key: string]: any },
    Q extends ParsedUrlQuery = ParsedUrlQuery
> = GetServerSideProps<P & { lang: string }, Q>;

export function getRequestLanguage(req: IncomingMessage): string {
    const acceptLanguageFallback = () => {
        if (req.headers["accept-language"]) {
            if (
                req.headers["accept-language"].startsWith("zh-Hans-CN") ||
                req.headers["accept-language"].startsWith("zh-Hans") ||
                req.headers["accept-language"].startsWith("zh-Hant")
            ) {
                return "zh-cn";
            } else if (req.headers["accept-language"].startsWith("en-GB")) {
                return "en";
            } else {
                return null;
            }
        }
        return null;
    };

    const lang = (
        ((req as any)?.query?.lang as string) ||
        ((req as any)?.cookies?.lang as string) ||
        acceptLanguageParser.pick(
            ["zh", "zh-cn", "zh-TW", "zh-HK", "en", "en-US", "vi", "vi-VN"],
            req.headers["accept-language"] || ""
        ) ||
        acceptLanguageFallback() ||
        defaultLanguage
    ).toLowerCase();

    if (supportLanguages.includes(lang)) {
        return lang;
    }

    return defaultLanguage;
}
