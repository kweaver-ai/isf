package api

import (
	"strconv"
)

type MultiStatus struct {
	ID     string      `json:"id"`               // 数据的唯一标识，必须参数
	Status int         `json:"status"`           // HTTP状态码，必须参数
	Header interface{} `json:"header,omitempty"` // HTTP报头，类型为object，可选参数
	Body   interface{} `json:"body,omitempty"`   // HTTP主体，类型为object或者array，可选参数
}

// MultiStatusObject 生成多状态对象
func MultiStatusObject(id string, header interface{}, body interface{}, status ...int) *MultiStatus {
	var statusCode int
	// 如果没有传status， 则解析错误码中的code
	if len(status) == 0 {
		if apiErr, ok := body.(*Error); ok {
			statusStr := strconv.Itoa(apiErr.Code)
			statusCode, _ = strconv.Atoi(statusStr[:3])
		} else {
			statusCode = 500
		}
	} else {
		statusCode = status[0]
	}
	return &MultiStatus{
		ID:     id,
		Status: statusCode,
		Header: header,
		Body:   body,
	}
}
