package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayRemoveDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "空数组",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "无重复元素",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "有重复元素",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "c", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ArrayRemoveDuplicate(tt.input)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestGetDocLibIDByDocID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "正常文档ID",
			input:    "gns://domain/folder/doc",
			expected: "gns://domain",
		},
		{
			name:     "无效格式",
			input:    "invalid-format",
			expected: "",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDocLibIDByDocID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInArray(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		list     []string
		expected bool
	}{
		{
			name:     "存在的元素",
			value:    "apple",
			list:     []string{"banana", "apple", "orange"},
			expected: true,
		},
		{
			name:     "不存在的元素",
			value:    "grape",
			list:     []string{"banana", "apple", "orange"},
			expected: false,
		},
		{
			name:     "空列表",
			value:    "apple",
			list:     []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InArray(tt.value, tt.list)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		precision int
		expected  string
	}{
		{
			name:      "0字节",
			size:      0,
			precision: 2,
			expected:  "0B",
		},
		{
			name:      "字节级别",
			size:      500,
			precision: 2,
			expected:  "500B",
		},
		{
			name:      "KB级别",
			size:      1024 * 2,
			precision: 2,
			expected:  "2.00KB",
		},
		{
			name:      "MB级别",
			size:      1024 * 1024 * 3,
			precision: 1,
			expected:  "3.0MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFileSize(tt.size, tt.precision)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	assert.Len(t, id1, 32)
	assert.Len(t, id2, 32)
	assert.NotEqual(t, id1, id2)
	assert.Regexp(t, "^[A-F0-9]{32}$", id1)
	assert.Regexp(t, "^[A-F0-9]{32}$", id2)
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int64
		b        int64
		expected int64
	}{
		{
			name:     "a小于b",
			a:        1,
			b:        2,
			expected: 1,
		},
		{
			name:     "a大于b",
			a:        10,
			b:        5,
			expected: 5,
		},
		{
			name:     "相等",
			a:        7,
			b:        7,
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}
