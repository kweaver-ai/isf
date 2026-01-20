// Package logics avatar AnyShare 用户头像业务逻辑层
package logics

import (
	"context"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/pborman/uuid"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type avatar struct {
	avatarDB          interfaces.DBAvatar
	ossGateway        interfaces.DnOSSGateWay
	pool              *sqlx.DB
	logger            common.Logger
	cleanOffSet       int
	cleanChan         chan struct{}
	maxAvatarSize     int64
	maxContentTypeLen int
	trace             observable.Tracer
}

var (
	avOnce   sync.Once
	avLogics *avatar

	// serviceID 服务ID，OSS存储第一级
	serviceID = "D0AC51CBCECA8357D8805B08DEB7B5D5"
)

// NewAvatar 创建新的avatar对象
func NewAvatar() *avatar {
	avOnce.Do(func() {
		config := common.SvcConfig
		avLogics = &avatar{
			avatarDB:          dbAvatar,
			ossGateway:        dnOSSGateWay,
			pool:              dbPool,
			logger:            common.NewLogger(),
			cleanOffSet:       config.CleanAvatarOffsetTime,
			cleanChan:         make(chan struct{}, 1),
			maxAvatarSize:     48 * 1024,
			maxContentTypeLen: 50,
			trace:             common.SvcARTrace,
		}

		// 启动删除线程
		go avLogics.cleanAvatarThread()

		// 启动时先删除一次无效头像
		avLogics.notifyCleanAvatarThread()
	})

	return avLogics
}

// Get 根据用户ID 获取用户URL
func (a *avatar) Get(ctx context.Context, visitor *interfaces.Visitor, userID string) (url string, err error) {
	a.trace.SetInternalSpanName("业务逻辑-获取用户头像")
	newCtx, span := a.trace.AddInternalTrace(ctx)
	defer func() { a.trace.TelemetrySpanEnd(span, err) }()

	// 获取头像存储信息
	info, err := a.avatarDB.Get(userID)
	if err != nil {
		return
	}

	// 如果没有设置头像，url返回空
	if info.Key == "" {
		return
	}

	// 获取 下载头像URL
	return a.ossGateway.GetDownloadURL(newCtx, visitor, info.OSSID, info.Key)
}

// Update 更新用户头像
func (a *avatar) Update(ctx context.Context, visitor *interfaces.Visitor, typ string, data []byte) (err error) {
	a.trace.SetInternalSpanName("业务逻辑-更新用户头像")
	newCtx, span := a.trace.AddInternalTrace(ctx)
	defer func() { a.trace.TelemetrySpanEnd(span, err) }()

	// 判断用户是否可用更新头像，只支持普通用户
	if visitor.Type != interfaces.RealName || AdminIDMap[visitor.ID] {
		err = rest.NewHTTPError("only support normal user", rest.BadRequest, nil)
		return
	}

	// 检查图片类型
	if len(typ) > a.maxContentTypeLen {
		err = rest.NewHTTPError("invalid params, file type error", rest.BadRequest, nil)
		return
	}

	// 检查图片大小
	if int64(len(data)) > a.maxAvatarSize {
		err = rest.NewHTTPError("invalid params, file is too big", rest.BadRequest, nil)
		return
	}

	// 获取OSSID
	ossID, err := a.getAvaliableOSS(newCtx, visitor)
	if err != nil {
		return
	}

	// 设置KEY
	key := strings.Replace(uuid.New(), "-", "", -1)
	key = strings.ToUpper(key)
	key = path.Join(serviceID, key)

	// 新增用户头像信息
	avatarInfo := interfaces.AvatarOSSInfo{
		Type:    typ,
		OSSID:   ossID,
		Key:     key,
		UserID:  visitor.ID,
		BUseful: false,
		Time:    common.Now().UnixNano() / 1000,
	}
	err = a.avatarDB.Add(&avatarInfo)
	if err != nil {
		return
	}

	// OSS上传图片
	err = a.ossGateway.UploadFile(newCtx, visitor, ossID, key, data)
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
				a.logger.Errorf("avatar Update Transaction Commit Error:%v", err)
				return
			}

			// 触发删除线程
			a.notifyCleanAvatarThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				a.logger.Errorf("avatar Update Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 更新原有头像状态
	err = a.avatarDB.SetAvatarUnableByID(visitor.ID, tx)
	if err != nil {
		return
	}

	// 更新现有头像状态
	err = a.avatarDB.UpdateStatusByKey(avatarInfo.Key, true, tx)
	return
}

// notify 触发无效头像删除线程
func (a *avatar) notifyCleanAvatarThread() {
	select {
	case a.cleanChan <- struct{}{}:
	default:
	}
}

// getAvaliableOSS 获取可用存储
func (a *avatar) getAvaliableOSS(ctx context.Context, visitor *interfaces.Visitor) (id string, err error) {
	// 获取本地可用存储信息
	infos, err := a.ossGateway.GetLocalEnabledOSSInfo(ctx, visitor)
	if err != nil {
		return
	}

	// 获取默认存储
	for _, v := range infos {
		if v.BDefault {
			id = v.ID
			return
		}
	}

	// 如果没有默认存储，且无有效存储
	if len(infos) == 0 {
		err = rest.NewHTTPError("no available oss", rest.InternalServerError, nil)
	} else {
		id = infos[0].ID
	}
	return
}

// cleanAvatarThread 删除无用头像
func (a *avatar) cleanAvatarThread() {
	const t1 time.Duration = 24

	for {
		now := common.Now()
		// 计算下一个1点
		next := now.Add(time.Hour * t1)
		next = time.Date(next.Year(), next.Month(), next.Day(), 1, 0, 0, 0, next.Location())

		t := time.NewTimer(next.Sub(now))

		// 到达清理时间或者被触发时，清理无效头像
		select {
		case <-t.C:
		case <-a.cleanChan:
		}

		a.logger.Debugf("clean avatar start")
		// 删除头像
		a.deleteUselessAvatar()

		a.logger.Debugf("clean avatar end")
	}
}

// deleteUselessAvatar 删除无效连接
func (a *avatar) deleteUselessAvatar() {
	// 获取时间戳
	endTime := common.Now().Add(time.Second * time.Duration(-1*a.cleanOffSet))
	timeStamp := endTime.UnixNano() / 1000 // 获取时间戳

	// 获取需要删除的头像信息
	out, err := a.avatarDB.GetUselessAvatar(timeStamp)
	if err != nil {
		a.logger.Errorf("deleteUselessAvatar GetUselessAvatar err: %v", err)
		return
	}

	// 删除头像
	ctx := context.Background()
	visitor := &interfaces.Visitor{
		ErrorCodeType: interfaces.Number,
	}
	for _, v := range out {
		// 删除OSS文件
		err = a.ossGateway.DeleteFile(ctx, visitor, v.OSSID, v.Key)
		if err != nil {
			a.logger.Errorf("deleteUselessAvatar ossGateway DeleteFile err: %v", err)
			continue
		}

		// 删除数据库信息
		err = a.avatarDB.Delete(v.Key)
		if err != nil {
			a.logger.Errorf("deleteUselessAvatar avatarDB Delete err: %v", err)
		}
	}
}
