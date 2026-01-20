import { BigNumber } from 'bignumber.js'

/**
 * 验证数字
 * @params input 输入值
 * @return 返回是否是数字
 */
export function number(input: any): boolean {
    return /^[\-]?[0-9]+(\.[0-9])?$/.test(String(input));
}

/**
 * 验证自然数
 * @params input 输入值
 * @return 返回是否是自然数
 */
export function natural(input: any): boolean {
    return /^[0-9]+$/.test(String(input));
}

/**
 * 验证正数（包括0）
 * @params input 输入值
 * @return 返回是否是正数
 */
export function positive(input: any): boolean {
    return /^[0-9]+(\.[0-9]+)?$/.test(String(input));
}

/**
 * 验证正整数
 * @params input 输入值
 * @return 返回是否是正整数
 */
export function positiveInteger(input: any): boolean {
    return /^[1-9]\d*$/.test(String(input));
}

/**
 * 验证邮箱
 * @params input 输入值
 * @return 返回是否是邮箱
 */
export function mail(input: any): boolean {
    return /^[\w\-]+(\.[\w\-]+)*@[\w\-]+(\.[\w\-]+)+$/.test(input);
}

/**
 * 验证用户显示名 显示名不能包含\ / : * ? " < > | 特殊字符
 * @params input 输入值
 * @return 返回显示名是否合法
 */
