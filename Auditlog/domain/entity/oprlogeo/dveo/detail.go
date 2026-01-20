package dveo

import (
	"AuditLog/common/enums/oprlogenums/dvenums"
	"AuditLog/common/utils"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
)

type Detail struct {
	FromObject *FromObject `json:"from_object,omitempty"`
	Object     *Object     `json:"object,omitempty"`
}

func NewDetail() *Detail {
	return &Detail{}
}

func (d *Detail) LoadByInterface(i interface{}) (err error) {
	if i == nil {
		return
	}

	//    通过json来实现
	jsonStr, err := utils.JSON().Marshal(i)
	if err != nil {
		return
	}

	err = utils.JSON().Unmarshal(jsonStr, d)
	if err != nil {
		return
	}

	return
}

type FromObject struct {
	ID     string                 `json:"id"`
	Path   string                 `json:"path"`
	Type   dvenums.FromObjectType `json:"type"`
	DocLib *oprlogeo.DocLib       `json:"doc_lib,omitempty"`
}

type Object struct {
	FavCategory *FavCategory `json:"fav_category,omitempty"`
}

type FavCategory struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IDPath   string `json:"id_path"`
	NamePath string `json:"name_path"`
}
