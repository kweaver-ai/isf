package opensearchcmp

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/opensearch-project/opensearch-go/opensearchapi"

	"AuditLog/common/utils"
)

func (o *OpsCmp) Create(ctx context.Context, index, docID string, docReader io.Reader) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: docID,
		Body:       docReader,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, o.client)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		err = fmt.Errorf("[OpsCmp][Create]:document creation error: %s", res.String())
	}

	return
}

func (o *OpsCmp) CreateInterfaceNoID(ctx context.Context, index string, i interface{}) (err error) {
	ctx, span := o.arTrace.AddInternalTrace(ctx)
	defer func() { o.arTrace.TelemetrySpanEnd(span, err) }()

	bys, err := utils.JSON().Marshal(i)
	if err != nil {
		return
	}

	err = o.Create(ctx, index, "", bytes.NewReader(bys))

	return
}
