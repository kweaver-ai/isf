/* 定义 输入适配器——核心逻辑层 有关的struct */
package models

// 发送日志
type SendLogVo struct {
	Language   string
	LogType    string
	LogContent *AuditLog
}

type ReceiveLogVo struct {
	Language   string
	LogType    string
	LogContent *AuditLog
}

// 角色成员信息
type RoleMemberInfo struct {
	Role    string       `json:"role"`
	Members []MemberInfo `json:"members"`
}

// 角色成员信息
type MemberInfo struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// OssInfo 对象存储信息
type OSSInfo struct {
	Default bool   `json:"default"`
	Enabled bool   `json:"enabled"`
	ID      string `json:"id"`
	Name    string `json:"name"`
}

// OSSUploadInfo 上传信息结构体
type OSSUploadInfo struct {
	UploadID   string `json:"upload_id"`
	UploadType string `json:"upload_type"`
	PartSize   int    `json:"partsize"`
	MaxNum     int    `json:"max_num"`
}

// OSSRequestInfo 上传请求信息
type OSSRequestInfo struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	RequestBody string            `json:"request_body"`
}

// OSSUploadPartInfo 上传分片信息
type OSSUploadPartInfo struct {
	Etag string `json:"etag"`
	Size int    `json:"size"`
}

type KcUserInfo struct {
	ID          string   `json:"user_id"`
	Name        string   `json:"user_name"`
	Birthday    string   `json:"birthday"`
	Category    []string `json:"category"`
	Address     string   `json:"address"`
	Native      string   `json:"native"`
	Graduated   []string `json:"graduated"`
	Certificate []string `json:"certificate"`
}

type KcUserInfoRes struct {
	Code    int          `json:"code"`
	Message string       `json:"msg"`
	Data    []KcUserInfo `json:"data"`
}

type ServiceModuleInfo struct {
	Name     string        `json:"name"`
	Version  string        `json:"version"`
	Services []ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Type     int    `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Path     string `json:"path"`
	Extra    any    `json:"extra"`
}
