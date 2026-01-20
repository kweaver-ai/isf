package persrecpo

type SvcConfigPo struct {
	ID        int64  `json:"id" db:"f_id"`
	Key       string `json:"key" db:"f_key"`
	Value     string `json:"value" db:"f_value"`
	CreatedAt int64  `json:"created_at" db:"f_created_at"`
	UpdatedAt int64  `json:"updated_at" db:"f_updated_at"`
}

func (p *SvcConfigPo) TableName() string {
	return "t_pers_rec_svc_config"
}
