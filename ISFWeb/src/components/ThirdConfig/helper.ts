import { trim, findIndex, reduceRight } from 'lodash';
import __ from './locale';

/**
 * 允许输入的最大数字
 */
const MAX_SAFE_INTEGER = 999999999999999;

/**
 * 允许输入的最小数字
 */
const MIN_SAFE_INTEGER = -999999999999999;

/**
 * 客户端参数配置默认配置项
 */
export enum DefaultClientConfigParameter {
    // /**
    //  * 第三方认证唯一的key
    //  */
    // AppKey = "appKey",

    // /**
    //  * 第三方认证唯一的id
    //  */
    // AppId = "appId",

    // /**
    //  * 第三方用户点击注销之后的跳转页面
    //  */
    // LogoutUrl = "logoutUrl",

    // /**
    //  * 第三方认证按钮文字内容
    //  */
    // LoginButtonText = "loginButtonText",

    /**
     * url匹配字符串，url匹配之后才会进行参数获取
     */
    MatchUrl = 'matchUrl',

    // /**
    //  * 自动跳转到“authServer”进行认证
    //  */
    // AutoCasRedirect = "autoCasRedirect",

    /**
     * 单点登录客户端认证服务器跳转地址
     */
    AuthServer = 'authServer',

    /**
     * 第三方认证界面显示、隐藏开关
     */
    HideThirdLogin = 'hideThirdLogin',

    /**
     * 登录界面显示、隐藏开关
     */
    HideLogin = 'hideLogin',
}

/**
 * 服务端参数配置默认配置项
 */
export enum DefaultInternalConfigParameter {
    // /**
    //  * 用户账号有效期，单位天，默认-1，表示不限制，6.0.6及以后版本支持
    //  */
    // ValidPeriod = "validPeriod",

    // /**
    //  * 用户创建后默认状态。设置true为启用，false为禁用
    //  */
    // UserCreateStatus = "userCreateStatus",

    // /**
    //  * 用户同步间隔时间，单位秒
    //  */
    // SyncInterval = "syncInterval",

    // /**
    //  * 认证集成模块
    //  */
    // AuthModule = "authModule",

    /**
     * 用户同步框架，固定BaseSyncer，不允许修改
     */
    SyncModule = 'syncModule',

    // /**
    //  * 用户同步模块
    //  */
    // OuModule = "ouModule",
}

/**
 * 终端
 */
export enum Terminal {
    /**
     * 客户端
     */
    Client = 'client',

    /**
     * 服务端
     */
    Internal = 'internal',
}

/**
 * 验证状态
 */
export enum ValidateStatus {
    /**
     * 正常
     */
    Normal = 0,

    /**
     * 输入项为空
     */
    Empty = 1,

    /**
     * 名称重复
     */
    NameRepeat = 2,

    /**
     * 接口报错
     */
    InterfaceError = 3,

    /**
     * ID长度超过限制（128）
     */
    IDLengthExceedsLimit = 4,

    /**
    * Name长度超过限制（128）
    */
    NameLengthExceedsLimit = 5,

    /**
     * 名称长度超过限制，主要用于限定客户端和服务端参数名称和值的长度，安全考虑
     */
    NameToLong = 6,

    /**
     * 值长度超过限制，主要用于限定客户端和服务端参数名称和值的长度，安全考虑
     */
    ValueToLong = 7,
}

/**
 * 输入验证提示
 */
export const ValidateMessage = {
    [ValidateStatus.Empty]: __('此输入项不允许为空。'),
    [ValidateStatus.NameRepeat]: __('该名称已存在，请您重新输入。'),
    [ValidateStatus.IDLengthExceedsLimit]: __('认证服务ID长度不能超过128个字符。'),
    [ValidateStatus.NameLengthExceedsLimit]: __('认证服务名称长度不能超过128个字符。'),
    [ValidateStatus.NameToLong]: __('名称长度不能超过800个字符。'),
    [ValidateStatus.ValueToLong]: __('值长度不能超过800个字符。'),
};

/**
 * 插件类型
 */
export enum NcTPluginType {
    /**
     * 认证插件
     */
    Authentication = 0,

    /**
     * 消息推送插件
     */
    Message = 1,
}

/**
 * 配置参数类型
 */
export enum Types {
    /**
     * 字符串
     */
    StringType = 'string',

    /**
     * 数字
     */
    NumberType = 'number',

