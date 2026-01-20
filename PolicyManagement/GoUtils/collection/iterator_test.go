package collection

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestForEachField(t *testing.T) {
	assert := assert.New(t)

	Convey("ForEachField 遍历结构体字段", t, func() {
		Convey("支持遍历基本的数据类型", func() {
			type Student struct {
				Name   string
				Age    int
				Score  float32
				Passed bool
			}

			peter := Student{
				Name:   "Zak",
				Age:    34,
				Score:  59.5,
				Passed: false,
			}

			ForEachField(peter, func(field string, value interface{}) {
				switch field {
				case "Name":
					assert.Equal("Zak", value)
				case "Age":
					assert.Equal(int64(34), value)
				case "Score":
					assert.Equal(59.5, value)
				case "Passed":
					assert.Equal(false, value)
				}
			})
		})

		Convey("不支持的类型触发panic", func() {
			somewhere := struct {
				Coord []float32
			}{
				Coord: []float32{35.1132121, 26.3212131},
			}

			assert.PanicsWithValue(`Type "slice" of field "Coord" is not supported`, func() {
				ForEachField(somewhere, func(field string, value interface{}) {
					switch field {
					case "Coord":
						assert.Equal("Zak", value)
					}
				})
			})
		})
	})
}
