#!/usr/bin/python3
#    coding:utf-8
import json
from Crypto.Cipher import AES
from binascii import b2a_hex, a2b_hex

class AESUtil(object):
    key_len = 32
    @classmethod
    def get_32_bit_key(cls, key):
        if isinstance(key, str):
            key = key.encode('utf-8')
        if len(key) >= cls.key_len:
            key = key[:cls.key_len]
        else:
            key = "%s%s" % (key, '0' * (cls.key_len - len(key)))
        return key

    @classmethod
    def zero_padding(cls, data):
        if data is None:
            data = ""
        if isinstance(data, str):
            data = data.encode('utf-8')
        if isinstance(data, (dict, list, set, tuple)):
            data = json.dumps(data)
        count = 16 - len(data) % 16
        if count:
            data += '\0' * count
        return data

    @classmethod
    def encrypt_data(cls, key, data):
        return b2a_hex(AES.new(cls.get_32_bit_key(key), AES.MODE_ECB).encrypt(cls.zero_padding(data)))

    @classmethod
    def descrypt_data(cls, key, hex_str):
        return AES.new(cls.get_32_bit_key(key), AES.MODE_ECB).decrypt(a2b_hex(hex_str)).rstrip('\0')
