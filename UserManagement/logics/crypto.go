// Package logics Anyshare 业务逻辑层 -密码相关
package logics

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/daknob/ntlm"
)

// DesCBCPadType des cbc 加密填充类型
type DesCBCPadType int

const (
	// 使用iota，增加新的枚举时，必须按照顺序添加，不得在中间插入
	_ DesCBCPadType = iota

	// PadNormal 填充
	PadNormal

	// PKCS5Padding 填充
	PKCS5Padding
)

// DecodeType 解密类型
type DecodeType int

const (
	// 使用iota，增加新的枚举时，必须按照顺序添加，不得在中间插入
	_ DecodeType = iota

	// RSA1024 解密
	RSA1024

	// RSA2048 解密
	RSA2048
)

var (
	strRSA1024PrivateKey = `
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
	`

	strRSA2048PrivateKey = `
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
	`
)

// decodeRSA rsa 解密
func decodeRSA(data string, typ DecodeType) (out string, err error) {
	// 检查rsa解密类型
	var strPrivateKey string
	switch typ {
	case RSA1024:
		strPrivateKey = strRSA1024PrivateKey
	case RSA2048:
		strPrivateKey = strRSA2048PrivateKey
	default:
		return out, errors.New("decode type error")
	}

	// pem 解码
	block, _ := pem.Decode([]byte(strPrivateKey))
	// X509解码
	var privateKey *rsa.PrivateKey
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	// base64 解密 go base64包可以处理\r\n 无需额外处理
	var tempData []byte
	tempData, err = base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	// rsa解码
	temp, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, tempData)
	if err != nil {
		return
	}
	out = string(temp)
	return
}

// encodeMD5 对密码明文进行MD5加密，空字符串不加密
func encodeMD5(pwd string) (out string, err error) {
	h := md5.New()
	_, err = h.Write([]byte(pwd))
	if err != nil {
		return
	}
	out = hex.EncodeToString(h.Sum(nil))
	return
}

// encodeNtlm ntlm 加密
func encodeNtlm(pwd string) (out string) {
	return ntlm.FromASCIIStringToHex(pwd)
}

// pkcs5Padding PKCS5 填充
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// padNormal padNormal 填充
func padNormal(ciphertext []byte, blockSize int) []byte {
	ciphertext = append(ciphertext, make([]byte, blockSize-len(ciphertext)%blockSize)...)
	return ciphertext
}

// encodeDes des 加密
func encodeDes(src string, typ DesCBCPadType) (string, error) {
	// string编码
	strDeskey := "Ea8ek&ah"
	key := []byte(strDeskey)
	data := []byte(src)

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	nBlockSize := block.BlockSize()
	switch typ {
	case PKCS5Padding:
		data = pkcs5Padding(data, nBlockSize)
	case PadNormal:
		data = padNormal(data, nBlockSize)
	default:
		return "", fmt.Errorf("des cbc pad type error:%d", typ)
	}
	// 获取CBC加密模式
	mode := cipher.NewCBCEncrypter(block, key)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)

	return base64.StdEncoding.EncodeToString(out), nil
}

// encodeSha2 sha256 加密
func encodeSha2(src string) string {
	// 字符串转化字节数组
	message := []byte(src)
	// 创建一个基于SHA256算法的hash.Hash接口的对象
	// sha-256加密
	hash := sha256.New()
	// 输入数据
	hash.Write(message)
	// 计算哈希值
	hashbytes := hash.Sum(nil)
	// 将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(hashbytes)
	// 返回哈希值
	return hashCode
}

// MD5Encrypt md5加密
func MD5Encrypt(plainText string) (encryptedPassword string, err error) {
	h := md5.New()
	_, err = h.Write([]byte(plainText))
	if err != nil {
		return
	}

	md5Passwd := hex.EncodeToString(h.Sum(nil))

	return md5Passwd, nil
}
