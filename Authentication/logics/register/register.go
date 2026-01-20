// Package register 逻辑层
package register

import (
	"encoding/json"
	"strings"
	"sync"

	"Authentication/interfaces"
	"Authentication/logics"

	"github.com/kweaver-ai/go-lib/rest"
	"golang.org/x/crypto/bcrypt"
)

var (
	rOnce sync.Once
	r     *register
)

type register struct {
	hydraAdmin interfaces.DnHydraAdmin
	db         interfaces.DBRegister
	cost       int
	conf       interfaces.DBConf
}

// NewRegister 新建注册接口操作对象
func NewRegister() *register {
	rOnce.Do(func() {
		r = &register{
			hydraAdmin: logics.DnHydraAdmin,
			db:         logics.DBRegister,
			cost:       12,
			conf:       logics.DBConf,
		}
	})

	return r
}

func (r *register) PublicRegister(client *interfaces.RegisterInfo) (clientInfo interfaces.ClientInfo, err error) {
	// 检查客户端的redirect url是否符合配置
	conf, err := r.conf.GetConfig(map[interfaces.ConfigKey]bool{
		interfaces.LimitRedirectURI: true,
	})
	if err != nil {
		return clientInfo, err
	}

	if len(conf.LimitRedirectURI) > 0 {
		for _, redirectURI := range client.RedirectURIs {
			// 前缀匹配
			bChecked := false
			for limitRedirectURI := range conf.LimitRedirectURI {
				if strings.HasPrefix(redirectURI, limitRedirectURI) {
					bChecked = true
					break
				}
			}

			if !bChecked {
				return clientInfo, rest.NewHTTPError("redirect url is not allowed: "+redirectURI, rest.BadRequest, nil)
			}
		}
	}
	// 注册
	info, err := r.hydraAdmin.PublicRegister(client)
	if err != nil {
		return
	}

	secret, err := r.hash([]byte(info.ClientSecret))
	if err != nil {
		return
	}

	metadata, err := json.Marshal(client.Metadata)
	if err != nil {
		return
	}

	dbClientData := &interfaces.DBRegisterInfo{
		ClientID:               info.ClientID,
		ClientName:             client.ClientName,
		ClientSecret:           string(secret),
		GrantTypes:             strings.Join(client.GrantTypes, "|"),
		ResponseTypes:          strings.Join(client.ResponseTypes, "|"),
		Scope:                  client.Scope,
		RedirectURIs:           strings.Join(client.RedirectURIs, "|"),
		PostLogoutRedirectURIs: strings.Join(client.PostLogoutRedirectURIs, "|"),
		Metadata:               metadata,
	}

	err = r.db.CreateClient(dbClientData)
	if err != nil {
		return
	}

	clientInfo = *info
	return
}

func (r *register) hash(data []byte) ([]byte, error) {
	s, err := bcrypt.GenerateFromPassword(data, r.cost)
	if err != nil {
		return nil, err
	}
	return s, nil
}