    /**
     * 布尔
     */
    BooleanType = 'boolean',

    /**
     * 对象
     */
    ObjectType = 'object',
}

/**
 * 配置参数信息
 */
export interface ConfigItem {
    /**
     * 数组下标，主要用于标识默认配置项的位置，自定义的配置项全部设置为null
     */
    index: number | null;

    /**
     * 配置参数名称
     */
    configName: string;

    /**
     * 配置参数类型
     */
    configType: any;

    /**
     * 配置参数值
     */
    configValue: any;

    /**
     * 是否为默认参数（非默认参数的名称为必填项）
     */
    defaultConfig: boolean;

    /**
     * 值是否为必填项
     */
    valueRequired: boolean;

    /**
     * 配置参数名称验证状态
     */
    nameValidateStatus: ValidateStatus;

    /**
     * 配置参数值验证状态
     */
    valueValidateStatus: ValidateStatus;
}

/**
 * 第三方插件信息
 */
export interface PluginInfo {
    /**
     * 索引Id，用于确定插件在磁盘的位置
     */
    indexId: number;

    /**
     * 唯一标识第三方认证系统
     */
    thirdPartyId: string;

    /**
     * 文件名
     */
    filename: string;

    /**
     * 文件内容
     */
    data: null | string;

    /**
     * 插件类型 0:认证 1:消息推送
     */
    type: NcTPluginType;

    /**
     * 对象存储ID，用于确定插件在存储的位置
     */
    objectId: string;
}

/**
 * 阶段
 */
export enum Stage {
    /**
     * 初始化
     */
    Initialize,

    /**
     * 高级配置
     */
    Advanced,
}

/**
 * 根据类型转换为相应的提示文字
 * @param configType
 */
export function formatterTypeName(configType) {
    switch (configType) {
        case Types.StringType:
            return __('字符')

        case Types.NumberType:
            return __('数字')

        case Types.BooleanType:
            return __('开关')

        default:
            return __('未知类型')
    }
}

/**
 * 错误码
 */
export enum ErrorCode {
    /**
    * 格式错误
    */
    WrongFormat = 2,

    /**
     * 第三方应用id已存在
     */
    IDExists = 21111,

    /**
     * 未配置存储
     */
    NoStorage = 21104,

    /**
     * 配置错误
     */
    InvalidConfig = 20602,

    /**
     * 非法的第三方插件
     */
    InvalidPlugin = 20615,

    /**
     * 不合法的文件名
     */
    IllegalFileName = 20616,
}

/**
 * 转换错误提示
 */
export function formatterError(errID: ErrorCode, errMsg?: string): string {
    switch (errID) {
        case ErrorCode.IDExists:
            return __('第三方应用APPID已存在。')

        case ErrorCode.NoStorage:
            return __('系统存储异常，请检查存储配置。')

        case ErrorCode.InvalidPlugin:
        case ErrorCode.InvalidConfig:
            return __('非法的第三方插件。')

        case ErrorCode.WrongFormat:
            return __('请上传规定格式的文件')

        case ErrorCode.IllegalFileName:
            return __('文件名不能包含 \\ / @ # $ % ^ & * ( ) [ ] 字符，且首位不能使用 + - .字符，长度不能超过255个字符，请重新上传。')

        default:
            return errMsg || __('未知的错误码。')
    }
}

// 客户端默认配置表
export const defaultClientConfig = {
    [DefaultClientConfigParameter.HideLogin]: { index: 1, configType: Types.BooleanType, valueRequired: false },
    [DefaultClientConfigParameter.HideThirdLogin]: { index: 2, configType: Types.BooleanType, valueRequired: false },
    [DefaultClientConfigParameter.AuthServer]: { index: 3, configType: Types.StringType, valueRequired: false },
    // [DefaultClientConfigParameter.AutoCasRedirect]: { index: 4, configType: Types.BooleanType, valueRequired: false },
    [DefaultClientConfigParameter.MatchUrl]: { index: 5, configType: Types.StringType, valueRequired: false },
    // [DefaultClientConfigParameter.LoginButtonText]: { index: 6, configType: Types.StringType, valueRequired: false },
    // [DefaultClientConfigParameter.LogoutUrl]: { index: 7, configType: Types.StringType, valueRequired: false },
    // [DefaultClientConfigParameter.AppId]: { index: 8, configType: Types.StringType, valueRequired: false },
    // [DefaultClientConfigParameter.AppKey]: { index: 9, configType: Types.StringType, valueRequired: false },
}

