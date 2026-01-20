package drivenadapters

import "Authentication/interfaces"

var (
	ErrCodeTypeToStr = map[interfaces.ErrorCodeType]string{
		interfaces.Number: "number",
		interfaces.Str:    "string",
	}
)
