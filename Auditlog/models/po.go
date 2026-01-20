/* 定义 核心逻辑层——输出适配器 有关的struct */
package models

type OutboxPo struct {
	ID          uint64 `gorm:"column:f_id;primary_key:not null"`
	Type        string `gorm:"column:f_mail_type;type:varchar(20)"`
	Destination []byte `gorm:"column:f_destination;type:text"`
	Content     []byte `gorm:"column:f_mail_content;type:text"`
	CreateTime  int64  `gorm:"column:f_create_time"`
}

// http://{host}:{port}/api/document/v1/batch-doc-libs/{fields} 请求体
type GetBatchDocLibRequest struct {
	Method string   `json:"method"`
	IDs    []string `json:"ids"`
}

type LogPO struct {
	LogID          string `gorm:"column:f_log_id;primary_key;not null" json:"logId"`                           // 日志id
	UserID         string `gorm:"column:f_user_id;type:char(40);not null" json:"userId"`                       // 用户id
	UserName       string `gorm:"column:f_user_name;type:char(128);not null" json:"userName"`                  // 用户显示名
	ObjID          string `gorm:"column:f_obj_id;type:char(40);not null" json:"objId"`                         // 对象id
	AdditionalInfo string `gorm:"column:f_additional_info;type:text;not null" json:"additionalInfo"`           // 附加信息
	Level          int    `gorm:"column:f_level;type:tinyint(4);not null" json:"level"`                        // 日志级别，1：信息，2：警告
	OpType         int    `gorm:"column:f_op_type;type:tinyint(4);not null" json:"opType"`                     // 日志类型
	Date           int64  `gorm:"column:f_date;type:bigint(20);not null" json:"date"`                          // 日志记录时间，微秒的时间戳
	IP             string `gorm:"column:f_ip;type:char(40);not null" json:"ip"`                                // 访问者的ip
	MAC            string `gorm:"column:f_mac;type:char(40);not null;default:''" json:"mac"`                   // 文档入口所属站点
	Msg            string `gorm:"column:f_msg;type:text;not null" json:"msg"`                                  // 日志描述
	ExMsg          string `gorm:"column:f_exmsg;type:text;not null" json:"exMsg"`                              // 日志附加描述
	UserAgent      string `gorm:"column:f_user_agent;type:varchar(1024);not null;default:''" json:"userAgent"` // 用户代理
	UserPaths      string `gorm:"column:f_user_paths;type:text" json:"userPaths"`                              // 用户类型
	ObjName        string `gorm:"column:f_obj_name;type:char(128);not null;default:''" json:"objName"`         // 对象名称
	ObjType        int    `gorm:"column:f_obj_type;type:tinyint(4);not null;default:0" json:"objType"`         // 对象类型
}

type HistoryPO struct {
	ID       string `gorm:"column:f_id;type:char(128);not null"`
	Name     string `gorm:"column:f_name;type:char(128);not null"`
	Size     int64  `gorm:"column:f_size;type:bigint;not null"`
	Type     int8   `gorm:"column:f_type;type:tinyint;not null"`
	Date     int64  `gorm:"column:f_date;type:bigint;not null"`
	DumpDate int64  `gorm:"column:f_dump_date;type:bigint;not null"`
	OssID    string `gorm:"column:f_oss_id;type:char(40);not null"`
}

type PersonalizedFeatureConfigPo struct {
	ID             int64  `gorm:"column:f_id;type:bigint(20);not null" json:"id"`
	Key            string `gorm:"column:f_key;type:char(36);not null" json:"key"`
	Type           string `gorm:"column:f_type;type:varchar(64);not null" json:"type"`
	Name           string `gorm:"column:f_name;type:char(128);not null" json:"name"`
	BuildIn        int    `gorm:"column:f_is_build_in;type:tinyint(1);not null" json:"is_built_in"`
	ApplicableType string `gorm:"column:f_applicable_type;type:varchar(32);not null" json:"applicable_type"`
	Structure      string `gorm:"column:f_structure;type:text;not null" json:"structure"`
	TopNum         int    `gorm:"column:f_top_num;type:int(11);not null" json:"top_num"`
	PeriodTime     int    `gorm:"column:f_period_time;type:int(11);not null" json:"period_time"`
	CreateBy       string `gorm:"column:f_created_by;type:char(36);not null" json:"created_by"`
	CreateAt       int64  `gorm:"column:f_created_at;type:bigint(20);not null" json:"created_at"`
	UpdateBy       string `gorm:"column:f_updated_by;type:char(36);not null" json:"updated_by"`
	UpdateAt       int64  `gorm:"column:f_updated_at;type:bigint(20);not null" json:"updated_at"`
}

type PersonalizedFeaturePo struct {
	ID       int64  `gorm:"column:f_id;type:bigint(20);not null"`
	ObjID    string `gorm:"column:f_obj_id;type:varchar(64);not null"`
	ObjType  string `gorm:"column:f_obj_type;type:varchar(20);not null"`
	Key      string `gorm:"column:f_key;type:char(36);not null"`
	Value    []byte `gorm:"column:f_value;type:text;not null"`
	CreateAt int64  `gorm:"column:f_created_at;type:bigint(20);not null"`
	UpdateAt int64  `gorm:"column:f_updated_at;type:bigint(20);not null"`
}
