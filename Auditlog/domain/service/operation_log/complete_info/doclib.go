package complete_info

import (
	"context"
	"path/filepath"

	"AuditLog/common"
	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/enums/oprlogenums/dvenums"
	"AuditLog/common/enums/oprlogenums/mbcenums"
	"AuditLog/common/utils"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
	"AuditLog/drivenadapters/httpaccess/document/docaccret"
)

func (l *CompleteInfo) completeDocLibInfo(c context.Context, maps []map[string]interface{}, eos []*oprlogeo.LogEntry, bizType oprlogenums.BizType) (err error) {
	completeBizTypes := []oprlogenums.BizType{
		oprlogenums.DirVisit,
		oprlogenums.MenuButtonClick,
	}

	// 判断是否是需要补全文档库信息的业务类型
	if !utils.ExistsGeneric(completeBizTypes, bizType) {
		return
	}

	// 1. 批量获取文档库信息
	var (
		docLibInfoMap   map[string]*docaccret.DocLibItem
		pathMap         map[string]string
		objectDocIDsMap map[string]string
	)

	docLibInfoMap, pathMap, objectDocIDsMap, err = l.getDocLibInfo(c, eos, bizType)
	if err != nil {
		return
	}

	// 2. 填充文档库、文档path信息 到eos
	err = l.doCompleteDocLib(eos, bizType, pathMap, objectDocIDsMap, docLibInfoMap)
	if err != nil {
		return
	}

	// 3. 合并eos到maps
	err = l.mergeEosToMap(maps, eos)

	return
}

// isDocLibObject 判断是否是文档库对象
func (l *CompleteInfo) isDocLibObject(bizType oprlogenums.BizType, object *oprlogeo.ObjectInfo) bool {
	c1 := bizType == oprlogenums.DirVisit && dvenums.IsDocLibOTString(object.Type)
	c2 := bizType == oprlogenums.MenuButtonClick && mbcenums.IsDocLibOTString(object.Type)

	return c1 || c2
}

// getDocLibInfo 批量获取文档库对象信息 1. 文档库信息 2. path信息（文档库、文档）
func (l *CompleteInfo) getDocLibInfo(c context.Context,
	eos []*oprlogeo.LogEntry, bizType oprlogenums.BizType) (
	docLibInfoMap map[string]*docaccret.DocLibItem, pathMap, objectDocIDsMap map[string]string, err error,
) {
	var (
		docLibInfos []*docaccret.DocLibItem
		docLibIDs   []string
		docIDs      []string
		objectIDs   []string
	)

	// pathMap 初始化
	pathMap = make(map[string]string)

	// 1. 收集doc_id和object_id
	docIDs, objectIDs, err = l.collectDocAndObjIDs(eos, bizType)
	if err != nil {
		return
	}

	// 2. 批量获取path信息
	if len(docIDs) != 0 {
		// 去重
		docIDs = utils.DeduplGeneric(docIDs)

		pathMap, err = l.eFastHttpAcc.Gns2Path(c, docIDs)
		if err != nil {
			return
		}
	}

	// 3. 批量获取object信息
	if len(objectIDs) != 0 {
		// 去重
		objectIDs = utils.DeduplGeneric(objectIDs)

		var objectPathMap map[string]string

		objectPathMap, err = l.eFastHttpAcc.Gns2Path(c, objectIDs)
		if err != nil {
			return
		}

		for docID, path := range objectPathMap {
			pathMap[docID] = path

			docIDs = append(docIDs, docID)
			objectID := filepath.Base(docID)
			objectDocIDsMap = make(map[string]string)
			objectDocIDsMap[objectID] = docID
		}
	}

	// 4. 获取doc_lib_id
	for _, docID := range docIDs {
		docLibID := common.GetDocLibIDByDocID(docID)
		if docLibID != "" {
			docLibIDs = append(docLibIDs, docLibID)
		}
	}

	// 3. 批量获取文档库信息
	if len(docLibIDs) != 0 {
		// 去重
		docLibIDs = utils.DeduplGeneric(docLibIDs)

		docLibInfos, err = l.documentHttpAcc.GetBatchDocLibInfos(docLibIDs)
		if err != nil {
			return
		}
	}

	docLibInfoMap = make(map[string]*docaccret.DocLibItem, len(docLibInfos))
	for _, info := range docLibInfos {
		docLibInfoMap[info.ID] = info
	}

	return
}
