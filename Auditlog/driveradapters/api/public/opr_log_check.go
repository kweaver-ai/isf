package driveradapters

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"AuditLog/common/constants/oprlogconsts"
	"AuditLog/common/enums"
	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/enums/oprlogenums/dvenums"
	"AuditLog/common/enums/oprlogenums/mbcenums"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
	"AuditLog/domain/entity/oprlogeo/dveo"
	"AuditLog/domain/entity/oprlogeo/mbceo"
	"AuditLog/errors"
	"AuditLog/infra/json_schema/jsc_opr_log"
)

// logHandler 日志校验 入口
func (l *logHandler) checkLogs(c *gin.Context, bodyData []byte, bizType oprlogenums.BizType) (err error) {
	// 1. 校验json schema
	err = l.checkJsonSchema(c, bodyData, bizType)
	if err != nil {
		return
	}

	var entryEos []*oprlogeo.LogEntry

	err = utils.JSON().Unmarshal(bodyData, &entryEos)
	if err != nil {
		return
	}

	if len(entryEos) == 0 {
		err = errors.NewCtx(c, errors.BadRequestErr, "请求体不能为空", "")
		return
	}

	// 2. 校验日志数量
	err = l.checkSize(c, entryEos)
	if err != nil {
		return
	}

	// 3. 校验批量日志的operator.id、operator.type、operation、recorder、referer.current、referer.previous是否一致
	err = l.checkBatchLogs(c, entryEos)
	if err != nil {
		return
	}

	// 4. 校验dir_visit字段
	if bizType == oprlogenums.DirVisit {
		err = l.checkDirVisit(c, entryEos)
		if err != nil {
			return
		}
	}

	// 5. 校验menu_button_click字段
	if bizType == oprlogenums.MenuButtonClick {
		err = l.checkMenuButtonClick(c, entryEos)
		if err != nil {
			return
		}
	}

	return
}

// checkJsonSchema  验证json schema
func (l *logHandler) checkJsonSchema(c *gin.Context, logBys []byte, bizType oprlogenums.BizType) (err error) {
	_, invalidFields, err := jsc_opr_log.ValidateOprLogJSONSchema(logBys, bizType)
	if err != nil {
		return
	}

	if len(invalidFields) != 0 {
		err = errors.NewCtx(c,
			errors.BadRequestErr,
			"json schema check failed",
			invalidFields)

		if helpers.IsOprLogShowLogForDebug() {
			// 验证失败时打印格式化的json
			formattedLogJSON, _err := utils.FormatJSONString(string(logBys))
			if _err != nil {
				err = _err
				return
			}

			fmt.Printf("[api report json] biz_type[%v]:\n%v\n",
				bizType, formattedLogJSON)
		}

		return
	}

	return
}

// checkSize 验证日志数量
func (l *logHandler) checkSize(c *gin.Context, entryEos []*oprlogeo.LogEntry) (err error) {
	// 检查日志数量
	if len(entryEos) > oprlogconsts.MaxLogSizePerReport {
		err = errors.NewCtx(c, errors.BadRequestErr, fmt.Sprintf("log size exceed %d", oprlogconsts.MaxLogSizePerReport), "")
	}

	return
}

// checkBatchLogs 验证批量日志的operator.id、operator.type、operation、recorder、referer.current、referer.previous是否一致
func (l *logHandler) checkBatchLogs(c *gin.Context, entryEos []*oprlogeo.LogEntry) (err error) {
	// 1. 验证entryEos中的每一个entry的operator.id、operator.type、operation、recorder、referer.current、referer.previous是否一致
	// 1.1. 获取第一个entry的operator.id、operator.type、operation、recorder、referer.current、referer.previous
	var (
		operatorId      string
		operatorType    enums.UserType
		refererCurrent  string
		refererPrevious string
	)

	if entryEos[0].Operator != nil {
		operatorId = entryEos[0].Operator.ID
		operatorType = entryEos[0].Operator.Type
	}

	operation := entryEos[0].Operation
	recorder := entryEos[0].Recorder

	if entryEos[0].Referer != nil {
		refererCurrent = entryEos[0].Referer.Current
		refererPrevious = entryEos[0].Referer.Previous
	}

	// 1.2. 遍历entryEos，验证operator.id、operator.type、operation、recorder、referer.current、referer.previous是否一致
	for _, entry := range entryEos {
		var (
			_refererCurrent  string
			_refererPrevious string

			_operatorId   string
			_operatorType enums.UserType
		)

		if entry.Operator != nil {
			_operatorId = entry.Operator.ID
			_operatorType = entry.Operator.Type
		}

		if entry.Referer != nil {
			_refererCurrent = entry.Referer.Current
			_refererPrevious = entry.Referer.Previous
		}

		isNotEqual := _operatorId != operatorId ||
			_operatorType != operatorType ||
			entry.Operation != operation ||
			entry.Recorder != recorder ||
			_refererCurrent != refererCurrent ||
			_refererPrevious != refererPrevious

		if isNotEqual {
			err = errors.NewCtx(c, errors.BadRequestErr, "批量接口中的每条日志的operator.id、operator.type、operation、recorder、referer.current、referer.previous必须一致", "")
			return
		}
	}

	return
}

