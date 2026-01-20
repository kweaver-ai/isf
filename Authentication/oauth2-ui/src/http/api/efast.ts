/* eslint-disable @typescript-eslint/ban-types */
/**
 * API to access AnyShare
 *
 * 如有任何疑问，可到开发者社区提问：https://developers.aishu.cn
 * # Authentication
 * - 调用需要鉴权的API，必须将token放在HTTP heÏader中："Authorization: Bearer ACCESS_TOKEN"
 * - 对于GET请求，除了将token放在HTTP header中，也可以将token放在URL query string中："tokenid=ACCESS_TOKEN"
 *
 */

export interface EFAST {
    "/eacp/v1/auth1/modifypassword": {
        /**
         * 备注：用户旧密码和新密码采用RSA加密
         */
        POST: {
            body: Auth1ModifypasswordReq;
        };
    };
    "/eacp/v1/auth1/sendauthvcode": {
        /**
         * 备注：向用户发送登录验证码
         */
        POST: {
            body: Auth1SendauthvcodeReq;
            response: Auth1SendauthvcodeRes;
        };
    };
    "/eacp/v1/auth1/pwd-retrieval-vcode": {
        POST: {
            body: SendVCodeInfo;
            response: VCodeUUID;
        };
    };
    "/eacp/v1/auth1/getvcode": {
        POST: {
            body: Auth1GetvcodeReq;
            response: Auth1GetvcodeRes;
        };
    };
    "/eacp/v1/auth1/login-configs": {
        GET: {
            response: Auth1GetLoginConfigsRes;
        };
    };
}

/**
 * 设备信息
 */
export interface Auth1SendauthvcodeReqDeviceinfo {
    /**
     * 设备名称
     */
    name?: string;
    /**
     * 客户端类型
     */
    client_type:
        | "unknown"
        | "ios"
        | "android"
        | "windows_phone"
        | "windows"
        | "mac_os"
        | "web"
        | "mobile_web"
        | "console_web"
        | "deploy_web"
        | "nas"
        | "linux";
    /**
     * 设备硬件类型，自定义。如：
     * iphone5s
     * ipad
     * 联想一体机
     * mac
     */
    description?: string;
    /**
     * 设备唯一标识号，
     * windows下取mac地址
     * ios取udid
     * web为空
     */
    udids?: string[];
    [k: string]: unknown;
}

/**
 * @example `{"account":"user01","password":"xxxxxxxxxx","oldtelnum":"","device":{"name":"eisoo测试iphone","client_type":"web","description":"IPhone","udids":["0a-23-fd-dd-aa-dd-xc"]}}`
 */
export interface Auth1SendauthvcodeReq {
    /**
     * 用户登录账号
     */
    account: string;
    /**
     * 加密后的密文
     */
    password: string;
    /**
     * 上一次的获取的手机号（处理管理员修改手机号的情况）
     */
    oldtelnum: string;
    /**
     * 设备信息
     */
    device: Auth1SendauthvcodeReqDeviceinfo;
    [k: string]: unknown;
}

/**
 * 验证码信息
 */
export interface Auth1ModifypasswordReqVcodeinfo {
    /**
     * 验证码唯一标识
     */
    uuid: string;
    /**
     * 验证码字符串
     */
    vcode: string;
    [k: string]: unknown;
}

/**
 * @example `{"account":"eisoo_user1","oldpwd":"xxxxxxxxxx ","newpwd":"xxxxxxxxxxxxx","vcodeinfo":{"uuid":"5501ebf8-a2e3-11e7-9b66-005056af48ce","vcode":"6PEd"},"isforgetpwd":false,"emailaddress":"abc@163.com","telnumber":"13166668888"}`
 */
export interface Auth1ModifypasswordReq {
    /**
     * 用户登录名
     */
    account: string;
    /**
     * 用户旧密码
     */
    oldpwd: string;
    /**
     * 用户新密码
     */
    newpwd: string;
    /**
     * 验证码信息
     */
    vcodeinfo?: Auth1ModifypasswordReqVcodeinfo;
    /**
     * 忘记密码时的获取邮箱
     */
    emailaddress?: string;
    /**
     * 忘记密码时的获取手机号
     */
    telnumber?: string;
    /**
     * 是否忘记密码
     */
    isforgetpwd?: boolean;
    [k: string]: unknown;
}


