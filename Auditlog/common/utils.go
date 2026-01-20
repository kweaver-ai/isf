package common

import (
	"encoding/json"
	stdErr "errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	"AuditLog/common/helpers"
	"AuditLog/errors"
)

// ParseBody 解析body中的json数据，如果数据类型不符，报错
// 错误信息中包含解析失败的field
func ParseBody(c *gin.Context, params interface{}) (err error) {
	err = c.ShouldBindJSON(&params)
	if err != nil {
		var invalideParams []string
		var field string
		if err1, ok := err.(*json.UnmarshalTypeError); ok {
			field = err1.Field
		}
		invalideParams = append(invalideParams, field)
		err = errors.New("", errors.BadRequestErr, "Invalid Params", map[string]interface{}{"invalid_params": invalideParams})
		return
	}
	return
}

// ErrorResponse 返回错误
func ErrorResponse(c *gin.Context, err error) {
	apiErr, ok := err.(*errors.ErrorResp)
	if !ok {
		apiErr = errors.New(c.GetHeader("x-language"), errors.InternalErr, "Internal Error", "")
	}

	tmpErr := rest.NewHTTPErrorV2(apiErr.Code(), apiErr.Cause())
	tmpErr.Description = apiErr.Description()

	tempDetail, ok := apiErr.Detail().(map[string]interface{})
	if ok {
		tmpErr.Detail = tempDetail
	} else {
		tmpErr.Detail = map[string]interface{}{"detail": apiErr.Detail()}
	}
	tmpErr.Solution = apiErr.Solution()
	tmpErr.Message = apiErr.Description()
	rest.ReplyError(c, tmpErr)
	c.Abort()
}

func ErrResponse(c *gin.Context, err error) {
	httpErr := CustomErrBuild(c, err)

	tmpErr := rest.NewHTTPErrorV2(httpErr.Code(), httpErr.Cause())
	tmpErr.Description = httpErr.Description()

	tempDetail, ok := httpErr.Detail().(map[string]interface{})
	if ok {
		tmpErr.Detail = tempDetail
	} else {
		tmpErr.Detail = map[string]interface{}{"detail": httpErr.Detail()}
	}
	tmpErr.Solution = httpErr.Solution()
	tmpErr.Message = httpErr.Description()
	rest.ReplyError(c, tmpErr)
	c.Abort()
}

func CustomErrBuild(c *gin.Context, err error) (httpErr *errors.ErrorResp) {
	// var httpErr *errors.ErrorResp
	ok := stdErr.As(err, &httpErr)
	if !ok {
		lang := helpers.GetLangFromCtx(c)
		httpErr = errors.New(string(lang), errors.InternalErr, err.Error(), nil)
	}

	return
}

// 字符串数组/切片去重(数据量大建议使用指针的的方式)
func ArrayRemoveDuplicate(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func Configure(obj interface{}, file string) (err error) {
	ioStream, path, err := GetConfigFileContent(file)
	if err != nil {
		fmt.Printf("Configure readfile error: %v\n", err)
		return
	}

	err = yaml.Unmarshal(ioStream, obj)
	ct := time.Now().Format("2006-01-02 15:04:05.000")
	if err != nil {
		fmt.Printf("[ERROR] %s configure decode error: %s\n", ct, err.Error())
		return
	} else {
		fmt.Printf("[INFO] %s configure decode: %s\n", ct, path+"/"+file)
	}

	return
}

// GetConfigFileContent 获取配置文件内容
func GetConfigFileContent(file string) (bys []byte, path string, err error) {
	path = "/sysvol/conf"
	if confPath := os.Getenv("CONFIG_PATH"); confPath != "" {
		path = confPath
	}

	bys, err = os.ReadFile(path + "/" + file)
	if err != nil {
		fmt.Printf("[GetConfigFileContent]: readfile error: %v\n", err)
		return
	}

	return
}

// GetDocLibIDByDocID 获取文档库ID
// isDocLib: docID是否为文档库ID
func GetDocLibIDByDocID(docID string) (docLibID string) {
	// docID 格式为 gns://xxx/ttt/ddd
	// docLibID 的 doclibid 为 gns://xxx
	// 获取gns://xxx
	parts := strings.Split(docID, "://")
	if len(parts) != 2 {
		return
	}

	leftParts := strings.Split(parts[1], "/")

	docLibID = parts[0] + "://" + leftParts[0]

	return
}

// 判断字符串是否在列表中
func InArray(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}

// 定义文件大小单位
var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "BB", "NB", "DB"}

// FormatFileSize 将字节大小转换为可读的格式，可配置小数点位数
func FormatFileSize(size int, precision int) string {
	if size == 0 {
		return "0B"
	}

	sizeF := float64(size)
	index := 0

	// 1024为基数进行换算
	for sizeF >= 1024 && index < len(sizeUnits)-1 {
		sizeF /= 1024
		index++
	}

	// 字节单位不显示小数点，四舍五入到整数
	if index == 0 {
		return fmt.Sprintf("%d%s", int(math.Round(sizeF)), sizeUnits[index])
	}

	// 其他单位根据指定精度保留小数并四舍五入
	factor := math.Pow10(precision)
	return fmt.Sprintf("%.*f%s", precision, math.Round(sizeF*factor)/factor, sizeUnits[index])
}

// GenerateID 生成32位大写无连字符的ID
func GenerateID() string {
	// 生成UUID并转换格式
	return strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
}

// Min 返回两个 int64 中的较小值
func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// SliceDifference 计算两个切片的差集
func SliceDifference(slice1, slice2 []string) []string {
	// 使用 map 记录 slice2 中的元素
	exists := make(map[string]bool)
	for _, item := range slice2 {
		exists[item] = true
	}

	// 遍历 slice1，找出不在 slice2 中的元素
	var diff []string
	for _, item := range slice1 {
		if !exists[item] {
			diff = append(diff, item)
		}
	}

	return diff
}
