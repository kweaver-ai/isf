// Package logics perm Anyshare 业务逻辑层 -文档权限
package logics

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/xeipuuv/gojsonschema"

	"Authorization/common"
	"Authorization/interfaces"
)

var (
	obligationOnce      sync.Once
	obligationSingleton *obligation
)

type obligation struct {
	db             interfaces.DBObligation
	obligationType interfaces.ObligationType
	userMgnt       interfaces.DrivenUserMgnt
	logger         common.Logger
}

// NewObligation 创建新的对象
func NewObligation() *obligation {
	obligationOnce.Do(func() {
		obligationSingleton = &obligation{
			db:             dbObligation,
			userMgnt:       dnUserMgnt,
			obligationType: NewObligationType(),
			logger:         common.NewLogger(),
		}
	})
	return obligationSingleton
}

func (ob *obligation) checkVisitorType(ctx context.Context, visitor *interfaces.Visitor) (err error) {
	// 获取访问者角色
	var roleTypes []interfaces.SystemRoleType

	// 实名用户获取对应角色信息
	if visitor.Type == interfaces.RealName {
		// 获取用户角色信息
		roleTypes, err = ob.userMgnt.GetUserRolesByUserID(ctx, visitor.ID)
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

func (ob *obligation) checkJsonValueValid(_ context.Context, schemaTmp, value any) (err error) {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(schemaTmp))
	if err != nil {
		err = gerrors.NewError(gerrors.PublicBadRequest, "schema is invalid")
		ob.logger.Errorf("checkJsonValueValid: %v", err)
		return
	}

	result, err := schema.Validate(gojsonschema.NewGoLoader(value))
	if err != nil {
		return gerrors.NewError(gerrors.PublicBadRequest, err.Error())
	}

	if !result.Valid() {
		msgList := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			msgList = append(msgList, err.String())
		}
		return gerrors.NewError(gerrors.PublicBadRequest, strings.Join(msgList, "; "))
	}

	return
}

// Add 添加义务
func (ob *obligation) Add(ctx context.Context, visitor *interfaces.Visitor, info *interfaces.ObligationInfo) (id string, err error) {
	// 角色检查 visitor
	err = ob.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}

	obligationTypeIDs := make(map[string]bool)
	obligationTypeIDs[info.TypeID] = true
	obligationTypeInfos, err := ob.obligationType.GetByIDSInternal(ctx, obligationTypeIDs)
	if err != nil {
		ob.logger.Errorf("Add: %v", err)
		return
	}

	if len(obligationTypeInfos) == 0 {
		ob.logger.Errorf("Add:obligation_type_id %s not found", info.TypeID)
		return "", gerrors.NewError(gerrors.PublicBadRequest, fmt.Sprintf("obligation_type_id %s not found", info.TypeID))
	}

	// 校验jsonschema
	err = ob.checkJsonValueValid(ctx, obligationTypeInfos[0].Schema, info.Value)
	if err != nil {
		ob.logger.Errorf("Add:checkJsonValueValid %v", err)
		return
	}

	// 服务端生成ID
	info.ID = uuid.New().String()
	id = info.ID
	err = ob.db.Add(ctx, info)
	if err != nil {
		ob.logger.Errorf("Add: %v", err)
		return
	}
	return
}

func (ob *obligation) Delete(ctx context.Context, visitor *interfaces.Visitor, obligationID string) (err error) {
	// 角色检查 visitor
	err = ob.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}

	err = ob.db.Delete(ctx, obligationID)
	if err != nil {
		ob.logger.Errorf("Delete: %v", err)
		return
	}
	return
}

func (ob *obligation) Update(ctx context.Context, visitor *interfaces.Visitor, obligationID string, name string, nameChanged bool,
	description string, descriptionChanged bool, value any, valueChanged bool,
) (err error) {
	// 角色检查 visitor
	err = ob.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}

	// 检查义务是否存在
	obligationInfo, err := ob.db.GetByID(ctx, obligationID)
	if err != nil {
		ob.logger.Errorf("Update: %v", err)
		return
	}

	if obligationInfo.ID == "" {
		ob.logger.Errorf("Update:obligation_id %s not found", obligationID)
		return gerrors.NewError(gerrors.PublicNotFound, fmt.Sprintf("obligation_id %s not found", obligationID))
	}

	// 校验schema
	if valueChanged {
		var obligationTypeInfos []interfaces.ObligationTypeInfo
		obligationTypeIDs := make(map[string]bool)
		obligationTypeIDs[obligationInfo.TypeID] = true

		obligationTypeInfos, err = ob.obligationType.GetByIDSInternal(ctx, obligationTypeIDs)
		if err != nil {
			ob.logger.Errorf("Update: %v", err)
			return
		}

		if len(obligationTypeInfos) == 0 {
			ob.logger.Errorf("Update:obligation_type_id %s not found", obligationInfo.TypeID)
			return gerrors.NewError(gerrors.PublicNotFound, fmt.Sprintf("obligation_type_id %s not found", obligationInfo.TypeID))
		}
		// 校验jsonschema
		err = ob.checkJsonValueValid(ctx, obligationTypeInfos[0].Schema, value)
		if err != nil {
			ob.logger.Errorf("Update:checkJsonValueValid %v", err)
			return err
		}
	}

	err = ob.db.Update(ctx, obligationID, name, nameChanged, description, descriptionChanged, value, valueChanged)
	if err != nil {
		ob.logger.Errorf("Update: %v", err)
		return
	}
	return
}

