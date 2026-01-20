// Package logics combine Anyshare 应用账户业务逻辑层
package logics

import (
	"context"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"

	"UserManagement/common"
	uerrors "UserManagement/errors"
	"UserManagement/interfaces"
)

type app struct {
	db            interfaces.DBApp
	userDB        interfaces.DBUser
	ob            interfaces.LogicsOutbox
	role          interfaces.LogicsRole
	hydra         interfaces.DrivenHydra
	eacpLog       interfaces.DrivenEacpLog
	messageBroker interfaces.DrivenMessageBroker
	logger        common.Logger
	pool          *sqlx.DB
	trace         observable.Tracer
	i18n          *common.I18n
}

var (
	aOnce   sync.Once
	aLogics *app

	maxLen  = 8
	maxYear = 100
)

// NewApp 新建注册接口操作对象
func NewApp() *app {
	aOnce.Do(func() {
		aLogics = &app{
			db:            dbApp,
			userDB:        dbUser,
			ob:            NewOutbox(OutboxBusinessApp),
			role:          NewRole(),
			hydra:         dnHydra,
			eacpLog:       dnEacpLog,
			messageBroker: dnMessageBroker,
			logger:        common.NewLogger(),
			pool:          dbPool,
			trace:         common.SvcARTrace,
			i18n: common.NewI18n(common.I18nMap{
				i18nIDObjectsInAppNotFound: {
					interfaces.SimplifiedChinese:  "此应用账户已不存在。",
					interfaces.TraditionalChinese: "此應用賬號已不存在。",
					interfaces.AmericanEnglish:    "This application account no longer exists.",
				},
			}),
		}

		aLogics.ob.RegisterHandlers(outboxDeleteApp, func(content interface{}) error {
			info := content.(map[string]interface{})
			if err := aLogics.hydra.Delete(info["id"].(string)); err != nil {
				return err
			}
			return nil
		})

		aLogics.ob.RegisterHandlers(oubtoxAppDeleted, func(content interface{}) error {
			info := content.(map[string]interface{})
			if err := aLogics.messageBroker.Publish(interfaces.AppDeleted, info["id"].(string)); err != nil {
				return err
			}
			return nil
		})

		aLogics.ob.RegisterHandlers(outboxUpdateApp, func(content interface{}) error {
			info := content.(map[string]interface{})
			if err := aLogics.hydra.Update(info["id"].(string), info["name"].(string), info["password"].(string)); err != nil {
				return err
			}
			return nil
		})

		aLogics.ob.RegisterHandlers(outboxAppNameChanged, func(content interface{}) error {
			contentJSON := content.(map[string]interface{})
			tmpAppInfo := interfaces.AppInfo{
				ID:   contentJSON["id"].(string),
				Name: contentJSON["name"].(string),
			}

			if err := aLogics.messageBroker.Publish(interfaces.AppNameChanged, tmpAppInfo); err != nil {
				return err
			}
			return nil
		})

		aLogics.ob.RegisterHandlers(outboxAppRegisteredLog, aLogics.sendAppRegisterAuditLog)
		aLogics.ob.RegisterHandlers(outboxAppDeletedLog, aLogics.sendAppDeletedAuditLog)
		aLogics.ob.RegisterHandlers(outboxAppModifiedLog, aLogics.sendAppModifiedAuditLog)
		aLogics.ob.RegisterHandlers(outboxAppTokenGeneratedLog, aLogics.sendAppTokenGeneratedAuditLog)
	})

	return aLogics
}

