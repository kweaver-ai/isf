package recvars

import (
	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/utils"
)

var RecBizTypes = []oprlogenums.BizType{
	oprlogenums.DocOperation,
	oprlogenums.DirVisit,
	oprlogenums.MenuButtonClick,
	oprlogenums.KcOperation,
	oprlogenums.UserLogin,
	oprlogenums.DocFlow,
	oprlogenums.FileCollector,
	oprlogenums.DocumentDomainSync,
	oprlogenums.Antivirus,
	// oprlogenums.SensitiveContent,
	oprlogenums.UserFeedback,
	oprlogenums.Sap,
	oprlogenums.Search,
	oprlogenums.ContentProcessFinish,
	oprlogenums.ContentAutomation,
	// oprlogenums.IntelliQaChat,
	oprlogenums.ClientOperation,
}

func IsRecBizType(bizType oprlogenums.BizType) bool {
	return utils.ExistsGeneric(RecBizTypes, bizType)
}
