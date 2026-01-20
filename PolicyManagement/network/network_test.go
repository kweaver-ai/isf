package network

import (
	"strings"
	"testing"

	"policy_mgnt/tapi/ethriftexception"
	"policy_mgnt/tapi/sharemgnt"
	"policy_mgnt/test"
	"policy_mgnt/test/mock_dependency"
	"policy_mgnt/test/mock_thrift"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func mockMgnt(t *testing.T) (*gomock.Controller, *mock_thrift.MockShareMgnt, *mock_dependency.MockAbstractDriven) {
	ctrl := gomock.NewController(t)
	mgnt := mock_thrift.NewMockShareMgnt(ctrl)
	driven := mock_dependency.NewMockAbstractDriven(ctrl)
	return ctrl, mgnt, driven
}

func addTestNetwork() (ids []string) {
	var net1 = models.NetworkRestriction{
		Name:    "existedName1",
		StartIP: "1.1.1.1",
		EndIP:   "2.1.1.1",
		Type:    "ip_segment",
		IpType:  "ipv4",
	}

	var net2 = models.NetworkRestriction{
		Name:      "existedName2",
		IPAddress: "1.1.1.1",
		IPMask:    "0.0.0.0",
		Type:      "ip_mask",
		IpType:    "ipv4",
	}

	var net3 = models.NetworkRestriction{
		Name:      "existedName3",
		IPAddress: "5.1.1.1",
		IPMask:    "0.0.0.0",
		Type:      "ip_mask",
		IpType:    "ipv4",
	}
	var net4 = models.NetworkRestriction{
		Name:    "existedName4",
		StartIP: "2001:db8:0:1::101",
		EndIP:   "2001:db8:0:1::104",
		Type:    "ip_segment",
		IpType:  "ipv6",
	}
	m, _ := NewManagement()
	netID1, _ := m.AddNetwork(&net1)
	netID2, _ := m.AddNetwork(&net2)
	netID3, _ := m.AddNetwork(&net3)
	netID4, _ := m.AddNetwork(&net4)
	ids = append(ids, netID1, netID2, netID3, netID4)
	return
}

func addTestAccessor(netID string) {
	db, _ := api.ConnectDB()

	r1 := models.NetworkAccessorRelation{
		NetworkId:    netID,
		AccessorId:   "user1",
		AccessorType: "user",
	}
	r2 := models.NetworkAccessorRelation{
		NetworkId:    netID,
		AccessorId:   "department1",
		AccessorType: "department",
	}

	db.Save(&r1)
	db.Save(&r2)
}

// 检查网段参数
func TestTrimLeftZero(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"01.0001.010.01", "1.1.10.1"},
		{"0001.010.01", "1.10.1"},
		{"01", "1"},
		{"ssss1", "ssss1"},
		{"10.0.000.1", "10.0.0.1"},
		{"000", "0"},
	}

	for _, test := range tests {
		res := trimLeftZero(test.input)
		assert.Equal(t, test.want, res)
	}
}

// 去掉网段左边的0
func TestTrimNet(t *testing.T) {
	Convey("去除成功", t, func() {
		input := &models.NetworkRestriction{
			StartIP:   "10.0.0.01",
			EndIP:     "0.000.0.01",
			IPAddress: "00.0.0.010",
			IPMask:    "10.0.0.0",
		}
		want := models.NetworkRestriction{
			StartIP:   "10.0.0.1",
			EndIP:     "0.0.0.1",
			IPAddress: "0.0.0.10",
			IPMask:    "10.0.0.0",
		}
		trimNet(input)
		assert.Equal(t, want, *input)
	})
}

