package drivenadapters

import "UserManagement/interfaces"

var (
	ErrCodeTypeToStr = map[interfaces.ErrorCodeType]string{
		interfaces.Number: "number",
		interfaces.Str:    "string",
	}
)
