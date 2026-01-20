package simu

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var errMock = errors.New("mock")

func mockMarshalFn(v interface{}) bool {
	_, err := jsoniter.Marshal(v)
	return err == nil
}

func mockUnmarshalFn(v interface{}) bool {
	err := jsoniter.Unmarshal([]byte{}, v)
	return err == nil
}

func TestMiddleware(t *testing.T) {
	type A struct {
		A string
	}
	ConveyTest("TestMiddleware:Marshal:True", t, func() {
		ok := mockMarshalFn("")
		So(ok, ShouldBeTrue)
	}, AsJSONMarshal([]byte{}, nil))
	ConveyTest("TestMiddleware:Marshal:False", t, func() {
		ok := mockMarshalFn("")
		So(ok, ShouldBeFalse)
	}, AsJSONMarshal([]byte{}, errMock))
	ConveyTest("TestMiddleware:Unmarshal:True", t, func() {
		ok := mockUnmarshalFn(&A{})
		So(ok, ShouldBeTrue)
	}, AsJSONUnmarshal(nil))
	ConveyTest("TestMiddleware:Unmarshal:False", t, func() {
		ok := mockUnmarshalFn(&A{})
		So(ok, ShouldBeFalse)
	}, AsJSONUnmarshal(errMock))
}
