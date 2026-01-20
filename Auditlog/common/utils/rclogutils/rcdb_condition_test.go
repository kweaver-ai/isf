package rclogutils

import (
	"testing"

	"AuditLog/models/rcvo"
)

// TestSetDefaultSortField 测试设置默认排序字段
func TestSetDefaultSortField(t *testing.T) {
	tests := []struct {
		name     string
		orderBy  rcvo.OrderFields
		ids      []string
		expected rcvo.OrderFields
	}{
		{
			name:    "空排序字段和空ID列表",
			orderBy: rcvo.OrderFields{},
			ids:     []string{},
			expected: rcvo.OrderFields{
				{Field: "date", Direction: "desc"},
			},
		},
		{
			name: "已有排序字段",
			orderBy: rcvo.OrderFields{
				{Field: "level", Direction: "asc"},
			},
			ids: []string{},
			expected: rcvo.OrderFields{
				{Field: "level", Direction: "asc"},
			},
		},
		{
			name:     "有ID列表",
			orderBy:  rcvo.OrderFields{},
			ids:      []string{"1", "2"},
			expected: rcvo.OrderFields{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderBy := tt.orderBy
			SetDefaultSortField(&orderBy, tt.ids)

			if len(orderBy) != len(tt.expected) {
				t.Errorf("期望长度 %v, 实际长度 %v", len(tt.expected), len(orderBy))
			}

			if len(orderBy) > 0 && (orderBy[0].Field != tt.expected[0].Field || orderBy[0].Direction != tt.expected[0].Direction) {
				t.Errorf("期望 %v, 实际 %v", tt.expected, orderBy)
			}
		})
	}
}

// TestAddWhereToSql 测试添加查询条件
func TestAddWhereToSql(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		where    map[string]any
		expected string
		hasError bool
	}{
		{
			name:     "空查询条件",
			sql:      "SELECT * FROM table [where]",
			where:    map[string]any{},
			expected: "SELECT * FROM table ",
			hasError: false,
		},
		{
			name:     "日期范围条件",
			sql:      "SELECT * FROM table [where]",
			where:    map[string]any{"date": []interface{}{float64(1000), float64(2000)}},
			expected: "SELECT * FROM table WHERE (f_date BETWEEN 1000 AND 2000)",
			hasError: false,
		},
		{
			name:     "级别条件",
			sql:      "SELECT * FROM table [where]",
			where:    map[string]any{"level": 1},
			expected: "SELECT * FROM table WHERE f_level = 1",
			hasError: false,
		},
		{
			name:     "模糊匹配条件",
			sql:      "SELECT * FROM table [where]",
			where:    map[string]any{"user": "test"},
			expected: "SELECT * FROM table WHERE f_user LIKE '%test%'",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := tt.sql

			err := AddWhereToSql(&sql, tt.where)
			if (err != nil) != tt.hasError {
				t.Errorf("期望错误 %v, 实际错误 %v", tt.hasError, err != nil)
			}

			if err == nil && sql != tt.expected {
				t.Errorf("期望 %v, 实际 %v", tt.expected, sql)
			}
		})
	}
}

// TestAddInToSql 测试添加枚举条件
func TestAddInToSql(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		enums    map[string][]string
		expected string
	}{
		{
			name:     "空枚举条件",
			sql:      "SELECT * FROM table [in]",
			enums:    map[string][]string{},
			expected: "SELECT * FROM table ",
		},
		{
			name:     "单个枚举条件",
			sql:      "SELECT * FROM table [in]",
			enums:    map[string][]string{"user_id": {"1", "2"}},
			expected: "SELECT * FROM table WHERE f_user_id IN ('1','2')",
		},
		{
			name:     "已有WHERE条件",
			sql:      "SELECT * FROM table WHERE 1=1 [in]",
			enums:    map[string][]string{"user_id": {"1", "2"}},
			expected: "SELECT * FROM table WHERE 1=1 AND f_user_id IN ('1','2')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := tt.sql
			AddInToSql(&sql, tt.enums)

			if sql != tt.expected {
				t.Errorf("期望 %v, 实际 %v", tt.expected, sql)
			}
		})
	}
}

// TestBuildActiveCondition 测试构建活跃日志条件
func TestBuildActiveCondition(t *testing.T) {
	tests := []struct {
		name      string
		condition map[string]any
		orderBy   rcvo.OrderFields
		ids       []string
		inUserIDs []string
		exUserIDs []string
		expected  string
		hasError  bool
	}{
		{
			name:      "基本查询",
			condition: map[string]any{},
			orderBy:   rcvo.OrderFields{{Field: "date", Direction: "desc"}},
			ids:       []string{},
			inUserIDs: []string{},
			exUserIDs: []string{},
			expected:  "   ORDER BY f_date DESC",
			hasError:  false,
		},
		{
			name:      "带ID条件",
			condition: map[string]any{},
			orderBy:   rcvo.OrderFields{{Field: "date", Direction: "desc"}},
			ids:       []string{"1", "2"},
			inUserIDs: []string{},
			exUserIDs: []string{},
			expected:  " WHERE f_log_id IN ('1','2')  ORDER BY f_date DESC",
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, err := BuildActiveCondition(tt.condition, tt.orderBy, tt.ids, tt.inUserIDs, tt.exUserIDs)
			if (err != nil) != tt.hasError {
				t.Errorf("期望错误 %v, 实际错误 %v", tt.hasError, err != nil)
			}

			if err == nil && sql != tt.expected {
				t.Errorf("期望 %v, 实际 %v", tt.expected, sql)
			}
		})
	}
}
