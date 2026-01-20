package enums

// UserType 这个现在应该不是从API来的了
// "authenticated_user" : 实名用户
// "anonymous_user"：匿名用户
// "app" ：应用账户
// ”internal_service“：内部服务
type UserType string

func (u UserType) String() string {
	return string(u)
}

const (
	AuthUser        UserType = "authenticated_user"
	AnonyUser       UserType = "anonymous_user"
	App             UserType = "app"
	InternalService UserType = "internal_service"
)
