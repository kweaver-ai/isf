package oprlogenums

import "AuditLog/common/constants/oprlogconsts"

// BizType 运营日志 业务类型
type BizType string

func (t BizType) ToTopic() OperationLogTopic {
	return OperationLogTopic(oprlogconsts.OLTPrefix + string(t))
}

func (t BizType) Check() bool {
	for _, v := range AllBizType {
		if v == t {
			return true
		}
	}

	return false
}

func (t BizType) IsClientBizType() bool {
	for _, v := range AllClientBizType {
		if v == t {
			return true
		}
	}

	return false
}

func (t BizType) IsServerBizType() bool {
	for _, v := range AllServerBizType {
		if v == t {
			return true
		}
	}

	return false
}
