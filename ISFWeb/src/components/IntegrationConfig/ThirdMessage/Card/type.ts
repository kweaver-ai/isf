export enum DefaultKeys {
    /**
     * 第三方app名字
     */
    ThirdpartyName = 'thirdparty_name',

    /**
     * 第三方配置开关
     */
    Enabled = 'enabled',

    /**
     * 消息插件类名
     */
    ClassName = 'class_name',

    /**
     * 消息类型
     */
    Channels = 'channels',

    /**
     * 插件需要其他配置，透传给第三方插件
     */
    Config = 'config',
}

export const defaultKeys = [
    DefaultKeys.ThirdpartyName,
    DefaultKeys.Enabled,
    DefaultKeys.ClassName,
    DefaultKeys.Channels,
    DefaultKeys.Config,
];