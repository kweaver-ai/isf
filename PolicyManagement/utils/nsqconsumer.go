package utils

import (
	"encoding/json"
	"fmt"
	"sync"

	"policy_mgnt/utils/gocommon/api"

	msqclient "github.com/kweaver-ai/proton-mq-sdk-go"
)

var (
	mqClient msqclient.ProtonMQClient
	mqOnce   sync.Once
)

func initMQClient() msqclient.ProtonMQClient {
	mqOnce.Do(func() {
		mqSDK, err := msqclient.NewProtonMQClientFromFile("/sysvol/conf/mq_config.yaml")
		if err != nil {
			panic(fmt.Sprintf("ERROR: new mq client failed: %v\n", err))
		}
		mqClient = mqSDK
	})
	return mqClient
}

type MQHandler interface {
	HandleMessage(msg []byte) error
}

// 引擎nsq传入message格式
type msgBodyDelete struct {
	ID string `json:"id"`
}

// 添加用户到部门message格式
type MsgDepartAddUser struct {
	ID          string   `json:"id"`
	DepartPaths []string `json:"dept_paths"`
}

// 用户移动message格式
type MsgUserMoved struct {
	ID            string `json:"id"`
	OldDepartPath string `json:"old_dept_path"`
	NewDepartPath string `json:"new_dept_path"`
}

// 添加用户到部门message格式
type MsgDepartRemoveUser struct {
	ID          string   `json:"id"`
	DepartPaths []string `json:"dept_paths"`
}

// TaskUserHandler 处理引擎消息的handler
type TaskUserHandler struct {
	DeleteUserChan chan string
}

// HandleMessage 是需要实现的处理消息的方法
func (h *TaskUserHandler) HandleMessage(msg []byte) error {
	l := api.NewLogger()
	l.Infof("Recv msg:%v\n", string(msg))

	var msgBody msgBodyDelete
	if err := json.Unmarshal(msg, &msgBody); err != nil {
		l.Errorf("Get wrong formate message:%v, err:%v\n", string(msg), err)
		return nil
	}

	h.DeleteUserChan <- msgBody.ID
	return nil
}

// TaskDepartmentHandler 处理引擎消息的handler
type TaskDepartmentHandler struct {
	DeleteDepartmentChan chan string
}

// HandleMessage 是需要实现的处理消息的方法
func (h *TaskDepartmentHandler) HandleMessage(msg []byte) error {
	l := api.NewLogger()
	l.Infof("Recv msg:%v\n", string(msg))

	var msgBody msgBodyDelete
	if err := json.Unmarshal(msg, &msgBody); err != nil {
		l.Errorf("ERROR: Get wrong formate message:%v, err:%v\n", string(msg), err)
		return nil
	}

	// 写入用户id
	h.DeleteDepartmentChan <- msgBody.ID
	return nil
}

// TaskDepartmentAddUserHandler 处理引擎消息的handler
type TaskDepartmentAddUserHandler struct {
	DepartmentAddUserChan chan MsgDepartAddUser
}

// HandleMessage 是需要实现的处理消息的方法
func (h *TaskDepartmentAddUserHandler) HandleMessage(msg []byte) error {
	l := api.NewLogger()
	l.Infof("INFO: Recv msg:%v\n", string(msg))

	var msgBody MsgDepartAddUser
	if err := json.Unmarshal(msg, &msgBody); err != nil {
		l.Errorf("ERROR: Get wrong format message:%v, err:%v\n", string(msg), err)
		return nil
	}
	h.DepartmentAddUserChan <- msgBody
	return nil
}

// TaskUserMovedHandler 处理引擎消息的handler
type TaskUserMovedHandler struct {
	UserMovedChan chan MsgUserMoved
}

