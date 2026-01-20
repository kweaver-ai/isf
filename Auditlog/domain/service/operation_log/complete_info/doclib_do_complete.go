package complete_info

import (
	"path/filepath"
	"strings"

	"AuditLog/common"
	"AuditLog/common/enums/oprlogenums"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
	"AuditLog/domain/entity/oprlogeo/dveo"
	"AuditLog/domain/entity/oprlogeo/mbceo"
	"AuditLog/drivenadapters/httpaccess/document/docaccret"
)

func (l *CompleteInfo) doCompleteDocLib(eos []*oprlogeo.LogEntry, bizType oprlogenums.BizType, pathMap, objectDocIDsMap map[string]string, docLibInfoMap map[string]*docaccret.DocLibItem) (err error) {
	for i := range eos {
		// 1 填充object信息
		err = l.completeObject(eos[i], bizType, pathMap, objectDocIDsMap, docLibInfoMap)
		if err != nil {
			return
		}

		// 2 填充detail信息
		switch bizType {
		case oprlogenums.DirVisit:
			err = l.completeDvDetail(eos[i], pathMap, objectDocIDsMap, docLibInfoMap)
		case oprlogenums.MenuButtonClick:
			err = l.completeMbcDetail(eos[i], pathMap, objectDocIDsMap)
		}

		if err != nil {
			return
		}
	}

	return
}

// completeObject 填充object信息
func (l *CompleteInfo) completeObject(eo *oprlogeo.LogEntry, bizType oprlogenums.BizType, pathMap map[string]string, objectDocIDsMap map[string]string, docLibInfoMap map[string]*docaccret.DocLibItem) (err error) {
	object := eo.Object

	if object == nil {
		return
	}

	b := l.isDocLibObject(bizType, object)

	if !b {
		return
	}

	if !strings.HasPrefix(object.ID, "gns://") {
		if docID, ok := objectDocIDsMap[object.ID]; ok {
			object.ID = docID
		}
	}

	// 1 填充object path信息
	if _, ok := pathMap[object.ID]; ok {
		object.Path = pathMap[object.ID]
		object.Name = filepath.Base(object.Path)
		// 替换description中的中占位符{{$object_path}}为pathMap[object.ID]
		eo.Description = strings.Replace(eo.Description, "{{$object_path}}", pathMap[object.ID], -1)
	}

	docLibID := common.GetDocLibIDByDocID(object.ID)

	if object.DocLib == nil {
		object.DocLib = &oprlogeo.DocLib{}
		object.DocLib.ID = docLibID
	}

	// 2 填充object 文档库信息
	if info, ok := docLibInfoMap[docLibID]; ok {
		object.DocLib.Name = info.Name
		object.DocLib.Type = info.Type
	}

	return
}

// completeDvDetail 填充dir_visit的detail信息
func (l *CompleteInfo) completeDvDetail(eo *oprlogeo.LogEntry, pathMap map[string]string, objectDocIDsMap map[string]string, docLibInfoMap map[string]*docaccret.DocLibItem) (err error) {
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
		if !strings.HasPrefix(fo.ID, "gns://") {
			if docID, ok := objectDocIDsMap[fo.ID]; ok {
				fo.ID = docID
			}
		}

		docLibID := common.GetDocLibIDByDocID(fo.ID)

		if _, ok := pathMap[fo.ID]; ok {
			fo.Path = pathMap[fo.ID]
			// 替换description中的中占位符{{$detail_from_object_path}}为pathMap[object.ID]
			eo.Description = strings.Replace(eo.Description, "{{$detail_from_object_path}}", pathMap[fo.ID], -1)
		}

		if fo.DocLib == nil {
			fo.DocLib = &oprlogeo.DocLib{}
			fo.DocLib.ID = docLibID
		}

		if info, ok := docLibInfoMap[docLibID]; ok {
			fo.DocLib.Name = info.Name
			fo.DocLib.Type = info.Type
		}
	}

	eo.Detail = detail

	return
}

// completeMbcDetail 填充menu_button_click的detail信息
func (l *CompleteInfo) completeMbcDetail(eo *oprlogeo.LogEntry, pathMap, objectDocIDsMap map[string]string) (err error) {
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
		if strings.HasPrefix(position.PathID, "gns://") {
			if docID, ok := objectDocIDsMap[position.PathID]; ok {
				position.PathID = docID
			}
		}

		if _, ok := pathMap[position.PathID]; ok {
			position.PathName = pathMap[position.PathID]
		}
	}

	eo.Detail = detail

	return
}
