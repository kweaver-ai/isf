// Package logics group AnyShare 内部组业务逻辑层
package logics

import (
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-playground/assert"
	"github.com/satori/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func TestNewInternalGroup(t *testing.T) {
	Convey("NewInternalGroup", t, func() {
		sqlDB, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		dbPool = sqlDB

		data := NewInternalGroup()
		assert.NotEqual(t, data, nil)
	})
}

func TestInternalGroupAddGroup(t *testing.T) {
	Convey("AddGroup", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		groupDB := mock.NewMockDBInternalGroup(ctrl)
		groupMemberDB := mock.NewMockDBInternalGroupMember(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		group := &internalGroup{
			groupMemberDB: groupMemberDB,
			groupDB:       groupDB,
			userDB:        userDB,
			logger:        common.NewLogger(),
			messageBroker: msg,
		}

		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey("groupDB Add err", func() {
			groupDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(testErr)
			_, err := group.AddGroup()

			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			groupDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(nil)
			out, err := group.AddGroup()

			assert.Equal(t, err, nil)
			_, err = uuid.FromString(out)
			assert.Equal(t, err, nil)
		})
	})
}

func TestInternalGroupDeleteGroup(t *testing.T) {
	Convey("DeleteGroup", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		groupDB := mock.NewMockDBInternalGroup(ctrl)
		groupMemberDB := mock.NewMockDBInternalGroupMember(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		group := &internalGroup{
			groupMemberDB: groupMemberDB,
			groupDB:       groupDB,
			userDB:        userDB,
			logger:        common.NewLogger(),
			messageBroker: msg,
			ob:            ob,
			pool:          dPool,
		}

		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey("groupDB get err", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := group.DeleteGroup([]string{strID})

			assert.Equal(t, err, testErr)
		})

		temp1 := make(map[string]interfaces.InternelGroup)
		Convey("exist len = 0,return nil", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(temp1, nil)

			err := group.DeleteGroup([]string{strID})

			assert.Equal(t, err, nil)
		})

		temp1[strID] = interfaces.InternelGroup{}
		Convey("pool begin err", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(temp1, nil)
			txMock.ExpectBegin().WillReturnError(testErr)
			err := group.DeleteGroup([]string{strID})

			assert.Equal(t, err, testErr)
		})

		Convey("groupMemberDB DeleteAll err", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(temp1, nil)
			txMock.ExpectBegin()
			groupMemberDB.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()
			err := group.DeleteGroup([]string{strID})

			assert.Equal(t, err, testErr)
		})

		Convey("groupDB Delete err", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(temp1, nil)
			txMock.ExpectBegin()
			groupMemberDB.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupDB.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()
			err := group.DeleteGroup([]string{strID})

			assert.Equal(t, err, testErr)
		})

		Convey("AddOutboxInfo err", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(temp1, nil)
			txMock.ExpectBegin()
			groupMemberDB.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupDB.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback().WillReturnError(testErr)
			err := group.DeleteGroup([]string{strID})

			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(temp1, nil)
			txMock.ExpectBegin()
			groupMemberDB.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupDB.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread()
			err := group.DeleteGroup([]string{strID})

			assert.Equal(t, err, nil)
		})
	})
}

func TestInternalGroupGetGroupMemberByID(t *testing.T) {
	Convey("GetGroupMemberByID", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		groupDB := mock.NewMockDBInternalGroup(ctrl)
		groupMemberDB := mock.NewMockDBInternalGroupMember(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		group := &internalGroup{
			groupMemberDB: groupMemberDB,
			groupDB:       groupDB,
			userDB:        userDB,
			logger:        common.NewLogger(),
			messageBroker: msg,
			ob:            ob,
		}

		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey("groupDB get err", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := group.GetGroupMemberByID(strID)

			assert.Equal(t, err, testErr)
		})

		data := make(map[string]interfaces.InternelGroup)
		Convey("internal group do not exist", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(data, nil)
			_, err := group.GetGroupMemberByID(strID)

			assert.Equal(t, err, rest.NewHTTPError("internal group do not exist", rest.URINotExist, nil))
		})

		var temp interfaces.InternelGroup
		data[strID] = temp
		var outInfos []interfaces.InternalGroupMember
		Convey("success", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(data, nil)
			groupMemberDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(outInfos, nil)
			out, err := group.GetGroupMemberByID(strID)

			assert.Equal(t, out, outInfos)
			assert.Equal(t, err, nil)
		})
	})
}

