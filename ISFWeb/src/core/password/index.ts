import { getRandom as secureRandom } from '@/util/random'

/**
 * 产生随机密码函数
 */
export function generateRandomPwd(passwordLength: number): string {
    let pwdType,
        rdmStr = '';

    while (rdmStr.length < passwordLength - 3) {
        pwdType = Math.floor(secureRandom(0, 1, 2) * 10) % 3;
        if (pwdType === 0) {
            rdmStr += (Math.floor(secureRandom(0, 1, 2) * 10))
        } else if (pwdType === 1) {
            rdmStr += (String.fromCharCode((Math.floor(secureRandom(0, 1, 3) * 100) + 4) % 26 + 65))
        } else {
            rdmStr += (String.fromCharCode((Math.floor(secureRandom(0, 1, 3) * 100) + 4) % 26 + 97))
        }
    }
    rdmStr += (Math.floor(secureRandom(0, 1, 2) * 10))
    rdmStr += (String.fromCharCode((Math.floor(secureRandom(0, 1, 3) * 100) + 4) % 26 + 65));
    rdmStr += (String.fromCharCode((Math.floor(secureRandom(0, 1, 3) * 100) + 4) % 26 + 97));

    return rdmStr.split('').sort(function () {
        if (secureRandom(0, 1, 1) > 0.5) {
            return 1
        } else {
            return -1
        }
    }).join('')
}

export const SpecialChars = '~!%#$@-_.'

/**
 * 生成带有特殊字符的随机密码
 */
export function generatePwdWithSpecialChar(passwordLength: number): string {
    const pwdStr = generateRandomPwd(passwordLength - 1) + SpecialChars[Math.floor(secureRandom(0, 1, 2) * SpecialChars.length)]
    return pwdStr.split('').sort(() => secureRandom(0, 1, 1) > 0.5 ? 1 : -1).join('')
}