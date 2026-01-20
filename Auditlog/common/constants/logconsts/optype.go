package logconsts

var (
	LogLevel *logLevel
	OpType   *opType
)

// 日志级别常量枚举
type logLevel struct {
	INFO int
	WARN int
}

// 操作类型常量枚举集合
type opType struct {
	ManagementType *mgntType
}

// 管理日志操作类型常量枚举
type mgntType struct {
	CREATE int
	ADD    int
	SET    int
	DELETE int
	EXPORT int
	EDIT   int
	OTHER  int
}

func init() {
	// 初始化日志级别枚举
	LogLevel = &logLevel{
		INFO: 1,
		WARN: 2,
	}
	// 初始化操作类型枚举
	OpType = &opType{
		ManagementType: &mgntType{
			CREATE: 1,
			ADD:    2,
			SET:    3,
			DELETE: 4,
			EXPORT: 9,
			EDIT:   20,
			OTHER:  127,
		},
	}
}
