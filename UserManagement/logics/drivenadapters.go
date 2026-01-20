// Package logics AnyShare
package logics

import "UserManagement/interfaces"

var dnEacpLog interfaces.DrivenEacpLog

// SetDnEacpLog 设置实例
func SetDnEacpLog(i interfaces.DrivenEacpLog) {
	dnEacpLog = i
}

var dnMessageBroker interfaces.DrivenMessageBroker

// SetDnMessageBroker 设置实例
func SetDnMessageBroker(i interfaces.DrivenMessageBroker) {
	dnMessageBroker = i
}

var dnHydra interfaces.DrivenHydra

// SetDnHydra 设置实例
func SetDnHydra(i interfaces.DrivenHydra) {
	dnHydra = i
}

var dnOSSGateWay interfaces.DnOSSGateWay

// SetDnOSSGateWay 设置实例
func SetDnOSSGateWay(i interfaces.DnOSSGateWay) {
	dnOSSGateWay = i
}
