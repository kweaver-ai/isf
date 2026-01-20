package persistence

import uniqidenums "AuditLog/common/enums/uniqueid"

type UniqueID struct {
	ID   string                   `json:"id" db:"f_id"`
	Flag uniqidenums.UniqueIDFlag `json:"flag" db:"f_flag"`
}

func (p *UniqueID) TableName() string {
	return "t_pers_rec_unique_id"
}
