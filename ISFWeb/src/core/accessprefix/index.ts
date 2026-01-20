import Cookies from 'js-cookie'

/**
 * 获取访问前缀
 */
const getAccessPrefix = (): string => {
    return Cookies.get('X-Forwarded-Prefix') || ''
}

export {
    getAccessPrefix,
}