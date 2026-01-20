// Package dbaccess avatar AnyShare 用户头像数据层
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

const (
	userID = "user_id"
)

func TestNewAvatar(t *testing.T) {
	Convey("Get, db is available", t, func() {
		data := NewAvatar()
		assert.NotEqual(t, data, nil)
	})
}

func TestGet(t *testing.T) {
	Convey("Get, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ava := avatar{
			db:     db,
			logger: common.NewLogger(),
		}

		fields := []string{"f_oss_id", "f_key", "f_type", "f_time"}
		Convey("unassigned users", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			infos, httpErr := ava.Get(userID)
			assert.Equal(t, infos.Key, "")
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			tempInfo := interfaces.AvatarOSSInfo{
				OSSID:   "ossid1",
				Key:     "key1",
				Type:    "png",
				Time:    1223,
				BUseful: true,
				UserID:  userID,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(tempInfo.OSSID, tempInfo.Key, "png", tempInfo.Time))
			info, httpErr := ava.Get(userID)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, info, tempInfo)
		})
	})
}

func TestAdd(t *testing.T) {
	Convey("Add, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ava := avatar{
			db:     db,
			logger: common.NewLogger(),
		}

		info := interfaces.AvatarOSSInfo{}
		Convey("unassigned users", func() {
			testErr := errors.New("xxxx")
			mock.ExpectExec("").WillReturnError(testErr)
			httpErr := ava.Add(&info)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			httpErr := ava.Add(&info)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestUpdateStatus(t *testing.T) {
	Convey("UpdateStatusByKey, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ava := avatar{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("unassigned users", func() {
			testErr := errors.New("xxxx")
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(testErr)
			mock.ExpectCommit()
			nTx, err := db.Begin()
			assert.Equal(t, err, nil)

			httpErr := ava.UpdateStatusByKey("xx", true, nTx)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			nTx, err := db.Begin()
			assert.Equal(t, err, nil)

			httpErr := ava.UpdateStatusByKey("xx", true, nTx)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestSetAvatarUnableByID(t *testing.T) {
	Convey("UpdateStatusSetAvatarUnableByIDByKey, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ava := avatar{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("unassigned users", func() {
			testErr := errors.New("xxxx")
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(testErr)
			mock.ExpectCommit()
			nTx, err := db.Begin()
			assert.Equal(t, err, nil)

			httpErr := ava.SetAvatarUnableByID("xx", nTx)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			nTx, err := db.Begin()
			assert.Equal(t, err, nil)

			httpErr := ava.SetAvatarUnableByID("xx", nTx)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestDelete(t *testing.T) {
	Convey("Delete, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ava := avatar{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("key is empty", func() {
			httpErr := ava.Delete("")
			assert.Equal(t, httpErr, nil)
		})

		Convey("unassigned users", func() {
			testErr := errors.New("xxxx")
			mock.ExpectExec("").WillReturnError(testErr)
			httpErr := ava.Delete("xxx")
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			httpErr := ava.Delete("xxx")
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestGetUselessAvatar(t *testing.T) {
	Convey("GetUselessAvatar, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ava := avatar{
			db:     db,
			logger: common.NewLogger(),
		}

		userID := "user_di1"
		fields := []string{"f_user_id", "f_oss_id", "f_key", "f_type", "f_status", "f_time"}
		Convey("unassigned users", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			infos, httpErr := ava.GetUselessAvatar(12)
			assert.Equal(t, len(infos), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			tempInfo := interfaces.AvatarOSSInfo{
				OSSID:   "ossid1",
				Key:     "key1",
				Type:    "png",
				Time:    1223,
				BUseful: true,
				UserID:  userID,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(tempInfo.UserID, tempInfo.OSSID, tempInfo.Key, "png", 1, tempInfo.Time))
			info, httpErr := ava.GetUselessAvatar(11)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(info), 1)
			assert.Equal(t, info[0], tempInfo)
		})
	})
}
