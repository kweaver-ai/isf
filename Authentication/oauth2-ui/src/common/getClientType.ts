export function getIsMobile(clientType?: string, userAgent?: string): boolean {
  return (
    (clientType && clientType === "webmobile") ||
    (userAgent
      ? /mobile|nokia|iphone|ipad|android|samsung|htc|blackberry/i.test(
          userAgent
        )
      : false)
  );
}

export function getIsElectronOpenExternal(device?: any): boolean {
  return (
    device?.client_type === "windows" ||
    device?.client_type === "mac_os" ||
    device?.client_type === "linux"
  );
}
