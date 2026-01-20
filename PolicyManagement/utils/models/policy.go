package models

import "github.com/open-policy-agent/opa/storage"

// Policy 策略表
type Policy[T []byte | string] struct {
	Name    string `gorm:"column:f_name;size:255;primary_key;not null"`
	Default T      `gorm:"column:f_default;type:text;not null"`
	Value   T      `gorm:"column:f_value;type:text;not null"`
	Locked  bool   `gorm:"column:f_locked"`
}

// TableName 表名
func (Policy[T]) TableName() string {
	return "t_policies"
}

// 非数据库结构体
// ClientRestrictionConfig 客户端登录选项配置值
type ClientRestrictionConfig struct {
	Unknown    bool `json:"unknown"`
	PcWeb      bool `json:"web"`
	MobileWeb  bool `json:"mobile_web"`
	Windows    bool `json:"windows"`
	Mac        bool `json:"mac_os"`
	Android    bool `json:"android"`
	IOS        bool `json:"ios"`
	WinPhone   bool `json:"windows_phone"`
	ConsoleWeb bool `json:"console_web"`
	DeployWeb  bool `json:"deploy_web"`
	NAS        bool `json:"nas"`
	Linux      bool `json:"linux"`
}

type OPADataPatch struct {
	PatchOP storage.PatchOp
	Path    string
	Value   interface{}
}
