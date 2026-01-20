package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 单元测试用例
func TestDifference(t *testing.T) {
	t.Parallel()

	// 测试差集为空的情况
	a := []string{"a", "b", "c"}
	b := []string{}
	expected := []string{"a", "b", "c"}
	result := Difference(a, b)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Difference failed, expected %v but got %v", expected, result)
	}

	// 测试差集为非空的情况
	a = []string{"a", "b", "c"}
	b = []string{"b", "c"}
	expected = []string{"a"}
	result = Difference(a, b)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Difference failed, expected %v but got %v", expected, result)
	}
}

// 单元测试用例
func TestIntersection(t *testing.T) {
	t.Parallel()

	// 测试交集为空的情况
	a := []string{"a", "b", "c"}
	b := []string{}
	expected := []string{}
	result := Intersection(a, b)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersection failed, expected %v but got %v", expected, result)
	}

	// 测试交集为非空的情况
	a = []string{"a", "b", "c"}
	b = []string{"b", "c"}
	expected = []string{"b", "c"}
	result = Intersection(a, b)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersection failed, expected %v but got %v", expected, result)
	}

	//	测试a为空的情况
	a = []string{}
	b = []string{"b", "c"}
	expected = []string{}
	result = Difference(a, b)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Difference failed, expected %v but got %v", expected, result)
	}
}

// 单元测试用例
func TestExists(t *testing.T) {
	t.Parallel()

	// 测试元素存在的情况
	a := []string{"a", "b", "c"}
	v := "b"
	expected := true
	result := Exists(a, v)

	if result != expected {
		t.Errorf("Exists failed, expected %v but got %v", expected, result)
	}

	// 测试元素不存在的情况
	v = "d"
	expected = false
	result = Exists(a, v)

	if result != expected {
		t.Errorf("Exists failed, expected %v but got %v", expected, result)
	}
}

// IsDuplicationGeneric ut
func TestIsDuplicationGeneric(t *testing.T) {
	t.Parallel()

	// 1. 测试有重复元素的情况
	a := []int{1, 2, 3, 1, 2, 3}
	result := IsDuplicationGeneric(a)

	assert.Equal(t, result, true)

	// 2. 测试无重复元素的情况
	a = []int{1, 2, 3}
	result = IsDuplicationGeneric(a)

	assert.Equal(t, result, false)

	// 3. 测试string类型
	a1 := []string{"a", "b", "c", "a", "b", "c"}
	result1 := IsDuplicationGeneric(a1)

	assert.Equal(t, result1, true)

	// 4. 测试int64类型
	a2 := []int64{1, 2, 3, 1, 2, 3}
	result2 := IsDuplicationGeneric(a2)

	assert.Equal(t, result2, true)

	// 5. 测试float64类型
	a3 := []float64{1.1, 2.2, 3.3, 1.1, 2.2, 3.3}
	result3 := IsDuplicationGeneric(a3)

	assert.Equal(t, result3, true)

	// 6. 测试bool类型
	a4 := []bool{true, false, true, false}
	result4 := IsDuplicationGeneric(a4)

	assert.Equal(t, result4, true)

	// 7. 测试空切片
	var a5 []string
	result5 := IsDuplicationGeneric(a5)

	assert.Equal(t, result5, false)

	// 8. 测试nil切片
	var a6 []string

	result6 := IsDuplicationGeneric(a6)

	assert.Equal(t, result6, false)

	// 9. 测试结构体切片
	type Person struct {
		Name string
		Age  int
	}

	a7 := []Person{{"a", 1}, {"b", 2}, {"c", 3}, {"a", 1}, {"b", 2}, {"c", 3}}
	result7 := IsDuplicationGeneric(a7)

	assert.Equal(t, result7, true)
}

// TestDeduplGeneric 测试 DeduplGeneric 函数
func TestDeduplGeneric(t *testing.T) {
	t.Parallel()

	// 测试空切片
	var a1 []int

	expected1 := []int{}
	result1 := DeduplGeneric(a1)

	assert.Equal(t, expected1, result1, "DeduplGeneric failed on empty slice")

	// 测试无重复元素
	a2 := []int{1, 2, 3}
	expected2 := []int{1, 2, 3}
	result2 := DeduplGeneric(a2)

	assert.Equal(t, expected2, result2, "DeduplGeneric failed on no duplicates")

	// 测试所有元素重复
	a3 := []int{1, 1, 1}
	expected3 := []int{1}
	result3 := DeduplGeneric(a3)

	assert.Equal(t, expected3, result3, "DeduplGeneric failed on all elements are the same")

	// 测试部分重复元素
	a4 := []int{1, 2, 2, 3, 1}
	expected4 := []int{1, 2, 3}
	result4 := DeduplGeneric(a4)

	assert.Equal(t, expected4, result4, "DeduplGeneric failed on some duplicates")
}

// TestDifferenceGeneric 测试 DifferenceGeneric 函数
func TestDifferenceGeneric(t *testing.T) {
	t.Parallel()

	a := []int{1, 2, 3, 4}
	b := []int{3, 4}

	expected := []int{1, 2}
	result := DifferenceGeneric(a, b)

	assert.Equal(t, expected, result, "DifferenceGeneric failed")
}

// TestExistsGeneric 测试 ExistsGeneric 函数
func TestExistsGeneric(t *testing.T) {
	t.Parallel()

	// 测试元素存在的情况
	a := []int{1, 2, 3}
	v := 2
	expected := true
	result := ExistsGeneric(a, v)

	assert.Equal(t, expected, result, "ExistsGeneric failed for existing element")

	// 测试元素不存在的情况
	v = 4
	expected = false
	result = ExistsGeneric(a, v)

	assert.Equal(t, expected, result, "ExistsGeneric failed for non-existing element")
}

// TestSliceToPtrSlice 测试 SliceToPtrSlice 函数
func TestSliceToPtrSlice(t *testing.T) {
	t.Parallel()

	a := []int{1, 2, 3}
	expected := []*int{&a[0], &a[1], &a[2]}
	result := SliceToPtrSlice(a)

	assert.Equal(t, expected, result, "SliceToPtrSlice failed")
}
