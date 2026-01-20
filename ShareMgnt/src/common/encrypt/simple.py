# coding: utf-8
"""
加密类简单调用库
"""
from hashlib import md5
from .pyDes import des, CBC, PAD_PKCS5, ECB, PAD_NORMAL
from base64 import (b64encode,
                    b64decode,
                    urlsafe_b64encode,
                    urlsafe_b64decode,
                    encodestring,
                    decodestring
                    )
from M2Crypto import RSA, BIO


def md5_hash(value):
    """
    返回指定值的md5 hash
    """
    if not value:
        return value

    return md5(value).hexdigest()


def des_encrypt(key, value, init_value='initialization_value'):
    """
    返回指定值des之后的值
    """

    if isinstance(key, str):
        key = key.encode('utf-8')
    if isinstance(value, str):
        value = value.encode('utf-8')
    if isinstance(init_value, str):
        init_value = init_value.encode('utf-8')

    key = des(key, CBC, init_value, padmode=PAD_PKCS5)
    # 加密后的数据经过base64编码
    return b64encode(key.encrypt(value))

def des_encrypt_with_padzero(key, value, init_value='initialization_value'):
    """
    返回指定值des之后的值
    """

    if isinstance(key, str):
        key = key.encode('utf-8')
    if isinstance(value, str):
        value = value.encode('utf-8')
    if isinstance(init_value, str):
        init_value = init_value.encode('utf-8')

    key = des(key, CBC, init_value, pad=chr(0), padmode=PAD_NORMAL)
    # 加密后的数据经过base64编码
    return b64encode(key.encrypt(value))


def des_decrypt(key, encrypt_value, init_value='initialization_value'):
    """
    根据key值解密加密的数据encrypt_value,
    """
    if not encrypt_value:
        return b''

    if isinstance(key, str):
        key = key.encode('utf-8')
    if isinstance(init_value, str):
        init_value = init_value.encode('utf-8')

    key = des(key, CBC, init_value, padmode=PAD_PKCS5)

    # 加密后的数据先解码在界面
    return key.decrypt(b64decode(encrypt_value), padmode=PAD_PKCS5)

def des_decrypt_with_padzero(key, encrypt_value, init_value='initialization_value'):
    """
    根据key值解密加密的数据encrypt_value,
    """
    if not encrypt_value:
        return b''

    if isinstance(key, str):
        key = key.encode('utf-8')
    if isinstance(init_value, str):
        init_value = init_value.encode('utf-8')

    key = des(key, CBC, init_value, pad=chr(0), padmode=PAD_NORMAL)

    # 加密后的数据先解码在界面
    return key.decrypt(b64decode(encrypt_value), padmode=PAD_NORMAL)


def der_length(length):
    """DER encoding of a length"""
    if length < 128:
        return chr(length)
    prefix = 0x80
    result = ''
    while length > 0:
        result = chr(length & 0xff) + result
        length >>= 8
        prefix += 1
    return chr(prefix) + result


def rsa_encrypt(public_key, data):
    """
    rsa加密(public key为pkcs#1格式),再用safe_base64编码
    """
    pk = public_key.split('\n')
    pk = '\0' + decodestring("".join(pk[1:-2]))
    pk = '\x30\x0d\x06\x09\x2a\x86\x48\x86\xf7\x0d\x01\x01\x01\x05\x00\x03' + der_length(len(pk)) + pk
    pk = '\x30' + der_length(len(pk)) + pk
    pk = '-----BEGIN PUBLIC KEY-----\n' + encodestring(pk) + '-----END PUBLIC KEY-----'
    bio = BIO.MemoryBuffer(pk)
    key = RSA.load_pub_key_bio(bio)
    en_data = key.public_encrypt(data, RSA.pkcs1_padding)
    result = urlsafe_b64encode(en_data)
    return result


# def eisoo_rsa_encrypt(data):
#     """
#     rsa加密， 先用私钥加密，pkcs1_padding填充，在进行base64编码
#     """
#     pub_key = b"""
# -----BEGIN PUBLIC KEY-----
# MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC7JL0DcaMUHumSdhxXTxqiABBC
# DERhRJIsAPB++zx1INgSEKPGbexDt1ojcNAc0fI+G/yTuQcgH1EW8posgUni0mcT
# E6CnjkVbv8ILgCuhy+4eu+2lApDwQPD9Tr6J8k21Ruu2sWV5Z1VRuQFqGm/c5vaT
# OQE5VFOIXPVTaa25mQIDAQAB
# -----END PUBLIC KEY-----
#     """
#     bio = BIO.MemoryBuffer(pub_key)
#     rsa = RSA.load_pub_key_bio(bio)

