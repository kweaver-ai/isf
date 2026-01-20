package logics

import (
	"time"

	gerrors "github.com/kweaver-ai/go-lib/error"

	"Authorization/interfaces"
)

// checkVisitorType 检测访问者类型
func checkVisitorType(visitor *interfaces.Visitor, roleTypes []interfaces.SystemRoleType, acceptVisitorTypes []interfaces.VisitorType, acceptRoleTypes []interfaces.SystemRoleType) (err error) {
	isVisitorTypePass := false
	isRoleTypePass := false
	// 访问者类型是否在检测范围内
	for _, val := range acceptVisitorTypes {
		if visitor.Type == val {
			isVisitorTypePass = true
		}
	}
	if !isVisitorTypePass {
		err = gerrors.NewError(gerrors.PublicForbidden, "Unsupported user type")
		return
	}

	if visitor.Type == interfaces.RealName {
		// 访问者角色是否在检测范围内
		for _, roleType := range roleTypes {
			for _, val := range acceptRoleTypes {
				if roleType == val {
					isRoleTypePass = true
				}
			}
		}
		if !isRoleTypePass {
			err = gerrors.NewError(gerrors.PublicForbidden, "Unsupported user role type")
			return
		}
	}
	return nil
}

func checkEndTime(endTime int64) (err error) {
	// 检查 截止时间
	curTime := time.Now().UnixNano() / 1000 // 获取当前时间
	if endTime != -1 && endTime <= curTime {
		return gerrors.NewError(gerrors.PublicBadRequest, "The expiration time of permission cannot be earlier than the current time")
	}
	return
}
