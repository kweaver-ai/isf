package common

// Number 定义一个泛型类型约束，支持更多的数值类型
type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// Min 求最小值
func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}