// 客户端默认配置数组
export const defaultClientConfigArray = [
    {
        index: defaultClientConfig[DefaultClientConfigParameter.HideLogin].index,
        configName: DefaultClientConfigParameter.HideLogin,
        configType: defaultClientConfig[DefaultClientConfigParameter.HideLogin].configType,
        configValue: false,
        defaultConfig: true,
        valueRequired: defaultClientConfig[DefaultClientConfigParameter.HideLogin].valueRequired,
        nameValidateStatus: ValidateStatus.Normal,
        valueValidateStatus: ValidateStatus.Normal,
    },
    {
        index: defaultClientConfig[DefaultClientConfigParameter.HideThirdLogin].index,
        configName: DefaultClientConfigParameter.HideThirdLogin,
        configType: defaultClientConfig[DefaultClientConfigParameter.HideThirdLogin].configType,
        configValue: true,
        defaultConfig: true,
        valueRequired: defaultClientConfig[DefaultClientConfigParameter.HideThirdLogin].valueRequired,
        nameValidateStatus: ValidateStatus.Normal,
        valueValidateStatus: ValidateStatus.Normal,
    },
    {
        index: defaultClientConfig[DefaultClientConfigParameter.AuthServer].index,
        configName: DefaultClientConfigParameter.AuthServer,
        configType: defaultClientConfig[DefaultClientConfigParameter.AuthServer].configType,
        configValue: '',
        defaultConfig: true,
        valueRequired: defaultClientConfig[DefaultClientConfigParameter.AuthServer].valueRequired,
        nameValidateStatus: ValidateStatus.Normal,
        valueValidateStatus: ValidateStatus.Normal,
    },
    // {
    //     index: defaultClientConfig[DefaultClientConfigParameter.AutoCasRedirect].index,
    //     configName: DefaultClientConfigParameter.AutoCasRedirect,
    //     configType: defaultClientConfig[DefaultClientConfigParameter.AutoCasRedirect].configType,
    //     configValue: false,
    //     defaultConfig: true,
    //     valueRequired: defaultClientConfig[DefaultClientConfigParameter.AutoCasRedirect].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    {
        index: defaultClientConfig[DefaultClientConfigParameter.MatchUrl].index,
        configName: DefaultClientConfigParameter.MatchUrl,
        configType: defaultClientConfig[DefaultClientConfigParameter.MatchUrl].configType,
        configValue: '',
        defaultConfig: true,
        valueRequired: defaultClientConfig[DefaultClientConfigParameter.MatchUrl].valueRequired,
        nameValidateStatus: ValidateStatus.Normal,
        valueValidateStatus: ValidateStatus.Normal,
    },
    // {
    //     index: defaultClientConfig[DefaultClientConfigParameter.LoginButtonText].index,
    //     configName: DefaultClientConfigParameter.LoginButtonText,
    //     configType: defaultClientConfig[DefaultClientConfigParameter.LoginButtonText].configType,
    //     configValue: '',
    //     defaultConfig: true,
    //     valueRequired: defaultClientConfig[DefaultClientConfigParameter.LoginButtonText].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    // {
    //     index: defaultClientConfig[DefaultClientConfigParameter.LogoutUrl].index,
    //     configName: DefaultClientConfigParameter.LogoutUrl,
    //     configType: defaultClientConfig[DefaultClientConfigParameter.LogoutUrl].configType,
    //     configValue: '',
    //     defaultConfig: true,
    //     valueRequired: defaultClientConfig[DefaultClientConfigParameter.LogoutUrl].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    // {
    //     index: defaultClientConfig[DefaultClientConfigParameter.AppId].index,
    //     configName: DefaultClientConfigParameter.AppId,
    //     configType: defaultClientConfig[DefaultClientConfigParameter.AppId].configType,
    //     configValue: '',
    //     defaultConfig: true,
    //     valueRequired: defaultClientConfig[DefaultClientConfigParameter.AppId].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    // {
    //     index: defaultClientConfig[DefaultClientConfigParameter.AppKey].index,
    //     configName: DefaultClientConfigParameter.AppKey,
    //     configType: defaultClientConfig[DefaultClientConfigParameter.AppKey].configType,
    //     configValue: '',
    //     defaultConfig: true,
    //     valueRequired: defaultClientConfig[DefaultClientConfigParameter.AppKey].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
]

