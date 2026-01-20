import { debounce, throttle } from 'lodash'
import localStorage from '@/util/local'

const LastTimeKey = 'console.active.lastTime'
const waitTime = 2 * 1000

const updateLocalStorage = () => {
    localStorage.set(LastTimeKey, new Date().getTime())
}

/**
 * 接口调用时更新操作时间，主要用于防止批量调用接口时超时退出
 */
export const apiUpdateActivityStatus = throttle(() => {
    updateLocalStorage()
}, waitTime)

/**
 * 用户操作更新操作时间
 */
export const updateActivityStatus = debounce(() => {
    updateLocalStorage()
}, waitTime)
