package persmodels

type PersFeatureConfigReq struct {
	FeatureType    string `json:"feature_type"`
	ApplicableType string `json:"applicable_type"`
	Keyword        string `json:"keyword"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
}

type PersFeatureConfig struct {
	ID             int64  `json:"id"`
	Key            string `json:"key"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	IsBuiltIn      int    `json:"is_built_in"`
	ApplicableType string `json:"applicable_type"`
	Structure      string `json:"structure"`
}

type PersFeatureConfigResp struct {
	TotalCount int64                `json:"total_count"`
	Entries    []*PersFeatureConfig `json:"entries"`
}
