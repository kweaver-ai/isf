//nolint:dupl
package dbaccess

import (
	"errors"
	"testing"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
)

func TestAddReservedName(t *testing.T) {
	Convey("添加保留名称", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		rn := &reservedName{
			db:     dbPool,
			logger: common.NewLogger(),
		}

		Convey("添加保留名称失败", func() {
			txMock.ExpectBegin()
			txMock.ExpectExec("insert into").WillReturnError(errors.New("test"))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			err = rn.AddReservedName(interfaces.ReservedNameInfo{}, tx)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("添加保留名称成功", func() {
			txMock.ExpectBegin()
			txMock.ExpectExec("insert into").WillReturnResult(sqlmock.NewResult(1, 1))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			err = rn.AddReservedName(interfaces.ReservedNameInfo{}, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestUpdateReservedName(t *testing.T) {
	Convey("更新保留名称", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		rn := &reservedName{
			db:     dbPool,
			logger: common.NewLogger(),
		}

		Convey("更新保留名称失败", func() {
			txMock.ExpectBegin()
			txMock.ExpectExec("update").WillReturnError(errors.New("test"))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			err = rn.UpdateReservedName(interfaces.ReservedNameInfo{}, tx)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("更新保留名称成功", func() {
			txMock.ExpectBegin()
			txMock.ExpectExec("update").WillReturnResult(sqlmock.NewResult(1, 1))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			err = rn.UpdateReservedName(interfaces.ReservedNameInfo{}, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetReservedNameByID(t *testing.T) {
	Convey("根据ID获取保留名称", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		rn := &reservedName{
			db:     dbPool,
			logger: common.NewLogger(),
		}

		Convey("根据ID获取保留名称失败", func() {
			txMock.ExpectBegin()
			txMock.ExpectQuery("select").WillReturnError(errors.New("test"))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			_, ok, err := rn.GetReservedNameByID("test", tx)
			assert.Equal(t, ok, false)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("根据ID获取保留名称成功", func() {
			txMock.ExpectBegin()
			txMock.ExpectQuery("select").WillReturnRows(sqlmock.NewRows([]string{"f_id", "f_name", "f_create_time", "f_update_time"}).AddRow("test", "test", time.Now().Unix(), time.Now().Unix()))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			_, ok, err := rn.GetReservedNameByID("test", tx)
			assert.Equal(t, ok, true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetReservedNameByName(t *testing.T) {
	Convey("根据名称获取保留名称", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		rn := &reservedName{
			db:     dbPool,
			logger: common.NewLogger(),
		}

		Convey("根据名称获取保留名称失败", func() {
			txMock.ExpectBegin()
			txMock.ExpectQuery("select").WillReturnError(errors.New("test"))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			_, ok, err := rn.GetReservedNameByName("test", tx)
			assert.Equal(t, ok, false)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("根据名称获取保留名称成功", func() {
			txMock.ExpectBegin()
			txMock.ExpectQuery("select").WillReturnRows(sqlmock.NewRows([]string{"f_id", "f_name", "f_create_time", "f_update_time"}).AddRow("test", "test", time.Now().Unix(), time.Now().Unix()))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			_, ok, err := rn.GetReservedNameByName("test", tx)
			assert.Equal(t, ok, true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteReservedName(t *testing.T) {
	Convey("删除保留名称", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		rn := &reservedName{
			db:     dbPool,
			logger: common.NewLogger(),
		}

		Convey("删除保留名称失败", func() {
			txMock.ExpectExec("delete from").WillReturnError(errors.New("test"))
			err = rn.DeleteReservedName("test")
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("删除保留名称成功", func() {
			txMock.ExpectExec("delete from").WillReturnResult(sqlmock.NewResult(1, 1))
			err = rn.DeleteReservedName("test")
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetLock(t *testing.T) {
	Convey("获取锁", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		rn := &reservedName{
			db:     dbPool,
			logger: common.NewLogger(),
		}

		Convey("获取锁失败", func() {
			txMock.ExpectBegin()
			txMock.ExpectExec("select").WillReturnError(errors.New("test"))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			err = rn.GetLock(tx)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("获取锁成功", func() {
			txMock.ExpectBegin()
			txMock.ExpectExec("select").WillReturnResult(sqlmock.NewResult(1, 1))
			tx, err := dbPool.Begin()
			assert.Equal(t, err, nil)
			err = rn.GetLock(tx)
			assert.Equal(t, err, nil)
		})
	})
}
