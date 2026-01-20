package logics

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"

	"github.com/kweaver-ai/go-lib/rest"
)

// DecodeAndDecrypt BASE64解码 && RSA解密
func DecodeAndDecrypt(cipherText string, privateKey *rsa.PrivateKey) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", rest.NewHTTPErrorV2(rest.BadRequest, err.Error())
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedData)
	if err != nil {
		return "", rest.NewHTTPErrorV2(rest.BadRequest, err.Error())
	}
	return string(plainText), nil
}
