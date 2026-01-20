package thrift

import (
	"context"
	"sync"

	"github.com/kweaver-ai/go-lib/tclient"

	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra/config"
	"AuditLog/interfaces"
	"AuditLog/tapi/sharemgnt"
)

var (
	smOnce sync.Once
	sm     *shareMgnt
)

type shareMgnt struct {
	Host   string
	Port   int
	logger api.Logger
}

func NewShareMgnt() interfaces.ShareMgntRepo {
	smOnce.Do(func() {
		conf := config.GetShareMgntConf()
		sm = &shareMgnt{
			Host:   conf.Host,
			Port:   conf.Port,
			logger: drivenadapters.Logger,
		}
	})
	return sm
}

// GetUsrRolemDeptIDs 获取用户ID和角色ID获取管辖范围内的部门信息
func (s *shareMgnt) GetRoleMemberInfos(roleID string) (res []*sharemgnt.NcTRoleMemberInfo, err error) {
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, s.Host, s.Port)
	if err != nil {
		s.logger.Errorf("[GetUsrRolemDeptIDs]: NewShareMgntClientFactory: %v", err)
		return
	}

	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorf("[GetUsrRolemDeptIDs]: ShareMgntTransport.Close: %v", closeErr)
		}
	}()

	var UserID string

	status, err := s.GetTriSystemStatus()
	if err != nil {
		s.logger.Errorf("[GetUsrRolemDeptIDs]: GetTriSystemStatus: %v", err)
		return nil, err
	}

	if status {
		UserID = sharemgnt.NCT_USER_SECURIT
	} else {
		UserID = sharemgnt.NCT_USER_ADMIN
	}

	// 三权分立下，security可调此接口，全责集中下admin可调
	if res, err = shareMgntClient.UsrRolem_GetMember(context.Background(), UserID, roleID); err != nil {
		s.logger.Errorf("[GetUsrRolemDeptIDs]: UsrRolem_GetMember: %v", err)
	}
	return
}

// GetTriSystemStatus 获取是否开启三权分立
func (s *shareMgnt) GetTriSystemStatus() (res bool, err error) {
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, s.Host, s.Port)
	if err != nil {
		s.logger.Errorf("[GetTriSystemStatus]: NewShareMgntClientFactory: %v", err)
		return
	}

	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorf("[GetTriSystemStatus]: ShareMgntTransport.Close: %v", closeErr)
		}
	}()

	if res, err = shareMgntClient.Usrm_GetTriSystemStatus(context.Background()); err != nil {
		s.logger.Errorf("[GetTriSystemStatus]: Usrm_GetTriSystemStatus: %v", err)
	}
	return
}