/**
 * @example `{"authway":"138******111","sendinterval":60,"isduplicatesended":false}`
 */
export interface Auth1SendauthvcodeRes {
    /**
     * 用户手机号
     */
    authway: string;
    /**
     * 短信发送验证码的间隔（单位/秒）
     */
    sendinterval: number;
    /**
     * 是否重复发送了
     * true表示在时间间隔内重复发送了
     * false表示未重复发送
     */
    isduplicatesended: boolean;
    [k: string]: unknown;
}

/**
 * 发送验证码信息
 */
export interface SendVCodeInfo {
    /**
     * 账户名
     */
    account: string;
    /**
     * 发送验证码类型
     * - `telephone` 通过手机短信发送
     * - `email` 通过邮箱发送
     *
     */
    type: "telephone" | "email";
    [k: string]: unknown;
}

/**
 * 验证码唯一标识
 */
export interface VCodeUUID {
    /**
     * 验证码唯一标识
     */
    uuid: string;
    [k: string]: unknown;
}

/**
 * @example `{"uuid":""}`
 */
export interface Auth1GetvcodeReq {
    /**
     * 上一条验证码的唯一标识
     */
    uuid: string;
    [k: string]: unknown;
}

/**
 * @example `{"uuid":"441857eb-d170-4f4b-ba60-8e2bb5676284","vcode":"iVBORw0KGgoAAAANSUhEUgAAAEQAAAA6CAYAAAAJO\/8DAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAuVSURBVGhD7VoLcJTVFd4EAiSQ54ZHCQlSLHkIA4hFcNqKUkaCYHgUkFYIiCAtCGKdoQplELE1FGzHwQJNUQSS8AjIBGh4JAWMgQzYsaBYqAEKZF\/J7ia7yW6yeezpPffef\/f\/\/71Z8nIJDGfmm81\/77n3nvP955577mY1QdrH4CG8eEiICg8JUeG+JCQ4pn0QzSnhISEq3NdbJjh2qAdSmyYm5a6Qz6HGA0WIyHkR1PPI8ZAQFR4IQqRnkfMiyOdQ44EkRK7TWnQIIaHEiBm9fgh\/6R4HJ7r2gW+DYuFqUHRA8B+yVkHXvrCFrP1iz0EQEZ0stLGlaDchr4cNhK+CY0Cv6RzAl7EqNEFoa0vQZkKiY5LhKIkGkVGdAf\/sEgsJ0UlC2\/2hTYREa1PgfHBvoSGdCd+QaOkvIEWeYNU5p02EFJI9KzKgM+JfwVoIVtnfoYQsJTlDtHBnxjs94hU+tJoQ9QAJISRv3CSZXbRoZ0cE2ToinxAK3+UPnkbBIMS0XoOEi90PWEwiW+QTQuG7\/OFu+Iic9aLF7gfsC+nr8aPDCMGjTLRY26EliAKdphuUaTRwhwA\/dZqeKr324yI5FSU\/OowQrAxFi7UescTpcLhNnEfoNBFg7JsEpoHDwRAxkJJyiwAJYqQhRPO0HGi75IdfQjQxmGz8I0iL5XAyXO8AQgya3tRRJMIyeQY4dmVDQ+kNcLvqAcVdZQNXyQWwb9gIxoQkSoxOE0bGti86\/ysgRE6Ep11EgBoSIe2LEC0hI5ZFREQsOA8cogT4E3edCypfXUrHlGm6UDLFc98dogiRk+FpFxGgRscQEkvftjFxGDTeMXCXmbguXIDqLVvBnrEZHDk50GSq4D1MHJ9mc1J60HnE8\/tHJyMklmyTLqALiwG31c7dBGi4dBXKn51AtxCSJUEXFgVVv32LazGp2byV9ulJvhGv4R9ou0REJyAkkjrj3JHD3SNkXLsBZV16wv+QAJojJN1omlRvElinvsS1mVQ8NpZGSluipFMRghFgikvhbhFxu8EYn0LJ0AvzgpaeQkhKzQfb+SCAuiMn2xwlaLvC8XtJCDphX7WeuwXg\/DiLkxHJdfBYjSQkdCXoTv6OIsBtRnS6RIPb5t1mxr6JrF21xt3QiQiJpGFee+gIdwnAkjrLJ\/TptgnCyAijfTSXEODfzl17+UiAyvRlfKx6Hf9oOyHRicpnGVpGCCukdJoQT7JEB+ovX+EuAZgGDFW85TuaYLBMnAmNt+5Aw5Wr4Mw+AJULloCuWwyUEr3qdRv5SELI3CVwnbTheIwo75pyG3whJ8QfOpQQrDOkKhOrz\/LhPwPzxDQw9O4PrtNfcJdI2PdPURBCt9TaP\/JerzTevgPmF+eROYPAMnUO2NdsgKqFy6EiNQ0qRo4j24mNxTVZtDVPTIAJQWOiqXH6sFiwrVlHqs9S7hZA1fwVYF+9gT8BmMe9wMOeOVBG8kZ58mioyfwYagsKAOobuCYT6+x0mnMarn3HW5g03r4N9vcyQK8dwNamOUl8AgWQEHQqihpknjIdGnXKogvFkjqTYDZ\/IjXFlkyeVDF5snK+jGwbJAlhiHsEqt\/\/gGszwbxjTVMew5I0WSrB+qt0agNGpihS0HYpd\/iFj+NtIISS8fwvuHlMHDv3EIKmgil5FJSFdAcTOR0kcdfXgyEynjmvOHZx\/gi6nfDIrZy9kI9gYowaBPohiWBKGgXm1ClQsy2T9zCpnL+ER4qvnQEihFWg+vC+5DVxq4hY57xMIwCNY8mvG3W+dl8e1yDl+tkS2sYcUL9VvA33pKQ4MnfzEeS43pkDN0gbzoljcQ1z6nTey8QYl0j71FsHbRdtETXaSQjbKo6PPuHmkHyx6A3qiE7Ti\/RLRiFxweRS1wfcdgfXJKQUFoNpxCgPMcwRJAbHsTrE0PMH4HbWUn0psjDJMr0IVtFOT6f9KLW5R2Uke20NCCFomCFmEDeFOFh0nhqjU1zCmIN4EuDb1Uf1g6ZyCx\/BpDb\/ONjf\/QNYZ80nY8MJJGdYye+Q1SHW2Ys4cWxujCTUqc3L5xrkWCdVsfwUQ7SCEN995KPE29WEoGHWtLncDLKHX1nO3448XHFMJKk9hpGTZCzoQomzfQZCXcFZPkophh79KHlsLMtPtjfX8l4A++\/W8zWk+dkt2jJ5DtcgUbr4DRlpDP4IkfvebkKqXlvFzSDH6U8n+RiClSc6WH\/xK64FoAuKpPvfMv2X4Nx7gNYbKHhC6UO0RN9bcFFCVq6m\/Sj21e+pCGF3pfKUsVwDoPqdjT46ASPEtsw\/IdKp4TpdxLXIETp+OjUYdal+UAyYBo0AY59EumW8Y1ltU\/Ph3\/hIzFErfdZghIzhGkhIxr0jxPqCtzaoXLCMG6LM8NhW\/e4mrkXqkE1baIQYaD9+ydyVOsW2CuYcBvyWDE+oJlltUz74CRUhfMtMmsU1xKQFhBAspgxR8eBubKSGeJNqKOlHp5geGlee\/BTVQXFX14C+O\/s6Ufy1II6NpieIbcUaPorMX\/Ild1T6Vh6Tai+6Zm3eca6FN+IhrUqqct\/bRYgnpDdv5aY0d+wyo+uOkbKcS11BEW1jDuK8EoHoZCido3y4l0QUy7g0OobNi\/qRVK9ymvfYdWYd5DrSVwsMASJES94EL8wavJWZZcpsRWGGwCMXt4BUU6DUHc4HfcIjVE8iR\/o0Pz+V6Lq4Jql8P8mit1zsw\/nwE9eoePo5rsHEGDeEk6zctgEiBMFOgopnUrlJTPDfC+a0GaTMHgnlw0aTyCEXvIxNUDnv19BkKOdaROqbwLF7D1QuXQHW+YvB9vt1UPfFOd7JxLErBwy94sGSnk5uuU+DKXEkISwNarb\/nWswscycS21hEae0s8WEqDvlv9lS40JBIZSUlMD5c+eguLgYioqK4POzn8PZM2fg9OnTUFhYCAWnCuDUyZNw4vgJKMCbq0rqv7wExsHDoa7wDG\/xL1Xr1oOx\/4+gqdLGW3wFk6558jROBm4Vb\/7yR4iCCALaplbyRwiS0dDQ0GL849gxbjIR8uwqKgbb22s9YW978y1wXbhIsqzsIkSkyWwG58HDNOpwq5WPGQf1V74Fd613u6E0lF4H+5\/+TLesNzKUW0WC\/B9VEtpNCEYGOlpLDPum1AgTVu6HsUuyfOB0Oqne0SNHwPzUc1Dx5AQw9kukx6pkOJb3+De2mRKGgvknk0h9Mg0qHn+G3F\/iKGkIPJpRBwk0agdDxRPjyZwTwfToKDJHuGeO5oiQcE1FiJoMD+RKCCREqEiA20QiBPF1qQF+\/vpeGPPqLgWQkHpyEcvLy\/MkVTQa7x3KcGbfoWKfRADTxcubPPTxM9xDjASsUbz90pxi4C+JRD5popMZpOfWEII5Q03I+OVZ8OSinQo4HA5KyOHPDguNuxfIDekn9KlFhKjbJGAClQj5+js9PPvabhj9yg4lFu6AmpoacLlc8NmhQ0Lj7gWWhyb4Oi89y9vVTvsjBE8TiZBF7+fBjxdm+uLlTKiurqaEHMzNFRp3L9A7MlHpuJwIOdRO+yMEj1YkpMJqg1ELtnNso8gtvEyBf9vtdqirq4Pc\/fuFxgUam7v1B01UkhcSIVGcBPmnyPHmgHUGEnKjrAIeT9\/qQW7hJbpNMDL2n\/o32Gw2Ssj+vfuEBgYSV4O00AOjI8qLIPztKkJ6ln+KHG8OWHRJW2bSyk9h5LwtMPPtbA8ZGBlIRlVVFSVkb3aO0MhAAX\/gMyT8UQUZiGBCBgIJUKNVhGAFioSgs8WXbsCIlz6EBesP0FNFTgoCScvOyhIaGgjg7++Twgf7kHFXQjwZtwXAchwJQeCxev7yTVj913yaQJEkCUgGYs\/u3UJjv0\/cIsjoFgfhkUOEZCA6jJD8\/HxajmMFikUX1hl4tOJpggkUcwZuE4wMJGPXtu30DvH9QgvngnvDHlJn\/CY0HgZE+G4RNZonJBH+D9vGuGIv62pAAAAAAElFTkSuQmCC"}`
 */
