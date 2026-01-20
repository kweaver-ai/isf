import { insertStyle } from '@/util/browser'
import { ClassName } from '@/ui/helper'

/**
 * 根据OEM配置设定对应皮肤样式
 */
export function apply(theme) {
    return insertStyle({
        [`.${ClassName.BorderColor}`]: {
            borderColor: `${theme}!important;`,
        },
        [`.${ClassName.BorderColor__Focus}:focus`]: {
            borderColor: `${theme}!important;`,
        },
        [`.${ClassName.BorderTopColor}`]: {
            borderTopColor: `${theme}!important;`,
        },
        [`.${ClassName.BorderRightColor}`]: {
            borderRightColor: `${theme}!important;`,
        },
        [`.${ClassName.BorderBottomColor}`]: {
            borderBottomColor: `${theme}!important;`,
        },
        [`.${ClassName.BorderLeftColor}`]: {
            borderLeftColor: `${theme}!important;`,
        },
        [`.${ClassName.BackgroundColor}`]: {
            backgroundColor: `${theme}!important;`,
        },
        [`.${ClassName.Color}`]: {
            color: `${theme}!important;`,
        },
        [`.${ClassName.Color__Hover}:hover`]: {
            color: `${theme}!important;`,
        },
    })
}