package umcmp

import (
	"fmt"

	"AuditLog/common/utils"
)

func (u *Um) getPrivateURLPrefix() string {
	return fmt.Sprintf("%s://%s:%d/api/user-management", u.umConf.Private.Protocol, utils.ParseHost(u.umConf.Private.Host), u.umConf.Private.Port)
}

//nolint:unused
func (u *Um) getPublicURLPrefix() string {
	return fmt.Sprintf("%s://%s:%d/api/user-management", u.umConf.Public.Protocol, utils.ParseHost(u.umConf.Public.Host), u.umConf.Public.Port)
}
