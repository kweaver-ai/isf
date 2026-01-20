// Package logics AnyShare
package logics

import "Authentication/interfaces"

var (
	// DnHydraAdmin 接口实例
	DnHydraAdmin interfaces.DnHydraAdmin
	// DnHydraPublic 接口实例
	DnHydraPublic interfaces.DnHydraPublic
	// DnUserManagement 接口实例
	DnUserManagement interfaces.DnUserManagement
	// DnEacp 接口实例
	DnEacp interfaces.DnEacp
	// DnShareMgnt 接口实例
	DnShareMgnt interfaces.DnShareMgnt
	// DnEacpLog 接口实例
	DnEacpLog interfaces.DnEacpLog
	// DnMessageBroker 接口实例
	DnMessageBroker interfaces.DrivenMessageBroker
)

// SetDnHydraAdmin 设置实例
func SetDnHydraAdmin(i interfaces.DnHydraAdmin) {
	DnHydraAdmin = i
}

// SetDnHydraPublic 设置实例
func SetDnHydraPublic(i interfaces.DnHydraPublic) {
	DnHydraPublic = i
}

// SetDnUserManagement 设置实例
func SetDnUserManagement(i interfaces.DnUserManagement) {
	DnUserManagement = i
}

// SetDnEacp 设置实例
func SetDnEacp(i interfaces.DnEacp) {
	DnEacp = i
}

// SetDnShareMgnt 设置实例
func SetDnShareMgnt(i interfaces.DnShareMgnt) {
	DnShareMgnt = i
}

// SetDnEacpLog 设置实例
func SetDnEacpLog(i interfaces.DnEacpLog) {
	DnEacpLog = i
}

// SetDnMessageBroker 设置实例
func SetDnMessageBroker(i interfaces.DrivenMessageBroker) {
	DnMessageBroker = i
}
