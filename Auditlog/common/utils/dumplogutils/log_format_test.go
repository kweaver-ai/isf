package dumplogutils

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"AuditLog/common"
	"AuditLog/models"
)

func TestLogInfo2CSVString(t *testing.T) {
	// 准备测试数据
	now := time.Now()
	testLog := &models.LogPO{
		LogID:          "test_log_id",
		Date:           now.UnixMicro(),
		UserName:       "test_user",
		UserPaths:      "/test/path",
		Level:          1,
		OpType:         1,
		IP:             "127.0.0.1",
		MAC:            "00:00:00:00:00:00",
		Msg:            "test message",
		ExMsg:          "test ex message",
		UserAgent:      "test user agent",
		AdditionalInfo: `{"key":"value"}`,
		ObjID:          "test_obj_id",
	}

	t.Run("CSV格式转换", func(t *testing.T) {
		csvStr, err := LogInfo2CSVString(testLog, common.Operation)
		if err != nil {
			t.Errorf("转换CSV失败: %v", err)
		}

		// 验证CSV格式
		fields := strings.Split(csvStr, ",")
		if len(fields) != 14 {
			t.Errorf("期望14个字段，实际获得%d个字段", len(fields))
		}

		// 验证字段是否被正确引用
		for _, field := range fields {
			if !strings.HasPrefix(field, `"`) || !strings.HasSuffix(field, `"`) {
				t.Errorf("字段 %s 未被正确引用", field)
			}
		}
	})

	t.Run("包含特殊字符的CSV转换", func(t *testing.T) {
		testLog.Msg = `test,message"with,special"chars`

		csvStr, err := LogInfo2CSVString(testLog, common.Operation)
		if err != nil {
			t.Errorf("转换包含特殊字符的CSV失败: %v", err)
		}

		if !strings.Contains(csvStr, `""`) {
			t.Error("特殊字符未被正确转义")
		}
	})
}

func TestLogInfo2XMLString(t *testing.T) {
	// 准备测试数据
	now := time.Now()
	testLog := &models.LogPO{
		LogID:          "test_log_id",
		Date:           now.UnixMicro(),
		UserName:       "test_user",
		UserPaths:      "/test/path",
		Level:          1,
		OpType:         1,
		IP:             "127.0.0.1",
		MAC:            "00:00:00:00:00:00",
		Msg:            "test message",
		ExMsg:          "test ex message",
		UserAgent:      "test user agent",
		AdditionalInfo: `{"key":"value"}`,
		ObjID:          "test_obj_id",
	}

	t.Run("XML格式转换", func(t *testing.T) {
		xmlStr, err := LogInfo2XMLString(testLog, common.Operation)
		if err != nil {
			t.Errorf("转换XML失败: %v", err)
		}

		// 验证XML格式
		if !strings.HasPrefix(xmlStr, "<log-id") {
			t.Error("XML格式不正确，缺少开始标签")
		}

		if !strings.HasSuffix(xmlStr, "</log-id>") {
			t.Error("XML格式不正确，缺少结束标签")
		}
	})

	t.Run("包含特殊字符的XML转换", func(t *testing.T) {
		testLog.Msg = `test<message>with&special"chars`

		xmlStr, err := LogInfo2XMLString(testLog, common.Operation)
		if err != nil {
			t.Errorf("转换包含特殊字符的XML失败: %v", err)
		}

		if !strings.Contains(xmlStr, "&lt;") || !strings.Contains(xmlStr, "&gt;") {
			t.Error("XML特殊字符未被正确转义")
		}
	})
}

func TestFormatAdditionalInfo(t *testing.T) {
	t.Run("正常JSON格式", func(t *testing.T) {
		additionalInfo := `{"key":"value"}`
		objID := "test_obj_id"

		info, err := formatAdditionalInfo(additionalInfo, objID)
		if err != nil {
			t.Errorf("格式化additionalInfo失败: %v", err)
		}

		// 解析结果验证
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(info), &result); err != nil {
			t.Errorf("解析结果失败: %v", err)
		}

		if result["obj_id"] != objID {
			t.Errorf("期望obj_id为%s，实际为%v", objID, result["obj_id"])
		}
	})

	t.Run("空JSON处理", func(t *testing.T) {
		info, err := formatAdditionalInfo("", "test_obj_id")
		if err != nil {
			t.Errorf("处理空JSON失败: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(info), &result); err != nil {
			t.Errorf("解析结果失败: %v", err)
		}
	})

	t.Run("无效JSON格式", func(t *testing.T) {
		_, err := formatAdditionalInfo("invalid json", "test_obj_id")
		if err == nil {
			t.Error("期望无效JSON返回错误，但没有")
		}
	})
}

func TestEscapeCSVField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"普通文本", "normal text", "normal text"},
		{"包含逗号", "text,with,comma", "text,with,comma"},
		{"包含引号", `text"with"quotes`, `text""with""quotes`},
		{"包含换行", "text\nwith\nbreaks", "text\nwith\nbreaks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeCSVField(tt.input)
			if tt.name == "包含引号" && !strings.Contains(result, `""`) {
				t.Errorf("引号未被正确转义: 期望包含 \"\"，实际为 %s", result)
			}
		})
	}
}
