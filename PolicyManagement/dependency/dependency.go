package dependency

import (
	"sync"

	svcac "policy_mgnt/adapter/driven/service_access"
)

//go:generate mockgen -source=./dependency.go -package=mock_dependency --destination=../test/mock_dependency/dependency_mock.go

var (
	abstractDrivenOnce      sync.Once
	abstractDrivenSingleton *abstractDriven
)

type AbstractDriven interface {
	GetBelongDepartByUserId(userId string) ([]string, error)
	GetDepartUserIds(departmentId string) ([]string, error)
}

type abstractDriven struct {
	userMgmt svcac.UsermgntDriven
}

func NewAbstractDriven() AbstractDriven {
	abstractDrivenOnce.Do(func() {
		abstractDrivenSingleton = &abstractDriven{
			userMgmt: svcac.NewUsermgntDriven(),
		}
	})
	return abstractDrivenSingleton
}

var _ AbstractDriven = (*abstractDriven)(nil)

func (a *abstractDriven) GetBelongDepartByUserId(userId string) ([]string, error) {
	return a.userMgmt.GetBelongDepartByUserId(userId)
}

func (a *abstractDriven) GetDepartUserIds(departmentId string) ([]string, error) {
	return a.userMgmt.GetDepartUserIds(departmentId)
}
