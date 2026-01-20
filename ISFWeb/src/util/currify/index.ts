/**
 * 函数柯里化工具
 * @params fn 函数
 * @params ...fixed 柯里化参数
 * @return function
 */
export function currify(fn: Function, ...fixed: Array<any>): (...args: Array<any>) => any {
    return function (...args) {
        return fn.apply(this, fixed.concat(args));
    }
}