// 服务端默认配置表
export const defaultInternalConfig = {
    // [DefaultInternalConfigParameter.OuModule]: { index: 1, configType: Types.StringType, valueRequired: false },
    [DefaultInternalConfigParameter.SyncModule]: { index: 2, configType: Types.StringType, valueRequired: true },
    // [DefaultInternalConfigParameter.AuthModule]: { index: 3, configType: Types.StringType, valueRequired: false },
    // [DefaultInternalConfigParameter.SyncInterval]: { index: 4, configType: Types.NumberType, valueRequired: false },
    // [DefaultInternalConfigParameter.UserCreateStatus]: { index: 5, configType: Types.BooleanType, valueRequired: false },
    // [DefaultInternalConfigParameter.ValidPeriod]: { index: 6, configType: Types.NumberType, valueRequired: false },
}

// 服务端默认配置数组
export const defaultInternalConfigArray = [
    // {
    //     index: defaultInternalConfig[DefaultInternalConfigParameter.OuModule].index,
    //     configName: DefaultInternalConfigParameter.OuModule,
    //     configType: defaultInternalConfig[DefaultInternalConfigParameter.OuModule].configType,
    //     configValue: '',
    //     defaultConfig: true,
    //     valueRequired: defaultInternalConfig[DefaultInternalConfigParameter.OuModule].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    {
        index: defaultInternalConfig[DefaultInternalConfigParameter.SyncModule].index,
        configName: DefaultInternalConfigParameter.SyncModule,
        configType: defaultInternalConfig[DefaultInternalConfigParameter.SyncModule].configType,
        configValue: 'BaseSyncer',
        defaultConfig: true,
        valueRequired: defaultInternalConfig[DefaultInternalConfigParameter.SyncModule].valueRequired,
        nameValidateStatus: ValidateStatus.Normal,
        valueValidateStatus: ValidateStatus.Normal,
    },
    // {
    //     index: defaultInternalConfig[DefaultInternalConfigParameter.AuthModule].index,
    //     configName: DefaultInternalConfigParameter.AuthModule,
    //     configType: defaultInternalConfig[DefaultInternalConfigParameter.AuthModule].configType,
    //     configValue: '',
    //     defaultConfig: true,
    //     valueRequired: defaultInternalConfig[DefaultInternalConfigParameter.AuthModule].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    // {
    //     index: defaultInternalConfig[DefaultInternalConfigParameter.SyncInterval].index,
    //     configName: DefaultInternalConfigParameter.SyncInterval,
    //     configType: defaultInternalConfig[DefaultInternalConfigParameter.SyncInterval].configType,
    //     configValue: 1800,
    //     defaultConfig: true,
    //     valueRequired: defaultInternalConfig[DefaultInternalConfigParameter.SyncInterval].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    // {
    //     index: defaultInternalConfig[DefaultInternalConfigParameter.UserCreateStatus].index,
    //     configName: DefaultInternalConfigParameter.UserCreateStatus,
    //     configType: defaultInternalConfig[DefaultInternalConfigParameter.UserCreateStatus].configType,
    //     configValue: true,
    //     defaultConfig: true,
    //     valueRequired: defaultInternalConfig[DefaultInternalConfigParameter.UserCreateStatus].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
    // {
    //     index: defaultInternalConfig[DefaultInternalConfigParameter.ValidPeriod].index,
    //     configName: DefaultInternalConfigParameter.ValidPeriod,
    //     configType: defaultInternalConfig[DefaultInternalConfigParameter.ValidPeriod].configType,
    //     configValue: -1,
    //     defaultConfig: true,
    //     valueRequired: defaultInternalConfig[DefaultInternalConfigParameter.ValidPeriod].valueRequired,
    //     nameValidateStatus: ValidateStatus.Normal,
    //     valueValidateStatus: ValidateStatus.Normal,
    // },
]

/**
 * 检查配置，并对非法项添加对应的状态同时trim处理configName和configValue
 * @param config
 */
