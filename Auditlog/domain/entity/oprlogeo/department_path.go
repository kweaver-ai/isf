package oprlogeo

import (
	"strings"

	"AuditLog/common/utils"
	"AuditLog/infra/cmp/umcmp"
)

// DepartmentPath 部门路径信息
type DepartmentPath struct {
	IDPath   string `json:"id_path"`   // 用户所属部门id全路径，示例："4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0"
	NamePath string `json:"name_path"` // 用户所属部门名称全路径，示例："爱数/数据智能产品BG/AnyShare研发线/智能搜索研发部"
}

func ParseDepartments(depts [][]umcmp.ObjectBaseInfo) (deptInfos []*DepartmentPath) {
	deptInfos = make([]*DepartmentPath, 0, len(depts))

	for _, dept := range depts {
		deptAllIDs := make([]string, 0, len(depts))
		deptAllNames := make([]string, 0, len(depts))

		for _, item := range dept {
			deptAllIDs = append(deptAllIDs, item.ID)
			deptAllNames = append(deptAllNames, item.Name)
		}

		deptInfo := &DepartmentPath{
			IDPath:   strings.Join(deptAllIDs, "/"),
			NamePath: strings.Join(deptAllNames, "/"),
		}

		deptInfos = append(deptInfos, deptInfo)
	}

	return
}

// GetDepartmentIDsByLevel 根据deptInfos []*DepartmentPath获取指定层级的部门id
// level表示要返回的从底层往上的层数，例如level=2表示返回每个部门路径的最后两层部门id
// Example:
//
//	Input:
//	  deptInfos = []*DepartmentPath{
//	    {
//	      IDPath: "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
//	      NamePath: "爱数/数据智能产品BG/AnyShare研发线",
//	    },
//	    {
//	      IDPath: "4e8bfbda-d99c-11eb-35b9-24e8e0506805/abcdae8b-d9c9-1eb1-5b39-506802412345",
//	      NamePath: "爱数/运营管理部",
//	    },
//	  }
//	  level = 2
//	Output:
//	  []string{"4bfdae8b-d9c9-1eb1-5b39-5068024e8e05", "e8bfbda4-d31c-12ab-34c9-50680524e8e0", "4e8bfbda-d99c-11eb-35b9-24e8e0506805", "abcdae8b-d9c9-1eb1-5b39-506802412345"}
func GetDepartmentIDsByLevel(deptInfos []*DepartmentPath, level int) []string {
	result := make([]string, 0, len(deptInfos)*level)

	for _, dept := range deptInfos {
		if dept == nil || dept.IDPath == "" {
			continue
		}

		ids := strings.Split(dept.IDPath, "/")

		if level <= 0 {
			continue
		}

		// 获取每个部门路径的最后level层ID
		start := len(ids)
		if level < start {
			start = level
		}

		for i := 0; i < start; i++ {
			result = append(result, ids[len(ids)-start+i])
		}
	}

	result = utils.DeduplGeneric(result)

	return result
}
