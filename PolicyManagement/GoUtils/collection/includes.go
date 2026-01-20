package collection

// IncludesString 查找切片是否包含特定字符串
// list 可以是包含任意类型元素的切片
// elem 是要检查的元素
func IncludesString(colleciton []string, value string) bool {
	if len(colleciton) == 0 {
		return false
	}

	for _, v := range colleciton {
		if v == value {
			return true
		}
	}

	return false
}
