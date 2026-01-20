package efastcmp

import (
	"fmt"

	"AuditLog/common/utils"
)

func (e *EFast) getUrlPrefix() string {
	return fmt.Sprintf("%s://%s:%d/api/efast", e.privateScheme, utils.ParseHost(e.privateHost), e.privatePort)
}

func (e *EFast) getPublicUrlPrefix() string {
	return fmt.Sprintf("%s://%s:%d/api/efast", e.publicScheme, utils.ParseHost(e.publicHost), e.publicPort)
}