// RegisterApp 应用账户注册
func (a *app) RegisterApp(visitor *interfaces.Visitor, name, password string, appType interfaces.AppType) (id string, err error) {
	credentialType := interfaces.CredentialTypePassword
	if password == "" {
		credentialType = interfaces.CredentialTypeToken
	}

	// 名称检查
	err = a.checkName(name)
	if err != nil {
		return
	}

	// 密码检查，如果是外部接口，则需解密
	strPwd := password
	if visitor != nil && credentialType == interfaces.CredentialTypePassword && visitor.ID != "" {
		strPwd, err = decodeRSA(password, RSA2048)
		if err != nil {
			err = rest.NewHTTPError(err.Error(), rest.BadRequest, nil)
			return
		}
	}

	if credentialType == interfaces.CredentialTypePassword {
		err = a.checkPassword(strPwd)
		if err != nil {
			return
		}
	}

	// 如果是外部接口，则判断权限，否则不判断
	if visitor != nil && visitor.ID != "" {
		err = checkManageAuthority(a.role, visitor.ID)
		if err != nil {
			return
		}
	}

	// 重名检查
	err = a.duplicateNameCheck(name)
	if err != nil {
		return
	}

	lifespan := 0
	if credentialType == interfaces.CredentialTypeToken {
		// 100年有效
		now := common.Now()
		future := now.AddDate(maxYear, 0, 0) // 加100年
		diff := future.Sub(now)
		lifespan = int(diff.Hours())
	}

	id, err = a.hydra.Register(name, strPwd, lifespan)
	if err != nil {
		return "", err
	}

	// 获取事务处理器
	tx, err := a.pool.Begin()
	if err != nil {
		return "", err
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				a.logger.Errorf("RegisterApp Transaction Commit Error:%v", err)
				return
			}

			a.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				a.logger.Errorf("RegisterApp Rollback err:%v", rollbackErr)
			}
		}
	}()

	info := &interfaces.AppCompleteInfo{
		AppInfo: interfaces.AppInfo{
			ID:             id,
			Name:           name,
			CredentialType: credentialType,
		},
		Type: appType,
	}

	info.Password, err = a.hash([]byte(strPwd))
	if err != nil {
		return "", err
	}

	err = a.db.RegisterApp(info, tx)
	if err != nil {
		return
	}

	// 如果是内部应用账号或者外部专用账户，则使用eacplog的系统id作为用户ID
	if visitor == nil {
		visitor = &interfaces.Visitor{
			ID: common.EacpLogSystemID,
		}
	}

	// 记录审计日志
	content := make(map[string]interface{})
	content["visitor"] = *visitor
	content["name"] = name
	err = a.ob.AddOutboxInfo(outboxAppRegisteredLog, content, tx)
	return
}

// sendAppRegisterAuditLog 注册应用账户发送审计消息
func (a *app) sendAppRegisterAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	name := info["name"].(string)
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		a.logger.Errorf("sendAppRegisterAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = a.eacpLog.EacpLog(&v, interfaces.OpAppRegister, name)
	if err != nil {
		a.logger.Errorf("sendAppRegisterAuditLog err:%v", err)
	}
	return err
}

// DeleteApp 删除应用账户
func (a *app) DeleteApp(visitor *interfaces.Visitor, id string) (err error) {
	// 如果是外部接口，则判断权限，否则不判断
	if visitor != nil && visitor.ID != "" {
		err = checkManageAuthority(a.role, visitor.ID)
		if err != nil {
			return
		}
	}

	appInfo, err := a.db.GetAppByID(id)
	if err != nil || appInfo == nil {
		return
	}

	// 获取事务处理器
	tx, err := a.pool.Begin()
	if err != nil {
		return
	}
	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				a.logger.Errorf("DeleteApp Transaction Commit Error:%v", err)
				return
			}

			// notify outbox推送线程
			a.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				a.logger.Errorf("DeleteApp Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 插入outbox 信息
	contentJSON := make(map[string]interface{})
	contentJSON["id"] = id

	err = a.ob.AddOutboxInfo(outboxDeleteApp, contentJSON, tx)
	if err != nil {
		a.logger.Errorf("Add Outbox Info err:%v", err)
		return
	}

	err = a.ob.AddOutboxInfo(oubtoxAppDeleted, contentJSON, tx)
	if err != nil {
		a.logger.Errorf("Add Outbox Info err:%v", err)
		return
	}

	err = a.db.DeleteApp(id, tx)
	if err != nil {
		return
	}

	// 记录审计日志
	if visitor != nil && visitor.ID != "" {
		content := make(map[string]interface{})
		content["visitor"] = *visitor
		content["name"] = appInfo.Name
		err = a.ob.AddOutboxInfo(outboxAppDeletedLog, content, tx)
	}
	return
}

