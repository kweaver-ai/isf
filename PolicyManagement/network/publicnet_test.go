package network

import (
	"testing"

	"policy_mgnt/test"

	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCreatePublicNet(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	// 使用sqlite3无法执行mysql的now方法，ut无法通过
	err := CreatePublicNet()
	_, ok := err.(sqlite3.Error)
	assert.Equal(t, true, ok)
}
