package oprlogenums

var AllBizType = []BizType{
	UserLogin, DocOperation, DocFlow, FileCollector, DocumentDomainSync, Antivirus, KcOperation, SensitiveContent, UserFeedback, Sap, Search, ContentProcessFinish, ContentAutomation, IntelliQaChat, ClientOperation, DirVisit, MenuButtonClick,
}

var AllServerBizType = []BizType{
	UserLogin, DocOperation, DocFlow, FileCollector, DocumentDomainSync, Antivirus, KcOperation, SensitiveContent, UserFeedback, Sap, Search, ContentProcessFinish, ContentAutomation, IntelliQaChat,
}

// AllClientBizType 所有的客户端日志模型，客户端上报的业务类型
var AllClientBizType = []BizType{
	DirVisit, MenuButtonClick, ClientOperation,
}
