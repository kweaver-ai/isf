package helpers

import (
	"runtime/debug"

	"AuditLog/gocommon/api"
)

// GoSafe 安全地执行一个 goroutine，会自动捕获和处理 panic
func GoSafe(logger api.Logger, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("goroutine panic: %v\nstack:\n%s\n", r, debug.Stack())
			}
		}()

		f()
	}()
}
