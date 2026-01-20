package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func Rsa2048Decrypt(encrypted string) (decrypted string, err error) {
	privateKey := `
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

	// pem 解码
	block, _ := pem.Decode([]byte(privateKey))
	// X509解码
	var priv *rsa.PrivateKey

	priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	// base64 解密 go base64包可以处理\r\n 无需额外处理
	var tempData []byte

	tempData, err = base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return
	}

	// rsa解码
	temp, err := rsa.DecryptPKCS1v15(rand.Reader, priv, tempData)
	if err != nil {
		return
	}

	return string(temp), nil
}
