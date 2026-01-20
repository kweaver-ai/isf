package jsonschema

import (
	_ "embed"
)

var (
	// 访问者网段决策
	//go:embed decision/acccessor_network.json
	AcccessorNetworkSchema string
)
