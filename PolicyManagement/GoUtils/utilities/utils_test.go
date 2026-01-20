package utilities

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestTrimStruct(t *testing.T) {
	type student struct {
		Name     string
		NickName string
		hobby    string
		Age      int
	}

	Convey("忽略非结构体", t, func() {
		s1 := "123123  "
		TrimStruct(&s1)
		assert.Equal(t, s1, s1)
	})

	Convey("去除单层结构体空格", t, func() {
		s1 := student{"bob ", " little bob ", " play game", 10}
		want := student{"bob", "little bob", "play game", 10}
		TrimStruct(&s1)
		assert.Equal(t, want, s1)
	})
}

func TestRemoveDuplicateStruct(t *testing.T) {
	type student struct {
		Name     string
		NickName string
		hobby    string
		Age      int
	}

	Convey("去除成功", t, func() {
		s1 := student{"bob", "little bob", "play game", 10}
		s2 := student{"bob", "little bob", "play game", 10}
		s3 := student{"sss", "little bob", "play game", 10}
		s := []interface{}{s1, s2, s3}
		res := RemoveDuplicateStruct(s)
		want := []interface{}{s1, s3}
		assert.Equal(t, want, res)
	})
}

func TestTrimDupStr(t *testing.T) {
	dupValues := []string{"v1", "v2", "v3", "v1", "v3", "v4"}
	result := []string{"v1", "v2", "v3", "v4"}
	assert.Equal(t, TrimDupStr(dupValues), result)
}

func TestInStrSlice(t *testing.T) {
	tests := []struct {
		target      string
		targetArray []string
		want        bool
	}{
		{"a", []string{"a", "b", "a"}, true},
		{"ss", []string{"a a", "b", "a a"}, false},
	}

	for _, test := range tests {
		res := InStrSlice(test.target, test.targetArray)
		assert.Equal(t, test.want, res)
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		array1 []string
		array2 []string
		want   []string
	}{
		{[]string{"a"}, []string{"b"}, []string{}},
		{[]string{"a"}, []string{"b", "a"}, []string{"a"}},
		{[]string{"b", "a"}, []string{"b", "a"}, []string{"b", "a"}},
	}

	for _, test := range tests {
		res := Intersection(test.array1, test.array2)
		assert.Equal(t, test.want, res)
	}
}

func TestEqualSlice(t *testing.T) {
	tests := []struct {
		array1 []string
		array2 []string
		want   bool
	}{
		{[]string{"a"}, []string{"b"}, false},
		{[]string{"a"}, []string{}, false},
		{[]string{"a"}, []string{"a", "b"}, false},
		{[]string{"b", "a"}, []string{"a", "b"}, false},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
	}

	for _, test := range tests {
		res := EqualSlice(test.array1, test.array2)
		assert.Equal(t, test.want, res)
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		array1 []string
		array2 []string
		want   []string
	}{
		{[]string{"a"}, []string{"b"}, []string{"b"}},
		{[]string{"a"}, []string{"a", "b"}, []string{"b"}},
		{[]string{"a", "b"}, []string{"a"}, []string{}},
		{[]string{"a", "b"}, []string{}, []string{}},
	}

	for _, test := range tests {
		res := Difference(test.array1, test.array2)
		assert.Equal(t, test.want, res)
	}
}

func TestConvertIPAoti(t *testing.T) {
	Convey("成功转换", t, func() {
		tests := []struct {
			input string
			want  int64
		}{
			{"10.2.20.30", 167908382},
			{"0.0.0.0", 0},
			{"255.255.255.255", 4294967295},
		}

		for _, test := range tests {
			res, err := ConvertIPAoti(test.input)
			assert.Equal(t, nil, err)
			assert.Equal(t, test.want, res)
		}
	})

	Convey("转换失败", t, func() {
		res, err := ConvertIPAoti("hehe")
		assert.Equal(t, int64(0), res)
		_, ok := err.(*strconv.NumError)
		assert.Equal(t, true, ok)
	})
}

func TestGetStringSubnet(t *testing.T) {
	var tests = []struct {
		input int
		excp  string
	}{
		{2, "192.0.0.0"},
		{22, "255.255.252.0"},
		{24, "255.255.255.0"},
		{30, "255.255.255.252"},
	}

	for _, test := range tests {
		assert.Equal(t, test.excp, GetStringSubnet(test.input, "%d", "."))
	}
}
func TestSubnetMap(t *testing.T) {
	maskMap := SubnetMap()
	// check 5 k,v
	assert.Equal(t, 0, maskMap["0.0.0.0"])
	assert.Equal(t, 5, maskMap["248.0.0.0"])
	assert.Equal(t, 10, maskMap["255.192.0.0"])
	assert.Equal(t, 23, maskMap["255.255.254.0"])
	assert.Equal(t, 32, maskMap["255.255.255.255"])
}

