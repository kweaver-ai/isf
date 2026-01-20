package network

import (
	"bytes"
	"database/sql"
	cerrors "errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"policy_mgnt/dependency"
	"policy_mgnt/tapi/ethriftexception"
	"policy_mgnt/tapi/sharemgnt"
	"policy_mgnt/thrift"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	cnetip "policy_mgnt/utils/gocommon/v2/net/netip"

	"github.com/kweaver-ai/GoUtils/utilities"
	"github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Management
type Management interface {
	AddNetwork(*models.NetworkRestriction) (string, error)
	GetNetworkByID(string) (models.WebNetworkRestriction, error)
	EditNetwork(string, *models.NetworkRestriction) error
	DeleteNetwork(string) error
	SearchNetwork(string, int, int) ([]models.WebNetworkRestriction, int, error)
	SearchAccessors(string, string, int, int) ([]models.AccessorInfo, int, error)
	GetNetworksByAccessorID(string, int, int) ([]models.WebNetworkRestriction, int, error)
	AddAccessors(string, []models.AccessorInfo) []*api.MultiStatus
	DeleteAccessors(string, []string) []*api.MultiStatus
	GetNetworkData(map[string]interface{}) error
	FlushAccessors([]string) error
}

type mgnt struct {
	sharemgnt      thrift.ShareMgnt
	abstractDriven dependency.AbstractDriven
}

// NewManagement 初始化管理实例
func NewManagement() (Management, error) {
	driven := dependency.NewAbstractDriven()
	client, err := thrift.NewShareMgnt()
	if err != nil {
		return nil, err
	}
	return NewManagementWithClient(client, driven), nil
}

// NewManagementWithClient 初始化管理实例
func NewManagementWithClient(client thrift.ShareMgnt, driven dependency.AbstractDriven) Management {
	return &mgnt{sharemgnt: client, abstractDriven: driven}
}

// 检查网段名称
func (m *mgnt) checkNetworkName(name string) bool {
	// 可以为空字符串
	if len(name) == 0 {
		return true
	}
	re := regexp.MustCompile(`^[^\\/:*?"<>|]{1,128}$`)
	return re.MatchString(name)
}

// 检查网段是否重复
func (m *mgnt) checkDuplicateNet(params *models.NetworkRestriction, id ...string) (err error) {
	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	var count int64
	query := db.Model(&models.NetworkRestriction{})
	// 如果传入id，则排除该id的网段
	if len(id) != 0 {
		query = query.Not("f_id = ?", id[0])
	}

	if params.Type == "ip_segment" {
		err = query.Where("f_start_ip = ? AND f_end_ip = ? AND f_type = ?", params.StartIP, params.EndIP, params.Type).Count(&count).Error
		if count != 0 {
			err = errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"start_ip", "end_ip"}}})
			return err
		}
	} else if params.Type == "ip_mask" {
		err = query.Where("f_ip_address = ? AND f_ip_mask = ? AND f_type = ?", params.IPAddress, params.IPMask, params.Type).Count(&count).Error
		if count != 0 {
			err = errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"ip_address", "netmask"}}})
			return err
		}
	}
	return
}

// 检查传入参数
func (m *mgnt) checkNetworkParams(params *models.NetworkRestriction) error {
	trimNet(params)
	ok := m.checkNetworkName(params.Name)
	if !ok {
		return errors.ErrInvalideName(&api.ErrorInfo{Detail: map[string]interface{}{"names": []string{params.Name}}})
	}

	if params.Type == "ip_segment" {
		// type=ip_segment, 采用start_ip+end_ip
		if params.IpType == "ipv4" {
			startIP := net.ParseIP(params.StartIP)
			if startIP.To4() == nil {
				return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip"}}})
			}

			endIP := net.ParseIP(params.EndIP)
			if endIP.To4() == nil {
				return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})
			}
			if bytes.Compare(startIP, endIP) == 1 {
				return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip", "end_ip"}}, Cause: "Start_ip can not greater than end_ip."})
			}
		}
		if params.IpType == "ipv6" {
			startIP := net.ParseIP(params.StartIP)
			if startIP.To4() != nil || startIP.To16() == nil {
				return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip"}}})
			}

			endIP := net.ParseIP(params.EndIP)
			if endIP.To4() != nil || endIP.To16() == nil {
				return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"end_ip"}}})
			}
			if bytes.Compare(startIP, endIP) == 1 {
				return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"start_ip", "end_ip"}}, Cause: "Start_ip can not greater than end_ip."})
			}
		}
	} else if params.Type == "ip_mask" && params.IpType == "ipv4" {
		// type=ip_mask, 采用ip+mask
		ipAdress := net.ParseIP(params.IPAddress)
		if ipAdress.To4() == nil {
			return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address"}}})
		}

		subnetMap := utilities.SubnetMap()
		_, ok := subnetMap[params.IPMask]
		if !ok {
			return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"netmask"}}})
		}
	} else if params.Type == "ip_mask" && params.IpType == "ipv6" {
		_, _, err := net.ParseCIDR(params.IPAddress + "/" + params.IPMask)
		if err != nil {
			return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address", "netmask"}}})
		}
	} else {
		return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"net_type", "ip_type"}}})
	}
	if params.IpType != "ipv4" && params.IpType != "ipv6" {
		return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_type"}}})
	}
	return nil
}

