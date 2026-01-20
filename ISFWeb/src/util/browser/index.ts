import { kebabCase, has, get, last, find,  every, assign } from 'lodash';
import session from '../session';

export interface UserAgent {
    // 系统架构 32/64 位
    platform?: number;

    // 浏览器详细版本
    version?: number;

    // 浏览器型号
    app?: number;

    // 是否是移动设备
    mobile?: boolean;

    // 操作系统
    os?: any;
}

// declare const { safari };
const safari = undefined;

/**
 * AnyShare支持的语言
 */
enum SupportedLanguage {
    ZH_CN = 'zh-cn',
    ZH_TW = 'zh-tw',
    EN_US = 'en-us',
}

/**
 * 浏览器类型
 */
export enum Browser {
    MSIE,
    Edge,
    Safari,
    Firefox,
    Chrome,
    WeChat,
}

/**
 * 系统类型
 */
export enum OSType {
    Windows,
    Mac,
    Android,
    IOS,
    WindowsPhone,
}

/**
 * 打开新窗口
 */
export function open(url: string): Window {
    if (typeof safari !== 'undefined' && has(safari, 'application.browserWindows.openTab')) {
        return safari && safari.application.browserWindows.openTab(url);
    } else {
        // IE 8 需要将“安全设置”中的“自动提示下载”设置为启用，否则会提示下载被阻止
        return window.open(url);
    }
}

/**
 * 是否使用了HTTPS
 */
export function useHTTPS(): boolean {
    return location.protocol === 'https:';
}

/**
 * 获取UserAgent信息
 */
export function userAgent(): UserAgent {
    const userAgent = window.navigator.userAgent;
    const os = (() => {
        switch (true) {
            case /Iphone|Ipad|Ipod/i.test(userAgent):
                return OSType.IOS;

            case /Android/i.test(userAgent):
                return OSType.Android;

            default:
                return undefined;
        }
    })();
    const app = (() => {
        switch (true) {
            // 微信中含有Safari
            case /MicroMessenger/i.test(userAgent):
                return Browser.WeChat;

            case /Edge/i.test(userAgent):
                return Browser.Edge;

            case /Trident/i.test(userAgent):
                return Browser.MSIE;

            // Chrome的UA包含Safari字样，Safari的UA不包含Chrome字样，因此要区分Chrome还是Safari必须先判断Chrome字样
            case /Chrome\/[\d\.]+/i.test(userAgent):
                return Browser.Chrome;

            case /Safari\/[\d\.]+$/i.test(userAgent):
                return Browser.Safari;

            case /Firefox\/[\d\.]+$/i.test(userAgent):
                return Browser.Firefox;
        }
    })();
    const platform = /win64|wow64/i.test(userAgent.toLowerCase()) || navigator.platform === 'MacIntel' ? 64 : 32;
    const mobile = /Mobile/i.test(navigator.userAgent);
    const version = (() => {
        let verMatch;

        switch (app) {
            case Browser.MSIE:
                switch (true) {
                    case userAgent.indexOf('Trident/4.0') !== -1:
                        return 8;

                    case userAgent.indexOf('Trident/5.0') !== -1:
                        return 9;

                    case userAgent.indexOf('Trident/6.0') !== -1:
                        return 10;

                    case userAgent.indexOf('Trident/7.0') !== -1:
                        return 11;

                    default:
                        return undefined;
                }

            case Browser.Edge:
                verMatch = userAgent.match(/Edge\/([\d\.]+)/i);
                return verMatch ? Number(verMatch[1]) : undefined;

            case Browser.Safari:
                verMatch = userAgent.match(/Safari\/([\d\.]+)/i);
                return verMatch ? Number(verMatch[1]) : undefined;

            case Browser.Firefox:
                verMatch = userAgent.match(/Firefox\/([\d\.]+)/i);
                return verMatch ? Number(verMatch[1]) : undefined;

            case Browser.Chrome:
                verMatch = userAgent.match(/Chrome\/([\d\.]+)/i);
                return verMatch ? Number(verMatch[1]) : undefined;

            case Browser.WeChat:
                verMatch = userAgent.match(/MicroMessenger\/([\d\.]+)/i);
                return verMatch ? Number(verMatch[1]) : undefined;
        }
    })();

    return {
        app,
        mobile,
        version,
        platform,
        os,
    }
}

