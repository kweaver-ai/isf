import { merge, reduce, keys, isArray, isObject, isFunction } from 'lodash';

export const hasOwnProperty = Object.prototype.hasOwnProperty

/**
 * 对象访问
 * @param obj 要访问的对象
 * @param ...args 传递单个参数key返回obj[key]的值
 * @param ...args 传递[key, value] 设置obj并返回obj
 * @param ...args 传递对象设置obj并返回obj
 * @return any 传递单个参数返回对应的值，传递两个参数或对象返回设置的对象
 */
export function access(obj: Record<string, any>, ...args: Array<any>): any {
    if (args.length === 1) {
        const arg = args[0];

        if (typeof arg === 'string') {
            return obj[arg];
        } else if (typeof arg === 'object') {
            return merge(obj, arg);
        }
    } else if (args.length === 2) {
        const [key, value] = args;
        obj[key] = value;
        return obj;
    }
}

/**
 * 计算表达式
 * @param expr 表达式，如果是函数则返回计算结果，否则返回原值
 * @return 返回表达式计算结果
 */
export function evaluate(expr) {
    if (expr instanceof Function) {
        return expr();
    } else {
        return expr;
    }
}

/**
 * 对两个对象进行浅比较
 * 参考实现 https://github.com/facebook/fbjs/blob/c69904a511b900266935168223063dd8772dfc40/packages/fbjs/src/core/shallowEqual.js
 */
export function shallowEqual(objA: any, objB: any): boolean {
    if (objA === objB) {
        return true;
    } else {
        if (
            typeof objA !== 'object' || objA === null ||
            typeof objB !== 'object' || objB === null
        ) {
            return false;
        }

        const keysA = Object.keys(objA);
        const keysB = Object.keys(objB);

        if (keysA.length !== keysB.length) {
            return false;
        }

        // Test for A's keys different from B.
        for (let i = 0; i < keysA.length; i++) {
            if (
                !hasOwnProperty.call(objB, keysA[i]) ||
                !(objA[keysA[i]] === objB[keysA[i]])
            ) {
                return false;
            }
        }

        return true;
    }
}

/**
 * 判断输入是否返回假值
 * @param input 测试值
 */
export function isExist(input: any): boolean {
    return input !== null && input !== undefined;
}

/**
 * 判断值为不存在／空数组／空对象
 */
export function isEmpty(input: any): boolean {
    if (!isExist(input)) {
        return true;
    } else if (isArray(input)) {
        return !input.length;
    } else if (isObject(input)) {
        return !keys(input).length;
    } else {
        return false;
    }
}

/**
 * 将任何输入Promise化
 * @param input 输入值
 */
export function promisify(input: any): Promise<any> {
    return isFunction(input && input.then) ?
        input :
        Promise.resolve(input)
}

interface Chain {
    (callback: ChainCallback): ChainInstance;
}

interface ChainCallback {
    (o: any, index: number, arr: Array<any>): any;
}

interface ChainInstance {
    (data: Array<any>): Promise<any>;
}

/**
 * Promise 队列调用
 * @param callback 回调
 * @return 返回调用队列函数，传入数据开始执行
 */
export const chain: Chain = function chain(callback) {
    return function (data) {
        return reduce(data, (prev, o, i, arr) => {
            return prev.then(() => callback(o, i, arr))
        }, Promise.resolve())
    }
}

/**
 * 按位加
 * @param x 被增加的比特位
 * @param y 增加的比特位
 */
export function bitSum(x?: number, ...adders: Array<number>): number | undefined {
    if (x === undefined) {
        return;
    }

    return reduce(adders, (result, adder) => result | adder, x);
}

/**
 * 按位减
 * @param x 被减的比特位
 * @param subs 被减的比特位
 */
export function bitSub(x?: number, ...subs: Array<number>): number | undefined {
    if (x === undefined) {
        return;
    }

    return reduce(subs, (result, sub) => result & ~sub, x);
}

/**
 * 测试比特位x是否包含比特位y
 * @param x 比特位值x
 * @param y 比特位值y
 */
export function bitTest(x: number, y: number): boolean {
    return (x & y) === y
}

/**
 * 测试输入值是否是undefined或者null
 * @param x 测试值
 */
export function isNil(x: any): boolean {
    return x === undefined || x === null;
}

/**
 * 深度改变object值并返回新object
 * @param object
 * @param path
 * @param value
 */
export function setCopy(object, path, value) {
    if (typeof path === 'string') {
        path = path.split('.')
    }
    const [key, ...subPath] = path
    if (path.length === 1) {
        return { ...object, [key]: value }
    }
    return { ...object, [key]: setCopy(object[key] || {}, subPath, value) }
}

/**
 * 安全的eval操作
 * @description 将字符串类型转换为对应的JS数据类型，如，safeEval('true') => true
 * @param input 输入值
 * @returns 返回JS类型的数据
 * @example
 *  ```ts
 *   safeEval('true') // true
 *  ```
 */
export function safeEvalPrimitive(input) {
    if (typeof input !== 'string') {
        return input
    }

    const lonelyValues = {
        'true': true,
        'false': false,
        NaN: NaN,
        'null': null,
        undefined: undefined,
    }

    if (lonelyValues.hasOwnProperty(input)) {
        return lonelyValues[input]
    }

    if (/^[+-]?\d+(\.\d+)?$/.test(input)) {
        return Number(input)
    }

    return input

}

/**
 * 删除数组中的某一项，并返回删除后的新数组
 * @param arr 操作的数组
 * @param index 指定要删除的下标值
 * @returns 删除后的新数组
 */
export function deleteArrElementByIndex<T>(arr: Array<T>, index: number): Array<T> {
    if (index < 0 || index >= arr.length) {
        return arr
    } else {
        return [...arr.slice(0, index), ...arr.slice(index + 1)]
    }
}

/**
 * 使用指定的回调方法修改数组的某一项并返回新数组
 * @param arr 操作的数组
 * @param index 指定修改项的下标值
 * @param callback 使用指定回调方法修改数组的某一项
 * @returns 修改后的新数组
 */
export function modifyArrElementByIndex(arr: Array<any>, index: number, callback: () => void): Array<any> {
    return reduce(arr, (result, element, i, arr) => {
        return i === index && isFunction(callback) ?
            [...result, callback(element, i, arr)]
            : [...result, element]
    }, [])
}