func (m *mgnt) getNetworkByName(name string) (*models.NetworkRestriction, error) {
	db, err := api.ConnectDB()
	if err != nil {
		return nil, err
	}
	var res models.NetworkRestriction
	err = db.Where("f_name = ?", name).First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// 搜索网段
func (m *mgnt) SearchNetwork(keyWord string, start, limit int) (data []models.WebNetworkRestriction, count int, err error) {
	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	var count64 int64
	var res []models.NetworkRestriction
	query := db.Order("f_created_at desc, f_segment_start, f_segment_end").Where("f_id not in (?)", PublicNetDBIDs)
	switch {
	case len(keyWord) > 0:
		keyWord = fmt.Sprintf("%%%s%%", keyWord)
		query = query.Where("f_type = ? AND (f_name like ? OR f_start_ip like ? OR f_end_ip like ?)", "ip_segment",
			keyWord, keyWord, keyWord).Or("f_type = ? AND (f_name like ? OR f_ip_address like ?)", "ip_mask", keyWord, keyWord)
		err = query.Model(&models.NetworkRestriction{}).Count(&count64).Error
		if err != nil {
			return
		}
		count = int(count64)

		if limit != -1 {
			query = query.Offset(start).Limit(limit)
		} else {
			query = query.Offset(start).Limit(count)
		}
		err = query.Find(&res).Error
		if err != nil {
			return
		}
	default:
		// 如果key_word为空，需要显示外网网段, count+1
		err = db.Model(&models.NetworkRestriction{}).Where("f_id not in (?)", PublicNetDBIDs).Count(&count64).Error
		if err != nil {
			return
		}
		count = int(count64)

		if start > 0 {
			// start > 0时，start-1,limit不变
			if limit != -1 {
				query = query.Offset(start - 1).Limit(limit)
			} else {
				query = query.Offset(start - 1).Limit(count)
			}

			err = query.Find(&res).Error
			if err != nil {
				return
			}
		} else {
			// start <= 0时，start不变，limit-1
			if limit > 0 {
				query = query.Offset(start).Limit(limit - 1)
			} else if limit == 0 {
				query = query.Offset(start).Limit(0)
			} else if limit == -1 {
				query = query.Offset(start).Limit(count)
			} else {
				query = query.Offset(start).Limit(0)
			}

			err = query.Find(&res).Error
			if err != nil {
				return
			}

			// 首部增加外网网段
			if limit > 0 || limit == -1 {
				publicNet := models.NetworkRestriction{
					ID: PublicNetID,
				}
				res = append([]models.NetworkRestriction{publicNet}, res...)
			}
		}
		count += 1
	}

	data = make([]models.WebNetworkRestriction, 0)
	for i := range res {
		webNet := models.WebNetworkRestriction{NetworkRestriction: &res[i]}
		data = append(data, webNet)
	}

	return
}

func (m *mgnt) getNetworksByAccessorID(accessorID string, start, limit int) (res []models.NetworkRestriction, count int, err error) {
	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	query := db.Table("t_network_restriction as net").
		Select("net.f_id, net.f_name, net.f_start_ip, net.f_end_ip, net.f_ip_address, net.f_ip_mask, net.f_type, net.f_ip_type").
		Joins("inner join t_network_accessor_relation as relation on relation.f_network_id = net.f_id").
		Order("net.f_created_at desc, net.f_segment_start, net.f_segment_end")

	// 搜索含有公网id的网段
	query2 := query.Where("relation.f_accessor_id = ? AND relation.f_network_id in (?) ", accessorID, PublicNetDBIDs)
	var count64, count2 int64
	err = query2.Count(&count2).Error
	if err != nil {
		return
	}

	query = db.Table("t_network_restriction as net").
		Select("net.f_id, net.f_name, net.f_start_ip, net.f_end_ip, net.f_ip_address, net.f_ip_mask, net.f_type, net.f_ip_type").
		Joins("inner join t_network_accessor_relation as relation on relation.f_network_id = net.f_id").
		Order("net.f_created_at desc, net.f_segment_start, net.f_segment_end")

	// 搜索不含有公网id的网段
	query1 := query.Where("relation.f_accessor_id = ? AND relation.f_network_id not in (?) ", accessorID, PublicNetDBIDs)
	err = query1.Count(&count64).Error
	if err != nil {
		return
	}
	count = int(count64)

	if count2 > 0 {
		// 如果此访问者绑定了外网网段
		if start > 0 {
			// start > 0时，start-1,limit不变
			if limit != -1 {
				query1 = query1.Offset(start - 1).Limit(limit)
			} else {
				query1 = query1.Offset(start - 1).Limit(count)
			}
		} else {
			// start <= 0时，start不变，limit-1
			if limit > 0 {
				query1 = query1.Offset(start).Limit(limit - 1)
			} else if limit == 0 {
				query1 = query1.Offset(start).Limit(0)
			} else if limit == -1 {
				query1 = query1.Offset(start).Limit(count)
			} else {
				query1 = query1.Offset(start).Limit(0)
			}
		}
		count += 1

	} else {
		if limit != -1 {
			query1 = query1.Offset(start).Limit(limit)
		} else {
			query1 = query1.Offset(start).Limit(count)
		}
	}

	rows, err := query1.Rows()
	for rows.Next() {
		var net models.NetworkRestriction
		// 防止name为null时转化数据报错
		var name sql.NullString
		err = rows.Scan(&net.ID, &name, &net.StartIP, &net.EndIP, &net.IPAddress, &net.IPMask, &net.Type, &net.IpType)
		if err != nil {
			return
		}
		if name.Valid {
			net.Name = name.String
		}
		res = append(res, net)
	}

	// 如果此访问者绑定了外网网段
	if count2 > 0 && start == 0 && (limit > 0 || limit == -1) {
		publicNet := models.NetworkRestriction{
			ID: PublicNetID,
		}
		res = append([]models.NetworkRestriction{publicNet}, res...)
	}

	return
}

// 根据访问者id搜索已绑定的网段
func (m *mgnt) GetNetworksByAccessorID(id string, start, limit int) (data []models.WebNetworkRestriction, count int, err error) {
	// 检查访问者id是否存在
	data = make([]models.WebNetworkRestriction, 0)
	if err := m.checkAccessor(&models.AccessorInfo{AccessorId: id, AccessorType: "any"}); err != nil {
		return data, 0, err
	}

	res, count, err := m.getNetworksByAccessorID(id, start, limit)
	if err != nil {
		return
	}

	for i := range res {
		webNet := models.WebNetworkRestriction{NetworkRestriction: &res[i]}
		data = append(data, webNet)
	}

	return
}

// 根据id获取网段信息
func (m *mgnt) getNetworkByID(id string) (*models.NetworkRestriction, error) {
	db, err := api.ConnectDB()
	if err != nil {
		return nil, err
	}

	var res models.NetworkRestriction
	err = db.Where("f_id = ?", id).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// 检查网段是否存在
// 不存在，报404
func (m *mgnt) checkNetworkExist(id string) (err error) {
	// 如果是外网id("public-net"), 存在
	if id == PublicNetID {
		return nil
	}

	db, err := api.ConnectDB()
	if err != nil {
		return
	}
	var count int64
	err = db.Model(&models.NetworkRestriction{}).Where("f_id = ?", id).Count(&count).Error
	if err != nil {
		return
	}

	if count == 0 {
		err = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})
		return
	}
	return
}

