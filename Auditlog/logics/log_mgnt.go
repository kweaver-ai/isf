package logics

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	gerror "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-redis/redis/v8"

	"AuditLog/common"
	"AuditLog/errors"
	"AuditLog/gocommon/api"
	gettext "AuditLog/infra"
	"AuditLog/interfaces"
	"AuditLog/locale"
	"AuditLog/models"
)

var (
	l     *logMgnt
	lOnce sync.Once
)

type logMgnt struct {
	logger        api.Logger
	loginLogRepo  interfaces.LogRepo // 数据库对象
	mgntLogRepo   interfaces.LogRepo // 数据库对象
	operLogRepo   interfaces.LogRepo // 数据库对象
	userMgntRepo  interfaces.UserMgntRepo
	cache         redis.Cmdable
	uniqueCacheID string
	mqClient      api.MQClient
	outbox        interfaces.Outbox
	dbPool        *sqlx.DB
	cacheTimeout  time.Duration
}

func NewLogMgnt() interfaces.LogMgnt {
	lOnce.Do(func() {
		o := NewOutbox("client_log")
		l = &logMgnt{
			logger:        logger,
			loginLogRepo:  loginLogRepo,
			operLogRepo:   operLogRepo,
			mgntLogRepo:   mgntLogRepo,
			userMgntRepo:  userMgntRepo,
			cache:         redisClient,
			uniqueCacheID: "as:audit_log:unique_id:",
			mqClient:      mqClient,
			outbox:        o,
			dbPool:        dbPool,
			cacheTimeout:  time.Second * 300,
		}
		o.RegisterHandlers(common.AuditLoginLogTopic, l.WriteAuditLoginLog)
		o.RegisterHandlers(common.AuditManagementTopic, l.WriteAuditManagementLog)
		o.RegisterHandlers(common.AuditOperationTopic, l.WriteAuditOperationLog)
	})

	return l
}

