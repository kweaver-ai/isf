package mbcenums

type OpType string

func (ot OpType) String() string {
	return string(ot)
}

func (ot OpType) IsConfigured() bool {
	return ot == OpConfigedBuiltInOpT || ot == OpConfigedFromAppstoreOpT
}

const (
	OpConfigedBuiltInOpT      OpType = "op_configed_built-in"      // 操作配置-内置
	OpConfigedFromAppstoreOpT OpType = "op_configed_from-appstore" // 操作配置-来自appstore
	OtherOpT                  OpType = "other"                     // 其他 非操作配置（非op-config接口返回的）
)
