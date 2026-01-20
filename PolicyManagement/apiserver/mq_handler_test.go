package apiserver

import (
	"policy_mgnt/common"
	"policy_mgnt/interfaces/mock"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ory/gojsonschema"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newMQHandler() *mqHandler {
	userCreatedSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userCreatedSchemaStr))
	if err != nil {
		common.NewLogger().Fatalln(err)
	}
	userStatusChangedSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userStatusChangedSchemaStr))
	if err != nil {
		common.NewLogger().Fatalln(err)
	}
	userDeletedSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userDeletedSchemaStr))
	if err != nil {
		common.NewLogger().Fatalln(err)
	}
	return &mqHandler{
		userCreatedSchema:       userCreatedSchema,
		userStatusChangedSchema: userStatusChangedSchema,
		userDeletedSchema:       userDeletedSchema,
	}
}

func TestOnUserDeleted(t *testing.T) {
	Convey("TestOnUserDeleted", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("policy-mgnt")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eventMock := mock.NewMockLogicsEvent(ctrl)
		log := common.NewLogger()

		mq := newMQHandler()
		mq.log = log
		mq.event = eventMock
		Convey("收到的消息缺少id", func() {
			message := []byte(`{"id1": "123"}`)
			err := mq.onUserDeleted(message)

			assert.Equal(t, err, nil)
		})

		Convey("正常情况", func() {
			eventMock.EXPECT().UserDeleted(gomock.Any()).Return(nil)

			message := []byte(`{"id": "123", "name": "test", "email": "test@test.com"}`)
			err := mq.onUserDeleted(message)
			assert.Equal(t, err, nil)
		})
	})
}

func TestOnUserCreated(t *testing.T) {
	Convey("TestOnUserCreated", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("policy-mgnt")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eventMock := mock.NewMockLogicsEvent(ctrl)
		log := common.NewLogger()

		mq := newMQHandler()
		mq.log = log
		mq.event = eventMock
		Convey("收到的消息缺少id", func() {
			message := []byte(`{"name": "test", "email": "test@test.com"}`)
			err := mq.onUserCreated(message)

			assert.Equal(t, err, nil)
		})

		Convey("正常情况", func() {
			eventMock.EXPECT().UserCreated(gomock.Any()).Return(nil)

			message := []byte(`{"id": "123", "name": "test", "email": "test@test.com"}`)
			err := mq.onUserCreated(message)
			assert.Equal(t, err, nil)
		})
	})
}

func TestOnUserStatusChanged(t *testing.T) {
	Convey("TestOnUserStatusChanged", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("policy-mgnt")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eventMock := mock.NewMockLogicsEvent(ctrl)
		log := common.NewLogger()

		mq := newMQHandler()
		mq.log = log
		mq.event = eventMock

		Convey("收到的消息缺少user_id", func() {
			message := []byte(`{"status": true}`)
			err := mq.onUserStatusChanged(message)
			assert.Equal(t, err, nil)
		})

		Convey("收到的消息缺少status", func() {
			message := []byte(`{"user_id": "123"}`)
			err := mq.onUserStatusChanged(message)
			assert.Equal(t, err, nil)
		})

		Convey("正常情况", func() {
			message := []byte(`{"user_id": "123", "status": true}`)
			eventMock.EXPECT().UserStatusChanged(gomock.Any(), gomock.Any()).Return(nil)
			err := mq.onUserStatusChanged(message)
			assert.Equal(t, err, nil)
		})
	})
}