func (l *logMgnt) ReceiveLog(info *models.ReceiveLogVo) (err error) {
	// 读缓存
	ctx := context.Background()
	uniqueCacheID := l.uniqueCacheID + info.LogContent.OutBizID
	// 我想获取到redis中指定key的值，命中则返回，不存在和发生未知error则继续执行
	_, err = l.cache.Get(ctx, uniqueCacheID).Result()
	if err == redis.Nil {
		l.logger.Infof("[%v_log] cache not found, write log", info.LogType)
		err = nil
	} else if err != nil {
		// 访问redis发生未知error
		l.logger.Warnf("[%v_log]  get cache error, error is %v", info.LogType, err)
		return
	} else {
		// 已存在缓存，直接返回
		l.logger.Infof("[%v_log] cache found", info.LogType)
		return
	}

	// 实名用户
	if info.LogContent.UserType == common.AuthenticatedUser {
		// user_name dept_paths,additional_info任一参数为空，则从user_management获取用户信息
		if info.LogContent.UserName == "" || info.LogContent.DeptPaths == "" || info.LogContent.AdditionalInfo == "" {
			var statusCode int
			var userInfo []models.User
			// userInfo := make([]models.User, 0)
			userInfo, statusCode, err = l.userMgntRepo.GetUserInfoByID([]string{info.LogContent.UserID})
			if err != nil {
				l.logger.Errorf("[%v_log] get user info error: %v", info.LogType, err)
				return
			}
			// 用户不存在，消费消息
			if statusCode == http.StatusNotFound {
				l.logger.Warnf("[%v_log] not found user, user_id is: %v, err is: %v", info.LogType, info.LogContent.UserID, err)
			}

			// 响应码不对，直接返回
			// if statusCode != http.StatusOK {
			// 	return
			// }

			if len(userInfo) != 0 {
				// user_name 为空
				if info.LogContent.UserName == "" {
					info.LogContent.UserName = userInfo[0].Name
				}

				// dept_paths 为空
				if info.LogContent.DeptPaths == "" {
					// 用户部门信息
					deptPathsSlice := make([]string, 0)

					for _, deptInfo := range userInfo[0].ParentDeps {
						deptAllPaths := make([]string, 0)
						for _, item := range deptInfo.([]interface{}) {
							deptAllPaths = append(deptAllPaths, item.(map[string]interface{})["name"].(string))
						}

						deptPath := strings.Join(deptAllPaths, "/")
						deptPathsSlice = append(deptPathsSlice, deptPath)
					}
					deptPaths := strings.Join(deptPathsSlice, ", ")

					if strings.Contains(strings.Join(userInfo[0].Roles, ","), "normal_user") {
						// 用户没有部门，设置为未分配组
						if deptPaths == "" {
							deptPaths = gettext.TextDomain("IDS_UNDISTRIBUTED_GROUP")
						}
					}

					info.LogContent.DeptPaths = deptPaths
				}
			}
		}
		// 应用账户
	} else if info.LogContent.UserType == common.APP {
		var appInfo models.App
		var statusCode int
		appInfo, statusCode, err = l.userMgntRepo.GetAppInfoByID(info.LogContent.UserID)
		if err != nil {
			l.logger.Errorf("[%v_log] get app info error: %v", info.LogType, err)
			return
		}

		// 应用账户不存在，消费消息
		if statusCode == http.StatusNotFound {
			l.logger.Warnf("[%v_log] not found app, app_id is: %v, err is: %v", info.LogType, info.LogContent.UserID, err)
		}

		// 响应码不对，直接返回
		if statusCode != http.StatusOK {
			return
		}

		if info.LogContent.UserName == "" {
			info.LogContent.UserName = appInfo.Name
		}
	}

	// user_name 为空 统一赋值为user_id
	if info.LogContent.UserName == "" {
		info.LogContent.UserName = info.LogContent.UserID
	}

	// 根据日志类型入库
	logType := info.LogType
	logContent := info.LogContent
	if logType == "login" {
		_, err = l.loginLogRepo.New(logContent)
	} else if logType == "management" {
		_, err = l.mgntLogRepo.New(logContent)
	} else if logType == "operation" {
		_, err = l.operLogRepo.New(logContent)
	} else {
		l.logger.Warnf("category is invalid")
		return
	}

	if err != nil {
		l.logger.Errorf("insert log error: %v", err)
		return err
	}

	// 写缓存
	err = l.cache.Set(ctx, uniqueCacheID, true, l.cacheTimeout).Err()
	if err != nil {
		l.logger.Warnf("set redis cache failed, err is: %v", err)
		return err
	}
	return
}

