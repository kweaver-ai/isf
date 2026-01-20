package panichelper

import (
	"fmt"
	"log"

	"AuditLog/common/helpers"
	"AuditLog/gocommon/api"
)

func Recovery(logger api.Logger) {
	if err := recover(); err != nil {
		// 1、记录日志
		panicLogMsg := PanicTraceErrLog(err)

		if helpers.IsDebugMode() {
			log.Println(panicLogMsg)
		}

		logger.Errorln(panicLogMsg)

		return
	}
}

func RecoveryAndSetErr(logger api.Logger, err *error) {
	if r := recover(); r != nil {
		// 1、记录日志
		panicLogMsg := PanicTraceErrLog(r)
		logger.Errorln(panicLogMsg)

		// 2、设置错误
		if e, ok := r.(error); ok {
			*err = e
		} else {
			*err = fmt.Errorf("%v", r)
		}

		return
	}
}
