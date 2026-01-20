/**
 * 终端类型枚举
 */
export enum EndpointType {
    /**
     * 同步盘/富客户端
     */
    Windows = 'windows',

    /**
     * iOS客户端
     */
    IOS = 'ios',

    /**
     * Android客户端
     */
    Android = 'android',

    /**
     * Mac客户端
     */
    MacOS = 'mac',

    /**
     * 桌面Web客户端
     */
    Web = 'web',

    /**
     * 移动Web客户端
     */
    MobileWeb = 'mobile_web',

    /**
     * Linux客户端
     */
    Linux = 'linux',

    /**
     * 未知类型的客户端
     */
    Unknown = 'unknown',

    /**
     * 管理控制台
     */
    ConsoleWeb = 'console_web',

    /**
     * 部署控制台
     */
    DeployWeb = 'deploy_web',

    /**
     * NAS
     */
    NAS = 'nas',
}