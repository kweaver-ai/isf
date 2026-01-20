package sqlhelper

type IClause interface {
	Build() (sqlStr string, args []interface{}, err error)
}
