package utils

import (
	"net"

	"errors"

	"github.com/sony/sonyflake"
)

// TODO: 优化

var (
	sf *sonyflake.Sonyflake
)

// 使用sonyflake获取唯一、自增id
// 传入ip，使用传入的ip作为机器码
// 不传入ip，使用ipv4作为机器码
// GetSonyflakeID 获取SonyFlakeID
func GetSonyflakeID() (uint64, error) {
	id, err := sf.NextID()
	if err != nil {
		return id, err
	}
	return id, nil
}

// 初始化sonyflake
func InitSonyflake(podIP string) error {
	var st sonyflake.Settings

	// https://github.com/tinrab/makaroni/tree/master/utilities/unique-id
	// 根据ip获取唯一id
	st.MachineID = func() (uint16, error) {
		ip := net.ParseIP(podIP)
		ip = ip.To16()
		if len(ip) < 4 {
			return 0, errors.New("invalid IP")
		}
		return uint16(ip[14])<<8 + uint16(ip[15]), nil
	}

	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		return errors.New("sonyflake initialization failed")
	}

	if _, err := GetSonyflakeID(); err != nil {
		return err
	}
	return nil
}
