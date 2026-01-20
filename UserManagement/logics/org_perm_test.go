package logics

import (
	"context"
	"errors"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"

	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func newOrgPerm(o interfaces.DBOrgPerm) *orgPerm {
	return &orgPerm{
		db: o,
	}
}

func TestNewOrgPerm(t *testing.T) {
	Convey("NewOrgPerm, db is available", t, func() {
		out := NewOrgPerm()
		assert.NotEqual(t, out, nil)
	})
}

func TestUpdateName(t *testing.T) {
	Convey("UpdateName, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mock.NewMockTraceClient(ctrl)

		o := mock.NewMockDBOrgPerm(ctrl)
		perm := newOrgPerm(o)
		perm.trace = trace

		Convey("UpdateName err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			o.EXPECT().UpdateName(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := perm.UpdateName("", "")

			assert.Equal(t, err, testErr)
		})

		Convey("Success err", func() {
			o.EXPECT().UpdateName(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			err := perm.UpdateName("", "")

			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteOrgPerm(t *testing.T) {
	Convey("DeleteOrgPerm, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPerm(ctrl)
		u := mock.NewMockDBUser(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		perm := newOrgPerm(o)
		perm.user = u
		perm.trace = trace
		perm.logger = common.NewLogger()

		objects := []interfaces.OrgType{interfaces.Department, interfaces.User, interfaces.Department}
		userID := strID
		ctx := context.Background()
		testErr := errors.New("xxx111")

		Convey("org type not realname", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			err := perm.DeleteOrgPerm(ctx, userID, interfaces.Anonymous, objects)
			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "subject type is not supported"))
		})

		objects = []interfaces.OrgType{interfaces.Department, interfaces.Department}
		Convey("org type are not uniqued", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			err := perm.DeleteOrgPerm(ctx, userID, interfaces.RealName, objects)
			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "org type is not uniqued"))
		})

		objects = []interfaces.OrgType{interfaces.Department, interfaces.User}
		Convey("DeleteOrgPermInfo error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			o.EXPECT().DeleteOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)

			err := perm.DeleteOrgPerm(ctx, userID, interfaces.RealName, objects)
			assert.Equal(t, err, testErr)
		})
	})
}

func TestSetOrgPerm(t *testing.T) {
	Convey("SetOrgPerm, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dPool, poolMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		o := mock.NewMockDBOrgPerm(ctrl)
		u := mock.NewMockDBUser(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		perm := newOrgPerm(o)
		perm.user = u
		perm.trace = trace
		perm.logger = common.NewLogger()
		perm.pool = dPool

		ctx := context.Background()
		testErr := errors.New("xxx111")

		permInfo := interfaces.OrgPerm{}
		Convey("visitor has not authority", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			err := perm.SetOrgPerm(ctx, strID, interfaces.Anonymous, []interfaces.OrgPerm{permInfo})
			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "subject type is not supported"))
		})

		Convey("GetUserDBInfo2 error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{permInfo})
			assert.Equal(t, err, testErr)
		})

		Convey("GetUserDBInfo2 no user", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)

			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{permInfo})
			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.URINotExist, "user not found"))
		})

		var userInfo1 interfaces.UserDBInfo
		Convey("GetPermByID error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{userInfo1}, nil)
			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{permInfo})
			assert.Equal(t, err, testErr)
		})

		permInfo.SubjectID = "xxxxxxx1"
		permInfo.Object = interfaces.User
		permInfo.Value = 1
		testObejcts := interfaces.OrgPerm{}
		testObejcts.SubjectID = strID
		testObejcts.Object = interfaces.User
		testObejcts.Value = 2
		objects := make(map[interfaces.OrgType]interfaces.OrgPerm, 0)
		objects[interfaces.User] = testObejcts
		Convey("handlePermInfo error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{userInfo1}, nil)
			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(objects, nil)

			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{permInfo})
			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "subject is not same as subject id"))
		})

		permInfo.SubjectID = strID
		Convey("pool begin error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{userInfo1}, nil)
			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(objects, nil)

			poolMock.ExpectBegin().WillReturnError(testErr)
			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{permInfo})
			assert.Equal(t, err, testErr)
		})

		Convey("UpdateOrgPerm error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{userInfo1}, nil)
			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(objects, nil)

			poolMock.ExpectBegin()
			o.EXPECT().UpdateOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			poolMock.ExpectRollback()
			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{permInfo})
			assert.Equal(t, err, testErr)
		})

		testObejct1 := interfaces.OrgPerm{}
		testObejct1.SubjectID = strID
		testObejct1.Object = interfaces.Group
		testObejct1.Value = 2
		Convey("AddOrgPerm error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{userInfo1}, nil)
			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(objects, nil)

			poolMock.ExpectBegin()
			o.EXPECT().AddOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			o.EXPECT().UpdateOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			poolMock.ExpectRollback()
			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{testObejct1})
			assert.Equal(t, err, testErr)
			assert.Equal(t, true, true)
		})

		Convey("success", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			u.EXPECT().GetUserDBInfo2(gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{userInfo1}, nil)
			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(objects, nil)

			poolMock.ExpectBegin()
			o.EXPECT().AddOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			o.EXPECT().UpdateOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			poolMock.ExpectCommit()
			err := perm.SetOrgPerm(ctx, strID, interfaces.RealName, []interfaces.OrgPerm{testObejct1})
			assert.Equal(t, err, nil)
		})
	})
}

