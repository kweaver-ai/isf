package logics

import (
	"errors"
	"strings"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

//nolint:funlen
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

		dbReservedName := mock.NewMockDBReservedName(ctrl)
		dbUser := mock.NewMockDBUser(ctrl)

		rn := &reservedName{
			pool:           dbPool,
			dbReservedName: dbReservedName,
			dbUser:         dbUser,
			logger:         common.NewLogger(),
		}

		name := interfaces.ReservedNameInfo{
			ID:   strings.Repeat("a", 32),
			Name: "test",
		}

		existIDInfo := interfaces.ReservedNameInfo{
			ID:   strings.Repeat("b", 32),
			Name: "test2",
		}

		existNameInfo := interfaces.ReservedNameInfo{
			ID:   strings.Repeat("c", 32),
			Name: "test2",
		}
		userInfo := interfaces.UserDBInfo{
			ID:   strings.Repeat("a", 32),
			Name: "test",
		}

		Convey("id非法", func() {
			Convey("id长度小于32", func() {
				name.ID = strings.Repeat("a", 31)
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("id长度大于32", func() {
				name.ID = strings.Repeat("a", 33)
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("id包含非法字符", func() {
				name.ID = strings.Repeat("a", 31) + "*"
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})
		})

		Convey("name非法", func() {
			Convey("name长度小于1", func() {
				name.Name = ""
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("name长度大于128", func() {
				name.Name = strings.Repeat("a", 129)
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("name以空格开头", func() {
				name.Name = " test"
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("name以空格结尾", func() {
				name.Name = "test "
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("name包含特殊字符", func() {
				name.Name = "test*"
				err := rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)

				name.Name = `test|`
				err = rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)

				name.Name = `test\`
				err = rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)

				name.Name = `test/`
				err = rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)

				name.Name = `test?`
				err = rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)

				name.Name = `test"`
				err = rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)

				name.Name = `test<`
				err = rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)

				name.Name = `test>`
				err = rn.UpdateReservedName(name)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})
		})

		Convey("开启事务失败", func() {
			txMock.ExpectBegin().WillReturnError(errors.New("test"))
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("获取锁失败", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetLock(gomock.Any()).Return(errors.New("test"))
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		dbReservedName.EXPECT().GetLock(gomock.Any()).AnyTimes().Return(nil)

		Convey("根据id获取保留名称失败", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).Return(existIDInfo, false, errors.New("test"))
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("根据name获取保留名称失败", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).Return(existIDInfo, false, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).Return(existNameInfo, false, errors.New("test"))
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("根据name获取用户信息失败", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).Return(existIDInfo, false, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).Return(existNameInfo, false, nil)
			dbUser.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).Return(userInfo, errors.New("test"))
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("用户名已存在", func() {
			dbUser.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).Return(existIDInfo, false, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).Return(existNameInfo, false, nil)
			txMock.ExpectBegin()
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			retErr := err.(*rest.HTTPError)
			assert.Equal(t, retErr.Code, rest.Conflict)
			assert.Equal(t, retErr.Detail["conflict_object"].(map[string]interface{})["type"], "user")
		})

		userInfo.ID = ""
		dbUser.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)

		Convey("名称和已有的文档库名重复", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).Return(existIDInfo, false, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).Return(existNameInfo, true, nil)
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			retErr := err.(*rest.HTTPError)
			assert.Equal(t, retErr.Code, rest.Conflict)
			assert.Equal(t, retErr.Detail["conflict_object"].(map[string]interface{})["type"], "doclib")
		})

		Convey("添加保留名称失败", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).AnyTimes().Return(existIDInfo, false, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).AnyTimes().Return(existNameInfo, false, nil)
			dbReservedName.EXPECT().AddReservedName(gomock.Any(), gomock.Any()).Return(errors.New("test"))
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("添加保留名称成功", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).AnyTimes().Return(existIDInfo, false, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).AnyTimes().Return(existNameInfo, false, nil)
			dbReservedName.EXPECT().AddReservedName(gomock.Any(), gomock.Any()).Return(nil)
			txMock.ExpectCommit()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, nil)
		})

		Convey("名称相同，无需修改", func() {
			existIDInfo.Name = name.Name
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).AnyTimes().Return(existIDInfo, true, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).AnyTimes().Return(existNameInfo, false, nil)
			txMock.ExpectCommit()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, nil)
		})

		Convey("更新后的名称和已有的名称重复", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).AnyTimes().Return(existIDInfo, true, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).AnyTimes().Return(existNameInfo, true, nil)
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			retErr := err.(*rest.HTTPError)
			assert.Equal(t, retErr.Code, rest.Conflict)
			assert.Equal(t, retErr.Detail["conflict_object"].(map[string]interface{})["type"], "doclib")
		})

		Convey("更新保留名称失败", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).AnyTimes().Return(existIDInfo, true, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).AnyTimes().Return(existNameInfo, false, nil)
			dbReservedName.EXPECT().UpdateReservedName(gomock.Any(), gomock.Any()).Return(errors.New("test"))
			txMock.ExpectRollback()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("更新保留名称成功", func() {
			txMock.ExpectBegin()
			dbReservedName.EXPECT().GetReservedNameByID(gomock.Any(), gomock.Any()).AnyTimes().Return(existIDInfo, true, nil)
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).AnyTimes().Return(existNameInfo, false, nil)
			dbReservedName.EXPECT().UpdateReservedName(gomock.Any(), gomock.Any()).Return(nil)
			txMock.ExpectCommit()
			err := rn.UpdateReservedName(name)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteReservedName(t *testing.T) {
	Convey("删除保留名称", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		dbReservedName := mock.NewMockDBReservedName(ctrl)
		dbUser := mock.NewMockDBUser(ctrl)

		rn := &reservedName{
			pool:           dbPool,
			dbReservedName: dbReservedName,
			dbUser:         dbUser,
			logger:         common.NewLogger(),
		}

		Convey("id非法", func() {
			Convey("id长度小于32", func() {
				id := strings.Repeat("a", 31)
				err := rn.DeleteReservedName(id)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("id长度大于32", func() {
				id := strings.Repeat("a", 33)
				err := rn.DeleteReservedName(id)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})

			Convey("id包含非法字符", func() {
				id := strings.Repeat("a", 31) + "*"
				err := rn.DeleteReservedName(id)
				assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			})
		})

		Convey("删除保留名称失败", func() {
			id := strings.Repeat("a", 32)
			dbReservedName.EXPECT().DeleteReservedName(gomock.Any()).Return(errors.New("test"))
			err := rn.DeleteReservedName(id)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("删除保留名称成功", func() {
			id := strings.Repeat("a", 32)
			dbReservedName.EXPECT().DeleteReservedName(gomock.Any()).Return(nil)
			err := rn.DeleteReservedName(id)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetReservedName(t *testing.T) {
	Convey("获取保留名称", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbReservedName := mock.NewMockDBReservedName(ctrl)
		dbUser := mock.NewMockDBUser(ctrl)

		rn := &reservedName{
			dbReservedName: dbReservedName,
			dbUser:         dbUser,
			logger:         common.NewLogger(),
		}
		name := "name"

		Convey("获取保留名称失败", func() {
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).Return(interfaces.ReservedNameInfo{}, false, errors.New("test"))
			_, err := rn.GetReservedName(name)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("获取保留名称成功", func() {
			dbReservedName.EXPECT().GetReservedNameByName(gomock.Any(), gomock.Any()).Return(interfaces.ReservedNameInfo{}, true, nil)
			_, err := rn.GetReservedName(name)
			assert.Equal(t, err, nil)
		})
	})
}
