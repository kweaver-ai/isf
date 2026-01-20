package service

import "policy_mgnt/utils/gocommon/v2/utils"

// Access 服务访问通用配置
type Access struct {
	Protocol string
	Host     string
	Port     string
}

// 一般来说可以通过如下方式使用 fmt.Sprintf("%s/your/path", Access{})
func (a Access) String() string {
	return utils.ContactStr(a.Protocol, "://", a.Host, ":", a.Port)
}
