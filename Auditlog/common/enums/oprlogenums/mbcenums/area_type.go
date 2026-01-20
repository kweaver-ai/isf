package mbcenums

type AreaType int

const (
	AreaTypeTopFixed      AreaType = 1 // 1: 顶部栏 固定展示区
	AreaTypeTopFold       AreaType = 2 // 2: 顶部栏 折叠区
	AreaTypeTopFoldMore   AreaType = 3 // 3: 顶部栏 折叠区 “更多”
	AreaTypeTopFoldIntegr AreaType = 4 // 4: 顶部栏 折叠区 “集成”

	AreaTypeRightFixed  AreaType = 5 // 5: 右键菜单 固定展示区
	AreaTypeRightMore   AreaType = 6 // 6: 右键菜单 “更多”
	AreaTypeRightIntegr AreaType = 7 // 7: 右键菜单 “集成”

)
