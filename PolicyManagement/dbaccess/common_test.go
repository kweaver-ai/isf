package dbaccess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFindInSetSQL(t *testing.T) {
	sql, args := GetFindInSetSQL([]string{"1", "2", "3"})
	assert.Equal(t, sql, "?,?,?")
	assert.Equal(t, args, []interface{}{"1", "2", "3"})
}
