import { findDOMNode } from 'react-dom';

/**
 * 皮肤相关类名
 */
export const ClassName = {
    BorderColor: 'skin-border-color',

    BorderColor__Focus: 'skin-border-color--focus',

    BorderTopColor: 'skin-border-top-color',

    BorderRightColor: 'skin-border-right-color',

    BorderBottomColor: 'skin-border-bottom-color',

    BorderLeftColor: 'skin-border-left-color',

    BackgroundColor: 'skin-background-color',

    BackgroundColor__Hover: 'skin-background-color--hover',

    BackgroundColor__Active: 'skin-background-color--active',

    Color: 'skin-color',

    Color__Hover: 'skin-color--hover',

    OemBackground: 'oem-background',

    OemWelcome: 'oem-welcome',

    OemProduct: 'oem-product',

    OemLogo: 'oem-logo',
}

/**
 * 获取局中坐标
 */
export function getCenterCoordinate(component: HTMLElement) {
    const container = findDOMNode(component);

    if (!container) {
        return {}
    } else {
        return {
            top: (document.documentElement.clientHeight - container.clientHeight) / 2,
            left: (document.documentElement.clientWidth - container.clientWidth) / 2,
        }
    }
}

/**
 * 根据字体大小获取文字样式
 */
export function getTextStyle({ fontSize }: { fontSize: number }): { fontSize: number; lineHeight: string } {
    return {
        [13]: { fontSize: '13px', lineHeight: '21px' },

    }[fontSize]
}