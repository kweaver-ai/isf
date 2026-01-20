package logics

import (
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

var (
	strID    string = "2114a570-66a9-11eb-ad9d-0050568274c4"
	strName  string = "Name1"
	strName1 string = "Name2"
	strID2   string = "ID2"
	strName2 string = "Name3"
)

func newOrgPermApp(o interfaces.DBOrgPermApp) *orgPermApp {
	return &orgPermApp{
		db: o,
	}
}

func TestNewOrgPermApp(t *testing.T) {
	Convey("NewOrgPermApp, db is available", t, func() {
		sqlDB, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		dbPool = sqlDB

		out := NewOrgPermApp()
		assert.NotEqual(t, out, nil)
	})
}

func TestUpdateAppName(t *testing.T) {
	Convey("UpdateAppName, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPermApp(ctrl)
		perm := newOrgPermApp(o)

		var appData interfaces.AppInfo
		Convey("GetRolesByUserID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			o.EXPECT().UpdateAppName(gomock.Any()).AnyTimes().Return(testErr)
			err := perm.UpdateAppName(&appData)

			assert.Equal(t, err, testErr)
		})

		Convey("Success err", func() {
			o.EXPECT().UpdateAppName(gomock.Any()).AnyTimes().Return(nil)
			err := perm.UpdateAppName(&appData)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetAppOrgPerm(t *testing.T) {
	Convey("GetAppOrgPerm, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPermApp(ctrl)
		r := mock.NewMockLogicsRole(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		perm := newOrgPermApp(o)
		perm.role = r
		perm.app = a

		visitor := interfaces.Visitor{}
		objects := []interfaces.OrgType{interfaces.Department, interfaces.User, interfaces.Department}
		userID := strID
		Convey("org type are not uniqued", func() {
			_, err := perm.GetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, rest.NewHTTPError("org type is not uniqued", rest.BadRequest, nil))
		})

		objects = []interfaces.OrgType{interfaces.Department, interfaces.User}
		Convey("vistor type is not realname", func() {
			_, err := perm.GetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", rest.Forbidden, nil))
		})

		visitor.Type = interfaces.RealName
		visitor.ID = strID
		userRole := make(map[string]map[interfaces.Role]bool)
		roles := make(map[interfaces.Role]bool)
		roles[interfaces.SystemRoleSuperAdmin] = true
		userRole[strID] = roles
		testErr := errors.New("xxxx1")
		Convey("GetAppName error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := perm.GetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		appInfo := interfaces.AppInfo{}
		Convey("GetAppPermByID error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := perm.GetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		out := make(map[interfaces.OrgType]interfaces.AppOrgPerm)
		out[interfaces.User] = interfaces.AppOrgPerm{}
		Convey("success", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(out, nil)
			out, err := perm.GetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 1)
		})

		objects = []interfaces.OrgType{interfaces.Department}
		Convey("success1", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(out, nil)
			out, err := perm.GetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 0)
		})
	})
}

func TestDeleteAppOrgPerm(t *testing.T) {
	Convey("DeleteAppOrgPerm, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqlDB, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)

		o := mock.NewMockDBOrgPermApp(ctrl)
		r := mock.NewMockLogicsRole(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		perm := newOrgPermApp(o)
		perm.role = r
		perm.app = a
		perm.logger = common.NewLogger()
		perm.ob = ob
		perm.pool = sqlDB

		visitor := interfaces.Visitor{}
		objects := []interfaces.OrgType{interfaces.Department, interfaces.User, interfaces.Department}
		userID := strID
		Convey("org type are not uniqued", func() {
			err := perm.DeleteAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, rest.NewHTTPError("org type is not uniqued", rest.BadRequest, nil))
		})

		objects = []interfaces.OrgType{interfaces.Department, interfaces.User}
		Convey("vistor type is not realname", func() {
			err := perm.DeleteAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", rest.Forbidden, nil))
		})

		visitor.Type = interfaces.RealName
		visitor.ID = strID
		userRole := make(map[string]map[interfaces.Role]bool)
		roles := make(map[interfaces.Role]bool)
		roles[interfaces.SystemRoleSuperAdmin] = true
		testErr := errors.New("xxx111")
		userRole[strID] = roles
		Convey("app id not exist", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := perm.DeleteAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		appInfo := interfaces.AppInfo{}
		Convey("GetAppPermByID error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := perm.DeleteAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteAppOrgPermInfo error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, nil)
			o.EXPECT().DeleteAppOrgPerm(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := perm.DeleteAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, nil)
			o.EXPECT().DeleteAppOrgPerm(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().AnyTimes()

			err := perm.DeleteAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, nil)
		})
	})
}

