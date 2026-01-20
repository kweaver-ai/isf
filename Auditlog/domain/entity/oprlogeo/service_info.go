package oprlogeo

import "AuditLog/common/helpers"

// ServiceInfo 服务信息
type ServiceInfo struct {
	Instance struct {
		ID string `json:"id"` // 示例：实例的ID
	} `json:"instance"`
	Name    string `json:"name"`    // 示例：服务名称
	Version string `json:"version"` // 示例：服务版本

	BuildInfo *BuildInfo `json:"build_info,omitempty"` // 构建信息
}

// BuildInfo 构建信息
type BuildInfo struct {
	BranchName string `json:"branch_name"` // 分支名称
	BuildTime  string `json:"build_time"`  // 构建时间
	CommitID   string `json:"commit_id"`   // 提交ID
}

func ToArLogBi(info helpers.BuildInfo) (bi *BuildInfo) {
	bi = &BuildInfo{
		BranchName: info.BranchName,
		BuildTime:  info.BuildTime,
		CommitID:   info.CommitID,
	}

	return
}
