// Package logics perm Anyshare 业务逻辑层 -文档权限
package logics

import (
	"context"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"

	gerrors "github.com/kweaver-ai/go-lib/error"

	"Authorization/common"
	"Authorization/interfaces"
)

var (
	resourceOnce      sync.Once
	resourceSingleton *resourceType
)

const (
	nameMaxLength = 50
)

type resourceType struct {
	db       interfaces.DBResourceType
	userMgmt interfaces.DrivenUserMgnt
	logger   common.Logger
}

// NewResourceType 创建新的NewResourceType对象
func NewResourceType() *resourceType {
	resourceOnce.Do(func() {
		resourceSingleton = &resourceType{
			db:       dbResourceType,
			logger:   common.NewLogger(),
			userMgmt: dnUserMgnt,
		}
	})
	return resourceSingleton
}

func (r *resourceType) checkVisitorType(ctx context.Context, visitor *interfaces.Visitor) (err error) {
	// 获取访问者角色
	var roleTypes []interfaces.SystemRoleType

	// 实名用户获取对应角色信息
	if visitor.Type == interfaces.RealName {
		// 获取用户角色信息
		roleTypes, err = r.userMgmt.GetUserRolesByUserID(ctx, visitor.ID)
		if err != nil {
			return err
		}
	}

	return checkVisitorType(
		visitor,
		roleTypes,
		[]interfaces.VisitorType{interfaces.RealName},
		[]interfaces.SystemRoleType{interfaces.SuperAdmin, interfaces.SystemAdmin, interfaces.SecurityAdmin},
	)
}

// GetPagination 获取资源
func (r *resourceType) GetPagination(ctx context.Context, visitor *interfaces.Visitor, params interfaces.ResourceTypePagination) (count int, resources []interfaces.ResourceType, err error) {
	// 权限检查 visitor
	err = r.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}

	count, resources, err = r.db.GetPagination(ctx, params)
	if err != nil {
		r.logger.Errorf("GetPagination: %v", err)
		return 0, nil, err
	}
	return count, resources, nil
}

// Set 设置资源
func (r *resourceType) Set(ctx context.Context, visitor *interfaces.Visitor, resourceType *interfaces.ResourceType) (err error) {
	// 权限检查 visitor
	err = r.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}
	// id 长度 不超过50 ，只能是数字或者字母或者下划线
	if len(resourceType.ID) > nameMaxLength {
		err = gerrors.NewError(gerrors.PublicBadRequest, "id length must less than 50")
		return
	}
	// 数字或者字母或者下划线
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(resourceType.ID) {
		err = gerrors.NewError(gerrors.PublicBadRequest, "invalid id, only number or letter or underline")
		return
	}

	err = r.db.Set(ctx, resourceType)
	if err != nil {
		r.logger.Errorf("Set: %v", err)
		return
	}
	return
}

// checkResourceTypeChange 检查资源类型是否发生变化， 有变化返回true
func (r *resourceType) checkResourceTypeChange(old, newInfo *interfaces.ResourceType) bool {
	if old.Name != newInfo.Name {
		return true
	}
	if old.InstanceURL != newInfo.InstanceURL {
		return true
	}
	if old.DataStruct != newInfo.DataStruct {
		return true
	}
	if old.Description != newInfo.Description {
		return true
	}
	if old.Hidden != newInfo.Hidden {
		return true
	}
	if !reflect.DeepEqual(old.Operation, newInfo.Operation) {
		return true
	}
	return false
}

// InitResourceTypes 资源类型批量添加
func (r *resourceType) InitResourceTypes(ctx context.Context, resourceTypes []interfaces.ResourceType) (err error) {
	for i := range resourceTypes {
		var info interfaces.ResourceType
		// 如果之前存在，则不添加
		// 所有用户可以调用
		var infoMap map[string]interfaces.ResourceType
		infoMap, err = r.db.GetByIDs(ctx, []string{resourceTypes[i].ID})
		if err != nil {
			r.logger.Errorf("InitResourceTypes GetByIDs: %v", err)
			return
		}
		info, ok := infoMap[resourceTypes[i].ID]
		if ok && !r.checkResourceTypeChange(&info, &resourceTypes[i]) {
			continue
		}
		err = r.db.Set(ctx, &resourceTypes[i])
		if err != nil {
			r.logger.Errorf("InitResourceTypes: %v", err)
			return
		}
	}
	return nil
}

