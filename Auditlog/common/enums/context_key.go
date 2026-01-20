package enums

type ContextKey string

func (ck ContextKey) String() string {
	return string(ck)
}

const ContextKeyPrefix = "ctx_k_"

const (
	VisitLangCtxKey ContextKey = ContextKeyPrefix + "visit_language"

	VisitUserIDCtxKey ContextKey = ContextKeyPrefix + "visit_user_id" // 访问用户ID（包括普通用户id和应用账号ID）

	VisitUserTokenCtxKey ContextKey = ContextKeyPrefix + "visit_user_token" // 访问用户token

	TraceIDCtxKey ContextKey = ContextKeyPrefix + "trace_id" // 调用链 ID

	VisitUserInfoCtxKey ContextKey = ContextKeyPrefix + "visit_user_info" // 访问用户信息
)
