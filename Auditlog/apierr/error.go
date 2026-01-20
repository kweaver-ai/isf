package apierr

import (
	"context"
	stdErrors "errors"

	"AuditLog/errors"
)

func NewCtx(ctx context.Context, code int, cause string, detail interface{}) *errors.ErrorResp {
	return errors.NewCtx(ctx, code, cause, detail)
}

func ToErrorResp(err error) (errorResp *errors.ErrorResp, ok bool) {
	ok = stdErrors.As(err, &errorResp)
	return
}