// Delete 删除资源
func (r *resourceType) Delete(ctx context.Context, visitor *interfaces.Visitor, resourceTypeID string) (err error) {
	// 权限检查 visitor
	err = r.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}
	err = r.db.Delete(ctx, resourceTypeID)
	if err != nil {
		r.logger.Errorf("Delete: %v", err)
		return err
	}
	return nil
}

// GetByID 获取资源
func (r *resourceType) GetByID(ctx context.Context, _ *interfaces.Visitor, resourceTypeID string) (resourceType interfaces.ResourceType, err error) {
	// 所有用户可以调用
	resourceTypeMap, err := r.db.GetByIDs(ctx, []string{resourceTypeID})
	if err != nil {
		r.logger.Errorf("GetByIDs: %v", err)
		return
	}
	resourceType, ok := resourceTypeMap[resourceTypeID]
	if !ok {
		return interfaces.ResourceType{}, nil
	}
	return
}

// GetAllOperation 获取资源类型所有操作
func (r *resourceType) GetAllOperation(ctx context.Context, visitor *interfaces.Visitor, resourceTypeID string, scope interfaces.OperationScopeType) (
	operations []interfaces.ResourceTypeOperationResponse, err error,
) {
	// 所有用户可以调用
	resourceTypeMap, err := r.db.GetByIDs(ctx, []string{resourceTypeID})
	if err != nil {
		r.logger.Errorf("GetAllOperation GetByIDs err: %v", err)
		return
	}
	resourceTypeInfo, ok := resourceTypeMap[resourceTypeID]
	if !ok {
		return
	}
	operations = make([]interfaces.ResourceTypeOperationResponse, 0, len(resourceTypeInfo.Operation))
	for _, operation := range resourceTypeInfo.Operation {
		opeInfo := interfaces.ResourceTypeOperationResponse{
			ID:          operation.ID,
			Description: operation.Description,
		}
		if slices.Contains(operation.Scope, scope) {
			// 根据国际化获取名称
			opeInfo.Name = r.getOperationNameByLanguage(visitor.Language, operation.Name)
			operations = append(operations, opeInfo)
		}
	}
	return
}

func (r *resourceType) getOperationNameByLanguage(language string, operationName []interfaces.OperationName) (name string) {
	for _, name := range operationName {
		if strings.EqualFold(name.Language, language) {
			return name.Value
		}
	}
	return
}

// GetByIDsInternal 批量获取资源
func (r *resourceType) GetByIDsInternal(ctx context.Context, resourceTypeIDs []string) (resourceMap map[string]interfaces.ResourceType, err error) {
	resourceMap, err = r.db.GetByIDs(ctx, resourceTypeIDs)
	if err != nil {
		r.logger.Errorf("GetByIDsInternal: %v", err)
		return nil, err
	}
	return
}

func (r *resourceType) GetAllInternal(ctx context.Context) (resourceTypes []interfaces.ResourceType, err error) {
	resourceTypes, err = r.db.GetAllInternal(ctx)
	if err != nil {
		r.logger.Errorf("GetAllInternal: %v", err)
		return nil, err
	}
	return
}

// SetPrivate 设置资源
func (r *resourceType) SetPrivate(ctx context.Context, resourceType *interfaces.ResourceType) (err error) {
	// id 长度 不超过50 ，只能是数字或者字母或者下划线
	if len(resourceType.ID) > nameMaxLength {
		err = gerrors.NewError(gerrors.PublicBadRequest, "id length must less than 50")
		return
	}
	// 数字或者字母或者下划线
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(resourceType.ID) {
		err = gerrors.NewError(gerrors.PublicBadRequest, "invalid id, only number or letter or underline")
		return
	}

	// 如果之前存在，则不添加
	// 所有用户可以调用
	var infoMap map[string]interfaces.ResourceType
	infoMap, err = r.db.GetByIDs(ctx, []string{resourceType.ID})
	if err != nil {
		r.logger.Errorf("InitResourceTypes GetByIDs: %v", err)
		return
	}
	// 存在且没变化
	info, ok := infoMap[resourceType.ID]
	if ok && !r.checkResourceTypeChange(&info, resourceType) {
		return
	}

	err = r.db.Set(ctx, resourceType)
	if err != nil {
		r.logger.Errorf("Set: %v", err)
		return
	}
	return
}
