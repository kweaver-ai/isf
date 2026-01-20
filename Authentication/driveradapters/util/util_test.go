package util

import (
	"testing"

	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIsEmpty(t *testing.T) {
	Convey("string is empty", t, func() {
		str := map[string]string{
			"aa": "",
			"bb": "bbb",
		}
		err := IsEmpty(str)
		assert.NotEqual(t, err, nil)
	})

	Convey("success", t, func() {
		str := map[string]string{
			"aa": "aaa",
			"bb": "bbb",
		}
		err := IsEmpty(str)
		assert.Equal(t, err, nil)
	})
}
