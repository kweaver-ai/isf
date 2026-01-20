package simu

import (
	"github.com/agiledragon/gomonkey"
	jsoniter "github.com/json-iterator/go"
)

// AsJSONMarshal jsoniter Marshal模拟
func AsJSONMarshal(byt []byte, err error) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func() {
			p := gomonkey.ApplyFunc(jsoniter.Marshal, func(v interface{}) ([]byte, error) {
				return byt, err
			})
			defer p.Reset()
			next()
		}
	}
}

// AsJSONUnmarshal jsoniter unmarshal 模拟
func AsJSONUnmarshal(err error) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func() {
			p := gomonkey.ApplyFunc(jsoniter.Unmarshal, func(data []byte, v interface{}) error {
				return err
			})
			defer p.Reset()
			next()
		}
	}
}
