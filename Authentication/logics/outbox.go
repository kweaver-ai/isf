// Package logics outbox Anyshare 业务逻辑层 -outbox发件箱
package logics

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"github.com/kweaver-ai/go-lib/rest"
	jsoniter "github.com/json-iterator/go"

	"Authentication/common"
	"Authentication/interfaces"
)

var (
	outboxHandlers map[int]func(interface{}) error = make(map[int]func(interface{}) error)
	outboxThreads  map[int]*outbox                 = make(map[int]*outbox)
)

type outbox struct {
	db           interfaces.DBOutbox
	pushChan     chan struct{}
	logger       common.Logger
	pool         *sqlx.DB
	businessType int
}

// NewOutbox 创建新的outbox对象
func NewOutbox(busType int) *outbox {
	// 判断业务类型outbox线程是否存在，若存在直接返回
	ob, ok := outboxThreads[busType]
	if !ok {
		ob = &outbox{
			db:           DBOutbox,
			pushChan:     make(chan struct{}, 1),
			logger:       common.NewLogger(),
			pool:         DBPool,
			businessType: busType,
		}
		go ob.startPushOutboxThread()
		ob.NotifyPushOutboxThread()
		outboxThreads[busType] = ob
	}

	return ob
}

func (o *outbox) RegisterHandlers(opType int, op func(interface{}) error) {
	_, ok := outboxHandlers[opType]
	if ok {
		o.logger.Fatalf("RegisterHandler: OpType Exist: %v", opType)
	}
	outboxHandlers[opType] = op
}

// 添加outbox消息
func (o *outbox) AddOutboxInfo(opType int, content interface{}, tx *sql.Tx) error {
	msg := interfaces.OutboxMsg{
		Type:    opType,
		Content: content,
	}
	return o.AddOutboxInfos([]interfaces.OutboxMsg{msg}, tx)
}

// 批量添加outbox消息
func (o *outbox) AddOutboxInfos(msgs []interfaces.OutboxMsg, tx *sql.Tx) error {
	outboxMsgs := make([]string, 0, len(msgs))
	for i := range msgs {
		_, ok := outboxHandlers[msgs[i].Type]
		if !ok {
			return rest.NewHTTPError(fmt.Sprintf("AddOutboxInfo: OpType Not Exist: %v", msgs[i].Type), rest.InternalServerError, nil)
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
		o.logger.Infof("Info : finish push all businessType %d messages", o.businessType)
		isAllFinished = true
		return
	}

	// 进行推送
	o.logger.Infof("Info : start push outbox message: f_id is %d, businessType is %d", messageID, o.businessType)
	var messageJSON interface{}
	if err = jsoniter.UnmarshalFromString(message, &messageJSON); err != nil {
		o.logger.Errorf("PushOutbox：jsoniter UnmarshalFromString Failed:%v", err)
		return
	}
	messageType := int(messageJSON.(map[string]interface{})["type"].(float64))
	content := messageJSON.(map[string]interface{})["content"]
	handler, ok := outboxHandlers[messageType]
	if !ok {
		o.logger.Errorf("PushOutbox：msg type error:%v", err)
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
