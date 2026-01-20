package config

import "AuditLog/common/utils"

type UmConf struct {
	Private SvcConf `yaml:"private"`
	Public  SvcConf `yaml:"public"`
}

func (u *UmConf) loadConf() {
	// 1. private
	u.Private.Protocol = utils.GetEnv("USER_MANAGEMENT_PRIVATE_PROTOCOL", "http")
	u.Private.Host = utils.GetEnv("USER_MANAGEMENT_PRIVATE_HOST", "user-management-private.anyshare")
	u.Private.Port = utils.GetEnvMustInt("USER_MANAGEMENT_PRIVATE_PORT", 30980)

	// 2. public
	u.Public.Protocol = utils.GetEnv("USER_MANAGEMENT_PUBLIC_PROTOCOL", "http")
	u.Public.Host = utils.GetEnv("USER_MANAGEMENT_PUBLIC_HOST", "user-management-public.anyshare")
	u.Public.Port = utils.GetEnvMustInt("USER_MANAGEMENT_PUBLIC_PORT", 30980)
}