// 获取网段
func (m *mgnt) GetNetworkByID(id string) (result models.WebNetworkRestriction, err error) {
	// 如果是外网，直接返回
	if utilities.InStrSlice(id, append(PublicNetDBIDs, PublicNetID)) {
		publicNet := models.NetworkRestriction{
			ID: id,
		}
		return models.WebNetworkRestriction{NetworkRestriction: &publicNet}, nil
	}
	res, err := m.getNetworkByID(id)
	if err != nil {
		if cerrors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})
		}
		return
	}

	result = models.WebNetworkRestriction{NetworkRestriction: res}
	return
}

// 修改网段前，处理数据
// type为ip_segment：修改start_ip+end_ip,
// type为ip_mask：修改ip+子网掩码
func (m *mgnt) modifyBeforeEdit(dbParams, params *models.NetworkRestriction) {
	switch params.Type {
	case "ip_segment":
		dbParams.StartIP = params.StartIP
		dbParams.EndIP = params.EndIP
	case "ip_mask":
		dbParams.IPAddress = params.IPAddress
		dbParams.IPMask = params.IPMask
	}
	dbParams.Name = params.Name
	dbParams.Type = params.Type
	dbParams.IpType = params.IpType
}

// 修改网段
func (m *mgnt) EditNetwork(id string, params *models.NetworkRestriction) (err error) {
	// 外网不能编辑、删除
	if utilities.InStrSlice(id, append(PublicNetDBIDs, PublicNetID)) {
		return errors.ErrNoPermission(&api.ErrorInfo{Cause: "Public net can only be read."})
	}

	if err := m.checkNetworkParams(params); err != nil {
		return err
	}

	dbParams, err := m.getNetworkByID(id)
	if err != nil {
		if cerrors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}})
		}
		return
	}

	if len(params.Name) != 0 {
		// check duplicate name
		res, err := m.getNetworkByName(params.Name)
		if err != nil {
			if !cerrors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		if res != nil && id != res.ID {
			err = errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"name"}}})
			return err
		}
	}

	m.modifyBeforeEdit(dbParams, params)
	// 检查重复网段
	err = m.checkDuplicateNet(params, id)
	if err != nil {
		return
	}

	// 把ip转换为int, 增加到params中
	switch params.IpType {
	case "ipv4":
		err = m.fillIntIP(dbParams)
		if err != nil {
			return
		}
	case "ipv6":
		err = m.fillIntIPV6(dbParams)
		if err != nil {
			return
		}
	}

	db, err := api.ConnectDB()
	if err != nil {
		return
	}
	tx := db.Begin()
	var netName interface{}
	if dbParams.Name == "" {
		netName = nil
	} else {
		netName = dbParams.Name
	}
	// 网段名为""时，数据库中的f_name应该为null，使用原生sql
	sqlString := "UPDATE t_network_restriction SET f_name=?, f_start_ip=?, f_end_ip=?, f_ip_address=?, f_ip_mask=?,f_segment_start=?, f_segment_end=?,f_type=?, f_ip_type=? where f_id = ?"
	values := []interface{}{netName, dbParams.StartIP, dbParams.EndIP, dbParams.IPAddress, dbParams.IPMask, dbParams.SegmentStart, dbParams.SegmentEnd, dbParams.Type, dbParams.IpType, dbParams.ID}
	if err := db.Exec(sqlString, values...).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return
}

// 删除网段
func (m *mgnt) DeleteNetwork(id string) (err error) {
	// 外网不能编辑、删除
	if utilities.InStrSlice(id, append(PublicNetDBIDs, PublicNetID)) {
		return errors.ErrNoPermission(&api.ErrorInfo{Cause: "Public net can only be read."})
	}

	if err := m.checkNetworkExist(id); err != nil {
		return err
	}

	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	err = db.Where("f_id = ?", id).Delete(&models.NetworkRestriction{}).Error
	if err != nil {
		return
	}

	err = db.Where("f_network_id = ?", id).Delete(&models.NetworkAccessorRelation{}).Error
	if err != nil {
		return
	}

	return
}

// 添加网段前，清空不需要保存的数据
// ip_segment为start_ip+end_ip, ip_mask为ip+子网掩码
func (m *mgnt) modifyBeforeAdd(params *models.NetworkRestriction) {
	switch params.Type {
	case "ip_segment":
		params.IPAddress = ""
		params.IPMask = ""
	case "ip_mask":
		params.StartIP = ""
		params.EndIP = ""
	}
}

// 把ip转换为int, 增加到params中
func (m *mgnt) fillIntIP(params *models.NetworkRestriction) (err error) {
	var segmentStart, segmentEnd int64
	switch params.Type {
	case "ip_segment":
		segmentStart, err = utilities.ConvertIPAoti(params.StartIP)
		if err != nil {
			return
		}
		params.SegmentStart = strconv.FormatInt(segmentStart, 10)
		segmentEnd, err = utilities.ConvertIPAoti(params.EndIP)
		if err != nil {
			return
		}
		params.SegmentEnd = strconv.FormatInt(segmentEnd, 10)
	case "ip_mask":
		ipRange := utilities.ConvertNetToRange(params.IPAddress, params.IPMask)
		if len(ipRange) < 1 {
			if err != nil {
				return errors.ErrInternalServerErrorPublic(&api.ErrorInfo{Cause: "Convert ip+mask to ip range error."})
			}
		}

		startIP := ipRange[0]
		endIP := ipRange[1]
		// 去除网络地址和广播地址
		start, err := utilities.ConvertIPAoti(startIP)
		if err != nil {
			return err
		}
		if params.IPMask == "255.255.255.255" || params.IPMask == "255.255.255.254" {
			params.SegmentStart = strconv.FormatInt(start, 10)
		} else {
			params.SegmentStart = strconv.FormatInt(start+1, 10)
		}

		end, err := utilities.ConvertIPAoti(endIP)
		if err != nil {
			return err
		}
		if params.IPMask == "255.255.255.255" {
			params.SegmentEnd = strconv.FormatInt(start, 10)
		} else if params.IPMask == "255.255.255.254" {
			params.SegmentEnd = strconv.FormatInt(end, 10)
		} else {
			params.SegmentEnd = strconv.FormatInt(end-1, 10)
		}
	}
	return
}

// 把ipv6转换为int, 增加到params中
func (m *mgnt) fillIntIPV6(params *models.NetworkRestriction) (err error) {
	switch params.Type {
	case "ip_segment":
		segmentStart := cnetip.IpToInt(net.ParseIP(params.StartIP))
		segmentEnd := cnetip.IpToInt(net.ParseIP(params.EndIP))
		params.SegmentStart = segmentStart.String()
		params.SegmentEnd = segmentEnd.String()
	case "ip_mask":
		_, ipnet, err := net.ParseCIDR(params.IPAddress + "/" + params.IPMask)
		if err != nil {
			return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"ip_address", "netmask"}}})
		}
		startIP, endIP := cnetip.IPNetToIPRange(ipnet)
		segmentStart := cnetip.IpToInt(startIP)
		segmentEnd := cnetip.IpToInt(endIP)
		params.SegmentStart = segmentStart.String()
		params.SegmentEnd = segmentEnd.String()
	}
	return
}

