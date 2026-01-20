export function replaceUrlParams(url: string, params: { [key: string]: string | number }): string {
    return url.replace(/\{\s*([^\{\}]+)\s*\}/g, (_, p) => {
        if (params[p] === undefined) {
            throw new Error(`params["${p}"] is undefined`);
        }
        return params[p] as string;
    });
}
