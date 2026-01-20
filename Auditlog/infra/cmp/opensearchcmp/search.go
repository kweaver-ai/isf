package opensearchcmp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go/opensearchapi"

	"AuditLog/models"
)

// ToJSON 将接口转换为JSON字符串
func ToJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("marshal to json failed: %v", err)
		return "{}"
	}

	return string(data)
}

// 聚合查询
func (o *OpsCmp) Query(ctx context.Context, dslQuery, index string) (responseBody *models.OSResp, err error) {
	// 创建查询请求
	req := opensearchapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(dslQuery),
	}

	// 执行查询
	var result *opensearchapi.Response

	for {
		result, err = req.Do(ctx, o.client)
		if err != nil {
			o.logger.Warnf("failed to execute search request: %v", err)
			time.Sleep(1 * time.Second)

			continue
		}

		break
	}

	defer result.Body.Close()

	// 检查响应状态
	if result.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(result.Body).Decode(&e); err != nil {
			o.logger.Warnf("failed to decode response body: %v", err)
			return
		}

		return
	}

	// 解析响应
	responseBody = &models.OSResp{}
	if err := json.NewDecoder(result.Body).Decode(responseBody); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	return
}