func (ob *obligation) GetByID(ctx context.Context, visitor *interfaces.Visitor, obligationID string) (info interfaces.ObligationInfo, err error) {
	// 角色检查 visitor
	err = ob.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}
	info, err = ob.db.GetByID(ctx, obligationID)
	if err != nil {
		ob.logger.Errorf("GetByID: %v", err)
		return
	}

	// 义务不存在
	if info.ID == "" {
		err = gerrors.NewError(gerrors.PublicNotFound, fmt.Sprintf("obligation_id %s not found ", obligationID))
		return
	}
	return
}

func (ob *obligation) Get(ctx context.Context, visitor *interfaces.Visitor, info *interfaces.ObligationSearchInfo) (count int, resultInfos []interfaces.ObligationInfo, err error) {
	// 角色检查 visitor
	err = ob.checkVisitorType(ctx, visitor)
	if err != nil {
		return
	}
	count, resultInfos, err = ob.db.Get(ctx, info)
	if err != nil {
		ob.logger.Errorf("Get: %v", err)
		return
	}
	return
}

// Query 查询义务
//
//	queryInfo 是查询条件
//	resultInfos 是查询结果 key是OperationID, value是ObligationInfo列表
func (o *obligation) Query(ctx context.Context, visitor *interfaces.Visitor, queryInfo *interfaces.QueryObligationInfo) (resultInfos map[string][]interfaces.ObligationInfo, err error) {
	o.logger.Debugf("Query: queryInfo %+v", queryInfo)
	// 获取义务类型
	queryTypeInfo := interfaces.QueryObligationTypeInfo{
		ResourceType: queryInfo.ResourceType,
		Operation:    queryInfo.Operation,
	}

	// 获取操作对应的义务类型
	// operationAndObligationTypes key是OperationID, value是ObligationTypeInfo列表
	operationAndObligationTypes, err := o.obligationType.Query(ctx, visitor, &queryTypeInfo)
	if err != nil {
		o.logger.Errorf("Query: %v", err)
		return
	}

	resultInfos = make(map[string][]interfaces.ObligationInfo, 0)
	for operationID := range operationAndObligationTypes {
		resultInfos[operationID] = []interfaces.ObligationInfo{}
	}

	// 过滤义务类型
	allObligationTypeMap := make(map[string]bool)
	if len(queryInfo.ObligationTypeIDs) != 0 {
		typeMap := make(map[string]bool)
		for _, obligationTypeID := range queryInfo.ObligationTypeIDs {
			typeMap[obligationTypeID] = true
		}
		// 过滤义务类型
		tmp := make(map[string][]interfaces.ObligationTypeInfo, 0)
		for operationID, typeInfos := range operationAndObligationTypes {
			for i := range typeInfos {
				if typeMap[typeInfos[i].ID] {
					allObligationTypeMap[typeInfos[i].ID] = true
					tmp[operationID] = append(tmp[operationID], typeInfos[i])
				}
			}
		}
		operationAndObligationTypes = tmp
	} else {
		for _, typeInfos := range operationAndObligationTypes {
			for i := range typeInfos {
				allObligationTypeMap[typeInfos[i].ID] = true
			}
		}
	}

	o.logger.Debugf("Query: operationAndObligationTypes %+v", operationAndObligationTypes)
	if len(allObligationTypeMap) == 0 {
		return
	}

	// 获取义务类型对应的义务
	// resultTypeAndObligation key是ObligationTypeID, value是ObligationInfo列表
	resultTypeAndObligation, err := o.db.GetByObligationTypeIDs(ctx, allObligationTypeMap)
	if err != nil {
		o.logger.Errorf("Query: %v", err)
		return
	}
	o.logger.Debugf("Query: resultTypeAndObligation %+v", resultTypeAndObligation)

	// 从 operationAndObligationTypes, resultTypeAndObligation 中合并义务类型
	for OperationID, typeInfos := range operationAndObligationTypes {
		for i := range typeInfos {
			resultInfos[OperationID] = append(resultInfos[OperationID], resultTypeAndObligation[typeInfos[i].ID]...)
		}
	}
	return
}

// 指定ID批量获取义务
func (o *obligation) GetByIDSInternal(ctx context.Context, obligationIDs map[string]bool) (infos []interfaces.ObligationInfo, err error) {
	o.logger.Debugf("GetByIDSInternal obligationIDs: %v", obligationIDs)
	IDs := make([]string, 0, len(obligationIDs))
	for ID := range obligationIDs {
		IDs = append(IDs, ID)
	}
	// 从数据库获取所有义务
	infos, err = o.db.GetByIDs(ctx, IDs)
	if err != nil {
		o.logger.Errorf("GetByIDSInternal GetByIDs err: %v", err)
		return
	}
	return
}
