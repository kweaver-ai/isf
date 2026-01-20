package receo

import (
	"time"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/domain/entity/oprlogeo"
)

func GetEosByOprLogEos(eos []*oprlogeo.LogEntry) (recEos []*RecLogEntry) {
	createdTime := time.Now().UnixNano() / 1e9

	recEos = make([]*RecLogEntry, 0, len(eos))

	for _, eo := range eos {
		reo := NewRecLogEntry(eo, createdTime)
		recEos = append(recEos, reo)
	}

	return
}

func GetDocs(logMap []map[string]interface{}) (rets []map[string]interface{}) {
	rets = make([]map[string]interface{}, 0)

	now := time.Now()
	timeUtcStr := now.UTC().Format("2006-01-02T15:04:05.000Z")

	for _, log := range logMap {
		bodyMap := make(map[string]interface{})
		bizType := log["biz_type"].(string)

		switch bizType {
		case string(oprlogenums.DocOperation):
			bodyMap["Type"] = oprlogenums.DocOperation
			bodyMap[string(oprlogenums.DocOperation)] = log
		case string(oprlogenums.KcOperation):
			bodyMap["Type"] = oprlogenums.KcOperation
			bodyMap[string(oprlogenums.KcOperation)] = log
		case string(oprlogenums.UserLogin):
			bodyMap["Type"] = oprlogenums.UserLogin
			bodyMap[string(oprlogenums.UserLogin)] = log
		case string(oprlogenums.DocFlow):
			bodyMap["Type"] = oprlogenums.DocFlow
			bodyMap[string(oprlogenums.DocFlow)] = log
		case string(oprlogenums.FileCollector):
			bodyMap["Type"] = oprlogenums.FileCollector
			bodyMap[string(oprlogenums.FileCollector)] = log
		case string(oprlogenums.DocumentDomainSync):
			bodyMap["Type"] = oprlogenums.DocumentDomainSync
			bodyMap[string(oprlogenums.DocumentDomainSync)] = log
		case string(oprlogenums.Antivirus):
			bodyMap["Type"] = oprlogenums.Antivirus
			bodyMap[string(oprlogenums.Antivirus)] = log
		case string(oprlogenums.ClientOperation):
			bodyMap["Type"] = oprlogenums.ClientOperation
			bodyMap[string(oprlogenums.ClientOperation)] = log
		case string(oprlogenums.UserFeedback):
			bodyMap["Type"] = oprlogenums.UserFeedback
			bodyMap[string(oprlogenums.UserFeedback)] = log
		case string(oprlogenums.Sap):
			bodyMap["Type"] = oprlogenums.Sap
			bodyMap[string(oprlogenums.Sap)] = log
		case string(oprlogenums.Search):
			bodyMap["Type"] = oprlogenums.Search
			bodyMap[string(oprlogenums.Search)] = log
		case string(oprlogenums.ContentAutomation):
			bodyMap["Type"] = oprlogenums.ContentAutomation
			bodyMap[string(oprlogenums.ContentAutomation)] = log
		case string(oprlogenums.ContentProcessFinish):
			bodyMap["Type"] = oprlogenums.ContentProcessFinish
			bodyMap[string(oprlogenums.ContentProcessFinish)] = log
		case string(oprlogenums.MenuButtonClick):
			bodyMap["Type"] = oprlogenums.MenuButtonClick
			bodyMap[string(oprlogenums.MenuButtonClick)] = log
		case string(oprlogenums.DirVisit):
			bodyMap["Type"] = oprlogenums.DirVisit
			bodyMap[string(oprlogenums.DirVisit)] = log
		}

		ret := make(map[string]interface{})
		ret["@timestamp"] = timeUtcStr
		ret["timestamp_int_nano"] = now.UnixNano()
		ret["Body"] = bodyMap
		rets = append(rets, ret)
	}

	return
}
