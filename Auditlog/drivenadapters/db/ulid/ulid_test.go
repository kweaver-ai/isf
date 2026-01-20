package dbaulid

import (
	"context"
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	uniqidenums "AuditLog/common/enums/uniqueid"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces/drivenadapter/idbaccess"
	"AuditLog/test/mock_log"
)

func initDb(logger api.Logger) (repoImpl idbaccess.UlidRepo, db *sqlx.DB, sqlMock sqlmock.Sqlmock) {
	var err error

	db, sqlMock, err = sqlx.New()
	if err != nil {
		panic(err)
	}

	repoImpl = &ulidRepo{
		db:     db,
		logger: logger,
	}

	return
}

func TestGenUniqID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	logger := mock_log.NewMockLogger(ctrl)

	repoImpl, _, sqlMock := initDb(logger)

	convey.Convey("GenUniqID", t, func() {
		convey.Convey("成功", func() {
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 1))

			id, err := repoImpl.GenUniqID(ctx, uniqidenums.UniqueIDFlagRedisDlm)
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(id, convey.ShouldNotEqual, "")
		})

		convey.Convey("第一次失败，第一次重试成功", func() {
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnError(errors.New("db error"))

			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 1))

			id, err := repoImpl.GenUniqID(ctx, uniqidenums.UniqueIDFlagRedisDlm)
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(id, convey.ShouldNotEqual, "")
		})

		convey.Convey("第一次失败，所有重试都失败", func() {
			for i := 0; i < 5; i++ {
				sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
					WillReturnError(errors.New("db error"))
			}

			id, err := repoImpl.GenUniqID(ctx, uniqidenums.UniqueIDFlagRedisDlm)
			convey.So(err, convey.ShouldNotEqual, nil)
			convey.So(err.Error(), convey.ShouldEqual, "[GenUniqID]: failed to generate unique id, err: db error")
			convey.So(id, convey.ShouldEqual, "")
		})
	})
}

func TestGenDBID(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockLogger(ctrl)

	repoImpl, db, sqlMock := initDb(logger)

	convey.Convey("GenDBID", t, func() {
		convey.Convey("成功", func() {
			sqlMock.ExpectBegin()
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectCommit()

			tx, err := db.Begin()
			convey.So(err, convey.ShouldEqual, nil)

			id, err := repoImpl.GenDBID(ctx, tx)
			convey.So(err, convey.ShouldEqual, nil)

			err = tx.Commit()
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(id, convey.ShouldNotEqual, "")
		})

		convey.Convey("第一次失败，第一次重试成功", func() {
			sqlMock.ExpectBegin()
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnError(errors.New("db error"))
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectCommit()

			tx, err := db.Begin()
			convey.So(err, convey.ShouldEqual, nil)

			id, err := repoImpl.GenDBID(ctx, tx)
			convey.So(err, convey.ShouldEqual, nil)

			err = tx.Commit()
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(id, convey.ShouldNotEqual, "")
		})

		convey.Convey("第一次失败，所有重试都失败", func() {
			sqlMock.ExpectBegin()

			for i := 0; i < 5; i++ {
				sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
					WillReturnError(errors.New("db error"))
			}
			sqlMock.ExpectRollback()

			tx, err := db.Begin()
			convey.So(err, convey.ShouldEqual, nil)

			id, err := repoImpl.GenDBID(ctx, tx)
			convey.So(err, convey.ShouldNotEqual, nil)
			convey.So(err.Error(), convey.ShouldEqual, "[GenDBID]: failed to generate unique id, err: db error")
			convey.So(id, convey.ShouldEqual, "")

			err = tx.Rollback()
			convey.So(err, convey.ShouldEqual, nil)
		})
	})
}

func TestBatchGenDBID(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockLogger(ctrl)

	repoImpl, db, sqlMock := initDb(logger)

	convey.Convey("BatchGenDBID", t, func() {
		convey.Convey("成功", func() {
			sqlMock.ExpectBegin()
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 500))
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 500))
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").WillReturnResult(sqlmock.NewResult(1, 2))
			sqlMock.ExpectCommit()

			tx, err := db.Begin()
			convey.So(err, convey.ShouldEqual, nil)

			ids, err := repoImpl.BatchGenDBID(ctx, tx, 1002)
			convey.So(err, convey.ShouldEqual, nil)

			err = tx.Commit()
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(len(ids), convey.ShouldEqual, 1002)
		})

		convey.Convey("第一次失败，第一次重试成功", func() {
			sqlMock.ExpectBegin()
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnError(errors.New("db error"))
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 500))
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
				WillReturnResult(sqlmock.NewResult(1, 500))
			sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").WillReturnResult(sqlmock.NewResult(1, 2))
			sqlMock.ExpectCommit()

			tx, err := db.Begin()
			convey.So(err, convey.ShouldEqual, nil)

			ids, err := repoImpl.BatchGenDBID(ctx, tx, 1002)
			convey.So(err, convey.ShouldEqual, nil)

			err = tx.Commit()
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(len(ids), convey.ShouldEqual, 1002)
		})

		convey.Convey("第一次失败，所有重试都失败。分批执行的第一次失败", func() {
			sqlMock.ExpectBegin()

			for i := 0; i < 5; i++ {
				sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
					WillReturnError(errors.New("db error"))
			}
			sqlMock.ExpectRollback()

			tx, err := db.Begin()
			convey.So(err, convey.ShouldEqual, nil)

			ids, err := repoImpl.BatchGenDBID(ctx, tx, 1002)
			convey.So(err, convey.ShouldNotEqual, nil)
			convey.So(err.Error(), convey.ShouldEqual, "[BatchGenDBID]: failed to batch generate unique id, err: db error")
			convey.So(len(ids), convey.ShouldEqual, 0)
			convey.So(ids, convey.ShouldEqual, []string(nil))

			err = tx.Rollback()
			convey.So(err, convey.ShouldEqual, nil)
		})

		convey.Convey("第一次失败，所有重试都失败。分批执行的第二次失败", func() {
			sqlMock.ExpectBegin()

			for i := 0; i < 5; i++ {
				sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
					WillReturnResult(sqlmock.NewResult(1, 500))
				sqlMock.ExpectExec("insert into t_pers_rec_unique_id ").
					WillReturnError(errors.New("db error"))
			}
			sqlMock.ExpectRollback()

			tx, err := db.Begin()
			convey.So(err, convey.ShouldEqual, nil)

			ids, err := repoImpl.BatchGenDBID(ctx, tx, 1002)
			convey.So(err, convey.ShouldNotEqual, nil)
			convey.So(err.Error(), convey.ShouldEqual, "[BatchGenDBID]: failed to batch generate unique id, err: db error")
			convey.So(len(ids), convey.ShouldEqual, 0)
			convey.So(ids, convey.ShouldEqual, []string(nil))

			err = tx.Rollback()
			convey.So(err, convey.ShouldEqual, nil)
		})
	})
}
