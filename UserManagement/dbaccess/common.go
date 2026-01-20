// Package dbaccess Anyshare 数据访问层 -通用
package dbaccess

import (
	"os"
	"strings"
	"time"

	"UserManagement/common"
	"UserManagement/interfaces"
)

// GetFindInSetSQL 获取集合查询部分的SQL
func GetFindInSetSQL(value []string) (sql string, args []interface{}) {
	set := make([]string, 0)
	for _, v := range value {
		set = append(set, "?")
		args = append(args, v)
	}
	sql = strings.Join(set, ",")

	return
}

// GetUUIDStringBySlice 获取字符串数组的string
func GetUUIDStringBySlice(value []string) (sql string) {
	set := make([]string, 0, len(value))
	for _, v := range value {
		set = append(set, "'"+v+"'")
	}
	sql = strings.Join(set, ",")

	return
}

// SplitArray 拆分数组，保证in值列表限制在500以内
func SplitArray(arr []string) [][]string {
	length := 500
	total := len(arr)
	count := total / length
	if total%length != 0 {
		count++
	}

	resArray := make([][]string, 0)
	start := 0
	end := 0
	for i := 0; i < count; i++ {
		end = (i + 1) * length
		if i != (count - 1) {
			resArray = append(resArray, arr[start:end])
		} else {
			resArray = append(resArray, arr[start:])
		}
		start = end
	}

	return resArray
}

// 用户禁用状态int type转换map
var disableIntToTypeMap = map[int]interfaces.DisableType{
	0: interfaces.Enabled,
	1: interfaces.Disabled,
	2: interfaces.Deleted,
}

// 用户自动禁用状态int type转换map
var autoDisableIntToTypeMap = map[int]interfaces.AutoDisableType{
	0: interfaces.AEnabled,
	1: interfaces.ADisabled,
	2: interfaces.ExpireDisabled,
}

// 用户认证类型int type转换map
var userTypeIntToTypeMap = map[int]interfaces.AuthType{
	1: interfaces.Local,
	2: interfaces.Domain,
	3: interfaces.Third,
}

// LDAP类型int type转换map
var ldapTypeIntToTypeMao = map[int]interfaces.LDAPServerType{
	1: interfaces.WindowAD,
	2: interfaces.OtherLDAP,
}

const (
	windowAD  = 1
	otherLDAP = 2
)

var (
	// 获取当前环境数据库类型，MariaDB、MySQL、GoldenDB、DM8、TiDB
	// 详情参考: https://confluence.aishu.cn/pages/viewpage.action?pageId=200147371
	dbType = os.Getenv("DB_TYPE")
)

// handlerUserDBData 数据切换
func handlerUserDBData(data *userDBData) (out interfaces.UserDBInfo) {
	out.ID = data.ID
	out.Name = data.Name
	out.Account = data.Account
	out.CSFLevel = data.CSFLevel
	out.CSFLevel2 = data.CSFLevel2
	out.Priority = data.Priority
	out.DisableStatus = disableIntToTypeMap[data.DisableStatus]
	out.AutoDisableStatus = autoDisableIntToTypeMap[data.AutoDisableStatus]
	out.Email = data.Email
	out.AuthType = userTypeIntToTypeMap[data.AuthType]
	out.Password = data.Password
	out.DesPassword = data.DesPassword
	out.NtlmPassword = data.NtlmPassword
	out.Sha2Password = data.Sha2Password
	out.Frozen = data.Frozen != 0
	out.Authenticated = data.Authenticated != 0
	out.TelNumber = data.TelNumber
	out.ThirdAttr = data.ThirdAttr
	out.ThirdID = data.ThirdID
	out.PWDControl = data.PWDControl
	if data.PWDTimeStamp != "" {
		out.PWDTimeStamp = parseStrToTime(data.PWDTimeStamp).Unix()
	}
	out.PWDErrCnt = data.PWDErrCnt
	if data.PWDErrLatestTime != "" {
		out.PWDErrLatestTime = parseStrToTime(data.PWDErrLatestTime).Unix()
	}

	// 根据ShareMgnt域认证接口逻辑，暂且分为 windowAD 和 otherLDAP 两种情况，后续有需要再扩展。
	if data.LDAPType != windowAD {
		data.LDAPType = otherLDAP
	}
	out.LDAPType = ldapTypeIntToTypeMao[data.LDAPType]
	out.DomainPath = data.DomainPath
	out.OssID = data.OssID
	out.ManagerID = data.ManagerID
	out.Remark = data.Remark
	out.Code = data.Code
	out.Position = data.Position

	if data.CreatedAtTimeStamp != "" {
		out.CreatedAtTimeStamp = parseStrToTime(data.CreatedAtTimeStamp).Unix()
	}

	return
}

// parseStrToTime 解析时间字符串
// 不同数据库返回的时间字符串类型可能不同
func parseStrToTime(timeStr string) (parsedTime time.Time) {
	var err error
	if dbType == "DM8" {
		// 达梦 返回时间的格式是 RFC3339
		parsedTime, err = time.ParseInLocation(time.RFC3339, timeStr, time.Local)
		if err != nil {
			common.NewLogger().Errorf("parse time as RFC3339 failed, timeStr=%s, err=%v\n", timeStr, err)
		}
	} else if strings.HasPrefix(dbType, "KDB") {
		// KDB 返回时间的格式是 RFC3339
		parsedTime, err = time.ParseInLocation(time.RFC3339, timeStr, time.Local)
		if err != nil {
			common.NewLogger().Errorf("KDB parse time as RFC3339 failed, timeStr=%s, err=%v\n", timeStr, err)
		}
	} else {
		parsedTime, err = time.ParseInLocation(time.DateTime, timeStr, time.Local)
		if err != nil {
			common.NewLogger().Errorf("parse time as DateTime failed, timeStr=%s, err=%v\n", timeStr, err)
		}
	}

	return parsedTime
}
