package opshttp

import (
	"context"

	recconsts "AuditLog/common/constants/recenums"
	"AuditLog/common/helpers"
)

func (o *opsHttpAcc) CreateInterfaceNoID(ctx context.Context, index string, doc interface{}) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	err = o.opsCmp.CreateInterfaceNoID(ctx, index, doc)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.CreateInterfaceNoID")
		return
	}

	return
}

func (o *opsHttpAcc) BatchCreate(ctx context.Context,
	index string, data []map[string]interface{}, isWithID bool,
) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	err = o.opsCmp.BatchCreate(ctx, index, data, isWithID)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.BatchCreate")
		return
	}

	return
}

func (o *opsHttpAcc) BatchCreateInterface(ctx context.Context, index string, docs interface{}, isWithID bool) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	size := recconsts.CreateDocBatchSize
	if helpers.IsLocalDev() {
		size = 2
	}

	err = o.opsCmp.BatchCreateInterface(ctx, index, docs, isWithID, size)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.BatchCreateInterface")
		return
	}

	return
}

func (o *opsHttpAcc) CreateIndex(ctx context.Context, index string, mapping, setting string) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	b, err := o.opsCmp.IndexExists(ctx, index)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.CreateIndex", "IndexExists")
		return
	}

	if b {
		return
	}

	err = o.opsCmp.CreateIndex(ctx, index, mapping, setting)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.CreateIndex", "CreateIndex")
		return
	}

	return
}

// DeleteIndex 删除索引，不存在则不处理
func (o *opsHttpAcc) DeleteIndex(ctx context.Context, index string) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	b, err := o.opsCmp.IndexExists(ctx, index)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.CreateIndex", "IndexExists")
		return
	}

	if !b {
		return
	}

	err = o.opsCmp.DeleteIndex(ctx, index)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.DeleteIndex")
		return
	}

	return
}

func (o *opsHttpAcc) DeleteDocByField(ctx context.Context, index string, field string, value interface{}) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	err = o.opsCmp.DeleteDocByField(ctx, index, field, value)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.DeleteDocByField")
		return
	}

	return
}

func (o *opsHttpAcc) DeleteDocsByFieldRange(ctx context.Context, index string, field string, from, to interface{}) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	err = o.opsCmp.DeleteDocsByFieldRange(ctx, index, field, from, to)
	if err != nil {
		helpers.RecordErrLogWithPos(o.logger, err, "opsHttpAcc.DeleteDocsByFieldRange")
		return
	}

	return
}
