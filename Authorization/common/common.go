// Package common 通用函数模块
package common

import (
	"fmt"

	"Authorization/interfaces"
)

// AccessorTypeToStr 访问者类型转换
func AccessorTypeToStr(accessorType interfaces.AccessorType) (accessorTypeToResStr string) {
	var ok bool
	if accessorTypeToResStr, ok = AccessorTypeToStringMap[accessorType]; ok {
		return
	}
	return "unknow"
}

// AccessorTypeToStringMap 访问者类型转换图
var AccessorTypeToStringMap = map[interfaces.AccessorType]string{
	interfaces.AccessorUser:       "user",
	interfaces.AccessorDepartment: "department",
	interfaces.AccessorGroup:      "group",
}

// StrToAccessorTypeMap 权限给予类型常量
var StrToAccessorTypeMap = map[string]interfaces.AccessorType{
	"user":          interfaces.AccessorUser,
	"department":    interfaces.AccessorDepartment,
	"contactor":     interfaces.AccessorContactor,
	"anonymoususer": interfaces.AccessorAnonymous,
	"group":         interfaces.AccessorGroup,
}

// AccessorTypeStrToInt 访问者类型字符串转int
func AccessorTypeStrToInt(accessorType string) interfaces.AccessorType {
	switch accessorType {
	case "user":
		return interfaces.AccessorUser
	case "department":
		return interfaces.AccessorDepartment
	case "contactor":
		return interfaces.AccessorContactor
	case "anonymoususer":
		return interfaces.AccessorAnonymous
	case "group":
		return interfaces.AccessorGroup
	default:
		panic(fmt.Sprintf("invalid accessor type %s", accessorType))
	}
}
