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
 * 服务端参数配置默认配置项
 */
export enum DefaultInternalConfigParameter {
    /**
     * 消息推送模块名称
     */
    MsgModule = 'msgModule',
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
     * 名称不合法（输入的名称为msg_type_list）
     */
    NameInvalid = 2,

    /**
     * 名称重复
     */
    NameRepeat = 3,

    /**
     * 接口报错
     */
    InterfaceError = 4,

    /**
    * Name长度超过限制（128）或者有特殊字符
    */
    NameLengthExceedsLimitOrHasInvalid = 6,

    /**
    * 名称长度超过限制，主要用于限定客户端和服务端参数名称和值的长度，安全考虑
    */
    NameToLong = 7,

    /**
     * 值长度超过限制，主要用于限定客户端和服务端参数名称和值的长度，安全考虑
     */
    ValueToLong = 8,

    /**
     * 消息类型值长度超过限制
     */
    MessagesValueToLong = 9,

    /**
     * 消息类型值不合法
     */
    MessagesValueInvalid = 10,

    /**
     * 第三方消息服务名称重复
     */
    ThirdPartyNameRepeat = 11,

    /**
     * 插件类名长度超过限制或者有特殊字符
     */
    PluginClassNameValueToLongOrHasInvalid = 12,
}

/**
 * 输入验证提示
 */
export const ValidateMessage = {
    [ValidateStatus.Empty]: __('此输入项不允许为空。'),
    [ValidateStatus.NameRepeat]: __('该名称已存在，请您重新输入。'),
    [ValidateStatus.NameLengthExceedsLimitOrHasInvalid]: __('消息服务名称不能包含 \\ / : * ? " < > |，长度不能超过128个字符。'),
    [ValidateStatus.PluginClassNameValueToLongOrHasInvalid]: __('插件类名名称不能包含 \\ / : * ? " < > |，长度不能超过128个字符。'),
    [ValidateStatus.NameToLong]: __('名称长度不能超过800个字符。'),
    [ValidateStatus.ValueToLong]: __('值长度不能超过800个字符。'),
    [ValidateStatus.MessagesValueToLong]: __('值长度不能超过200个字符。'),
    [ValidateStatus.MessagesValueInvalid]: __('值不能包含//'),
    [ValidateStatus.ThirdPartyNameRepeat]: __('消息服务名称已存在，请重新输入。'),
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
 * 一条配置参数信息，用于数组填充
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

export interface MessageConfigItem {
    /**
     * 配置参数值
     */
    configValue: any;

    /**
    * 配置参数值验证状态
    */
    valueValidateStatus: ValidateStatus;
}

/**
 * 接口返回数据经转换后的配置参数数据
 */
export interface ConfigInfo {
    /**
     * 第三方消息服务名称
     */
    thirdPartyName: string;

    /**
     * 第三方认证服务启用状态
     */
    enabled: boolean;

    /**
     * 插件类名
     */
    pluginClassName: string;

    /**
     * 插件参数配置
     */
    internalConfig: string;

    /**
     * 消息类型
     */
    messages: ReadonlyArray<MessageConfigItem>;

    /**
     * 插件
     */
    plugin: null | PluginInfo;

    /**
     * 第三方认证服务唯一索引
     */
    indexId: null | number;
}

/**
 * 高级配置接收所有配置项参数
 */
export interface Params {
    /**
     * 第三方app名字
     */
    thirdparty_name: string;

    /**
     * 第三方配置开关
     */
    enabled: boolean;

    /**
     * 消息插件类名
     */
    class_name: string;

    /**
     * 消息类型
     */
    channels: ReadonlyArray<MessageConfigItem>;

    /**
     * 插件需要其他配置，透传给第三方插件
     */
    config: ReadonlyArray<ConfigItem>;
}

/**
 * 接口传参
 */
export interface InterfaceParams {
    /**
     * 第三方app名字
     */
    thirdparty_name: string;

    /**
     * 第三方配置开关
     */
    enabled: boolean;

    /**
     * 消息插件类名
     */
    class_name: string;

    /**
     * 消息类型
     */
    channels: ReadonlyArray<string>;

    /**
     * 插件需要其他配置，透传给第三方插件
     */
    config: Record<string, any>;
}

/**
 * 第三方插件信息
 */
export interface PluginInfo {
    /**
     * 文件名
     */
    filename: string;
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
     * 消息为空
     */
    MessageEmpty = 1,

    /**
     * 格式错误
     */
    WrongFormat = 2,

    /**
     * 第三方应用已不存在
     */
    AppNotExist = 3,

    /**
     * 无法连接服务
     */
    ServerError = 4,

    /**
     * 未配置存储
     */
    NoStorage = 21104,

    /**
     * 不合法的文件名
     */
    IllegalFileName = 20616,

    /**
     * 未知错误导致插件上传失败
     */
    UnknownError = 20000,
}

/**
 * 转换错误提示
 */
export function formatterError(errID: ErrorCode, errMsg?: string): string {
    switch (errID) {
        case ErrorCode.MessageEmpty:
            return __('请您至少填写一种消息类型。')

        case ErrorCode.NoStorage:
            return __('系统存储异常，请检查存储配置。')

        case ErrorCode.WrongFormat:
            return __('请上传规定格式的文件。')

        case ErrorCode.IllegalFileName:
            return __('文件名不能包含 \\ / @ # $ % ^ & * ( ) [ ] 字符，且首位不能使用 + - .字符，长度不能超过255个字符，请重新上传。')

        case ErrorCode.UnknownError:
            return __('未知错误，上传失败。')

        case ErrorCode.ServerError:
            return __('无法连接服务，请稍后再试。')

        case ErrorCode.AppNotExist:
            return __('保存失败，该第三方应用已不存在。')

        default:
            return errMsg || __('未知的错误码。')
    }
}

/**
 * 服务端默认配置表
 */
export const defaultInternalConfig = {
    [DefaultInternalConfigParameter.MsgModule]: { index: 1, configType: Types.StringType, valueRequired: true },
}

/**
 * 服务端默认配置数组
 */
export const defaultInternalConfigArray = [
    {
        index: defaultInternalConfig[DefaultInternalConfigParameter.MsgModule].index,
        configName: DefaultInternalConfigParameter.MsgModule,
        configType: defaultInternalConfig[DefaultInternalConfigParameter.MsgModule].configType,
        configValue: '',
        defaultConfig: true,
        valueRequired: defaultInternalConfig[DefaultInternalConfigParameter.MsgModule].valueRequired,
        nameValidateStatus: ValidateStatus.Normal,
        valueValidateStatus: ValidateStatus.Normal,
    },
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
                valueValidateStatus: (configType === Types.StringType && trimConfigValue.length > 800) ? ValidateStatus.ValueToLong : ValidateStatus.Normal,
            },
        ]
    }, [])

    return arr.slice().reverse()
}

