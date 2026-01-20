package httpclientcmp

import (
	"AuditLog/common/utils"
)

type DetailMap = map[string]interface{}

// CommonResp 通用响应
// 参考：https://confluence.aishu.cn/pages/viewpage.action?pageId=190114672
type CommonResp struct {
	Code        int       `json:"code"`        // 错误码（前三位：标准http错误码，中间三位为服务器特定码，后三位服务中自定义码）
	Cause       string    `json:"cause"`       // 错误原因，产生错误的具体原因
	Message     string    `json:"message"`     // 错误信息
	Description string    `json:"description"` // 符合国际化要求的错误描述
	Solution    string    `json:"solution"`    // 符合国际化要求的针对当前错误的操作提示
	Detail      DetailMap `json:"detail"`
	Debug       string    `json:"debug,omitempty"` // CAPP部分项目用到，其它项目可能没有这个
}

type CommonRespError CommonResp

func (r *CommonRespError) Error() string {
	bys, err := utils.JSON().Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bys)
}

const (
	RetryInterval = 5
)
