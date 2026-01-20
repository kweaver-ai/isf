package serviceaccess

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	chttp "policy_mgnt/utils/gocommon/v2/http"
	clog "policy_mgnt/utils/gocommon/v2/log"

	"policy_mgnt/common/config"
)

var (
	usermgntOnce      sync.Once
	usermgntSingleton *usermgntDriven
)

// UsermgntDriven usermgnt操作接口
type UsermgntDriven interface {
	// GetBelongDepartByUserId 通过用户id获取所属部门id
	GetBelongDepartByUserId(userId string) (departmentIds []string, err error)
	// GetDepartUserIds 获取部门下所有用户
	GetDepartUserIds(departmentId string) (userIds []string, err error)
}

type usermgntDriven struct {
	PrivateBaseURL string
	HTTPClient     chttp.Client
	logger         clog.Logger
}

// NewUsermgntDriven 新建 usermgnt httpclient 操作对象
func NewUsermgntDriven() UsermgntDriven {
	usermgntOnce.Do(func() {
		usermgntSingleton = &usermgntDriven{
			PrivateBaseURL: config.Config.UserMgmtPvt.Protocol + "://" + config.Config.UserMgmtPvt.Host + ":" + config.Config.UserMgmtPvt.Port + "/api/user-management/v1",
			HTTPClient:     chttp.NewClient(),
			logger:         clog.NewLogger(),
		}
	})
	return usermgntSingleton
}

// GetBelongDepartByUserId 通过用户id获取所属部门id
func (d *usermgntDriven) GetBelongDepartByUserId(userId string) (departmentIds []string, err error) {
	var resp *http.Response
	var respBodyByte []byte
	departUrl := d.PrivateBaseURL + fmt.Sprintf("/users/%v/department_ids", userId)
	resp, err = d.HTTPClient.Get(departUrl)
	if err != nil {
		err = fmt.Errorf("ERROR: GetDepartmentIds:%v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("ERROR: GetDepartmentIds: Response = %v", resp.StatusCode)
		return
	}
	respBodyByte, err = io.ReadAll(resp.Body)
	if err = json.Unmarshal(respBodyByte, &departmentIds); err != nil {
		err = fmt.Errorf("ERROR: GetDepartmentIds: %v", err)
		return
	}
	return
}

// GetDepartUserIds 获取部门下所有用户
func (d *usermgntDriven) GetDepartUserIds(departmentId string) (userIds []string, err error) {
	var resp *http.Response
	var respBodyByte []byte
	departUrl := d.PrivateBaseURL + fmt.Sprintf("/departments/%v/all_user_ids", departmentId)
	resp, err = d.HTTPClient.Get(departUrl)
	if err != nil {
		err = fmt.Errorf("ERROR: GetDepartmentIds:%v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("ERROR: GetDepartmentIds: Response = %v", resp.StatusCode)
		return
	}
	respBodyByte, err = io.ReadAll(resp.Body)
	var allUserIds struct {
		AllUserIds []string `json:"all_user_ids"`
	}
	if err = json.Unmarshal(respBodyByte, &allUserIds); err != nil {
		err = fmt.Errorf("ERROR: GetDepartmentIds: %v", err)
		return
	}
	userIds = allUserIds.AllUserIds
	return
}
