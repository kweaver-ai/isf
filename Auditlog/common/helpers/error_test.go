package helpers

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"AuditLog/test/mock_log"
)

func TestRecordErrLogWithPos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	logger.EXPECT().Errorln(gomock.Any()).AnyTimes()

	err := errors.New("test error")
	RecordErrLogWithPos(logger, err, "test")
}
