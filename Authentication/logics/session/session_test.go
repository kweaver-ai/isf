package session

import (
	"database/sql"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func newSession(db interfaces.DBSession) *session {
	return &session{
		db: db,
	}
}

func TestGet(t *testing.T) {
	Convey("get session", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBSession(ctrl)
		s := newSession(db)

		Convey("session_id not exist", func() {
			sessionID := "b1c7384f-78cc-4565-8186-4e4b76d5fe32"
			testErr := sql.ErrNoRows
			db.EXPECT().Get(gomock.Any()).AnyTimes().Return(nil, testErr)
			context, err := s.Get(sessionID)
			assert.Equal(t, context, interfaces.Context{})
			assert.NotEqual(t, err, nil)
		})
	})
}

func TestPut(t *testing.T) {
	Convey("put session", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBSession(ctrl)
		s := newSession(db)

		Convey("sessionID empty", func() {
			ctx := interfaces.Context{}
			err := s.Put(ctx)
			assert.Equal(t, err, rest.NewHTTPError("session_id is invalid", rest.URINotExist, nil))
		})

		Convey("success", func() {
			ctx := interfaces.Context{}
			ctx.SessionID = "xxxx"
			db.EXPECT().Put(gomock.Any()).AnyTimes().Return(nil)
			err := s.Put(ctx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDelete(t *testing.T) {
	Convey("delete session", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBSession(ctrl)
		s := newSession(db)

		Convey("success", func() {
			sessionID := "b1c7384f-78cc-4565-8186-4e4b76d5fe32"
			db.EXPECT().Delete(gomock.Any()).AnyTimes().Return(nil)
			err := s.Delete(sessionID)
			assert.Equal(t, err, nil)
		})
	})
}

func TestEcronDelete(t *testing.T) {
	Convey("delete session", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBSession(ctrl)
		s := newSession(db)

		Convey("success", func() {
			var exp int64 = 1597742762762165184
			db.EXPECT().EcronDelete(gomock.Any()).AnyTimes().Return(nil)
			err := s.EcronDelete(exp)
			assert.Equal(t, err, nil)
		})
	})
}
