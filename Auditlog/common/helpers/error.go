package helpers

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"AuditLog/gocommon/api"
)

//go:noinline
func RecordErrLogWithPos(logger api.Logger, err error, positions ...string) {
	sb := strings.Builder{}
	for _, pos := range positions {
		sb.WriteString("[")
		sb.WriteString(pos)
		sb.WriteString("]")
	}

	pc, file, line, ok := runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		fnName := filepath.Base(fn.Name())

		msg := fmt.Errorf("%v error: %v, loc: %s:%d, func:%s\n", sb.String(), err, file, line, fnName)
		if IsAaronLocalDev() {
			log.Println(msg)
		}

		logger.Errorln(msg)
	} else {
		msg := fmt.Errorf("%v error: %v\n", sb.String(), err)

		if IsAaronLocalDev() {
			log.Println(msg)
		}

		logger.Errorln(msg)
	}
}