func (l *logMgnt) ReceiveAuditLog(info *models.ReceiveLogVo) (err error) {
	// 读缓存
	ctx := context.Background()
	uniqueCacheID := l.uniqueCacheID + info.LogContent.OutBizID
	// 我想获取到redis中指定key的值，命中则返回，不存在和发生未知error则继续执行
	_, err = l.cache.Get(ctx, uniqueCacheID).Result()
	if err == redis.Nil {
		l.logger.Infof("[%v_log] cache not found, write log", info.LogType)
		err = nil
	} else if err != nil {
		// 访问redis发生未知error
		l.logger.Warnf("[%v_log]  get cache error, error is %v", info.LogType, err)
		return
	} else {
		// 已存在缓存，直接返回
		l.logger.Infof("[%v_log] cache found", info.LogType)
		return
	}

	// 实名用户
	if info.LogContent.UserType == common.AuthenticatedUser {
		// user_name dept_paths 任一参数为空，则从user_management获取用户信息
		if info.LogContent.UserName == "" || info.LogContent.DeptPaths == "" {
			var statusCode int
			var userInfo []models.User
			// userInfo := make([]models.User, 0)
			userInfo, statusCode, err = l.userMgntRepo.GetUserInfoByID([]string{info.LogContent.UserID})
			if err != nil {
				l.logger.Errorf("[%v_log] get user info error: %v", info.LogType, err)
				return
			}
			// 用户不存在，消费消息
			if statusCode == http.StatusNotFound {
				l.logger.Warnf("[%v_log] not found user, user_id is: %v, err is: %v", info.LogType, info.LogContent.UserID, err)
			}

			if len(userInfo) != 0 {
				// user_name 为空
				if info.LogContent.UserName == "" {
					info.LogContent.UserName = userInfo[0].Name
				}

				// dept_paths 为空
				if info.LogContent.DeptPaths == "" {
					// 用户部门信息
					deptPathsSlice := make([]string, 0)

					for _, deptInfo := range userInfo[0].ParentDeps {
						deptAllPaths := make([]string, 0)
						for _, item := range deptInfo.([]interface{}) {
							deptAllPaths = append(deptAllPaths, item.(map[string]interface{})["name"].(string))
						}

						deptPath := strings.Join(deptAllPaths, "/")
						deptPathsSlice = append(deptPathsSlice, deptPath)
					}
					deptPaths := strings.Join(deptPathsSlice, ", ")

					if strings.Contains(strings.Join(userInfo[0].Roles, ","), "normal_user") {
						// 用户没有部门，设置为未分配组
						if deptPaths == "" {
							deptPaths = gettext.TextDomain("IDS_UNDISTRIBUTED_GROUP")
						}
					}

					info.LogContent.DeptPaths = deptPaths
				}
			}
		}
		// 应用账户
	} else if info.LogContent.UserType == common.APP {
		var appInfo models.App
		var statusCode int
		appInfo, statusCode, err = l.userMgntRepo.GetAppInfoByID(info.LogContent.UserID)
		if err != nil {
			l.logger.Errorf("[%v_log] get app info error: %v", info.LogType, err)
			return
		}

		// 应用账户不存在，消费消息
		if statusCode == http.StatusNotFound {
			l.logger.Warnf("[%v_log] not found app, app_id is: %v, err is: %v", info.LogType, info.LogContent.UserID, err)
		}

		// 响应码不对，直接返回
		if statusCode != http.StatusOK {
			return
		}

		if info.LogContent.UserName == "" {
			info.LogContent.UserName = appInfo.Name
		}
	} else if info.LogContent.UserType == common.AnonymousUser {
		if info.LogContent.UserName == "" {
			info.LogContent.UserName = locale.GetI18nCtx(ctx, common.AnonymousUser)
		}
	}

	// user_name 为空 统一赋值为user_id
	if info.LogContent.UserName == "" {
		info.LogContent.UserName = info.LogContent.UserID
	}

	// 根据日志类型入库
	logType := info.LogType
	logContent := info.LogContent
	if logType == "login" {
		_, err = l.loginLogRepo.New(logContent)
	} else if logType == "management" {
		_, err = l.mgntLogRepo.New(logContent)
	} else if logType == "operation" {
		_, err = l.operLogRepo.New(logContent)
	} else {
		l.logger.Warnf("category is invalid")
		return
	}

	if err != nil {
		l.logger.Errorf("insert log error: %v", err)
		return err
	}

	// 写缓存
	err = l.cache.Set(ctx, uniqueCacheID, true, l.cacheTimeout).Err()
	if err != nil {
		l.logger.Warnf("set redis cache failed, err is: %v", err)
		return err
	}
	return
}

