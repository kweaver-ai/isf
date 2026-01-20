import { useEffect, useRef } from 'react'
import { isEqual } from 'lodash'

/**
 * 实现深比较函数
 */
function deepCompareEquals(a: any, b: any): boolean {
    return isEqual(a, b)
}

function useDeepCompareMemoize(value: any): any {
    const ref = useRef<any>()

    if (!deepCompareEquals(value, ref.current)) {
        ref.current = value
    }

    return ref.current
}

export function useDeepCompareEffect(callback: () => void, dependencies: any[]): void {
    useEffect(
        callback,
        dependencies.map(useDeepCompareMemoize),
    )
}