// sendAppDeletedAuditLog 删除应用账户发送审计消息
func (a *app) sendAppDeletedAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	name := info["name"].(string)
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		a.logger.Errorf("sendAppDeletedAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = a.eacpLog.EacpLog(&v, interfaces.OpDeleteApp, name)
	if err != nil {
		a.logger.Errorf("sendAppDeletedAuditLog err:%v", err)
	}
	return err
}

// 检查名称和密码，并且对密码解码
func (a *app) checkNameAndPWD(bName bool, name string, bPwd bool, pwd string) (strPwd string, err error) {
	// 检查名称是否合法
	if bName {
		err = a.checkName(name)
		if err != nil {
			return "", err
		}
	}

	// 检查密码是否合法
	if bPwd {
		strPwd, err = decodeRSA(pwd, RSA2048)
		if err != nil {
			err = rest.NewHTTPError(err.Error(), rest.BadRequest, nil)
			return "", err
		}
		err = a.checkPassword(strPwd)
		if err != nil {
			return "", err
		}
	}
	return strPwd, nil
}

// UpdateApp 更新应用账户
func (a *app) UpdateApp(visitor *interfaces.Visitor, id string, bName bool, name string, bPwd bool, pwd string) (err error) {
	// 检查名称,密码是否合法,并且对密码解密
	var strPwd string
	strPwd, err = a.checkNameAndPWD(bName, name, bPwd, pwd)
	if err != nil {
		return
	}

	// 判断是否具有更新权限
	err = checkManageAuthority(a.role, visitor.ID)
	if err != nil {
		return
	}

	// 验证应用账户ID存在
	appInfo, err := a.GetApp(id)
	if err != nil {
		return
	}

	if appInfo.CredentialType == interfaces.CredentialTypeToken && bPwd {
		return rest.NewHTTPErrorV2(rest.BadRequest, "app credential type is token, password cannot be updated")
	}

	appName := name
	if !bName {
		appName = appInfo.Name
	} else if appInfo.Name != name {
		// 重名检查
		err = a.duplicateNameCheck(name)
		if err != nil {
			return
		}
	}

	// 密码加密
	var password string
	if bPwd {
		password, err = a.hash([]byte(strPwd))
		if err != nil {
			return
		}
	}

	// 获取事务处理器
	tx, err := a.pool.Begin()
	if err != nil {
		return
	}
	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				a.logger.Errorf("UpdateApp Transaction Commit Error:%v", err)
				return
			}

			// notify outbox推送线程
			a.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				a.logger.Errorf("UpdateApp Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 插入outbox 信息
	contentJSON := make(map[string]interface{})
	contentJSON["id"] = id
	contentJSON["name"] = appName
	contentJSON["password"] = strPwd

	err = a.ob.AddOutboxInfo(outboxUpdateApp, contentJSON, tx)
	if err != nil {
		a.logger.Errorf("Add Outbox Info err:%v", err)
		return
	}

	// 判断应用账户名修改
	if appInfo.Name != appName {
		err = a.ob.AddOutboxInfo(outboxAppNameChanged, contentJSON, tx)
		if err != nil {
			a.logger.Errorf("Add Outbox Info err:%v", err)
			return
		}
	}

	err = a.db.UpdateApp(id, bName, name, bPwd, password, tx)
	if err != nil {
		a.logger.Errorf("UpdateApp Info err:%v", err)
		return
	}

	// 记录审计日志
	if visitor != nil && visitor.ID != "" {
		content := make(map[string]interface{})
		content["visitor"] = *visitor
		content["name"] = appName
		err = a.ob.AddOutboxInfo(outboxAppModifiedLog, content, tx)
	}
	return
}