func TestInternalGroupUpdateMembers(t *testing.T) {
	Convey("UpdateMembers", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		groupDB := mock.NewMockDBInternalGroup(ctrl)
		groupMemberDB := mock.NewMockDBInternalGroupMember(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		group := &internalGroup{
			groupMemberDB: groupMemberDB,
			groupDB:       groupDB,
			userDB:        userDB,
			logger:        common.NewLogger(),
			messageBroker: msg,
			ob:            ob,
			pool:          dPool,
		}

		testErr := rest.NewHTTPError("error", 503000000, nil)
		infos := make(map[string]interfaces.InternelGroup)
		members := make([]interfaces.InternalGroupMember, 0)
		Convey("groupDB get err", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, testErr)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, testErr)
		})

		Convey("internal group not exist", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, rest.NewHTTPError("internal group do not exist", rest.URINotExist, nil))
		})

		temp1 := interfaces.InternalGroupMember{
			ID: strID,
		}
		temp2 := interfaces.InternalGroupMember{
			ID: strID,
		}
		members = append(members, temp1, temp2)
		infos[strID] = interfaces.InternelGroup{}
		Convey("internal group members not uniqued", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, rest.NewHTTPError("internal group member do not unique", rest.BadRequest, nil))
		})

		members = append([]interfaces.InternalGroupMember{}, temp1)
		Convey("internal group members check user exist error", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(nil, nil, testErr)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, testErr)
		})

		Convey("internal group members user not exist", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(nil, nil, nil)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, rest.NewHTTPError("some members are not existing", rest.BadRequest, map[string]interface{}{"ids": []string{strID}}))
		})

		Convey("pool begin error", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(nil, []string{strID}, nil)
			txMock.ExpectBegin().WillReturnError(testErr)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, testErr)
		})

		Convey("groupMemberDB deleteall error", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(nil, []string{strID}, nil)
			txMock.ExpectBegin()
			groupMemberDB.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, testErr)
		})

		Convey("groupMemberDB add error", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(nil, []string{strID}, nil)
			txMock.ExpectBegin()
			groupMemberDB.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupMemberDB.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			groupDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(infos, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(nil, []string{strID}, nil)
			txMock.ExpectBegin()
			groupMemberDB.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupMemberDB.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()

			err := group.UpdateMembers(strID, members)

			assert.Equal(t, err, nil)
		})
	})
}

func TestSendInternalGroupDeletedInfo(t *testing.T) {
	Convey("sendInternalGroupDeletedInfo", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		groupDB := mock.NewMockDBInternalGroup(ctrl)
		groupMemberDB := mock.NewMockDBInternalGroupMember(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		group := &internalGroup{
			groupMemberDB: groupMemberDB,
			groupDB:       groupDB,
			userDB:        userDB,
			logger:        common.NewLogger(),
			messageBroker: msg,
			ob:            ob,
		}

		testErr := rest.NewHTTPError("error", 503000000, nil)
		context := make(map[string]interface{})
		context["ids"] = []interface{}{strID}

		Convey("groupDB get err", func() {
			msg.EXPECT().InternalGroupDeleted(gomock.Any()).AnyTimes().Return(testErr)
			err := group.sendInternalGroupDeletedInfo(context)

			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			msg.EXPECT().InternalGroupDeleted(gomock.Any()).AnyTimes().Return(nil)
			err := group.sendInternalGroupDeletedInfo(context)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetBelongGroupsByID(t *testing.T) {
	Convey("GetBelongGroups", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		groupDB := mock.NewMockDBInternalGroup(ctrl)
		groupMemberDB := mock.NewMockDBInternalGroupMember(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		group := &internalGroup{
			groupMemberDB: groupMemberDB,
			groupDB:       groupDB,
			userDB:        userDB,
			logger:        common.NewLogger(),
			messageBroker: msg,
			ob:            ob,
		}

		var info interfaces.InternalGroupMember
		outInfos := []string{strID}
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey("success", func() {
			groupMemberDB.EXPECT().GetBelongGroups(gomock.Any()).AnyTimes().Return(outInfos, testErr)
			out, err := group.GetBelongGroups(info)

			assert.Equal(t, err, testErr)
			assert.Equal(t, out, outInfos)
		})
	})
}
