package rest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/kweaver-ai/go-lib/simu"
	. "github.com/smartystreets/goconvey/convey"
)

func validationOK(err error) {
	So(err, ShouldBeNil)
}
func validationFail(err error) {
	fmt.Println(err)
	So(err, ShouldNotBeNil)
}

func TestValidationTrans(t *testing.T) {
	type tArgs struct {
		data   string
		object interface{}
	}
	type tCase struct {
		title string
		args  tArgs
		mws   []simu.MiddlewareFunc
		exp   func(error)
	}
	validate := validator.New()
	type element struct {
		A string `json:"a" validate:"required"`
	}
	cases := []tCase{
		{
			title: "错误的json",
			args: tArgs{
				data: "{",
				object: &struct {
					A string `json:"a"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "json类型错误",
			args: tArgs{
				data: `{"a":"1"}`,
				object: &struct {
					A int `json:"a"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "简单属性required, 未赋值不通过",
			args: tArgs{
				data: `{}`,
				object: &struct {
					A int `json:"a" validate:"required"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "简单属性required, 零值不通过",
			args: tArgs{
				data: `{"a":0}`,
				object: &struct {
					A int `json:"a" validate:"required"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "简单属性OK",
			args: tArgs{
				data: `{"a":1}`,
				object: &struct {
					A int `json:"a" validate:"required"`
				}{},
			},
			exp: validationOK,
		},
		{
			title: "简单属性指针类型required, 未赋值不通过",
			args: tArgs{
				data: `{}`,
				object: &struct {
					A *int `json:"a" validate:"required"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "简单属性指针类型required, 零值OK",
			args: tArgs{
				data: `{"a":0}`,
				object: &struct {
					A *int `json:"a" validate:"required"`
				}{},
			},
			exp: validationOK,
		},
		{
			title: "枚举类型不匹配",
			args: tArgs{
				data: `{"a":"c"}`,
				object: &struct {
					A string `json:"a" validate:"required,oneof=a b"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "Slice required 未赋值不通过",
			args: tArgs{
				data: `{}`,
				object: &struct {
					Slice []*element `json:"slice" validate:"required,gt=0,dive"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "Slice 元素个数为0 不通过",
			args: tArgs{
				data: `{"slice":[]}`,
				object: &struct {
					Slice []*element `json:"slice" validate:"required,gt=0,dive"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "Slice 元素简单属性required, 未赋值不通过",
			args: tArgs{
				data: `{"slice":[{}]}`,
				object: &struct {
					Slice []*element `json:"slice" validate:"required,gt=0,dive"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "Slice OK",
			args: tArgs{
				data: `{"slice":[{"a":"1"}]}`,
				object: &struct {
					Slice []*element `json:"slice" validate:"required,gt=0,dive"`
				}{},
			},
			exp: validationOK,
		},
		{
			title: "Map required 未赋值不通过",
			args: tArgs{
				data: `{}`,
				object: &struct {
					Map *element `json:"map" validate:"required"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "Map 简单属性required, 未赋值不通过",
			args: tArgs{
				data: `{"map":{}}`,
				object: &struct {
					Map *element `json:"map" validate:"required"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "Map OK",
			args: tArgs{
				data: `{"map":{"a":"1"}}`,
				object: &struct {
					Map *element `json:"map" validate:"required"`
				}{},
			},
			exp: validationOK,
		},
		{
			title: "其他错误, int 小于设定值",
			args: tArgs{
				data: `{"a":1}`,
				object: &struct {
					A int `json:"map" validate:"gt=2"`
				}{},
			},
			exp: validationFail,
		},
		{
			title: "其他错误, int 大于设定值",
			args: tArgs{
				data: `{"a":1}`,
				object: &struct {
					A int `json:"map" validate:"lt=0"`
				}{},
			},
			exp: validationFail,
		},
	}

	for _, v := range cases {
		tc := v
		simu.ConveyTest(tc.title, t, func() {
			err := json.Unmarshal([]byte(tc.args.data), tc.args.object)
			if err != nil {
				tc.exp(ValidationTrans(err))
				return
			}
			tc.exp(ValidationTrans(validate.Struct(tc.args.object)))
		}, tc.mws...)
	}
}
