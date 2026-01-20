package thrift

import (
	"context"
	"os"

	"policy_mgnt/tapi/sharemgnt"

	"github.com/kweaver-ai/go-lib/tclient"
	"github.com/kweaver-ai/go-lib/thrift"
	// "policy_mgnt/utils/errors"
	// "policy_mgnt/utils/gocommon/api"
	// "github.com/silenceper/pool"
)

// ShareMgnt sharemgnt thrift 接口封装
//
//go:generate mockgen -package mock_thrift -source ./sharemgnt.go -destination ../test/mock_thrift/sharemgnt_mock.go
type ShareMgnt interface {
	GetHostName() (string, error)
	ApplyPasswordStrengthMeter(sharemgnt.NcTUsrmPasswordConfig) error
	GetPasswordConfig() (*sharemgnt.NcTUsrmPasswordConfig, error)
	ApplyImageVcode(sharemgnt.NcTVcodeConfig) error
	ApplyMultiFactorAuth(string, string) error
	GetUserByID(string) (*sharemgnt.NcTUsrmGetUserInfo, error)
	GetDepartmentByID(string) (*sharemgnt.NcTUsrmDepartmentInfo, error)
	GetOrganizationByID(string) (*sharemgnt.NcTUsrmOrganizationInfo, error)
	GetAllOrgs() ([]*sharemgnt.NcTRootOrgInfo, error)
	SearchSupervisoryUsers(string, string, int, int) ([]*sharemgnt.NcTSearchUserInfo, error)
	SearchDepartments(string, string, int, int) ([]*sharemgnt.NcTUsrmDepartmentInfo, error)
}

type sharemgntImpl struct {
	client    *sharemgnt.NcTShareMgntClient
	transport *thrift.TBufferedTransport
}

// NewShareMgnt init sharemgnt
func NewShareMgnt() (ShareMgnt, error) {
	aShareMgntClient, err := newShareMgntClient()
	if err != nil {
		return nil, err
	}
	return NewShareMgntWithClient(aShareMgntClient), nil
}

// NewShareMgntWithClient 使用指定客户端
func NewShareMgntWithClient(aShareMgntClient *AShareMgntClient) ShareMgnt {
	return &sharemgntImpl{client: aShareMgntClient.client, transport: aShareMgntClient.transport}
}

