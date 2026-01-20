package netip

import (
	"math/big"
	"net"
)

// IPNetToIPRange 将net.IPNet转换成起止net.IP, 支持IPV4和IPV6
func IPNetToIPRange(ipnet *net.IPNet) (net.IP, net.IP) {
	endIP := make(net.IP, len(ipnet.IP))
	for i := 0; i < len(ipnet.IP); i++ {
		endIP[i] = ipnet.IP[i] | ^ipnet.Mask[i]
	}
	return ipnet.IP, endIP
}

// IpToInt 将net.IP转换成*big.Int, 支持IPV4和IPV6
// FIXME 当前使用big.Int来表示能正常表示出来，但是实际上应为无符号128位整数，未来标准库有所支持或将考虑修改返回值类型
// TODO type uint128 struct { hi, lo uint64 }
func IpToInt(ip net.IP) *big.Int {
	return big.NewInt(0).SetBytes(ip)
}
