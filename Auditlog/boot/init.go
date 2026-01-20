package boot

import (
	"os"

	// "AuditLog/boot/oprlogboot"
	"AuditLog/common/enums"
	"AuditLog/infra/cmp/langcmp"
	"AuditLog/infra/cmp/logcmp"
	"AuditLog/infra/config"
	"AuditLog/infra/json_schema/jsc_opr_log"
)

func Init() {
	// 1.初始化配置
	config.InitConfig()

	//	2. json schema 相关初始化
	jsc_opr_log.Init()

	// 3. 设置系统默认语言
	defaultLang := os.Getenv(enums.LangEnv)
	if defaultLang == "" {
		defaultLang = string(langcmp.ZhCN)
	}

	langcmp.NewLangCmp().SetSysDefLang(defaultLang)

	//	4. 初始化日志
	logcmp.InitLogger(config.GetLogConf())

	//	5. 运营日志 boot
	// oprlogboot.Init()
}
