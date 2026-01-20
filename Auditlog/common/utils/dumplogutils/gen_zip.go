package dumplogutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/yeka/zip"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// GenZipFile 生成zip文件
func GenZipFile(fileContent []byte, fileName string, pwd string) (zipContent []byte, err error) {
	// 验证输入参数
	if len(fileContent) == 0 {
		return nil, errors.New("[GenZipFile] file content is empty")
	}

	if fileName == "" {
		return nil, errors.New("[GenZipFile] file name is empty")
	}

	// 将文件名转换为GBK编码
	gbkFileName, err := UTF8ToGBK(fileName)
	if err != nil {
		return nil, fmt.Errorf("[GenZipFile] convert filename to GBK error: %w", err)
	}

	// 创建一个预估大小的缓冲区，避免频繁扩容
	estimatedSize := len(fileContent) + 1024
	buf := bytes.NewBuffer(make([]byte, 0, estimatedSize))

	zipWriter := zip.NewWriter(buf)
	defer func() {
		closeErr := zipWriter.Close()
		if closeErr != nil {
			err = closeErr
			return
		}

		zipContent = buf.Bytes()
	}()

	var writer io.Writer
	if pwd != "" {
		writer, err = zipWriter.Encrypt(gbkFileName, pwd, zip.StandardEncryption)
	} else {
		writer, err = zipWriter.Create(gbkFileName)
	}

	if err != nil {
		return nil, fmt.Errorf("[GenZipFile] create zip entry error: %w", err)
	}

	if _, err = writer.Write(fileContent); err != nil {
		return nil, fmt.Errorf("[GenZipFile] write file content error: %w", err)
	}

	return
}

// UTF8ToGBK 将UTF8文本转换为GBK编码
func UTF8ToGBK(text string) (string, error) {
	encoder := simplifiedchinese.GBK.NewEncoder()

	gbkBytes, err := encoder.String(text)
	if err != nil {
		return "", err
	}

	return gbkBytes, nil
}

// SplitFile 将文件分块
func SplitFile(fileContent []byte, size int64) (parts [][]byte, err error) {
	parts = make([][]byte, 0)
	fileLen := len(fileContent)

	for i := 0; i < fileLen; i += int(size) {
		end := i + int(size)
		if end > fileLen {
			end = fileLen
		}

		parts = append(parts, fileContent[i:end])
	}

	return parts, nil
}