// HandleMessage 是需要实现的处理消息的方法
func (h *TaskUserMovedHandler) HandleMessage(msg []byte) error {
	l := api.NewLogger()
	l.Infof("INFO: Recv msg:%v\n", string(msg))

	var msgBody MsgUserMoved
	if err := json.Unmarshal(msg, &msgBody); err != nil {
		l.Errorf("ERROR: Get wrong format message:%v, err:%v\n", string(msg), err)
		return nil
	}

	// 写入文档库id,name
	h.UserMovedChan <- msgBody
	return nil
}

// TaskDepartmentAddUserHandler 处理引擎消息的handler
type TaskDepartmentRemoveUserHandler struct {
	DepartmentRemoveUserChan chan MsgDepartRemoveUser
}

// HandleMessage 是需要实现的处理消息的方法
func (h *TaskDepartmentRemoveUserHandler) HandleMessage(msg []byte) error {
	l := api.NewLogger()
	l.Infof("INFO: Recv msg:%v\n", string(msg))

	var msgBody MsgDepartRemoveUser
	if err := json.Unmarshal(msg, &msgBody); err != nil {
		l.Errorf("ERROR: Get wrong format message:%v, err:%v\n", string(msg), err)
		return nil
	}
	h.DepartmentRemoveUserChan <- msgBody
	return nil
}

// 初始化sharemgnt消费者
func InitSharemgntConsumer(userDeleteChan, deptDeleteChan chan string, departmentAddUserChan chan MsgDepartAddUser, userMovedChan chan MsgUserMoved, departmentRemoveUserChan chan MsgDepartRemoveUser) {

	// 初始化mq连接
	initMQClient()

	topicDelteUser := "core.user.delete"
	topicDelteDepartment := "core.dept.delete"
	topicDepartmentAddUser := "user_management.department.user.added"
	topicUserMoved := "user_management.user.moved"
	topicDepartmentRemoveUser := "user_management.department.user.removed"
	channel := "sharemgnt"
	pollIntervalMilliseconds := 100
	maxInFlight := 16
	l := api.NewLogger()

	uHandler := &TaskUserHandler{
		DeleteUserChan: userDeleteChan,
	}
	go func() {
		if err := mqClient.Sub(topicDelteUser, channel, uHandler.HandleMessage, int64(pollIntervalMilliseconds), maxInFlight); err != nil {
			l.Errorf("ERROR: Topic: %v, message: %v", topicDelteUser, err)
		}
	}()

	dHandler := &TaskDepartmentHandler{
		DeleteDepartmentChan: deptDeleteChan,
	}
	go func() {
		if err := mqClient.Sub(topicDelteDepartment, channel, dHandler.HandleMessage, int64(pollIntervalMilliseconds), maxInFlight); err != nil {
			l.Errorf("ERROR: Topic: %v, message: %v", topicDelteDepartment, err)
		}
	}()

	departmentAddUserHandler := &TaskDepartmentAddUserHandler{
		DepartmentAddUserChan: departmentAddUserChan,
	}
	go func() {
		if err := mqClient.Sub(topicDepartmentAddUser, channel, departmentAddUserHandler.HandleMessage, int64(pollIntervalMilliseconds), maxInFlight); err != nil {
			l.Errorf("ERROR: Topic: %v, message: %v", topicDepartmentAddUser, err)
		}
	}()

	lUpdateHandler := &TaskUserMovedHandler{
		UserMovedChan: userMovedChan,
	}
	go func() {
		if err := mqClient.Sub(topicUserMoved, channel, lUpdateHandler.HandleMessage, int64(pollIntervalMilliseconds), maxInFlight); err != nil {
			l.Errorf("ERROR: Topic: %v, message: %v", topicUserMoved, err)
		}
	}()

	departmentRemoveUserHandler := &TaskDepartmentRemoveUserHandler{
		DepartmentRemoveUserChan: departmentRemoveUserChan,
	}
	go func() {
		if err := mqClient.Sub(topicDepartmentRemoveUser, channel, departmentRemoveUserHandler.HandleMessage, int64(pollIntervalMilliseconds), maxInFlight); err != nil {
			l.Errorf("ERROR: Topic: %v, message: %v", topicDepartmentRemoveUser, err)
		}
	}()

}
