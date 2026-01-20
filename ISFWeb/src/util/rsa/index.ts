import { RSAKey, hex2b64 } from '@/libs/rsa.min';
import * as JSEncrypt from '@/libs/jsencrypt.min';

export function rsaEncrypt(input, publicKey) {
    let rsaKey = new RSAKey();

    rsaKey.setPublic(publicKey, '10001');

    return hex2b64(rsaKey.encrypt(input)).replace(/(.{64})/g, function (match) {
        return match += '\n';
    });
}

/**
 * ����
 */
export function jsEncrypt(input, publicKey) {
    let encrypt = new JSEncrypt();

    encrypt.setPublicKey(publicKey);

    return encrypt.encrypt(input)
}

/**
 * ����
 */
export function jsDecrypt(input, privatekey) {
    let decrypt = new JSEncrypt();

    decrypt.setPrivateKey(privatekey);

    return decrypt.decrypt(input)
}