// sendAppModifiedAuditLog 更新应用账户发送审计消息
func (a *app) sendAppModifiedAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	name := info["name"].(string)
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		a.logger.Errorf("sendAppModifiedAuditLog mapstructure.Decode err:%v", err)
		return
	}
	err = a.eacpLog.EacpLog(&v, interfaces.OpUpdateApp, name)
	if err != nil {
		a.logger.Errorf("sendAppModifiedAuditLog err:%v", err)
	}
	return err
}

// AppList 应用账户列表
func (a *app) AppList(visitor *interfaces.Visitor, searchInfo *interfaces.SearchInfo) (info *[]interfaces.AppInfo, num int, err error) {
	// 判断权限
	err = checkGetInfoAuthority(a.role, visitor.ID)
	if err != nil {
		return
	}

	num, err = a.db.AppListCount(searchInfo)
	if err != nil {
		return
	}

	info, err = a.db.AppList(searchInfo)
	if err != nil {
		return
	}

	return
}

// GetApp 获取应用账户
func (a *app) GetApp(id string) (appInfo *interfaces.AppInfo, err error) {
	appInfo, err = a.db.GetAppByID(id)
	if err != nil {
		return
	} else if appInfo == nil {
		return nil, rest.NewHTTPErrorV2(uerrors.NotFound, "id does not exist", rest.SetCodeStr(uerrors.StrNotFoundAppNotFound))
	}

	return
}

// GenerateAppToken 生成应用账户令牌
func (a *app) GenerateAppToken(ctx context.Context, visitor *interfaces.Visitor, appID string) (token string, err error) {
	// trace
	a.trace.SetInternalSpanName("业务逻辑-生成应用账户令牌")
	_, span := a.trace.AddInternalTrace(ctx)
	defer func() { a.trace.TelemetrySpanEnd(span, err) }()

	// 权限检查
	if visitor != nil && visitor.ID != "" {
		err = checkManageAuthority(a.role, visitor.ID)
		if err != nil {
			return
		}
	}

	// 判断应用账户是否存在
	appInfo, err := a.db.GetAppByID(appID)
	if err != nil {
		return
	} else if appInfo == nil {
		return "", gerrors.NewError(uerrors.StrBadRequestAppNotFound, a.i18n.Load(i18nIDObjectsInAppNotFound, visitor.Language),
			gerrors.SetDetail(map[string]interface{}{"id": appID}))
	}

	// 账户密码类型不允许
	if appInfo.CredentialType == interfaces.CredentialTypePassword {
		return "", gerrors.NewError(gerrors.PublicBadRequest, "app credential type is password, token cannot be generated")
	}

	// 修改应用账户密码
	strPwd := a.generateRandomPassword()
	err = a.hydra.Update(appInfo.ID, appInfo.Name, strPwd)
	if err != nil {
		return
	}

	// 删除应用账户token
	if err = a.hydra.DeleteClientToken(appInfo.ID); err != nil {
		return
	}

	// 生成令牌
	token, err = a.hydra.GenerateToken(appInfo.ID, strPwd)
	if err != nil {
		return
	}

	// 获取事务处理器
	tx, err := a.pool.Begin()
	if err != nil {
		return
	}
	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				a.logger.Errorf("GenerateAppToken Transaction Commit Error:%v", err)
				return
			}

			// notify outbox推送线程
			a.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				a.logger.Errorf("GenerateAppToken Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 记录审计日志
	content := make(map[string]interface{})
	content["visitor"] = *visitor
	content["name"] = appInfo.Name
	err = a.ob.AddOutboxInfo(outboxAppTokenGeneratedLog, content, tx)
	return token, err
}

// sendAppTokenGeneratedAuditLog 生成应用账户令牌发送审计消息
func (a *app) sendAppTokenGeneratedAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	name := info["name"].(string)
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		a.logger.Errorf("sendAppModifiedAuditLog mapstructure.Decode err:%v", err)
		return
	}
	err = a.eacpLog.EacpLog(&v, interfaces.OpAppTokenGenerated, name)
	if err != nil {
		a.logger.Errorf("sendAppModifiedAuditLog err:%v", err)
	}
	return err
}

