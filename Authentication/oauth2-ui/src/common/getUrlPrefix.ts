import Cookies from "js-cookie";
export const WebPortalUrlBasePathName = Cookies.get("X-Forwarded-Web-Client-Basepath") || "/anyshare";
export const getUrlPrefix = () => {
    const prefixPath = Cookies.get("X-Forwarded-Prefix");
    const urlPrefix = !prefixPath || prefixPath === "/" ? "" : prefixPath;
    return urlPrefix;
};

export const getServerPrefix = (req: any) => {
    let urlPrefix = "";
    try {
        const prefixPath = (req as any).cookies["X-Forwarded-Prefix"];
        urlPrefix = !prefixPath || prefixPath === "/" ? "" : prefixPath;
    } catch {
        urlPrefix = "";
    }
    return urlPrefix;
};
