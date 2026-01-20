package helpers

import (
	"os"
	"testing"
)

func TestIsLocalDev(t *testing.T) {
	// 测试用例1: 环境变量为 true
	t.Run("env var is true", func(t *testing.T) {
		os.Setenv(EnvIsLocalDev, "true")
		defer os.Unsetenv(EnvIsLocalDev)

		if !IsLocalDev() {
			t.Error("IsLocalDev() should return true when env var is 'true'")
		}
	})

	// 测试用例2: 环境变量为 false
	t.Run("env var is false", func(t *testing.T) {
		os.Setenv(EnvIsLocalDev, "false")
		defer os.Unsetenv(EnvIsLocalDev)

		if IsLocalDev() {
			t.Error("IsLocalDev() should return false when env var is 'false'")
		}
	})

	// 测试用例3: 测试 mock 功能
	t.Run("mock is true", func(t *testing.T) {
		mockIsLocalDev = true
		defer func() { mockIsLocalDev = false }()

		if !IsLocalDev() {
			t.Error("IsLocalDev() should return true when mockIsLocalDev is true")
		}
	})
}

func TestIsDebugMode(t *testing.T) {
	// 测试用例1: debug 模式开启
	t.Run("debug mode is true", func(t *testing.T) {
		os.Setenv(isDebugMode, "true")
		defer os.Unsetenv(isDebugMode)

		if !IsDebugMode() {
			t.Error("IsDebugMode() should return true when env var is 'true'")
		}
	})

	// 测试用例2: debug 模式关闭
	t.Run("debug mode is false", func(t *testing.T) {
		os.Setenv(isDebugMode, "false")
		defer os.Unsetenv(isDebugMode)

		if IsDebugMode() {
			t.Error("IsDebugMode() should return false when env var is 'false'")
		}
	})
}

func TestIsOprLogShowLogForDebug(t *testing.T) {
	// 测试用例1: debug 模式开启
	t.Run("debug mode is true", func(t *testing.T) {
		os.Setenv(isDebugMode, "true")
		defer os.Unsetenv(isDebugMode)

		if !IsOprLogShowLogForDebug() {
			t.Error("IsOprLogShowLogForDebug() should return true when debug mode is true")
		}
	})

	// 测试用例2: debug 模式关闭
	t.Run("debug mode is false", func(t *testing.T) {
		os.Setenv(isDebugMode, "false")
		defer os.Unsetenv(isDebugMode)

		if IsOprLogShowLogForDebug() {
			t.Error("IsOprLogShowLogForDebug() should return false when debug mode is false")
		}
	})
}
