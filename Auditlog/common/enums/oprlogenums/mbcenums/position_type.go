package mbcenums

type PositionType string

func (pot PositionType) String() string {
	return string(pot)
}

const (
	DirStarredPoT              PositionType = "dir/starred"              // 文档中心-收藏夹
	DirShareManagePoT          PositionType = "dir/sharemanage"          // 文档中心-共享管理
	AgentMessagePoT            PositionType = "agentmessage"             // 待办消息
	AutoSheetsPoT              PositionType = "autosheets"               // 表格中心
	AutoSheetsNewFilePoT       PositionType = "autosheetsNewFile"        // 表格中心-新建表格
	AutoSheetsNewFormPoT       PositionType = "autosheetsNewForm"        // 表格中心-新建表单
	CognitiveAssistantPoT      PositionType = "cognitiveAssistant"       // 认知助手
	ContentAutomationPoT       PositionType = "content-automation"       // 工作中心
	DirPoT                     PositionType = "dir"                      // 文档中心 (文档列表页面)
	DirFilelockPoT             PositionType = "dir/filelock"             // 文档中心-文件锁管理
	DirQuarantinePoT           PositionType = "dir/quarantine"           // 文档中心-文件隔离区
	DirRecyclePoT              PositionType = "dir/recycle"              // 文档中心-回收站
	DocLibTypePoT              PositionType = "doc_lib_type"             // 库类型
	HomePoT                    PositionType = "home"                     // 首页
	KnowledgeCenterPoT         PositionType = "knowledge-center"         // 知识中心
	NotificationMessagePoT     PositionType = "notificationmessage"      // 审核消息
	PersonalHomepageFollowPoT  PositionType = "personalhomepage/follow"  // 个人主页-我的关注
	PersonalHomepageProfilePoT PositionType = "personalhomepage/profile" // 个人主页-个人资料

	PortalPoT             PositionType = "portal"              // 门户
	PreviewPoT            PositionType = "preview"             // 预览
	SettingsMenuPoT       PositionType = "settingsmenu"        // 设置
	ShareLinkAnonymousPoT PositionType = "sharelink_anonymous" // 匿名共享页面
	SmartSearchPoT        PositionType = "smartSearch"         // 智能搜索
	SyncDiskPoT           PositionType = "syncdisk"            // 同步盘
	WorkCenterPoT         PositionType = "work-center"         // 应用
)
