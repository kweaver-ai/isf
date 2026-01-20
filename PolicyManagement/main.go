package main

import (
	"context"

	"policy_mgnt/apiserver"
	"policy_mgnt/common"
	"policy_mgnt/common/config"
	"policy_mgnt/general"
	"policy_mgnt/infra/log"
	"policy_mgnt/infra/outbox"
	"policy_mgnt/network"
	"policy_mgnt/utils"
	cdb "policy_mgnt/utils/gocommon/v2/db"
	cutils "policy_mgnt/utils/gocommon/v2/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := utils.InitConfig(); err != nil {
		panic(err)
	}
	config.InitConfig()

	if err := cdb.InitGormDB(&config.Config.DB); err != nil {
		panic(err)
	}

	// 初始化SonyFlake
	if err := cutils.InitSonyflake(config.Config.PodIP); err != nil {
		panic(err)
	}

	// 初始化日志配置
	log.InitLogger()

	// 初始化Outbox
	if err := outbox.InitOutBoxer(); err != nil {
		panic(err)
	}

	if err := general.CreateDefaultPolicy(); err != nil {
		panic(err)
	}
	if err := network.CreatePublicNet(); err != nil {
		panic(err)
	}
	if err := outbox.StartOutBoxer(context.Background()); err != nil {
		panic(err)
	}

	// 初始化ARTrace
	common.InitARTrace("policy-management")

	// 初始化log
	common.NewLogger().SetLevel(int(logrus.InfoLevel))

	// 启动MQ处理器
	mqHandler := apiserver.NewMQHandler()
	mqHandler.Subscribe()

	// 启动API服务器
	apiserver.StartAPIServer()
}