// 检查网段参数
func TestCheckNetworkParams(t *testing.T) {
	getString := func(n int, s string) string {
		var res []string
		for i := 0; i < n; i++ {
			res = append(res, s)
		}
		return strings.Join(res, "")
	}

	var tests = []struct {
		input   models.NetworkRestriction
		wanterr error
	}{
		{models.NetworkRestriction{Name: "", StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, nil},
		{models.NetworkRestriction{Name: getString(128, "a"), StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, nil},
		{models.NetworkRestriction{Name: getString(128, "我"), StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, nil},
		{models.NetworkRestriction{Name: `\`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`\`}}})},
		{models.NetworkRestriction{Name: `/`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`/`}}})},
		{models.NetworkRestriction{Name: `:`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`:`}}})},
		{models.NetworkRestriction{Name: `*`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`*`}}})},
		{models.NetworkRestriction{Name: `?`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`?`}}})},
		{models.NetworkRestriction{Name: `"`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`"`}}})},
		{models.NetworkRestriction{Name: `<`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`<`}}})},
		{models.NetworkRestriction{Name: `>`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`>`}}})},
		{models.NetworkRestriction{Name: `|`, StartIP: "0.0.0.0", EndIP: "0.0.0.0", Type: "ip_segment", IpType: "ipv4"}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{`|`}}})},
		{models.NetworkRestriction{Name: getString(129, "a")}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{getString(129, "a")}}})},
		{models.NetworkRestriction{Name: getString(129, "你")}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{getString(129, "你")}}})},
		{models.NetworkRestriction{Name: getString(12, ":")}, errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{getString(12, ":")}}})},
		{models.NetworkRestriction{}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"net_type", "ip_type"}}})},
		{models.NetworkRestriction{Type: "wrong"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"net_type", "ip_type"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "255.255.255.255", EndIP: "255.255.255.255", Type: "ip_segment", IpType: "ipv4"}, nil},
		{models.NetworkRestriction{Name: "", StartIP: "66.1.1.1", EndIP: "66.1.1.1", Type: "ip_segment", IpType: "ipv4"}, nil},
		{models.NetworkRestriction{Name: "", StartIP: "1", Type: "ip_segment", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "266.1.1.1", Type: "ip_segment", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "155.155.155.155", Type: "ip_segment", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "255.255.255.255", EndIP: "255.255.255.355", Type: "ip_segment", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "66.1.1.1", EndIP: "-1", Type: "ip_segment", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "255.255.255.255", EndIP: "155.255.255.255", Type: "ip_segment", IpType: "ipv4"},
			errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip", "end_ip"}}, Cause: "Start_ip can not greater than end_ip."})},
		{models.NetworkRestriction{Name: "", Type: "ip_mask"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"net_type", "ip_type"}}})},
		{models.NetworkRestriction{Name: "", IPAddress: "255.255.255.255", IPMask: "255.255.255.255", Type: "ip_mask", IpType: "ipv4"}, nil},
		{models.NetworkRestriction{Name: "", IPAddress: "255.255.255.255", IPMask: "0.0.0.0", Type: "ip_mask", IpType: "ipv4"}, nil},
		{models.NetworkRestriction{Name: "", IPAddress: "555.255.255.255", IPMask: "255.255.255.255", Type: "ip_mask", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address"}}})},
		{models.NetworkRestriction{Name: "", IPAddress: "255", IPMask: "255.255.255.255", Type: "ip_mask", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address"}}})},
		{models.NetworkRestriction{Name: "", IPAddress: "255.255.255.255", IPMask: "255", Type: "ip_mask", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"netmask"}}})},
		{models.NetworkRestriction{Name: "", IPAddress: "255.255.255.255", IPMask: "255.55.255.255", Type: "ip_mask", IpType: "ipv4"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"netmask"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "2001:db8:0:1::100", EndIP: "2001:db8:0:1::101", Type: "ip_segment", IpType: "ipv6"}, nil},
		{models.NetworkRestriction{Name: "", StartIP: "1", Type: "ip_segment", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "123443", Type: "ip_segment", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "2001:db8:0:1::100", Type: "ip_segment", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "2001:db8:0:1::100", EndIP: "sssss", Type: "ip_segment", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})},
		{models.NetworkRestriction{Name: "", StartIP: "2001:db8:0:1::100", EndIP: "-1", Type: "ip_segment", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})},
		{models.NetworkRestriction{Name: "", IPAddress: "312412351325", IPMask: "123", Type: "ip_mask", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address", "netmask"}}})},
		{models.NetworkRestriction{Name: "", IPAddress: "2001:db8:0:1::100", IPMask: "255", Type: "ip_mask", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address", "netmask"}}})},
		{models.NetworkRestriction{Name: "", IPAddress: "2001:db8:0:1::100", IPMask: "2.4", Type: "ip_mask", IpType: "ipv6"}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address", "netmask"}}})},
	}

	m := &mgnt{}
	for _, test := range tests {
		err := m.checkNetworkParams(&test.input)
		assert.Equal(t, test.wanterr, err)
	}
}

// 测试添加网段
func TestAddNetwork(t *testing.T) {
	m := &mgnt{}
	teardown := test.SetUpDB(t)
	defer teardown(t)

	addTestNetwork()

	// case: error situations
	var tests = []struct {
		input *models.NetworkRestriction
		err   error
	}{
		{&models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", Type: "ip_segment", IpType: "ipv4"}, nil},
		{&models.NetworkRestriction{Name: "existedName1", StartIP: "1.1.1.2", EndIP: "2.1.1.1", Type: "ip_segment", IpType: "ipv4"},
			errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"name"}}})},
		{&models.NetworkRestriction{Name: "dupSegNet", StartIP: "1.1.1.1", EndIP: "2.1.1.1", Type: "ip_segment", IpType: "ipv4"},
			errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"start_ip", "end_ip"}}})},
		{&models.NetworkRestriction{Name: "dupipmask", IPAddress: "1.1.1.1", IPMask: "0.0.0.0", Type: "ip_mask", IpType: "ipv4"},
			errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"ip_address", "netmask"}}})},
		{&models.NetworkRestriction{Name: "22", StartIP: "2001:db8:0:1::100", EndIP: "2001:db8:0:1::101", Type: "ip_segment", IpType: "ipv6"}, nil},
	}

	for _, test := range tests {
		_, err := m.AddNetwork(test.input)
		assert.Equal(t, test.err, err)
	}

	// case: comepare data with db
	var net = models.NetworkRestriction{
		Name:      "test",
		IPAddress: "5.1.1.1",
		IPMask:    "255.255.0.0",
		Type:      "ip_mask",
		IpType:    "ipv4",
	}
	netID, err := m.AddNetwork(&net)
	assert.Equal(t, err, nil)

	db, _ := api.ConnectDB()
	var dbNet models.NetworkRestriction
	db.Where("f_id = ?", netID).First(&dbNet)
	assert.Equal(t, net.Name, dbNet.Name)
	assert.Equal(t, net.IPAddress, dbNet.IPAddress)
	assert.Equal(t, net.IPMask, dbNet.IPMask)
	// ip to int
	assert.Equal(t, "83951617", dbNet.SegmentStart)
	assert.Equal(t, "84017150", dbNet.SegmentEnd)
}

// 测试ipv4转换
func TestFillIPV4(t *testing.T) {
	var tests = []struct {
		input   *models.NetworkRestriction
		ipSeg   []string
		wanterr error
	}{
		{&models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", Type: "ip_segment"}, []string{"167837954", "335610113"}, nil},
		{&models.NetworkRestriction{Name: "ok", IPAddress: "10.1.1.2", IPMask: "255.255.255.0", Type: "ip_mask"}, []string{"167837953", "167838206"}, nil},
		{&models.NetworkRestriction{Name: "ok", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_mask"}, []string{"1", "131070"}, nil},
	}

	m := &mgnt{}
	for _, test := range tests {
		err := m.fillIntIP(test.input)
		if err != nil {
			assert.Equal(t, test.wanterr, err)
		}
		assert.Equal(t, test.ipSeg[0], test.input.SegmentStart)
		assert.Equal(t, test.ipSeg[1], test.input.SegmentEnd)
	}
}

// 测试ipv6转换
func TestFillIPV6(t *testing.T) {
	var tests = []struct {
		input   *models.NetworkRestriction
		ipSeg   []string
		wanterr error
	}{
		{&models.NetworkRestriction{Name: "ok", StartIP: "2001:db8:0:1::100", EndIP: "2001:db8:0:1::101", Type: "ip_segment"}, []string{"42540766411282592875350729025363378432", "42540766411282592875350729025363378433"}, nil},
		{&models.NetworkRestriction{Name: "ok", IPAddress: "2001:db8:0:1::101", IPMask: "12", Type: "ip_mask"}, []string{"42535295865117307932921825928971026432", "42618372614853865174978313870238547967"}, nil},
	}

	m := &mgnt{}
	for _, test := range tests {
		err := m.fillIntIPV6(test.input)
		if err != nil {
			assert.Equal(t, test.wanterr, err)
		}
		assert.Equal(t, test.ipSeg[0], test.input.SegmentStart)
		assert.Equal(t, test.ipSeg[1], test.input.SegmentEnd)
	}
}

// 测试根据网段id获取网段信息
func TestGetNetworkByID(t *testing.T) {
	m := &mgnt{}
	teardown := test.SetUpDB(t)
	defer teardown(t)
	// 添加测试用网段
	ids := addTestNetwork()

	var tests = []struct {
		input string
		data  models.WebNetworkRestriction
		err   error
	}{
		{ids[0], models.WebNetworkRestriction{
			NetworkRestriction: &models.NetworkRestriction{
				Name:    "existedName1",
				StartIP: "1.1.1.1",
				EndIP:   "2.1.1.1",
				Type:    "ip_segment",
				IpType:  "ipv4",
			},
			SegmentStart: nil,
			SegmentEnd:   nil,
			CreatedAt:    nil,
		}, nil},
		{ids[1], models.WebNetworkRestriction{
			NetworkRestriction: &models.NetworkRestriction{
				Name:      "existedName2",
				IPAddress: "1.1.1.1",
				IPMask:    "0.0.0.0",
				Type:      "ip_mask",
				IpType:    "ipv4",
			},
			SegmentStart: nil,
			SegmentEnd:   nil,
			CreatedAt:    nil,
		}, nil},
		{"not_existed_id", models.WebNetworkRestriction{}, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})},
	}

	for _, test := range tests {
		res, err := m.GetNetworkByID(test.input)
		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.data.Name, res.Name)
			assert.Equal(t, test.data.StartIP, res.StartIP)
			assert.Equal(t, test.data.EndIP, res.EndIP)
			assert.Equal(t, test.data.IPAddress, res.IPAddress)
			assert.Equal(t, test.data.IPMask, res.IPMask)
			assert.Equal(t, test.data.Type, res.Type)
		}
	}
}

// 测试根据访问者id获取网段信息
func TestGetNetworksByAccessorID(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)
	// 添加测试用网段
	// Name:    "existedName1",
	// StartIP: "1.1.1.1",
	// EndIP:   "2.1.1.1",
	// Type:    "ip_segment",
	// AND
	// Name:      "existedName2",
	// IPAddress: "1.1.1.1",
	// IPMask:    "0.0.0.0",
	// Type:      "ip_mask",
	// AND
	// Name:      "existedName3",
	// IPAddress: "5.1.1.1",
	// IPMask:    "0.0.0.0",
	// Type:      "ip_mask",
	// 添加测试用访问者
	// AccessorId:   "user1",
	// AccessorType: "user",
	// AccessorId:   "department1",
	// AccessorType: "department",
	ids := addTestNetwork()
	addTestAccessor(ids[0])

	ctrl, client, driven := mockMgnt(t)
	defer ctrl.Finish()
	mgnt := NewManagementWithClient(client, driven)

	var tests = []struct {
		id         string
		s          int
		l          int
		matchIDs   []string
		matchCount int
		wanterr    error
	}{
		{"user1", 0, 20, ids[:1], 1, nil},
		{"user1", 0, 0, []string{}, 1, nil},
		{"user1", 0, -1, ids[:1], 1, nil},
		{"not_existed_acc", 0, -1, []string{}, 0, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"accessor_id"}}})},
		{"user_without_net", 0, -1, []string{}, 0, nil},
		{"department_without_net", 0, -1, []string{}, 0, nil},
		{"org_without_net", 0, -1, []string{}, 0, nil},
	}

	for _, test := range tests {
		switch test.id {
		// 用户存在，没有绑定网段
		case "user_without_net":
			mockUserName := "mock_name"
			user := sharemgnt.NcTUsrmUserInfo{
				DisplayName:   &mockUserName,
				DepartmentIds: []string{"depart1"},
				PwdControl:    false,
			}
			mockRes := sharemgnt.NcTUsrmGetUserInfo{
				ID:   test.id,
				User: &user,
			}
			client.EXPECT().GetUserByID(test.id).Return(&mockRes, nil)

		// 部门存在，没有绑定网段
		case "department_without_net":
			mockUserName := "mock_name"
			user := sharemgnt.NcTUsrmUserInfo{
				DisplayName:   &mockUserName,
				DepartmentIds: []string{"depart1"},
				PwdControl:    false,
			}
			mockRes := sharemgnt.NcTUsrmGetUserInfo{
				ID:   test.id,
				User: &user,
			}
			mockError := &ethriftexception.NcTException{
				ErrID: int32(sharemgnt.NcTShareMgntError_NCT_USER_NOT_EXIST),
			}
			client.EXPECT().GetUserByID(test.id).Return(&mockRes, mockError)

			mockRes2 := sharemgnt.NcTUsrmDepartmentInfo{
				DepartmentName: "mock_name",
			}
			client.EXPECT().GetDepartmentByID(test.id).Return(&mockRes2, nil)

		// 组织存在，没有绑定网段
		case "org_without_net":
			mockUserName := "mock_name"
			user := sharemgnt.NcTUsrmUserInfo{
				DisplayName:   &mockUserName,
				DepartmentIds: []string{"depart1"},
				PwdControl:    false,
			}
			mockRes := sharemgnt.NcTUsrmGetUserInfo{
				ID:   test.id,
				User: &user,
			}
			mockError := &ethriftexception.NcTException{
				ErrID: int32(sharemgnt.NcTShareMgntError_NCT_USER_NOT_EXIST),
			}
			client.EXPECT().GetUserByID(test.id).Return(&mockRes, mockError)

			mockRes2 := sharemgnt.NcTUsrmDepartmentInfo{
				DepartmentName: "mock_name",
			}
			mockError2 := &ethriftexception.NcTException{
				ErrID: int32(sharemgnt.NcTShareMgntError_NCT_DEPARTMENT_NOT_EXIST),
			}
			client.EXPECT().GetDepartmentByID(test.id).Return(&mockRes2, mockError2)

			mockRes1 := sharemgnt.NcTUsrmOrganizationInfo{
				OrganizationName: "mock_name",
			}
			client.EXPECT().GetOrganizationByID(test.id).Return(&mockRes1, nil)

		// 访问者不存在
		case "not_existed_acc":
			mockUserName := "mock_name"
			user := sharemgnt.NcTUsrmUserInfo{
				DisplayName:   &mockUserName,
				DepartmentIds: []string{"depart1"},
				PwdControl:    false,
			}
			mockRes := sharemgnt.NcTUsrmGetUserInfo{
				ID:   test.id,
				User: &user,
			}
			mockError := &ethriftexception.NcTException{
				ErrID: int32(sharemgnt.NcTShareMgntError_NCT_USER_NOT_EXIST),
			}
			client.EXPECT().GetUserByID(test.id).Return(&mockRes, mockError)

			mockRes2 := sharemgnt.NcTUsrmDepartmentInfo{
				DepartmentName: "mock_name",
			}
			mockError2 := &ethriftexception.NcTException{
				ErrID: int32(sharemgnt.NcTShareMgntError_NCT_DEPARTMENT_NOT_EXIST),
			}
			client.EXPECT().GetDepartmentByID(test.id).Return(&mockRes2, mockError2)

			mockRes1 := sharemgnt.NcTUsrmOrganizationInfo{
				OrganizationName: "mock_name",
			}
			mockError1 := &ethriftexception.NcTException{
				ErrID: int32(sharemgnt.NcTShareMgntError_NCT_ORGNIZATION_NOT_EXIST),
			}
			client.EXPECT().GetOrganizationByID(test.id).Return(&mockRes1, mockError1)
		default: // 默认为用户存在
			// case "user1": // 默认为用户存在
			mockUserName := "mock_name"
			user := sharemgnt.NcTUsrmUserInfo{
				DisplayName:   &mockUserName,
				DepartmentIds: []string{"depart1"},
				PwdControl:    false,
			}
			mockRes := sharemgnt.NcTUsrmGetUserInfo{
				ID:   test.id,
				User: &user,
			}
			client.EXPECT().GetUserByID(test.id).Return(&mockRes, nil)
		}
		res, c, err := mgnt.GetNetworksByAccessorID(test.id, test.s, test.l)

		if err != nil {
			assert.Equal(t, test.wanterr, err)
		}
		assert.Equal(t, test.matchCount, c)
		for _, i := range res {
			assert.Contains(t, test.matchIDs, i.ID)
		}
	}
}

// 测试EditNetwork
func TestEditNetwork(t *testing.T) {
	m := &mgnt{}
	teardown := test.SetUpDB(t)
	defer teardown(t)
	// 添加测试用网段
	ids := addTestNetwork()

	var tests = []struct {
		inputid   string
		inputdata *models.NetworkRestriction
		wantdata  *models.NetworkRestriction
		wanterr   error
	}{
		// 外网无法编辑
		{PublicNetDBIDs[0], &models.NetworkRestriction{}, &models.NetworkRestriction{}, errors.ErrNoPermission(&api.ErrorInfo{Cause: "Public net can only be read."})},
		{PublicNetID, &models.NetworkRestriction{}, &models.NetworkRestriction{}, errors.ErrNoPermission(&api.ErrorInfo{Cause: "Public net can only be read."})},
		{"not_existed_id", &models.NetworkRestriction{}, &models.NetworkRestriction{}, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"net_type", "ip_type"}}})},
		{"not_existed_id", &models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_segment", IpType: "ipv4"},
			&models.NetworkRestriction{}, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})},
		{ids[0], &models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_segment", IpType: "ipv4"},
			&models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "", IPMask: "", Type: "ip_segment"}, nil},
		{ids[0], &models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_mask", IpType: "ipv4"},
			&models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_mask", IpType: "ipv4"}, nil},
		{ids[0], &models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_mask", IpType: "ipv4"},
			&models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_mask", IpType: "ipv4"},
			errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"ip_address", "netmask"}}})},
		{ids[0], &models.NetworkRestriction{Name: "existedName2", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_mask", IpType: "ipv4"},
			&models.NetworkRestriction{Name: "ok", StartIP: "10.1.1.2", EndIP: "20.1.1.1", IPAddress: "0.0.0.0", IPMask: "255.254.0.0", Type: "ip_mask", IpType: "ipv4"},
			errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"name"}}})},
	}

	db, _ := api.ConnectDB()
	for _, test := range tests {
		err := m.EditNetwork(test.inputid, test.inputdata)
		if err != nil {
			assert.Equal(t, test.wanterr, err)
		}
		if test.wantdata != nil {
			var dbNet models.NetworkRestriction
			db.Where("f_id = ?", test.inputid).First(&dbNet)
			assert.Equal(t, (*test.wantdata).Name, dbNet.Name)
			assert.Equal(t, (*test.wantdata).StartIP, dbNet.StartIP)
			assert.Equal(t, (*test.wantdata).EndIP, dbNet.EndIP)
			assert.Equal(t, (*test.wantdata).IPAddress, dbNet.IPAddress)
			assert.Equal(t, (*test.wantdata).IPMask, dbNet.IPMask)
			assert.Equal(t, (*test.wantdata).Type, dbNet.Type)
		}
	}
}

// 测试DeleteNetwork
func TestDeleteNetwork(t *testing.T) {
	m := &mgnt{}
	teardown := test.SetUpDB(t)
	defer teardown(t)
	// 添加测试用网段
	ids := addTestNetwork()

	err := m.checkNetworkExist(ids[0])
	assert.Equal(t, nil, err)
	err = m.checkNetworkExist(ids[1])
	assert.Equal(t, nil, err)

	var tests = []struct {
		input   string
		wanterr error
	}{
		// 外网无法删除
		{PublicNetDBIDs[0], errors.ErrNoPermission(&api.ErrorInfo{Cause: "Public net can only be read."})},
		{PublicNetID, errors.ErrNoPermission(&api.ErrorInfo{Cause: "Public net can only be read."})},
		{"not_existed", errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})},
		{ids[0], nil},
		{ids[1], nil},
	}

	for _, test := range tests {
		err = m.DeleteNetwork(test.input)
		assert.Equal(t, test.wanterr, err)
		if err == nil {
			// 检查网段是否成功删除
			err = m.checkNetworkExist(test.input)
			assert.Equal(t, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}}), err)
			// 检查该网段下的访问者是否成功删除
			res, _, _ := m.getAccessors(test.input, 0, -1)
			assert.Equal(t, len(res), 0)
		}
	}
}

// 测试SearchNetwork
func TestSearchNetwork(t *testing.T) {
	m := &mgnt{}
	teardown := test.SetUpDB(t)
	defer teardown(t)
	// 添加测试用网段
	// Name:    "existedName1",
	// StartIP: "1.1.1.1",
	// EndIP:   "2.1.1.1",
	// Type:    "ip_segment",
	// AND
	// Name:      "existedName2",
	// IPAddress: "1.1.1.1",
	// IPMask:    "0.0.0.0",
	// Type:      "ip_mask",
	// AND
	// Name:      "existedName3",
	// IPAddress: "5.1.1.1",
	// IPMask:    "0.0.0.0",
	// Type:      "ip_mask",
	ids := addTestNetwork()

	var tests = []struct {
		k          string
		s          int
		l          int
		matchIDs   []string
		matchCount int
	}{
		// key_word="",start=0时,返回外网网段
		{"", 0, 20, append([]string{PublicNetID}, ids...), 5},
		{"", 0, 1, append([]string{PublicNetID}, ids[2:]...), 5},
		{"", 0, 0, []string{}, 5},
		{"existedName3", 0, 20, ids[2:], 1},
		{"3", 0, 20, ids[2:], 1},
		{"3", 0, 0, []string{}, 1},
		{"3", 1, 1, []string{}, 1},
		{"0.0", 0, -1, []string{}, 0},
		{"1.1.1.1", 0, -1, ids[:2], 2},
		{"", 1, -1, ids, 5},
		{"", -1, -1, append([]string{PublicNetID}, ids...), 5},
	}

	for _, test := range tests {
		res, c, err := m.SearchNetwork(test.k, test.s, test.l)
		assert.Equal(t, nil, err)
		assert.Equal(t, test.matchCount, c)
		if test.k == "" && test.s == 0 && test.l != 0 {
			publicNet := models.NetworkRestriction{
				ID: PublicNetID,
			}
			wantFirstNet := models.WebNetworkRestriction{NetworkRestriction: &publicNet}
			if len(res) > 0 {
				assert.Equal(t, wantFirstNet, res[0])
			}
		}
		for _, i := range res {
			assert.Contains(t, test.matchIDs, i.ID)
		}
	}
}

// 测试SearchAccessors
func TestSearchAccessors(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	ctrl, client, driven := mockMgnt(t)
	defer ctrl.Finish()
	mgnt := NewManagementWithClient(client, driven)

	// 添加测试数据
	ids := addTestNetwork()
	addTestAccessor(ids[0])

	var tests = []struct {
		netID     string
		k         string
		s         int
		l         int
		wantdata  []models.AccessorInfo
		wantcount int
		wanterr   error
	}{
		// 网段不存在
		{"not_exist", "", 0, 1, nil, 0, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})},
		{ids[1], "", 0, 1, []models.AccessorInfo{}, 0, nil},
		{ids[0], "", 0, -1, []models.AccessorInfo{{AccessorName: "mock_name", AccessorId: "department1", AccessorType: "department"},
			{AccessorName: "mock_name", AccessorId: "user1", AccessorType: "user"}}, 2, nil},
		{ids[0], "", 0, 4, []models.AccessorInfo{{AccessorName: "mock_name", AccessorId: "department1", AccessorType: "department"},
			{AccessorName: "mock_name", AccessorId: "user1", AccessorType: "user"}}, 2, nil},
	}

	for _, test := range tests {
		if test.wantcount != 0 {
			for _, i := range test.wantdata {
				switch i.AccessorType {
				case "user":
					mockUserName := "mock_name"
					user := sharemgnt.NcTUsrmUserInfo{
						DisplayName:   &mockUserName,
						DepartmentIds: []string{"depart1"},
						PwdControl:    false,
					}
					mockRes := sharemgnt.NcTUsrmGetUserInfo{
						ID:   i.AccessorId,
						User: &user,
					}
					client.EXPECT().GetUserByID(i.AccessorId).AnyTimes().Return(&mockRes, nil)
				case "department":
					mockRes := sharemgnt.NcTUsrmDepartmentInfo{
						DepartmentName: "mock_name",
					}
					client.EXPECT().GetDepartmentByID(i.AccessorId).Return(&mockRes, nil)

					mockRes1 := sharemgnt.NcTUsrmOrganizationInfo{
						OrganizationName: "mock_name",
					}
					client.EXPECT().GetOrganizationByID(i.AccessorId).AnyTimes().Return(&mockRes1, nil)
				}
			}
		}
		r, c, err := mgnt.SearchAccessors(test.netID, test.k, test.s, test.l)
		if err != nil {
			assert.Equal(t, test.wanterr, err)
		} else {
			assert.Equal(t, test.wantdata, r)
			assert.Equal(t, test.wantcount, c)
		}
	}
}

// 测试AddAccessors
func TestAddAccessors(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	ctrl, client, driven := mockMgnt(t)
	defer ctrl.Finish()
	mgnt := NewManagementWithClient(client, driven)

	// 添加测试用网段
	ids := addTestNetwork()

	var tests = []struct {
		inputid   string
		inputdata []models.AccessorInfo
		wanterr   error
	}{
		// 为外网网段具体id
		{PublicNetDBIDs[0], []models.AccessorInfo{{AccessorId: "userid_01", AccessorName: "username_01", AccessorType: "user"}},
			errors.ErrNoPermission(&api.ErrorInfo{Cause: "Not allowed to add accessors in this network."})},
		// 网段不存在
		{"not_exist", []models.AccessorInfo{{AccessorId: "userid_01", AccessorName: "username_01", AccessorType: "user"}},
			errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})},
		// 访问者类型错误
		{ids[0], []models.AccessorInfo{{AccessorId: "id_01", AccessorName: "mock_name", AccessorType: "userxxxx"}},
			errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"accessor_type"}}})},
		// 添加访问者主程序中使用的原生mysql语句，使用now()获取当前时间；UT使用的sqlite3,使用datetime('now')获取当前时间
		// 二者sql语句不兼容，注释掉成功添加情况
		// // 成功添加用户
		// {ids[0], []models.AccessorInfo{{AccessorId: "id_01", AccessorName: "mock_name", AccessorType: "user"}}, nil},
		// // 重复添加
		// {ids[0], []models.AccessorInfo{{AccessorId: "id_01", AccessorName: "mock_name", AccessorType: "user"}}, nil},
		// // 成功添加部门
		// {ids[0], []models.AccessorInfo{{AccessorId: "id_02", AccessorName: "mock_name", AccessorType: "department"}}, nil},
	}
	db, _ := api.ConnectDB()
	for _, test := range tests {
		if test.wanterr == nil {
			switch test.inputdata[0].AccessorType {
			case "user":
				mockUserName := "mock_name"
				user := sharemgnt.NcTUsrmUserInfo{
					DisplayName:   &mockUserName,
					DepartmentIds: []string{"depart1"},
					PwdControl:    false,
				}
				mockRes := sharemgnt.NcTUsrmGetUserInfo{
					ID:   test.inputdata[0].AccessorId,
					User: &user,
				}
				client.EXPECT().GetUserByID(test.inputdata[0].AccessorId).Return(&mockRes, nil)
			case "department":
				mockRes := sharemgnt.NcTUsrmDepartmentInfo{
					DepartmentName: "mock_name",
				}
				client.EXPECT().GetDepartmentByID(test.inputdata[0].AccessorId).Return(&mockRes, nil)
			}
		}
		res := mgnt.AddAccessors(test.inputid, test.inputdata)
		switch test.inputid {
		case "not_exist", PublicNetDBIDs[0]:
			// 网段不存在或为外网具体id，body中的id为网段id
			assert.Equal(t, test.inputid, res[0].ID)
		default:
			// 其余情况，body中id为访问者id
			assert.Equal(t, test.inputdata[0].AccessorId, res[0].ID)
		}
		if test.wanterr != nil {
			assert.Equal(t, test.wanterr, res[0].Body)
		} else {
			var count int64
			db.Model(&models.NetworkAccessorRelation{}).Where("f_network_id =? AND f_accessor_id = ? AND f_accessor_type = ?",
				test.inputid, test.inputdata[0].AccessorId, test.inputdata[0].AccessorType).Count(&count)
			assert.Equal(t, 1, count)
		}
	}
}

// 测试DeleteAccessors
func TestDeleteAccessors(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	// 添加测试数据
	ids := addTestNetwork()
	addTestAccessor(ids[0])

	var tests = []struct {
		netID   string
		accID   []string
		wanterr error
	}{
		{"unexist_net_id", []string{""}, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})},
		{ids[0], []string{"unexist_acc"}, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"accessor_id"}}})},
		{ids[0], []string{"user1"}, nil},
	}

	m := &mgnt{}
	db, _ := api.ConnectDB()
	for _, test := range tests {
		res := m.DeleteAccessors(test.netID, test.accID)
		// 网段不存在，body中的id为网段id
		if test.netID == "unexist_net_id" {
			assert.Equal(t, test.netID, res[0].ID)
		} else {
			// 其余情况，body中id为访问者id
			assert.Equal(t, test.accID[0], res[0].ID)
		}

		if test.wanterr != nil {
			assert.Equal(t, test.wanterr, res[0].Body)
		} else {
			var count int64
			db.Model(&models.NetworkAccessorRelation{}).Where("f_network_id =? AND f_accessor_id = ?",
				test.netID, test.accID).Count(&count)
			assert.Equal(t, int64(0), count)
		}
	}
}

// 测试FlushAccessors
func TestFlushAccessors(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	// 添加测试数据
	ids := addTestNetwork()
	addTestAccessor(ids[0])

	var tests = []struct {
		input []string
	}{
		{[]string{"unexist_accessor"}},
		{[]string{"user1"}},
	}

	m := &mgnt{}
	db, _ := api.ConnectDB()
	for _, test := range tests {
		err := m.FlushAccessors(test.input)
		assert.Equal(t, nil, err)
		var count int64
		db.Model(&models.NetworkAccessorRelation{}).Where("f_accessor_id = ?",
			test.input).Count(&count)
		assert.Equal(t, int64(0), count)
	}
}

// 测试GetNetworkData
func TestGetNetworkData(t *testing.T) {
	teardown := test.SetUpDB(t)
	defer teardown(t)

	ctrl, client, driven := mockMgnt(t)
	defer ctrl.Finish()
	mgnt := NewManagementWithClient(client, driven)

	// test department tree
	// -department1
	//         - department2
	//                  - user1
	//         - user2

	driven.EXPECT().GetDepartUserIds(gomock.Any()).Return([]string{"user1"}, nil)
	//var mockRes3 []*sharemgnt.NcTDepartmentInfo

	// 添加测试数据
	ids := addTestNetwork()
	addTestAccessor(ids[0])

	//allpolicy --> convertdata
	res := make(map[string]interface{})
	err := mgnt.GetNetworkData(res)
	assert.Equal(t, nil, err)

	// 实际结果如下，由于返回的是interface,且其中有指针，不能整体比较
	// 	wantRes := `{
	//     "departments": {
	//         "department1": [
	//             {
	//                 "start_ip": 16843009,
	//                 "end_ip": 33620225
	//             }
	//         ]
	//     },
	//     "users": {
	//         "user2": {
	//             "departments": [
	//                 "1"
	//             ],
	//         "user1": {
	//             "nets": [
	//                 {
	//                     "start_ip": 16843009,
	//                     "end_ip": 33620225
	//                 }
	//             ]
	//         }
	//     }
	// }`
	// wantDepartments := map[string][]models.NetSegment{
	// 	"department1": []models.NetSegment{models.NetSegment{StartIP: big.NewInt(16843009), EndIP: big.NewInt(33620225)}},
	// }
	// wantUsers := map[string]interface{}{
	// 	"user1": &models.UserNetworkInfo{Departments: []string(nil),NetSegements: []models.NetSegment{models.NetSegment{StartIP: big.NewInt(16843009), EndIP: big.NewInt(33620225)}}},
	// 	"user2": &models.UserNetworkInfo{Departments: []string{"department1"}, NetSegements: []models.NetSegment(nil)},
	// }
	// resultIpv4 := map[string]interface{}{
	// 	"departments": wantDepartments,
	// 	"users": wantUsers,
	// }
	// assert.Equal(t, resultIpv4, res["ipv4"])

	// users, ok := res["users"].(map[string]*models.UserNetworkInfo)
	// assert.Equal(t, true, ok)

	// wantUser1 := models.UserNetworkInfo{Departments: []string(nil),
	// 	NetSegements: []models.NetSegment{models.NetSegment{StartIP: 16843009, EndIP: 33620225}}}
	// assert.Equal(t, wantUser1, *users["user1"])
	// wantUser2 := models.UserNetworkInfo{Departments: []string{"department1"}, NetSegements: []models.NetSegment(nil)}
	// assert.Equal(t, wantUser2, *users["user2"])
}

// // 测试UpdateUserDepartmentRelation
// func TestUpdateUserDepartmentRelation(t *testing.T) {
// 	teardown := test.SetUpDB(t)
// 	defer teardown(t)

// 	ctrl, client, driven := mockMgnt(t)
// 	defer ctrl.Finish()
// 	mgnt := NewManagementWithClient(client, driven)

// 	driven.EXPECT().GetBelongDepartByUserId(gomock.Any()).Return([]string{"department1", "department2"}, nil).Times(2)
// 	driven.EXPECT().IncrementalUpdateData(gomock.Any(), gomock.Any()).Return(nil).Times(4)

// 	// 添加测试数据
// 	ids := addTestNetwork()
// 	addTestAccessor(ids[0])
// 	addTestAccessor(ids[3])

// 	// case1
// 	err := mgnt.UpdateUserDepartmentRelation("user2", []string{"department1/department2"})
// 	assert.Equal(t, nil, err)

// 	// case2
// 	err = mgnt.UpdateUserDepartmentRelation("user1", []string{"department1/department2"})
// 	assert.Equal(t, nil, err)

// 	// case3
// 	driven.EXPECT().GetBelongDepartByUserId(gomock.Any()).Return([]string{"department1", "department2"}, fmt.Errorf("test 500"))
// 	err = mgnt.UpdateUserDepartmentRelation("user1", []string{"department1/department2"})
// 	assert.Equal(t, err.Error(), "GetBelongDepartByUserId: test 500")

// 	// case4
// 	driven.EXPECT().GetBelongDepartByUserId(gomock.Any()).Return([]string{"department1", "department2"}, nil)
// 	driven.EXPECT().IncrementalUpdateData(gomock.Any(), gomock.Any()).Return(fmt.Errorf("test 500"))
// 	err = mgnt.UpdateUserDepartmentRelation("user1", []string{"department1/department2"})
// 	assert.Equal(t, err.Error(), "IncrementalUpdateData: test 500")

// }
