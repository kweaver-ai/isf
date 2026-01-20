package logics

const (
	// 使用iota，增加新的枚举时，必须按照顺序添加，不得在中间插入
	// 否则会改变枚举的值，造成outbox handler与db中记录的type对应不上
	_ = iota
	// outboxProductAuthorizedAddedLog 产品授权新增记录日志
	outboxProductAuthorizedAddedLog
	// outboxProductAuthorizedUpdatedLog 产品授权更新记录日志
	outboxProductAuthorizedUpdatedLog
	// outboxProductAuthorizedDeletedLog 产品授权删除记录日志
	outboxProductAuthorizedDeletedLog
)

const (
	_ = iota
	// i18nIDUserHasNoAuthUserProduct 用户暂未获得产品授权
	i18nIDUserHasNoAuthUserProduct
	// i18nIDHasNoLicense 产品无有效授权
	i18nIDHasNoLicense
	// i18nIDProductAuthorizedNotValid 产品授权无效
	i18nIDProductAuthorizedNotValid
	// i18nIDProductAuthorizedOverQuota 产品授权超过限额
	i18nIDProductAuthorizedOverQuota
)

func RemoveDuplicate(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