/**
 * 判断浏览器是否匹配查询结果
 */
export function isBrowser(query: UserAgent): boolean {
    const ua = userAgent();

    return every(Object.entries(query), ([k, v]) => ua[k] === v);
}

/**
 * 检测是否是IE8
 */
export function isIE8() {
    return isBrowser({ app: Browser.MSIE, version: 8 })
}

/**
 * 检测是否是IE11
 */
export function isIE11() {
    return isBrowser({ app: Browser.MSIE, version: 11 })
}

/**
 * 获取页面的origin
 */
export function origin() {
    return location.origin || `${location.protocol}//${location.host}`;
}

/**
 * 映射keyCode到字符串
 * Safari, XP下的Chrome对部分按键返回键码而不是ASCII码
 */
export function mapKeyCode(input: number) {
    const KeyCodeMap = {
        '186': ';',
        '187': '=',
        '188': ',',
        '189': '-',
        '190': '.',
        '191': '/',
        '192': '`',
        '219': '[',
        '220': '\\',
        '221': ']',
        '222': '\'',
    }

    return KeyCodeMap[input] || String.fromCharCode(input);
}

/**
 * 获取当前脚本的绝对路径
 */
function getScriptSrc(): string {
    // Chrome / FF / IE10+
    if (document.currentScript) {
        return document.currentScript.src;
    }

    // IE 10-
    return get(find(document.scripts, (script) => script.readyState === 'interactive'), 'src');
}

/**
 * 获取当前脚本的目录
 */
function getScriptDir(): string {
    const [file, dir] = getScriptSrc().match(/(.+)\/.+$/);
    return dir;
}

/**
 * 获取当前执行中的脚本相对的路径
 */
export function toRelative(relative: string): string {
    return `${getScriptDir()}/${relative}`;
}

/**
 * 加载css
 */
export function loadCSS(href, { rel = 'stylesheet', type = 'text/css', media = '' } = {}) {
    let link = document.createElement('link')
    assign(link, { href, rel, type, media })
    document.querySelector('head').appendChild(link)
}

/**
 * 加载js
 */
export function loadScript(src, { type = 'text/javascript' } = {}) {
    let script = document.createElement('script')
    assign(script, { src, type })
    document.querySelector('head').appendChild(script)
}

/**
 * 根据文件后缀加载js或css
 */
export function load(url) {
    const ext = last(url.split('.'))

    switch (ext) {
        case 'js':
            return loadScript(url)
        case 'css':
            return loadCSS(url)
    }
}

/**
 * 移除事件
 * @param target 绑定事件的DOM
 * @param event 事件名称
 * @param listener 事件处理函数
 */
export function unbindEvent<T extends Event>(target: Window | HTMLDocument | HTMLElement, event: string, listener: (event: T) => any) {
    if (isBrowser({ app: Browser.MSIE, version: 8 })) {
        target.detachEvent('on' + event, listener)
    } else {
        target.removeEventListener(event, listener)
    }
}

/**
 * 绑定事件
 * @param target 绑定事件的DOM
 * @param event 事件名称
 * @param listener 事件处理函数
 * @return 返回移除事件的函数
 */
export function bindEvent<T extends Event>(target: Window | HTMLDocument | HTMLElement, event: string, listener: (event: T) => any): () => void {
    if (isBrowser({ app: Browser.MSIE, version: 8 })) {
        target.attachEvent('on' + event, listener)
    } else {
        target.addEventListener(event, listener)
    }

    return unbindEvent.bind(null, target, event, listener)
}

