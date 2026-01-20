package sqlhelper

import (
	"strings"
)

type AndClauses []*Clause

// Build 生成sql语句
// 示例：
//
//	clause := AndClauses{
//	 &Clause{
//	   Key:      "id",
//	   Value:    []int{1, 2, 3},
//	   Operator: OperatorIn,
//	 },
//	 &Clause{
//	 Key:      "updated_at",
//	 Value:    "2022-01-01",
//	 Operator: OperatorLtEq,
//	 },
//	}
//
// sqlStr, args := clause.Build()
// sqlStr: "id in (?,?,?) and updated_at <= ?"
// args: []interface{}{1, 2, 3, "2022-01-01"}
func (a *AndClauses) Build() (sqlStr string, args []interface{}, err error) {
	args = make([]interface{}, 0)

	if len(*a) == 0 {
		return
	}

	clauseStrs := make([]string, 0, len(*a))

	for _, clause := range *a {
		_sqlStr, _args, _err := clause.Build()
		if _err != nil {
			err = _err
			return
		}

		clauseStrs = append(clauseStrs, _sqlStr)
		args = append(args, _args...)
	}

	sqlStr = strings.Join(clauseStrs, " and ")

	return
}

type OrClauses []*Clause

// Build 生成sql语句
// 示例：
//
//	clause := OrClauses{
//	 &Clause{
//	   Key:      "id",
//	   Value:    []int{1, 2, 3},
//	   Operator: OperatorIn,
//	 },
//	 &Clause{
//	 Key:      "updated_at",
//	 Value:    "2022-01-01",
//	 Operator: OperatorLtEq,
//	 },
//	}
//
// sqlStr, args := clause.Build()
// sqlStr: "id in (?,?,?) or updated_at <= ?"
// args: []interface{}{1, 2, 3, "2022-01-01"}
func (o *OrClauses) Build() (sqlStr string, args []interface{}, err error) {
	args = make([]interface{}, 0)

	if len(*o) == 0 {
		return
	}

	clauseStrs := make([]string, 0, len(*o))

	for _, clause := range *o {
		_sqlStr, _args, _err := clause.Build()
		if _err != nil {
			err = _err
			return
		}

		clauseStrs = append(clauseStrs, _sqlStr)
		args = append(args, _args...)
	}

	sqlStr = strings.Join(clauseStrs, " or ")

	return
}

type WhereBuilder struct {
	Raw        []string
	RawArgs    []interface{}
	AndClauses AndClauses
	OrClauses  OrClauses
	// AndGroups  []IClause
	// OrGroups   []IClause
}

func NewWhereBuilder() *WhereBuilder {
	return &WhereBuilder{
		AndClauses: AndClauses{},
		OrClauses:  OrClauses{},
		// AndGroups: make([]IClause, 0),
		// OrGroups:   make([]IClause, 0),
	}
}

func (w *WhereBuilder) Where(key string, op Operator, value interface{}) *WhereBuilder {
	if !op.Check() {
		panic("operator not support")
	}

	clause := &Clause{
		Key:      key,
		Value:    value,
		Operator: op,
	}

	w.AndClauses = append(w.AndClauses, clause)

	return w
}

func (w *WhereBuilder) WhereEqual(key string, value interface{}) *WhereBuilder {
	return w.Where(key, OperatorEq, value)
}

func (w *WhereBuilder) WhereNotEqual(key string, value interface{}) *WhereBuilder {
	return w.Where(key, OperatorNeq, value)
}

func (w *WhereBuilder) Like(key string, value interface{}) *WhereBuilder {
	return w.Where(key, OperatorLike, value)
}

func (w *WhereBuilder) Or(key string, op Operator, value interface{}) *WhereBuilder {
	if !op.Check() {
		panic("operator not support")
	}

	clause := &Clause{
		Key:      key,
		Value:    value,
		Operator: op,
	}

	w.OrClauses = append(w.OrClauses, clause)

	return w
}

func (w *WhereBuilder) OrEqual(key string, value interface{}) *WhereBuilder {
	return w.Or(key, OperatorEq, value)
}

func (w *WhereBuilder) OrNotEqual(key string, value interface{}) *WhereBuilder {
	return w.Or(key, OperatorNeq, value)
}

func (w *WhereBuilder) OrLike(key string, value interface{}) *WhereBuilder {
	return w.Or(key, OperatorLike, value)
}

func (w *WhereBuilder) In(key string, value interface{}) *WhereBuilder {
	clause := &Clause{
		Key:      key,
		Value:    value,
		Operator: OperatorIn,
	}

	w.AndClauses = append(w.AndClauses, clause)

	return w
}

func (w *WhereBuilder) NotIn(key string, value interface{}) *WhereBuilder {
	clause := &Clause{
		Key:      key,
		Value:    value,
		Operator: OperatorNotIn,
	}

	w.AndClauses = append(w.AndClauses, clause)

	return w
}

func (w *WhereBuilder) WhereRaw(condition string, arg ...interface{}) *WhereBuilder {
	w.Raw = append(w.Raw, condition)
	w.RawArgs = append(w.RawArgs, arg...)

	return w
}

// ToWhereSQL 生成sql语句
// 子句优先级：and子句 > and分组 > or子句 > or分组
// and子句和and分组之间是and关系
// and分组 和 or子句 之间是or关系
// or子句和or分组之间是or关系
// 示例：见common/sqlhelper/where_builder_test.go
func (w *WhereBuilder) ToWhereSQL() (whereStr string, args []interface{}, err error) {
	args = make([]interface{}, 0)

	// and子句
	andSQLStr, andArgs, err := w.AndClauses.Build()
	if err != nil {
		return
	}

	// andsqlStr切片，用于存放and子句和and分组
	andSQLStrSlice := make([]string, 0)

	// 1.1: 先拼接and子句，加入到andSqlStrSlice
	if andSQLStr != "" {
		andSQLStrSlice = append(andSQLStrSlice, andSQLStr)
		args = append(args, andArgs...)
	}

	// 将and子句和and分组拼接起来
	if len(andSQLStrSlice) > 0 {
		whereStr = strings.Join(andSQLStrSlice, " and ")
	}

	// 2.1: 再拼接or子句
	// or子句
	orSQLStr, orArgs, err := w.OrClauses.Build()
	if err != nil {
		return
	}

	if orSQLStr != "" {
		if whereStr != "" {
			whereStr = whereStr + " or " + orSQLStr
		} else {
			whereStr = orSQLStr
		}

		args = append(args, orArgs...)
	}

	// 3.1: 再拼接原生sql
	if len(w.Raw) > 0 {
		if whereStr != "" {
			whereStr = whereStr + " and " + strings.Join(w.Raw, " and ")
		} else {
			whereStr = strings.Join(w.Raw, " and ")
		}

		args = append(args, w.RawArgs...)
	}

	return
}
