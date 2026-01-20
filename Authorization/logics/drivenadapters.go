// Package logics Anyshare 业务逻辑层
package logics

import "Authorization/interfaces"

// dnUserMgnt 实例
var dnUserMgnt interfaces.DrivenUserMgnt

// SetDnUserMgnt 设置实例
func SetDnUserMgnt(i interfaces.DrivenUserMgnt) {
	dnUserMgnt = i
}