func TestSetAppOrgPerm(t *testing.T) {
	Convey("SetAppOrgPerm, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPermApp(ctrl)
		r := mock.NewMockLogicsRole(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		dPool, poolMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		perm := newOrgPermApp(o)
		perm.role = r
		perm.app = a
		perm.pool = dPool
		perm.logger = common.NewLogger()
		perm.ob = ob

		visitor := interfaces.Visitor{}
		Convey("visitor has not authority", func() {
			err := perm.SetAppOrgPerm(&visitor, strID, nil)
			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", rest.Forbidden, nil))
		})

		visitor.Type = interfaces.RealName
		visitor.ID = strID
		userRole := make(map[string]map[interfaces.Role]bool)
		roles := make(map[interfaces.Role]bool)
		roles[interfaces.SystemRoleSuperAdmin] = true
		userRole[strID] = roles

		objects := []interfaces.AppOrgPerm{}
		userID := strID
		testErr := errors.New("xxx111")
		Convey("GetAppName error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		appInfo := interfaces.AppInfo{}
		Convey("GetAppPermByID error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		testObejcts := interfaces.AppOrgPerm{}
		testObejcts.Subject = "xxxxxxx1"
		testObejcts.Object = interfaces.User
		objects = append(objects, testObejcts)
		Convey("handlePermInfo error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, nil)
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, rest.NewHTTPError("subject is not same as app id", rest.BadRequest, nil))
		})

		objects[0].Subject = userID
		Convey("pool begin error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, nil)
			poolMock.ExpectBegin().WillReturnError(testErr)
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		Convey("AddAppOrgPermInfo error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, nil)
			poolMock.ExpectBegin()
			o.EXPECT().AddAppOrgPerm(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			poolMock.ExpectRollback()
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		Convey("AddAppOrgPermInfo outbox error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, nil)
			poolMock.ExpectBegin()
			o.EXPECT().AddAppOrgPerm(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			poolMock.ExpectRollback()
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		outPerms := make(map[interfaces.OrgType]interfaces.AppOrgPerm)
		outPerms[interfaces.User] = interfaces.AppOrgPerm{}
		Convey("UpdateAppOrgPermInfo error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(outPerms, nil)
			poolMock.ExpectBegin()
			o.EXPECT().UpdateAppOrgPerm(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			poolMock.ExpectRollback()
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		Convey("UpdateAppOrgPermInfo  add outbox error", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(outPerms, nil)
			poolMock.ExpectBegin()
			o.EXPECT().UpdateAppOrgPerm(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			poolMock.ExpectRollback()
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, testErr)
		})

		objects = make([]interfaces.AppOrgPerm, 0)
		Convey("success", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(&appInfo, nil)
			o.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, nil)
			ob.EXPECT().NotifyPushOutboxThread().AnyTimes()
			poolMock.ExpectBegin()
			poolMock.ExpectCommit()
			err := perm.SetAppOrgPerm(&visitor, userID, objects)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCheckAppPermManageAuthority(t *testing.T) {
	Convey("checkManageAuthority, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPermApp(ctrl)
		r := mock.NewMockLogicsRole(ctrl)
		a := mock.NewMockLogicsApp(ctrl)

		perm := newOrgPermApp(o)
		perm.role = r
		perm.app = a
		perm.logger = common.NewLogger()

		visitor := interfaces.Visitor{}
		Convey("visitor has not authority", func() {
			err := perm.checkManageAuthority(&visitor)
			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", rest.Forbidden, nil))
		})

		visitor.Type = interfaces.RealName
		visitor.ID = strID
		userRole := make(map[string]map[interfaces.Role]bool)
		roles := make(map[interfaces.Role]bool)
		roles[interfaces.SystemRoleSuperAdmin] = true
		userRole[strID] = roles
		Convey("success", func() {
			r.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRole, nil)
			err := perm.checkManageAuthority(&visitor)
			assert.Equal(t, err, nil)
		})
	})
}

func TestHandlePermInfo(t *testing.T) {
	Convey("handlePermInfo, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockDBOrgPermApp(ctrl)
		r := mock.NewMockLogicsRole(ctrl)
		a := mock.NewMockLogicsApp(ctrl)

		perm := newOrgPermApp(o)
		perm.role = r
		perm.app = a
		perm.logger = common.NewLogger()

		id := "userID1"
		name := "userName"
		currentPerms := make(map[interfaces.OrgType]interfaces.AppOrgPerm)
		insertPerms := make([]interfaces.AppOrgPerm, 0)

		insertPerm1 := interfaces.AppOrgPerm{}
		insertPerm1.Subject = "xxx"
		insertPerm1.Object = interfaces.User
		insertPerms = append(insertPerms, insertPerm1)
		Convey("subject is not same with id", func() {
			_, _, err := perm.handlePermInfo(id, name, currentPerms, insertPerms)
			assert.Equal(t, err, rest.NewHTTPError("subject is not same as app id", rest.BadRequest, nil))
		})

		insertPerms[0].Subject = id
		insertPerm2 := interfaces.AppOrgPerm{}
		insertPerm2.Subject = id
		insertPerm2.Object = interfaces.Department
		insertPerms = append(insertPerms, insertPerm2)

		curPerm1 := interfaces.AppOrgPerm{}
		curPerm1.Object = interfaces.Department
		curPerm2 := interfaces.AppOrgPerm{}
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

func TestSendDeleteOrgPermAppAuditLog(t *testing.T) {
	Convey("sendDeleteOrgPermAppAuditLog, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacpLog := mock.NewMockDrivenEacpLog(ctrl)
		perm := newOrgPermApp(nil)
		perm.eacpLog = eacpLog
		perm.logger = common.NewLogger()

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor": make(map[string]interface{}),
			"perm":    make(map[string]interface{}),
		}
		Convey("error", func() {
			eacpLog.EXPECT().OpDeleteOrgPermAppLog(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := perm.sendDeleteOrgPermAppAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}

func TestSendAddOrgPermAppAuditLog(t *testing.T) {
	Convey("sendAddOrgPermAppAuditLog, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacpLog := mock.NewMockDrivenEacpLog(ctrl)
		perm := newOrgPermApp(nil)
		perm.eacpLog = eacpLog
		perm.logger = common.NewLogger()

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor": make(map[string]interface{}),
			"perm":    make(map[string]interface{}),
		}
		Convey("error", func() {
			eacpLog.EXPECT().OpAddOrgPermAppLog(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := perm.sendAddOrgPermAppAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}

func TestSendUpdateOrgPermAppAuditLog(t *testing.T) {
	Convey("sendUpdateOrgPermAppAuditLog, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacpLog := mock.NewMockDrivenEacpLog(ctrl)
		perm := newOrgPermApp(nil)
		perm.eacpLog = eacpLog
		perm.logger = common.NewLogger()

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor": make(map[string]interface{}),
			"perm":    make(map[string]interface{}),
		}
		Convey("error", func() {
			eacpLog.EXPECT().OpUpdateOrgPermAppLog(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := perm.sendUpdateOrgPermAppAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}
