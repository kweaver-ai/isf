package logics

import (
	"testing"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/magiconair/properties/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"AuditLog/common"
	"AuditLog/interfaces"
	"AuditLog/interfaces/mock"
	"AuditLog/models"
	"AuditLog/test/mock_log"
	mock_msqclient "AuditLog/test/mock_mqclient"
)

func newdepend(t *testing.T) (*mock_log.MockLogger, *mock.MockLogRepo, *mock.MockLogRepo, *mock.MockLogRepo, *mock.MockUserMgntRepo,
	*mock_msqclient.MockMQClient, *mock.MockOutbox,
) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return mock_log.NewMockLogger(ctrl), mock.NewMockLogRepo(ctrl), mock.NewMockLogRepo(ctrl), mock.NewMockLogRepo(ctrl), mock.NewMockUserMgntRepo(ctrl),
		mock_msqclient.NewMockMQClient(ctrl), mock.NewMockOutbox(ctrl)
}

func newLogMgnt(logger *mock_log.MockLogger, loginLog *mock.MockLogRepo, operLog *mock.MockLogRepo, mgntLog *mock.MockLogRepo, userMgnt *mock.MockUserMgntRepo,
	dbPool *sqlx.DB, cache *redis.Client, mqClient *mock_msqclient.MockMQClient, outbox *mock.MockOutbox,
) interfaces.LogMgnt {
	logmgnt := &logMgnt{
		logger:        logger,
		loginLogRepo:  loginLog,
		operLogRepo:   operLog,
		mgntLogRepo:   mgntLog,
		userMgntRepo:  userMgnt,
		cache:         cache,
		uniqueCacheID: "as:audit_log:unique_id:",
		mqClient:      mqClient,
		outbox:        outbox,
		dbPool:        dbPool,
		cacheTimeout:  300 * time.Second,
	}

	return logmgnt
}

