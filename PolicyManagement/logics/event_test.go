package logics

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserCreated(t *testing.T) {
	Convey("UserCreated, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eve := &event{
			userCreatedHandlers:       make([]func(string) error, 0),
			userStatusChangedHandlers: make([]func(string, bool) error, 0),
		}

		Convey("execute error", func() {
			eve.userCreatedHandlers = []func(string) error{func(string) error {
				return errors.New("test")
			}}
			err := eve.UserCreated("test")
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("execute success", func() {
			eve.userCreatedHandlers = []func(string) error{func(string) error {
				return nil
			}}
			err := eve.UserCreated("test")
			assert.Equal(t, err, nil)
		})
	})
}

func TestUserStatusChanged(t *testing.T) {
	Convey("UserStatusChanged, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eve := &event{
			userCreatedHandlers:       make([]func(string) error, 0),
			userStatusChangedHandlers: make([]func(string, bool) error, 0),
		}

		Convey("execute error", func() {
			eve.userStatusChangedHandlers = []func(string, bool) error{func(string, bool) error {
				return errors.New("test")
			}}
			err := eve.UserStatusChanged("test", true)
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("execute success", func() {
			eve.userStatusChangedHandlers = []func(string, bool) error{func(string, bool) error {
				return nil
			}}
			err := eve.UserStatusChanged("test", true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestRegisterUserCreated(t *testing.T) {
	Convey("RegisterUserCreated, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eve := &event{
			userCreatedHandlers:       make([]func(string) error, 0),
			userStatusChangedHandlers: make([]func(string, bool) error, 0),
		}

		Convey("register user created", func() {
			eve.RegisterUserCreated(func(string) error {
				return nil
			})
			assert.Equal(t, len(eve.userCreatedHandlers), 1)
		})
	})
}

func TestRegisterUserStatusChanged(t *testing.T) {
	Convey("RegisterUserStatusChanged, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eve := &event{
			userCreatedHandlers:       make([]func(string) error, 0),
			userStatusChangedHandlers: make([]func(string, bool) error, 0),
		}

		Convey("register user status changed", func() {
			eve.RegisterUserStatusChanged(func(string, bool) error {
				return nil
			})
			assert.Equal(t, len(eve.userStatusChangedHandlers), 1)
		})
	})
}