func (l *logMgnt) AddAuditLog(visitor interfaces.Visitor, logType interfaces.LogType, info *models.ReceiveLogVo) (logID string, err error) {
	// 检查 只有实名用户有权限
	if visitor.Type != interfaces.RealName {
		err = gerror.NewError(gerror.PublicForbidden, "user is not authenticated")
		return
	}

	// 获取用户信息
	userInfo, _, err := l.userMgntRepo.GetUserInfoByID([]string{visitor.ID})
	if err != nil {
		l.logger.Errorf("[%v_log] get user info error: %v", info.LogType, err)
		return
	}

	info.LogContent.UserName = userInfo[0].Name
	info.LogContent.UserID = userInfo[0].ID
	info.LogContent.UserType = common.AuthenticatedUser
	info.LogContent.IP = visitor.IP
	info.LogContent.Mac = visitor.Mac
	info.LogContent.UserAgent = visitor.UserAgent

	// 用户部门信息
	deptPathsSlice := make([]string, 0)
	for _, deptInfo := range userInfo[0].ParentDeps {
		deptAllPaths := make([]string, 0)
		for _, item := range deptInfo.([]interface{}) {
			deptAllPaths = append(deptAllPaths, item.(map[string]interface{})["name"].(string))
		}

		deptPath := strings.Join(deptAllPaths, "/")
		deptPathsSlice = append(deptPathsSlice, deptPath)
	}
	deptPaths := strings.Join(deptPathsSlice, ", ")

	if strings.Contains(strings.Join(userInfo[0].Roles, ","), "normal_user") {
		// 用户没有部门，设置为未分配组
		if deptPaths == "" {
			deptPaths = gettext.TextDomain("IDS_UNDISTRIBUTED_GROUP")
		}
	}
	info.LogContent.DeptPaths = deptPaths

	// 根据日志类型入库
	logContent := info.LogContent
	if logType == interfaces.LogType_Login {
		logID, err = l.loginLogRepo.New(logContent)
	} else if logType == interfaces.LogType_Management {
		logID, err = l.mgntLogRepo.New(logContent)
	} else if logType == interfaces.LogType_Operation {
		logID, err = l.operLogRepo.New(logContent)
	} else {
		l.logger.Warnf("category is invalid")
		return "", err
	}

	if err != nil {
		l.logger.Errorf("insert log error: %v", err)
		return "", err
	}
	return logID, nil
}

func (l *logMgnt) SendLog(info *models.SendLogVo) (err error) {
	logType := info.LogType
	logContent := info.LogContent
	tx, err := l.dbPool.Begin()
	if err != nil {
		return err
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			err = tx.Commit()
			if err != nil {
				l.logger.Errorf("recordLog Transaction Commit Error:%v", err)
				return
			}
			// 触发outbox消息推送线程
			l.outbox.NotifyPushOutboxThread()
		default:
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				l.logger.Errorf("recordLog Transaction Rollback Error:%v", rollbackErr)
			}
		}
	}()

	if logType == "login" {
		err = l.outbox.AddOutboxInfo(common.AuditLoginLogTopic, logContent, tx)
	} else if logType == "management" {
		err = l.outbox.AddOutboxInfo(common.AuditManagementTopic, logContent, tx)
	} else if logType == "operation" {
		err = l.outbox.AddOutboxInfo(common.AuditOperationTopic, logContent, tx)
	} else {
		err = errors.New(info.Language, errors.BadRequestErr, "category is invalid", nil)
	}

	return
}

// 记录审计登录日志
func (l *logMgnt) WriteAuditLoginLog(entity interface{}) (err error) {
	msgByte, err := json.Marshal(entity)
	if err != nil {
		l.logger.Warnf("json marshal failed, err is %v", err)
		return
	}
	if err = l.mqClient.Publish(common.AuditLoginLogTopic, msgByte); err != nil {
		return
	}
	return
}

// 记录审计管理日志
func (l *logMgnt) WriteAuditManagementLog(entity interface{}) (err error) {
	msgByte, err := json.Marshal(entity)
	if err != nil {
		l.logger.Warnf("json marshal failed, err is %v", err)
		return
	}
	if err = l.mqClient.Publish(common.AuditManagementTopic, msgByte); err != nil {
		return
	}
	return
}

// 记录审计操作日志
func (l *logMgnt) WriteAuditOperationLog(entity interface{}) (err error) {
	msgByte, err := json.Marshal(entity)
	if err != nil {
		l.logger.Warnf("json marshal failed, err is %v", err)
		return
	}
	if err = l.mqClient.Publish(common.AuditOperationTopic, msgByte); err != nil {
		return
	}
	return
}
