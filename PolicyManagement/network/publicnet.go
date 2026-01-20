package network

import (
	"fmt"
	"strings"

	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/kweaver-ai/GoUtils/utilities"

	"gorm.io/gorm"
)

const (
	// 展示给前端的外网id
	PublicNetID = "public-net"
	// 数据库中公网网段id
	PublicNetID1 = "1114e5f4-efb3-49ab-8f4f-df846900e6a2"
	PublicNetID2 = "2114e5f4-efb3-49ab-8f4f-df846900e6a2"
	PublicNetID3 = "3114e5f4-efb3-49ab-8f4f-df846900e6a2"
	PublicNetID4 = "4114e5f4-efb3-49ab-8f4f-df846900e6a2"
)

var (
	PublicNetDBIDs = []string{PublicNetID1, PublicNetID2, PublicNetID3, PublicNetID4}
)

// NetRange 表示一个网段
type NetRange [2]int64

// 获取PublicNetMap
// 排除以下私有地址
// 10.0.0.0--10.255.255.255、
// 172.16.0.0--172.31.255.255、
// 192.168.0.0--192.168.255.255
func getPublicNetMap() map[string]NetRange {
	publicNetStringMap := map[string][]string{
		PublicNetID1: []string{"0.0.0.0", "9.255.255.255"},
		PublicNetID2: []string{"11.0.0.0", "172.15.255.255"},
		PublicNetID3: []string{"172.32.0.0", "192.167.255.255"},
		PublicNetID4: []string{"192.169.0.0", "255.255.255.255"},
	}
	publicNetMap := make(map[string]NetRange, 4)
	for i, v := range publicNetStringMap {
		start, _ := utilities.ConvertIPAoti(v[0])
		end, _ := utilities.ConvertIPAoti(v[1])
		publicNetMap[i] = NetRange{start, end}
	}
	return publicNetMap
}

// batchInsertNets 自行构造批量插入的语句
func batchInsertNets(db *gorm.DB, nets map[string]NetRange) error {
	if len(nets) == 0 {
		return nil
	}
	// 存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(nets))
	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(nets)*3)
	// 遍历nets准备相关数据
	for k, v := range nets {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, `(?, '', '', '', '', ?, ?, '', now())`)
		valueArgs = append(valueArgs, k, v[0], v[1])
	}
	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO t_network_restriction (f_id, f_start_ip, f_end_ip, "+
		"f_ip_address, f_ip_mask, f_segment_start, f_segment_end, f_type, f_created_at) VALUES %s",
		strings.Join(valueStrings, ","))
	if err := db.Exec(stmt, valueArgs...).Error; err != nil {
		return err
	}
	return nil
}

// CreatePublicNet add publicnet to db if not exists
func CreatePublicNet() error {
	db, err := api.ConnectDB()
	if err != nil {
		return err
	}

	unexistedPublicNetMap := make(map[string]NetRange)
	for k, v := range getPublicNetMap() {
		var count int64
		if err := db.Model(&models.NetworkRestriction{}).Where("f_id = ?", k).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			unexistedPublicNetMap[k] = v
		}
	}

	if err := batchInsertNets(db, unexistedPublicNetMap); err != nil {
		return err
	}

	return nil
}
