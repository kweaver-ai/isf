import './public-path'
import React from 'react'
import { createRoot } from 'react-dom/client'
import { ConfigProvider as AntdConfigProvider, message } from 'antd'
import zhCN from 'antd/es/locale/zh_CN'
import enUS from 'antd/es/locale/en_US'
import zhTW from 'antd/es/locale/zh_TW'
import AppConfigContext from '../context/AppConfigContext'
import { initData } from './init'
import { ContextType } from './types'
import { Lang } from '../context/type'
import session from '@/util/session';
import { generate } from "@ant-design/colors";
import './global.css'

const getUILocale = (lang: Lang): typeof enUS | typeof zhTW | typeof zhCN => {
    const langs = {
        [Lang.us]: enUS,
        [Lang.tw]: {
            ...zhTW, 
            DatePicker: { 
                ...zhTW.DatePicker, 
                lang: { 
                    ...zhTW.DatePicker!.lang, 
                    monthFormat: 'M月', 
                    shortWeekDays: ['日', '一', '二', '三', '四', '五', '六']
                }
            } 
        },
        [Lang.zh]: {
            ...zhCN, 
            DatePicker: { 
                ...zhCN.DatePicker, 
                lang: { 
                    ...zhCN.DatePicker!.lang, 
                    monthFormat: 'M月', 
                    shortWeekDays: ['日', '一', '二', '三', '四', '五', '六']
                }
            } 
        },
    }

    return langs[lang] || zhCN
}

let removeThemeStyleFn
let root;

const bootstrap = async () => {
}

const formatProps = (props: any) => {
    session.set('isf.userInfo', props?.config.userInfo)
    session.set('isf.userid', props?.userid)
    return { 
        ...props, 
        host: props.config.systemInfo.location.hostname, 
        port: props.config.systemInfo.location.port, 
        protocol: props.config.systemInfo.location.protocol, 
        theme: props.oemConfigs.theme, 
        getUserInfo: () => props.config.userInfo, 
        getToken: () => props.token.getToken.access_token, 
        onTokenExpired: props.token.onTokenExpired, 
        refreshToken: props.token.refreshOauth2Token
    }
}

const getMount = (Component: React.ReactNode) => {
    /**
     * 挂载插件
     */
    const mount = async ({ container, ...otherProps }: ContextType) => {
        const props = formatProps(otherProps)

        const element = container?.querySelector('#isf-web-plugins') || document.getElementById('isf-web-plugins')
        const colorPalette = generate(props.theme)
       
        const oemColor = {
            // 主色浅色背景色
            colorPrimaryBg: colorPalette[0],
            // 主色浅色背景悬浮态
            colorPrimaryBgHover: colorPalette[1],
            // 主色描边色
            colorPrimaryBorder: colorPalette[2],
            // 主色描边色悬浮态
            colorPrimaryBorderHover: colorPalette[3],
            // 主色悬浮态
            colorPrimaryHover: colorPalette[4],
            // 主色
            colorPrimary: colorPalette[5], 
            // 主色激活态
            colorPrimaryActive: colorPalette[6],
            // 主色文本悬浮态
            colorPrimaryTextHover: colorPalette[4],
            // 主色文本
            colorPrimaryText: colorPalette[5], 
            // 主色文本激活态
            colorPrimaryTextActive: colorPalette[6],
        }
        const { removeThemeStyle } = initData({...props, container: element})
        removeThemeStyleFn = removeThemeStyle
        root = createRoot(element);

        message.config({
            getContainer() {
                return element as HTMLElement
            },
        })
        
        root.render(
            <AntdConfigProvider
                prefixCls='isf'
                theme={{
                    token: {
                        colorPrimary: props.theme,
                        colorInfo: props.theme,
                        colorLink: props.theme,
                        colorSuccess: '#52c41a',
                        colorWarning: '#faad14',
                        colorError: '#f5222d',
                    },
                }}
                locale={getUILocale(props.lang)}
                getPopupContainer={() => element}
            >
                <AppConfigContext.Provider value={{ element, ...props, oemColor}}>
                    <style type="text/css">
                        {`
                            html a {
                                color: ${props.theme};
                                
                                &:hover {
                                    color: ${oemColor.colorPrimaryHover};
                                }

                                &:active {
                                   color: ${oemColor.colorPrimaryTextActive}; 
                                }
                            }
                            
                            a[disabled] {
                                color: rgba(0,0,0,0.25) !important; 
                            }
                        `}
                    </style>
                    {Component}
                </AppConfigContext.Provider>
            </AntdConfigProvider>
        )
    }

    return mount
}

/**
 * 卸载插件
 */
const unmount = async () => {
    try {
        root.unmount();
        removeThemeStyleFn?.()
    } catch {}
}

export {
    bootstrap,
    getMount,
    unmount,
}