package infra

import (
	"errors"
	"net"
	"os"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

// https://github.com/tinrab/makaroni/tree/master/utilities/unique-id
// 根据ip获取唯一id
func NewMachineID() func() (uint16, error) {
	return func() (uint16, error) {
		ipStr := os.Getenv("POD_IP")
		ip := net.ParseIP(ipStr)
		ip = ip.To16()

		if len(ip) < 4 {
			return 0, errors.New("invalid IP")
		}

		return uint16(ip[14])<<8 + uint16(ip[15]), nil
	}
}

// 使用sonyflake获取唯一、自增id
// 传入ip，使用传入的ip作为机器码
// 不传入ip，使用ipv4作为机器码
func GetUniqueID() (uint64, error) {
	return sf.NextID()
}

// 初始化sonyflake
func init() {
	var st sonyflake.Settings
	// st.StartTime = time.Now()
	st.MachineID = NewMachineID()

	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}
