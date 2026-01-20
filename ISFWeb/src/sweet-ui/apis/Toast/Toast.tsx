import Toaster, { ToastInfo } from '../../components/Toaster/Toaster'

/**
 * 默认持续时间(ms)
 */
let defaultDuration = 2000

/**
 * toaster实例
 */
let toasterInstance = null

/**
 * toastChunk唯一标识
*/
let key = 1

/**
 * 一个界面最多显示toastChunk个数
 */
let maxCount = 10

/**
 * 获取toaster实例
 */
function getToasterInstance(callback, element) {
    // 如果已存在实例则无需重新创建
    if (toasterInstance) {
        callback(toasterInstance)
        return
    }

    Toaster.newInstance(
        {
            maxCount,
            holder: element,
        },
        (instance) => {
            if (toasterInstance) {
                callback(toasterInstance)
                return
            }

            toasterInstance = instance
            callback(instance)
        },
    )
}

/**
 * 添加ToastChunk
 */
function toast(args: ToastInfo): void {
    const contextElement = document.querySelector('#isf-web-plugins')
    const duration = args.duration || defaultDuration

    getToasterInstance((instance) => {
        instance.addToast({
            ...args,
            duration,
            key: args.key || key++,
            content: args.content,
            onClose: typeof args.onClose === 'function' ? args.onClose : () => Promise.resolve(true),
        })
    }, args.element || contextElement || window.document.body)
}

/**
 * 导出Toast.open方法，深色
 */
export const open = (content: JSX.Element | string, options = {}): void => {
    toast({ content, ...options })
}

/**
 * toast成功，深色
 */
export const success = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'success', ...options })
}

/**
 * toast失败，深色
 */
export const error = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'error', ...options })
}

/**
 * toast警告，深色
 */
export const warning = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'warning', ...options })
}

/**
 * toast一般消息提示，深色
 */
export const info = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'info', ...options })
}

/**
 * toast无图标，浅色
 */
export const lightOpen = (content: JSX.Element | string, options = {}): void => {
    toast({ content, mode: 'light', ...options })
}

/**
 * toast成功，浅色
 */
export const lightSuccess = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'success', mode: 'light', ...options })
}

/**
 * toast失败，浅色
 */
export const lightError = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'error', mode: 'light', ...options })
}

/**
 * toast警告，浅色
 */
export const lightWarning = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'warning', mode: 'light', ...options })
}

/**
 * toast一般消息提示，浅色
 */
export const lightInfo = (content: JSX.Element | string, options = {}): void => {
    toast({ content, type: 'info', mode: 'light', ...options })
}

/**
 * 导出Toast.destroy方法
 */
export const destory = (): void => {
    if (toasterInstance) {
        toasterInstance.destroy()

        toasterInstance = null
    }
}