export interface Auth1GetvcodeRes {
    /**
     * 验证码唯一标识
     */
    uuid: string;
    /**
     * 编码后的验证码字符串
     */
    vcode: string;
    [k: string]: unknown;
}

/**
 * 双因子认证配置信息
 */
export interface Auth1GetconfigResDualfactorauthserverstatus {
    /**
     * 是否使用动态密保验证
     */
    auth_by_OTP: boolean;
    /**
     * 是否使用Ukey验证
     */
    auth_by_Ukey: boolean;
    /**
     * 是否使用邮箱验证
     */
    auth_by_email: boolean;
    /**
     * 是否使用短信验证
     */
    auth_by_sms: boolean;
    [k: string]: unknown;
}


/**
 * 限速配置信息
 */
export interface Auth1GetconfigResLimitrateconfig {
    /**
     * 是否开启网络限速
     */
    isenabled: boolean;
    /**
     * 限速类型
     * 0 用户级别限速
     * 1 用户组总体限速
     */
    limittype: number;
    [k: string]: unknown;
}

/**
 * 登录验证码配置信息
 */
export interface Auth1GetconfigResVcodelonginconfig {
    /**
     * 是否启用登录验证码功能
     */
    isenable: boolean;
    /**
     * 达到开启登录验证码的密码出错次数
     */
    passwderrcnt: number;
    [k: string]: unknown;
}

