export enum LoginWay {
    Account,
    AccountAndImgCaptcha,
    AccountAndSmsCaptcha,
    AccountAndDynamicPassword,
}

export enum LoginAuthType {
    /**
     * 账号密码
     */
    account = 'account',

    /**
     * 账号密码+图形验证码
     */
    accountAndImageCaptcha = 'account_and_image_captcha',

    /**
     * 账号密码+短信验证码
     */
    accountAndSMSCaptcha = 'account_and_SMS_captcha',

    /**
     * 账号密码+动态密码
     */
    accountAndDynamicPassword = 'account_and_dynamic_password',
}