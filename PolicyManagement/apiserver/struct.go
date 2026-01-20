package apiserver

import "encoding/json"

type policyParam struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}
