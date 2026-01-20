// Package rsa RSA加解密
package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// 公钥和私钥为公司统一RSA公私钥
var privateKeystr = `
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC7JL0DcaMUHumSdhxXTxqiABBCDERhRJIsAPB++zx1INgSEKPG
bexDt1ojcNAc0fI+G/yTuQcgH1EW8posgUni0mcTE6CnjkVbv8ILgCuhy+4eu+2l
ApDwQPD9Tr6J8k21Ruu2sWV5Z1VRuQFqGm/c5vaTOQE5VFOIXPVTaa25mQIDAQAB
AoGAe3ZNTExX7hpGtd097Uu+okmwcCJvqkv2sxkbkGpnBE7avXBE29ABItt/mAoB
AkJvshH8m+hhjwuaD62VkO7qsppTg4yL98Z0ZZ4kPqxJaIVU8FDmJyz1Png4ywg9
mw57saoZ+7GFQSITA7Kb5BeMP2xNeLIWjN2s29fhWMxTskECQQD4C4hcT1nWmA8i
j27eK/XDbVeceb8x1fKZ0wor+fQ5wTwmnhVPrRe4AgBVin8kR9kw+fu55Kk0U02o
uZTMh6U1AkEAwSUzIeX3R8TpaRsKnjV0GWGNdxesAg/IZw5k0A9JyiEvFFhNHfp9
00zvutxdK0D1tdzrn+UHwFrVzlN4zqtDVQJBAJT8OFdZwhhHFTAo/uqrdN6BGpJ9
/f0tCJ6kSAPKCot2KW74nMxSp2B6s0CuA1gDX80vGae6VHd9YbPqZBnFj9ECQBQO
HcoWS9/q5WWhhi+5Uy3TgFHuZlDsfJ2e0/76p2nSmkXdiVxkhy4qnfXkLdRw8VKJ
9vlqWayygeLjrfaft+UCQQDH7HQNcfZOcIYV0kpdztP/0bAeh7VHpjDrga7R+W+K
PR2M7HVDPSmmLtrm5snlUwVfu24VF4SO113z24zy6IRn
-----END RSA PRIVATE KEY-----
`

var publicKeystr = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC7JL0DcaMUHumSdhxXTxqiABBC
DERhRJIsAPB++zx1INgSEKPGbexDt1ojcNAc0fI+G/yTuQcgH1EW8posgUni0mcT
E6CnjkVbv8ILgCuhy+4eu+2lApDwQPD9Tr6J8k21Ruu2sWV5Z1VRuQFqGm/c5vaT
OQE5VFOIXPVTaa25mQIDAQAB
-----END PUBLIC KEY-----
`

// RsaInter RSA加解密接口
type RsaInter interface {
	Encrypt(data []byte) (string, error)
	Decrypt(secretData []byte) (string, error)
}

type rsaInter struct {
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

// NewRsa 返回rsa结构体
func NewRsa() RsaInter {
	rsaObj := &rsaInter{}
	rsaObj.init()
	return rsaObj
}

// initial 初始化
func (r *rsaInter) init() {
	if privateKeystr != "" {
		block, _ := pem.Decode([]byte(privateKeystr))
		r.rsaPrivateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	if publicKeystr != "" {
		block, _ := pem.Decode([]byte(publicKeystr))
		publickKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
		switch p := publickKey.(type) {
		case *rsa.PublicKey:
			r.rsaPublicKey = p
		default:
			r.rsaPublicKey = nil
		}
	}
}

// Encrypt rsa加密
func (r *rsaInter) Encrypt(data []byte) (string, error) {
	blockLength := r.rsaPublicKey.N.BitLen()/8 - 11
	if len(data) <= blockLength {
		encrypt, err := rsa.EncryptPKCS1v15(rand.Reader, r.rsaPublicKey, data)
		return string(encrypt), err
	}

	buffer := bytes.NewBufferString("")

	pages := len(data) / blockLength

	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(data) {
				continue
			}
			end = len(data)
		}

		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, r.rsaPublicKey, data[start:end])
		if err != nil {
			return "", err
		}
		buffer.Write(chunk)
	}
	return buffer.String(), nil
}

// Decrypt rsa解密
func (r *rsaInter) Decrypt(secretData []byte) (string, error) {
	blockLength := r.rsaPublicKey.N.BitLen() / 8
	if len(secretData) <= blockLength {
		decrypt, err := rsa.DecryptPKCS1v15(rand.Reader, r.rsaPrivateKey, secretData)
		return string(decrypt), err
	}

	buffer := bytes.NewBufferString("")

	pages := len(secretData) / blockLength
	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(secretData) {
				continue
			}
			end = len(secretData)
		}

		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, r.rsaPrivateKey, secretData[start:end])
		if err != nil {
			return "", err
		}
		buffer.Write(chunk)
	}
	return buffer.String(), nil
}
