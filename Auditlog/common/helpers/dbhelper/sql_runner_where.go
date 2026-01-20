package dbhelper

import (
	"AuditLog/common/helpers/sqlhelper"
)

func (q *SQLRunner) Where(key string, op sqlhelper.Operator, value interface{}) *SQLRunner {
	q.sb.Where(key, op, value)
	return q
}

func (q *SQLRunner) WhereEqual(key string, value interface{}) *SQLRunner {
	q.sb.Where(key, sqlhelper.OperatorEq, value)
	return q
}

func (q *SQLRunner) WhereNotEqual(key string, value interface{}) *SQLRunner {
	q.sb.Where(key, sqlhelper.OperatorNeq, value)
	return q
}

func (q *SQLRunner) Or(key string, op sqlhelper.Operator, value interface{}) *SQLRunner {
	q.sb.Or(key, op, value)
	return q
}

func (q *SQLRunner) OrEqual(key string, value interface{}) *SQLRunner {
	q.sb.Or(key, sqlhelper.OperatorEq, value)
	return q
}

func (q *SQLRunner) In(key string, value interface{}) *SQLRunner {
	q.sb.In(key, value)
	return q
}

func (q *SQLRunner) Like(key string, value interface{}) *SQLRunner {
	q.sb.Like(key, value)
	return q
}

func (q *SQLRunner) NotIn(key string, value interface{}) *SQLRunner {
	q.sb.NotIn(key, value)
	return q
}

func (q *SQLRunner) WhereRaw(condition string, arg ...interface{}) *SQLRunner {
	q.sb.WhereRaw(condition, arg...)

	return q
}

func (q *SQLRunner) WhereByWhereBuilder(wb *sqlhelper.WhereBuilder) (err error) {
	condition, args, err := wb.ToWhereSQL()
	if err != nil {
		return
	}

	q.sb.WhereRaw("("+condition+")", args...)

	return
}
