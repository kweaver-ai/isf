package oprlogenums

import (
	"strings"

	"AuditLog/common/constants/oprlogconsts"
)

// OperationLogTopic 运营日志 Topic
// 【注意】：此处的OperationLog为“运营日志”，不是操作日志
type OperationLogTopic string

func (l OperationLogTopic) GetBizType() BizType {
	return BizType(strings.TrimPrefix(string(l), oprlogconsts.OLTPrefix))
}

func GetAllOLT() (tps []OperationLogTopic) {
	tps = make([]OperationLogTopic, 0, len(AllBizType))

	for _, v := range AllBizType {
		tps = append(tps, v.ToTopic())
	}

	return
}
