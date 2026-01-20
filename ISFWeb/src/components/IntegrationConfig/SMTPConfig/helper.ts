/**
 * 输入不合法状态
 */
export enum ValidateState {
    /**
     * 正常
     */
    OK,

    /**
     * 空值
     */
    Empty,

    /**
     * 邮件服务器输入错误
     */
    ServerError,

    /**
     * 端口输入有误
     */
    PortError,

    /**
     * 邮件地址输入错误
     */
    EmailError,

}

/**
 * 安全连接选项
 */
export enum SafeMode {
    /**
     * 默认值，无
     */
    Default,

    /**
     * SSL/TLS
     */
    SslOrTsl,

    /**
     * STARTTLS
     */
    Starttls,
}

/**
 * 端口值
 */
export enum Port {
    /**
     * 默认值
     */
    Default = 25,

    /**
     *  SSL/TSL对应的端口值
     */
    SslOrTsl = 465,

    /**
     * STARTTLS对应的端口
     */
    Starttls = 587,
}

/**
 * 测试状态
 */
export enum TestStatus {
    /**
     * 未开始
     */
    NoStart,

    /**
     * 正在测试
     */
    Testing,

    /**
     * 测试完成
     */
    Tested,
}