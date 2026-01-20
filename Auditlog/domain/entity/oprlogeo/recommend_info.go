package oprlogeo

// RecommendInfo 推荐相关信息
type RecommendInfo struct {
	NotUseForRec bool   `json:"not_use_for_rec"` // 是否不用于推荐，默认：false
	ExtInfoJSON  string `json:"ext_info_json"`   // 推荐用扩展信息，示例：{"k1":{"k2::"v2"}}
}
