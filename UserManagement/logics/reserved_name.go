package logics

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type reservedName struct {
	logger         common.Logger
	dbReservedName interfaces.DBReservedName
	pool           *sqlx.DB
	dbUser         interfaces.DBUser
}

var (
	rnOnce sync.Once
	rn     *reservedName
)

// NewReservedName 创建
func NewReservedName() *reservedName {
	rnOnce.Do(func() {
		rn = &reservedName{
			logger:         common.NewLogger(),
			pool:           dbPool,
			dbReservedName: dbReservedName,
			dbUser:         dbUser,
		}
	})
	return rn
}

// UpdateReservedName 添加、更新保留名称
//
//nolint:gocyclo
func (r *reservedName) UpdateReservedName(name interfaces.ReservedNameInfo) error {
	var err error
	// 检查id是否合法
	reg := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	if !reg.MatchString(name.ID) {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "invalid id")
		return err
	}

	// 获取Unicode字符长度，以字符为单位而不是字节
	nameLen := utf8.RuneCountInString(name.Name)
	// 检查名称是否合法
	if nameLen < 1 || nameLen > 128 {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "'name' length must be between 1 and 128")
		return err
	}

	if strings.HasPrefix(name.Name, " ") || strings.HasSuffix(name.Name, " ") {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "'name' cannot start or end with space")
		return err
	}

	reg = regexp.MustCompile(`[\\/*?"<>|]`)
	if reg.MatchString(name.Name) {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "'name' cannot contain \\ / * ? \" < > |")
		return err
	}

	tx, err := r.pool.Begin()
	if err != nil {
		r.logger.Errorf("UpdateReservedName start transaction error: %v", err)
		return err
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			err = tx.Commit()
			if err != nil {
				r.logger.Errorf("UpdateReservedName Transaction Commit Error:%v", err)
				return
			}
		default:
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				r.logger.Errorf("UpdateReservedName Transaction Rollback Error:%v", rollbackErr)
			}
		}
	}()

	err = r.dbReservedName.GetLock(tx)
	if err != nil {
		r.logger.Errorln("UpdateReservedName GetLock Error")
		return err
	}

	existIDInfo, idExist, err := r.dbReservedName.GetReservedNameByID(name.ID, tx)
	if err != nil {
		r.logger.Errorln("failed to GetReservedNameByID when UpdateReservedName")
		return err
	}

	existNameInfo, nameExist, err := r.dbReservedName.GetReservedNameByName(name.Name, tx)
	if err != nil {
		r.logger.Errorln("failed to GetReservedNameByName when UpdateReservedName")
		return err
	}

	userInfo, err := r.dbUser.GetUserInfoByName(context.Background(), name.Name)
	if err != nil {
		r.logger.Errorf("failed to GetUserInfoByName when UpdateReservedName ,err: %v", err)
		return err
	}

	// ID不为空字符串说明和用户名重复
	if userInfo.ID != "" {
		return rest.NewHTTPErrorV2(rest.Conflict, "name already exists", rest.SetDetail(map[string]interface{}{
			"conflict_object": map[string]interface{}{
				"id":   userInfo.ID,
				"type": "user",
			},
		}))
	}

	if !idExist {
		// 名称和已有的文档库名重复
		if nameExist {
			return rest.NewHTTPErrorV2(rest.Conflict, "name already exists", rest.SetDetail(map[string]interface{}{
				"conflict_object": map[string]interface{}{
					"id":   existNameInfo.ID,
					"type": "doclib",
				},
			}))
		}
		err = r.dbReservedName.AddReservedName(name, tx)
		if err != nil {
			r.logger.Errorln("failed to add reserved name when UpdateReservedName")
		}
		return err
	}

	// 已存在记录，判断是否需要更新

	// 名称相同，无需修改
	if existIDInfo.Name == name.Name {
		return nil
	}

	if nameExist && existNameInfo.ID != name.ID {
		return rest.NewHTTPErrorV2(rest.Conflict, "name already exists", rest.SetDetail(map[string]interface{}{
			"conflict_object": map[string]interface{}{
				"id":   existNameInfo.ID,
				"type": "doclib",
			},
		}))
	}

	// 更新名称
	err = r.dbReservedName.UpdateReservedName(name, tx)
	if err != nil {
		r.logger.Errorln("failed to update reserved name when UpdateReservedName")
	}
	return err
}

// DeleteReservedName 删除保留名称
func (r *reservedName) DeleteReservedName(id string) error {
	// 检查id是否合法
	reg := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	if !reg.MatchString(id) {
		return rest.NewHTTPErrorV2(rest.BadRequest, "invalid id")
	}

	err := r.dbReservedName.DeleteReservedName(id)
	if err != nil {
		r.logger.Errorln("failed to DeleteReservedName")
	}
	return err
}

// GetReservedName 获取保留名称
func (r *reservedName) GetReservedName(name string) (info interfaces.ReservedNameInfo, err error) {
	info, _, err = r.dbReservedName.GetReservedNameByName(name, nil)
	if err != nil {
		r.logger.Errorln("failed to GetReservedNameByName")
	}
	return info, err
}