func TestConvertNetToRange(t *testing.T) {
	tests := []struct {
		inputip   string
		inputmask string
		wantRange []string
	}{
		{"192.168.0.1", "255.255.255.0", []string{"192.168.0.0", "192.168.0.255"}},
		{"266.168.0.1", "255.255.255.0", []string{"266.168.0.0", "266.168.0.255"}},
	}

	for _, test := range tests {
		res := ConvertNetToRange(test.inputip, test.inputmask)
		assert.Equal(t, test.wantRange, res)
	}
}

func TestTarZstd(t *testing.T) {
	Convey("打包Zstd格式文件", t, func() {
		f1, _ := ioutil.TempFile("", "f1")
		f1Path := f1.Name()
		f1Content := []byte("f1 content")
		f1.Write(f1Content)
		defer os.Remove(f1Path)
		tempDir, _ := ioutil.TempDir("", "")
		tempPath := path.Join(tempDir, "test.tar.zst")
		defer os.RemoveAll(tempDir)
		_, err := os.Stat(tempPath)
		_, ok := err.(*os.PathError)
		assert.Equal(t, true, ok) // file not exist

		err = TarZstd(tempPath, f1Path, 3)
		assert.Equal(t, nil, err)

		_, err = os.Stat(tempPath)
		assert.Equal(t, nil, err) // file exist
	})
}
func TestTarGz(t *testing.T) {
	Convey("打包多个文件", t, func() {
		f1, _ := ioutil.TempFile("", "f1")
		f1Path := f1.Name()
		f1Content := []byte("f1 content")
		f1.Write(f1Content)
		defer os.Remove(f1Path)

		f2, _ := ioutil.TempFile("", "f2")
		f2Content := []byte("f2 content")
		f2Path := f2.Name()
		f2.Write(f2Content)
		defer os.Remove(f2Path)

		tempDir, _ := ioutil.TempDir("", "")
		tempPath := path.Join(tempDir, "test.tar.gz")
		defer os.RemoveAll(tempDir)

		_, err := os.Stat(tempPath)
		_, ok := err.(*os.PathError)
		assert.Equal(t, true, ok) // file not exist

		err = TarGz([]string{f1Path, f2Path}, tempPath)
		assert.Equal(t, nil, err)

		_, err = os.Stat(tempPath)
		assert.Equal(t, nil, err) // file exist
	})

	Convey("打包多个文件和文件夹", t, func() {
		f1, _ := ioutil.TempFile("", "f1")
		f1Path := f1.Name()
		f1Content := []byte("f1 content")
		f1.Write(f1Content)
		defer os.Remove(f1Path)

		f2, _ := ioutil.TempFile("", "f2")
		f2Content := []byte("f2 content")
		f2Path := f2.Name()
		f2.Write(f2Content)
		defer os.Remove(f2Path)

		tempDir1, _ := ioutil.TempDir("", "")
		f3, _ := os.Create(path.Join(tempDir1, "f3"))
		defer os.Remove(f3.Name())
		defer os.RemoveAll(tempDir1)

		tempDir, _ := ioutil.TempDir("", "")
		tempPath := path.Join(tempDir, "test.tar.gz")
		defer os.RemoveAll(tempDir)

		_, err := os.Stat(tempPath)
		_, ok := err.(*os.PathError)
		assert.Equal(t, true, ok) // file not exist

		err = TarGz([]string{f1Path, f2Path, tempDir1}, tempPath)
		assert.Equal(t, nil, err)

		_, err = os.Stat(tempPath)
		assert.Equal(t, nil, err) // file exist
	})
}

func TestOverWrite(t *testing.T) {
	Convey("文件不存在", t, func() {
		destFile := "./test.txt"
		_, err := os.Stat(destFile)
		_, ok := err.(*os.PathError)
		assert.Equal(t, true, ok) // file not exist

		content := []byte("hello world")
		err = OverWrite(content, destFile)
		assert.Equal(t, nil, err)
		defer os.Remove(destFile)

		_, err = os.Stat(destFile)
		assert.Equal(t, nil, err) // file exist

		content1, _ := ioutil.ReadFile(destFile)
		assert.Equal(t, content, content1)
	})

	Convey("文件已存在,覆盖原内容", t, func() {
		destFile := "./test.txt"
		file, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, 0666)
		assert.Equal(t, nil, err)
		defer file.Close()
		defer os.Remove(destFile)
		_, err = os.Stat(destFile)
		assert.Equal(t, nil, err) // file exist

		file.WriteString("init test.txt")

		content := []byte("hello world")
		err = OverWrite(content, destFile)
		assert.Equal(t, nil, err)

		content1, _ := ioutil.ReadFile(destFile)
		assert.Equal(t, content, content1)
	})
}