export const checkConfig = (config: ReadonlyArray<ConfigItem>): ReadonlyArray<ConfigItem> => {
    const arr = reduceRight(config, (result, item) => {
        const { configName, configValue, configType, defaultConfig, valueRequired } = item
        const trimConfigName = trim(configName)
        const trimConfigValue = (configType === Types.StringType ? trim(configValue) : configValue)
        const nameValidateStatus = (!defaultConfig && trimConfigName === '') ?
            ValidateStatus.Empty
            :
            trimConfigName.length > 800 ?
                ValidateStatus.NameToLong
                :
                findIndex(result, (i) => i.configName === trimConfigName) === -1 ?
                    ValidateStatus.Normal
                    :
                    ValidateStatus.NameRepeat

        return [
            ...result,
            {
                ...item,
                configName: trimConfigName,
                configValue: trimConfigValue,
                nameValidateStatus,
                valueValidateStatus: (valueRequired && trimConfigValue === '') ? ValidateStatus.Empty : (configType === Types.StringType && trimConfigValue.length > 800) ? ValidateStatus.ValueToLong : ValidateStatus.Normal,
            },
        ]
    }, [])

    return arr.slice().reverse()
}

/**
 * 检查字符串是否为 json 对象
 * @param {*} str 字符串
 * @returns boolean
 */
export function isJsonObject(str: string): boolean {
    try {
        // 如果可以转换，则是合法的json格式，如果转换失败，则catch到错误，返回false
        const authorizedJson = JSON.parse(str)
        // 判断转出来的对象，而非其他类型
        return (typeof authorizedJson === 'object' && !(authorizedJson instanceof Array) && authorizedJson !== null)
    } catch (error) {
        return false
    }
}

/**
 * json字符串（config,internalConfig）转换为对象数组
 * @param config 待转换的配置参数字符串
 */
export function transformJsonToArray(config: string, terminal?: Terminal, stage?: Stage): ReadonlyArray<ConfigItem> {
    // TODO:此处未完整处理升级场景
    // 当需要考虑升级时，除了判断config是不是对象格式的json,非对象格式JSON用默认数组填充，还需要考虑原来使用对象格式JSON但只填写了部分默认参数时，其他默认项自动补充
    if (isJsonObject(config)) {
        let arr = []
        const objConfig = JSON.parse(config)
        for (const key in objConfig) {
            // 高级配置时，当设置的数字大于最大限制时候，自动还原为最大数
            const numberType = stage === Stage.Advanced && typeof objConfig[key] === Types.NumberType
            arr = (terminal === Terminal.Client ?
                [
                    ...arr,
                    {
                        index: (defaultClientConfig[key] && defaultClientConfig[key].index) || null,
                        configName: key,
                        configType: (stage === Stage.Initialize && defaultClientConfig[key] && defaultClientConfig[key].configType) || typeof objConfig[key],
                        configValue: (numberType && objConfig[key] > MAX_SAFE_INTEGER) ?
                            MAX_SAFE_INTEGER
                            :
                            (numberType && objConfig[key] < MIN_SAFE_INTEGER) ?
                                MIN_SAFE_INTEGER
                                :
                                objConfig[key],
                        defaultConfig: !!defaultClientConfig[key],
                        valueRequired: !!(defaultClientConfig[key] && defaultClientConfig[key].valueRequired),
                        nameValidateStatus: ValidateStatus.Normal,
                        valueValidateStatus: ValidateStatus.Normal,
                    },
                ] :
                [
                    ...arr,
                    {
                        index: (defaultInternalConfig[key] && defaultInternalConfig[key].index) || null,
                        configName: key,
                        configType: (stage === Stage.Initialize && defaultInternalConfig[key] && defaultInternalConfig[key].configType) || typeof objConfig[key],
                        configValue: (numberType && objConfig[key] > MAX_SAFE_INTEGER) ?
                            MAX_SAFE_INTEGER
                            :
                            (numberType && objConfig[key] < MIN_SAFE_INTEGER) ?
                                MIN_SAFE_INTEGER
                                :
                                objConfig[key],
                        defaultConfig: !!defaultInternalConfig[key],
                        valueRequired: !!(defaultInternalConfig[key] && defaultInternalConfig[key].valueRequired),
                        nameValidateStatus: ValidateStatus.Normal,
                        valueValidateStatus: ValidateStatus.Normal,
                    },
                ])
        }
        return arr
    } else {
        return terminal === Terminal.Client ? defaultClientConfigArray : defaultInternalConfigArray
    }
}

/**
 * 将数组转换成对象
 * @param {*} configArray 配置参数数组
 * @returns 配置参数数组转换成的对象
 */
export function transformArrayToObject(configArray: ReadonlyArray<ConfigItem>): Record<string, any> {
    let obj = {}
    configArray.forEach(({ configName, configValue }) => {
        obj[configName] = configValue
    })
    return obj
}