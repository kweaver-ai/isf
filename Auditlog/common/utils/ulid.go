package utils

import "github.com/oklog/ulid/v2"

// UlidMake 生成ulid
func UlidMake() string {
	return ulid.Make().String()
}
