package utilities

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/datadog/zstd"

	"github.com/brotherpowers/ipsubnet"
	"github.com/sony/sonyflake"
	"github.com/xeipuuv/gojsonschema"
)

var (
	sf *sonyflake.Sonyflake
)

type zstdEncoder struct {
	File       *os.File
	ZstdWriter *zstd.Writer
}

// 去除struct各个field的前后空格
// 暂时只支持单层struct
func TrimStruct(input interface{}) {
	rv := reflect.ValueOf(input).Elem()
	// 忽略非结构体指针
	if rv.Kind() != reflect.Struct {
		return
	}
	n := rv.NumField()
	for i := 0; i < n; i++ {
		fv := rv.Field(i)
		fKind := fv.Kind()
		if fKind != reflect.String {
			continue
		}
		// 访问私有field
		fv = reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem()
		newStr := strings.TrimSpace(fv.String())
		fv.SetString(newStr)
	}
}

// deduplicate slice of struct
func RemoveDuplicateStruct(a []interface{}) (ret []interface{}) {
	n := len(a)
	for i := 0; i < n; i++ {
		state := false
		for j := i + 1; j < n; j++ {
			if j > 0 && reflect.DeepEqual(a[i], a[j]) {
				state = true
				break
			}
		}
		if !state {
			ret = append(ret, a[i])
		}
	}
	return
}

