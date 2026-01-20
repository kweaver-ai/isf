import { getPloicyInfo } from '../apis/console/loginsecuritypolicy'

export enum ProtLevel {
    /**
     * 通用级
     */
    Common = 0,

    /**
     * 秘密级
     */
    Classified = 1,

    /**
     * 机密级一级
     */
    Confidential = 2,

    /**
     * 机密级增强
     */
    MoreConfidential = 3
}

/**
 * 系统保护等级-对应密码有效期
 */
export const ProtLevelPassExpire = {
    /**
     * 普通级
     */
    [ProtLevel.Common]: -1,

    /**
     * 秘密级 - 小于30天
     */
    [ProtLevel.Classified]: 30,

    /**
     * 机密级一般 - 小于7天
     */

    [ProtLevel.Confidential]: 7,

    /**
     * 机密级增强 - 小于3天
     */
    [ProtLevel.MoreConfidential]: 3,
}

/**
 * 获取系统保护等级
 */
export const getSystemProtectionLevel = async (): Promise<ProtLevel> => {
    try {
        const { data: [{ value: { level } }] } = await getPloicyInfo({ mode: 'current', name: 'system_protection_levels' })
        return level
    } catch {
        return ProtLevel.Common
    }
}
