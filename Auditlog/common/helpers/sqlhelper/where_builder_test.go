package sqlhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhereBuilder_ToWhereSql(t *testing.T) {
	wb := NewWhereBuilder()
	wb.Where("name", OperatorEq, "John")
	wb.Where("city", OperatorIn, []string{"New York", "Los Angeles"})
	wb.Or("age", OperatorGt, 30)
	wb.Or("country", OperatorIn, []string{"USA", "Canada"})

	wb.In("city2", []string{"New York2", "Los Angeles2"})

	sqlStr, args, err := wb.ToWhereSQL()
	assert.Equal(t, nil, err)

	assert.Equal(t, "name = ? and city in (?,?) and city2 in (?,?) or age > ? or country in (?,?)", sqlStr)

	assert.Equal(t, []interface{}{"John", "New York", "Los Angeles", "New York2", "Los Angeles2", 30, "USA", "Canada"}, args)
}

func TestWhereBuilder_WhereEqual_WhereNotEqual(t *testing.T) {
	wb := NewWhereBuilder()
	wb.WhereEqual("name", "John")
	wb.WhereNotEqual("country", "USA")

	sqlStr, args, err := wb.ToWhereSQL()
	assert.Equal(t, nil, err)
	assert.Equal(t, "name = ? and country <> ?", sqlStr)
	assert.Equal(t, []interface{}{"John", "USA"}, args)
}

func TestWhereBuilder_WhereRaw(t *testing.T) {
	wb := NewWhereBuilder()
	wb.WhereRaw("name = ? and country <> ?", "John", "USA")

	sqlStr, args, err := wb.ToWhereSQL()
	assert.Equal(t, nil, err)
	assert.Equal(t, "name = ? and country <> ?", sqlStr)
	assert.Equal(t, []interface{}{"John", "USA"}, args)
}

func TestWhereBuilder_WhereRaw2(t *testing.T) {
	wb := NewWhereBuilder()
	wb.Where("age", OperatorEq, 30)
	wb.WhereRaw("name = ? and country <> ?", "John", "USA")

	sqlStr, args, err := wb.ToWhereSQL()
	assert.Equal(t, nil, err)
	assert.Equal(t, "age = ? and name = ? and country <> ?", sqlStr)
	assert.Equal(t, []interface{}{30, "John", "USA"}, args)
}
