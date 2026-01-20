/* 定义一些通用的struct*/
package models

type AuditLog struct {
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	UserType       string `json:"user_type"`
	Level          int    `json:"level"`
	OpType         int    `json:"op_type"`
	Date           int64  `json:"date"`
	IP             string `json:"ip"`
	Mac            string `json:"mac"`
	Msg            string `json:"msg"`
	Exmsg          string `json:"ex_msg"`
	UserAgent      string `json:"user_agent"`
	ObjID          string `json:"obj_id"`
	AdditionalInfo string `json:"additional_info"`
	OutBizID       string `json:"out_biz_id"`
	DeptPaths      string `json:"dept_paths"`
	ObjType        int
	ObjName        string
}

type Group struct {
	ID   string `json:"id"`   // 用户组ID
	Name string `json:"name"` // 用户组名称
	Type string `json:"type"` // 用户组类型
}

type User struct {
	ID         string        `json:"id"`          // 用户ID
	Roles      []string      `json:"roles"`       // 角色
	Name       string        `json:"name"`        // 名称
	Account    string        `json:"account"`     // 账户名
	ParentDeps []interface{} `json:"parent_deps"` // 部门路径
	Telephone  string        `json:"telephone"`   // 电话
	Email      string        `json:"email"`       // 邮箱
	Groups     []Group       `json:"groups"`      // 用户组
	Level      int           `json:"csf_level"`   // 级别
	Frozen     bool          `json:"frozen"`      // 是否冻结
	Enabled    bool          `json:"enabled"`     // 是否启用
}

type DeptInfo struct {
	ID         string    `json:"department_id"` // 部门ID
	Name       string    `json:"name"`          // 部门名称
	Type       string    `json:"type"`          // 部门类型
	Managers   []Manager `json:"managers"`      // 部门管理员
	ParentDeps []DepInfo `json:"parent_deps"`   // 部门路径
}

type DepInfo struct {
	ID   string `json:"id"`   // 部门ID
	Name string `json:"name"` // 部门名称
	Type string `json:"type"` // 部门类型
}

type Manager struct {
	ID   string `json:"id"`   // 用户ID
	Name string `json:"name"` // 用户名称
}

type App struct {
	ID   string `json:"id"`   // 应用账户ID
	Name string `json:"name"` // 应用账户名
}

// OutboxMsg outbox消息结构体
type OutboxMsg struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// AccountInfo 账户信息
type AccountInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Account     string `json:"account"`
}

// Visitor 请求访问者对象
type Visitor struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	CsfLevel  float64  `json:"csf_level"`
	IP        string   `json:"ip"`
	Mac       string   `json:"mac"`
	Udid      string   `json:"udid"`
	AgentType string   `json:"client_type"`
	Roles     []string `json:"roles"`
	Email     string   `json:"email"`
	Token     string   `json:"token"`
	Type      string   `json:"type"` // 用户代理类型
}

type OSResp struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{}   `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations"`
}

type UserFeature struct {
	Behavior []map[string]interface{} `json:"behavior"`
	// Content  []map[string][]map[string]interface{} `json:"content"`
}
