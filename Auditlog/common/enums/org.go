package enums

type OrgObjType string

// 类型 1-部门（或组织） 2-用户 3-用户组类型 1-部门（或组织） 2-用户 3-用户组
const (
	OrgObjTypeDep   OrgObjType = "dept"       // 部门（或组织）
	OrgObjTypeUser  OrgObjType = "user"       // 用户
	OrgObjTypeGroup OrgObjType = "user_group" // 用户组
)

func (e OrgObjType) String() string {
	switch e {
	case OrgObjTypeUser:
		return "user"
	case OrgObjTypeDep:
		return "dept"
	case OrgObjTypeGroup:
		return "user_group"
	}
	return ""
}

func Contains(ts []string, t OrgObjType) bool {
	for _, v := range ts {
		if v == t.String() {
			return true
		}
	}
	return false
}
