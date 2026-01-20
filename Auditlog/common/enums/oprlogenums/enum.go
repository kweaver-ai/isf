package oprlogenums

const (
	UserLogin            BizType = "user_login"             // 用户登录
	DocOperation         BizType = "doc_operation"          // 文档操作
	DocFlow              BizType = "doc_flow"               // 文档流转
	FileCollector        BizType = "file_collector"         // 文档收集
	DocumentDomainSync   BizType = "document_domain_sync"   // 文档域同步
	Antivirus            BizType = "antivirus"              // 杀毒
	KcOperation          BizType = "kc_operation"           // 知识中心的操作记录
	SensitiveContent     BizType = "sensitive_content"      // 敏感文件
	UserFeedback         BizType = "user_feedback"          // 用户反馈
	Sap                  BizType = "sap"                    // SAP
	Search               BizType = "search"                 // 搜索
	ContentProcessFinish BizType = "content_process_finish" // 内容处理
	ContentAutomation    BizType = "content_automation"     // 工作中心
	IntelliQaChat        BizType = "intelli_qa_chat"        // 认知助手
	ClientOperation      BizType = "client_operation"       // 客户端日志模型
	DirVisit             BizType = "dir_visit"              // 目录访问
	MenuButtonClick      BizType = "menu_button_click"      // 菜单按钮点击
)
