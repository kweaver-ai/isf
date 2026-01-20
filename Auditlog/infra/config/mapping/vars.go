package mapping

import (
	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/types/rectypes"
)

var RecMappingMap map[rectypes.RecLogIndex]string

func init() {
	RecMappingMap = make(map[rectypes.RecLogIndex]string)

	// // 1. doc operation
	// docOp := rectypes.GetIndexByOprBizType(oprlogenums.DocOperation, "")
	// RecMappingMap[docOp] = GetDocOperationCommonMapping()

	// // 2. dir visit
	// dirVisit := rectypes.GetIndexByOprBizType(oprlogenums.DirVisit, "cd")
	// RecMappingMap[dirVisit] = GetDirVisitCdMapping()

	// // 3. menu button click
	// mbc := rectypes.GetIndexByOprBizType(oprlogenums.MenuButtonClick, "")
	// RecMappingMap[mbc] = GetMbcMapping()

	// 4. docOperation
	docOperation := rectypes.GetIndexByOprBizType(oprlogenums.DocOperation, "")
	RecMappingMap[docOperation] = GetDocOperationMapping()

	// 5. kcOperation
	kcOperation := rectypes.GetIndexByOprBizType(oprlogenums.KcOperation, "")
	RecMappingMap[kcOperation] = GetKcOperationMapping()

	antivirus := rectypes.GetIndexByOprBizType(oprlogenums.Antivirus, "")
	RecMappingMap[antivirus] = GetAntivirusMapping()

	clientOperation := rectypes.GetIndexByOprBizType(oprlogenums.ClientOperation, "")
	RecMappingMap[clientOperation] = GetClientOperationMapping()

	contentAutomation := rectypes.GetIndexByOprBizType(oprlogenums.ContentAutomation, "")
	RecMappingMap[contentAutomation] = GetContentAutomationMapping()

	contentProcessFinish := rectypes.GetIndexByOprBizType(oprlogenums.ContentProcessFinish, "")
	RecMappingMap[contentProcessFinish] = GetContentProcessFinishMapping()

	dirVisit := rectypes.GetIndexByOprBizType(oprlogenums.DirVisit, "")
	RecMappingMap[dirVisit] = GeDirVisitMapping()

	docFlow := rectypes.GetIndexByOprBizType(oprlogenums.DocFlow, "")
	RecMappingMap[docFlow] = GetDocFlowMapping()

	documentDomainSync := rectypes.GetIndexByOprBizType(oprlogenums.DocumentDomainSync, "")
	RecMappingMap[documentDomainSync] = GetDocumentDomainSyncMapping()

	fileCollector := rectypes.GetIndexByOprBizType(oprlogenums.FileCollector, "")
	RecMappingMap[fileCollector] = GetFileCollectorMapping()

	menuButtonClick := rectypes.GetIndexByOprBizType(oprlogenums.MenuButtonClick, "")
	RecMappingMap[menuButtonClick] = GetMenuButtonClickMapping()

	sap := rectypes.GetIndexByOprBizType(oprlogenums.Sap, "")
	RecMappingMap[sap] = GetSapMapping()

	search := rectypes.GetIndexByOprBizType(oprlogenums.Search, "")
	RecMappingMap[search] = GetSearchMapping()

	userFeedback := rectypes.GetIndexByOprBizType(oprlogenums.UserFeedback, "")
	RecMappingMap[userFeedback] = GetUserFeedbackMapping()

	userLogin := rectypes.GetIndexByOprBizType(oprlogenums.UserLogin, "")
	RecMappingMap[userLogin] = GetUserLoginMapping()
}
