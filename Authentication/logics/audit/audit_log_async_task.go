package audit

import (
	"context"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
)

var (
	auditLogOnce sync.Once
	auditLog     *auditLogAsyncTask
)

type auditLogAsyncTask struct {
	logger   common.Logger
	pushChan chan struct{}
	unOb     interfaces.DBUnorderedOutbox
	eacpLog  interfaces.DnEacpLog
}

// NewAuditLogAsyncTask 创建新的审计日志处理对象
func NewAuditLogAsyncTask() *auditLogAsyncTask {
	auditLogOnce.Do(func() {
		auditLog = &auditLogAsyncTask{
			unOb:     logics.DBUnorderedOutbox,
			pushChan: make(chan struct{}, 1),
			eacpLog:  logics.DnEacpLog,
			logger:   common.NewLogger(),
		}

		// 更新status状态为1的数据，使之重新发送审计日志
		go auditLog.restartUnAsyncTask()

		go auditLog.startAuditLogAsyncTaskThread()
	})
	return auditLog
}

func (auditLog *auditLogAsyncTask) Log(topic string, message interface{}) (err error) {
	contentJSON := make(map[string]interface{})
	contentJSON["topic"] = topic
	contentJSON["message"] = message

	messageStr, err := jsoniter.MarshalToString(contentJSON)
	if err != nil {
		auditLog.logger.Errorln("Log MarshalToString error: %#v", err)
		return
	}
	unorderedOutboxInfo := interfaces.UnorderedOutbox{
		Message:   messageStr,
		Status:    interfaces.OutboxNotStarted,
		CreatedAt: common.Now().UnixNano() / 1e3,
		UpdatedAt: common.Now().UnixNano() / 1e3,
	}

	err = auditLog.unOb.AddUnorderedOutboxInfo(unorderedOutboxInfo)
	if err != nil {
		auditLog.logger.Errorln("Log AddUnorderedOutboxInfo error: %#v", err)
		return
	}

	auditLog.notifyAuditLogAsyncTaskThread()
	return nil
}

// startAuditLogAsyncTaskThread 开启异步审计日志更新线程
func (auditLog *auditLogAsyncTask) startAuditLogAsyncTaskThread() {
	// 每隔30s处理异步任务
	const timeDuration int = 30
	ticker := time.NewTicker(time.Second * time.Duration(timeDuration))
	for {
		select {
		case <-ticker.C:
		case <-auditLog.pushChan:
		}

		// 逐条进行推送
		for {
			isAllFinished := auditLog.push()
			// 出错或消息都推送完毕时退出当前循环
			if isAllFinished {
				break
			}
		}
	}
}

// notifyAuditLogAsyncTaskThread 唤醒异步审计日志更新线程
func (auditLog *auditLogAsyncTask) notifyAuditLogAsyncTaskThread() {
	select {
	case auditLog.pushChan <- struct{}{}:
	default:
	}
}

func (auditLog *auditLogAsyncTask) push() (isAllFinished bool) {
	// 标记是否全部推送完成
	isAllFinished = false

	// 每次处理异步审计日志计算任务前，将 update时间超过60s未更新 且 state为1的任务 更新为0
	updatedTime := common.Now().Add(-1*time.Minute).UnixNano() / (1e3)

	err := auditLog.unOb.RestartUnorderedOutboxInfo(updatedTime)
	if err != nil {
		auditLog.logger.Errorf("restart un async audit log task failed, err: %#v", err)
		return
	}

	// 获取审计日志发送任务
	auditLogAsyncInfo, exist, err := auditLog.unOb.GetUnorderedOutboxInfo()
	if err != nil {
		auditLog.logger.Errorf("GetUnorderedOutboxInfo error: %#v", err)
		return
	}
	if !exist {
		isAllFinished = true
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动自动续期线程
	go auditLog.startAutoRenewalThread(ctx, auditLogAsyncInfo.ID)

	// 出错一直重试，直至处理完成
	tryCount := 0
	for {
		err := func() (errs error) {
			contentJSON := make(map[string]interface{})
			errs = jsoniter.UnmarshalFromString(auditLogAsyncInfo.Message, &contentJSON)
			if errs != nil {
				auditLog.logger.Warnf("push UnmarshalFromString error: %#v", errs)
				return
			}
			message := contentJSON["message"]
			topic, _ := contentJSON["topic"].(string)

			errs = auditLog.eacpLog.Publish(topic, message)
			if errs != nil {
				return
			}

			// 删除任务
			errs = auditLog.unOb.DeleteUnorderedOutboxInfoByID(auditLogAsyncInfo.ID)
			if errs != nil {
				return
			}
			return
		}()
		tryCount++
		if err != nil {
			auditLog.logger.Errorf("process audit log async task retry %d failed, err: %v", tryCount, err)
			// 每次休眠，默认30s
			time.Sleep(time.Second * (30))
		} else {
			break
		}
	}

	// 关闭续期
	cancel()
	return
}

// 自动续期线程
func (auditLog *auditLogAsyncTask) startAutoRenewalThread(ctx context.Context, id string) {
	// 出错时每隔10s重试，直至成功
	const timeDuration int = 10
	ticker := time.NewTicker(time.Second * time.Duration(timeDuration))
	for {
		select {
		case <-ticker.C:
			// 根据ID做自动续期操作
			runAble, err := auditLog.autoRenewal(id)
			if err != nil {
				continue
			}
			if !runAble {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// autoRenewal 自动续期
func (auditLog *auditLogAsyncTask) autoRenewal(id string) (runAble bool, err error) {
	runAble, err = auditLog.unOb.UpdateUnorderedOutboxUpdateTimeByID(id)
	if err != nil {
		return runAble, err
	}
	return runAble, nil
}

func (auditLog *auditLogAsyncTask) restartUnAsyncTask() {
	// 获取当前时间 - 1分钟
	updatedTime := common.Now().Add(-20*time.Second).UnixNano() / 1e3

	err := auditLog.unOb.RestartUnorderedOutboxInfo(updatedTime)
	if err != nil {
		auditLog.logger.Errorf("restart un async audit log task failed, err: %#v", err)
		return
	}

	auditLog.logger.Infoln("finish restart un async audit log task")
}