// ConvertAppName 根据应用账户ID获取应用账户名
func (a *app) ConvertAppName(ids []string, bV2, bStrict bool) (nameInfo []interfaces.NameInfo, err error) {
	nameInfo = make([]interfaces.NameInfo, 0)
	if len(ids) == 0 {
		return nameInfo, nil
	}

	nameInfo, exsitIDs, err := a.db.GetAppName(ids)
	if err != nil {
		return nameInfo, err
	}

	// 如果严格模式， 且有应用账户不存在，则返回错误
	if bStrict && len(exsitIDs) != len(ids) {
		// 获取不存在的应用账户id
		notExistIDs := Difference(ids, exsitIDs)
		if bV2 {
			err = rest.NewHTTPErrorV2(uerrors.AppNotFound, "app does not exist",
				rest.SetDetail(map[string]interface{}{"ids": notExistIDs}),
				rest.SetCodeStr(uerrors.StrBadRequestAppNotFound))
		} else {
			err = rest.NewHTTPErrorV2(uerrors.NotFound, "app does not exist",
				rest.SetDetail(map[string]interface{}{"ids": notExistIDs}))
		}

		return nil, err
	}

	return
}

// generateRandomPassword 生成8位随机密码
func (a *app) generateRandomPassword() (pwd string) {
	// 定义字符集
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	password := make([]byte, maxLen)

	// 生成剩余的6位随机字符
	for i := 0; i < maxLen; i++ {
		password[i] = charset[rand.Intn(len(charset))]
	}

	return string(password)
}

func (a *app) hash(data []byte) (pwd string, err error) {
	hasLen := 12
	tmp, err := bcrypt.GenerateFromPassword(data, hasLen)
	if err != nil {
		return
	}

	return string(tmp), nil
}

func (a *app) duplicateNameCheck(name string) (err error) {
	// 检查与其他应用账户名重复
	appInfo, err := a.db.GetAppByName(name)
	if err != nil {
		return
	} else if appInfo != nil {
		err = rest.NewHTTPErrorV2(uerrors.Conflict, "name already exists",
			rest.SetDetail(map[string]interface{}{"type": "app", "id": appInfo.ID}),
			rest.SetCodeStr(uerrors.StrConflictApp))
		return err
	}

	// 应用账户名不能与用户名（非显示名）重复
	exist, err := a.userDB.CheckNameExist(name)
	if err != nil {
		return
	}
	if exist {
		return rest.NewHTTPErrorV2(uerrors.Conflict, "name already exists",
			rest.SetDetail(map[string]interface{}{"type": "user"}),
			rest.SetCodeStr(uerrors.StrConflictApp))
	}

	return nil
}

// checkName 检查名称是否合法
func (a *app) checkName(name string) (err error) {
	illegalChars := []string{" ", "|", "\\", "/", ":", "*", "?", "\"", ">", "<"}
	for _, char := range illegalChars {
		if strings.ContainsAny(name, char) {
			return rest.NewHTTPError("param name is illegal", rest.BadRequest, nil)
		}
	}

	if utf8.RuneCountInString(name) > 128 || utf8.RuneCountInString(name) < 1 {
		return rest.NewHTTPError("param name is illegal", rest.BadRequest, nil)
	}

	return nil
}

// checkName 检查密码是否合法
func (a *app) checkPassword(pwd string) (err error) {
	var isStandard = regexp.MustCompile(`^[\w~!%#$@\-.]+$`).MatchString
	ret := !isStandard(pwd) ||
		utf8.RuneCountInString(pwd) < 6 ||
		utf8.RuneCountInString(pwd) > 100
	if ret {
		return rest.NewHTTPError("param password is illegal", rest.BadRequest, nil)
	}
	return nil
}