// 去除网段左边的0,比如"010.0.000.090"->"10.0.0.90"
func trimLeftZero(netString string) string {
	if len(netString) < 1 {
		return ""
	}
	netSegs := strings.Split(netString, ".")
	var newNetSegs []string
	for _, i := range netSegs {
		if len(i) > 0 {
			newNetSegs = append(newNetSegs, strings.TrimLeft(i[:len(i)-1], "0")+i[len(i)-1:])
		}
	}
	return strings.Join(newNetSegs, ".")
}

// 去除网段信息左边的0
func trimNet(params *models.NetworkRestriction) {
	params.StartIP = trimLeftZero(params.StartIP)
	params.EndIP = trimLeftZero(params.EndIP)
	params.IPAddress = trimLeftZero(params.IPAddress)
	params.IPMask = trimLeftZero(params.IPMask)
}

// 添加网段，返回成功添加的网段id
func (m *mgnt) AddNetwork(params *models.NetworkRestriction) (networkID string, err error) {
	trimNet(params)
	err = m.checkNetworkParams(params)
	if err != nil {
		return
	}

	if len(params.Name) != 0 {
		res, err := m.getNetworkByName(params.Name)
		if err != nil {
			if !cerrors.Is(err, gorm.ErrRecordNotFound) {
				return "", err
			}
		}
		if res != nil {
			err = errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"name"}}})
			return "", err
		}
	}

	byteID, err := uuid.NewV4()
	if err != nil {
		return
	}
	networkID = byteID.String()

	params.ID = networkID
	m.modifyBeforeAdd(params)

	// 检查重复网段
	err = m.checkDuplicateNet(params)
	if err != nil {
		return
	}

	// 把ip转换为int, 增加到params中
	switch params.IpType {
	case "ipv4":
		err = m.fillIntIP(params)
		if err != nil {
			return
		}
	case "ipv6":
		err = m.fillIntIPV6(params)
		if err != nil {
			return
		}
	}

	db, err := api.ConnectDB()
	if err != nil {
		return
	}
	tx := db.Begin()
	err = db.Save(params).Error
	if err != nil {
		if sqlErr, ok := err.(*mysql.MySQLError); ok {
			// Error 1062: Duplicate entry
			// 键重复sql错误重新捕获，改为409抛给前端
			if sqlErr.Number == 1062 {
				err = errors.ErrConflict(&api.ErrorInfo{Detail: map[string]interface{}{"conflict_params": []string{"name"}}})
				return "", err
			}
		}
		tx.Rollback()
		return
	}
	tx.Commit()

	return
}

// 检查访问者在引擎的数据库中是否存在
// 若存在，补充访问者的名称
func (m *mgnt) checkThriftUser(accessorInfo *models.AccessorInfo) error {
	userInfo, userErr := m.sharemgnt.GetUserByID(accessorInfo.AccessorId)
	if userErr != nil {
		if ncuserErr, ok := userErr.(*ethriftexception.NcTException); ok {
			if ncuserErr.ErrID == int32(sharemgnt.NcTShareMgntError_NCT_USER_NOT_EXIST) {
				return errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"accessor_id"}}})
			}
		}
		return userErr
	}
	accessorInfo.AccessorName = *userInfo.User.DisplayName
	return nil
}

// 检查访问者在引擎的数据库中是否存在
// 若存在，补充访问者的名称
func (m *mgnt) checkThriftDepartment(accessorInfo *models.AccessorInfo) error {
	//根据id查找部门
	departmentInfo, departmentErr := m.sharemgnt.GetDepartmentByID(accessorInfo.AccessorId)
	if departmentErr != nil {
		if ncdepartmentErr, ok := departmentErr.(*ethriftexception.NcTException); ok {
			if ncdepartmentErr.ErrID == int32(sharemgnt.NcTShareMgntError_NCT_DEPARTMENT_NOT_EXIST) {
				return errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"accessor_id"}}})
			}
		}
		return departmentErr
	}
	accessorInfo.AccessorName = departmentInfo.DepartmentName
	return nil
}

// 检查访问者在引擎的数据库中是否存在
// 若存在，补充访问者的名称
func (m *mgnt) checkThriftOrg(accessorInfo *models.AccessorInfo) error {
	//根据id查找组织
	orgInfo, orgErr := m.sharemgnt.GetOrganizationByID(accessorInfo.AccessorId)
	if orgErr != nil {
		if ncorgErr, ok := orgErr.(*ethriftexception.NcTException); ok {
			if ncorgErr.ErrID == int32(sharemgnt.NcTShareMgntError_NCT_ORGNIZATION_NOT_EXIST) {
				return errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"accessor_id"}}})
			}
		}
		return orgErr
	}
	accessorInfo.AccessorName = orgInfo.OrganizationName
	return nil
}

func (m *mgnt) errorHandle(accID string, terrors ...error) error {
	for _, v := range terrors {
		// 跳过空指针
		if v == nil {
			continue
		}

		if _, ok := v.(*api.Error); !ok {
			// 如果不是APIError，抛错
			return errors.ErrInternalServerErrorPublic(&api.ErrorInfo{Cause: "Thrift server error."})
		}
	}

	// 只要全部是APIError，报404
	return errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"accessor_id"}}})
}

// 检查访问者在引擎的数据库中是否存在
// 若存在，补充访问者的名称
func (m *mgnt) checkAccessor(accessorInfo *models.AccessorInfo) (err error) {
	accessorType := accessorInfo.AccessorType
	var userErr, departmentErr, orgErr error
	if accessorType != "user" && accessorType != "department" && accessorType != "any" {
		return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"accessor_type"}}})
	}

	if accessorType == "user" || accessorType == "any" {
		userErr = m.checkThriftUser(accessorInfo)
		if userErr == nil {
			return
		}
	}

	if accessorType == "department" || accessorType == "any" {
		//根据id查找部门、组织
		departmentErr = m.checkThriftDepartment(accessorInfo)
		if departmentErr == nil {
			return
		}
		orgErr = m.checkThriftOrg(accessorInfo)
		if orgErr == nil {
			return
		}
	}

	err = m.errorHandle(accessorInfo.AccessorId, userErr, departmentErr, orgErr)
	return
}