func TestReceiveLog(t *testing.T) {
	Convey("ReceiveLog", t, func() {
		// ctx := context.Background()
		dbPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		redisClient, redisMock := redismock.NewClientMock()
		logger, loginLog, operLog, mgntLog, userMgnt, mqClient, outbox := newdepend(t)
		logmgnt := newLogMgnt(logger, loginLog, operLog, mgntLog, userMgnt, dbPool, redisClient, mqClient, outbox)
		Convey("用户 记录登录日志成功, 只有user_id 200", func() {
			info := &models.ReceiveLogVo{
				Language: "zh-cn",
				LogType:  "login",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "",
					UserType:       "authenticated_user",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			uniqueCacheID := "as:audit_log:unique_id:" + info.LogContent.OutBizID
			redisMock.ExpectGet(uniqueCacheID).SetErr(redis.Nil)
			redisMock.Regexp().ExpectSet("as:audit_log:unique_id:"+info.LogContent.OutBizID, true, 300*time.Second).SetVal("OK")
			userInfo := []models.User{
				{
					ID:      "111",
					Roles:   []string{""},
					Name:    "user1",
					Account: "user1",
					ParentDeps: []interface{}{
						[]interface{}{
							map[string]interface{}{
								"id":   "f114a570-66a9-11eb-ad9d-0050568274c4",
								"name": "xx公司",
								"type": "department",
							},
							map[string]interface{}{
								"id":   "f214a570-66a9-11eb-ad9d-0050568274c4",
								"name": "研发部",
								"type": "department",
							},
						},
						[]interface{}{
							map[string]interface{}{
								"id":   "f114a570-66a9-11eb-ad9d-0050568274c4",
								"name": "bb公司",
								"type": "department",
							},
							map[string]interface{}{
								"id":   "f214a570-66a9-11eb-ad9d-0050568274c4",
								"name": "设计部",
								"type": "department",
							},
						},
					},
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{"111"}).Return(userInfo, 200, nil)
			logger.EXPECT().Infof(gomock.Any(), gomock.Any())
			logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			loginLog.EXPECT().New(gomock.Any()).Return("", nil)
			err = logmgnt.ReceiveLog(info)
			assert.Equal(t, err, nil)
		})
		Convey("用户 记录操作日志成功, 只有user_id 200", func() {
			info := &models.ReceiveLogVo{
				Language: "zh-cn",
				LogType:  "operation",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "",
					UserType:       "authenticated_user",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			redisMock.ExpectGet("as:audit_log:unique_id:" + info.LogContent.OutBizID).RedisNil()
			redisMock.Regexp().ExpectSet("as:audit_log:unique_id:", true, 300*time.Second).SetVal("OK")
			userInfo := []models.User{
				{
					ID:      "111",
					Roles:   []string{""},
					Name:    "user1",
					Account: "user1",
					ParentDeps: []interface{}{
						[]interface{}{
							map[string]interface{}{
								"id":   "f114a570-66a9-11eb-ad9d-0050568274c4",
								"name": "xx公司",
								"type": "department",
							},
							map[string]interface{}{
								"id":   "f214a570-66a9-11eb-ad9d-0050568274c4",
								"name": "研发部",
								"type": "department",
							},
						},
						[]interface{}{
							map[string]interface{}{
								"id":   "f114a570-66a9-11eb-ad9d-0050568274c4",
								"name": "bb公司",
								"type": "department",
							},
							map[string]interface{}{
								"id":   "f214a570-66a9-11eb-ad9d-0050568274c4",
								"name": "设计部",
								"type": "department",
							},
						},
					},
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{"111"}).Return(userInfo, 200, nil)
			logger.EXPECT().Infof(gomock.Any(), gomock.Any())
			logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			operLog.EXPECT().New(gomock.Any()).Return("", nil)
			err = logmgnt.ReceiveLog(info)
			assert.Equal(t, err, nil)
		})
		Convey("用户 记录管理日志成功, 只有user_id 200", func() {
			info := &models.ReceiveLogVo{
				Language: "zh-cn",
				LogType:  "management",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "",
					UserType:       "authenticated_user",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			redisMock.ExpectGet("as:audit_log:unique_id:" + info.LogContent.OutBizID).RedisNil()
			redisMock.Regexp().ExpectSet("as:audit_log:unique_id:", true, 300*time.Second).SetVal("OK")
			userInfo := []models.User{
				{
					ID:      "111",
					Roles:   []string{""},
					Name:    "user1",
					Account: "user1",
					ParentDeps: []interface{}{
						[]interface{}{
							map[string]interface{}{
								"id":   "f114a570-66a9-11eb-ad9d-0050568274c4",
								"name": "xx公司",
								"type": "department",
							},
							map[string]interface{}{
								"id":   "f214a570-66a9-11eb-ad9d-0050568274c4",
								"name": "研发部",
								"type": "department",
							},
						},
						[]interface{}{
							map[string]interface{}{
								"id":   "f114a570-66a9-11eb-ad9d-0050568274c4",
								"name": "bb公司",
								"type": "department",
							},
							map[string]interface{}{
								"id":   "f214a570-66a9-11eb-ad9d-0050568274c4",
								"name": "设计部",
								"type": "department",
							},
						},
					},
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{"111"}).Return(userInfo, 200, nil)
			logger.EXPECT().Infof(gomock.Any(), gomock.Any())
			logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			mgntLog.EXPECT().New(gomock.Any()).Return("", nil)
			err = logmgnt.ReceiveLog(info)
			assert.Equal(t, err, nil)
		})
		Convey("用户 无需记录日志, 有缓存", func() {
			info := &models.ReceiveLogVo{
				Language: "zh-cn",
				LogType:  "login",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "user1",
					UserType:       "authenticated_user",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			redisMock.ExpectGet("as:audit_log:unique_id:" + info.LogContent.OutBizID).SetVal("OK")
			logger.EXPECT().Infof(gomock.Any(), gomock.Any())
			logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			err = logmgnt.ReceiveLog(info)
			assert.Equal(t, err, nil)
		})
		// Convey("用户 记录日志失败 user_id不存在", func() {
		// 	dbPool, _, err := sqlx.New()
		// 	assert.Equal(t, err, nil)
		// 	defer func() {
		// 		if closeErr := dbPool.Close(); closeErr != nil {
		// 			assert.Equal(t, 1, 1)
		// 		}
		// 	}()

		// 	redisClient, redisMock := redismock.NewClientMock()
		// 	logger, loginLog, operLog, mgntLog, userMgnt, mqClient, outbox := newdepend(t)
		// 	logmgnt := newLogMgnt(logger, loginLog, operLog, mgntLog, userMgnt, dbPool, redisClient, mqClient, outbox)
		// 	info := &models.ReceiveLogVo{
		// 		Language: "zh-cn",
		// 		LogType:  "operation",
		// 		LogContent: &models.AuditLog{
		// 			UserID:         "111",
		// 			UserName:       "",
		// 			UserType:       "authenticated_user",
		// 			Level:          1,
		// 			OpType:         1,
		// 			Date:           1717670099640000,
		// 			IP:             "127.0.0.1",
		// 			Mac:            "",
		// 			Msg:            "我是一个消息",
		// 			Exmsg:          "",
		// 			UserAgent:      "Mozilla/5.0",
		// 			ObjID:          "",
		// 			AdditionalInfo: "",
		// 			OutBizID:       "123",
		// 			DeptPaths:      "",
		// 		},
		// 	}
		// 	redisMock.ExpectGet("as:audit_log:unique_id:" + info.LogContent.OutBizID).RedisNil()
		// 	redisMock.Regexp().ExpectSet("as:audit_log:unique_id:", true, 300*time.Second).SetVal("OK")
		// 	userInfo := []models.User{}
		// 	userMgnt.EXPECT().GetUserInfoByID([]string{"111"}).Return(userInfo, 404, nil)
		// 	logger.EXPECT().Infof(gomock.Any(), gomock.Any())
		// 	logger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		// 	operLog.EXPECT().New(gomock.Any()).Return(nil)
		// 	err = logmgnt.ReceiveLog(info)
		// 	assert.Equal(t, err, nil)
		// })
		Convey("应用账户 记录日志成功 只有user_id", func() {
			info := &models.ReceiveLogVo{
				Language: "zh-cn",
				LogType:  "management",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "",
					UserType:       "app",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			redisMock.ExpectGet("as:audit_log:unique_id:" + info.LogContent.OutBizID).RedisNil()
			redisMock.Regexp().ExpectSet("as:audit_log:unique_id:", true, 300*time.Second).SetVal("OK")
			userInfo := models.App{
				ID:   "111",
				Name: "user1",
			}
			userMgnt.EXPECT().GetAppInfoByID("111").Return(userInfo, 200, nil)
			logger.EXPECT().Infof(gomock.Any(), gomock.Any())
			logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			mgntLog.EXPECT().New(gomock.Any()).Return("", nil)
			err = logmgnt.ReceiveLog(info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestSendLog(t *testing.T) {
	Convey("SendLog", t, func() {
		dbPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		redisClient, _ := redismock.NewClientMock()
		logger, loginLog, operLog, mgntLog, userMgnt, mqClient, outbox := newdepend(t)
		logmgnt := newLogMgnt(logger, loginLog, operLog, mgntLog, userMgnt, dbPool, redisClient, mqClient, outbox)
		Convey("发送登录日志成功", func() {
			info := &models.SendLogVo{
				Language: "zh-cn",
				LogType:  "login",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "",
					UserType:       "app",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			txMock.ExpectBegin()
			txMock.ExpectCommit()
			outbox.EXPECT().AddOutboxInfo(common.AuditLoginLogTopic, gomock.Any(), gomock.Any()).Return(nil)
			outbox.EXPECT().NotifyPushOutboxThread()
			err = logmgnt.SendLog(info)
			assert.Equal(t, err, nil)
		})
		Convey("发送操作日志成功", func() {
			info := &models.SendLogVo{
				Language: "zh-cn",
				LogType:  "operation",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "",
					UserType:       "app",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			txMock.ExpectBegin()
			txMock.ExpectCommit()
			outbox.EXPECT().AddOutboxInfo(common.AuditOperationTopic, gomock.Any(), gomock.Any()).Return(nil)
			outbox.EXPECT().NotifyPushOutboxThread()
			err = logmgnt.SendLog(info)
			assert.Equal(t, err, nil)
		})
		Convey("发送管理日志成功", func() {
			info := &models.SendLogVo{
				Language: "zh-cn",
				LogType:  "management",
				LogContent: &models.AuditLog{
					UserID:         "111",
					UserName:       "",
					UserType:       "app",
					Level:          1,
					OpType:         1,
					Date:           1717670099640000,
					IP:             "127.0.0.1",
					Mac:            "",
					Msg:            "我是一个消息",
					Exmsg:          "",
					UserAgent:      "Mozilla/5.0",
					ObjID:          "",
					AdditionalInfo: "",
					OutBizID:       "123",
					DeptPaths:      "",
				},
			}
			txMock.ExpectBegin()
			txMock.ExpectCommit()
			outbox.EXPECT().AddOutboxInfo(common.AuditManagementTopic, gomock.Any(), gomock.Any()).Return(nil)
			outbox.EXPECT().NotifyPushOutboxThread()
			err = logmgnt.SendLog(info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestWriteAuditLoginLog(t *testing.T) {
	Convey("WriteAuditLoginLog", t, func() {
		dbPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		redisClient, _ := redismock.NewClientMock()
		logger, loginLog, operLog, mgntLog, userMgnt, mqClient, outbox := newdepend(t)
		logmgnt := newLogMgnt(logger, loginLog, operLog, mgntLog, userMgnt, dbPool, redisClient, mqClient, outbox)
		Convey("记录登录日志成功", func() {
			entity := &models.AuditLog{
				UserID:         "111",
				UserName:       "",
				UserType:       "app",
				Level:          1,
				OpType:         1,
				Date:           1717670099640000,
				IP:             "127.0.0.1",
				Mac:            "",
				Msg:            "我是一个消息",
				Exmsg:          "",
				UserAgent:      "Mozilla/5.0",
				ObjID:          "",
				AdditionalInfo: "",
				OutBizID:       "123",
				DeptPaths:      "",
			}
			mqClient.EXPECT().Publish(gomock.Any(), gomock.Any())
			err = logmgnt.WriteAuditLoginLog(entity)
			assert.Equal(t, err, nil)
		})
	})
}

func TestWriteAuditOperationLog(t *testing.T) {
	Convey("WriteAuditOperationLog", t, func() {
		dbPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		redisClient, _ := redismock.NewClientMock()
		logger, loginLog, operLog, mgntLog, userMgnt, mqClient, outbox := newdepend(t)
		logmgnt := newLogMgnt(logger, loginLog, operLog, mgntLog, userMgnt, dbPool, redisClient, mqClient, outbox)
		Convey("记录操作日志成功", func() {
			entity := &models.AuditLog{
				UserID:         "111",
				UserName:       "",
				UserType:       "app",
				Level:          1,
				OpType:         1,
				Date:           1717670099640000,
				IP:             "127.0.0.1",
				Mac:            "",
				Msg:            "我是一个消息",
				Exmsg:          "",
				UserAgent:      "Mozilla/5.0",
				ObjID:          "",
				AdditionalInfo: "",
				OutBizID:       "123",
				DeptPaths:      "",
			}
			mqClient.EXPECT().Publish(gomock.Any(), gomock.Any())
			err = logmgnt.WriteAuditOperationLog(entity)
			assert.Equal(t, err, nil)
		})
	})
}

func TestWriteAuditManagementLog(t *testing.T) {
	Convey("WriteAuditManagementLog", t, func() {
		dbPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dbPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		redisClient, _ := redismock.NewClientMock()
		logger, loginLog, operLog, mgntLog, userMgnt, mqClient, outbox := newdepend(t)
		logmgnt := newLogMgnt(logger, loginLog, operLog, mgntLog, userMgnt, dbPool, redisClient, mqClient, outbox)
		Convey("记录管理日志成功", func() {
			entity := &models.AuditLog{
				UserID:         "111",
				UserName:       "",
				UserType:       "app",
				Level:          1,
				OpType:         1,
				Date:           1717670099640000,
				IP:             "127.0.0.1",
				Mac:            "",
				Msg:            "我是一个消息",
				Exmsg:          "",
				UserAgent:      "Mozilla/5.0",
				ObjID:          "",
				AdditionalInfo: "",
				OutBizID:       "123",
				DeptPaths:      "",
			}
			mqClient.EXPECT().Publish(gomock.Any(), gomock.Any())
			err = logmgnt.WriteAuditManagementLog(entity)
			assert.Equal(t, err, nil)
		})
	})
}