// checkDirVisit 验证dir_visit业务类型字段
func (l *logHandler) checkDirVisit(c *gin.Context, eos []*oprlogeo.LogEntry) (err error) {
	for _, eo := range eos {
		object := eo.Object
		detail := eo.Detail

		detailEo := dveo.NewDetail()

		err = detailEo.LoadByInterface(detail)
		if err != nil {
			return
		}

		// 1. 验证object相关字段
		err = l.checkDvObject(c, object, detailEo)
		if err != nil {
			return
		}

		//	2. 验证detail.from_object相关字段
		fromObject := detailEo.FromObject
		switch fromObject.Type {
		case dvenums.NormalDirFOT, dvenums.AssociationFileFOT: // 文件夹|文档库
			if fromObject.ID == "" {
				err = errors.NewCtx(c, errors.BadRequestErr, "当from_object.type为normal_dir或associationFile时，from_object.id不能为空", "")
				return
			}

			// 校验from_object.ID是否为gns://xxx/xxx/xxx格式
			if !strings.HasPrefix(fromObject.ID, "gns://") && len(fromObject.ID) != 32 {
				err = errors.NewCtx(c, errors.BadRequestErr, "当from_object.type为normal_dir或associationFile时，from_object.id格式不正确，应该为gns://xxx格式或者objectID", "")
				return
			}
		}
	}

	return
}

// checkDvObject 验证dir_visit业务类型的object相关字段
func (l *logHandler) checkDvObject(c *gin.Context, object *oprlogeo.ObjectInfo, detailEo *dveo.Detail) (err error) {
	switch object.Type {
	case dvenums.FolderOT.String(): // 文件夹|文档库
		if object.ID == "" {
			err = errors.NewCtx(c, errors.BadRequestErr, "当object.type为folder时，object.id不能为空", "")
			return
		}
		// 校验object.ID是否为gns://xxx/xxx/xxx格式
		if !strings.HasPrefix(object.ID, "gns://") && len(object.ID) != 32 {
			err = errors.NewCtx(c, errors.BadRequestErr, "当object.type为folder时，object.id格式不正确，应该为gns://xxx格式或者objectID", "")
			return
		}
	case dvenums.FavCategoryOT.String(): // 收藏夹分类
		// 1.1. 验证fav_category相关字段
		if detailEo.Object == nil || detailEo.Object.FavCategory == nil {
			err = errors.NewCtx(c, errors.BadRequestErr, "当object.type为fav_category时，detail.object.fav_category不能为空", "")
			return
		}

		fc := detailEo.Object.FavCategory
		if fc.ID == "" || fc.Name == "" || fc.IDPath == "" || fc.NamePath == "" {
			err = errors.NewCtx(c, errors.BadRequestErr, "当object.type为fav_category时，detail.object.fav_category的id、name、id_path、name_path不能为空", "")
			return
		}

		// 1.2. 验证doc_lib相关字段
		if object.DocLib != nil && object.DocLib.ID != "" {
			err = errors.NewCtx(c, errors.BadRequestErr, "当object.type为fav_category时，不应该传object.doc_lib相关信息", "")
			return
		}
	}

	return
}

// checkMenuButtonClick 验证menu_button_click业务类型字段
func (l *logHandler) checkMenuButtonClick(c *gin.Context, eos []*oprlogeo.LogEntry) (err error) {
	for _, eo := range eos {
		// 1. object相关字段校验
		object := eo.Object

		if object != nil {
			switch object.Type {
			case mbcenums.FolderOT.String(), mbcenums.FileOT.String():
				if object.ID == "" {
					err = errors.NewCtx(c, errors.BadRequestErr, "当object.type为folder或file时，object.id不能为空", "")
					return
				}

				// 校验object.ID是否为gns://xxx/xxx/xxx格式
				if !strings.HasPrefix(object.ID, "gns://") && len(object.ID) != 32 {
					err = errors.NewCtx(c, errors.BadRequestErr, "当object.type为folder或file时，object.id格式不正确，应该为gns://xxx格式或者objectID", "")
					return
				}
			}
		}

		// 2. detail相关字段校验
		err = l.checkMbcDetail(c, eo)
		if err != nil {
			return
		}
	}

	return
}

// checkMbcDetail 验证menu_button_click业务类型的detail字段
func (l *logHandler) checkMbcDetail(c *gin.Context, entry *oprlogeo.LogEntry) (err error) {
	detail := entry.Detail

	detailEo := mbceo.NewDetail()

	err = detailEo.LoadByInterface(detail)
	if err != nil {
		return
	}

	// 2.1. 验证detailEo.Position.PathID
	if detailEo.Position.IsDocLibPosition() {
		pathID := detailEo.Position.PathID
		if pathID == "" {
			err = errors.NewCtx(c, errors.BadRequestErr, "当position.type为dir时，position.path_id不能为空", "")
			return
		}

		// 校验position.path_id是否为gns://xxx/xxx/xxx格式
		if !strings.HasPrefix(pathID, "gns://") && len(pathID) == 32 {
			err = errors.NewCtx(c, errors.BadRequestErr, "当position.type为dir时，position.path_id格式不正确，应该为gns://xxx格式或者objectID", "")
			return
		}
	}

	// 2.2. 验证detailEo.Position.AreaType
	if detailEo.OpType.IsConfigured() && detailEo.Position.AreaType == 0 {
		err = errors.NewCtx(c, errors.BadRequestErr, "当op_type为op_configed_built-in或op_configed_from-appstore时，position.area_type不能为空", "")
		return
	}

	return
}
