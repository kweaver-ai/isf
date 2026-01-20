package complete_info

import (
	"strings"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/domain/entity/oprlogeo"
	"AuditLog/domain/entity/oprlogeo/dveo"
	"AuditLog/domain/entity/oprlogeo/mbceo"
)

// collectDocAndDocLibIDs 收集doc_id和object_id
func (l *CompleteInfo) collectDocAndObjIDs(eos []*oprlogeo.LogEntry, bizType oprlogenums.BizType) (docIDs []string, objectIDs []string, err error) {
	for i := range eos {
		// 1. 收集object的doc_id和object_id
		var objDocIDs, objIDs []string

		objDocIDs, objIDs, err = l.collectObjectDocAndDocLibIDs(eos[i], bizType)
		if err != nil {
			return
		}

		docIDs = append(docIDs, objDocIDs...)
		objectIDs = append(objectIDs, objIDs...)

		// 2. 收集detail中的的文档和文档库id
		var detailDocIDs, detailObjIDs []string

		switch bizType {
		case oprlogenums.DirVisit:
			detailDocIDs, detailObjIDs, err = l.collectDvDetailDocAndDocLibIDs(eos[i])
		case oprlogenums.MenuButtonClick:
			detailDocIDs, detailObjIDs, err = l.collectMbcDetailDocAndDocLibIDs(eos[i])
		}

		if err != nil {
			return
		}

		docIDs = append(docIDs, detailDocIDs...)
		objectIDs = append(objectIDs, detailObjIDs...)
	}

	return
}

// collectObjectDocAndDocLibIDs 收集object的doc_id和object_id
func (l *CompleteInfo) collectObjectDocAndDocLibIDs(eo *oprlogeo.LogEntry, bizType oprlogenums.BizType) (docIDs []string, objectIDs []string, err error) {
	object := eo.Object
	if object == nil {
		return
	}

	b := l.isDocLibObject(bizType, object)
	if !b {
		return
	}

	if object.ID == "" {
		return
	}

	if strings.HasPrefix(object.ID, "gns://") {
		// 1.1 收集object的doc_id
		docIDs = append(docIDs, object.ID)
	} else {
		objectIDs = append(objectIDs, object.ID)
	}

	// 1.2收集object的doc_lib_id
	// 获取object中的id，取第一层文档库id作为doc_lib_id
	// docLibID := common.GetDocLibIDByDocID(object.ID)
	// if docLibID != "" {
	// 	docLibIDs = append(docLibIDs, docLibID)
	// }

	return
}

// collectDvDetailDocAndDocLibIDs 收集dirVisit的detail中文档和文档库id
func (l *CompleteInfo) collectDvDetailDocAndDocLibIDs(eo *oprlogeo.LogEntry) (docIDs []string, objIDs []string, err error) {
	if eo.Detail == nil {
		return
	}

	detail := dveo.NewDetail()

	err = detail.LoadByInterface(eo.Detail)
	if err != nil {
		return
	}

	fo := detail.FromObject
	if fo == nil {
		return
	}

	if fo.Type.IsDocLibFOT() {
		// 1.3 收集from_object的doc_id
		if strings.HasPrefix(fo.ID, "gns://") {
			docIDs = append(docIDs, fo.ID)
		} else {
			objIDs = append(objIDs, fo.ID)
		}
		// 1.4 收集from_object的doc_lib_id
		// docLibID := common.GetDocLibIDByDocID(fo.ID)
		//
		//	if docLibID != "" {
		//		docLibIDs = append(docLibIDs, docLibID)
		//	}
	}

	return
}

// collectMbcDetailDocAndDocLibIDs 收集menuButtonClick的detail中的文档和文档库id
func (l *CompleteInfo) collectMbcDetailDocAndDocLibIDs(eo *oprlogeo.LogEntry) (docIDs []string, objectIDs []string, err error) {
	if eo.Detail == nil {
		return
	}

	detail := mbceo.NewDetail()

	err = detail.LoadByInterface(eo.Detail)
	if err != nil {
		return
	}

	position := detail.Position
	if position == nil {
		return
	}

	if position.IsDocLibPosition() {
		// 收集position的doc_id
		if strings.HasPrefix(position.PathID, "gns://") {
			docIDs = append(docIDs, position.PathID)
		} else {
			objectIDs = append(objectIDs, position.PathID)
		}
	}

	return
}
