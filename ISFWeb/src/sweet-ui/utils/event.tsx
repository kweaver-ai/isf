import { isFunction } from 'lodash'

export class SweetUIEvent<P> {
    constructor(detail: P) {
        this.detail = detail;
    }

    /**
     * 事件额外参数
     */
    public detail: P;

    /**
     * 标识事件默认行为是否被阻止
     */
    public defaultPrevented: boolean = false;

    /**
     * 是否冒泡事件到上层组件
     */
    public cancelBubble: boolean = false;

    /**
     * 阻止事件默认行为
     */
    preventDefault(): void {
        this.defaultPrevented = true;
    }

    /**
     * 阻止事件冒泡
     */
    stopPropagation(): void {
        this.cancelBubble = true;
    }
}

/**
 * 创建事件触发函数，允许使用同一个事件对象触发多个事件，返回一个函数
 * 利用该返回函数执行执行，如果在事件中调用了了 event.preventDefault()，则不会执行.then()定义的默认行为
 * @param detail 触发事件时传递的参数
 *
 * ```typescript
 * const dispatchEvent = createEventDispatcher(onEventFired, (event) => console.log(event.detail))
 *
 * dispatchEvent(true) // 打印 true
 *
 * // 如果执行了 event.preventDefault()，则不会打印true
 * const onEventFired = event => event.preventDefault()
 *
 * dispatchEvent(true) // 不会打印
 * ```
 */
export function createEventDispatcher(eventListener?: (event: SweetUIEvent<any>) => void, defaultHandler?: (event: SweetUIEvent<any>) => void) {
    return function dispatchEvent(detailOrEvent?: SweetUIEvent<any> | any) {
        // 如果事件由上层组件传递过来并且取消了冒泡，则在此中断
        if (detailOrEvent instanceof SweetUIEvent && detailOrEvent.cancelBubble) {
            return
        }

        let event

        if (detailOrEvent instanceof SweetUIEvent) {
            event = new SweetUIEvent(detailOrEvent.detail)
        } else {
            event = new SweetUIEvent(detailOrEvent)
        }

        if (isFunction(eventListener)) {
            eventListener(event);
        }

        if (!event.defaultPrevented) {
            if (isFunction(defaultHandler)) {
                defaultHandler(event)
            }
        }
    }
}