/**
 * 检查消息类型配置是否正确
 */
export const checkMessages = (messages: ReadonlyArray<MessageConfigItem>): { newMessages: ReadonlyArray<MessageConfigItem>; validateIndex: number } => {
    let validateIndex = -1;

    const newMessages = messages.map((item, index) => {
        const { configValue } = item;
        const trimConfigValue = configValue.trim();

        const newItem = {
            configValue: trimConfigValue,
            valueValidateStatus:
                trimConfigValue ?
                    trimConfigValue.length > 200 ?
                        ValidateStatus.MessagesValueToLong
                        : /\/{2,}/.test(trimConfigValue) ?
                            ValidateStatus.MessagesValueInvalid
                            : ValidateStatus.Normal
                    : ValidateStatus.Empty,

        };

        if (validateIndex === -1) {
            validateIndex = newItem.valueValidateStatus === ValidateStatus.Normal ? -1 : index;
        }

        return newItem;
    })

    return {
        newMessages,
        validateIndex,
    }
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
export function transformJsonToArray(config: string, stage: Stage): ReadonlyArray<ConfigItem> {
    // TODO:此处未完整处理升级场景
    // 当需要考虑升级时，除了判断config是不是对象格式的json,非对象格式JSON用默认数组填充，还需要考虑原来使用对象格式JSON但只填写了部分默认参数时，其他默认项自动补充
    if (isJsonObject(config)) {
        let arr = []
        const objConfig = JSON.parse(config)
        for (const key in objConfig) {
            // 高级配置时，当设置的数字大于最大限制时候，自动还原为最大数
            const numberType = stage === Stage.Advanced && typeof objConfig[key] === Types.NumberType
            arr = [
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
                    defaultConfig: false,
                    valueRequired: false,
                    nameValidateStatus: ValidateStatus.Normal,
                    valueValidateStatus: ValidateStatus.Normal,
                },
            ]
        }
        return arr
    } else {
        return defaultInternalConfigArray
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