package collection

import (
	"fmt"
	"reflect"
)

// ForEachField 遍历结构体对象
// 只支持string、int、float、bool、slice，其他类型会引起panic
func ForEachField(s interface{}, fn func(name string, value interface{})) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := val.FieldByName(fieldName)

		switch fieldValue.Kind() {
		case reflect.String:
			fn(fieldName, fieldValue.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fn(fieldName, fieldValue.Int())
		case reflect.Float32, reflect.Float64:
			fn(fieldName, fieldValue.Float())
		case reflect.Bool:
			fn(fieldName, fieldValue.Bool())
		default:
			panic(fmt.Sprintf(`Type "%s" of field "%s" is not supported`, fieldValue.Kind(), fieldName))
		}
	}
}
