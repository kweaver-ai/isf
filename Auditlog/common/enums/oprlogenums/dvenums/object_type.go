package dvenums

type ObjectType string

func (ot ObjectType) String() string {
	return string(ot)
}

const (
	FolderOT      ObjectType = "folder"       // 文件夹|文档库
	FavCategoryOT ObjectType = "fav_category" // 收藏夹
	DocLibTypeOT  ObjectType = "doc_lib_type" // 库类型
)

func IsDocLibOT(ot ObjectType) bool {
	return ot == FolderOT
}

func IsDocLibOTString(ot string) bool {
	return IsDocLibOT(ObjectType(ot))
}
