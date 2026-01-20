package logics

import (
	"errors"
	"testing"

	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testPWD = "123123"
)

func TestDecodeRSA(t *testing.T) {
	Convey("RSA1024解密测试", t, func() {
		pwd := "KnSqccgPMx\r\nQN+ob5XXq0uh7Yitp\rqYktCT0c5ql/XtRacyD8Pv\npCyMsc4V14cngr2QxnAuzFCzfOVGoqPypAWKiVzfC1RcvF9lHK8wB7qI8+cw7yPN25J490ibMMBIGNuy5W+aM4gy8eFVVAxuAC8Gb4wNzQ6pKlMQ4MnEcI3X48="
		out, err := decodeRSA(pwd, RSA1024)
		assert.Equal(t, out, testPWD)
		assert.Equal(t, err, nil)
	})

	Convey("RSA2048解密测试", t, func() {
		pwd := "zhMytQF/dmSfreO1Qgkdr8wBEtzi/2QcFwroQV8y+AnFjqhS6aAkVExtgk1VpjwtBk6DlmtSedTFngRLbc61aDQKKJhXTtGYucnRGgOOqD2uu" +
			"I+MaxxAk5t7Vys29XzEyHXB5OAvETjfxkNV/5jAxmQ8k29NDraxpz/yhZ/SsnviskBaGE+l/n+7EvhL2VIVhf9Yp3FB96tOxrjfApf+7a0iIN5NgM+5YjazKnN8nHAJ5Em" +
			"SINBbb+nK+7ciC+IkEBLXRms5Hv5KWpUdPP23iw55Nl3ffjXvqtUyCfVBOqItWDAd32DA7U8Qg8Ver7Tn3wScuGokNijmis0dlEbRVw=="
		out, err := decodeRSA(pwd, RSA2048)
		assert.Equal(t, out, "1111")
		assert.Equal(t, err, nil)
	})

	Convey("type error", t, func() {
		pwd := ""
		out, err := decodeRSA(pwd, 34)
		assert.Equal(t, out, "")
		assert.Equal(t, err, errors.New("decode type error"))
	})
}

func TestEncodeMD5(t *testing.T) {
	Convey("MD5加密测试 ", t, func() {
		pwd := testPWD
		out, err := encodeMD5(pwd)
		assert.Equal(t, out, "4297f44b13955235245b2497399d7a93")
		assert.Equal(t, err, nil)
	})
}

func TestEncodeNtlm(t *testing.T) {
	Convey("ntlm加密测试 ", t, func() {
		pwd := testPWD
		out := encodeNtlm(pwd)
		assert.Equal(t, out, "579110c49145015c47ecd267657d3174")
	})
}

func TestEncodeDes(t *testing.T) {
	Convey("des加密测试 PKCS5Padding", t, func() {
		pwd := testPWD
		out, err := encodeDes(pwd, PKCS5Padding)
		assert.Equal(t, out, "LT8h+r7l8nw=")
		assert.Equal(t, err, nil)
	})

	Convey("des加密测试 ", t, func() {
		pwd := "111111111111111111"
		out, err := encodeDes(pwd, PadNormal)
		assert.Equal(t, out, "4eHALvs2fwKMU5Df/J2WC86i4WMH1sv7")
		assert.Equal(t, err, nil)
	})
}
