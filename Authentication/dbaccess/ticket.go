package dbaccess

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
)

type ticket struct {
	dbTrace     *sqlx.DB
	batchNumber int // batchNumber 一次清理失效记录的数量
	logger      common.Logger
	trace       observable.Tracer
}

var (
	tOnce sync.Once
	t     *ticket
)

// NewTicket 创建ticket对象
func NewTicket() *ticket {
	tOnce.Do(func() {
		t = &ticket{
			dbTrace:     dbTracePool,
			batchNumber: 100,
			logger:      common.NewLogger(),
			trace:       common.SvcARTrace,
		}
	})
	return t
}

// Create 创建新的ticket
func (t *ticket) Create(ctx context.Context, info *interfaces.TicketInfo) (err error) {
	t.trace.SetClientSpanName("数据访问层-生成新的单点登录凭据")
	newCtx, span := t.trace.AddClientTrace(ctx)
	defer func() { t.trace.TelemetrySpanEnd(span, err) }()

	sqlStr := "insert into t_ticket(`f_id`, `f_user_id`, `f_client_id`, `f_create_time`) values(?, ?, ?, ?)"
	if _, err = t.dbTrace.ExecContext(newCtx, sqlStr, info.ID, info.UserID, info.ClientID, info.CreateTime); err != nil {
		return err
	}
	return nil
}

// GetTicketByID 根据ID获取ticket信息
func (t *ticket) GetTicketByID(ctx context.Context, id string) (info *interfaces.TicketInfo, err error) {
	t.trace.SetClientSpanName("数据访问层-根据ID获取单点登录凭据信息")
	newCtx, span := t.trace.AddClientTrace(ctx)
	defer func() { t.trace.TelemetrySpanEnd(span, err) }()

	info = &interfaces.TicketInfo{}
	sqlStr := "select f_user_id, f_client_id, f_create_time from t_ticket where f_id = ?"
	if err := t.dbTrace.QueryRowContext(newCtx, sqlStr, id).Scan(&info.UserID, &info.ClientID, &info.CreateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, rest.NewHTTPErrorV2(rest.Unauthorized, "invalid ticket")
		}
		return nil, err
	}
	info.ID = id
	return info, nil
}

// GetExpiredRecords 获取已失效的单点登录凭据的ID集合
func (t *ticket) GetExpiredRecords(ctx context.Context, expiration time.Duration) (ids []string, err error) {
	t.trace.SetClientSpanName("数据访问层-获取已失效的单点登录凭据的ID集合")
	newCtx, span := t.trace.AddClientTrace(ctx)
	defer func() { t.trace.TelemetrySpanEnd(span, err) }()

	sqlStr := "select f_id from t_ticket where f_create_time < ? limit ?"
	rows, err := t.dbTrace.QueryContext(newCtx, sqlStr, time.Now().Add(-expiration).Unix(), t.batchNumber)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				t.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				t.logger.Errorln(closeErr)
			}
		}
	}()
	if err != nil {
		t.logger.Errorln(err, sqlStr)
		return nil, err
	}

	ids = make([]string, 0)
	id := ""
	for rows.Next() {
		if scanErr := rows.Scan(&id); scanErr != nil {
			t.logger.Errorln(scanErr, sqlStr)
			return nil, scanErr
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// DeleteByIDs 根据ID批量删除单点登录凭据
func (t *ticket) DeleteByIDs(ctx context.Context, ids []string) (err error) {
	t.trace.SetClientSpanName("数据访问层-根据ID批量删除单点登录凭据")
	newCtx, span := t.trace.AddClientTrace(ctx)
	defer func() { t.trace.TelemetrySpanEnd(span, err) }()

	if len(ids) == 0 {
		return nil
	}

	set, argIDs := GetFindInSetSQL(ids)
	sqlStr := "delete from t_ticket where f_id in (" + set + ")"
	_, err = t.dbTrace.ExecContext(newCtx, sqlStr, argIDs...)
	if err != nil {
		return err
	}

	return nil
}
