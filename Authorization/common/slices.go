package common

import "slices"

// Distinct 自适应去重函数，考虑切片长度为0或1的特殊情况
//
//nolint:mnd
func Distinct[T comparable](slice []T) []T {
	// 如果切片长度为0或1，直接返回原切片
	if len(slice) < 2 {
		return slice
	}

	// 对于小切片，使用简单循环去重
	if len(slice) < 10 { // 假设切片长度小于10时为小切片
		return distinctSmallSlice(slice)
	}
	// 对于大切片，使用map去重
	return distinctLargeSlice(slice)
}

// distinctSmallSlice 适用于小切片的去重
func distinctSmallSlice[T comparable](slice []T) []T {
	var result []T
	for _, v := range slice {
		if !slices.Contains(result, v) {
			result = append(result, v)
		}
	}
	return result
}

// distinctLargeSlice 适用于大切片的去重，使用map优化
func distinctLargeSlice[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
