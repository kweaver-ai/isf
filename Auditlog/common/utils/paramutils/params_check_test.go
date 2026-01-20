package paramutils

import (
	"context"
	"os"
	"testing"

	"AuditLog/common"
	"AuditLog/infra/cmp/langcmp"
	"AuditLog/models/rcvo"
)

func TestMain(m *testing.M) {
	langcmp.NewLangCmp().SetSysDefLang(string(langcmp.ZhCN))
	os.Exit(m.Run())
}

func TestCategoryCheck(t *testing.T) {
	ctx := context.Background()

	// 测试有效的category
	if err := CategoryCheck(ctx, common.AllLogType[0]); err != nil {
		t.Errorf("Expected nil error for valid category, but got %v", err)
	}

	// 测试无效的category
	if err := CategoryCheck(ctx, "invalid_category"); err == nil {
		t.Error("Expected error for invalid category, but got nil")
	}
}

func TestLimitCheck(t *testing.T) {
	ctx := context.Background()
	defaultLimit := 100

	// 测试有效的limit
	if err := LimitCheck(ctx, 50, defaultLimit); err != nil {
		t.Errorf("Expected nil error for valid limit, but got %v", err)
	}

	// 测试limit <= 0
	if err := LimitCheck(ctx, 0, defaultLimit); err == nil {
		t.Error("Expected error for limit <= 0, but got nil")
	}

	// 测试limit > defaultLimit
	if err := LimitCheck(ctx, defaultLimit+1, defaultLimit); err == nil {
		t.Error("Expected error for limit > defaultLimit, but got nil")
	}
}

func TestParamsCheck(t *testing.T) {
	ctx := context.Background()
	candidates := []string{"field1", "field2", "field3"}
	tag := "test"

	// 测试有效的参数
	validEntry := []string{"field1", "field2"}
	if err := ParamsCheck(ctx, validEntry, candidates, tag); err != nil {
		t.Errorf("Expected nil error for valid entry, but got %v", err)
	}

	// 测试空字符串参数
	emptyEntry := []string{"field1", ""}
	if err := ParamsCheck(ctx, emptyEntry, candidates, tag); err != nil {
		t.Errorf("Expected nil error for entry with empty string, but got %v", err)
	}

	// 测试无效的参数
	invalidEntry := []string{"field1", "invalid_field"}
	if err := ParamsCheck(ctx, invalidEntry, candidates, tag); err == nil {
		t.Error("Expected error for invalid entry, but got nil")
	}
}

func TestGetAvaliableParams(t *testing.T) {
	// 模拟getParams函数
	mockGetParams := func() (*rcvo.ReportMetadataRes, error) {
		return &rcvo.ReportMetadataRes{
			Fields: []rcvo.ReportField{
				{
					Field:       "field1",
					IsCanSearch: 1,
					IsCanSort:   1,
					SearchFieldConfig: rcvo.ReportSearchFieldConfig{
						IsCanSearchByApi: true,
					},
				},
				{
					Field:       "field2",
					IsCanSearch: 0,
					IsCanSort:   1,
				},
				{
					Field:       "field3",
					IsCanSearch: 1,
					IsCanSort:   0,
					SearchFieldConfig: rcvo.ReportSearchFieldConfig{
						IsCanSearchByApi: false,
					},
				},
			},
		}, nil
	}

	params := GetAvaliableParams(mockGetParams)

	// 验证DataFields
	expectedDataFields := []string{"field1", "field2", "field3"}
	if !compareSlices(params.DataFields, expectedDataFields) {
		t.Errorf("Expected DataFields %v, but got %v", expectedDataFields, params.DataFields)
	}

	// 验证SearchFields
	expectedSearchFields := []string{"field1", "field3"}
	if !compareSlices(params.SearchFields, expectedSearchFields) {
		t.Errorf("Expected SearchFields %v, but got %v", expectedSearchFields, params.SearchFields)
	}

	// 验证KeyWordFields
	expectedKeyWordFields := []string{"field1"}
	if !compareSlices(params.KeyWordFields, expectedKeyWordFields) {
		t.Errorf("Expected KeyWordFields %v, but got %v", expectedKeyWordFields, params.KeyWordFields)
	}

	// 验证OrderFields
	expectedOrderFields := []string{"field1", "field2"}
	if !compareSlices(params.OrderFields, expectedOrderFields) {
		t.Errorf("Expected OrderFields %v, but got %v", expectedOrderFields, params.OrderFields)
	}
}

// 辅助函数：比较两个字符串切片是否相等
func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
