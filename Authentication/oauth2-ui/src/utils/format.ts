const ProtocolString = "http";

function isIpv6(ip: string) {
    return ip.includes(":");
}

export function urlFormat(url: string | undefined) {
    try {
        if (!url) {
            return url;
        }

        if (!url.includes(ProtocolString)) {
            return isIpv6(url) ? `[${url}]` : url;
        }

        const { protocol, hostname, port } = new URL(url);

        if (isIpv6(hostname)) {
            if (/^\[(.+)\]$/.test(hostname)) {
                return url;
            }
            return `${protocol}//[${hostname}]${port ? ":" + port : ""}`;
        }

        return url;
    } catch (error) {
        console.error(`内部错误，格式化URL错误`, error);
    }
}