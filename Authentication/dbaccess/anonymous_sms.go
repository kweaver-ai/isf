package dbaccess

import (
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
)

type anonymousSMS struct {
	dbTrace     *sqlx.DB
	batchNumber int // batchNumber 一次清理失效记录的数量
	desKey      string
	logger      common.Logger
	trace       observable.Tracer
}

var (
	aSMSOnce sync.Once
	aSMS     *anonymousSMS
)

// NewAnonymousSMS 创建anonymousSMS对象
func NewAnonymousSMS() *anonymousSMS {
	aSMSOnce.Do(func() {
		aSMS = &anonymousSMS{
			dbTrace:     dbTracePool,
			batchNumber: 10,
			desKey:      "Ea8ek&ah",
			logger:      common.NewLogger(),
			trace:       common.SvcARTrace,
		}
	})
	return aSMS
}

// Create 创建新的匿名登录验证码
func (aSMS *anonymousSMS) Create(ctx context.Context, aSMSInfo *interfaces.AnonymousSMSInfo) (err error) {
	aSMS.trace.SetClientSpanName("数据访问层-创建新的匿名登录验证码")
	newCtx, span := aSMS.trace.AddClientTrace(ctx)
	defer func() { aSMS.trace.TelemetrySpanEnd(span, err) }()

	phoneNumber, err := aSMS.encryptDES([]byte(aSMSInfo.PhoneNumber), []byte(aSMS.desKey))
	if err != nil {
		return err
	}

	sqlStr := "insert into t_anonymous_sms_vcode(`f_id`, `f_phone_number`, `f_anonymity_id`, `f_content`, `f_create_time`) " +
		"values(?, ?, ?, ?, ?)"
	if _, err = aSMS.dbTrace.ExecContext(newCtx, sqlStr, aSMSInfo.ID, hex.EncodeToString(phoneNumber), aSMSInfo.AnonymityID, aSMSInfo.Content, time.Unix(aSMSInfo.CreateTime, 0)); err != nil {
		return err
	}
	return nil
}

func convertTime(req time.Time) (result time.Time) {
	dbType := os.Getenv("DB_TYPE")
	var err error
	if strings.HasPrefix(dbType, "KDB") {
		// KDB 读取出来是默认UTC时区
		result, err = time.ParseInLocation(time.DateTime, req.Format(time.DateTime), time.Local)
		if err != nil {
			common.NewLogger().Errorf("KDB parse time as DateTime failed, timeStr=%s, err=%v\n", req.UTC().String(), err)
		}
	} else {
		result = req
	}
	return result
}

// GetInfoByID 根据ID获取匿名验证码记录
func (aSMS *anonymousSMS) GetInfoByID(ctx context.Context, id string) (aSMSInfo *interfaces.AnonymousSMSInfo, err error) {
	aSMS.trace.SetClientSpanName("数据访问层-根据ID获取匿名验证码记录")
	newCtx, span := aSMS.trace.AddClientTrace(ctx)
	defer func() { aSMS.trace.TelemetrySpanEnd(span, err) }()

	aSMSInfo = &interfaces.AnonymousSMSInfo{}
	sqlStr := "select f_phone_number, f_anonymity_id, f_content, f_create_time from t_anonymous_sms_vcode where f_id = ?"
	phoneNumber := ""
	var createTime time.Time
	if err1 := aSMS.dbTrace.QueryRowContext(newCtx, sqlStr, id).Scan(&phoneNumber, &aSMSInfo.AnonymityID, &aSMSInfo.Content, &createTime); err1 != nil {
		return nil, err1
	}

	aSMSInfo.ID = id
	aSMSInfo.CreateTime = convertTime(createTime).Unix()
	cipherText, err := hex.DecodeString(phoneNumber)
	if err != nil {
		return nil, err
	}
	plainText, err := aSMS.decryptDES(cipherText, []byte(aSMS.desKey))
	if err != nil {
		return nil, err
	}
	aSMSInfo.PhoneNumber = string(plainText)

	return aSMSInfo, nil
}

