package persconsts

import "time"

const (
	PersSvcName = "personalization"

	// RedisKeyPrefix personalization 缓存redis key 前缀
	RedisKeyPrefix = "personalization"

	UnifiedRecRedisTTL = time.Second * 5 // 缓存时间

	UnifiedRecRedisOpTimeout = time.Millisecond * 100 // 缓存操作超时时间
)

// 特征类型
const (
	FeatureTypeStatic  string = "static"
	FeatureTypeDynamic string = "dynamic"
)

// 适用范围
const (
	ApplicableTypeUser  string = "user"
	ApplicableTypeDept  string = "dept"
	ApplicableTypeOther string = "other"
)

// 特征key
const (
	UserStaticFeatureKey    string = "user_static"
	DeptStaticFeatureKey    string = "dept_static"
	UserDocCenterFeatureKey string = "user_doc_center_behavior_ratio"
	UserKcFeatureKey        string = "user_kc_behavior_ratio"
	DeptDocCenterFeatureKey string = "dept_doc_center_behavior_ratio"
	DeptKcFeatureKey        string = "dept_kc_behavior_ratio"
)