func (s *sharemgntImpl) open() error {
	if !s.transport.IsOpen() {
		err := s.transport.Open()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sharemgntImpl) close() {
	s.transport.Close()
}

// GetHostName 获取集群域名、地址
func (s *sharemgntImpl) GetHostName() (string, error) {
	err := s.open()
	if err != nil {
		return "", err
	}
	defer s.close()
	return s.GetHostName()
}

// ApplyPasswordStrengthMeter 设置密码强度
func (s *sharemgntImpl) ApplyPasswordStrengthMeter(config sharemgnt.NcTUsrmPasswordConfig) (err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// err = s.client.Usrm_SetPasswordConfig(context.Background(), &config)
	// if err != nil {
	// 	return
	// }
	// return nil
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	err = shareMgntClient.Usrm_SetPasswordConfig(context.Background(), &config)
	if err != nil {
		return
	}
	return
}

// GetPasswordConfig 获取密码相关策略
func (s *sharemgntImpl) GetPasswordConfig() (res *sharemgnt.NcTUsrmPasswordConfig, err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// res, err = s.client.Usrm_GetPasswordConfig(context.Background())
	// if err != nil {
	// 	return
	// }
	// return res, nil
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	res, err = shareMgntClient.Usrm_GetPasswordConfig(context.Background())
	if err != nil {
		return
	}
	return
}

// ApplyImageVcode 设置多因子认证
func (s *sharemgntImpl) ApplyImageVcode(config sharemgnt.NcTVcodeConfig) (err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// err = s.client.Usrm_SetVcodeConfig(context.Background(), &config)
	// if err != nil {
	// 	return
	// }
	// return nil
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	err = shareMgntClient.Usrm_SetVcodeConfig(context.Background(), &config)
	if err != nil {
		return
	}
	return
}

func (s *sharemgntImpl) ApplyMultiFactorAuth(name, config string) (err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// err = s.client.SetCustomConfigOfString(context.Background(), name, config)
	// if err != nil {
	// 	return
	// }
	// return nil
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	err = shareMgntClient.SetCustomConfigOfString(context.Background(), name, config)
	if err != nil {
		return
	}
	return
}

// GetUserByID 根据用户id获取用户信息
func (s *sharemgntImpl) GetUserByID(userID string) (userInfo *sharemgnt.NcTUsrmGetUserInfo, err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// // //defer s.close()
	// userInfo, err = s.client.Usrm_GetUserInfo(context.Background(), userID)
	// if err != nil {
	// 	return
	// }
	// return
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	userInfo, err = shareMgntClient.Usrm_GetUserInfo(context.Background(), userID)
	if err != nil {
		return
	}
	return
}

// GetDepartmentByID根据部门id获取用户信息
func (s *sharemgntImpl) GetDepartmentByID(departmentID string) (departmentInfo *sharemgnt.NcTUsrmDepartmentInfo, err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// departmentInfo, err = s.client.Usrm_GetDepartmentById(context.Background(), departmentID)
	// if err != nil {
	// 	return
	// }
	// return
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	departmentInfo, err = shareMgntClient.Usrm_GetDepartmentById(context.Background(), departmentID)
	if err != nil {
		return
	}
	return
}

// GetOrganizationByID根据组织id获取用户信息
func (s *sharemgntImpl) GetOrganizationByID(orgID string) (orgInfo *sharemgnt.NcTUsrmOrganizationInfo, err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// orgInfo, err = s.client.Usrm_GetOrganizationById(context.Background(), orgID)
	// if err != nil {
	// 	return
	// }
	// return
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	orgInfo, err = shareMgntClient.Usrm_GetOrganizationById(context.Background(), orgID)
	if err != nil {
		return
	}
	return
}

// 获取所有组织
func (s *sharemgntImpl) GetAllOrgs() (orgs []*sharemgnt.NcTRootOrgInfo, err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// orgs, err = s.client.Usrm_GetSupervisoryRootOrg(context.Background(), sharemgnt.NCT_USER_ADMIN)
	// if err != nil {
	// 	return
	// }
	// return
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	orgs, err = shareMgntClient.Usrm_GetSupervisoryRootOrg(context.Background(), sharemgnt.NCT_USER_ADMIN)
	if err != nil {
		return
	}
	return
}

// 搜索管理的用户
func (s *sharemgntImpl) SearchSupervisoryUsers(managerid, keyWord string, start, limit int) (userInfos []*sharemgnt.NcTSearchUserInfo, err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// userInfos, err = s.client.Usrm_SearchSupervisoryUsers(context.Background(), managerid, keyWord, int32(start), int32(limit))
	// if err != nil {
	// 	return
	// }
	// return
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	userInfos, err = shareMgntClient.Usrm_SearchSupervisoryUsers(context.Background(), managerid, keyWord, int32(start), int32(limit))
	if err != nil {
		return
	}
	return
}

// 搜索管理的部门
func (s *sharemgntImpl) SearchDepartments(managerid, keyWord string, start, limit int) (userInfos []*sharemgnt.NcTUsrmDepartmentInfo, err error) {
	// err = s.open()
	// if err != nil {
	// 	return
	// }
	// defer s.close()
	// userInfos, err = s.client.Usrm_SearchDepartments(context.Background(), managerid, keyWord, int32(start), int32(limit))
	// if err != nil {
	// 	return
	// }
	// return
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	if err != nil {
		return
	}

	defer transport.Close()
	userInfos, err = shareMgntClient.Usrm_SearchDepartments(context.Background(), managerid, keyWord, int32(start), int32(limit))
	if err != nil {
		return
	}
	return
}