export function dispalyName(input: any): boolean {
    return !/[/\\:*?"<>|]/.test(input);
}

/**
 * 验证是否是时间
 * @params input 输入值
 * @return 返回是否是时间
 */
export function clock(input: any): boolean {
    return natural(input) && input < 60;
}

/**
 * 验证是否超过字数限制
 * @params input 输入值
 * @return 是否超出限制
 */
export function tweet(value: any) {
    return value.length <= 140;
}

/**
 * 限制输入最大长度
 * @params input 输入值
 * @return 是否超出限制
 */
export function maxLength(maxLength: number, trim: boolean = true) {
    return function (input) {
        input = String(input);
        return (trim ? input.trim() : input).length <= maxLength;
    }
}

/**
 * 限制输入范围(自然数)
 */
export function range(from: number, to: number) {
    return function (input) {
        return natural(input) && input >= from && input <= to;
    }
}

/**
 * 颜色
 */
export function validateColor(input: string | number) {
    return /^#?[0-9A-Fa-f]{6}$/.test(String(input))
}

/**
 * 子网掩码
 */
export function subNetMask(input: string | number) {
    return !!String(input) && /^((128|192|224|240|248|252|254|255)\.0\.0\.0)|(255\.(((0|128|192|224|240|248|252|254|255)\.0\.0)|(255\.(((0|128|192|224|240|248|252|254|255)\.0)|255\.(0|128|192|224|240|248|252|254|255)))))$/.test(String(input))
}

/**
 * IPv4
 */
export function IP(input: string | number) {
    return !!String(input) && /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(String(input))
}

/**
 * IPv6
 */
export function IPV6(input: string | number) {
    const inputString = String(input);

    return !!inputString &&
        /^\s*(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\s*$/.test(inputString);
}

/**
 * IPv6前缀
 */
export function IPV6Prefix(input: string | number) {
    const inputNumber = Number(input)

    return !Number.isNaN(inputNumber) && Number.isFinite(inputNumber) &&
        0 <= inputNumber && inputNumber <= 128
}

/**
 * 验证正整数并且有位数限制
 * @params input 输入值
 * @return 返回是否是正整数
 */
export function positiveIntegerAndMaxLength(maxLength: number) {
    return function (input) {
        input = String(input);
        return positiveInteger(input) && input.length <= maxLength;
    }
}

/**
 * 验证邮箱格式并且有长度限制
 * @params input 输入值
 * @return 返回是否是正整数
 */
export function mailAndLenth(input, minLength, maxLength): boolean {
    return /^[\w\-]+(\.[\w\-]+)*@[\w\-]+(\.[\w\-]+)+$/.test(input) && minLength < input.length && maxLength > input.length;

}

export function isURL(str) {
    const urlRegex = '^(?!mailto:)(?:(?:http|https|ftp)://)(?:\\S+(?::\\S*)?@)?(?:(?:(?:[1-9]\\d?|1\\d\\d|2[01]\\d|22[0-3])(?:\\.(?:1?\\d{1,2}|2[0-4]\\d|25[0-5])){2}(?:\\.(?:[0-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-4]))|(?:\\[)((?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,7}:|(?:[0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,5}(?::[0-9a-fA-F]{1,4}){1,2}|(?:[0-9a-fA-F]{1,4}:){1,4}(?::[0-9a-fA-F]{1,4}){1,3}|(?:[0-9a-fA-F]{1,4}:){1,3}(?::[0-9a-fA-F]{1,4}){1,4}|(?:[0-9a-fA-F]{1,4}:){1,2}(?::[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:(?:(?::[0-9a-fA-F]{1,4}){1,6})|:(?:(?::[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(?::[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(?:ffff(?::0{1,4}){0,1}:){0,1}(?:(?:25[0-5]|(?:2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3}(?:25[0-5]|(?:2[0-4]|1{0,1}[0-9]){0,1}[0-9])|(?:[0-9a-fA-F]{1,4}:){1,4}:(?:(?:25[0-5]|(?:2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3}(?:25[0-5]|(?:2[0-4]|1{0,1}[0-9]){0,1}[0-9]))(?:\\])|(?:(?:[a-z\\u00a1-\\uffff0-9]+-?)*[a-z\\u00a1-\\uffff0-9]+)(?:\\.(?:[a-z\\u00a1-\\uffff0-9]+-?)*[a-z\\u00a1-\\uffff0-9]+)*(?:\\.(?:[a-z\\u00a1-\\uffff]{2,})))|localhost)(?::\\d{2,5})?(?:(/|\\?|#)[^\\s]*)?$'
    const url = new RegExp(urlRegex, 'i')
    return str.length < 2083 && url.test(str);
}

/**
 * 匹配域名
 * @params input 输入值 (只能包含 英文、数字 及 -. 字符，长度范围 3~100 个字符)
 */
export function isDomain(input): boolean {
    return /^[a-zA-Z0-9\-(\.)?]+([a-zA-Z])+$/.test(input) && input.length >= 3 && input.length <= 100;
}

/**
 * 验证Mac地址
 * @export
 * @param {string} input 输入值（只能是由6组数字和字母（不区分大小写），每组一个数字一个字母组成，并且组与组之间使用-进行连接）
 * @returns {boolean}
 */
export function isMac(input: string): boolean {
    return /^[A-Fa-f0-9]{2}(\-[A-Fa-f0-9]{2}){5}$/.test(input);
}

/**
 * 验证文件后缀是否正确
 * @export
 * @param {string} input 输入值（不允许包含.\|/*?"<>:）
 * @returns {boolean}
 */
export function isSuffix(input: string): boolean {
    return /^\.([^\.\\\|\/\*\?"<>:])+$/.test(input);
}

/**
 * 匹配英文和数字
 */
export function isLetterOrNumber(input: string): boolean {
    return /^[a-zA-Z0-9]+$/.test(input);
}

/*
 * 验证文本输入
 */
export function isUserName(input: string) {
    return /^[^\/\\:*?"<>|\s]+$/.test(input);
}

/**
 * 验证授权码
 */
export function validLicense(input): boolean {
    return /^[A-Z0-9]{5}(\-[A-Z0-9]{5}){5}$/.test(input);
}

/**
 * 验证激活码
 */
export function validActiveCode(input): boolean {
    return /^[A-Z0-9]+$/.test(input);
}

/**
 * 验证手机号码
 */
export function cellphone(value) {
    return /^[\d]{11}$/i.test(value)
}

/**
 * 验证电话号码(包括手机和座机)  手机号码只能包含 数字，长度范围 1~20 个字符
 */
export function phoneNum(value) {
    return /^[\d][\d- ]{0,20}$/i.test(value);
}

export function decimal(input: any): boolean {
    return /^[0-9]+(\.[0-9]{0,2})?$/.test(String(input));
}

/**
 * 验证端口号  1~65535 之间的整数
 * @param val 端口号
 */
export function isPort(val: any): boolean {
    return (/^[1-9]\d{0,4}$/.test(String(val)) && parseInt(String(val)) <= 65535)
}

/**
 * 验证域名别名
 * @param display_name 域名别名
 */
export function isDisplayName(display_name: string): boolean {
    return !(/[\s\/\\:\*\?"<>\|()]/.test(display_name));
}

/**
 * 验证身份证
 * @params input 输入值
 * @return 返回是否是身份证
 */
export function idcard(input: any): boolean {
    return /^\d{17}[X0-9]$/.test(String(input));
}

/**
 * 验证身份证(包含港澳台等)
 * @params input 输入值
 * @return 返回是否是身份证
 */
export function variousIdCard(input: any): boolean {
    return /^[A-Za-z0-9/()-]{8,18}$/i.test(String(input));
}

/**
 * 验证域名
 * @param val 域名 (域名只能包含 英文、数字 及 -. 字符，每一级不能以“-”字符开头或结尾，每一级长度必需 1~63 个字符，且总长不能超过253个字符。)
 */
export function isDomainName(val: any): boolean {
    return /^(?=^.{3,253}$)(([a-zA-Z0-9]{1,63}|[a-zA-Z0-9][-a-zA-Z0-9]{0,61}[a-zA-Z0-9])\.)*([a-zA-Z0-9]{1,63}|[a-zA-Z0-9][-a-zA-Z0-9]{0,61}[a-zA-Z0-9])$/.test(String(val));
}

/**
 * 判断字符串是否可转换为json格式
 * @param str 字符串
 */
export function isJSON(str) {
    if (typeof str === 'string') {
        try {
            JSON.parse(str)
            return true
        } catch (e) {
            return false
        }
    }
}

/**
 * 验证名称
 * @param input 只允许输入中文英文数字和特殊字符~!%#$@-_.及空格
 */
export function isName(input: any): boolean {
    return /^[\u4E00-\u9FA5A-Za-z0-9~!%#$@\-_\.\s]+$/.test(String(input))
}

/**
 * 验证用户组织管理，用户名
 * @param input 不能包含 空格 或 \ / : * ? " < > | 特殊字符，长度不能超过128个字符。
 */
export function isLoginName(input: any): boolean {
    return /^[^\/\\:*?"<>|\s]{1,128}$/.test(String(input))
}

/**
 * 验证用户组织管理，用户名
 * @param input 不能包含 空格 或 \ / : * ? " < > | 特殊字符，长度不能超过128个字符。
 */
export function isUserLoginName(input: any): boolean {
    return /^[^\/\\*?"<>|\s]{1,128}$/.test(String(input))
}

/**
 * 验证用户组织管理，显示名，备注
 * @param input 不能包含 \ / : * ? " < > | 特殊字符，长度不能超过128个字符。
 */
export function isNormalName(input: any): boolean {
    return /^[^\/\\:*?"<>|]{1,128}$/.test(String(input))
}

/**
 * 验证用户组织管理，显示名，备注
 * @param input 不能包含 \ / : * ? " < > | 特殊字符，长度不能超过128个字符。
 */
export function isUserNormalName(input: any): boolean {
    return /^[^\/\\*?"<>|]{1,128}$/.test(String(input))
}

/**
 * 验证角色名称
 * @param input 不能包含 \ / : * ? " < > | 特殊字符，长度不能超过50个字符。
 */
export function isNormalShortName(input: any): boolean {
    return /^[^\/\\:*?"<>|]{1,50}$/.test(String(input))
}

/**
 * @param input 验证岗位
 * @param input 长度不能超过50个字符。
 */
export function isNormalPosition(input: any): boolean {
    return /^.{0,50}$/.test(String(input))
}

/**
 * 验证用户编码
 */
export function isNormalCode(input: any): boolean {
    return /^[a-zA-Z0-9_-]{1,255}$/.test(String(input))
}

/**
 * 验证用户组织管理，配额空间
 * @param input 配额空间值为不超过 1000000 的正数，支持小数点后两位。
 */
export function isSpace(input: any): boolean {
    return (Number(input) === Number(Number(input).toFixed(2)) && Number(input) <= 1000000 && Number(input) > 0)
}

/**
 * 验证是否为Dom节点
 */
export function isDom(node: HTMLElement): boolean {
    // 首先判断是否支持HTMLElement
    return (
        (typeof HTMLElement === 'function') ?
            (node instanceof HTMLElement)
            : (node && (typeof node === 'object') && (node.nodeType === 1) && (typeof node.nodeName === 'string'))
    )

}

/**
 * emoji表情包正则表达式
 */
export const emojiRule = /[\uD83C|\uD83D|\uD83E][\uDC00-\uDFFF][\u200D|\uFE0F]|[\uD83C|\uD83D|\uD83E][\uDC00-\uDFFF]|[0-9|*|#]\uFE0F\u20E3|[0-9|#]\u20E3|[\u203C-\u3299]\uFE0F\u200D|[\u203C-\u3299]\uFE0F|[\u2122-\u2B55]|\u303D|[\A9|\AE]\u3030|\u00A9|\u00AE|\u3030/ig

/**
 * 检查是否为boolean (例如, true, false)
 * @param value 被检查的值
 * @return 如果是boolean，返回true；否则返回false
 */
export function isBoolean(value: any): boolean {
    return typeof value === 'boolean'
}

/**
 * 检查是否为函数 (例如, function fn() {})
 * @param value 被检查的值
 * @return 如果是函数，返回true；否则返回false
 */
export function isFunction(value: any): boolean {
    return typeof value === 'function'
}

/**
 * 检查是否为对象 (例如, [], {}, new Number(0), and new String(''))
 * @param value 被检查的值
 * @return 如果是对象，返回true；否则返回false
 */
export function isObject(value: any): boolean {
    return value !== null && typeof value === 'object'
}

/**
 * 判断BigNumber对象是否相等
 * @param value checkValue 判断的值
 * @param checkValue 判断的另一个值
 * @return 如果相等，返回true；否则返回false
 */
export function isBigNumberEqual(value: BigNumber | number, checkValue: BigNumber | number) {
    return value.toString() === checkValue.toString()
}