func TestValidJson(t *testing.T) {
	accessorsSchema := `{
		"type": "array",
		"items": {
			"type": "object",
			"properties": {
				"accessor_id": {
					"type": "string"
				},
				"accessor_type": {
					"type": "string",
					"enum": ["user", "department"]
				}
			},
			"required": [
				"accessor_id",
				"accessor_type"
			]
		}
	}`

	Convey("传入内容不是json格式", t, func() {
		invalidInput := "[{sss]"
		invalide_params, cause := ValidJson(accessorsSchema, invalidInput)
		assert.Equal(t, []string([]string{"input json data"}), invalide_params)
		assert.Equal(t, "invalid character 's' looking for beginning of object key string", cause)
	})

	Convey("校验内容不是json格式", t, func() {
		invalidSchema := "{sss"
		invalide_params, cause := ValidJson(invalidSchema, "hehe")
		assert.Equal(t, []string([]string{"input json shema"}), invalide_params)
		assert.Equal(t, "invalid character 's' looking for beginning of object key string", cause)
	})

	Convey("传入内容第一个数据缺少指定字段", t, func() {
		input := `[{"accessor_id_11":"user_01","accessor_type":"user"}]`
		invalide_params, cause := ValidJson(accessorsSchema, input)
		assert.Equal(t, []string([]string{"0"}), invalide_params)
		assert.Equal(t, "0: accessor_id is required", cause)
	})

	Convey("传入内容第一个数据指定字段不为指定值", t, func() {
		input := `[{"accessor_id":"user_01","accessor_type":"user_0111"}]`
		invalide_params, cause := ValidJson(accessorsSchema, input)
		assert.Equal(t, []string([]string{"0.accessor_type"}), invalide_params)
		assert.Equal(t, "0.accessor_type: 0.accessor_type must be one of the following: \"user\", \"department\"", cause)
	})

	Convey("传入内容第一个数据有多个错误", t, func() {
		input := `[{"accessor_id":1,"accessor_type":"user_0111"}]`
		invalide_params, cause := ValidJson(accessorsSchema, input)
		assert.Equal(t, []string([]string{"0.accessor_id", "0.accessor_type"}), invalide_params)
		assert.Equal(t,
			"0.accessor_id: Invalid type. Expected: string, given: integer; 0.accessor_type: 0.accessor_type must be one of the following: \"user\", \"department\"",
			cause)
	})

	Convey("校验通过", t, func() {
		input := `[{"accessor_id":"user_01","accessor_type":"user"}]`
		invalide_params, cause := ValidJson(accessorsSchema, input)
		assert.Equal(t, 0, len(invalide_params))
		assert.Equal(t, "", cause)
	})
}

func TestNewMachineID(t *testing.T) {
	// Convey("未设置环境变量", t, func() {
	// 	f := NewMachineID()
	// 	_, err := f()
	// 	assert.Equal(t, errors.New("'POD_IP' environment variable not set"), err)
	// })

	Convey("非法ip", t, func() {
		os.Setenv("POD_IP", "1111")
		f := NewMachineID()
		_, err := f()
		assert.Equal(t, errors.New("invalid IP"), err)
		os.Unsetenv("POD_IP")
	})

	//var id1,id2 uint64
	var ipStr = "10.2.64.185"
	Convey("环境变量ip正常", t, func() {
		os.Setenv("POD_IP", ipStr)
		f := NewMachineID()
		id1, err := f()
		assert.Equal(t, nil, err)
		assert.NotEqual(t, 0, id1)
		os.Unsetenv("POD_IP")
	})

	// Convey("传入ip正常", t, func() {
	// 	f := NewMachineID()
	// 	id2, err := f()
	// 	assert.Equal(t, nil, err)
	// 	assert.NotEqual(t, 0, id2)
	// })
	// assert.Equal(t, id1, nil)
}

func TestGetUniqueID(t *testing.T) {
	Convey("环境变量MY_IP正常", t, func() {
		os.Setenv("POD_IP", "10.2.64.185")
		var last uint64
		for i := 0; i < 5; i++ {
			id, err := GetUniqueID()
			assert.Equal(t, nil, err)
			assert.Equal(t, true, (id-last) > 0)
			last = id
		}
		os.Unsetenv("POD_IP")
	})

	// Convey("传入变量ip正常", t, func() {
	// 	var last uint64
	// 	for i := 0; i < 5; i++ {
	// 		id, err := GetUniqueID("10.2.64.185")
	// 		assert.Equal(t, nil, err)
	// 		assert.Equal(t, true, (id-last) > 0)
	// 		last = id
	// 	}
	// })
}
