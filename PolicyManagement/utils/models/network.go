package models

import (
	"math/big"
	"time"
)

// NetworkRestriction 访问者白名单
type NetworkRestriction struct {
	ID           string    `gorm:"column:f_id;size:36;primary_key:not null" json:"id"`
	Name         string    `gorm:"column:f_name;size:128;unique_index;default: null" json:"name,omitempty"`
	StartIP      string    `gorm:"column:f_start_ip;size:40;index;not null" json:"start_ip,omitempty"`
	EndIP        string    `gorm:"column:f_end_ip;size:40;index;not null" json:"end_ip,omitempty"`
	IPAddress    string    `gorm:"column:f_ip_address;size:40;index;not null" json:"ip_address,omitempty"`
	IPMask       string    `gorm:"column:f_ip_mask;size:15;not null" json:"netmask,omitempty"`
	SegmentStart string    `gorm:"column:f_segment_start;size:128;not null" json:"start,omitempty"`
	SegmentEnd   string    `gorm:"column:f_segment_end;size:128;not null" json:"end,omitempty"`
	Type         string    `gorm:"column:f_type;size:15;index;not null" json:"net_type"`
	IpType       string    `gorm:"column:f_ip_type;size:15;not null;default:'ipv4'" json:"ip_type"`
	CreatedAt    time.Time `gorm:"column:f_created_at;precision:3" json:"f_created_at"`
}

// WebNetworkRestriction 返回给前端的访问者白名单
type WebNetworkRestriction struct {
	*NetworkRestriction
	SegmentStart *struct{} `json:"start,omitempty"`
	SegmentEnd   *struct{} `json:"end,omitempty"`
	CreatedAt    *struct{} `json:"f_created_at,omitempty"`
}

// TableName 表名
func (NetworkRestriction) TableName() string {
	return "t_network_restriction"
}

// NetworkAccessorRelation 白名单、访问者关系表
type NetworkAccessorRelation struct {
	ID           int64     `gorm:"column:f_id;primary_key;AUTO_INCREMENT;not null"`
	NetworkId    string    `gorm:"column:f_network_id;size:36;unique_index:idx_net_acc;not null"`
	AccessorId   string    `gorm:"column:f_accessor_id;size:36;unique_index:idx_net_acc;not null"`
	AccessorType string    `gorm:"column:f_accessor_type;size:10;index;not null"`
	CreatedAt    time.Time `gorm:"column:f_created_at;precision:3"`
}

// TableName 表名
func (NetworkAccessorRelation) TableName() string {
	return "t_network_accessor_relation"
}

// 非数据库结构体
// 访问者信息
type AccessorInfo struct {
	AccessorId   string `json:"accessor_id"`
	AccessorName string `json:"accessor_name"`
	AccessorType string `json:"accessor_type"` // user为用户，department为部门
}

// 网段的起始值
type UserNetworkInfo struct {
	Departments  []string     `json:"departments,omitempty"`
	NetSegements []NetSegment `json:"nets,omitempty"`
}

// 网段的起始值
type NetSegment struct {
	StartIP interface{} `json:"start_ip"`
	EndIP   interface{} `json:"end_ip"`
}

// 提供给OPA网段白名单信息
type OPANetworks struct {
	IsEnabled bool                   `json:"is_enabled"`
	Accessors map[string]interface{} `json:"accessors"`
}

// OPA增量更新所需查询数据
type OPAUpdateInfo struct {
	AccessorId string `json:"f_accessor_id"`
	IpType     string `json:"f_ip_type"`
}

// OPA决策请求访问者
type Accessor struct {
	AccessorId string   `json:"accessor_id"`
	IpType     string   `json:"ip_type"`
	Ip         *big.Int `json:"ip"`
}

// 用户移动结构体
type UserDepartRelation struct {
	UserId          string   `json:"user_id"`
	DepartmentPaths []string `json:"department_paths"`
}