// TrimDupStr trim duplicate string from slice
func TrimDupStr(values []string) []string {
	result := make([]string, 0, len(values))
	temp := map[string]struct{}{}
	for _, item := range values {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// InStrSlice check string in slice
func InStrSlice(value string, slice []string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Intersection 取字符串列表交集，
// 交集结果排序以`b`为准
// 算法来源：https://github.com/juliangruber/go-intersect
func Intersection(a, b []string) []string {
	flag := struct{}{}
	set := make([]string, 0)
	hash := make(map[string]struct{})

	for _, ai := range a {
		hash[ai] = flag
	}
	for _, bi := range b {
		if _, ok := hash[bi]; ok {
			set = append(set, bi)
		}
	}

	return set
}

// EqualSlice 比较两个 Slice
// 算法来源：https://stackoverflow.com/a/15312097/2400575
func EqualSlice(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Difference 取字符串列表差集(b - a)，
// Difference(["a"], ["a", "b"]) -> ["b"]
func Difference(a, b []string) []string {
	flag := struct{}{}
	set := make([]string, 0)
	hash := make(map[string]struct{})

	for _, ai := range a {
		hash[ai] = flag
	}
	for _, bi := range b {
		if _, ok := hash[bi]; !ok {
			set = append(set, bi)
		}
	}

	return set
}

// convert ip to int
func ConvertIPAoti(ip string) (intIp int64, err error) {
	var segs []int
	for _, seg := range strings.Split(ip, ".") {
		seg, err := strconv.Atoi(seg)
		if err != nil {
			return 0, err
		}
		segs = append(segs, seg)
	}
	intIp = int64(segs[3] + segs[2]*256 + segs[1]*65536 + segs[0]*16777216)
	return
}

// input:(24,"%d",".") output:"255.255.255.0"
func GetStringSubnet(size int, format, separator string) string {
	subnet_mask := 0xFFFFFFFF << uint(32-size)
	maskQuards := []string{}
	maskQuards = append(maskQuards, fmt.Sprintf(format, (subnet_mask>>24)&0xFF))
	maskQuards = append(maskQuards, fmt.Sprintf(format, (subnet_mask>>16)&0xFF))
	maskQuards = append(maskQuards, fmt.Sprintf(format, (subnet_mask>>8)&0xFF))
	maskQuards = append(maskQuards, fmt.Sprintf(format, (subnet_mask>>0)&0xFF))

	return strings.Join(maskQuards, separator)
}

// all cidr subnet map:
// {"0.0.0.0":0, ..., "255.255.255.255":32}
func SubnetMap() map[string]int {
	subnetMap := make(map[string]int)
	// 子网掩码范围0-32
	for i := 0; i < 33; i++ {
		subnet := GetStringSubnet(i, "%d", ".")
		subnetMap[subnet] = i
	}
	return subnetMap
}

// convert ip + mask to ip range
// input:("192.168.0.1", "255.255.255.0") output:["192.168.0.1", "192.168.0.254"]
func ConvertNetToRange(ipAddress, mask string) (ipRange []string) {
	subnetMap := SubnetMap()
	subInt := subnetMap[mask]
	sub := ipsubnet.SubnetCalculator(ipAddress, subInt)
	ipRange = sub.GetIPAddressRange()
	return
}

//Zstd打包、压缩文件、文件夹，格式为tar.zst
//level是压缩级别：1~19，越大越慢。默认为3

func newZstdEncoder(filename string, level int) (*zstdEncoder, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	w := zstd.NewWriterLevel(f, level)
	return &zstdEncoder{File: f, ZstdWriter: w}, nil
}

func (s *zstdEncoder) close() {

	s.ZstdWriter.Close()
	s.File.Close()
}

func TarZstd(dst, src string, compress_level int) error {
	//log.Println("Zstd ", src, "->", dst)
	s, err := newZstdEncoder(dst, compress_level)
	defer s.close()
	if err != nil {
		return err
	}
	src_file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer src_file.Close()
	_, err = io.Copy(s.ZstdWriter, src_file)
	if err != nil {
		return err
	}

	return nil
}

// 打包、压缩文件、文件夹，格式为tar.gz
func TarGz(srcs []string, dst string) error {
	// 创建文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// 创建 Tar.Writer 结构
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, src := range srcs {
		fileFunc := func(fileName string, fi os.FileInfo, err error) error {
			// 因为这个闭包会返回个 error ，所以先要处理一下这个
			if err != nil {
				return err
			}

			hdr, err := tar.FileInfoHeader(fi, "")
			if err != nil {
				return err
			}
			// strings.TrimPrefix 将 fileName 的最左侧的 / 去掉，
			hdr.Name = strings.TrimPrefix(fileName, string(filepath.Separator))

			// 写入文件信息
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			// 判断下文件是否是标准文件，如果不是就不处理了，
			// 如： 目录，这里就只记录了文件信息，不会执行下面的 copy
			if !fi.Mode().IsRegular() {
				return nil
			}

			// 打开文件
			fr, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer fr.Close()

			// copy 文件数据到 tw
			if _, err := io.Copy(tw, fr); err != nil {
				return err
			}

			return nil
		}

		if err := filepath.Walk(src, fileFunc); err != nil {
			return err
		}
	}

	return nil
}

// 覆盖文件
func OverWrite(data []byte, destPath string) error {
	df, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(df)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	df.Close()

	return nil
}

// 校验json数据，返回错误字段、错误原因
func ValidJson(shema, doc string) (invalide_params []string, cause string) {
	jsonLoader := gojsonschema.NewStringLoader(shema)
	schema, err := gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		return []string{"input json shema"}, err.Error()
	}

	documentLoader := gojsonschema.NewStringLoader(doc)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return []string{"input json data"}, err.Error()
	}

	if result.Valid() {
		return
	}
	for _, desc := range result.Errors() {
		invalide_params = append(invalide_params, desc.Field())
		if cause == "" {
			cause = desc.String()
			continue
		}
		cause = strings.Join([]string{cause, desc.String()}, "; ")
	}
	return
}

// https://github.com/tinrab/makaroni/tree/master/utilities/unique-id
// 根据ip获取唯一id
func NewMachineID() func() (uint16, error) {
	return func() (uint16, error) {
		ipStr := os.Getenv("POD_IP")
		ip := net.ParseIP(ipStr)
		ip = ip.To16()
		if ip == nil || len(ip) < 4 {
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