/**
 * 插入样式声明
 * @param rules 样式声明
 * @param parent 插入目标位置
 */
export function insertStyle(rules: Record<string, any>, parent: HTMLElement = document.querySelector('head'), { media = '' } = {}) {
    const style = document.createElement('style');
    if (media) {
        assign(style, { media });
    }
    parent.appendChild(style);

    const styleSheet = style.sheet

    let index = styleSheet.cssRules.length - 1

    Object.entries(rules).forEach(([selector, declartons]) => {
        index++
        const declartonString = Object.entries(declartons).map(([key, value]) => `${kebabCase(key)}:${value};`).join('');

        styleSheet.insertRule ?
            styleSheet.insertRule(`${selector} { ${declartonString} }`, index) :
            styleSheet.addRule(selector, declartonString)
    });

    // 移除rules，兼容插件ie11
    return function () {
        parent.removeChild(style)
    }
}

/**
 * 获取浏览器当前语言临时解决方案
 */
export function envLanguage() {
    const getLanguageHash = () => {
        const [match, language] = /\blang=([a-zA-Z\-]+)\b/.exec(window.location.hash) || [undefined, undefined];

        if (language) {
            session.set('lang', language.toLowerCase());

            return language;
        }
    }

    return (getLanguageHash() || session.get('lang') || window.navigator.language || window.navigator.browserLanguage).toLowerCase()

}

/**
 * 获取适配后的语言环境
 * @param lang 输入的语言
 * @return 返回适配后支持的语言
 */
export function getSupportedLanguage(lang = ''): SupportedLanguage {
    const [language, region] = lang.split(/[\-_]/).map((x) => x && x.toLowerCase());

    if (language === 'zh') {
        switch (region) {
            case 'cn':
                return SupportedLanguage.ZH_CN;

            case 'tw':
                return SupportedLanguage.ZH_TW;

            default:
                return SupportedLanguage.ZH_CN;

        }
    } else {
        return SupportedLanguage.EN_US;
    }
}

/**
 * 限制字符串长度为浏览器允许的最大长度
 * @param str 输入字符串
 * @param maxEncodedLength 允许的urlencoded后的最大长度
 * @example
 * ```
 *      substrByEncodedLength('Web客户端HTTPS证书安装' 60) // 'Web客户端HTTPS证书'
 * ```
 */
export function substrByEncodedLength(str: string, maxEncodedLength: number): string {
    let subIndex = 0;

    for (let encodedLength = 0; subIndex < str.length; subIndex++) {
        encodedLength += encodeURIComponent(str[subIndex]).length;

        if (encodedLength > maxEncodedLength) {
            break;
        }
    }

    return str.substr(0, subIndex);
}

/**
 * 全屏
 * @param element
 */
export function requestFullScreen(element: HTMLElement | Document) {
    (element.requestFullscreen || element.webkitRequestFullScreen || element.mozRequestFullScreen || element.msRequestFullscreen).call(element)

}

/**
 * 退出全屏
 */
export function exitFullScreen(el = document) {
    (el.exitFullscreen || el.webkitExitFullscreen || el.webkitCancelFullScreen || el.mozExitFullScreen || el.mozCancelFullScreen || el.msExitFullscreen).call(el)
}

/**
 * 设置tab页title
 * @param title
 */
export function setTitle(title) {
    // 当页面中嵌入了flash，并且页面的地址中含有“片段标识”（即网址#之后的文字）IE标签页标题被自动修改为网址片段标识
    if (userAgent().app === Browser.MSIE || userAgent().app === Browser.Edge) {
        setTimeout(function () {
            document.title = title
        }, 1000)
    } else {
        document.title = title
    }
}

/**
 * 设置favicon
 */
export function setFavicon(favicon) {
    const link = document.createElement('link');
    link.rel = 'shortcut icon';
    link.type = 'image/x-icon';
    link.href = favicon;
    document.querySelector('head').appendChild(link)
}
