package logics

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	apiErr "AuditLog/errors"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	"AuditLog/models"

	jsoniter "github.com/json-iterator/go"
)

var (
	outboxHandlers map[string]func(interface{}) error = make(map[string]func(interface{}) error)
	outboxThreads  map[string]*outbox                 = make(map[string]*outbox)
)

type outbox struct {
	db       interfaces.DBOutbox
	pushChan chan struct{}
	logger   api.Logger
	pool     *sqlx.DB
	// 打印日志使用
	businessType string
}

// NewOutbox 创建新的outbox对象
func NewOutbox(busType string) interfaces.Outbox {
	// 判断业务类型outbox线程是否存在，若存在直接返回
	ob, ok := outboxThreads[busType]
	if !ok {
		ob = &outbox{
			db:           dbOutbox,
			pushChan:     make(chan struct{}, 1),
			logger:       logger,
			pool:         dbPool,
			businessType: busType,
		}
		go ob.startPushOutboxThread()
		ob.NotifyPushOutboxThread()
		outboxThreads[busType] = ob
	}

	return ob
}

func (o *outbox) RegisterHandlers(opType string, op func(interface{}) error) {
	_, ok := outboxHandlers[opType]
	if ok {
		o.logger.Errorf("RegisterHandler: OpType Exist: %v", opType)
	}
	outboxHandlers[opType] = op
}

// 添加outbox消息
func (o *outbox) AddOutboxInfo(opType string, content interface{}, tx *sql.Tx) error {
	msg := models.OutboxMsg{
		Type:    opType,
		Content: content,
	}
	return o.AddOutboxInfos([]models.OutboxMsg{msg}, tx)
}

// 批量添加outbox消息
func (o *outbox) AddOutboxInfos(msgs []models.OutboxMsg, tx *sql.Tx) error {
	outboxMsgs := make([]string, 0, len(msgs))
	for i := range msgs {
		_, ok := outboxHandlers[msgs[i].Type]
		if !ok {
			return apiErr.New("", apiErr.InternalErr, "Internal Error", fmt.Sprintf("AddOutboxInfo: OpType Not Exist: %v", msgs[i].Type))
		}
		outboxMsg, _ := jsoniter.MarshalToString(msgs[i])
		outboxMsgs = append(outboxMsgs, outboxMsg)
	}
	return o.db.AddOutboxInfos(o.businessType, outboxMsgs, tx)
}

// notify outbox推送线程
func (o *outbox) NotifyPushOutboxThread() {
	select {
	case o.pushChan <- struct{}{}:
	default:
	}
}

// 开启推送outbox消息线程
func (o *outbox) startPushOutboxThread() {
	// 出错时每隔30s重试，直至成功
	const timeDuration int = 30
	ticker := time.NewTicker(time.Second * time.Duration(timeDuration))
	needTimer := false
	for {
		select {
		case <-ticker.C:
			// 无需定时器触发
			if !needTimer {
				continue
			}
		case <-o.pushChan:
		}

		// 逐条进行推送
		for {
			var isAllFinished bool
			needTimer, isAllFinished = o.push()
			// 出错或消息都推送完毕时退出循环
			if needTimer || isAllFinished {
				break
			}
		}
	}
}

// 推送一条outbox消息
func (o *outbox) push() (needTimer, isAllFinished bool) {
	// 标记是否需要定时推送
	needTimer = false
	// 标记是否全部推送完成
	isAllFinished = false
	// 获取DB事务
	o.logger.Infof("Info : start push outbox Transaction %v", o.businessType)
	tx, err := o.pool.Begin()
	if err != nil {
		needTimer = true
		o.logger.Errorf("PushOutbox %v：Get Transaction Failed:%v", o.businessType, err)
		return
	}
	// 异常时需要定时推送
	defer func() {
		switch err {
		case nil:
			// 提交事务
			err = tx.Commit()
			if err != nil {
				needTimer = true
				isAllFinished = false
				o.logger.Errorf("PushOutbox：Transaction Commit Failed:%v", err)
			}
		default:
			needTimer = true
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				o.logger.Errorf("PushOutbox：Transaction Rollback Error:%v", rollbackErr)
			}
		}
	}()

	// 获取推送消息
	messageID, message, err := o.db.GetPushMessage(o.businessType, tx)
	if err != nil {
		o.logger.Errorf("PushOutbox：Get PushMessage Failed:%v", err)
		return
	}
	// 若为空说明所有outbox消息均已推送
	if message == "" {
		o.logger.Infof("Info : finish push all businessType %v messages", o.businessType)
		isAllFinished = true
		return
	}

	// 进行推送
	o.logger.Infof("Info : start push outbox message: f_id is %d, businessType is %v", messageID, o.businessType)
	var messageJSON interface{}
	if err = jsoniter.UnmarshalFromString(message, &messageJSON); err != nil {
		o.logger.Warnf("PushOutbox：jsoniter UnmarshalFromString Failed:%v", err)
		return
	}
	messageType := messageJSON.(map[string]interface{})["type"].(string)
	content := messageJSON.(map[string]interface{})["content"]
	handler, ok := outboxHandlers[messageType]
	if !ok {
		o.logger.Warnf("PushOutbox：msg type error:%v", err)
		err = errors.New("msg type is wrong")
		return
	}
	err = handler(content)
	if err != nil {
		o.logger.Errorf("PushOutbox：Publish outbox info Failed:%v", err)
		return
	}

	// 推送成功，删除outbox表记录
	err = o.db.DeleteOutboxInfoByID(messageID, tx)
	if err != nil {
		o.logger.Errorf("Push outbox info err:%v", err)
		return
	}

	return
}