func (m *mgnt) getRelations(networkID string, start, limit int) (relations []models.NetworkAccessorRelation, count int, err error) {
	var count64 int64
	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	// 如果网段id为public-net，转化为数据库外网id的第一个
	if networkID == PublicNetID {
		networkID = PublicNetDBIDs[0]
	}
	err = db.Model(&models.NetworkAccessorRelation{}).Where("f_network_id = ?", networkID).Count(&count64).Error
	if err != nil {
		return
	}
	count = int(count64)
	if limit == -1 {
		limit = count
	}

	err = db.Where("f_network_id = ?", networkID).Offset(start).Limit(limit).Order("f_created_at desc, f_id desc").Find(&relations).Error
	if err != nil {
		return
	}

	return
}

// 获取已经添加过的访问者ID
func (m *mgnt) getAddedAccessorIDs(networkID string) (accessorIDs []string, err error) {
	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	var result []models.NetworkAccessorRelation
	err = db.Select("f_accessor_id").Where("f_network_id = ?", networkID).Find(&result).Error
	if err != nil {
		return
	}

	for _, each := range result {
		accessorIDs = append(accessorIDs, each.AccessorId)
	}

	return
}

// 获取已经添加过的访问者
func (m *mgnt) getAccessors(networkID string, start, limit int) (accessorInfos []models.AccessorInfo, count int, err error) {
	accessorInfos = make([]models.AccessorInfo, 0)
	relations, count, err := m.getRelations(networkID, start, limit)
	if err != nil {
		return
	}

	for _, relation := range relations {
		var accessorInfo models.AccessorInfo
		accessorInfo.AccessorId = relation.AccessorId
		accessorInfo.AccessorType = relation.AccessorType
		accessorInfos = append(accessorInfos, accessorInfo)
	}
	return
}

// 根据网段id、访问者id获取已经添加过的访问者id列表
func (m *mgnt) getAccessorsByNetIDAccessorIDs(networkID string, accessorIDs []string, start, limit int) ([]string, int, error) {
	wantIDs := make([]string, 0)
	var count int64
	// 访问者id为空，返回
	if len(accessorIDs) == 0 {
		return wantIDs, 0, nil
	}

	db, err := api.ConnectDB()
	if err != nil {
		return nil, 0, err
	}

	// 如果网段id为public-net，转化为数据库外网id的第一个
	if networkID == PublicNetID {
		networkID = PublicNetDBIDs[0]
	}
	err = db.Model(&models.NetworkAccessorRelation{}).Where("f_network_id = ? AND f_accessor_id in (?)", networkID, accessorIDs).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if limit == -1 {
		limit = int(count)
	}

	var relations []models.NetworkAccessorRelation
	err = db.Select("f_accessor_id").Where("f_network_id = ? AND f_accessor_id in (?)", networkID,
		accessorIDs).Offset(start).Limit(limit).Order("f_created_at desc").Find(&relations).Error
	if err != nil {
		return nil, 0, err
	}

	for _, relation := range relations {
		wantIDs = append(wantIDs, relation.AccessorId)
	}
	return wantIDs, int(count), nil
}

// 过滤重复添加的访问者
func (m *mgnt) trimDupAccessors(accessorInfos []models.AccessorInfo) (accessors []models.AccessorInfo) {
	var input []interface{}
	for _, v := range accessorInfos {
		input = append(input, v)
	}

	output := utilities.RemoveDuplicateStruct(input)

	for _, v := range output {
		t := v.(models.AccessorInfo)
		accessors = append(accessors, t)
	}
	return
}

// 忽略已经添加过的访问者
func (m *mgnt) filterAccessors(networkID string, accessors []models.AccessorInfo) (dbAccessors []models.AccessorInfo, err error) {
	existedAccessors, _, err := m.getAccessors(networkID, 0, -1)
	if err != nil {
		return
	}
	// 如果没有添加过用户，全部返回
	if len(existedAccessors) == 0 {
		dbAccessors = accessors
		return
	}

	for _, accessorInfo := range accessors {
		existed := false
		for _, existedAccessor := range existedAccessors {
			// 如果数据库关系表中存在相同的访问者，过滤
			if accessorInfo.AccessorId == existedAccessor.AccessorId {
				existed = true
				continue
			}
		}
		if !existed {
			dbAccessors = append(dbAccessors, accessorInfo)
		}
	}

	return
}