func TestOrgHandlePermInfo(t *testing.T) {
	Convey("handlePermInfo, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPerm(ctrl)

		perm := newOrgPerm(o)
		perm.logger = common.NewLogger()

		id := strID1
		name := strID1
		currentPerms := make(map[interfaces.OrgType]interfaces.OrgPerm)
		insertPerms := make([]interfaces.OrgPerm, 0)

		insertPerm1 := interfaces.OrgPerm{}
		insertPerm1.SubjectID = strID
		insertPerm1.Object = interfaces.User
		insertPerms = append(insertPerms, insertPerm1)
		Convey("subject is not same with id", func() {
			_, _, err := perm.handlePermInfo(id, name, currentPerms, insertPerms)
			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "subject is not same as subject id"))
		})

		insertPerms[0].SubjectID = id
		insertPerm2 := interfaces.OrgPerm{}
		insertPerm2.SubjectID = id
		insertPerm2.Object = interfaces.Department
		insertPerms = append(insertPerms, insertPerm2)

		curPerm1 := interfaces.OrgPerm{}
		curPerm1.Object = interfaces.Department
		curPerm2 := interfaces.OrgPerm{}
		curPerm2.Object = interfaces.Group
		currentPerms[interfaces.Department] = curPerm1
		currentPerms[interfaces.Group] = curPerm2
		Convey("success", func() {
			insertData, updateData, err := perm.handlePermInfo(id, name, currentPerms, insertPerms)
			assert.Equal(t, err, nil)

			assert.Equal(t, len(insertData), 1)
			assert.Equal(t, insertData[0].Object, interfaces.User)

			assert.Equal(t, len(updateData), 1)
			assert.Equal(t, updateData[0].Object, interfaces.Department)
		})
	})
}

func TestCheckPerms(t *testing.T) {
	Convey("CheckPerms, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPerm(ctrl)
		u := mock.NewMockDBUser(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		perm := newOrgPerm(o)
		perm.user = u
		perm.trace = trace
		perm.logger = common.NewLogger()

		ctx := context.Background()
		testErr := errors.New("xxx111")

		Convey("GetPermByID error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			result, err := perm.CheckPerms(ctx, strID, interfaces.User, interfaces.OPRead)
			assert.Equal(t, result, false)
			assert.Equal(t, err, testErr)
		})

		out := make(map[interfaces.OrgType]interfaces.OrgPerm)
		Convey("not org type", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(out, nil)

			result, err := perm.CheckPerms(ctx, strID, interfaces.User, interfaces.OPRead)
			assert.Equal(t, result, false)
			assert.Equal(t, err, nil)
		})

		out[interfaces.User] = interfaces.OrgPerm{
			Value: 2,
		}
		Convey("org type has no auth", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(out, nil)

			result, err := perm.CheckPerms(ctx, strID, interfaces.User, interfaces.OPRead)
			assert.Equal(t, result, false)
			assert.Equal(t, err, nil)
		})

		out[interfaces.User] = interfaces.OrgPerm{
			Value: 3,
		}
		Convey("success", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			o.EXPECT().GetPermByID(gomock.Any(), gomock.Any()).AnyTimes().Return(out, nil)

			result, err := perm.CheckPerms(ctx, strID, interfaces.User, interfaces.OPRead)
			assert.Equal(t, result, true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestOnUserDeleted(t *testing.T) {
	Convey("onUserDeleted, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mock.NewMockTraceClient(ctrl)

		o := mock.NewMockDBOrgPerm(ctrl)
		perm := newOrgPerm(o)
		perm.trace = trace

		Convey("DeleteOrgPermByID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			o.EXPECT().DeleteOrgPermByID(gomock.Any()).AnyTimes().Return(testErr)
			err := perm.onUserDeleted(strID)

			assert.Equal(t, err, testErr)
		})

		Convey("Success err", func() {
			o.EXPECT().DeleteOrgPermByID(gomock.Any()).AnyTimes().Return(nil)
			err := perm.onUserDeleted(strID)

			assert.Equal(t, err, nil)
		})
	})
}
