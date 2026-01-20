package rsa

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestRSANoLimitEncrypt(t *testing.T) {
	// 加解密正常字符串
	Convey("RSAEncrypt()", t, func() {
		data := "asAlqlTkWU0zqfxrLTed"
		rsaObj := NewRsa()
		encrypt, _ := rsaObj.Encrypt([]byte(data))
		decrypt, _ := rsaObj.Decrypt([]byte(encrypt))
		assert.Equal(t, data, decrypt)
	})
	// 加解密特殊字符
	Convey("RSADecrypt()", t, func() {
		data := "asAlqlTkWU0zqfxrLTed！@#￥%……&*（）—+魔鬼"
		rsaObj := NewRsa()
		encrypt, _ := rsaObj.Encrypt([]byte(data))
		decrypt, _ := rsaObj.Decrypt([]byte(encrypt))
		assert.Equal(t, data, decrypt)
	})
	// 超长密码测试
	Convey("RSADecryptLongString()", t, func() {
		data := strings.Repeat("H", 24270) + "e"
		rsaObj := NewRsa()
		encrypt, _ := rsaObj.Encrypt([]byte(data))
		decrypt, _ := rsaObj.Decrypt([]byte(encrypt))
		assert.Equal(t, data, decrypt)
	})
}
