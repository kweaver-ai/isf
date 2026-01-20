package dumplogutils

import (
	"bytes"
	"os"
	"testing"

	"github.com/yeka/zip"

	"AuditLog/infra/cmp/langcmp"
)

func TestMain(m *testing.M) {
	langcmp.NewLangCmp().SetSysDefLang(string(langcmp.ZhCN))
	os.Exit(m.Run())
}

func TestGenZipFile(t *testing.T) {
	// 测试数据准备
	fileContent := []byte("这是测试内容")
	fileName := "test.txt"

	// 测试场景1：不带密码的zip文件生成
	t.Run("无密码压缩", func(t *testing.T) {
		zipContent, err := GenZipFile(fileContent, fileName, "")
		if err != nil {
			t.Errorf("生成无密码zip文件失败: %v", err)
		}

		// 验证生成的zip文件
		reader, err := zip.NewReader(bytes.NewReader(zipContent), int64(len(zipContent)))
		if err != nil {
			t.Errorf("读取zip文件失败: %v", err)
		}

		if len(reader.File) != 1 {
			t.Errorf("期望1个文件，实际获得%d个文件", len(reader.File))
		}

		if reader.File[0].Name != fileName {
			t.Errorf("期望文件名%s，实际获得%s", fileName, reader.File[0].Name)
		}
	})

	// 测试场景2：带密码的zip文件生成
	t.Run("带密码压缩", func(t *testing.T) {
		password := "123456"

		zipContent, err := GenZipFile(fileContent, fileName, password)
		if err != nil {
			t.Errorf("生成带密码zip文件失败: %v", err)
		}

		// 验证生成的zip文件
		reader, err := zip.NewReader(bytes.NewReader(zipContent), int64(len(zipContent)))
		if err != nil {
			t.Errorf("读取zip文件失败: %v", err)
		}

		if !reader.File[0].IsEncrypted() {
			t.Error("期望文件被加密，但实际未加密")
		}
	})
}

func TestSplitFile(t *testing.T) {
	// 测试数据准备
	fileContent := []byte("abcdefghijklmnopqrstuvwxyz")

	// 测试场景1：正常分块
	t.Run("正常分块", func(t *testing.T) {
		size := int64(5)

		parts, err := SplitFile(fileContent, size)
		if err != nil {
			t.Errorf("文件分块失败: %v", err)
		}

		expectedParts := 6 // 26字节分成5字节一块，应该有6块
		if len(parts) != expectedParts {
			t.Errorf("期望%d个分块，实际获得%d个分块", expectedParts, len(parts))
		}

		// 验证最后一个分块的大小
		lastPartSize := len(parts[len(parts)-1])
		if lastPartSize != 1 {
			t.Errorf("期望最后一块大小为1，实际为%d", lastPartSize)
		}
	})

	// 测试场景2：分块大小大于文件大小
	t.Run("分块大小大于文件", func(t *testing.T) {
		size := int64(100)

		parts, err := SplitFile(fileContent, size)
		if err != nil {
			t.Errorf("文件分块失败: %v", err)
		}

		if len(parts) != 1 {
			t.Errorf("期望1个分块，实际获得%d个分块", len(parts))
		}

		if len(parts[0]) != len(fileContent) {
			t.Errorf("期望分块大小为%d，实际为%d", len(fileContent), len(parts[0]))
		}
	})

	// 测试场景3：空文件
	t.Run("空文件", func(t *testing.T) {
		emptyContent := []byte{}
		size := int64(5)

		parts, err := SplitFile(emptyContent, size)
		if err != nil {
			t.Errorf("空文件分块失败: %v", err)
		}

		if len(parts) != 0 {
			t.Errorf("期望0个分块，实际获得%d个分块", len(parts))
		}
	})
}
