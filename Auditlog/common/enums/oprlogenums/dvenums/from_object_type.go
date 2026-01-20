package dvenums

type FromObjectType string

func (fot FromObjectType) String() string {
	return string(fot)
}

func (fot FromObjectType) IsDocLibFOT() bool {
	return fot == NormalDirFOT || fot == AssociationFileFOT
}

const (
	DocLibTypeFOT              FromObjectType = "doc_lib_type"             // 库类型
	NormalDirFOT               FromObjectType = "normal_dir"               // 普通目录
	DirStarredFOT              FromObjectType = "dir/starred"              // 文档中心-收藏夹
	DirShareManageFOT          FromObjectType = "dir/sharemanage"          // 文档中心-共享管理
	AssociationFileFOT         FromObjectType = "associationFile"          // 关联文件
	NavDocCenterFOT            FromObjectType = "nav_doc_center"           // 导航-文档中心
	ByLinkFOT                  FromObjectType = "by_link"                  // 通过链接
	AgentMessageFOT            FromObjectType = "agentmessage"             // 待办消息
	AutoSheetsFOT              FromObjectType = "autosheets"               // 表格中心
	AutoSheetsNewFileFOT       FromObjectType = "autosheetsNewFile"        // 表格中心-新建表格
	AutoSheetsNewFormFOT       FromObjectType = "autosheetsNewForm"        // 表格中心-新建表单
	CognitiveAssistantFOT      FromObjectType = "cognitiveAssistant"       // 认知助手
	ContentAutomationFOT       FromObjectType = "content-automation"       // 工作中心
	DirFOT                     FromObjectType = "dir"                      // 文档中心
	DirFilelockFOT             FromObjectType = "dir/filelock"             // 文档中心-文件锁管理
	DirQuarantineFOT           FromObjectType = "dir/quarantine"           // 文档中心-文件隔离区
	DirRecycleFOT              FromObjectType = "dir/recycle"              // 文档中心-回收站
	HomeFOT                    FromObjectType = "home"                     // 首页
	KnowledgeCenterFOT         FromObjectType = "knowledge-center"         // 知识中心
	NotificationMessageFOT     FromObjectType = "notificationmessage"      // 审核消息
	PersonalHomepageFollowFOT  FromObjectType = "personalhomepage/follow"  // 个人主页-我的关注
	PersonalHomepageProfileFOT FromObjectType = "personalhomepage/profile" // 个人主页-个人资料

	PortalFOT             FromObjectType = "portal"              // 门户
	PreviewFOT            FromObjectType = "preview"             // 预览
	SettingsMenuFOT       FromObjectType = "settingsmenu"        // 设置
	ShareLinkAnonymousFOT FromObjectType = "sharelink_anonymous" // 匿名共享页面
	SmartSearchFOT        FromObjectType = "smartSearch"         // 智能搜索
	SyncDiskFOT           FromObjectType = "syncdisk"            // 同步盘
	WorkCenterFOT         FromObjectType = "work-center"         // 应用

	Unknown FromObjectType = "unknown" // 未知
)
