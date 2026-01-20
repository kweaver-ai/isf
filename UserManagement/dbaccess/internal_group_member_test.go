package dbaccess

import (
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"UserManagement/common"
	"UserManagement/interfaces"
)

func TestNewInternalGroupMember(t *testing.T) {
	Convey("NewInternalGroupMember", t, func() {
		data := NewInternalGroupMember()
		assert.NotEqual(t, data, nil)
	})
}

func TestInternalGroupMemberAdd(t *testing.T) {
	Convey("Add, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		internalGroupMember := &internalGroupMember{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("infos len = 0 ", func() {
			httpErr := internalGroupMember.Add(strID, nil, nil)
			assert.Equal(t, httpErr, nil)
		})

		testErr := errors.New("xxx")
		var info interfaces.InternalGroupMember
		infos := []interfaces.InternalGroupMember{info}
		Convey("exec fail", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(testErr)

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			httpErr := internalGroupMember.Add(strID, infos, tx)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("exec success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			httpErr := internalGroupMember.Add(strID, infos, tx)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestInternalGroupMemberDelete(t *testing.T) {
	Convey("Delete, db is available", t, func() {
		internalGroupMember := &internalGroupMember{
			db:     nil,
			logger: common.NewLogger(),
		}
		Convey("len = 0", func() {
			httpErr := internalGroupMember.DeleteAll([]string{}, nil)
			assert.Equal(t, httpErr, nil)
		})

		testErr := errors.New("xxx")
		Convey("exec fail", func() {
			db1, mock1, err := sqlx.New()
			assert.Equal(t, err, nil)

			mock1.ExpectBegin()
			mock1.ExpectExec("").WillReturnError(testErr)

			tx, err := db1.Begin()
			assert.Equal(t, err, nil)

			httpErr := internalGroupMember.DeleteAll([]string{strID}, tx)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("exec success", func() {
			db1, mock1, err := sqlx.New()
			assert.Equal(t, err, nil)

			mock1.ExpectBegin()
			mock1.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db1.Begin()
			assert.Equal(t, err, nil)

			httpErr := internalGroupMember.DeleteAll([]string{strID}, tx)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestInternalGroupMemberGet(t *testing.T) {
	Convey("Get, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		internalGroupMember := &internalGroupMember{
			db:     db,
			logger: common.NewLogger(),
		}

		testErr := errors.New("xxx")
		Convey("query fail", func() {
			mock.ExpectQuery("").WillReturnError(testErr)

			_, httpErr := internalGroupMember.Get(strID)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("query success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_member_id", "f_member_type"}).AddRow(strID, 1))

			out, httpErr := internalGroupMember.Get(strID)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(out), 1)
			assert.Equal(t, out[0].ID, strID)
			assert.Equal(t, out[0].Type, interfaces.User)
		})
	})
}

func TestInternalGroupMemberGetBelongGroupsByID(t *testing.T) {
	Convey("GetBelongGroups, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		internalGroupMember := &internalGroupMember{
			db:     db,
			logger: common.NewLogger(),
		}

		testErr := errors.New("xxx")
		var info interfaces.InternalGroupMember
		Convey("query fail", func() {
			mock.ExpectQuery("").WillReturnError(testErr)

			_, httpErr := internalGroupMember.GetBelongGroups(info)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("query success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_internal_group_id"}).AddRow(strID))

			outInfos, httpErr := internalGroupMember.GetBelongGroups(info)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(outInfos), 1)
			assert.Equal(t, outInfos[0], strID)
		})
	})
}
