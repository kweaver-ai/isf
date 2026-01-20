package rclogutils

import (
	"context"
	"os"
	"testing"

	"AuditLog/common/constants/rclogconsts"
	"AuditLog/infra/cmp/langcmp"
)

func TestMain(m *testing.M) {
	langcmp.NewLangCmp().SetSysDefLang(string(langcmp.ZhCN))
	os.Exit(m.Run())
}

// TestGetActiveMetadata 测试获取活跃日志元数据
func TestGetActiveMetadata(t *testing.T) {
	// 重置全局变量，确保测试的独立性
	acReportMetadata = nil

	// 第一次调用，测试元数据的创建
	meta1, err := GetActiveMetadata()
	if err != nil {
		t.Errorf("GetActiveMetadata() error = %v", err)
		return
	}

	// 验证关键字段
	if meta1.DefaultSortField != "date" {
		t.Errorf("Expected DefaultSortField to be 'date', got %s", meta1.DefaultSortField)
	}

	if meta1.DefaultSortDirection != "desc" {
		t.Errorf("Expected DefaultSortDirection to be 'desc', got %s", meta1.DefaultSortDirection)
	}

	if meta1.IdField != rclogconsts.LogID {
		t.Errorf("Expected IdField to be '%s', got %s", rclogconsts.LogID, meta1.IdField)
	}

	// 验证字段数量
	expectedFieldCount := 10
	if len(meta1.Fields) != expectedFieldCount {
		t.Errorf("Expected %d fields, got %d", expectedFieldCount, len(meta1.Fields))
	}

	// 第二次调用，测试缓存机制
	meta2, err := GetActiveMetadata()
	if err != nil {
		t.Errorf("Second GetActiveMetadata() error = %v", err)
		return
	}

	// 验证返回的是同一个对象
	if meta1 != meta2 {
		t.Error("Expected cached metadata to be returned on second call")
	}
}

// TestGetHistoryMetadata 测试获取历史日志元数据
func TestGetHistoryMetadata(t *testing.T) {
	// 重置全局变量，确保测试的独立性
	historyReportMetadata = nil
	ctx := context.Background()

	// 第一次调用，测试元数据的创建
	meta1, err := GetHistoryMetadata(ctx)
	if err != nil {
		t.Errorf("GetHistoryMetadata() error = %v", err)
		return
	}

	// 验证关键字段
	if meta1.DefaultSortField != "dump_date" {
		t.Errorf("Expected DefaultSortField to be 'dump_date', got %s", meta1.DefaultSortField)
	}

	if meta1.DefaultSortDirection != "desc" {
		t.Errorf("Expected DefaultSortDirection to be 'desc', got %s", meta1.DefaultSortDirection)
	}

	if meta1.IdField != rclogconsts.ID {
		t.Errorf("Expected IdField to be '%s', got %s", rclogconsts.ID, meta1.IdField)
	}

	// 验证字段数量
	expectedFieldCount := 5
	if len(meta1.Fields) != expectedFieldCount {
		t.Errorf("Expected %d fields, got %d", expectedFieldCount, len(meta1.Fields))
	}

	// 第二次调用，测试缓存机制
	meta2, err := GetHistoryMetadata(ctx)
	if err != nil {
		t.Errorf("Second GetHistoryMetadata() error = %v", err)
		return
	}

	// 验证返回的是同一个对象
	if meta1 != meta2 {
		t.Error("Expected cached metadata to be returned on second call")
	}
}

// TestMetadataFieldsContent 测试元数据字段内容
func TestMetadataFieldsContent(t *testing.T) {
	// 重置全局变量
	acReportMetadata = nil
	historyReportMetadata = nil

	// 测试活跃日志元数据字段内容
	activeMeta, _ := GetActiveMetadata()
	for _, field := range activeMeta.Fields {
		// 验证必要字段不为空
		if field.Field == "" {
			t.Error("Active metadata field name should not be empty")
		}

		if field.FieldTitle == "" {
			t.Error("Active metadata field title should not be empty")
		}
	}

	// 测试历史日志元数据字段内容
	ctx := context.Background()

	historyMeta, _ := GetHistoryMetadata(ctx)
	for _, field := range historyMeta.Fields {
		// 验证必要字段不为空
		if field.Field == "" {
			t.Error("History metadata field name should not be empty")
		}

		if field.FieldTitle == "" {
			t.Error("History metadata field title should not be empty")
		}
	}
}
