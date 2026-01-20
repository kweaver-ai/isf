package uniqidenums

type UniqueIDFlag int

const (
	UniqueIDFlagDB UniqueIDFlag = 1 // UniqueIDFlagDB is the flag for unique id in database

	UniqueIDFlagRedisDlm UniqueIDFlag = 2 // UniqueIDFlagRedisDlm is the flag for unique value in redis dlm
)