// 添加访问者至指定网段
func (m *mgnt) addAccessorsToNet(networkID string, accessorInfos []models.AccessorInfo) error {
	// 批量添加
	n := len(accessorInfos)
	if n == 0 {
		return nil
	}
	db, err := api.ConnectDB()
	if err != nil {
		return err
	}

	// 存放 (?, ?, ?) 的slice
	valueStrings := make([]string, 0, n)
	// 存放values的slice
	valueArgs := make([]interface{}, 0, n*3)
	// 遍历nets准备相关数据
	for _, v := range accessorInfos {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?, ?, now())")
		valueArgs = append(valueArgs, networkID, v.AccessorId, v.AccessorType)
	}
	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO t_network_accessor_relation (f_network_id, f_accessor_id, f_accessor_type, f_created_at) VALUES %s",
		strings.Join(valueStrings, ","))

	tx := db.Begin()
	if err := db.Exec(stmt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 添加访问者至数据库
func (m *mgnt) addAccessors(networkID string, accessorInfos []models.AccessorInfo) error {
	if networkID == PublicNetID {
		for _, v := range PublicNetDBIDs {
			if err := m.addAccessorsToNet(v, accessorInfos); err != nil {
				return err
			}
		}
	} else {
		if err := m.addAccessorsToNet(networkID, accessorInfos); err != nil {
			return err
		}
	}

	return nil
}

// search user and department from sharemgnt
func (m *mgnt) searchAccessorFromSharemgnt(keyWord string, start, limit int) ([]models.AccessorInfo, error) {
	sharemgntAccessors := make([]models.AccessorInfo, 0)
	// 根据名称从引擎获取所有匹配
	sharemgntUsers, err := m.sharemgnt.SearchSupervisoryUsers(sharemgnt.NCT_USER_ADMIN, keyWord, start, limit)
	if err != nil {
		return nil, err
	}
	for _, v := range sharemgntUsers {
		accessor := models.AccessorInfo{
			AccessorId:   v.ID,
			AccessorName: v.DisplayName,
			AccessorType: "user",
		}
		sharemgntAccessors = append(sharemgntAccessors, accessor)
	}

	sharemgntDepartments, err := m.sharemgnt.SearchDepartments(sharemgnt.NCT_USER_ADMIN, keyWord, start, limit)
	if err != nil {
		return nil, err
	}
	for _, v := range sharemgntDepartments {
		accessor := models.AccessorInfo{
			AccessorId:   v.DepartmentId,
			AccessorName: v.DepartmentName,
			AccessorType: "department",
		}
		sharemgntAccessors = append(sharemgntAccessors, accessor)
	}
	return sharemgntAccessors, nil
}

// 获取访问者信息
func (m *mgnt) SearchAccessors(networkID, keyWord string, start,
	limit int) (data []models.AccessorInfo, count int, err error) {
	err = m.checkNetworkExist(networkID)
	if err != nil {
		return
	}

	if len(keyWord) == 0 {
		// 如果keyWord为空，搜索规则为：
		// 1、根据网段id获取全部绑定了的访问者
		// 2、根据获取访问者的id，调用引擎根据id获取详细信息接口，获取显示名
		// 根据id获取已绑定的访问者
		data, count, err := m.getAccessors(networkID, start, limit)
		if err != nil {
			return nil, 0, err
		}

		for i := range data {
			accessor := &data[i]
			err = m.checkAccessor(accessor)
			if err != nil {
				return nil, 0, err
			}
		}
		return data, count, nil
	} else {
		// 如果keyWord不为空，搜索规则为：
		// 1、调用引擎搜索用户、部门接口，获取服务条件的访问者
		// 2、根据获取的访问者id，查找该网段下符合条件的数据

		// 从引擎所有符合条件的访问者
		sharemgntAccessors, err := m.searchAccessorFromSharemgnt(keyWord, 0, -1)
		if err != nil {
			return nil, 0, err
		}
		var accessorIDs []string
		idAccessorMap := make(map[string]models.AccessorInfo)
		for _, v := range sharemgntAccessors {
			accessorID := v.AccessorId
			accessorIDs = append(accessorIDs, accessorID)
			idAccessorMap[accessorID] = v
		}

		// 根据引擎返回的id，查询符合条件的访问者
		wantIDs, count, err := m.getAccessorsByNetIDAccessorIDs(networkID, accessorIDs, start, limit)
		if err != nil {
			return nil, 0, err
		}
		data := make([]models.AccessorInfo, 0)
		for _, v := range wantIDs {
			data = append(data, idAccessorMap[v])
		}
		return data, count, nil
	}
}

// 状态码为207时的返回内容
func (m *mgnt) appendErrMultiStatus(res *[]*api.MultiStatus, accessor models.AccessorInfo, err error, code int) {
	var st *api.MultiStatus
	if err == nil {
		// err为空，使用传入的状态码
		body := make(map[string]string)
		body["accessor_id"] = accessor.AccessorId
		body["accessor_name"] = accessor.AccessorName
		body["accessor_type"] = accessor.AccessorType
		st = api.MultiStatusObject(accessor.AccessorId, nil, body, code)
	} else {
		// err不为空，状态码为错误码前三位
		if serr, ok := err.(*api.Error); ok {
			st = api.MultiStatusObject(accessor.AccessorId, nil, serr)
		} else {
			st = api.MultiStatusObject(accessor.AccessorId, nil, errors.ErrInternalServerErrorPublic(&api.ErrorInfo{Cause: "Inter error."}))
		}
	}
	*res = append(*res, st)
}

// 添加访问者
func (m *mgnt) AddAccessors(networkID string, accessorInfos []models.AccessorInfo) []*api.MultiStatus {
	// 如果是外网的具体id，403
	if utilities.InStrSlice(networkID, PublicNetDBIDs) {
		mst := api.MultiStatusObject(networkID, nil, errors.ErrNoPermission(&api.ErrorInfo{Cause: "Not allowed to add accessors in this network."}))
		return []*api.MultiStatus{mst}
	}

	// 检查网段是否存在
	if err := m.checkNetworkExist(networkID); err != nil {
		mst := api.MultiStatusObject(networkID, nil, err)
		return []*api.MultiStatus{mst}
	}

	// 去重
	accessors := m.trimDupAccessors(accessorInfos)

	var svAccessors []models.AccessorInfo
	res := make([]*api.MultiStatus, 0)
	for i := range accessors {
		accessor := &accessors[i]
		// 检查访问者是否在引擎中存在
		serr := m.checkAccessor(accessor)
		// 通过检查，添加至待保存访问者列表中
		if serr == nil {
			svAccessors = append(svAccessors, *accessor)
		}
		m.appendErrMultiStatus(&res, *accessor, serr, http.StatusCreated)
	}

	// 忽略已经添加过的访问者
	dbAccessors, err := m.filterAccessors(networkID, svAccessors)
	if err != nil {
		mst := api.MultiStatusObject("", nil, err)
		return []*api.MultiStatus{mst}
	}

	// 写入数据库
	if err := m.addAccessors(networkID, dbAccessors); err != nil {
		mst := api.MultiStatusObject("", nil, err)
		return []*api.MultiStatus{mst}
	}

	return res
}

// 检查访问者是否在关系表中
// 若不存在，返回不存在的访问者id
func (m *mgnt) checkAccessorsExist(networkID string, accessorIDs []string) (notfoundIDs []string, serr error) {
	// 网段id为外网id，取外网实际id的首个
	if networkID == PublicNetID {
		networkID = PublicNetDBIDs[0]
	}
	existedAccessorIDs, err := m.getAddedAccessorIDs(networkID)
	if err != nil {
		return
	}

	notfoundIDs = utilities.Difference(existedAccessorIDs, accessorIDs)
	return
}

// 删除数据库中访问者
func (m *mgnt) deleteAccessors(networkID string, accessorIDs []string) error {
	db, err := api.ConnectDB()
	if err != nil {
		return err
	}

	delete := func(db *gorm.DB, networkID string, accessorIDs []string) error {
		if err := db.Where("f_network_id = ? AND f_accessor_id in (?)", networkID, accessorIDs).
			Delete(&[]models.NetworkAccessorRelation{}).Error; err != nil {
			return err
		}
		return nil
	}

	if networkID == PublicNetID {
		for _, v := range PublicNetDBIDs {
			if err := delete(db, v, accessorIDs); err != nil {
				return err
			}
		}
	} else {
		if err := delete(db, networkID, accessorIDs); err != nil {
			return err
		}
	}

	return nil
}

// 删除访问者
func (m *mgnt) DeleteAccessors(networkID string, accessorIDs []string) []*api.MultiStatus {
	// 如果是外网的具体id，403
	if utilities.InStrSlice(networkID, PublicNetDBIDs) {
		mst := api.MultiStatusObject(networkID, nil, errors.ErrNoPermission(&api.ErrorInfo{Cause: "Not allowed to add accessors in this network."}))
		return []*api.MultiStatus{mst}
	}

	if err := m.checkNetworkExist(networkID); err != nil {
		mst := api.MultiStatusObject(networkID, nil, err)
		return []*api.MultiStatus{mst}
	}

	notfoundIDs, err := m.checkAccessorsExist(networkID, accessorIDs)
	if err != nil {
		mst := api.MultiStatusObject("", nil, err)
		return []*api.MultiStatus{mst}
	}
	res := make([]*api.MultiStatus, 0)
	for _, v := range notfoundIDs {
		err404 := errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"accessor_id"}}})
		accessor := models.AccessorInfo{
			AccessorId: v,
		}
		m.appendErrMultiStatus(&res, accessor, err404, 0)
	}

	toDelIDS := utilities.Difference(notfoundIDs, accessorIDs)
	for _, v := range toDelIDS {
		accessor := models.AccessorInfo{
			AccessorId: v,
		}
		m.appendErrMultiStatus(&res, accessor, nil, http.StatusOK)
	}

	if err := m.deleteAccessors(networkID, toDelIDS); err != nil {
		mst := api.MultiStatusObject("", nil, err)
		return []*api.MultiStatus{mst}
	}

	return res
}

