/**
 * 高级定时器
 * @param callback 要执行的函数
 * @param timeout 调用间隔
 * @return Function 返回停止定时器的函数
 * @desc 如果callback返回Promise，则在Promise的resolve或reject中继续间隔调用，否则按间隔时间固定调用
 * @example
 * ```typescript
 *  let count = 5;
 *
 *  // 每隔 1s 计数 -1，计数为 0 时自动停止
 *  const stop = timer(() => (count -= 1) === 0 && stop());
 * ```
 */
export function timer(callback: Function, timeout: number = 1000): () => void {
    let timer;
    let caller = function () {
        let ret = callback();

        if (ret && ret.then) {
            ret.then(() => {
                if (timer !== null) {
                    timer = setTimeout(caller, timeout)
                }
            }, () => {
                if (timer !== null) {
                    timer = setTimeout(caller, timeout)
                }
            })
        } else {
            if (timer !== null) {
                timer = setTimeout(caller, timeout)
            }
        }
    }

    caller();

    return function () {
        clearTimeout(timer);
        timer = null;
    }
}

/**
 * 同步延迟
 * @param time 延迟的时间
 */
export function syncDelay(time: number = 1000): Promise<void> {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve()
        }, time)
    })
}