/**
 * 登录验证码配置信息
 */
export interface Auth1GetconfigResVcodeServerStatus {
    /**
     * 邮箱验证码服务器开关
     */
    send_vcode_by_email: boolean;
    /**
     * 短信验证码服务器开关
     */
    send_vcode_by_sms: boolean;
    [k: string]: unknown;
}

/**
 * windows ad单点登录相关配置信息
 */
export interface Auth1GetconfigResWindowsadsso {
    /**
     * 是否开启了windows ad单点登录
     */
    is_enabled: boolean;
    [k: string]: unknown;
}

/**
 * 表示开启了第三方认证，如果未开启，则不会有
 */
export interface Auth1GetconfigResThirdauth {
    /**
     * 唯一第三方认证系统
     */
    id: string;
    /**
     * 第三方认证系统的配置参数
     */
    config: {
        [k: string]: unknown;
    };
    [k: string]: unknown;
}

/**
 * 第三方标密系统配置，如果未开启，则不会有
 */
export interface Auth1GetconfigResThirdcsfsysconfig {
    /**
     * 第三方标密系统唯一标识
     */
    id: string;
    /**
     * 是否仅上传已标密文件
     */
    only_upload_classified: boolean;
    /**
     * 是否仅共享已标密文件
     */
    only_share_classified: boolean;
    /**
     * 是否自动识别文件密级
     */
    auto_match_doc_classfication: boolean;
    [k: string]: unknown;
}