// 根据类型获取访问者网段
func (m *mgnt) getNetsByType(accessorType string) (data []models.NetworkAccessorRelation, err error) {
	db, err := api.ConnectDB()
	if err != nil {
		return
	}

	err = db.Select("f_network_id, f_accessor_id").Where("f_accessor_type = ?", accessorType).Find(&data).Error
	if err != nil {
		return
	}

	return
}

// 返回3个map
// {用户id：[绑定的网段列表]}
// {部门id：[绑定的网段列表]}
// {用户id：[所在部门的列表]}
func (m *mgnt) getAllPolicies() (map[string]map[string]struct{}, map[string]map[string]struct{}, map[string]map[string]struct{}, error) {
	// 获取用户的访问者网段
	userRelations, err := m.getNetsByType("user")
	if err != nil {
		return nil, nil, nil, err
	}
	// 访问者信息表
	// key：用户id
	// value：访问者网段id列表
	useridNetsMap := make(map[string]map[string]struct{})
	for _, r := range userRelations {
		// 如果key不存在，新建
		userid := r.AccessorId
		netid := r.NetworkId
		if _, ok := useridNetsMap[userid]; !ok {
			netids := make(map[string]struct{})
			netids[netid] = struct{}{}
			useridNetsMap[userid] = netids
		} else { // 如果key已存在，增加
			useridNetsMap[userid][netid] = struct{}{}
		}
	}

	// 获取部门、组织的访问者网段
	departmentRelations, err := m.getNetsByType("department")
	if err != nil {
		return nil, nil, nil, err
	}
	// 访问者信息表
	// key：部门id
	// value：访问者网段id列表
	departmentidNetsMap := make(map[string]map[string]struct{})
	// 缓存部门和部门下所有用户id信息
	// key: 部门id
	// value： 部门下用户的id构成的map
	departmentUseridsMap := make(map[string]map[string]struct{})
	for _, r := range departmentRelations {
		// 如果key不存在，新建
		departmentid := r.AccessorId
		netid := r.NetworkId
		if _, ok := departmentidNetsMap[departmentid]; !ok {
			netids := make(map[string]struct{})
			netids[netid] = struct{}{}
			departmentidNetsMap[departmentid] = netids
		} else { // 如果key已存在，增加
			departmentidNetsMap[departmentid][netid] = struct{}{}
		}

		// 如果没有查找过，查找部门下的所有用户id
		if _, ok := departmentUseridsMap[departmentid]; !ok {
			users := make(map[string]struct{})
			if err := m.getUseridsByDepartIDs([]string{departmentid}, users); err != nil {
				return nil, nil, nil, err
			}
			departmentUseridsMap[departmentid] = users
		}
	}

	// 缓存部门和部门下所有用户id信息
	// key: 用户id
	// value： 用户所在部门的id信息
	useridDepartmentidsMap := make(map[string]map[string]struct{})
	for departmentid, userids := range departmentUseridsMap {
		for userid := range userids {
			// 如果key不存在，新建
			if _, ok := useridDepartmentidsMap[userid]; !ok {
				departmentids := make(map[string]struct{})
				departmentids[departmentid] = struct{}{}
				useridDepartmentidsMap[userid] = departmentids
			} else { // 如果key已存在，增加
				useridDepartmentidsMap[userid][departmentid] = struct{}{}
			}
		}
	}

	return useridNetsMap, departmentidNetsMap, useridDepartmentidsMap, nil
}

// 返回map{网段id:[start_ip, end_ip]}
func (m *mgnt) getNetidSegMap() (map[string]models.NetSegment, map[string]models.NetSegment, error) {
	db, err := api.ConnectDB()
	if err != nil {
		return nil, nil, err
	}

	var networks []models.NetworkRestriction
	err = db.Select("f_id, f_segment_start, f_segment_end, f_ip_type").Find(&networks).Error
	if err != nil {
		return nil, nil, err
	}

	// 示例：map[string]models.NetSegment{"net_01_id":models.NetSegment{StartIP:0, EndIP:167772159}}
	netIDSegmentIpv4Map := make(map[string]models.NetSegment)
	netIDSegmentIpv6Map := make(map[string]models.NetSegment)
	for _, n := range networks {
		if n.IpType == "ipv4" {
			startIP, _ := new(big.Int).SetString(n.SegmentStart, 10)
			endIP, _ := new(big.Int).SetString(n.SegmentEnd, 10)
			netIDSegmentIpv4Map[n.ID] = models.NetSegment{
				StartIP: startIP,
				EndIP:   endIP,
			}
		} else {
			startIP, _ := new(big.Int).SetString(n.SegmentStart, 10)
			endIP, _ := new(big.Int).SetString(n.SegmentEnd, 10)
			netIDSegmentIpv6Map[n.ID] = models.NetSegment{
				StartIP: startIP,
				EndIP:   endIP,
			}
		}
	}

	return netIDSegmentIpv4Map, netIDSegmentIpv6Map, nil
}

