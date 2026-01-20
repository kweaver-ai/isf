package rectypes

import (
	recconsts "AuditLog/common/constants/recenums"
	"AuditLog/common/enums/oprlogenums"
)

type RecLogIndex string

func GetIndexByOprBizType(bizType oprlogenums.BizType, operation string) (index RecLogIndex) {
	idx := recconsts.IndexPrefix + string(bizType)
	//
	//if operation != "" {
	//	if bizType == oprlogenums.DirVisit {
	//		idx = idx + "_" + operation
	//	}
	//}

	index = RecLogIndex(idx)

	return
}