// GetExpiredRecords 获取已失效的匿名验证码记录的ID集合
func (aSMS *anonymousSMS) GetExpiredRecords(ctx context.Context, aSMSExpiration time.Duration) (ids []string, err error) {
	aSMS.trace.SetClientSpanName("数据访问层-获取已失效的匿名验证码记录的ID集合")
	newCtx, span := aSMS.trace.AddClientTrace(ctx)
	defer func() { aSMS.trace.TelemetrySpanEnd(span, err) }()

	timeStr := time.Unix(time.Now().Add(-aSMSExpiration).Unix(), 0)
	sqlStr := "select f_id from t_anonymous_sms_vcode where f_create_time < ? limit ?"

	rows, err := aSMS.dbTrace.QueryContext(newCtx, sqlStr, timeStr, aSMS.batchNumber)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				aSMS.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				aSMS.logger.Errorln(closeErr)
			}
		}
	}()
	if err != nil {
		aSMS.logger.Errorln(err, sqlStr, timeStr)
		return nil, err
	}

	ids = make([]string, 0)
	id := ""
	for rows.Next() {
		if scanErr := rows.Scan(&id); scanErr != nil {
			aSMS.logger.Errorln(scanErr, sqlStr)
			return nil, scanErr
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// DeleteByIDs 根据ID删除匿名验证码记录
func (aSMS *anonymousSMS) DeleteByIDs(ctx context.Context, ids []string) (err error) {
	aSMS.trace.SetClientSpanName("数据访问层-根据ID删除匿名验证码记录")
	newCtx, span := aSMS.trace.AddClientTrace(ctx)
	defer func() { aSMS.trace.TelemetrySpanEnd(span, err) }()

	if len(ids) == 0 {
		return nil
	}

	set, argIDs := GetFindInSetSQL(ids)
	sqlStr := "delete from t_anonymous_sms_vcode where f_id in (" + set + ")"

	_, err = aSMS.dbTrace.ExecContext(newCtx, sqlStr, argIDs...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecordWithinValidityPeriod 删除仍在有效期内的匿名验证码
func (aSMS *anonymousSMS) DeleteRecordWithinValidityPeriod(ctx context.Context, phoneNumber, anonymityID string, aSMSExpiration time.Duration) (err error) {
	aSMS.trace.SetClientSpanName("数据访问层-根据ID删除匿名验证码记录")
	newCtx, span := aSMS.trace.AddClientTrace(ctx)
	defer func() { aSMS.trace.TelemetrySpanEnd(span, err) }()

	encryptedData, err := aSMS.encryptDES([]byte(phoneNumber), []byte(aSMS.desKey))
	if err != nil {
		return err
	}

	timeStr := time.Unix(time.Now().Add(-aSMSExpiration).Unix(), 0)
	sqlStr := "delete from t_anonymous_sms_vcode where f_phone_number = ? and f_anonymity_id = ? and f_create_time >= ?"
	_, err = aSMS.dbTrace.ExecContext(newCtx, sqlStr, hex.EncodeToString(encryptedData), anonymityID, timeStr)
	if err != nil {
		return err
	}

	return nil
}

// encryptDES des 加密
func (aSMS *anonymousSMS) encryptDES(plaintext, key []byte) (ciphertext []byte, err error) {
	// 创建一个 DES 加密器
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个加密模式
	mode := cipher.NewCBCEncrypter(block, key)

	// 填充原始数据
	plaintext = pad(plaintext, block.BlockSize())

	// 加密数据
	ciphertext = make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// 填充数据
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// decryptDES des 解密
func (aSMS *anonymousSMS) decryptDES(ciphertext, key []byte) (text []byte, err error) {
	// 创建一个 DES 解密器
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个解密模式
	mode := cipher.NewCBCDecrypter(block, key)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	unpaddedPlaintext := unpad(plaintext)

	return unpaddedPlaintext, nil
}

// 去除填充
func unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