// 格式化需要传出的部门数据
// 示例如下
// "department01": [
//
//	{
//	    "start_ip": 123,
//	    "end_ip": 223
//	}
//
// ],
// "department02": [
//
//	{
//	    "start_ip": 123,
//	    "end_ip": 223
//	},
//	{
//	    "start_ip": 55,
//	    "end_ip": 66
//	}
//
// ]
func (m *mgnt) convertDepartmentData(departmentidNetsMap map[string]map[string]struct{},
	netIDSegmentMap map[string]models.NetSegment) (map[string][]models.NetSegment, error) {
	res := make(map[string][]models.NetSegment)
	for userID, netIDs := range departmentidNetsMap {
		var netSegments []models.NetSegment
		for netID := range netIDs {
			if netIDSegmentMap[netID].StartIP != nil && netIDSegmentMap[netID].EndIP != nil {
				netSegments = append(netSegments, netIDSegmentMap[netID])
			}
		}
		res[userID] = netSegments
	}

	return res, nil
}

// 格式化需要传出的用户数据
// 示例如下
//
//	"user03": {
//	    "nets": [
//	        {
//	            "start_ip": 32,
//	            "end_ip": 42
//	        }
//	    ]
//	},
//
//	"user02": {
//	    "departments": ["department01", "department02"]
//	},
//
//	"user01": {
//	    "departments": ["department01", "department02"],
//	    "nets": [
//	        {
//	            "start_ip": 123,
//	            "end_ip": 223
//	        },
//	        {
//	            "start_ip": 3,
//	            "end_ip": 4
//	        }
//	    ]
//	}
func (m *mgnt) convertUserData(useridNetsMap map[string]map[string]struct{}, useridDepartmentidsMap map[string]map[string]struct{},
	netIDSegmentMap map[string]models.NetSegment) (map[string]*models.UserNetworkInfo, error) {
	res := make(map[string]*models.UserNetworkInfo)
	// 转换用户所在的部门
	for userID, departmentMap := range useridDepartmentidsMap {
		var departmentIDs []string
		for departmentID := range departmentMap {
			departmentIDs = append(departmentIDs, departmentID)
		}

		var userNetInfo models.UserNetworkInfo
		userNetInfo.Departments = departmentIDs
		res[userID] = &userNetInfo
	}

	// 转换用户绑定的网段
	for userID, netIDs := range useridNetsMap {
		var netSegments []models.NetSegment
		for netID := range netIDs {
			if netIDSegmentMap[netID].StartIP != nil && netIDSegmentMap[netID].EndIP != nil {
				netSegments = append(netSegments, netIDSegmentMap[netID])
			}
		}

		if v, ok := res[userID]; ok { // 如果已存在部门信息，只需要给网段信息赋值
			v.NetSegements = make([]models.NetSegment, 0)
			v.NetSegements = netSegments
			fmt.Println()
		} else { // 不存在部门信息，新建
			var userNetInfo models.UserNetworkInfo
			userNetInfo.NetSegements = netSegments
			res[userID] = &userNetInfo
		}
	}
	return res, nil
}

// 根据部门id，获取用户ID
func (m *mgnt) getUseridsByDepartIDs(departmentIDs []string, users map[string]struct{}) (err error) {
	for _, departmentID := range departmentIDs {
		// 获取部门、组织下的所有用户ID
		getUserids, err := m.abstractDriven.GetDepartUserIds(departmentID)
		if err != nil {
			return err
		}
		for _, i := range getUserids {
			users[i] = struct{}{}
		}
	}
	return
}

// OPA需要的策略信息
// 示例如下
//
//	{
//		"users": {
//			"user01": {
//				"departments": ["department01", "department02"],
//				"nets": [
//					{
//						"start_ip": 123,
//						"end_ip": 223
//					},
//					{
//						"start_ip": 3,
//						"end_ip": 4
//					}
//				]
//			}
//		},
//		"departments": {
//			"department01": [
//				{
//					"start_ip": 523,
//					"end_ip": 623
//				}
//			],
//			"department02": [
//				{
//					"start_ip": 55,
//					"end_ip": 66
//				}
//			]
//		}
//	}
func (m *mgnt) GetNetworkData(result map[string]interface{}) error {
	// 获取{用户id：网段id列表}的map
	// 比如：data =  map[string][]string{"user_01_id":[]string{"net_01_id", "net_02_id"}}
	ipv4Result := make(map[string]interface{})
	ipv6Result := make(map[string]interface{})
	useridNetsMap, departmentidNetsMap, useridDepartmentidsMap, err := m.getAllPolicies()
	if err != nil {
		return err
	}

	// 获取网段id和具体网段区间的map
	netIDSegmentIpv4Map, netIDSegmentIpv6Map, err := m.getNetidSegMap()
	if err != nil {
		return err
	}

	if len(netIDSegmentIpv4Map) != 0 {
		// 转换用户信息
		userIpv4Res, err := m.convertUserData(useridNetsMap, useridDepartmentidsMap, netIDSegmentIpv4Map)
		if err != nil {
			return err
		}

		// 转换部门信息
		departmentIpv4Res, err := m.convertDepartmentData(departmentidNetsMap, netIDSegmentIpv4Map)
		if err != nil {
			return err
		}
		ipv4Result["users"] = userIpv4Res
		ipv4Result["departments"] = departmentIpv4Res
		result["ipv4"] = ipv4Result
	}

	if len(netIDSegmentIpv6Map) != 0 {
		// 转换用户信息
		userIpv6Res, err := m.convertUserData(useridNetsMap, useridDepartmentidsMap, netIDSegmentIpv6Map)
		if err != nil {
			return err
		}

		// 转换部门信息
		departmentIpv6Res, err := m.convertDepartmentData(departmentidNetsMap, netIDSegmentIpv6Map)
		if err != nil {
			return err
		}
		ipv6Result["users"] = userIpv6Res
		ipv6Result["departments"] = departmentIpv6Res
		result["ipv6"] = ipv6Result
	}

	return nil
}

// 清理引擎已经删除的用户、部门
func (m *mgnt) FlushAccessors(accessorIDs []string) (err error) {
	db, err := api.ConnectDB()
	if err != nil {
		return
	}
	err = db.Where("f_accessor_id in (?) ", accessorIDs).Delete(&[]models.NetworkAccessorRelation{}).Error
	if err != nil {
		return
	}
	return
}
