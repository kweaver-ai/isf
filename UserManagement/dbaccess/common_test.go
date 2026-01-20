// Package delcommon Anyshare 数据访问层 -通用
package dbaccess

import (
	"testing"

	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"UserManagement/interfaces"
)

func TestGetFindInSetSQL(t *testing.T) {
	Convey("GetFindInSetSQL", t, func() {
		strName := []string{"xxx", "yyy", "zzzz"}
		sql, args := GetFindInSetSQL(strName)

		interName := []interface{}{"xxx", "yyy", "zzzz"}
		assert.Equal(t, sql, "?,?,?")
		assert.Equal(t, args, interName)
	})
}

func TestSplitArray(t *testing.T) {
	Convey("SplitArray", t, func() {
		arr := []string{}
		i := 0
		for i < 1000 {
			arr = append(arr, "xxzz")
			i++
		}
		arr = append(arr, "xxx")

		results := SplitArray(arr)
		assert.Equal(t, len(results), 3)
		assert.Equal(t, results[2][0], "xxx")
	})
}

func TestHandlerUserDBData(t *testing.T) {
	Convey("HandlerUserDBData", t, func() {
		dbData := userDBData{
			ID:                "xxx",
			Name:              "zzz",
			Account:           "asdad",
			CSFLevel:          1,
			Priority:          2,
			DisableStatus:     1,
			AutoDisableStatus: 1,
			Email:             "sdada",
			AuthType:          2,
			Password:          "dasds",
			DesPassword:       "dasdda",
			NtlmPassword:      "123123",
			Frozen:            1,
			Authenticated:     0,
			PWDControl:        true,
		}

		out := handlerUserDBData(&dbData)
		assert.Equal(t, out.ID, dbData.ID)
		assert.Equal(t, out.Name, dbData.Name)
		assert.Equal(t, out.Account, dbData.Account)
		assert.Equal(t, out.CSFLevel, dbData.CSFLevel)
		assert.Equal(t, out.Priority, dbData.Priority)
		assert.Equal(t, out.DisableStatus, interfaces.Disabled)
		assert.Equal(t, out.AutoDisableStatus, interfaces.ADisabled)
		assert.Equal(t, out.Email, dbData.Email)
		assert.Equal(t, out.AuthType, interfaces.Domain)
		assert.Equal(t, out.Password, dbData.Password)
		assert.Equal(t, out.DesPassword, dbData.DesPassword)
		assert.Equal(t, out.NtlmPassword, dbData.NtlmPassword)
		assert.Equal(t, out.Frozen, true)
		assert.Equal(t, out.Authenticated, false)
		assert.Equal(t, out.PWDControl, true)
	})
}
