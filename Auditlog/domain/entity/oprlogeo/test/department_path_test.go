package oprlogeotest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"AuditLog/domain/entity/oprlogeo"
)

func TestGetDepartmentIDsByLevel(t *testing.T) {
	tests := []struct {
		name      string
		deptInfos []*oprlogeo.DepartmentPath
		level     int
		want      []string
	}{
		{
			name: "正常情况-返回最后两层",
			deptInfos: []*oprlogeo.DepartmentPath{
				{
					IDPath:   "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
					NamePath: "爱数/数据智能产品BG/AnyShare研发线",
				},
				{
					IDPath:   "4e8bfbda-d99c-11eb-35b9-24e8e0506805/abcdae8b-d9c9-1eb1-5b39-506802412345",
					NamePath: "爱数/运营管理部",
				},
			},
			level: 2,
			want: []string{
				"4bfdae8b-d9c9-1eb1-5b39-5068024e8e05",
				"e8bfbda4-d31c-12ab-34c9-50680524e8e0",
				"4e8bfbda-d99c-11eb-35b9-24e8e0506805",
				"abcdae8b-d9c9-1eb1-5b39-506802412345",
			},
		},
		{
			name: "level大于部门层级数",
			deptInfos: []*oprlogeo.DepartmentPath{
				{
					IDPath:   "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05",
					NamePath: "爱数/数据智能产品BG",
				},
			},
			level: 3,
			want: []string{
				"4e8bfbda-d99c-11eb-35b9-24e8e0506805",
				"4bfdae8b-d9c9-1eb1-5b39-5068024e8e05",
			},
		},
		{
			name: "level=1只返回最后一层",
			deptInfos: []*oprlogeo.DepartmentPath{
				{
					IDPath:   "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
					NamePath: "爱数/数据智能产品BG/AnyShare研发线",
				},
			},
			level: 1,
			want: []string{
				"e8bfbda4-d31c-12ab-34c9-50680524e8e0",
			},
		},
		{
			name: "空部门路径",
			deptInfos: []*oprlogeo.DepartmentPath{
				{
					IDPath:   "",
					NamePath: "",
				},
			},
			level: 2,
			want:  []string{},
		},
		{
			name:      "nil部门信息",
			deptInfos: []*oprlogeo.DepartmentPath{nil},
			level:     2,
			want:      []string{},
		},
		{
			name:      "空部门信息切片",
			deptInfos: []*oprlogeo.DepartmentPath{},
			level:     2,
			want:      []string{},
		},
		{
			name: "level<=0",
			deptInfos: []*oprlogeo.DepartmentPath{
				{
					IDPath:   "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05",
					NamePath: "爱数/数据智能产品BG",
				},
			},
			level: 0,
			want:  []string{},
		},
		{
			name: "包含重复ID",
			deptInfos: []*oprlogeo.DepartmentPath{
				{
					IDPath:   "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
					NamePath: "爱数/数据智能产品BG/AnyShare研发线",
				},
				{
					IDPath:   "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05",
					NamePath: "爱数/数据智能产品BG",
				},
			},
			level: 2,
			want: []string{
				"4e8bfbda-d99c-11eb-35b9-24e8e0506805",
				"4bfdae8b-d9c9-1eb1-5b39-5068024e8e05",
				"e8bfbda4-d31c-12ab-34c9-50680524e8e0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := oprlogeo.GetDepartmentIDsByLevel(tt.deptInfos, tt.level)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
