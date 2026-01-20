package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/common/enums"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	"AuditLog/domain/entity/oprlogeo"
	"AuditLog/drivenadapters/httpaccess/usermgnt"
	"AuditLog/errors"
	"AuditLog/gocommon/api"
)

type RolePermissionOption struct {
	IsFromOprLogAPI bool // 是否来自“运营日志”API的角色控制
}

// RolePermission 访问用户角色控制
func RolePermission(allowRoles []string, opt ...RolePermissionOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if helpers.IsAaronLocalDev() {
			c.Set("accountType", enums.AccountTypeUser.String())
			c.Set("userRoles", []string{common.NormalUser})
			// c.Set("userRoles", []string{})
		}

		// 是否来自“运营日志”API的角色控制
		var isFromOprLogAPI bool

		if len(opt) > 0 {
			isFromOprLogAPI = opt[0].IsFromOprLogAPI
		}

		var flag bool

		accountType := c.GetString("accountType")
		if accountType == enums.AccountTypeUser.String() {
			// 访问用户角色控制
			userRoles, _ := c.Get("userRoles")

			var userRolesSlice []string

			if reflect.TypeOf(userRoles).Kind() == reflect.Slice {
				roles := reflect.ValueOf(userRoles)
				for i := 0; i < roles.Len(); i++ {
					role := roles.Index(i)
					userRolesSlice = append(userRolesSlice, role.Interface().(string))
				}
			}

			// 用户不存在或是匿名用户
			if len(userRolesSlice) == 0 {
				if isFromOprLogAPI {
					userType, err := getUserTypeForOprLogAPI(c)
					if err != nil {
						common.ErrResponse(c, err)
						return
					}

					// 如果是匿名用户，直接返回（不检测角色）
					if userType == enums.AnonyUser {
						return
					}
				}
			}

			// 判断用户角色是否有权限
			for _, role := range userRolesSlice {
				_, exist := api.SliceFind(allowRoles, role)
				if exist {
					flag = true
					break
				} else {
					flag = false
					continue
				}
			}

			if !flag {
				err := errors.NewCtx(c, errors.ForbiddenErr, "The user's role has not permission to do this service", "")
				common.ErrResponse(c, err)

				return
			}
		} else if accountType == enums.AccountTypeApp.String() {
			// 判断应用账户是否存在
			umgnt := usermgnt.NewUserMgnt()

			_, statusCode, err := umgnt.GetAppInfoByID(c.GetString("userId"))
			if err != nil {
				err := errors.NewCtx(c, errors.InternalErr, "Failed to connect to user-management service", "")
				common.ErrResponse(c, err)

				return
			}

			if statusCode == http.StatusNotFound {
				err := errors.NewCtx(c, errors.ForbiddenErr, "The app is not exist", "")
				common.ErrResponse(c, err)

				return
			}

			if statusCode != http.StatusOK {
				err := errors.NewCtx(c, errors.ForbiddenErr, fmt.Sprintf("GetAppInfoByID Error, statusCode is %d", statusCode), "")
				common.ErrResponse(c, err)

				return
			}

			return
		} else {
			err := errors.NewCtx(c, errors.ForbiddenErr, "Don't have permission to do this service", "")
			common.ErrResponse(c, err)

			return
		}
	}
}

// getUserTypeForOprLogAPI 从日志中获取用户类型
func getUserTypeForOprLogAPI(c *gin.Context) (userType enums.UserType, err error) {
	// 1. 读取请求体
	originBodyData, err := c.GetRawData()
	if err != nil {
		return
	}

	// 2.统一成数组
	isBatch := strings.Contains(c.Request.URL.Path, "operation-log/batch")
	bodyData := originBodyData

	if !isBatch {
		bodyData = utils.JSONObjectToArray(bodyData)
	}

	// 3.解析请求体 to eos
	var eos []*oprlogeo.LogEntry

	err = utils.JSON().Unmarshal(bodyData, &eos)
	if err != nil {
		return
	}

	if len(eos) == 0 {
		err = errors.NewCtx(c, errors.BadRequestErr, "请求体不能为空", "")
		return
	}

	// 4. 获取用户类型
	userType = eos[0].Operator.Type

	// 5. 将请求体内容重新写回到请求中
	c.Request.Body = io.NopCloser(bytes.NewBuffer(originBodyData))

	return
}
