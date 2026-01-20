/**
 * 语言枚举
 */
export enum Lang {
    /**
     * 简体中文
     */
    zh = 'zh-cn',

    /**
     * 繁体中文
     */
    tw = 'zh-tw',

    /**
     * 英文
     */
    us = 'en-us',
}

/**
 * 协议枚举
 */
export enum Protocol {
    Https = 'https:',

    Http = 'http:'
}

export interface CommonType {
    /**
     * 协议
     */
    protocol: Protocol;

    /**
     * 域名或ip
     */
    host: string;

    /**
     * 端口
     */
    port: string | number;

    /**
     * 访问前缀
     */
    prefix: string;
    
    /**
     * 语言
     */
    lang: Lang;

    /**
     * 主题色（一些按钮/tab 需要有主题色背景或边框颜色）
     */
    theme: string;

    config: {
        getTheme: {
            normal: string;
            hover: string;
            active: string;
            disabled: string;
            normalRgba: string;
            hoverRgba: string;
            activeRgba: string;
            disabledRgba: string;
        };
        getMicroWidgets: () => any;
        getMicroWidgetByName: () => any;
        userInfo: {
            userid:string;
            account:string;
            name:string;
            mail:string;
            csflevel:string;
            usertype:number;
            telnumber:string;
            [k:string]:unknown;
        }
    };

    token: {
        /**
         * 获取当前登录用户的token（使用函数获取，才能保证获取到的是最新的）
         */
        getToken: () => {
            access_token: string;
            id_token:string;
            refresh_token:string;
        };

        /**
         * token过期的回调（token失效时，可调用此函数，然后管理控制台就会退出到登录页面）
         */
        onTokenExpired: () => void;

        /**
         * 刷新token函数（token失效时，可调用此函数刷新token，刷新成功后需重发token失效的请求）
         */
        refreshOauth2Token: () => Promise<string>;
    };

    history: {
        getBasePath: string;
        getBasePathByName:() => string;
        navigateToMicroWidget: () => void;
    };

    /**
     * 开放给插件的组件
     */
    components: Record<string, any>;

    /**
     * 使用此函数加载组件
     */
    mountComponent: (params: { component: any, props: Record<string,any>, element: HTMLElement }) => void;
    
    /**
     * 使用此函数卸载组件
     */
    unmountComponent: (element: HTMLElement) => void;

    /**
     * navigate
     */
    navigate: (path: string) => void;
}

export interface AppConfigContextType extends CommonType {
    /**
     * 挂载插件的元素
     */
    element: HTMLElement;

    /**
     * oem
     */
    oemColor: {
        // 主色浅色背景色
        colorPrimaryBg: string;
        // 主色浅色背景悬浮态
        colorPrimaryBgHover: string;
        // 主色描边色
        colorPrimaryBorder: string;
        // 主色描边色悬浮态
        colorPrimaryBorderHover: string;
        // 主色悬浮态
        colorPrimaryHover: string;
        // 主色
        colorPrimary: string;
        // 主色激活态
        colorPrimaryActive: string;
        // 主色文本悬浮态
        colorPrimaryTextHover: string;
        // 主色文本
        colorPrimaryText: string;
        // 主色文本激活态
        colorPrimaryTextActive: string;
    };
}