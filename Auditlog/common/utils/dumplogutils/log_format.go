package dumplogutils

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"

	"AuditLog/common"
	"AuditLog/common/constants"
	"AuditLog/common/utils"
	"AuditLog/locale"
	"AuditLog/models"
)

// LogInfo2CSVString 将日志信息转换为CSV格式字符串
func LogInfo2CSVString(log *models.LogPO, logType string) (string, error) {
	additionalInfo, err := formatAdditionalInfo(log.AdditionalInfo, log.ObjID)
	if err != nil {
		return "", fmt.Errorf("format additional info failed: %w", err)
	}

	// 预先准备所有字段
	fields := []string{
		utils.FormatTime(time.UnixMicro(log.Date).Local()),
		log.UserName,
		log.UserPaths,
		locale.GetRCLogLevelI18n(context.Background(), locale.LogLevelMap[log.Level]),
		formatOpType(logType, log.OpType),
		log.IP,
		log.MAC,
		log.Msg,
		log.ExMsg,
		log.UserAgent,
		additionalInfo,
		log.ObjName,
		locale.GetRCLogObjTypeI18n(context.Background(), log.ObjType),
	}

	var b strings.Builder

	for i, field := range fields {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString(`"`)
		b.WriteString(escapeCSVField(field))
		b.WriteString(`"`)
	}

	return b.String(), nil
}

// LogInfo2XMLString 将日志信息转换为XML格式字符串
func LogInfo2XMLString(log *models.LogPO, logType string) (string, error) {
	additionalInfo, err := formatAdditionalInfo(log.AdditionalInfo, log.ObjID)
	if err != nil {
		return "", fmt.Errorf("format additional info failed: %w", err)
	}

	// XML 模板
	const xmlTemplate = `<log-id Remark="%s">` +
		`<date Remark="%s"> </date>` +
		`<user-name Remark="%s"> </user-name>` +
		`<user-paths Remark="%s"> </user-paths>` +
		`<level Remark="%s"> </level>` +
		`<op-type Remark="%s"> </op-type>` +
		`<ip Remark="%s"> </ip>` +
		`<mac Remark="%s"> </mac>` +
		`<msg Remark="%s"> </msg>` +
		`<ex-msg Remark="%s"> </ex-msg>` +
		`<user-agent Remark="%s"> </user-agent>` +
		`<additional-info Remark="%s"> </additional-info>` +
		`<obj-name Remark="%s"> </obj-name>` +
		`<obj-type Remark="%s"> </obj-type>` +
		`</log-id>`

	// 格式化 XML 字符串
	return fmt.Sprintf(xmlTemplate,
		log.LogID,
		utils.FormatTime(time.UnixMicro(log.Date).Local()),
		html.EscapeString(log.UserName),
		html.EscapeString(log.UserPaths),
		locale.GetRCLogLevelI18n(context.Background(), locale.LogLevelMap[log.Level]),
		formatOpType(logType, log.OpType),
		html.EscapeString(log.IP),
		html.EscapeString(log.MAC),
		html.EscapeString(log.Msg),
		html.EscapeString(log.ExMsg),
		html.EscapeString(log.UserAgent),
		html.EscapeString(additionalInfo),
		html.EscapeString(log.ObjName),
		locale.GetRCLogObjTypeI18n(context.Background(), log.ObjType),
	), nil
}

// formatOpType 根据日志类型和操作类型返回对应的国际化字符串
func formatOpType(logType string, opType int) string {
	ctx := context.Background()

	switch logType {
	case common.Management:
		return locale.GetRCLogMgntI18n(ctx, opType)
	case common.Login:
		return locale.GetRCLogLoginI18n(ctx, opType)
	case common.Operation:
		return locale.GetRCLogOpI18n(ctx, opType)
	}

	return ""
}

// formatAdditionalInfo 处理 additionalInfo
func formatAdditionalInfo(additionalInfo string, objID string) (info string, err error) {
	// 处理 additionalInfo
	additionalInfoItems := make(map[string]interface{})

	if additionalInfo != "" {
		if err := json.Unmarshal([]byte(additionalInfo), &additionalInfoItems); err != nil {
			return "", fmt.Errorf("unmarshal additional info failed: %w", err)
		}
	}

	// 添加 obj_id
	additionalInfoItems["obj_id"] = objID

	// 转回 JSON 字符串
	additionalInfoBytes, err := json.Marshal(additionalInfoItems)
	if err != nil {
		return "", fmt.Errorf("marshal additional info failed: %w", err)
	}

	return string(additionalInfoBytes), nil
}

// escapeCSVField 处理CSV字段中的特殊字符
func escapeCSVField(field string) string {
	// 如果字段包含逗号、双引号或换行符，需要进行转义
	needsQuotes := strings.ContainsAny(field, ",\"\n\r")
	if !needsQuotes {
		return field
	}

	// 将字段中的双引号替换为两个双引号（CSV标准转义方式）
	field = strings.ReplaceAll(field, `"`, `""`)

	return field
}

// 获取System账户
func GetSystemAccount(ctx context.Context) *models.AccountInfo {
	return &models.AccountInfo{
		ID:          constants.SystemID,
		DisplayName: locale.GetI18nCtx(ctx, constants.SystemID),
		Account:     "system",
	}
}
