package utils

// Difference 计算差集(a - b)
// 返回差集，即在a中存在，但在b中不存在的元素
func Difference(a, b []string) (diff []string) {
	m := make(map[string]struct{}, len(b))
	for _, v := range b {
		m[v] = struct{}{}
	}

	diff = make([]string, 0)

	for _, v := range a {
		if _, ok := m[v]; !ok {
			diff = append(diff, v)
		}
	}

	return diff
}

// DifferenceGeneric 计算差集(a - b) 泛型版本
// 返回差集，即在a中存在，但在b中不存在的元素
func DifferenceGeneric[T comparable](a, b []T) (diff []T) {
	m := make(map[T]struct{}, len(b))
	for _, v := range b {
		m[v] = struct{}{}
	}

	diff = make([]T, 0)

	for _, v := range a {
		if _, ok := m[v]; !ok {
			diff = append(diff, v)
		}
	}

	return diff
}

// Intersection 计算交集
// 返回交集，即在a中存在，且在b中也存在的元素
func Intersection(a, b []string) (intersection []string) {
	m := make(map[string]struct{}, len(b))
	for _, v := range b {
		m[v] = struct{}{}
	}

	intersection = make([]string, 0)

	for _, v := range a {
		if _, ok := m[v]; ok {
			intersection = append(intersection, v)
		}
	}

	return intersection
}

func Exists(a []string, v string) bool {
	for _, item := range a {
		if item == v {
			return true
		}
	}

	return false
}

// ExistsGeneric 切片中是否存在某个元素 泛型版本
// 切片a中是否存在元素v
// 时间复杂度O(n)
func ExistsGeneric[T comparable](a []T, v T) bool {
	for i := range a {
		if a[i] == v {
			return true
		}
	}

	return false
}

// DeduplGeneric 泛型版本 切片去重
func DeduplGeneric[T comparable](a []T) (newSlice []T) {
	newSlice = make([]T, 0, len(a))

	if len(a) == 0 {
		return
	}

	m := make(map[T]struct{}, len(a))
	for _, v := range a {
		m[v] = struct{}{}
	}

	// 使用这种方法，可以保证去重后的切片顺序和原切片一致
	for _, v := range a {
		if _, ok := m[v]; ok {
			newSlice = append(newSlice, v)
			delete(m, v)
		}
	}

	return newSlice
}

// SliceToPtrSlice 切片转换为指针切片（使用原元素地址）
func SliceToPtrSlice[T comparable](a []T) []*T {
	b := make([]*T, len(a))
	for i := range a {
		b[i] = &a[i]
	}

	return b
}

// IsDuplicationGeneric 判断切片是否有重复元素 泛型版本
func IsDuplicationGeneric[T comparable](a []T) bool {
	m := make(map[T]struct{}, len(a))
	for _, v := range a {
		if _, ok := m[v]; ok {
			return true
		}

		m[v] = struct{}{}
	}

	return false
}