#     if isinstance(data, str):
#         data = data.encode('utf-8')
#     encrypted = rsa.public_encrypt(data, RSA.pkcs1_padding)
#     return b64encode(encrypted)


def eisoo_rsa_decrypt(data):
    """
    rsa解密， 先用safe_base64解密，再用私钥解密
    """
    private_key = b"""
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDB2fhLla9rMx+6LWTXajnK11Kdp520s1Q+TfPfIXI/7G9+L2YC
4RA3M5rgRi32s5+UFQ/CVqUFqMqVuzaZ4lw/uEdk1qHcP0g6LB3E9wkl2FclFR0M
+/HrWmxPoON+0y/tFQxxfNgsUodFzbdh0XY1rIVUIbPLvufUBbLKXHDPpwIDAQAB
AoGBALCM/H6ajXFs1nCR903aCVicUzoS9qckzI0SIhIOPCfMBp8+PAJTSJl9/ohU
YnhVj/kmVXwBvboxyJAmOcxdRPWL7iTk5nA1oiVXMer3Wby+tRg/ls91xQbJLVv3
oGSt7q0CXxJpRH2oYkVVlMMlZUwKz3ovHiLKAnhw+jEsdL2BAkEA9hA97yyeA2eq
f9dMu/ici99R3WJRRtk4NEI4WShtWPyziDg48d3SOzYmhEJjPuOo3g1ze01os70P
ApE7d0qcyQJBAMmt+FR8h5MwxPQPAzjh/fTuTttvUfBeMiUDrIycK1I/L96lH+fU
i4Nu+7TPOzExnPeGO5UJbZxrpIEUB7Zs8O8CQQCLzTCTGiNwxc5eMgH77kVrRudp
Q7nv6ex/7Hu9VDXEUFbkdyULbj9KuvppPJrMmWZROw04qgNp02mayM8jeLXZAkEA
o+PM/pMn9TPXiWE9xBbaMhUKXgXLd2KEq1GeAbHS/oY8l1hmYhV1vjwNLbSNrH9d
yEP73TQJL+jFiONHFTbYXwJAU03Xgum5mLIkX/02LpOrz2QCdfX1IMJk2iKi9osV
KqfbvHsF0+GvFGg18/FXStG9Kr4TjqLsygQJT76/MnMluw==
-----END RSA PRIVATE KEY-----
    """
    private_key = RSA.load_key_string(private_key)
    # 去除所有的\r和\n
    data = data.replace(r'\n', '')
    data = data.replace(r'\r', '')
    data = urlsafe_b64decode(data)
    return private_key.private_decrypt(data, RSA.pkcs1_padding)