/**
 * login-configs的oem配置信息
 */
export interface Auth1LoginConfigsResOemconfig {
    /**
     * 登录配置，是否允许记住用户名和密码
     */
    rememberpass: boolean;
}

/**
 * @example `{"csf_level_enum":{"内部":6,"机密":8,"秘密":7,"非密":5},"dualfactor_auth_server_status":{"auth_by_OTP":false,"auth_by_Ukey":false,"auth_by_email":false,"auth_by_sms":false},"enable_secret_mode":false,"enable_strong_pwd":false,"internal_link_prefix":"AnyShare:\/\/","oemconfig":{"clearcache":false,"clientlogouttime":-1,"hidecachesetting":false,"maxpassexpireddays":-1,"rememberpass":true},"smtp_server_exists":true,"strong_pwd_length":8,"tag_max_num":30,"vcode_login_config":{"isenable":false,"passwderrcnt":0},"windows_ad_sso":{"is_enabled":true}}`
 */

export interface Auth1GetLoginConfigsRes {
    /**
     * 双因子认证配置信息
     */
    dualfactor_auth_server_status: Auth1GetconfigResDualfactorauthserverstatus;
    /**
     * 是否开启涉密模式
     * true表示开启
     * false表示关闭
     */
    enable_secret_mode: boolean;
    /**
     * 是否开启强密码配置
     */
    enable_strong_pwd: boolean;
    /**
     * 强密码最小长度
     */
    strong_pwd_length: number;
    /**
     * 登录验证码配置信息
     */
    vcode_login_config: Auth1GetconfigResVcodelonginconfig;
    /**
     * 验证码服务器配置信息
     */
    vcode_server_status: Auth1GetconfigResVcodeServerStatus;
    /**
     * windows ad单点登录相关配置信息
     */
    windows_ad_sso: Auth1GetconfigResWindowsadsso;
    /**
     * oem配置信息
     */
    oemconfig: Auth1LoginConfigsResOemconfig;
    /**
     * 表示开启了第三方认证，如果未开启，则不会有
     */
    thirdauth?: Auth1GetconfigResThirdauth;
}