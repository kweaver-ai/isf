package helpers

import (
	"context"

	"AuditLog/common/enums"
	"AuditLog/common/types"
)

func GetUserIDFromCtx(ctx context.Context) (userID string) {
	vInter := ctx.Value(enums.VisitUserIDCtxKey.String())
	if vInter == nil {
		return
	}

	if v, ok := vInter.(string); ok {
		userID = v
	} else {
		panic("GetUserIDFromCtx:ctx.Value(enums.VisitUserIDCtxKey) is not string")
	}

	return
}

func GetUserTokenFromCtx(ctx context.Context) (userToken string) {
	vInter := ctx.Value(enums.VisitUserTokenCtxKey.String())
	if vInter == nil {
		return
	}

	if v, ok := vInter.(string); ok {
		userToken = v
	} else {
		panic("GetUserTokenFromCtx:ctx.Value(enums.VisitUserTokenCtxKey) is not string")
	}

	return
}

func GetTraceIDFromCtx(ctx context.Context) (traceID string) {
	vInter := ctx.Value(enums.TraceIDCtxKey.String())
	if vInter == nil {
		return
	}

	if v, ok := vInter.(string); ok {
		traceID = v
	} else {
		panic("GetTraceIDFromCtx:ctx.Value(enums.TraceIDCtxKey) is not string")
	}

	return
}

//func GetClientInfoFromCtx(ctx context.Context) (clientInfo *types.ClientInfo) {
//	vInter := ctx.Value(enums.ClientInfoCtxKey)
//	if vInter == nil {
//		return
//	}
//
//	if v, ok := vInter.(*types.ClientInfo); ok {
//		clientInfo = v
//	} else {
//		panic("GetClientInfoFromCtx:ctx.Value(enums.ClientInfoCtxKey) is not string")
//	}
//
//	return
//}

func GetVisitUserInfoFromCtx(ctx context.Context) (info *types.VisitUserInfo) {
	vInter := ctx.Value(enums.VisitUserInfoCtxKey.String())
	if vInter == nil {
		return
	}

	if v, ok := vInter.(*types.VisitUserInfo); ok {
		info = v
	} else {
		panic("GetVisitUserInfoFromCtx:ctx.Value(enums.VisitUserInfoCtxKey) is not string")
	}

	return
}

func NewBgCtxWithUserInfo(ctx context.Context) (newCtx context.Context) {
	info := GetVisitUserInfoFromCtx(ctx)

	newCtx = context.Background()

	ctxKey := enums.VisitUserInfoCtxKey.String()

	//nolint:staticcheck
	newCtx = context.WithValue(newCtx, ctxKey, info)

	return
}