def eisoo_rsa2048_decrypt(data):
    """
    rsa解密， 先用safe_base64解密，再用私钥解密
    """
    private_key = b"""
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAsyOstgbYuubBi2PUqeVjGKlkwVUY6w1Y8d4k116dI2SkZI8f
xcjHALv77kItO4jYLVplk9gO4HAtsisnNE2owlYIqdmyEPMwupaeFFFcg751oiTX
JiYbtX7ABzU5KQYPjRSEjMq6i5qu/mL67XTkhvKwrC83zme66qaKApmKupDODPb0
RRkutK/zHfd1zL7sciBQ6psnNadh8pE24w8O2XVy1v2bgSNkGHABgncR7seyIg81
JQ3c/Axxd6GsTztjLnlvGAlmT1TphE84mi99fUaGD2A1u1qdIuNc+XuisFeNcUW6
fct0+x97eS2eEGRr/7qxWmO/P20sFVzXc2bF1QIDAQABAoIBAACDungGYoJ87bLl
DUQUqtl0CRxODoWEUwxUz0XIGYrzu84nJBf5GOs9Xv6i9YbNgJN2xkJrtTU7VUJF
AfaSP4kZXqqAO9T1Id9zVc5oomuldSiLUwviwaMek1Yh9sFRqWNGGxBdd7Y1ckm8
Roy+kHZ7xXqlIxOmdCC+7DgQMVgSV64wzQY8p7L9kTLIkeDodEolkUkGsreF9I9S
kzlLjGU9flPt13319G0KSaQUWEpxF/UBr2gKJvQPQHSRzzl5HlRwznZkU4Hs6RID
ue6E68ZJNMRn3FUAvLMCRw9C4PQQR/x/50WH4BXJ9veVIOIpTVCJedI0QZjbVuBk
RPKHTMkCgYEA2XjGIw9Vp0qu/nCeo5Nk15xt/SJCn0jIhyRpckHtCidotkiZmFdU
vUK7IwbAUPqEJcgmS/zwREV8Gff8S324C2RoDN4FxFtBMZgQjqV1zYqGLQSbTJUh
GlpTe7jKVskuSPSf00OqqAIlYNtzZK3mWj8MadFD99Wo9gktXRAFdf0CgYEA0uBe
wuE007XLqb8ANS+4U0CkexeVDkDzI9yXN2CB+L5wmJ/WsNF8iD53xHxpwZWRiizX
ArBdhWL9yv4YkbryyD15NRSQhLanRcs0MqGh1GJJ9vpGzBjfJJ3Bw0hBfkwnf/C6
nNzGjNWNTeNKwlcFaVhBADyGYZt9Len9YYFNKrkCgYEAmsn7BYNprOxciCAy2i0U
Lt9Z7j3Pe757dK13HGtOQ9bvEie0o5ktaJSxzGmGw1y8aIQAtj9v6Lgob/dxrW3r
bLhn0xjItA1b5ufciRu+MLFzdWF9BFJ1QGOgXkSWSJVji2wKwn28X18/qaQpizS3
6+5KcJsRrLp4S78WedHogSUCgYEAomb5k8wtCv7vIoNefZeKtVMLWWEIAjozBmNU
cel5L0A7Js+yX+p1pde2FTRbniK6O1fdHs0EuT1Lh5G5CkKXx27QcfisdAjXOgEM
6hFguFgZ7oNBEt30vBZiqypyhfnQUc/rZ/L/VmcAtANgB9tM55x4Mt5p/7Hn7fxO
j1EtRMECgYEAp2sI035BcCR2kFW1vC9eXLAPZ0anyy1/T1dEgFJ/ELqmGEMEWZKA
9H1KH6YIkDdXabwfaSTRebaEescCxRtgmo5WEdZxw4Nz66SSomc24aD0iem7+VSl
x2qRWdif0jHG8fOdMey3NrY7NF4xQTzuO9jDnLpBTwFg3o7QlywIBlM=
-----END RSA PRIVATE KEY-----
    """
    private_key = RSA.load_key_string(private_key)
    # 去除所有的\r和\n
    data = data.replace(r'\n', '')
    data = data.replace(r'\r', '')
    data = urlsafe_b64decode(data)
    return private_key.private_decrypt(data, RSA.pkcs1_padding)


def des_encrypt_simple(data):
    # des 加密并用 base64 编码
    key = "01001110"
    # 用空格填充数据长度为8的倍数
    n = (8 - len(data) % 8) % 8
    for i in range(0, n):
        data += " "
    d = des(key, ECB, IV=None, pad=None, padmode=PAD_NORMAL)
    result = b64encode(d.encrypt(data))
    return result


def des_decrypt_simple(data):
    # base64 解密然后 des 解密
    key = "01001110"
    d = des(key, ECB, IV=None, pad=None, padmode=PAD_NORMAL)
    result = d.decrypt(b64decode(data))
    # 去除填充空格
    result = result.rstrip()
    return result


# def cmp_encrypt(pwd):
#     # cmp 租户密码加密
#     key = """
# -----BEGIN PUBLIC KEY-----
# MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD18CpBqwJS93lBKglAaFsyBjC0
# AgO2hrNJ4Yv3lliZ7fQoRJ9p1a+P/VcJFTIPTeRFV0TU02O5n28PkAtTXdIKbJMa
# v85PuLB7b5aIXKXOQS4K3veKR8PT/2Z+1iUM/dG8il6f3amXNk0mZA8MWPhSXKia
# 9G7pi8U/kLdN9pN9YwIDAQAB
# -----END PUBLIC KEY-----
#     """
#     pubkey = key.encode('utf8')
#     bio = BIO.MemoryBuffer(pubkey)
#     rsa = RSA.load_pub_key_bio(bio)
#     encrypted = rsa.public_encrypt(pwd.encode('utf8'), RSA.pkcs1_padding)
#     return b64encode(encrypted)
